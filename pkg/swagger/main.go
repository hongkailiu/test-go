// https://ops.tips/blog/a-swagger-golang-hello-world/
// main declares the CLI that spins up the server of
// our API.
// It takes some arguments, validates if they're valid
// and match the expected type and then intiialize the
// server.
package main

import (
	models2 "github.com/hongkailiu/test-go/pkg/experimental/swagger/swagger/models"
	restapi2 "github.com/hongkailiu/test-go/pkg/experimental/swagger/swagger/restapi"
	operations2 "github.com/hongkailiu/test-go/pkg/experimental/swagger/swagger/restapi/operations"
	"log"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

type cliArgs struct {
	Port int `arg:"-p,help:port to listen to"`
}

var (
	args = &cliArgs{
		Port: 8080,
	}
)

// getHostnameHandler implements the handler that
// takes a set of parameters as described in swagger.yml
// and then produces a response.
// This response might be an error of a succesfull response.
//	-	In case of failure we create the payload
//		that would indicate the failure.
//	-	In case of success, the payload the indicates
//		the success with the hostname.
func getHostnameHandler(params operations2.GetHostnameParams) middleware.Responder {
	payload, err := os.Hostname()

	if err != nil {
		errPayload := &models2.Error{
			Code:    500,
			Message: swag.String("failed to retrieve hostname"),
		}

		return operations2.NewGetHostnameDefault(500).
			WithPayload(errPayload)
	}

	return operations2.NewGetHostnameOK().WithPayload(payload)
}

func getUsersHandler(params operations2.GetUsersParams) middleware.Responder {
	id := int64(3)
	users := []*models2.User{{&id, "hongkliu"}}
	return operations2.NewGetUsersOK().WithPayload(users)
}

func getUserUserIDHandler(params operations2.GetUserUserIDParams) middleware.Responder {

	if params.UserID == 1 {
		id := int64(1)
		user := models2.User{&id, "mike"}
		return operations2.NewGetUserUserIDOK().WithPayload(&user)
	}
	return operations2.NewGetUserUserIDNotFound()

}

// main performs the main routine of the application:
//	1.	parses the args;
//	2.	analyzes the declaration of the API
//	3.	sets the implementation of the handlers
//	4.	listens on the port we want
func main() {
	arg.MustParse(args)

	// Load the JSON that corresponds to our swagger.yml
	// api definition.
	// This JSON is hardcoded as part of the generated code
	// that go-swagger creates.
	swaggerSpec, err := loads.Analyzed(restapi2.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// Load a dummy object that servers as an interface
	// that allows us to implement the API specification.
	api := operations2.NewHelloAPI(swaggerSpec)

	// Create the REST api server that will make use of
	// the object that will container our handler implementations.
	server := restapi2.NewServer(api)
	defer server.Shutdown()

	// Configure the server port
	server.Port = args.Port

	// Add our handler implementation
	api.GetHostnameHandler = operations2.GetHostnameHandlerFunc(
		getHostnameHandler)

	api.GetUsersHandler = operations2.GetUsersHandlerFunc(getUsersHandler)
	api.GetUserUserIDHandler = operations2.GetUserUserIDHandlerFunc(getUserUserIDHandler)

	// Let it run
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
