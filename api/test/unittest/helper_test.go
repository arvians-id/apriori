package unittest

import (
	"fmt"
	"github.com/arvians-id/apriori/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(t *testing.M) {
	fmt.Println("BEFORE: Begin Testing")

	t.Run()

	fmt.Println("AFTER: End Testing")
}

func TestConvertStrToInt(t *testing.T) {
	str := "34"
	result := util.StrToInt(str)
	assert.Equal(t, 34, result, "they should be equal integer")
}

func TestConvertIntToStr(t *testing.T) {
	integer := 34
	result := util.IntToStr(integer)
	assert.Equal(t, "34", result, "they should be equal string")
}

func TestUpperWords(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		expected string
	}{
		{
			name:     "the letters not being uppercase",
			request:  "in with a",
			expected: "in with a",
		},
		{
			name:     "the letters should be uppercase",
			request:  "Hello World",
			expected: "Hello World",
		},
	}

	for _, value := range tests {
		t.Run(value.name, func(t *testing.T) {
			str := util.UpperWords(value.request)
			assert.Equal(t, value.expected, str)
		})
	}
}
