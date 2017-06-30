package main

import (
	"bytes"
	"testing"

	"github.com/altairsix/eventsource"
	"github.com/stretchr/testify/assert"
)

func TestEncoderDecoder(t *testing.T) {
	e1 := eventsource.Event(&A{ID: "abc"})
	e2 := eventsource.Event(&B{ID: "def"})

	buf := bytes.NewBuffer(nil)
	encoder := NewEncoder(buf)

	// Encode

	n, err := encoder.WriteEvent(e1)
	assert.Nil(t, err)
	assert.Equal(t, 17, n)

	n, err = encoder.WriteEvent(e2)
	assert.Nil(t, err)
	assert.Equal(t, 17, n)

	// Decode

	decoder := NewDecoder(buf)

	event, err := decoder.ReadEvent()
	assert.Nil(t, err)
	assert.Equal(t, e1, event)

	event, err = decoder.ReadEvent()
	assert.Nil(t, err)
	assert.Equal(t, e2, event)
}
