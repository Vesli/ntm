package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserModel", func() {
	u := &User{
		FBID:     "10155",
		Name:     "test name",
		Email:    "test@email.com",
		Birthday: "1990-04-03 00:00:00",
	}

	u2 := &User{
		FBID:     "20005",
		Name:     "name other",
		Email:    "test2@email.com",
		Birthday: "1995-05-03 00:00:00",
	}

	Describe("Inserting a new user", func() {
		It("Should not return an error", func() {
			err := dbTest.Create(&u).Error
			Expect(err).To(BeNil())
		})
	})

	Describe("Testing User methods", func() {
		Context("User already exists", func() {
			It("Should return true", func() {
				b := u.userAlreadyExists(dbTest)
				Expect(b).To(BeTrue())
			})
		})
		Context("User do not exists", func() {
			It("Should return true", func() {
				b := u2.userAlreadyExists(dbTest)
				Expect(b).To(BeFalse())
			})
		})

		Context("Update user", func() {
			It("Should not return an error", func() {
				u2.Name = "globog"
				err := u2.updateUserInDB(dbTest)
				Expect(err).To(BeNil())
			})
		})
	})
})
