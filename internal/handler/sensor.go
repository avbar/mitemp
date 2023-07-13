package handler

type Sensor struct {
	Name string
	MAC  string
}

type Reading struct {
	Temperature float64
	Humidity    float64
	Voltage     float64
}
