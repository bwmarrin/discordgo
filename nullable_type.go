package discordgo

import (
	"encoding/json"
)

// Due to the complexity of Discord API rules in special cases,
// we need to handle the nullable types in a special way.
// This is the reason why we have the implemented a nullable type.
// The nullable type respect this interface (based on string implementation)
//
// If you needs to implement a nullable type, you should use this interface.
// with your type.
//
//   type NullableObject interface {
//     Elem() string
//     Ptr() *string
//     IsNil() bool
//     IsNull() bool
//     Null() NullableObject
//     MarshalJSON() ([]byte, error)
//     UnmarshalJSON(data []byte) error
//   }

// nullStringValue represents the value used to hard-set a NullableString
// to a "null" JSON Statement
const nullStringValue string = "null"

// NullableString is a nullable string
// With the complexity of new the Discord API version, it is necessary to use this type
// to represent nullable strings. This is because the Discord API needs to use null
// strings in certain situations to undefined values.
type NullableString struct {
	str *string
}

// NullString creates a NullableString from a string
func NullString(s string) NullableString {
	return NullableString{&s}
}

// Elem returns the string value assigned in NullableString
func (ns NullableString) Elem() string {
	return *ns.str
}

// Ptr returns a pointer to the string value assigned in NullableString
func (ns NullableString) Ptr() *string {
	return ns.str
}

// IsNil returns true if the string assigned is nil
func (ns NullableString) IsNil() bool {
	return ns.str == nil
}

// IsNull returns true if the string value assigned in equals to `"null"`
// /!\ When `str` is nil, this function will return false
func (ns NullableString) IsNull() bool {
	return !ns.IsNil() && ns.Elem() == nullStringValue
}

// Null sets the string value assigned to NullableString to null statement
func (ns NullableString) Null() NullableString {
	var nullString = nullStringValue
	ns.str = &nullString
	return ns
}

// MarshalJSON returns the serialized JSON of the NullableString
// if the string is equals to `null`, will override the serialization with `"null"`
func (ns *NullableString) MarshalJSON() ([]byte, error) {
	if ns.IsNull() {
		return []byte(nullStringValue), nil
	}

	return json.Marshal(ns.str)
}

// UnmarshalJSON extracts the string value from the JSON data
func (ns *NullableString) UnmarshalJSON(data []byte) error {
	if string(data) == nullStringValue {
		*ns = ns.Null()
		return nil
	}

	if err := json.Unmarshal(data, ns.str); err != nil {
		return err
	}

	return nil
}
