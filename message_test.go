package miraihttp_test

import (
	. "github.com/CuteReimu/mirai-sdk-http"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChain(t *testing.T) {
	assert.Equal(t, MessageChain(&Plain{Text: "123"}), []SingleMessage{&Plain{Type: "Plain", Text: "123"}})
	assert.Equal(t, MessageChain(&Poke{Name: "SixSixSix"}), []SingleMessage{&Poke{Type: "Poke", Name: "SixSixSix"}})
	assert.Equal(t, MessageChain(&Image{ImageId: "1", Url: "url"}), []SingleMessage{&Image{Type: "Image", ImageId: "1", Url: "url"}})
}
