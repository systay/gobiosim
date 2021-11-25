package main

import "testing"

var s *simulation

func init() {
	world := &World{
		StepsPerGeneration: STEPS_PER_GEN,
		XSize:              SIZE,
		YSize:              SIZE,
		cells:              make([]Cell, 100*100),
	}
	fillWithRandomPeeps(world)
	s = &simulation{
		world: world,
	}
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.step()
	}
}
