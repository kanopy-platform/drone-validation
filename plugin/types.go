package plugin

type plugin struct {
	policyPath string
}

type DroneConfig struct {
	Kind  string `yaml:"kind" json:"kind"`
	Type  string `yaml:"type" json:"type"`
	Name  string `yaml:"name" json:"name"`
	Steps []struct {
		Name      string   `yaml:"name" json:"name"`
		Image     string   `yaml:"image" json:"image"`
		Commands  []string `yaml:"commands" json:"commands"`
		DependsOn []string `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	} `yaml:"steps" json:"steps"`
	Trigger struct {
		Branch []string `yaml:"branch" json:"branch"`
	} `yaml:"trigger" json:"trigger"`
}
