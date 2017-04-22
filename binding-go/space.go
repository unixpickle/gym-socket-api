package gym

// Space defines an action or observation space.
type Space struct {
	// Space type, such as "Discrete", "Tuple", "MultiBinary",
	// "MultiDiscrete", or "Box".
	Type string `json:"type"`

	// Number of elements, used for MultiBinary and
	// Discrete spaces.
	N int `json:"n"`

	// Boundaries used for Box and MultiDiscrete.
	// For Box, these are flattened.
	Low  []float64 `json:"low"`
	High []float64 `json:"high"`

	// Shape for Box spaces.
	Shape []int `json:"shape"`

	// Subspaces for Tuple spaces.
	Subspaces []*Space `json:"subspaces"`
}
