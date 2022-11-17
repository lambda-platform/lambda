package DB

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Blob []byte

// GormDataType gorm common data type
func (date Blob) GormDataType() string {
	return "BLOB"
}
func (b Blob) GobEncode() ([]byte, error) {
	return b, nil
}

func (blob *Blob) GobDecode(b []byte) error {
	*blob = b
	return nil
}

// UnmarshalJSON Parses the json string in the custom format
func (ct *Blob) UnmarshalJSON(b []byte) (err error) {
	*ct = b
	return
}
func (b *Blob) Scan(value interface{}) (err error) {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	*b = bytes
	return err
}

func (b Blob) Value() (driver.Value, error) {
	if len(b) == 0 {
		return nil, nil
	}
	fmt.Println(string(b))
	fmt.Println(string(b))
	fmt.Println(string(b))

	return "utl_raw.cast_to_raw(\"" + string(b) + "\")", nil
}

// MarshalJSON writes a quoted string in the custom format
func (ct Blob) MarshalJSON() ([]byte, error) {
	return ct, nil
}

// String returns the time in the custom format
func (ct *Blob) String() string {
	if ct != nil {
		return string(*ct)
	} else {
		return ""
	}
}
