package data

type Parser interface {
	Parse([]byte) (interface{}, error)
}
