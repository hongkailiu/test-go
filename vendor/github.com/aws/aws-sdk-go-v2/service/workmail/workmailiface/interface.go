// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package workmailiface provides an interface to enable mocking the Amazon WorkMail service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package workmailiface

import (
	"github.com/aws/aws-sdk-go-v2/service/workmail"
)

// WorkMailAPI provides an interface to enable mocking the
// workmail.WorkMail service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // Amazon WorkMail.
//    func myFunc(svc workmailiface.WorkMailAPI) bool {
//        // Make svc.AssociateDelegateToResource request
//    }
//
//    func main() {
//        cfg, err := external.LoadDefaultAWSConfig()
//        if err != nil {
//            panic("failed to load config, " + err.Error())
//        }
//
//        svc := workmail.New(cfg)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mockWorkMailClient struct {
//        workmailiface.WorkMailAPI
//    }
//    func (m *mockWorkMailClient) AssociateDelegateToResource(input *workmail.AssociateDelegateToResourceInput) (*workmail.AssociateDelegateToResourceOutput, error) {
//        // mock response/functionality
//    }
//
//    func TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mockWorkMailClient{}
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
type WorkMailAPI interface {
	AssociateDelegateToResourceRequest(*workmail.AssociateDelegateToResourceInput) workmail.AssociateDelegateToResourceRequest

	AssociateMemberToGroupRequest(*workmail.AssociateMemberToGroupInput) workmail.AssociateMemberToGroupRequest

	CreateAliasRequest(*workmail.CreateAliasInput) workmail.CreateAliasRequest

	CreateGroupRequest(*workmail.CreateGroupInput) workmail.CreateGroupRequest

	CreateResourceRequest(*workmail.CreateResourceInput) workmail.CreateResourceRequest

	CreateUserRequest(*workmail.CreateUserInput) workmail.CreateUserRequest

	DeleteAliasRequest(*workmail.DeleteAliasInput) workmail.DeleteAliasRequest

	DeleteGroupRequest(*workmail.DeleteGroupInput) workmail.DeleteGroupRequest

	DeleteMailboxPermissionsRequest(*workmail.DeleteMailboxPermissionsInput) workmail.DeleteMailboxPermissionsRequest

	DeleteResourceRequest(*workmail.DeleteResourceInput) workmail.DeleteResourceRequest

	DeleteUserRequest(*workmail.DeleteUserInput) workmail.DeleteUserRequest

	DeregisterFromWorkMailRequest(*workmail.DeregisterFromWorkMailInput) workmail.DeregisterFromWorkMailRequest

	DescribeGroupRequest(*workmail.DescribeGroupInput) workmail.DescribeGroupRequest

	DescribeOrganizationRequest(*workmail.DescribeOrganizationInput) workmail.DescribeOrganizationRequest

	DescribeResourceRequest(*workmail.DescribeResourceInput) workmail.DescribeResourceRequest

	DescribeUserRequest(*workmail.DescribeUserInput) workmail.DescribeUserRequest

	DisassociateDelegateFromResourceRequest(*workmail.DisassociateDelegateFromResourceInput) workmail.DisassociateDelegateFromResourceRequest

	DisassociateMemberFromGroupRequest(*workmail.DisassociateMemberFromGroupInput) workmail.DisassociateMemberFromGroupRequest

	ListAliasesRequest(*workmail.ListAliasesInput) workmail.ListAliasesRequest

	ListGroupMembersRequest(*workmail.ListGroupMembersInput) workmail.ListGroupMembersRequest

	ListGroupsRequest(*workmail.ListGroupsInput) workmail.ListGroupsRequest

	ListMailboxPermissionsRequest(*workmail.ListMailboxPermissionsInput) workmail.ListMailboxPermissionsRequest

	ListOrganizationsRequest(*workmail.ListOrganizationsInput) workmail.ListOrganizationsRequest

	ListResourceDelegatesRequest(*workmail.ListResourceDelegatesInput) workmail.ListResourceDelegatesRequest

	ListResourcesRequest(*workmail.ListResourcesInput) workmail.ListResourcesRequest

	ListUsersRequest(*workmail.ListUsersInput) workmail.ListUsersRequest

	PutMailboxPermissionsRequest(*workmail.PutMailboxPermissionsInput) workmail.PutMailboxPermissionsRequest

	RegisterToWorkMailRequest(*workmail.RegisterToWorkMailInput) workmail.RegisterToWorkMailRequest

	ResetPasswordRequest(*workmail.ResetPasswordInput) workmail.ResetPasswordRequest

	UpdatePrimaryEmailAddressRequest(*workmail.UpdatePrimaryEmailAddressInput) workmail.UpdatePrimaryEmailAddressRequest

	UpdateResourceRequest(*workmail.UpdateResourceInput) workmail.UpdateResourceRequest
}

var _ WorkMailAPI = (*workmail.WorkMail)(nil)
