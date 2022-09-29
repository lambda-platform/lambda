package custom_model

import (
	"github.com/lambda-platform/lambda/DB"
	"gorm.io/gorm"
	"io"
	"strings"
	"time"
)

func MarshalGormDeletedAt(f gorm.DeletedAt) Marshaler {
	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, f.Time.String())
	})
}

func UnmarshalGormDeletedAt(v interface{}) (gorm.DeletedAt, error) {

	s := strings.Trim(v.(string), `"`)
	nt, err := time.Parse(DB.CtLayout, s)
	ct := time.Time(nt)

	dct := gorm.DeletedAt{}
	dct.Time = ct
	return dct, err
}
