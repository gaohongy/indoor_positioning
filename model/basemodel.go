package model

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}

type Coordinate struct {
	Coordinate_x float64 `json:"coordinate_x"`
	Coordinate_y float64 `json:"coordinate_y"`
	Coordinate_z float64 `json:"coordinate_z"`
}
