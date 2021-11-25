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
				weight:         100,
			},
		},
	}

	net2, err := genome.buildNet()
	require.NoError(t, err)
	fmt.Println(net2)
}