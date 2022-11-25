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

	if s != "null" && s != "" {
		nt, err := time.Parse(DB.CtLayout, s)
		if err != nil {
			return DB.Date{
				IsNotNull: false,
			}, err
		}
		return DB.Date{
			DateValue: nt,
			IsNotNull: true,
		}, err
	} else {
		return DB.Date{
			IsNotNull: false,
		}, nil
	}

}
