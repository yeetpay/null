package null

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type (

	// UUID is a nullable UUID. It supports SQL and JSON serialization.
	// It will marshal to null if null. Blank string input will be considered null.
	UUID struct {
		uuid.UUID
		Valid bool
	}
)

// UUIDFrom creates a new UUID that will never be blank.
func UUIDFrom(u uuid.UUID) UUID {
	return NewUUID(u, true)
}

// StringFromPtr creates a new String that be null if s is nil.
func UUIDFromPtr(u *uuid.UUID) UUID {
	if u == nil {
		return NewUUID(uuid.UUID{}, false)
	}
	return NewUUID(*u, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (u UUID) ValueOrZero() uuid.UUID {
	if !u.Valid {
		return uuid.UUID{}
	}
	return u.UUID
}

// NewString creates a new String
func NewUUID(u uuid.UUID, valid bool) UUID {
	return UUID{UUID: u, Valid: valid}
}

// UnmarshalJSON implements json.Unmarshaler.
// It will unmarshal to a null UUID if the input is blank.
// It will return an error if the input is not an string, blank, or "null".
func (u *UUID) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "" || str == "null" || str == "00000000-0000-0000-0000-000000000000" {
		u.Valid = false
		return nil
	}

	var err error
	u.UUID, err = uuid.ParseBytes(data)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal uuid: %w", err)
	}
	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null UUID if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (u *UUID) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" || str == "00000000-0000-0000-0000-000000000000" {
		u.Valid = false
		return nil
	}

	var err error
	u.UUID, err = uuid.ParseBytes(text)

	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	u.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this UUID is null.
func (u UUID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(u.UUID)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this UUID is null.
func (u UUID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return u.UUID.MarshalText()
}

// SetValid changes this UUID's value and also sets it to be non-null.
func (u *UUID) SetValid(n uuid.UUID) {
	u.UUID = n
	u.Valid = true
}

// Ptr returns a pointer to this UUID's value, or a nil pointer if this UUID is null.
func (u *UUID) Ptr() *uuid.UUID {
	if !u.Valid {
		return nil
	}
	return &u.UUID
}

// IsZero returns true for invalid UUID, for future omitempty support (Go 1.4?)
func (u UUID) IsZero() bool {
	return !u.Valid
}

// Equal returns true if both UUID have the same value or are both null.
func (u UUID) Equal(other UUID) bool {
	return u.Valid == other.Valid &&
		(!u.Valid || u.UUID.String() == other.UUID.String())
}

func (u *UUID) Scan(src interface{}) error {
	if err := u.UUID.Scan(src); err != nil {
		return err
	}
	u.Valid = u.UUID != uuid.UUID{}
	return nil
}
