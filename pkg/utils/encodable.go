package utils

type Encodable interface {
	ToMap() map[string]interface{}
}
