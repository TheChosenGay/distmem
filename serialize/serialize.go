package serialize

type Serializer interface {
	Marshal(v any) ([]byte, error)
}
