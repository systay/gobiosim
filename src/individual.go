package main

import (
	"math"
	"math/rand"
)

type (
	Individual struct {
		genome     Genome
		location   Coord
		birthPlace Coord
		age        uint16
		wasBlocked bool // will be true if this individual was not able to do an action last step because it was blocked
		brain      *NeuralNet
	}

	// Actions encodes the actions taken by an individual. The offset corresponds to the Action value,
	// and the float value at the index says how much an action is taken
	Actions = []float64
)

func createIndividual(world *World) *Individual {
	genome := makeRandomGenome(rand.Intn(20) + 2)
	brain, err := genome.buildNet()
	if err == TooSimple {
		return createIndividual(world)
	}
	if err != nil {
		panic(err)
	}
	place := world.randomCoord()
	peep := &Individual{
		genome:     genome,
		location:   place,
		birthPlace: place,
		age:        0,
		brain:      brain,
	}

	return peep
}

func (i *Individual) step(world *World) Actions {
	// First we build the sensor inputs that the brains uses into a slice
	inputs := make([]float64, 0, len(i.brain.Sensors))
	for _, sensor := range i.brain.Sensors {
		value := getSensorValue(i, world, sensor)
		inputs = append(inputs, value)
	}
	actions := make(Actions, NUM_ACTIONS)
	var neuronFirings []*Neuron

	// this is the function that will be called whenever there is a signal.
	// The recipient of the signal can be a neuron, or it can be an action sink
	handleFiring := func(to Sink, v float64) {
		switch dst := to.(type) {
		case ActionSink:
			actions[dst.action] += v
		case *Neuron:
			dst.value += v
			for dst.value > 1 {
				// a neuron will keep firing until it gets it's internal state under 1
				neuronFirings = append(neuronFirings, dst)
				dst.value -= 1
			}
		}
	}

	// Next step is to fire the connections to the sensor inputs
	for _, conn := range i.brain.Connections {
		sensor, ok := conn.From.(SensorInput)
		if !ok {
			continue
		}
		srcValue := inputs[sensor.idx]
		handleFiring(conn.To, conn.multiplier*srcValue)
	}

	// If neurons received signals in the last step, we could now have new signals that we need to handle
	// Since the neural net is not an acyclic graph, we limit the number of signals we allow per step and individual
	// We could deal with this in other ways, this method was chosen mostly because it is simple
	iterLeft := 10
	for len(neuronFirings) > 0 && iterLeft > 0 {
		iterLeft--
		current := neuronFirings[0]
		neuronFirings = neuronFirings[1:]
		for _, conn := range i.brain.Connections {
			if conn.From != current {
				continue
			}
			handleFiring(conn.To, conn.multiplier*1)
		}
	}
	i.age++

	for idx, action := range actions {
		actions[idx] = math.Tanh(action)
	}

	return actions
}

func plusMinusOne() int {
	if rand.Intn(2) == 0 {
		return -1
	}
	return 1
}

func (i *Individual) clone() *Individual {
	clone := *i
	clone.age = 0
	var mutant bool
	ready := false
	for !ready {
		clone.genome, mutant = clone.genome.clone()
		if mutant {
			net, err := clone.genome.buildNet()
			if err != nil {
				if err != TooSimple {
					panic(err)
				}
			} else {
				ready = true
			}
			clone.brain = net
		}
	}
	return &clone
}

func getSensorValue(i *Individual, w *World, s Sensor) float64 {
	switch s {
	case LOC_X:
		// map current X location to value between 0.0..1.0
		return float64(i.location.X) / float64(w.XSize)
	case LOC_Y:
		// map current Y location to value between 0.0..1.0
		return float64(i.location.Y) / float64(w.YSize)
	case BOUNDARY_DIST:
		// Finds the closest boundary, compares that to the max possible dist
		// to a boundary from the center, and converts that linearly to the
		// sensor range 0.0..1.0
		x := getSensorValue(i, w, BOUNDARY_DIST_X)
		y := getSensorValue(i, w, BOUNDARY_DIST_Y)

		return math.Min(x, y)

	case BOUNDARY_DIST_X:
		maxDist := float64(w.XSize / 2)
		return float64(min(i.location.X, w.XSize-i.location.X-1)) / maxDist

	case BOUNDARY_DIST_Y:
		maxDist := float64(w.YSize / 2)
		return float64(min(i.location.Y, w.YSize-i.location.Y-1)) / maxDist

	case AGE:
		// sets the age to a normalized value between 0 and 1
		return float64(i.age) / float64(w.StepsPerGeneration)

	case BLOCK:
		if i.wasBlocked {
			return 1
		}
		return 0

	}
	panic("oh noes")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
