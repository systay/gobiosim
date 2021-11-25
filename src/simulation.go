package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"sync"
	"time"
)

type simulation struct {
	world *World
}

const (
	MOVEMENT      = 3
	POPULATION    = 1000
	MUTATION_RATE = 10
	GENERATIONS   = 5000
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

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func main() {
	// _ = termbox.Init()
	// defer termbox.Close()
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

	generations := GENERATIONS
	for generations > 0 {
		generations--

		steps := s.world.StepsPerGeneration
		for steps > 0 {
			steps--
			s.step()
			// world.printIndividuals()
			// time.Sleep(100)
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

		fmt.Printf("%d %d\n", generations, len(survivors))
		// tbprint(0, 0, termbox.ColorWhite, termbox.ColorDefault, fmt.Sprintf("%d %d\n", generations, len(survivors)))
	}
	fmt.Println("done")
}

func cull(world *World) []*Individual {
	peeps := world.peeps
	world.clearAll()

	var survivors []*Individual
	for _, peep := range peeps {
		if peep.location.X > 40 && peep.location.X < 60 &&
			peep.location.Y > 40 && peep.location.Y < 60 {
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
