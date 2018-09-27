// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package machinelearningiface provides an interface to enable mocking the Amazon Machine Learning service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package machinelearningiface

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/machinelearning"
)

// MachineLearningAPI provides an interface to enable mocking the
// machinelearning.MachineLearning service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // Amazon Machine Learning.
//    func myFunc(svc machinelearningiface.MachineLearningAPI) bool {
//        // Make svc.AddTags request
//    }
//
//    func main() {
//        cfg, err := external.LoadDefaultAWSConfig()
//        if err != nil {
//            panic("failed to load config, " + err.Error())
//        }
//
//        svc := machinelearning.New(cfg)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mockMachineLearningClient struct {
//        machinelearningiface.MachineLearningAPI
//    }
//    func (m *mockMachineLearningClient) AddTags(input *machinelearning.AddTagsInput) (*machinelearning.AddTagsOutput, error) {
//        // mock response/functionality
//    }
//
//    func TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mockMachineLearningClient{}
//
//        myfunc(mockSvc)
//
//        // Verify myFunc's functionality
//    }
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters. Its suggested to use the pattern above for testing, or using
// tooling to generate mocks to satisfy the interfaces.
type MachineLearningAPI interface {
	AddTagsRequest(*machinelearning.AddTagsInput) machinelearning.AddTagsRequest

	CreateBatchPredictionRequest(*machinelearning.CreateBatchPredictionInput) machinelearning.CreateBatchPredictionRequest

	CreateDataSourceFromRDSRequest(*machinelearning.CreateDataSourceFromRDSInput) machinelearning.CreateDataSourceFromRDSRequest

	CreateDataSourceFromRedshiftRequest(*machinelearning.CreateDataSourceFromRedshiftInput) machinelearning.CreateDataSourceFromRedshiftRequest

	CreateDataSourceFromS3Request(*machinelearning.CreateDataSourceFromS3Input) machinelearning.CreateDataSourceFromS3Request

	CreateEvaluationRequest(*machinelearning.CreateEvaluationInput) machinelearning.CreateEvaluationRequest

	CreateMLModelRequest(*machinelearning.CreateMLModelInput) machinelearning.CreateMLModelRequest

	CreateRealtimeEndpointRequest(*machinelearning.CreateRealtimeEndpointInput) machinelearning.CreateRealtimeEndpointRequest

	DeleteBatchPredictionRequest(*machinelearning.DeleteBatchPredictionInput) machinelearning.DeleteBatchPredictionRequest

	DeleteDataSourceRequest(*machinelearning.DeleteDataSourceInput) machinelearning.DeleteDataSourceRequest

	DeleteEvaluationRequest(*machinelearning.DeleteEvaluationInput) machinelearning.DeleteEvaluationRequest

	DeleteMLModelRequest(*machinelearning.DeleteMLModelInput) machinelearning.DeleteMLModelRequest

	DeleteRealtimeEndpointRequest(*machinelearning.DeleteRealtimeEndpointInput) machinelearning.DeleteRealtimeEndpointRequest

	DeleteTagsRequest(*machinelearning.DeleteTagsInput) machinelearning.DeleteTagsRequest

	DescribeBatchPredictionsRequest(*machinelearning.DescribeBatchPredictionsInput) machinelearning.DescribeBatchPredictionsRequest

	DescribeDataSourcesRequest(*machinelearning.DescribeDataSourcesInput) machinelearning.DescribeDataSourcesRequest

	DescribeEvaluationsRequest(*machinelearning.DescribeEvaluationsInput) machinelearning.DescribeEvaluationsRequest

	DescribeMLModelsRequest(*machinelearning.DescribeMLModelsInput) machinelearning.DescribeMLModelsRequest

	DescribeTagsRequest(*machinelearning.DescribeTagsInput) machinelearning.DescribeTagsRequest

	GetBatchPredictionRequest(*machinelearning.GetBatchPredictionInput) machinelearning.GetBatchPredictionRequest

	GetDataSourceRequest(*machinelearning.GetDataSourceInput) machinelearning.GetDataSourceRequest

	GetEvaluationRequest(*machinelearning.GetEvaluationInput) machinelearning.GetEvaluationRequest

	GetMLModelRequest(*machinelearning.GetMLModelInput) machinelearning.GetMLModelRequest

	PredictRequest(*machinelearning.PredictInput) machinelearning.PredictRequest

	UpdateBatchPredictionRequest(*machinelearning.UpdateBatchPredictionInput) machinelearning.UpdateBatchPredictionRequest

	UpdateDataSourceRequest(*machinelearning.UpdateDataSourceInput) machinelearning.UpdateDataSourceRequest

	UpdateEvaluationRequest(*machinelearning.UpdateEvaluationInput) machinelearning.UpdateEvaluationRequest

	UpdateMLModelRequest(*machinelearning.UpdateMLModelInput) machinelearning.UpdateMLModelRequest

	WaitUntilBatchPredictionAvailable(*machinelearning.DescribeBatchPredictionsInput) error
	WaitUntilBatchPredictionAvailableWithContext(aws.Context, *machinelearning.DescribeBatchPredictionsInput, ...aws.WaiterOption) error

	WaitUntilDataSourceAvailable(*machinelearning.DescribeDataSourcesInput) error
	WaitUntilDataSourceAvailableWithContext(aws.Context, *machinelearning.DescribeDataSourcesInput, ...aws.WaiterOption) error

	WaitUntilEvaluationAvailable(*machinelearning.DescribeEvaluationsInput) error
	WaitUntilEvaluationAvailableWithContext(aws.Context, *machinelearning.DescribeEvaluationsInput, ...aws.WaiterOption) error

	WaitUntilMLModelAvailable(*machinelearning.DescribeMLModelsInput) error
	WaitUntilMLModelAvailableWithContext(aws.Context, *machinelearning.DescribeMLModelsInput, ...aws.WaiterOption) error
}

var _ MachineLearningAPI = (*machinelearning.MachineLearning)(nil)
