package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

type simulation struct {
	world *World
}

const (
	MOVEMENT      = 3
	POPULATION    = 1000
	MUTATION_RATE = 100 // x in 1000
	GENERATIONS   = 1000
	STEPS_PER_GEN = 250
	SIZE          = 500
	DUMP_EVERY    = 100
)

type act struct {
	peepID  int
	actions Actions
}

func init() {
	seed := time.Now().UnixNano()
	// seed := int64(1637951517777129656)
	fmt.Printf("rand seed: %d\n", seed)
	rand.Seed(seed)
}

func main() {
	world := &World{
		StepsPerGeneration: STEPS_PER_GEN,
		XSize:              SIZE,
		YSize:              SIZE,
		cells:              make([]Cell, SIZE*SIZE),
		survivalArea: Area{
			TopLeft:     Coord{0, 0},
			BottomRight: Coord{100, SIZE},
		},
		barriers: []Area{{
			TopLeft:     Coord{200, 0},
			BottomRight: Coord{210, 100},
		}, {
			TopLeft:     Coord{200, SIZE - 100},
			BottomRight: Coord{210, SIZE},
		}, {
			TopLeft:     Coord{220, 90},
			BottomRight: Coord{240, SIZE - 90},
		}},
	}
	world.fillBarriers()
	fillWithRandomPeeps(world)

	s := &simulation{
		world: world,
	}

	bar := pb.ProgressBarTemplate(`Generation {{counters . }} Survivors: {{string . "survivors"}} {{bar . }} {{percent . }} {{rtime . "ETA %s"}}`).Start(GENERATIONS)
	for generation := 0; generation < GENERATIONS; generation++ {
		bar.Increment()
		for step := 0; step < s.world.StepsPerGeneration; step++ {
			s.step()
			if generation%DUMP_EVERY == 0 {
				produceImage(generation, step, world)
			}
		}

		survivors := cull(world)
		bar.Set("survivors", fmt.Sprintf("%d", len(survivors)))
		if generation%DUMP_EVERY == 0 {
			dumpIndividuals(generation, survivors)
		}

		if len(survivors) == 0 {
			fmt.Println("extinction")
			os.Exit(0)
		}

		copies := POPULATION / len(survivors)

		// fair distribution of survivors
		for _, survivor := range survivors {
			for i := 0; i < copies; i++ {
				clone := survivor.clone()
				clone.location = world.randomCoord()
				clone.birthPlace = clone.location
				world.addPeep(clone)
			}
		}

		// random fill up of peeps until we reach desired population
		for len(world.peeps) < POPULATION {
			var clone *Individual
			if plusMinusOne() > 0 {
				// now and then we'll add a brand-new mutant to the mix, to try to get away from local minimum
				clone = createIndividual(world)
			} else {
				peep := survivors[rand.Intn(len(survivors))]
				clone = peep.clone()
				clone.location = world.randomCoord()
			}
			clone.birthPlace = clone.location
			world.addPeep(clone)
		}
	}
	bar.Finish()
	fmt.Println("done")
}

func dumpIndividuals(generation int, peeps []*Individual) {
	var data []string
	seen := map[string]int{}
	for _, peep := range peeps {
		brain := peep.brain.String() + "\n"
		if idx, ok := seen[brain]; ok {
			data[idx] += "*"
			continue
		}

		seen[brain] = len(data)
		data = append(data, brain)
	}
	output := strings.Join(data, "\n")
	err := os.WriteFile(fmt.Sprintf("%04d/peeps.txt", generation), []byte(output), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func produceImage(generation, step int, world *World) {
	// we copy the cells and write the image on a separate thread
	cells := make([]Cell, len(world.cells))
	copy(cells, world.cells)

	go func() {
		img := image.NewNRGBA(image.Rect(0, 0, world.XSize, world.YSize))
		for x := 0; x < world.XSize; x++ {
			for y := 0; y < world.XSize; y++ {
				offset := y*world.XSize + x
				switch cells[offset] {
				case EMPTY:
					if world.survivalArea.inside(x, y) {
						img.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 0xff})
					} else {
						img.Set(x, y, color.White)
					}
				case BARRIER:
					img.Set(x, y, color.RGBA{R: 200, G: 200, B: 200, A: 0xff})
				default: // here is an individual
					img.Set(x, y, color.Black)
				}
			}
		}
		directory := fmt.Sprintf("%04d", generation)
		err := mkdirIfNotExists(directory)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(fmt.Sprintf("%s/image%03d.png", directory, step))
		if err != nil {
			log.Fatal(err)
		}

		if err := png.Encode(f, img); err != nil {
			_ = f.Close()
			log.Fatal(err)
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
}

func mkdirIfNotExists(directory string) error {
	err := os.Mkdir(directory, os.ModePerm)
	if err != nil {
		pathErr, ok := err.(*fs.PathError)
		if ok {
			sysErr, ok := pathErr.Err.(syscall.Errno)
			if ok {
				if sysErr == syscall.EEXIST {
					return nil
				}
			}
		}
	}
	return err
}

func cull(world *World) []*Individual {
	peeps := world.peeps
	world.clearAll()

	var survivors []*Individual
	for _, peep := range peeps {
		if world.survivalArea.inside(peep.location.X, peep.location.Y) {
			survivors = append(survivors, peep)
		}
	}
	return survivors
}

func fillWithRandomPeeps(world *World) {
	for i := 0; i < POPULATION; i++ {
		individual := createIndividual(world)
		if len(individual.brain.Connections) < 3 {
			i--
			continue
		}
		world.addPeep(individual)
	}
}

func (w *World) randomCoord() Coord {
	location := Coord{
		X: rand.Intn(w.XSize),
		Y: rand.Intn(w.YSize),
	}

	newOffset := w.offset(location)
	if w.cells[newOffset] != EMPTY {
		return w.randomCoord()
	}

	return location
}

// Goes over all individuals and first lets their neural nets run and produce an action slice.
// This is done concurrently, and then the actions are actually performed in a single thread
func (s *simulation) step() {
	peepActions := s.startPeeking()

	for actions := range peepActions {
		individual := s.world.peeps[actions.peepID]
		for act, value := range actions.actions {
			if value != 0 {
				switch Action(act) {
				case MOVE_X:
					loc := individual.location
					loc.X += int(value * MOVEMENT)
					individual.wasBlocked = s.world.updateLocation(actions.peepID, loc)
				case MOVE_Y:
					loc := individual.location
					loc.Y += int(value * MOVEMENT)
					individual.wasBlocked = s.world.updateLocation(actions.peepID, loc)
				case MOVE_RANDOM:
					loc := individual.location
					if plusMinusOne() > 0 {
						loc.X += int(value * MOVEMENT)
					} else {
						loc.Y += int(value * MOVEMENT)
					}
					individual.wasBlocked = s.world.updateLocation(actions.peepID, loc)
				}
			}
		}
	}
}

// runs the neural nets concurrently and produces a channel with their action outputs
func (s *simulation) startPeeking() chan act {
	// we start all the individuals in separate goroutines, and then wait for them to finish
	peepActions := make(chan act, len(s.world.peeps))
	var wg sync.WaitGroup
	for id, peep := range s.world.peeps {
		wg.Add(1)
		go func(peep *Individual, id int) {
			peep.wasBlocked = false
			peepActions <- act{
				peepID:  id,
				actions: peep.step(s.world),
			}
			wg.Done()
		}(peep, id)
	}
	wg.Wait()
	close(peepActions)
	return peepActions
}
