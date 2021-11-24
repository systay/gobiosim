package main

import (
	"math"
	"math/rand"
)

type (
	// Gene specifies one synaptic connection in a neural net. Each
	// connection has an input (source) which is either a sensor or another neuron.
	// Each connection has an output, which is either an action or another neuron.
	// Each connection has a floating point weight derived from a signed 16-bit
	// value. The signed integer weight is scaled to a small range, then cubed
	// to provide fine resolution near zero.
	Gene struct {
		sourceIsSensor bool // sensor if true, otherwise neuron
		sourceID       uint8
		sinkIsAction   bool // action if true, otherwise neuron
		sinkID         uint8
		noOfNeurons    uint8 // determines how many neurons this individual has
		weight         int16
	}

	// Genome defines an individuals' genome. It consists of a set of genes.
	// The genome is used to build an individuals neural net brain
	// noOfNeurons should not be a
	Genome struct {
		genes       []Gene
		noOfNeurons int
	}
)

func (g *Gene) weightAsFloat() float32 {
	return float32(g.weight) / 8192.0
}

func randInt16() int16 {
	return int16(rand.Int63() >> 48)
}

func makeRandomGene() Gene {
	gene := Gene{}
	gene.sourceIsSensor = rand.Int()%2 == 0
	gene.sinkIsAction = rand.Int()%2 == 0
	gene.sourceID = randUint8()
	gene.sinkID = randUint8()
	gene.weight = randInt16()
	return gene
}

func makeRandomGenome(size int) Genome {
	genome := Genome{
		genes:       make([]Gene, 0, size),
		noOfNeurons: size,
	}
	for i := 0; i < size; i++ {
		genome.genes = append(genome.genes, makeRandomGene())
	}
	return genome
}

func randUint8() uint8 {
	return uint8(rand.Int31n(255))
}

type path = []Connection

func (g Genome) buildNet() NeuralNet {
	result := NeuralNet{
		Neurons: make([]*Neuron, g.noOfNeurons),
	}
	tentative := []path{}
	done := make([]bool, len(g.genes))
	doneCount := 0

	// we first find all action sinks and add them to the net if they connect a Sensor with an Action.
	// if the Action is connected through a Neuron, we add it to the tentative list, until we know that
	// there is a connection from some Sensor to the neuron, directly or indirectly
	for id, gene := range g.genes {
		if gene.sinkIsAction {
			done[id] = true
			doneCount++
			output := createActionSink(gene.sinkID)
			multiplier := float64(gene.weight) / float64(math.MaxInt16)

			if gene.sourceIsSensor {
				input := createSensorInput(gene.sourceID, result)

				// we have a connection straight from a sensor to an action sink
				// no need to do anything else - we just add this as is
				conn := Connection{
					From:       input,
					To:         output,
					multiplier: multiplier,
				}
				result.Connections = append(result.Connections, conn)
				continue
			}

			neuron := result.getNeuronByID(int(gene.sourceID) % g.noOfNeurons)
			conn := Connection{
				From:       neuron,
				To:         output,
				multiplier: multiplier,
			}
			tentative = append(tentative, path{conn})
		}
	}

	for doneCount < len(g.genes) {
		for id, gene := range g.genes {
			if done[id] {
				continue
			}

			conn := Connection{
				From:       nil,
				To:         nil,
				multiplier: float64(gene.weight) / float64(math.MaxInt16),
			}

			if gene.sourceIsSensor {
				conn.From = createSensorInput(gene.sourceID, result)
			} else {
				conn.From = result.getNeuronByID(int(gene.sourceID) % g.noOfNeurons)
			}

			conn.To = result.getNeuronByID(int(gene.sinkID) % g.noOfNeurons)

		}
	}
	return result
}

func createActionSink(sinkID uint8) ActionSink {
	action := Action(sinkID % uint8(NUM_ACTIONS))
	output := ActionSink{
		action: action,
	}
	return output
}

func createSensorInput(source uint8, result NeuralNet) SensorInput {
	sensor := Sensor(source % uint8(NUM_SENSES))
	sensorIdx := -1
	for i, s := range result.Sensors {
		if s == sensor {
			// the sensor is already added to the net - just store the offset
			sensorIdx = i
		}
	}
	if sensorIdx < 0 {
		// not already added to the net - let's add it
		result.Sensors = append(result.Sensors, sensor)
		sensorIdx = len(result.Sensors) - 1
	}
	input := SensorInput{
		s:   sensor,
		idx: sensorIdx,
	}
	return input
}
