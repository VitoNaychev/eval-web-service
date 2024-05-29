package assert

import (
	"reflect"
	"testing"
)

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
