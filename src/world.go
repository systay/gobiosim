package main

type (
	Cell = uint16
	World struct {
		StepsPerGeneration int
		XSize int
		YSize int
		cells []Cell
		peeps []*Individual
	}
)

const EMPTY uint16 = 0
const BARRIER uint16 = 0xffff

func executeStep(individual Individual, world World)  {
	
}