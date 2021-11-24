package main

type (
	Cell  = uint16
	World struct {
		StepsPerGeneration int
		XSize              int
		YSize              int
		cells              []Cell
		peeps              []*Individual
	}
)

const EMPTY uint16 = 0
const BARRIER uint16 = 0xffff

func (world *World) addPeep(individual *Individual) {
	id := len(world.peeps)
	world.peeps = append(world.peeps, individual)
	offset := world.offset(individual.birthPlace)
	world.cells[offset] = uint16(id)
}

func (world *World) offset(place Coord) int {
	return place.Y*world.XSize + place.X
}

func (world *World) updateLocation(peepIdx int, location Coord) {
	// first we check if the spot is taken. if it isn't, we just ignore the location change
	newOffset := world.offset(location)
	if world.cells[newOffset] != EMPTY {
		return
	}

	// if the spot is empty, we can move the peep to the new location
	oldOffset := world.offset(world.peeps[peepIdx].location)
	world.cells[oldOffset] = EMPTY
	world.cells[newOffset] = Cell(peepIdx)
	world.peeps[peepIdx].location = location
}
