package custom_model

import (
	"io"
)

func MarshalByte(f []byte) Marshaler {

	return WriterFunc(func(w io.Writer) {
		w.Write(f)
	})
}

func UnmarshalByte(v interface{}) ([]byte, error) {

	bytes, ok := v.([]byte)
	if !ok {
		return nil, nil
	} else {
		return bytes, nil
	}

}
