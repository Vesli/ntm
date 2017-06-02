package main

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("UserHandler", func() {
	Describe("Main UserHandler test", func() {
		var server *ghttp.Server

		BeforeEach(func() {
			server = ghttp.NewServer()

			testJSON := []byte(
				`{
					"id": "10152725637197590",
					"name": "Jay Cee",
					"birthday": "1990-04-03T18:25:43.511Z",
					"email": "test@mail.fr"
				}`,
			)
			server.RouteToHandler("GET", "/test-user", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, testJSON),
			))

			server.RouteToHandler("GET", "/ntm-api", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, "Welcome to the NTM API!"),
			))

			server.RouteToHandler("POST", "/ntm-api/user/subscribe", registerAndLogginUser)
		})

		AfterEach(func() {
			server.Close()
		})

		Describe("Test Get on main entry", func() {
			Context("Test server", func() {
				It("Should have a normal behaviour", func() {
					response, err := http.Get(server.URL() + "/ntm-api")
					Ω(server.ReceivedRequests()).Should(Not(BeNil()))

					Ω(err).ShouldNot(HaveOccurred())
					Ω(response).ShouldNot(BeNil())
				})
			})
		})

		Describe("Test getUserFromToken", func() {
			t := &Token{}

			Context("With correct url", func() {
				It("Should succeed", func() {
					confTest.FBURL = server.URL() + confTest.FBURL
					u, err := getUserFromToken(t, confTest)
					Ω(err).Should(BeNil())
					Ω(u).ShouldNot(BeNil())
				})
			})

			Context("With uncorrect url", func() {
				It("Should fail", func() {
					confTest.FBURL = server.URL() + "ok"
					u, err := getUserFromToken(t, confTest)
					Ω(err).ShouldNot(BeNil())
					Ω(u).Should(BeNil())
				})
			})
		})
	})
})
