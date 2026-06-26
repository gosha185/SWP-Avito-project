package storage_test

import (
	"regexp"
	"testing"

	"bonus-service/internal/validator"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Valid_Initially(t *testing.T) {
	v := validator.New()
	assert.True(t, v.Valid())
}

func TestValidator_Check_AddError(t *testing.T) {
	v := validator.New()
	v.Check(false, "field", "must not be empty")
	assert.False(t, v.Valid())
	assert.Equal(t, "must not be empty", v.Errors["field"])
}

func TestValidator_Check_NoErrorOnTrue(t *testing.T) {
	v := validator.New()
	v.Check(true, "field", "must not be empty")
	assert.True(t, v.Valid())
}

func TestValidator_AddError_DoesNotOverwrite(t *testing.T) {
	v := validator.New()
	v.AddError("field", "first error")
	v.AddError("field", "second error")
	assert.Equal(t, "first error", v.Errors["field"])
}

func TestValidator_MultipleErrors(t *testing.T) {
	v := validator.New()
	v.Check(false, "name", "must not be empty")
	v.Check(false, "age", "must be positive")
	assert.False(t, v.Valid())
	assert.Len(t, v.Errors, 2)
}

func TestIsPermitted(t *testing.T) {
	assert.True(t, validator.IsPermitted("a", "a", "b", "c"))
	assert.False(t, validator.IsPermitted("z", "a", "b", "c"))
	assert.True(t, validator.IsPermitted(42, 1, 42, 100))
}

func TestIsMatch(t *testing.T) {
	emailRx := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	assert.True(t, validator.IsMatch("user@example.com", emailRx))
	assert.False(t, validator.IsMatch("not-an-email", emailRx))
}

func TestIsUnique(t *testing.T) {
	assert.True(t, validator.IsUnique([]int{1, 2, 3}))
	assert.False(t, validator.IsUnique([]int{1, 2, 2}))
	assert.True(t, validator.IsUnique([]string{"a", "b", "c"}))
	assert.False(t, validator.IsUnique([]string{"a", "a"}))
}
