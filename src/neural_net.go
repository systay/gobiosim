package main

import (
	"fmt"
	"strings"
)

type (
	// An individual's "brain" is a neural net specified by a set
	// of Genes where each Gene specifies one connection in the neural net (see
	// Genome comments above). Each neuron has a single output which is
	// connected to a set of sinks where each sink is either an action output
	// or another neuron. Each neuron has a set of input sources where each
	// source is either a sensor or another neuron. There is no concept of
	// layers in the net: it's a free-for-all topology with forward, backwards,
	// and sideways connection allowed. Weighted connections are allowed
	// directly from any source to any action.

	// Currently, the genome does not specify the activation function used in
	// the neurons.

	// When the input is a sensor, the input value to the sink is the raw
	// sensor value of type float and depends on the sensor. If the output
	// is an action, the source's output value is interpreted by the action
	// node and whether the action occurs or not depends on the action's
	// implementation.

	// In the genome, neurons are identified by 15-bit unsigned indices,
	// which are reinterpreted as values in the range 0..p.genomeMaxLength-1
	// by taking the 15-bit index modulo the max number of allowed neurons.
	// In the neural net, the neurons that end up connected get new indices
	// assigned sequentially starting at 0.

	// NeuralNet encodes an individuals brain
	NeuralNet struct {
		// Contains the sensors that feed to something interesting
		Sensors []Sensor

		Neurons []*Neuron

		Connections []Connection
	}

	Connection struct {
		From Source // Either a sensor, or a neuron
		To   Sink   // either an action, or a neuron

		// multiplier is a value between -1.0..1.0,
		multiplier float64
	}

	Neuron struct {
		id int
		// when value reaches 1.0, the neuron will fire
		// until then it will accumulate into the value,
		// a state which survives between steps
		value float64
	}

	Source interface{ Get() }
	Sink   interface{ Set() }

	ActionSink struct {
		action Action
	}

	SensorInput struct {
		s   Sensor
		idx int
	}
)

func (s SensorInput) Get() {}
func (n *Neuron) Get()     {}
func (n *Neuron) Set()     {}
func (n ActionSink) Set()  {}

func (n *NeuralNet) getSensorOffset(s Sensor) int {
	for idx, sensor := range n.Sensors {
		if s == sensor {
			return idx
		}
	}

	n.Sensors = append(n.Sensors, s)
	return len(n.Sensors) - 1
}

func (n *NeuralNet) getNeuronByID(id int) *Neuron {
	neuron := n.Neurons[id]
	if neuron == nil {
		neuron = &Neuron{id: id}
		n.Neurons[id] = neuron
	}
	return neuron
}

func (n *NeuralNet) String() string {
	var sensors []string
	for idx, conn := range n.Connections {
		sensors = append(sensors, fmt.Sprintf("%d:%s", idx, conn))
	}
	return strings.Join(sensors, "\n")
}

func (conn Connection) String() (result string) {
	sensor, ok := conn.From.(SensorInput)
	if ok {
		result = sensor.s.String()
	} else {
		result = fmt.Sprintf("N%d", conn.From.(*Neuron).id)
	}

	result += fmt.Sprintf(" -[%03f]-> ", conn.multiplier)

	action, ok := conn.To.(ActionSink)
	if ok {
		result += action.action.String()
	} else {
		result = fmt.Sprintf("N%d", conn.To.(*Neuron).id)
	}

	return
}
