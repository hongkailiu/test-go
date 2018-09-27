// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package resourcegroups

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restjson"
)

// ResourceGroups provides the API operation methods for making requests to
// AWS Resource Groups. See this package's package overview docs
// for details on the service.
//
// ResourceGroups methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type ResourceGroups struct {
	*aws.Client
}

// Used for custom client initialization logic
var initClient func(*ResourceGroups)

// Used for custom request initialization logic
var initRequest func(*ResourceGroups, *aws.Request)

// Service information constants
const (
	ServiceName = "resource-groups" // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName       // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the ResourceGroups client with a config.
//
// Example:
//     // Create a ResourceGroups client from just a config.
//     svc := resourcegroups.New(myConfig)
func New(config aws.Config) *ResourceGroups {
	var signingName string
	signingName = "resource-groups"
	signingRegion := config.Region

	svc := &ResourceGroups{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2017-11-27",
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

// newRequest creates a new request for a ResourceGroups operation and runs any
// custom request initialization.
func (c *ResourceGroups) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(c, req)
	}

	return req
}
