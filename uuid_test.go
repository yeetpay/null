package null

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

var (
	uuidZero            = uuid.UUID{}
	UUIDText            = "14070757-48e6-4b4b-9cd8-98fbf065cf31"
	UUIDStringJSON      = []byte(`"14070757-48e6-4b4b-9cd8-98fbf065cf31"`)
	UUIDBlankStringJSON = []byte(`"00000000-0000-0000-0000-000000000000"`)
	nullUUIDJSON        = []byte(`{"UUID":"14070757-48e6-4b4b-9cd8-98fbf065cf31","Valid":true}`)
	uuidTestValue, _    = uuid.Parse("14070757-48e6-4b4b-9cd8-98fbf065cf31")
)

func TestUUIDFrom(t *testing.T) {
	i := UUIDFrom(uuidTestValue)
	assertUUID(t, i, "UUIDFrom()")

	zero := UUIDFrom(uuid.UUID{})
	if !zero.Valid {
		t.Error("UUIDFrom(0)", "is invalid, but should be valid")
	}
}

func TestUUIDFromPtr(t *testing.T) {
	ptr := &uuidTestValue
	i := UUIDFromPtr(ptr)
	assertUUID(t, i, "UUIDFromPtr()")

	null := UUIDFromPtr(nil)
	assertNullUUID(t, null, "UUIDFromPtr(nil)")
}

func TestUnmarshalUUID(t *testing.T) {
	var i UUID
	err := json.Unmarshal(UUIDStringJSON, &i)
	maybePanic(err)
	assertUUID(t, i, "uuid json")

	var ni UUID
	err = json.Unmarshal(nullUUIDJSON, &ni)
	if err == nil {
		panic("err should not be nill")
	}

	var null UUID
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUUID(t, null, "null json")

	var badType UUID
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUUID(t, badType, "wrong type json")

	var invalid UUID
	err = invalid.UnmarshalJSON(invalidJSON)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUUID(t, invalid, "invalid json")
}

func TestTextUnmarshalUUID(t *testing.T) {
	var i UUID
	err := i.UnmarshalText([]byte("14070757-48e6-4b4b-9cd8-98fbf065cf31"))
	maybePanic(err)
	assertUUID(t, i, "UnmarshalText() uuid")

	var blank UUID
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUUID(t, blank, "UnmarshalText() empty int")

	var null UUID
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullUUID(t, null, `UnmarshalText() "null"`)

	var zeroValue UUID
	err = null.UnmarshalText([]byte("00000000-0000-0000-0000-000000000000"))
	maybePanic(err)
	assertNullUUID(t, zeroValue, `UnmarshalText() "null"`)

	var invalid UUID
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalUUID(t *testing.T) {
	i := UUIDFrom(uuidTestValue)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, string(UUIDStringJSON), "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUUID(uuidZero, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalUUIDText(t *testing.T) {
	i := UUIDFrom(uuidTestValue)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, UUIDText, "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUUID(uuidZero, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUUIDPointer(t *testing.T) {
	i := UUIDFrom(uuidTestValue)
	ptr := i.Ptr()
	if *ptr != uuidTestValue {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, uuidTestValue)
	}

	null := NewUUID(uuidZero, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUUIDIsZero(t *testing.T) {
	i := UUIDFrom(uuidTestValue)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUUID(uuidZero, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUUID(uuidZero, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUUIDSetValid(t *testing.T) {
	change := NewUUID(uuidZero, false)
	assertNullUUID(t, change, "SetValid()")
	change.SetValid(uuidTestValue)
	assertUUID(t, change, "SetValid()")
}

func TestUUIDScan(t *testing.T) {
	var u UUID
	err := u.Scan(UUIDText)
	maybePanic(err)
	assertUUID(t, u, "scanned uuid")

	var null UUID
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUUID(t, null, "scanned null")
}

func TestUUIDValueOrZero(t *testing.T) {
	valid := NewUUID(uuidTestValue, true)
	if valid.ValueOrZero() != uuidTestValue {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewUUID(uuidTestValue, false)
	if invalid.ValueOrZero() != uuidZero {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestUUIDEqual(t *testing.T) {
	uuid1 := NewUUID(uuidTestValue, false)
	uuid2 := NewUUID(uuidTestValue, false)
	assertUUIDEqualIsTrue(t, uuid1, uuid2)

	uuid1 = NewUUID(uuidTestValue, false)
	uuid2 = NewUUID(uuid.New(), false)
	assertUUIDEqualIsTrue(t, uuid1, uuid2)

	uuid1 = NewUUID(uuidTestValue, true)
	uuid2 = NewUUID(uuidTestValue, true)
	assertUUIDEqualIsTrue(t, uuid1, uuid2)

	uuid1 = NewUUID(uuidTestValue, true)
	uuid2 = NewUUID(uuidTestValue, false)
	assertUUIDEqualIsFalse(t, uuid1, uuid2)

	uuid1 = NewUUID(uuidTestValue, false)
	uuid2 = NewUUID(uuidTestValue, true)
	assertUUIDEqualIsFalse(t, uuid1, uuid2)

	uuid1 = NewUUID(uuidTestValue, true)
	uuid2 = NewUUID(uuid.New(), true)
	assertUUIDEqualIsFalse(t, uuid1, uuid2)
}

func assertUUID(t *testing.T, i UUID, from string) {
	if i.UUID.String() != uuidTestValue.String() {
		t.Errorf("bad %s uuid: %s ≠ %s\n", from, i.UUID.String(), uuidTestValue.String())
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUUID(t *testing.T, i UUID, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertUUIDEqualIsTrue(t *testing.T, a, b UUID) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of UUID{%v, Valid:%t} and UUID{%v, Valid:%t} should return true", a.UUID, a.Valid, b.UUID, b.Valid)
	}
}

func assertUUIDEqualIsFalse(t *testing.T, a, b UUID) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of UUID{%v, Valid:%t} and UUID{%v, Valid:%t} should return false", a.UUID, a.Valid, b.UUID, b.Valid)
	}
}
