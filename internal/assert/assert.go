package assert

import (
	"errors"
	"reflect"
	"testing"
)

func NoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func ErrorIs(t testing.TB, err, target error) {
	t.Helper()

	switch {
	case err == nil:
		t.Fatalf("expected error to be %v, got nil", target)
	case !errors.Is(err, target):
		t.Fatalf("expected error to be %v, got %v", target, err)
	}
}

func Equal[T any](t testing.TB, expected, actual T) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}
