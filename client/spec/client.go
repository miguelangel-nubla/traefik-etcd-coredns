package spec

type Client interface {
	GetName() string
	Init() error
	Update(Record) error
	Delete(Record) error
	Close() error
}
