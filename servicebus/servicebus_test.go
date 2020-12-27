package servicebus_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()

	defer func() {
		os.Exit(code)
	}()
}

func TestSchemaValidation(t *testing.T) {
	t.Run("Should be successful if schema is valid", func(t *testing.T) {
		t.Fail()
	})
}
