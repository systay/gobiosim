package main

type (
	Cell  = uint16
	World struct {
		StepsPerGeneration int
		XSize              int
		YSize              int
		cells              []Cell
		peeps              []*Individual
		survivalArea       Area
		barriers           []Area
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
	return world.offsetXY(place.X, place.Y)
}

func (world *World) offsetXY(x, y int) int {
	return y*world.XSize + x
}

func limit(v, max int) int {
	if v > max {
		return max
	}
	if v < 0 {
		return 0
	}
	return v
}

func (world *World) updateLocation(peepIdx int, location Coord) (blocked bool) {
	// contain ourselves to the given world
	location.X = limit(location.X, world.XSize-1)
	location.Y = limit(location.Y, world.YSize-1)

	// first we check if the spot is taken. if it isn't, we just ignore the location change
	newOffset := world.offset(location)
	if world.cells[newOffset] != EMPTY {
		return true
	}

	// if the spot is empty, we can move the peep to the new location
	oldOffset := world.offset(world.peeps[peepIdx].location)
	world.cells[oldOffset] = EMPTY
	world.cells[newOffset] = Cell(peepIdx)
	world.peeps[peepIdx].location = location
	return false
}

func (world *World) clearAll() {
	for _, peep := range world.peeps {
		world.cells[world.offset(peep.location)] = EMPTY
	}
	world.peeps = nil
}

func (world *World) fillBarriers() {
	for _, barrier := range world.barriers {
		for x := barrier.TopLeft.X; x < barrier.BottomRight.X; x++ {
			for y := barrier.TopLeft.Y; y < barrier.BottomRight.Y; y++ {
				idx := world.offsetXY(x, y)
				world.cells[idx] = BARRIER
			}
		}
	}
}
