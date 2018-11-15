// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package appsync

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restjson"
)

// AppSync provides the API operation methods for making requests to
// AWS AppSync. See this package's package overview docs
// for details on the service.
//
// AppSync methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type AppSync struct {
	*aws.Client
}

// Used for custom client initialization logic
var initClient func(*AppSync)

// Used for custom request initialization logic
var initRequest func(*AppSync, *aws.Request)

// Service information constants
const (
	ServiceName = "appsync"   // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the AppSync client with a config.
//
// Example:
//     // Create a AppSync client from just a config.
//     svc := appsync.New(myConfig)
func New(config aws.Config) *AppSync {
	var signingName string
	signingName = "appsync"
	signingRegion := config.Region

	svc := &AppSync{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2017-07-25",
				JSONVersion:   "1.1",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc)
	}

	return svc
}

// newRequest creates a new request for a AppSync operation and runs any
// custom request initialization.
func (c *AppSync) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(c, req)
	}

	return req
}