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

func (g Gene) weightAsFloat() float32 {
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
	gene.weight = -randInt16() + randInt16()
	return gene
}

func makeRandomGenome(size int) Genome {
	genome := Genome{
		genes:       make([]Gene, 0, size),
		noOfNeurons: int(math.Sqrt(float64(size))),
	}
	for i := 0; i < size; i++ {
		genome.genes = append(genome.genes, makeRandomGene().normalize(genome.noOfNeurons))
	}
	return genome
}

func randUint8() uint8 {
	return uint8(rand.Int31n(255))
}

type identifiable struct {
	objects []interface{}
}

func (id *identifiable) idOf(obj interface{}) (int, bool) {
	for idx, item := range id.objects {
		if obj == item {
			return idx, false
		}
	}
	id.objects = append(id.objects, obj)
	return len(id.objects) - 1, true
}

func (g Genome) buildNet() (*NeuralNet, error) {
	graph, paths, err := buildGraphAndPaths(&g)
	if err != nil {
		return nil, err
	}

	result := NewNeuralNet(g.noOfNeurons)

	seen := map[int]interface{}{}
	for _, path := range paths {
		for _, vertix := range path {
			from := vertix.from
			to := vertix.to

			vIdx := from*graph.size + to
			if _, ok := seen[vIdx]; ok {
				continue
			}
			seen[vIdx] = nil
			con := Connection{
				From:       nil,
				To:         nil,
				multiplier: vertix.data.(float64),
			}

			switch src := graph.GetNode(from).(type) {
			case Sensor:
				con.From = SensorInput{
					s:   src,
					idx: result.getSensorOffset(src),
				}
			case int:
				con.From = result.getNeuronByID(src)
			}

			switch src := graph.GetNode(to).(type) {
			case Action:
				con.To = ActionSink{
					action: src,
				}
			case int:
				con.To = result.getNeuronByID(src)
			}

			result.Connections = append(result.Connections, con)
		}
	}

	return result, nil
}

func NewNeuralNet(noOfNeurons int) *NeuralNet {
	return &NeuralNet{
		Neurons: make([]*Neuron, noOfNeurons),
	}
}

func buildGraphAndPaths(g *Genome) (*Graph, []Path, error) {
	maxPossibleSize := len(g.genes) * 2
	graph := NewGraph(maxPossibleSize)
	nodes := &identifiable{}
	var sensors, actions []int
	for _, gene := range g.genes {
		var isSensor, isAction bool
		var obj interface{}
		if gene.sourceIsSensor {
			obj = getSensor(gene.sourceID)
			isSensor = true
		} else {
			if g.noOfNeurons == 0 {
				g.noOfNeurons = 1
			}
			obj = int(gene.sourceID) % g.noOfNeurons
		}
		srcID, added := nodes.idOf(obj)
		if added {
			err := graph.AddNode(srcID, obj)
			if err != nil {
				return nil, nil, err
			}
			if isSensor {
				sensors = append(sensors, srcID)
			}
		}

		if gene.sinkIsAction {
			obj = getAction(gene.sinkID)
			isAction = true
		} else {
			if g.noOfNeurons == 0 {
				g.noOfNeurons = 1
			}
			obj = int(gene.sourceID) % g.noOfNeurons
		}
		dstID, added := nodes.idOf(obj)
		if added {
			err := graph.AddNode(dstID, obj)
			if err != nil {
				return nil, nil, err
			}
			if isAction {
				actions = append(actions, dstID)
			}
		}

		weight := float64(gene.weight) / float64(math.MaxInt16)
		err := graph.AddVertix(srcID, dstID, weight)
		if err != nil {
			return nil, nil, err
		}
	}

	paths := graph.PathsBetween(sensors, actions)
	return graph, paths, nil
}

func (g Gene) normalize(neuronCount int) Gene {
	if g.sourceIsSensor {
		g.sourceID = g.sourceID % uint8(NUM_SENSES)
	} else {
		g.sourceID = g.sourceID % uint8(neuronCount)
	}
	if g.sinkIsAction {
		g.sinkID = g.sinkID % uint8(NUM_ACTIONS)
	} else {
		g.sinkID = g.sinkID % uint8(neuronCount)
	}

	return g
}

func getSensor(source uint8) Sensor {
	return Sensor(source % uint8(NUM_SENSES))
}

func getAction(source uint8) Action {
	return Action(source % uint8(NUM_ACTIONS))
}

func shouldMutate() bool {
	return rand.Intn(1000) < MUTATION_RATE
}

func (g Genome) clone() (output Genome, mutant bool) {
	output = g
	if len(g.genes) == 0 {
		if shouldMutate() {
			mutant = true
			output.genes = append(output.genes, makeRandomGene())
		}
		return
	}

	for idx, gene := range output.genes {
		if shouldMutate() {
			mutant = true
			switch rand.Intn(3) {
			case 0:
				gene.sourceID = uint8(int(gene.sourceID) + plusMinusOne())
			case 1:
				gene.sinkID = uint8(int(gene.sinkID) + plusMinusOne())
			case 2:
				gene.weight = int16(int(gene.weight) + plusMinusOne()*10)
			}
			normalize := gene.normalize(output.noOfNeurons)
			output.genes[idx] = normalize
		}
	}

	if shouldMutate() {
		// add a new gene
		mutant = true
		pos := rand.Intn(len(output.genes))
		output.genes = append(output.genes[:pos+1], output.genes[pos:]...)
		output.genes[pos] = makeRandomGene()
	}

	if shouldMutate() {
		// remove gene
		mutant = true
		pos := rand.Intn(len(output.genes))
		output.genes = append(output.genes[:pos], output.genes[pos+1:]...)
	}

	if shouldMutate() {
		// add/remove neuron
		mutant = true
		output.noOfNeurons += plusMinusOne()
	}
	return
}