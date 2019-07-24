package gql

import (
	"fmt"
	io "io"

	"github.com/vektah/gqlgen/graphql"
)

// nolint: vet, errcheck
func MarshalInt64Scalar(n int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		text := []byte(fmt.Sprintf(`%v`, n))
		w.Write(text)
	})
}

func UnmarshalInt64Scalar(v interface{}) (int64, error) {
	n, ok := v.(int64)
	if !ok {
		return 0, fmt.Errorf("Not int64")
	}
	return n, nil
}
