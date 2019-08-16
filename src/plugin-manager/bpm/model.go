package bpm

type Process struct {
	Name string `yaml:"name"`
	Executable string `yaml:"executable"`
	Args []string `yaml:"args"`
	Env map[string]string `yaml:"env"`
	Limits map[string]string `yaml:"limits"`
	// TODO: unsafe?
}

type Bpm struct {
	Processes []Process
}