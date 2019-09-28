package unittest

import (
	"fmt"
	proto2 "github.com/hongkailiu/test-go/pkg/experimental/probuf/gen/proto"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	. "github.com/onsi/gomega"
)

func TestPersonConstructor(t *testing.T) {
	o := NewGomegaWithT(t)

	p := proto2.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*proto2.Person_PhoneNumber{
			{Number: "555-4321", Type: proto2.Person_HOME},
		},
	}

	o.Expect(p).NotTo(BeNil())
}

func TestPB(t *testing.T) {
	o := NewGomegaWithT(t)

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	fileName := tmpFile.Name()
	fmt.Printf("======fileName: %s\n", fileName)
	defer os.Remove(tmpFile.Name())
	o.Expect(err).To(BeNil())

	p1 := proto2.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*proto2.Person_PhoneNumber{
			{Number: "555-4321", Type: proto2.Person_HOME},
		},
	}

	p2 := proto2.Person{
		Id:    5678,
		Name:  "Jane Doe",
		Email: "jdoe@example.com",
		Phones: []*proto2.Person_PhoneNumber{
			{Number: "888-4321", Type: proto2.Person_HOME},
		},
	}

	persons := []*proto2.Person{&p1, &p2}

	book := &proto2.AddressBook{People: persons}
	// ...

	// Write the new address book back to disk.
	out, err := proto.Marshal(book)
	o.Expect(err).To(BeNil())
	err = ioutil.WriteFile(fileName, out, 0644)
	o.Expect(err).To(BeNil())

	in, err := ioutil.ReadFile(fileName)
	o.Expect(err).To(BeNil())
	resultAddressBook := &proto2.AddressBook{}
	err = proto.Unmarshal(in, resultAddressBook)
	o.Expect(err).To(BeNil())

	// because XXX_* fields could be different
	for _, p := range resultAddressBook.People {
		r := getPersonFromBook(p, book)
		o.Expect(r).NotTo(BeNil())
		o.Expect(p.Name).To(Equal(r.Name))
		o.Expect(p.Email).To(Equal(r.Email))
		o.Expect(len(p.Phones)).To(Equal(len(r.Phones)))

		// TODO: delve into phones
	}
}

func getPersonFromBook(p *proto2.Person, book *proto2.AddressBook) *proto2.Person {
	for _, r := range book.People {
		if p.Id == r.Id {
			return p
		}
	}
	return nil
}
