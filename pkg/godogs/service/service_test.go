package service

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
)

var (
	opt         = godog.Options{Output: colors.Colored(os.Stdout)}
	c           = context{loginService: NewService()}
	loginResult string
)

type context struct {
	loginService LoginService
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {

	s.Step(`^I user username "([^"]*)" and password "([^"]*)" as input of the service$`, iUserUsernameAndPasswordAsInputOfTheService)
	s.Step(`^the service returns "([^"]*)"$`, theServiceReturns)
}

func iUserUsernameAndPasswordAsInputOfTheService(username, password string) error {
	loginResult = c.loginService.Login(username, password)
	return nil
}

func theServiceReturns(arg1 string) error {
	if loginResult == arg1 {
		return nil
	} else {
		return fmt.Errorf("loginResult is not \"OK\": %s", loginResult)
	}

}
