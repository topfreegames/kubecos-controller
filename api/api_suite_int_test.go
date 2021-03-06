// mystack-controller api
// +build integration
// https://github.com/topfreegames/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/topfreegames/mystack-controller/api"
	oTesting "github.com/topfreegames/mystack-controller/testing"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"testing"

	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var clientset kubernetes.Interface
var app *api.App
var conn *sqlx.DB
var db *sqlx.Tx
var config *viper.Viper

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Integration Suite")
}

var _ = BeforeSuite(func() {
	l := logrus.New()
	l.Level = logrus.FatalLevel

	var err error
	conn, err = oTesting.GetTestDB()
	Expect(err).NotTo(HaveOccurred())

	clientset := fake.NewSimpleClientset()

	config, err = oTesting.GetDefaultConfig()
	Expect(err).NotTo(HaveOccurred())
	app, err = api.NewApp("0.0.0.0", 8889, config, false, l, clientset)
	Expect(err).NotTo(HaveOccurred())
	app.DeploymentReadiness = &oTesting.MockReadiness{}
	app.JobReadiness = &oTesting.MockReadiness{}
})

var _ = BeforeEach(func() {
	var err error
	clientset = fake.NewSimpleClientset()
	app.Clientset = clientset
	db, err = conn.Beginx()
	Expect(err).NotTo(HaveOccurred())
	app.DB = db
})

var _ = AfterEach(func() {
	err := db.Rollback()
	Expect(err).NotTo(HaveOccurred())
	db = nil
	app.DB = conn
})

var _ = AfterSuite(func() {
	if conn != nil {
		err := conn.Close()
		Expect(err).NotTo(HaveOccurred())
		db = nil
	}
})
