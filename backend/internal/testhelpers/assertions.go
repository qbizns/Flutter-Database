package testhelpers

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Assert provides test assertion helpers
type Assert struct {
	t *testing.T
}

// NewAssert creates a new assertion helper
func NewAssert(t *testing.T) *Assert {
	return &Assert{t: t}
}

// Equal asserts two values are equal
func (a *Assert) Equal(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		a.t.Errorf("assertion failed: values not equal\nExpected: %+v\nActual:   %+v\n%v",
			expected, actual, formatMessage(msgAndArgs...))
	}
}

// NotEqual asserts two values are not equal
func (a *Assert) NotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()

	if reflect.DeepEqual(expected, actual) {
		a.t.Errorf("assertion failed: values should not be equal\nValue: %+v\n%v",
			expected, formatMessage(msgAndArgs...))
	}
}

// Nil asserts value is nil
func (a *Assert) Nil(value interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()

	if !isNil(value) {
		a.t.Errorf("assertion failed: expected nil\nActual: %+v\n%v",
			value, formatMessage(msgAndArgs...))
	}
}

// NotNil asserts value is not nil
func (a *Assert) NotNil(value interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()

	if isNil(value) {
		a.t.Errorf("assertion failed: expected not nil\n%v",
			formatMessage(msgAndArgs...))
	}
}

// True asserts value is true
func (a *Assert) True(value bool, msgAndArgs ...interface{}) {
	a.t.Helper()

	if !value {
		a.t.Errorf("assertion failed: expected true\n%v",
			formatMessage(msgAndArgs...))
	}
}

// False asserts value is false
func (a *Assert) False(value bool, msgAndArgs ...interface{}) {
	a.t.Helper()

	if value {
		a.t.Errorf("assertion failed: expected false\n%v",
			formatMessage(msgAndArgs...))
	}
}

// NoError asserts error is nil
func (a *Assert) NoError(err error, msgAndArgs ...interface{}) {
	a.t.Helper()

	if err != nil {
		a.t.Errorf("assertion failed: expected no error\nError: %v\n%v",
			err, formatMessage(msgAndArgs...))
	}
}

// Error asserts error is not nil
func (a *Assert) Error(err error, msgAndArgs ...interface{}) {
	a.t.Helper()

	if err == nil {
		a.t.Errorf("assertion failed: expected error\n%v",
			formatMessage(msgAndArgs...))
	}
}

// ErrorContains asserts error message contains substring
func (a *Assert) ErrorContains(err error, contains string, msgAndArgs ...interface{}) {
	a.t.Helper()

	if err == nil {
		a.t.Errorf("assertion failed: expected error\n%v",
			formatMessage(msgAndArgs...))
		return
	}

	if !containsString(err.Error(), contains) {
		a.t.Errorf("assertion failed: error should contain %q\nError: %v\n%v",
			contains, err, formatMessage(msgAndArgs...))
	}
}

// Len asserts slice/map/string length
func (a *Assert) Len(obj interface{}, length int, msgAndArgs ...interface{}) {
	a.t.Helper()

	actualLen := getLen(obj)
	if actualLen != length {
		a.t.Errorf("assertion failed: incorrect length\nExpected: %d\nActual:   %d\n%v",
			length, actualLen, formatMessage(msgAndArgs...))
	}
}

// Empty asserts slice/map/string is empty
func (a *Assert) Empty(obj interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()

	if getLen(obj) != 0 {
		a.t.Errorf("assertion failed: expected empty\nActual: %+v\n%v",
			obj, formatMessage(msgAndArgs...))
	}
}

// NotEmpty asserts slice/map/string is not empty
func (a *Assert) NotEmpty(obj interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()

	if getLen(obj) == 0 {
		a.t.Errorf("assertion failed: expected not empty\n%v",
			formatMessage(msgAndArgs...))
	}
}

// Contains asserts slice/string contains element
func (a *Assert) Contains(container, element interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()

	if !contains(container, element) {
		a.t.Errorf("assertion failed: container should contain element\nContainer: %+v\nElement: %+v\n%v",
			container, element, formatMessage(msgAndArgs...))
	}
}

// NotContains asserts slice/string doesn't contain element
func (a *Assert) NotContains(container, element interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()

	if contains(container, element) {
		a.t.Errorf("assertion failed: container should not contain element\nContainer: %+v\nElement: %+v\n%v",
			container, element, formatMessage(msgAndArgs...))
	}
}

// ValidUUID asserts value is a valid UUID
func (a *Assert) ValidUUID(value uuid.UUID, msgAndArgs ...interface{}) {
	a.t.Helper()

	if value == uuid.Nil {
		a.t.Errorf("assertion failed: expected valid UUID\n%v",
			formatMessage(msgAndArgs...))
	}
}

// TimeApproxEqual asserts two times are approximately equal (within delta)
func (a *Assert) TimeApproxEqual(expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	a.t.Helper()

	diff := expected.Sub(actual)
	if diff < 0 {
		diff = -diff
	}

	if diff > delta {
		a.t.Errorf("assertion failed: times not approximately equal\nExpected: %v\nActual:   %v\nDelta:    %v\nDiff:     %v\n%v",
			expected, actual, delta, diff, formatMessage(msgAndArgs...))
	}
}

// Helper functions

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	kind := v.Kind()
	return (kind == reflect.Chan || kind == reflect.Func ||
		kind == reflect.Interface || kind == reflect.Map ||
		kind == reflect.Ptr || kind == reflect.Slice) && v.IsNil()
}

func getLen(obj interface{}) int {
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len()
	default:
		return 0
	}
}

func contains(container, element interface{}) bool {
	containerValue := reflect.ValueOf(container)

	switch containerValue.Kind() {
	case reflect.String:
		return containsString(container.(string), element.(string))
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			if reflect.DeepEqual(containerValue.Index(i).Interface(), element) {
				return true
			}
		}
	}

	return false
}

func formatMessage(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	if len(msgAndArgs) == 1 {
		return msgAndArgs[0].(string)
	}
	return ""
}
