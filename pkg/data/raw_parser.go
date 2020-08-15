package data

type rawParser struct{}

var Raw Parser = &rawParser{}

func (rawParser) Parse(buf []byte) (interface{}, error) {
	return buf, nil
}
