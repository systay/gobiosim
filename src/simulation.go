package main

import (
	"fmt"
	"math/rand"
	"time"
)

type simulation struct {
	world World
}

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
					loc.X += int(value)
					s.world.updateLocation(peepIdx, loc)
				case MOVE_Y:
					loc := individual.location
					loc.Y += int(value)
					s.world.updateLocation(peepIdx, loc)
				}
			}
		}
	}
}

func init() {
	seed := time.Now().UnixNano()
	//nano := int64(1637780343848163000)
	fmt.Printf("rand seed: %d\n", seed)
	rand.Seed(seed)
}

func main() {
	world := World{
		StepsPerGeneration: 250,
		XSize:              100,
		YSize:              100,
		cells:              make([]Cell, 100*100),
		peeps:              []*Individual{},
	}
	for i := 0; i < 100; i++ {
		individual := createIndividual(world.XSize, world.YSize)
		world.addPeep(individual)
	}

	s := &simulation{
		world: world,
	}

	steps := s.world.StepsPerGeneration
	for steps > 0 {
		steps--
		s.step()
	}
}

func createIndividual(x, y int) *Individual {
	genome := makeRandomGenome(10)
	brain, err := genome.buildNet2()
	if err != nil {
		panic(err)
	}
	place := Coord{
		X: rand.Intn(x),
		Y: rand.Intn(y),
	}
	peep := &Individual{
		location:       place,
		birthPlace:     place,
		age:            0,
		brain:          *brain,
		responsiveness: 1,
		oscPeriod:      0,
		longProbeDist:  0,
		lastMoveDir:    0,
		challengeBits:  0,
	}
	return peep
}
