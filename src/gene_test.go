package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleGenome(t *testing.T) {
	genome := Genome{
		noOfNeurons: 1,
		genes:       []Gene{
			{
				sourceIsSensor: true,
				sourceID:       uint8(LOC_X),
				sinkIsAction:   true,
				sinkID:         uint8(MOVE_X),
				noOfNeurons:    0,
				weight:         100,
			},
		},
	}

	net2, err := genome.buildNet2()
	require.NoError(t, err)
	fmt.Println(net2)
}