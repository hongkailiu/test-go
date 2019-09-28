package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hongkailiu/test-go/pkg/experimental/test/mockery/service"
	"github.com/hongkailiu/test-go/pkg/experimental/test/mockery/service/mocks"
)

func TestMockMethodWithoutArgs(t *testing.T) {
	theDBMock := mocks.DB{}
	theDBMock.On("FetchDefaultMessage").Return("foofofofof", nil)    // mock the expectation
	g := service.Greeter{Database: &theDBMock, Lang: "en"}           // create Greeter object using mocked db
	assert.Equal(t, "Message is: foofofofof", g.GreetInDefaultMsg()) // assert what actual value that will come
	theDBMock.AssertNumberOfCalls(t, "FetchDefaultMessage", 1)       // we can assert how many times the mocked method will be called
	theDBMock.AssertExpectations(t)                                  // this method will ensure everything specified with On and Return was in fact called as expected
}

func TestMockMethodWithArgs(t *testing.T) {
	theDBMock := mocks.DB{}
	theDBMock.On("FetchMessage", "sg").Return("lah", nil) // if FetchMessage("sg") is called, then return "lah"
	g := service.Greeter{Database: &theDBMock, Lang: "sg"}
	assert.Equal(t, "Message is: lah", g.Greet())
	theDBMock.AssertExpectations(t)
}

func TestMockMethodWithArgsIgnoreArgs(t *testing.T) {
	theDBMock := mocks.DB{}
	theDBMock.On("FetchMessage", mock.Anything).Return("lah", nil) // if FetchMessage(...) is called with any argument, please also return lah
	g := service.Greeter{Database: &theDBMock, Lang: "in"}
	assert.Equal(t, "Message is: lah", g.Greet())
	theDBMock.AssertCalled(t, "FetchMessage", "in")
	theDBMock.AssertNotCalled(t, "FetchMessage", "ch")
	theDBMock.AssertExpectations(t)
	mock.AssertExpectationsForObjects(t, &theDBMock)
}

func TestMatchedBy(t *testing.T) {
	theDBMock := mocks.DB{}
	theDBMock.On("FetchMessage", mock.MatchedBy(func(lang string) bool { return lang[0] == 'i' })).Return("bzzzz", nil) // all of these call FetchMessage("iii"), FetchMessage("i"), FetchMessage("in") will match
	g := service.Greeter{Database: &theDBMock, Lang: "izz"}
	msg := g.Greet()
	assert.Equal(t, "Message is: bzzzz", msg)
	theDBMock.AssertExpectations(t)
}

// for test coverage only
func TestMatchedBy2(t *testing.T) {
	theDBMock := mocks.DB{}
	theDBMock.On("FetchMessage", mock.MatchedBy(func(lang string) bool { return lang[0] == 'i' })).Return("bzzzz", nil) // all of these call FetchMessage("iii"), FetchMessage("i"), FetchMessage("in") will match
	g := service.NewGreeter(&theDBMock, "izz")
	msg := g.Greet()
	assert.Equal(t, "Message is: bzzzz", msg)
	theDBMock.AssertExpectations(t)
}

func TestNewDB(t *testing.T) {
	r := service.NewDB()
	assert.NotNil(t, r)

	msg, err := r.FetchDefaultMessage()
	assert.Nil(t, err)
	assert.Equal(t, "default message", msg)

	msg, err = r.FetchMessage("en")
	assert.Nil(t, err)
	assert.Equal(t, "hello", msg)

	msg, err = r.FetchMessage("es")
	assert.Nil(t, err)
	assert.Equal(t, "holla", msg)

	msg, err = r.FetchMessage("un")
	assert.Nil(t, err)
	assert.Equal(t, "bzzzz", msg)
}
