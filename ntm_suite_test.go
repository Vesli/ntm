package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vesli/ntm/config"

	"testing"
)

var (
	s        *service
	confTest *config.Config
	dbTest   *gorm.DB
)

func TestNtm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ntm Suite")
}

var _ = BeforeSuite(func() {
	var err error

	Describe("Testing config JSON", func() {
		Context("Wrong config file", func() {
			s, err = newService("co.json")
			Expect(s).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
		Context("Error on config file", func() {
			s, err = newService("config-fail-test.json")
			Expect(s).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	s, err = newService("config-test.json")
	Expect(err).NotTo(HaveOccurred())

	confTest = s.Config
	dbTest = s.DB
})

var _ = AfterSuite(func() {
	fmt.Println("Closing!")
	s.Close()
})
