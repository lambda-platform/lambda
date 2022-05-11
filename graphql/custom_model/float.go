package custom_model

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func MarshalFloat(f float32) Marshaler {
	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%g", f))
	})
}


func UnmarshalFloat(v interface{}) (float32, error) {
	switch v := v.(type) {
	case string:
		floatv, _ := strconv.ParseFloat(v, 32)
		return float32(floatv), nil
	case int:
		return float32(v), nil
	case int64:
		return float32(v), nil
	case float64:
		return float32(v), nil
	case float32:
		return v, nil
	case json.Number:
		floatv, _ := strconv.ParseFloat(string(v), 32)
		return float32(floatv), nil
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}
}
func MarshalFloat64(f float64) Marshaler {
	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%g", f))
	})
}


func UnmarshalFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case json.Number:
		return strconv.ParseFloat(string(v), 64)
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}
}

func MarshalFloat32(f float32) Marshaler {
	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%g", f))
	})
}


func UnmarshalFloat32(v interface{}) (float32, error) {
	switch v := v.(type) {
	case string:
		floatv, _ := strconv.ParseFloat(v, 32)
		return float32(floatv), nil
	case int:
		return float32(v), nil
	case int64:
		return float32(v), nil
	case float64:
		return float32(v), nil
	case float32:
		return v, nil
	case json.Number:
		floatv, _ := strconv.ParseFloat(string(v), 32)
		return float32(floatv), nil
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}
}
type Marshaler interface {
	MarshalGQL(w io.Writer)
}
type WriterFunc func(writer io.Writer)

func (f WriterFunc) MarshalGQL(w io.Writer) {
	f(w)
}
