package topic_test

import (
	"github.com/scarabsoft/go-bus/internal/topic"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

func TestErrors(t *testing.T){
	t.Run("ErrDoesNotExists", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(topic.ErrDoesNotExists.Error(), is.EqualTo("topic does not exists"))
	})
	t.Run("ErrAlreadyExists", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(topic.ErrAlreadyExists.Error(), is.EqualTo("topic already exists"))
	})
	t.Run("ErrAlreadyClosed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(topic.ErrAlreadyClosed.Error(), is.EqualTo("topic already closed"))
	})
	t.Run("ErrAlreadySubscribed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(topic.ErrAlreadySubscribed.Error(), is.EqualTo("handler already subscribed"))
	})
}

