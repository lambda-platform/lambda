package DB

import (
	"fmt"
	"strings"
	"time"
	"database/sql"
	"database/sql/driver"
)

// CustomTime provides an example of how to declare a new time Type with a custom formatter.
// Note that time.Time methods are not available, if needed you can add and cast like the String method does
// Otherwise, only use in the json struct at marshal/unmarshal time.
type Date time.Time

const CtLayout = "2006-01-02"

func Time(ct *Date)  time.Time{
	return time.Time(*ct)
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}
func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}
// UnmarshalJSON Parses the json string in the custom format
func (ct *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(CtLayout, s)
	*ct = Date(nt)
	return
}
func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}
// MarshalJSON writes a quoted string in the custom format
func (ct Date) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *Date) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(CtLayout))
}
