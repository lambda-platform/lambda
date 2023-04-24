package custom_model

import (
	"io"
)

func MarshalByte(f []byte) Marshaler {

	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, string(f))
	})
}

func UnmarshalByte(v interface{}) (string, error) {

	bytes, ok := v.([]byte)
	if !ok {
		return string(bytes), nil
	} else {
		return "", nil
	}

}
