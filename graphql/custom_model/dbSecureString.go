package custom_model

import (
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/lambda-platform/lambda/DB"
	"io"
	"strconv"
)

func MarshalDBSecureString(f DB.SecureString) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(string(f)))
	})
}

func UnmarshalDBSecureString(v interface{}) (DB.SecureString, error) {
	str, ok := v.(string)
	if !ok {
		return "", errors.New("EncryptedString must be a string")
	}

	return DB.SecureString(str), nil
}
