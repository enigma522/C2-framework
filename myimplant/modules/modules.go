package modules

// Module interface that all modules must implement
type Module interface {
	Name() string
	Execute(command string, data []byte) (string, error)
}
