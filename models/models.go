package models

type Vector struct {
	X, Y, Z float64
}

type Motion struct {
	Acceleration    Vector
	AngularVelocity Vector
}

type Position struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	AngleX float64 `json:"angleX"`
	AngleY float64 `json:"angleY"`
	AngleZ float64 `json:"angleZ"`
}
