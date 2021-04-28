package plugin

type plugin struct {
	policyPath string
}

type DroneConfig struct {
	Kind string `yaml:"kind" json:"kind"`
	Type string `yaml:"type" json:"type"`
	Name string `yaml:"name" json:"name"`
	Data []byte `yaml:"-" json:"-"`
}

type Promotion struct {
	Build  int64  `json:"build,omitempty"`
	Deploy string `json:"deploy,omitempty"`
}
