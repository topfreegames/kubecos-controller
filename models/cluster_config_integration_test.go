// mystack-controller api
// +build integration
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

	mTest "github.com/topfreegames/mystack-controller/testing"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

var _ = Describe("ClusterConfig", func() {
	var (
		conn     runner.Connection
		db       *runner.Tx
		err      error
		services = map[string]*ClusterAppConfig{
			"postgres": &ClusterAppConfig{Image: "postgres:1.0"},
			"redis":    &ClusterAppConfig{Image: "redis:1.0"},
		}
		apps = map[string]*ClusterAppConfig{
			"app1": &ClusterAppConfig{
				Image: "app1",
				Port:  5000,
				Environment: []*EnvVar{
					&EnvVar{
						Name:  "DATABASE_URL",
						Value: "postgres://derp:1234@example.com",
					},
				},
			},
			"app2": &ClusterAppConfig{
				Image: "app2",
				Port:  5001,
			},
		}
	)

	BeforeSuite(func() {
		conn, err = mTest.GetTestDB()
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		db, err = conn.Begin()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err = db.Rollback()
		Expect(err).NotTo(HaveOccurred())
		db = nil
	})

	AfterSuite(func() {
		err = conn.(*runner.DB).DB.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("WriteClusterConfig", func() {
		It("should write cluster config", func() {
			err = WriteClusterConfig(db, "myCustomApps", apps, services)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("LoadClusterConfig", func() {
		It("should load cluster config", func() {
			err = WriteClusterConfig(db, "myCustomApps", apps, services)
			Expect(err).NotTo(HaveOccurred())

			returnApps, returnServices, err := LoadClusterConfig(db, "myCustomApps")
			Expect(err).NotTo(HaveOccurred())
			Expect(returnServices).To(BeEquivalentTo(services))
			Expect(returnApps).To(BeEquivalentTo(apps))
		})
	})
})
