package mongodata

type Decoder[M any] interface {
	Schema() string
	Decode(Document) (*M, error)
}

type Document interface {
	Decode(any) error
}

func Decode[M any](decoder Document) (*M, error) {
	var m M
	return &m, decoder.Decode(&m)
}
