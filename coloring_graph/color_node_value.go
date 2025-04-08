package coloringgraph

type ColorNodeValue string

func (c ColorNodeValue) Serialize() []byte {
	return []byte(c)
}

func DeserializeColorNodeValue(data []byte) (ColorNodeValue, error) {
	return ColorNodeValue(data), nil
}
