package custom_model

import (
	"io"
)

func MarshalInt64(f int32) Marshaler {
	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, string(f))
	})
}

func UnmarshalInt64(v interface{}) (int64, error) {
	return v.(int64), nil
}
