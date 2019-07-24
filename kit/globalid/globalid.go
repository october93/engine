package globalid

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strings"

	uuid "github.com/gofrs/uuid"
)

var Nil = ID("")

type ID string
type IDs []ID

type Tuple struct {
	A ID
	B ID
}

// Generates a global identifier using pseudo-random numbers.
func Next() ID {
	id, _ := uuid.NewV4() // nolint: gas

	return ID(id.String())
}

// Validate checks if the given string is a valid UUID v4 and returns the error
// otherwise.
func Validate(s string) error {
	_, err := Parse(s)
	return err
}

// Valid returns true if the given string is a valid UUID v4.
func (id ID) Valid() bool {
	return Validate(string(id)) == nil
}

// Parses a given string to ID given that it is a valid UUID v4.
func Parse(s string) (ID, error) {
	ss := strings.ToLower(s)
	u, err := uuid.FromString(ss)
	if err != nil {
		return Nil, err
	}
	if u.Version() != 4 {
		return Nil, fmt.Errorf("Parsed UUID is not of version 4 but %d", u.Version())
	}
	return ID(ss), nil
}

//Index allows globalid to implement queue.Indexable
func (id ID) Index() ID {
	return id
}

// Returns canonical string representation of the global identifier.
// Example: 123e4567-e89b-12d3-a456-426655440000.
func (id ID) String() string {
	return strings.ToLower(string(id))
}

// Bytes returns bytes slice representation of the global identifier.
func (id ID) Bytes() []byte {
	return []byte(id)
}

func (id ID) MarshalText() (text []byte, err error) {
	uui, err := uuid.FromString(id.String())
	if err != nil {
		return nil, err
	}
	return uui.MarshalText()
}
func (id *ID) UnmarshalText(text []byte) (err error) {
	res, err := uuid.FromString(string(text))
	if err != nil {
		return err
	}
	*id = ID(res.String())
	return nil
}

// MarshalJSON returns the byte representation of the global identifier
// enclosed in double quotes to form valid JSON.
func (id ID) MarshalJSON() (text []byte, err error) {
	text = []byte(fmt.Sprintf(`"%v"`, id.String()))
	return
}

func (id *ID) UnmarshalGQL(v interface{}) error {
	*id = ID(v.(string))
	return nil
}

// nolint: errcheck, gas
func (id ID) MarshalGQL(w io.Writer) {
	text, err := id.MarshalJSON()
	if err != nil {
		return
	}
	_, err = w.Write(text)
	if err != nil {
		return
	}
}

// UnmarshalJSON is based on the underlying UnmarshalText implementation.
func (id *ID) UnmarshalJSON(b []byte) error {
	s := string(b)
	// handle leading and trailing double quotes for JSON serialization
	s = strings.Trim(s, `"`)
	if err := Validate(s); err != nil {
		return err
	}
	u, err := uuid.FromString(s)
	if err != nil {
		return err
	}
	err = u.UnmarshalText([]byte(s))
	if err != nil {
		return err
	}
	*id = ID(u.String())
	return nil
}

func (id ID) Value() (driver.Value, error) {
	if id == Nil {
		return nil, nil
	}
	res, err := uuid.FromString(string(id))
	if err != nil {
		return nil, err
	}
	return res.Value()
}

func (id *ID) Scan(src interface{}) error {
	if src == nil {
		*id = Nil
		return nil
	}
	var uid uuid.UUID
	if err := uid.Scan(src); err != nil {
		return err
	}
	*id = ID(uid.String())
	return nil
}

func (ids IDs) Value() (driver.Value, error) {
	sa := make([]string, len(ids))
	for i, id := range ids {
		sa[i] = string(id)
	}
	return driver.Value(fmt.Sprintf("{%s}", strings.Join(sa, ","))), nil
}

func (ids *IDs) Scan(src interface{}) error {
	toStringArray := func(s string) []string {
		r := strings.Trim(s, "{}")
		return strings.Split(r, ",")
	}

	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	sa := toStringArray(string(b))
	lids := make(IDs, len(sa))
	for i, s := range sa {
		res, err := uuid.FromString(s)
		if err != nil {
			return err
		}
		lids[i] = ID(res.String())
	}
	(*ids) = lids
	return nil
}
