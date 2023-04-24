package custom_model

import (
	"github.com/pkg/errors"
	"io"
)

func MarshalByte(f string) Marshaler {

	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, string(f))
	})
}

func UnmarshalByte(v interface{}) ([]byte, error) {
	bytes, ok := v.(string)
	if !ok {
		return nil, errors.New("could not unmarshal byte, value is not a string")
	}

	return []byte(bytes), nil
}
