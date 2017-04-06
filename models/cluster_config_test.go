// mystack-controller api
// +build unit
// https://github.com/topfreegames/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	. "github.com/topfreegames/mystack-controller/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	yaml1 = `
services:
  postgres:
    image: postgres:1.0
  redis:
    image: redis:1.0
apps:
  app1:
    image: app1
    port: 5000
    env:
      - name: DATABASE_URL
        value: postgres://derp:1234@example.com
  app2:
    image: app2
    port: 5001
`
)

var _ = Describe("ClusterConfig", func() {
	Describe("ParseYaml", func() {
		It("should build correct struct form yaml", func() {
			services, apps, err := ParseYaml(yaml1)
			Expect(err).NotTo(HaveOccurred())

			Expect(services["postgres"].Image).To(Equal("postgres:1.0"))
			Expect(services["redis"].Image).To(Equal("redis:1.0"))

			Expect(apps["app1"].Image).To(Equal("app1"))
			Expect(apps["app1"].Port).To(Equal(5000))
			Expect(apps["app1"].Environment).To(BeEquivalentTo([]*EnvVar{
				&EnvVar{
					Name:  "DATABASE_URL",
					Value: "postgres://derp:1234@example.com",
				},
			}))

			Expect(apps["app2"].Image).To(Equal("app2"))
			Expect(apps["app2"].Port).To(Equal(5001))
			Expect(apps["app2"].Environment).To(BeNil())
		})
	})
})
