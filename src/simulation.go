package main

import (
	"fmt"
	"math/rand"
	"time"
)

type simulation struct {
	world World
}

const MOVEMENT = 3

func (s *simulation) step() {
	var peepActions []Actions
	for _, peep := range s.world.peeps {
		peepActions = append(peepActions, peep.step(s.world))
	}
	for peepIdx, actions := range peepActions {
		individual := s.world.peeps[peepIdx]
		for act, value := range actions {
			if value > 0 {
				switch Action(act) {
				case MOVE_X:
					loc := individual.location
					loc.X += int(value * MOVEMENT)
					s.world.updateLocation(peepIdx, loc)
				case MOVE_Y:
					loc := individual.location
					loc.Y += int(value * MOVEMENT)
					s.world.updateLocation(peepIdx, loc)
				}
			}
		}
	}
}

func init() {
	seed := time.Now().UnixNano()
	// nano := int64(1637780343848163000)
	fmt.Printf("rand seed: %d\n", seed)
	rand.Seed(seed)
}

const POPULATION = 80

func main() {

	world := World{
		StepsPerGeneration: 250,
		XSize:              100,
		YSize:              100,
		cells:              make([]Cell, 100*100),
		peeps:              []*Individual{},
	}
	for i := 0; i < POPULATION; i++ {
		individual := createIndividual(world.XSize, world.YSize)
		world.addPeep(individual)
	}

	s := &simulation{
		world: world,
	}

	generations := 1000
	for generations > 0 {
		generations--

		steps := s.world.StepsPerGeneration
		for steps > 0 {
			time.Sleep(10 * time.Millisecond)
			world.printIndividuals()
			fmt.Println(steps)
			steps--
			s.step()
		}

		peeps := world.peeps
		world.clearAll()

		var survivors []*Individual
		for _, peep := range peeps {
			if peep.location.X > 40 && peep.location.X < 60 &&
				peep.location.Y > 40 && peep.location.Y < 60 {
				survivors = append(survivors, peep)
			}
		}

		if len(survivors) == 0 {
			panic("extinction")
		}

		copies := POPULATION / len(survivors)

		for _, survivor := range survivors {
			for i := 0; i <= copies; i++ {
				clone := survivor.clone()
				clone.location = randomCoord(world.XSize, world.YSize)
				world.addPeep(clone)
			}
		}

		fmt.Printf("%d survivors for generation %d", len(survivors), generations)
		time.Sleep(3 * time.Second)
	}
}

func createIndividual(x, y int) *Individual {
	genome := makeRandomGenome(10)
	brain, err := genome.buildNet2()
	if err != nil {
		panic(err)
	}
	place := randomCoord(x, y)
	peep := &Individual{
		location:   place,
		birthPlace: place,
		age:        0,
		brain:      *brain,
	}
	return peep
}

func randomCoord(x int, y int) Coord {
	place := Coord{
		X: rand.Intn(x),
		Y: rand.Intn(y),
	}
	return place
}
