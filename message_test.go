package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage(t *testing.T) {
	h := NewHeader(Test, "this is a test")
	buf := h.Marshal()
	h2 := Unmarshal(buf)

	assert.Equal(t, h.Body, h2.Body)
}
