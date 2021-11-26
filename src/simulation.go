package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"math/rand"
	"os"
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
	MUTATION_RATE = 10
	GENERATIONS   = 1000
	STEPS_PER_GEN = 250
	SIZE          = 100
)

type act struct {
	peepID  int
	actions Actions
}

func init() {
	seed := time.Now().UnixNano()
	// nano := int64(1637780343848163000)
	fmt.Printf("rand seed: %d\n", seed)
	rand.Seed(seed)
}

func main() {
	world := &World{
		StepsPerGeneration: STEPS_PER_GEN,
		XSize:              SIZE,
		YSize:              SIZE,
		cells:              make([]Cell, 100*100),
	}
	fillWithRandomPeeps(world)

	s := &simulation{
		world: world,
	}

	for generation := 0; generation < GENERATIONS; generation++{
		for step :=0; step <  s.world.StepsPerGeneration; step++{
			s.step()
			if generation % 100 == 0 {
				// we only write an image every hundred generations
				produceImage(generation, step, world)
			}
		}

		survivors := cull(world)

		if len(survivors) == 0 {
			fmt.Println("extinction")
			os.Exit(0)
		}

		copies := POPULATION / len(survivors)

		// fair distribution of survivors
		for _, survivor := range survivors {
			for i := 0; i < copies; i++ {
				clone := survivor.clone()
				clone.location = randomCoord(world.XSize, world.YSize)
				clone.birthPlace = clone.location
				world.addPeep(clone)
			}
		}

		// random fill up of peeps until we reach desired population
		for len(world.peeps) < POPULATION {
			peep := survivors[rand.Intn(len(survivors))]
			clone := peep.clone()
			clone.location = randomCoord(world.XSize, world.YSize)
			clone.birthPlace = clone.location
			world.addPeep(clone)
		}

		fmt.Printf("%d %d\n", generation, len(survivors))
	}
	fmt.Println("done")
}

func produceImage(generation, step int, world *World) {
	img := image.NewNRGBA(image.Rect(0, 0, world.XSize, world.YSize))
	for x := 0; x < world.XSize; x++ {
		for y := 0; y < world.XSize; y++ {
			offset := y*world.XSize + x
			if world.cells[offset] == EMPTY {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}
	directory := fmt.Sprintf("%03d", generation)
	err := createIfNotExists(directory)
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
}

func createIfNotExists(directory string) error {
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
		if peep.location.X > 80 &&
			peep.location.Y > 80 {
			survivors = append(survivors, peep)
		}
	}
	return survivors
}

func fillWithRandomPeeps(world *World) {
	for i := 0; i < POPULATION; i++ {
		individual := createIndividual(world.XSize, world.YSize)
		world.addPeep(individual)
	}
}

func randomCoord(x int, y int) Coord {
	place := Coord{
		X: rand.Intn(x),
		Y: rand.Intn(y),
	}
	return place
}

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
					s.world.updateLocation(actions.peepID, loc)
				case MOVE_Y:
					loc := individual.location
					loc.Y += int(value * MOVEMENT)
					s.world.updateLocation(actions.peepID, loc)
				}
			}
		}
	}
}

func (s *simulation) startPeeking() chan act {
	// we start all the individuals in separate goroutines, and then wait for them to finish
	peepActions := make(chan act, len(s.world.peeps))
	var wg sync.WaitGroup
	for id, peep := range s.world.peeps {
		wg.Add(1)
		go func(peep *Individual, id int) {
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
