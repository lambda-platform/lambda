package custom_model

import (
	"io"
)

func MarshalByte(f []byte) Marshaler {

	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, string(f))
	})
}

func UnmarshalByte(v interface{}) ([]byte, error) {

	bytes, ok := v.([]byte)
	if !ok {
		return bytes, nil
	} else {
		return nil, nil
	}

}
