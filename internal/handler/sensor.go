package handler

type Sensor struct {
	Name string `yaml:"name"`
	MAC  string `yaml:"mac"`
}

type Reading struct {
	Temperature float64
	Humidity    float64
	Voltage     float64
}
