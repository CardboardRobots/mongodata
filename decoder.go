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

func DecodeModel[M any](decoder Document) (*M, error) {
	var m model[M]
	return &m.Value, decoder.Decode(&m)
}

type repositoryMapper[A any, B any] interface {
	Schema() string
	Map(A) (B, error)
}

type MapDecoder[A any, B any] struct {
	mapper repositoryMapper[A, B]
}

func NewMapDecoder[A any, B any](mapper repositoryMapper[A, B]) MapDecoder[A, B] {
	return MapDecoder[A, B]{
		mapper: mapper,
	}
}

func (rd MapDecoder[A, B]) Schema() string {
	return rd.mapper.Schema()
}

func (rd MapDecoder[A, B]) Decode(d Document) (*model[B], error) {
	a, err := Decode[model[A]](d)
	if err != nil {
		return nil, err
	}

	b, err := rd.mapper.Map(a.Value)
	if err != nil {
		return nil, err
	}

	return &model[B]{
		ID:      a.ID,
		Version: a.Version,
		Schema:  a.Schema,
		Events:  a.Events,
		Value:   b,
	}, nil
}
