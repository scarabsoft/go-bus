package topic_test

import (
	"github.com/scarabsoft/go-bus/internal/topic"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

const givenName = "givenName"

func TestErrors(t *testing.T) {
	t.Run("ErrorNotExists", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		testInstance := topic.ErrorNotExists{Name: givenName}
		assert.That(testInstance.Error(), is.EqualTo("givenName does not exists"))
	})
	t.Run("ErrAlreadyExists", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		testInstance := topic.ErrorAlreadyExists{Name: givenName}
		assert.That(testInstance.Error(), is.EqualTo("givenName already exists"))
	})
	t.Run("ErrAlreadyClosed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		testInstance := topic.ErrorAlreadyClosed{Name: givenName}
		assert.That(testInstance.Error(), is.EqualTo("givenName already closed"))
	})
	t.Run("ErrAlreadySubscribed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(topic.ErrAlreadySubscribed.Error(), is.EqualTo("handler already subscribed"))
	})
}
