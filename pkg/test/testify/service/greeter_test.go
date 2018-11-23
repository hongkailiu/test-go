package service_test

import (
	"github.com/hongkailiu/test-go/pkg/test/mockery/service/mocks"
	"testing"

	"github.com/hongkailiu/test-go/pkg/test/testify/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dbMock struct {
	mock.Mock
}

func (d *dbMock) FetchMessage(lang string) (string, error) {
	args := d.Called(lang)
	return args.String(0), args.Error(1)
}

func (d *dbMock) FetchDefaultMessage() (string, error) {
	args := d.Called()
	return args.String(0), args.Error(1)
}

func TestMockMethodWithoutArgs(t *testing.T) {
	theDBMock := dbMock{}                                            // create the mock
	theDBMock.On("FetchDefaultMessage").Return("foofofofof", nil)    // mock the expectation
	g := service.Greeter{Database: &theDBMock, Lang: "en"}           // create Greeter object using mocked db
	assert.Equal(t, "Message is: foofofofof", g.GreetInDefaultMsg()) // assert what actual value that will come
	theDBMock.AssertNumberOfCalls(t, "FetchDefaultMessage", 1)       // we can assert how many times the mocked method will be called
	theDBMock.AssertExpectations(t)                                  // this method will ensure everything specified with On and Return was in fact called as expected
}

func TestMockMethodWithArgs(t *testing.T) {
	theDBMock := dbMock{}
	theDBMock.On("FetchMessage", "sg").Return("lah", nil) // if FetchMessage("sg") is called, then return "lah"
	g := service.Greeter{Database: &theDBMock, Lang: "sg"}
	assert.Equal(t, "Message is: lah", g.Greet())
	theDBMock.AssertExpectations(t)
}

func TestMockMethodWithArgsIgnoreArgs(t *testing.T) {
	theDBMock := dbMock{}
	theDBMock.On("FetchMessage", mock.Anything).Return("lah", nil) // if FetchMessage(...) is called with any argument, please also return lah
	g := service.Greeter{Database: &theDBMock, Lang: "in"}
	assert.Equal(t, "Message is: lah", g.Greet())
	theDBMock.AssertCalled(t, "FetchMessage", "in")
	theDBMock.AssertNotCalled(t, "FetchMessage", "ch")
	theDBMock.AssertExpectations(t)
	mock.AssertExpectationsForObjects(t, &theDBMock)
}

func TestMatchedBy(t *testing.T) {
	theDBMock := dbMock{}
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
}
