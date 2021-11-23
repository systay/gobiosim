package main


type (
	Sensor uint8
	Action uint8
)

// Place the sensor neuron you want enabled prior to NUM_SENSES. Any
// that are after NUM_SENSES will be disabled in the simulator.
// If new items are added to this enum, also update the name functions
// in analysis.cpp.
// I means data about the individual, mainly stored in Indiv
// W means data about the environment, mainly stored in Peeps or Grid
const (
	LOC_X             Sensor = iota // I distance from left edge
	LOC_Y                          // I distance from bottom
	BOUNDARY_DIST_X                // I X distance to nearest edge of world
	BOUNDARY_DIST                  // I distance to nearest edge of world
	BOUNDARY_DIST_Y                // I Y distance to nearest edge of world
	AGE                            // I
	NUM_SENSES                     // <<------------------ END OF ACTIVE SENSES MARKER
)

// Place the action neuron you want enabled prior to NUM_ACTIONS. Any
// that are after NUM_ACTIONS will be disabled in the simulator.
// If new items are added to this enum, also update the name functions
// in analysis.cpp.
// I means the action affects the individual internally (Indiv)
// W means the action also affects the environment (Peeps or Grid)
const (
	MOVE_X                Action = iota // W +- X component of movement
	MOVE_Y                             // W +- Y component of movement
	NUM_ACTIONS                        // <<----------------- END OF ACTIVE ACTIONS MARKER
)

func (a Action) String() string {
	return actionNames[a]
}
func (s Sensor) String() string {
	return sensorNames[s]
}

var actionNames = map[Action]string{
	MOVE_X:                "MOVE_X",
	MOVE_Y:                "MOVE_Y",
}

var sensorNames = map[Sensor]string{
	LOC_X:             "LOC_X",
	LOC_Y:             "LOC_Y",
	BOUNDARY_DIST_X:   "BOUNDARY_DIST_X",
	BOUNDARY_DIST:     "BOUNDARY_DIST",
	BOUNDARY_DIST_Y:   "BOUNDARY_DIST_Y",
	AGE:               "AGE",
}
