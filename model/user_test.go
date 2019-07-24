package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestHashPassword(t *testing.T) {
	var hashPasswordTests = []struct {
		password string
	}{
		{"October: Secret Society"},
		{""},
	}
	for _, tt := range hashPasswordTests {
		hash, err := HashPassword(tt.password)
		if err != nil {
			t.Fatal(err)
		}
		if len(hash) != 80 {
			t.Errorf("HashPassword(%s): key length has deviated from 80 bytes: %d bytes", tt.password, len(hash))
		}
	}
}

func TestPasswordMatches(t *testing.T) {
	var passwordMatchesTests = []struct {
		name     string
		password string
		input    string
		expected bool
	}{
		{
			name:     "valid",
			password: "secret",
			input:    "secret",
			expected: true,
		},
		{
			name:     "incorrect password input",
			password: "secret",
			input:    "love",
			expected: false,
		},
		{
			name:     "empty password input",
			password: "secret",
			input:    "",
			expected: false,
		},
		{
			name:     "empty password and input",
			password: "",
			input:    "",
			expected: true,
		},
	}
	for _, tt := range passwordMatchesTests {
		tt := tt // capture range variable for parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := User{}
			u.SetPassword(tt.password)
			actual, err := u.PasswordMatches(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if actual != tt.expected {
				t.Errorf("PasswordMatches(): expected %t, actual: %t", tt.expected, actual)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	var validateUsernameTests = []struct {
		name      string
		username  string
		blacklist map[string]bool
		expected  error
	}{
		{
			name:     "valid username",
			username: "gopher",
			expected: nil,
		},
		{
			name:     "empty username",
			username: "",
			expected: errors.New("username is too short"),
		},
		{
			name:     "username too long",
			username: "franzkafka123456789999999",
			expected: errors.New("username is too long"),
		},
		{
			name:     "username with 20 characters",
			username: "_commanderjamespike_",
			expected: nil,
		},
		{
			name:     "username with 2 character",
			username: "aa",
			expected: nil,
		},
		{
			name:     "username with underscores",
			username: "___",
			expected: nil,
		},
		{
			name:     "username with only numbers",
			username: "12345",
			expected: nil,
		},
		{
			name:     "username with invalid characters",
			username: "franzkafka!$",
			expected: errors.New("username can only use letters, numbers and underscores"),
		},
		{
			name:      "blacklisted username",
			username:  "redrocket",
			blacklist: map[string]bool{"redrocket": true},
			expected:  errors.New("username is unavailable"),
		},
	}
	for _, tt := range validateUsernameTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			blacklist := tt.blacklist
			if blacklist == nil {
				blacklist = make(map[string]bool)
			}
			err := ValidateUsername(tt.username, blacklist)
			if tt.expected != err && (tt.expected == nil || err.Error() != tt.expected.Error()) {
				t.Errorf("ValidateUsername(%s): expected %v, actual: %v", tt.username, tt.expected, err)
			}
		})
	}
}

func TestNeverMarshalPassword(t *testing.T) {
	var neverMarshalPasswordTests = []struct {
		name     string
		user     *User
		b        []byte
		expected string
	}{
		{
			name:     "empty user",
			user:     &User{},
			expected: `{"displayname":"","profileimg_path":"","cover_image_path":"","userBio":"","username":"","botchedSignup":false}`,
		},
		{
			name:     "user with password",
			user:     &User{PasswordHash: string("55IGNhcm5hbCBwbGVhc3VyZS4=")},
			expected: `{"displayname":"","profileimg_path":"","cover_image_path":"","userBio":"","username":"","botchedSignup":false}`,
		},
	}
	for _, tt := range neverMarshalPasswordTests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.user.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON(%s): unexpected error: %v", tt.name, err)
			}
			if !reflect.DeepEqual(tt.expected, string(result)) {
				t.Errorf("MarshalJSON(%s): expected %v, actual: %v", tt.name, tt.expected, string(result))
			}
		})
	}
}

func TestMarshalJSONExportedUser(t *testing.T) {
	var marshalJSONExportedUser = []struct {
		name     string
		user     *ExportedUser
		b        []byte
		expected string
	}{
		{
			name:     "empty user",
			user:     &ExportedUser{},
			expected: `{"displayname":"","profileimg_path":"","cover_image_path":"","userBio":"","username":"","botchedSignup":false}`,
		},
		{
			name:     "admin user",
			user:     &ExportedUser{Admin: true},
			expected: `{"displayname":"","profileimg_path":"","cover_image_path":"","userBio":"","username":"","admin":true,"botchedSignup":false}`,
		},
	}
	for _, tt := range marshalJSONExportedUser {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := json.Marshal(tt.user)
			if err != nil {
				t.Fatalf("MarshalJSON(%s): unexpected error: %v", tt.name, err)
			}
			if !reflect.DeepEqual(tt.expected, string(result)) {
				t.Errorf("MarshalJSON(%s): expected %v, actual: %v", tt.name, tt.expected, string(result))
			}
		})
	}
}

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashPassword(fmt.Sprintf("password-%d", i+1))
	}
}
