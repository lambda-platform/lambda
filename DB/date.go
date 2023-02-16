package DB

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Date provides a custom time type that can represent a date or datetime with a custom formatter.
type Date struct {
	DateValue time.Time
	IsNotNull bool
}

// Supported date and datetime layouts
const (
	DateLayout     = "2006-01-02"
	DatetimeLayout = "2006-01-02 15:04:05"
)

// Time returns the underlying time value of the Date type.
func (date *Date) Time() time.Time {
	return date.DateValue
}

// GormDataType returns the data type that GORM should use for the Date type.
func (date Date) GormDataType() string {
	return "date"
}

// GobEncode encodes the Date type to a binary format for use with the gob package.
func (date *Date) GobEncode() ([]byte, error) {
	return date.DateValue.GobEncode()
}

// GobDecode decodes the binary data from GobEncode into the Date type.
func (date *Date) GobDecode(b []byte) error {
	return date.DateValue.GobDecode(b)
}

// UnmarshalJSON parses a JSON string in the custom format and sets the Date value.
func (date *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	if s != "null" && s != "" {
		// Try parsing the string as a date first
		nt, err := time.Parse(DateLayout, s)
		if err != nil {
			// If parsing as a date failed, try parsing as a datetime
			nt, err = time.Parse(DatetimeLayout, s)
			if err != nil {
				return err
			}
		}
		date.DateValue = nt
		date.IsNotNull = true
	} else {
		date.IsNotNull = false
	}

	return nil
}

// Scan scans the value from the database and sets the Date value.
func (date *Date) Scan(value interface{}) error {
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)

	if err == nil {
		date.DateValue = nullTime.Time
		date.IsNotNull = true
	}

	return err
}

// Value returns the database driver value for the Date type.
func (date Date) Value() (driver.Value, error) {
	if date.IsNotNull {
		return date.DateValue, nil
		//return time.Time(date.DateValue).Format(DatetimeLayout), nil
	} else {
		return nil, nil
	}
}

// MarshalJSON writes a quoted string in the custom format.
func (date Date) MarshalJSON() ([]byte, error) {
	if date.IsNotNull {
		return []byte(date.String()), nil
	} else {
		return json.Marshal(nil)
	}
}

// String returns the Date value in the custom format.
func (date *Date) String() string {
	if date.IsNotNull {
		return fmt.Sprintf("%q", date.DateValue.Format(DatetimeLayout))
	} else {
		return ""
	}
}

// StringValue returns the Date value as a string in the custom format.
func (date *Date) StringValue() string {
	if date.IsNotNull {
		return date.DateValue.Format(DatetimeLayout)
	} else {
		return ""
	}
}
