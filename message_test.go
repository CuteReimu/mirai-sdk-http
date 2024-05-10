package miraihttp

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

func TestMessage(t *testing.T) {
	assert.Equal(t, MessageChain{&Plain{Type: "Plain", Text: "123"}}, buildMessageChain(MessageChain{&Plain{Text: "123"}}))
	assert.Equal(t, MessageChain{&Poke{Type: "Poke", Name: "SixSixSix"}}, buildMessageChain(MessageChain{&Poke{Name: "SixSixSix"}}))
	assert.Equal(t, MessageChain{&Image{Type: "Image", ImageId: "1", Url: "url"}}, buildMessageChain(MessageChain{&Image{ImageId: "1", Url: "url"}}))
}

func TestMessageChain(t *testing.T) {
	content := `[{"type":"Plain","text":"123"},{"type":"Poke","name":"SixSixSix"},{"type":"Image","imageId":"1","url":"url"}]`
	assert.Equal(t, parseMessageChain(gjson.Parse(content).Array()), buildMessageChain(
		MessageChain{&Plain{Text: "123"}, &Poke{Name: "SixSixSix"}, &Image{ImageId: "1", Url: "url"}},
	))
}
