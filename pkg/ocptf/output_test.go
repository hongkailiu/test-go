package ocptf_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/hongkailiu/test-go/pkg/ocptf"
	. "github.com/onsi/gomega"
)

func TestGetListOutput1(t *testing.T) {
	o := NewGomegaWithT(t)

	r, err := GetListOutput(nil, nil)

	o.Expect(err).To(BeNil())
	o.Expect(r).NotTo(BeNil())
	o.Expect(r.GroupMap).NotTo(BeNil())
	o.Expect(r.GroupMap).ShouldNot(BeEmpty())

	bytes, err := json.Marshal(r.GroupMap)
	o.Expect(err).To(BeNil())
	jsonString := fmt.Sprintf(string(bytes))
	fmt.Println(jsonString)
	o.Expect(jsonString).Should(And(ContainSubstring(UnderlineMetaKey), ContainSubstring(HostVarsKey)))
}

func TestGetListOutput2(t *testing.T) {
	o := NewGomegaWithT(t)

	groups := []Group{{Name: Master, Hosts: []string{"ec2-34-219-178-152.us-west-2.compute.amazonaws.com", "ec2-34-219-178-153.us-west-2.compute.amazonaws.com"}, Vars: map[string]interface{}{"var1": true}, Children: []string{}}}

	hosts := []Host{{Name: "ec2-34-219-178-152.us-west-2.compute.amazonaws.com", VarMap: map[string]interface{}{"var001": "value1", "var002": "value2"}},
		{Name: "ec2-34-219-178-153.us-west-2.compute.amazonaws.com", VarMap: map[string]interface{}{}}}

	r, err := GetListOutput(groups, hosts)

	o.Expect(err).To(BeNil())
	o.Expect(r).NotTo(BeNil())
	o.Expect(r.GroupMap).NotTo(BeNil())
	o.Expect(r.GroupMap).ShouldNot(BeEmpty())

	bytes, err := json.Marshal(r.GroupMap)
	o.Expect(err).To(BeNil())
	jsonString := fmt.Sprintf(string(bytes))
	fmt.Println(jsonString)
	o.Expect(jsonString).Should(And(ContainSubstring(UnderlineMetaKey), ContainSubstring(HostVarsKey)))
}
