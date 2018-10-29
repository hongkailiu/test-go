package unittest

import (
	"testing"

	"github.com/hongkailiu/test-go/proto_buffer/gen/proto"
	. "github.com/onsi/gomega"
)

func TestPersonConstructor(t *testing.T) {
	o := NewGomegaWithT(t)

	p := proto.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*proto.Person_PhoneNumber{
			{Number: "555-4321", Type: proto.Person_HOME},
		},
	}

	o.Expect(p).NotTo(BeNil())
}
