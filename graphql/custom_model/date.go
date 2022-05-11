package custom_model

import (
	"github.com/lambda-platform/lambda/DB"
	"io"
	"strings"
	"time"
)

func MarshalDate(f DB.Date) Marshaler {
	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, f.String())
	})
}


func UnmarshalDate(v interface{}) (DB.Date, error) {


	s := strings.Trim(v.(string), `"`)
	nt, err := time.Parse(DB.CtLayout, s)
	ct := DB.Date(nt)

	return ct, err
}
