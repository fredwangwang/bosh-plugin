package executor

type Status struct {
	Running bool
}

type Executor interface {
	Prepare() error
	Start() error
	Stop() error
	Status() (Status, error)
}
