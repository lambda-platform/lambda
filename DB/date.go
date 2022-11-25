package DB

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Date CustomTime provides an example of how to declare a new time Type with a custom formatter.
// Note that time.Time methods are not available, if needed you can add and cast like the String method does
// Otherwise, only use in the json struct at marshal/unmarshal time.
type Date struct {
	DateValue time.Time
	IsNotNull bool
}

const CtLayout = "2006-01-02"

func Time(ct *Date) time.Time {
	return time.Time(ct.DateValue)

}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}
func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date.DateValue).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {

	return (*time.Time)(&date.DateValue).GobDecode(b)
}

// UnmarshalJSON Parses the json string in the custom format
func (date *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	if s != "null" && s != "" {
		nt, err := time.Parse(CtLayout, s)
		if err != nil {
			return err
		}
		date.DateValue = nt
		date.IsNotNull = true

	} else {
		date.IsNotNull = false
	}

	return
}
func (date *Date) Scan(value interface{}) error {
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)

	if err == nil {
		date.DateValue = nullTime.Time
		date.IsNotNull = true
	}

	return err
}

func (date Date) Value() (driver.Value, error) {

	if date.IsNotNull {

		y, m, d := time.Time(date.DateValue).Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date.DateValue).Location()), nil
	} else {
		return nil, nil
	}
}

// MarshalJSON writes a quoted string in the custom format
func (date Date) MarshalJSON() ([]byte, error) {

	if date.IsNotNull {

		return []byte(date.String()), nil
	} else {
		return json.Marshal(nil)
	}

}

// String returns the time in the custom format
func (date *Date) String() string {

	if date.IsNotNull {
		t := time.Time(date.DateValue)
		return fmt.Sprintf("%q", t.Format(CtLayout))
	} else {
		return ""
	}

	//return t.Format(CtLayout)
}

// String returns the time in the custom format
func (date *Date) StringValua() string {
	if !date.IsNotNull {
		return ""
	} else {
		t := time.Time(date.DateValue)

		return t.Format(CtLayout)
	}
}
