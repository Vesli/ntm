package helper_test

import (
	"bytes"

	. "github.com/vesli/ntm/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResponseDecode", func() {
	type Token struct {
		AccessToken string `json:"id"`
	}

	Describe("Test Decode function", func() {
		var t *Token

		Context("with incorrect formated JSON", func() {
			It("Should return an error", func() {
				errD := DecodeBody(&t, bytes.NewReader([]byte(`{"hello: "ok"}`)))

				Ω(errD).ShouldNot(BeNil())
				Ω(t).Should(BeNil())
			})
		})

		Context("with incorrect JSON content", func() {
			It("Should not return an access token", func() {
				err := DecodeBody(&t, bytes.NewReader([]byte(`{"hello": "ok"}`)))

				Ω(err).Should(BeNil())
				Ω(t.AccessToken).Should(BeEmpty())
			})
		})
	})
})
