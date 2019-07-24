package globalid

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestCaseInsensitive(t *testing.T) {
	id, err := Parse("CA4CB7D0-0C16-4F2A-BE97-960600B1FF6B")
	if err != nil {
		t.Fatal(err)
	}
	exp := "ca4cb7d0-0c16-4f2a-be97-960600b1ff6b"
	if id.String() != exp {
		t.Fatalf("Parse(): expected parsed string %s, actual: %s", exp, id.String())
	}
}

func TestNext(t *testing.T) {
	id := Next()
	if id == Nil {
		t.Errorf("Next(): expected non-nil UUID, actual: %v", id)
	}
}

func TestNextCollision(t *testing.T) {
	next1 := Next()
	next2 := Next()
	if next1 == next2 {
		t.Errorf("Next(): expected to never generate the same UUID twice")
	}
}

func TestNextReturnsUUIDV4(t *testing.T) {
	id := Next()
	u, err := uuid.FromString(string(id))
	if err != nil {
		t.Errorf("Next(): unexpected error: %v", err)
	}
	actual := u.Version()
	if actual != 4 {
		t.Errorf("Next(): expected %v to be UUID v4, actual: %v", id, actual)
	}
}

func TestValidate(t *testing.T) {
	var validateTests = []struct {
		n   string
		err bool
	}{
		{
			"3c9ba730-03ff-4912-9b16-cdcc1e82d33d",
			false,
		},
		{
			"abc",
			true,
		},
		{
			"",
			true,
		},
		{
			"00000000-0000-0000-0000-000000000000",
			true,
		},
		{
			"78d7c818-f63b-11e6-bc64-92361f002671", // UUID v1
			true,
		},
		{
			"6fa459ea-ee8a-3ca4-894e-db77e160355e", // UUID v3
			true,
		},
		{
			"90691cbc-b5ea-5826-ae98-951e30fc3b2d", // UUID v5
			true,
		},
	}

	for _, tt := range validateTests {
		err := Validate(tt.n)
		if tt.err && err == nil {
			t.Errorf("Validate(%v): expected an error but is %v", tt.n, err)
		}
		if !tt.err && err != nil {
			t.Errorf("Validate(%v): expected no error but is %v", tt.n, err)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	s := "3c9ba730-03ff-4912-9b16-cdcc1e82d33d"
	id, err := Parse(s)
	if err != nil {
		t.Errorf("MarshalJSON(): unexpected error: %v", err)
	}
	actual, err := id.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON(): unexpected error: %v", err)
	}
	expected := `"3c9ba730-03ff-4912-9b16-cdcc1e82d33d"`
	if string(actual) != expected {
		t.Errorf("MarshalJSON(): expected %vm actual: %v", expected, actual)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var unmarshalJSONTests = []struct {
		name     string
		n        []byte
		expected ID
	}{
		{
			name:     "valid",
			n:        []byte("3c9ba730-03ff-4912-9b16-cdcc1e82d33d"),
			expected: ID("3c9ba730-03ff-4912-9b16-cdcc1e82d33d"),
		},
		{
			name:     "invalid",
			n:        []byte("abc"),
			expected: Nil,
		},
		{
			name:     "empty",
			n:        []byte(""),
			expected: Nil,
		},
		{
			name:     "null UUID",
			n:        []byte("00000000-0000-0000-0000-000000000000"),
			expected: Nil,
		},
		{
			name:     "UUID v1",
			n:        []byte("78d7c818-f63b-11e6-bc64-92361f002671"), // UUID v1,
			expected: Nil,
		},
		{
			name:     "UUID v3",
			n:        []byte("6fa459ea-ee8a-3ca4-894e-db77e160355e"), // UUID v3,
			expected: Nil,
		},

		{
			name:     "UUID v5",
			n:        []byte("90691cbc-b5ea-5826-ae98-951e30fc3b2d"), // UUID v5,
			expected: Nil,
		},
		{
			name:     "JSON",
			n:        []byte(`"3c9ba730-03ff-4912-9b16-cdcc1e82d33d"`), // string enclosed with double quotes
			expected: ID("3c9ba730-03ff-4912-9b16-cdcc1e82d33d"),
		},
	}

	for _, tt := range unmarshalJSONTests {
		tt := tt // capture range variable for parallell testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var id ID
			err := id.UnmarshalJSON(tt.n)
			if tt.expected != Nil && err != nil {
				t.Fatalf("UnmarshalJSON(%v): unexpected error: %v", string(tt.n), err)
			}
			if id != tt.expected {
				t.Errorf("UnmarshalJSON(%v): expected %v, actual: %v", string(tt.n), tt.expected, id)
			}
		})
	}
}

func TestScan(t *testing.T) {
	var scanTests = []struct {
		name     string
		src      interface{}
		expected ID
	}{
		{
			name:     "nil",
			src:      nil,
			expected: Nil,
		},
		{
			name:     "valid",
			src:      ID("3c9ba730-03ff-4912-9b16-cdcc1e82d33d").Bytes(),
			expected: ID("3c9ba730-03ff-4912-9b16-cdcc1e82d33d"),
		},
	}
	for _, tt := range scanTests {
		tt := tt // capture range variable for parallell testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := Next()
			err := id.Scan(tt.src)
			if err != nil {
				t.Fatalf("Scan(): unexpected error: %v", err)
			}
			if id != tt.expected {
				t.Errorf("Scan(): unexpected result %v, expected %v", id, tt.expected)
			}
		})
	}
}

func BenchmarkNext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Next()
	}
}
