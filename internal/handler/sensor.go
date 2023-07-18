package handler

type Sensor struct {
	Name string `yaml:"name"`
	MAC  string `yaml:"mac"`
}

type Reading struct {
	Temperature float32
	Humidity    uint8
	Voltage     float32
	Battery     uint8
}
