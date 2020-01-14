package json

import "encoding/json"

var marshal = struct {
	marshaller   func(v interface{}) ([]byte, error)
	unMarshaller func(data []byte, v interface{}) error
}{
	marshaller:   json.Marshal,
	unMarshaller: json.Unmarshal,
}

func ReplaceHandler(marshaller func(v interface{}) ([]byte, error), unMarshaller func(data []byte, v interface{}) error) {
	marshal.marshaller = marshaller
	marshal.unMarshaller = unMarshaller
}

func Marshal(v interface{}) ([]byte, error) {
	return marshal.marshaller(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return marshal.unMarshaller(data, v)
}
