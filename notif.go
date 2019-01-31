package duckmail

type Notifier interface {
	Name() string
	Send(p Person) error
}
