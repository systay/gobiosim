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
		brain      NeuralNet
	}

	// Actions encodes the actions taken by an individual. The offset corresponds to the Action value,
	// and the float value at the index says how much an action is taken
	Actions = []float64
)

func (i *Individual) step(world World) Actions {
	inputs := make([]float64, 0, len(i.brain.Sensors))
	for _, sensor := range i.brain.Sensors {
		value := getSensorValue(i, world, sensor)
		inputs = append(inputs, value)
	}
	actions := make(Actions, NUM_ACTIONS)
	var neuronFirings []*Neuron

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

	for _, conn := range i.brain.Connections {
		sensor, ok := conn.From.(SensorInput)
		if !ok {
			continue
		}
		srcValue := inputs[sensor.idx]
		handleFiring(conn.To, conn.multiplier*srcValue)
	}

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

const MUTATION_RATE = 0.01 * math.MaxFloat64

func plusMinusOne() int {
	if rand.Intn(2) == 0 {
		return -1
	}
	return 1
}

func (i *Individual) clone() *Individual {
	clone := *i
	clone.age = 0
	for _, gene := range clone.genome.genes {
		if rand.NormFloat64() < MUTATION_RATE {
			switch rand.Intn(4) {
			case 0:
				gene.sourceID = uint8(int(gene.sourceID) + plusMinusOne())
			case 1:
				gene.sinkID = uint8(int(gene.sinkID) + plusMinusOne())
			case 2:
				gene.noOfNeurons = uint8(int(gene.noOfNeurons) + plusMinusOne())
			case 3:
				gene.weight = int16(int(gene.weight) + plusMinusOne())
			}
		}
	}
	return &clone
}

func getSensorValue(i *Individual, w World, s Sensor) float64 {
	switch s {
	case LOC_X:
		// map current X location to value between 0.0..1.0
		return float64(i.location.X / w.XSize)
	case LOC_Y:
		// map current Y location to value between 0.0..1.0
		return float64(i.location.Y / w.YSize)
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

	}
	panic("oh noes")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
