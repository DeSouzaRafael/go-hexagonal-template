package validator_test

import (
	"testing"

	webvalidator "github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/validator"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string `validate:"required"`
	Age  int    `validate:"min=0,max=150"`
}

func TestNew(t *testing.T) {
	v := webvalidator.New()
	assert.NotNil(t, v)
}

func TestValidate(t *testing.T) {
	v := webvalidator.New()

	t.Run("valid struct", func(t *testing.T) {
		err := v.Validate(&testStruct{Name: "Rafa", Age: 30})
		assert.NoError(t, err)
	})

	t.Run("missing required field", func(t *testing.T) {
		err := v.Validate(&testStruct{Age: 30})
		assert.Error(t, err)
	})

	t.Run("field out of range", func(t *testing.T) {
		err := v.Validate(&testStruct{Name: "Rafa", Age: 200})
		assert.Error(t, err)
	})
}
