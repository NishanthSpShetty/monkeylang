package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	assert.Equal(t, hello1.HashKey(), hello2.HashKey(), "strings with same content have different hash keys")
	assert.Equal(t, diff1.HashKey(), diff2.HashKey(), "strings with same content have different hash keys")
	assert.NotEqual(t, hello1.HashKey(), diff1.HashKey(), "strings with different content have same hash keys")

}
