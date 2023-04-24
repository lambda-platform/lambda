package custom_model

import (
	"io"
)

func MarshalInt32(f int32) Marshaler {
	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, string(f))
	})
}

func UnmarshalInt32(v interface{}) (int32, error) {
	return v.(int32), nil
}
