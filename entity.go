package mongodata

type Entity[I EntityID] interface {
	ID() I
	Version() int
}

type EntityID interface {
	String() string
}

type Event interface {
	Type() string
}
