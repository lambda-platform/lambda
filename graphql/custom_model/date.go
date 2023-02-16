package custom_model

import (
	"encoding/json"
	"github.com/lambda-platform/lambda/DB"
	"io"
	"strings"
	"time"
)

func MarshalDate(f DB.Date) Marshaler {

	return WriterFunc(func(w io.Writer) {
		if f.IsNotNull {
			w.Write([]byte(f.String()))
		} else {
			v, _ := json.Marshal(nil)
			w.Write(v)
		}

	})
}

func UnmarshalDate(v interface{}) (DB.Date, error) {

	s := strings.Trim(v.(string), `"`)

	if s != "null" && s != "" {
		nt, err := time.Parse(DB.DatetimeLayout, s)
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
