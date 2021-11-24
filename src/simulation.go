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
					individual.location.X += int(value)
				case MOVE_Y:
					individual.location.Y += int(value)
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
	genome := makeRandomGenome(10)
	brain, err := genome.buildNet2()
	if err != nil {
		panic(err)
	}
	peep := &Individual{
		location:       &Coord{},
		birthPlace:     &Coord{},
		age:            0,
		brain:          *brain,
		responsiveness: 1,
		oscPeriod:      0,
		longProbeDist:  0,
		lastMoveDir:    0,
		challengeBits:  0,
	}
	s := &simulation{
		world: World{
			StepsPerGeneration: 250,
			XSize:              100,
			YSize:              100,
			cells:              make([]Cell, 100*100),
			peeps:              []*Individual{peep},
		},
	}

	steps := s.world.StepsPerGeneration
	for steps > 0 {
		steps--
		s.step()
	}
}
