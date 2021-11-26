package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNeuralNet_String(t *testing.T) {
	it := makeRandomGenome(10)
	net, err := it.buildNet()
	require.NoError(t, err)
	fmt.Println(net.String())
}
