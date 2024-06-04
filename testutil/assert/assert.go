package assert

import (
	"errors"
	"reflect"
	"testing"
)

func RequireNotNil(t testing.TB, got interface{}) {
	t.Helper()

	if got == nil {
		t.Fatal("expected non-nil value, but got nil")
	}
}

func ErrorType[T error](t testing.TB, got error) {
	t.Helper()

	var want T
	if !errors.As(got, &want) {
		t.Errorf("got error with type %v want %v", reflect.TypeOf(got), reflect.TypeOf(want))
	}
}

func RequireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("didn't want error but got %v", err)
	}
}

func Equal(t testing.TB, got, want interface{}) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
