// mystack-controller api
// https://github.com/topfreegames/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/topfreegames/mystack-controller/models"
)

//ClusterHandler handles cluster creation and deletion
type ClusterHandler struct {
	App    *App
	Method string
}

func (c *ClusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch c.Method {
	case "create":
		c.create(w, r)
	case "delete":
		c.deleteCluster(w, r)
	case "routes":
		c.getRoutes(w, r)
	}
}

func (c *ClusterHandler) create(w http.ResponseWriter, r *http.Request) {
	logger := loggerFromContext(r.Context())
	email := emailFromCtx(r.Context())
	username := usernameFromEmail(email)

	log(logger, "Creating cluster for user %s", username)
	clusterName := GetClusterName(r)

	cluster, err := models.NewCluster(
		c.App.DB,
		username,
		clusterName,
		c.App.DeploymentReadiness,
		c.App.JobReadiness,
	)
	if err != nil {
		c.App.HandleError(w, Status(err), "create cluster error", err)
		return
	}

	err = cluster.Create(c.App.Logger, c.App.Clientset)
	if err != nil {
		c.App.HandleError(w, Status(err), "create cluster error", err)
		return
	}

	routes, err := cluster.Routes(c.App.AppsRoutesDomain, c.App.Clientset)
	if err != nil {
		c.App.HandleError(w, Status(err), "create cluster error", err)
		return
	}
	routesResponse := map[string][]string{
		"routes": routes,
	}
	bts, err := json.Marshal(&routesResponse)
	if err != nil {
		c.App.HandleError(w, Status(err), "create cluster error", err)
		return
	}

	WriteBytes(w, http.StatusOK, bts)
	log(logger, "Cluster successfully created for user %s", username)
}

func (c *ClusterHandler) deleteCluster(w http.ResponseWriter, r *http.Request) {
	logger := loggerFromContext(r.Context())
	email := emailFromCtx(r.Context())
	username := usernameFromEmail(email)

	log(logger, "Deleting cluster for user %s", username)
	clusterName := GetClusterName(r)

	cluster, err := models.NewCluster(
		c.App.DB,
		username,
		clusterName,
		c.App.DeploymentReadiness,
		c.App.JobReadiness,
	)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		cluster = &models.Cluster{Username: username}
	} else if err != nil {
		c.App.HandleError(w, Status(err), "retrieve cluster error", err)
		return
	}

	err = cluster.Delete(c.App.Clientset)
	if err != nil {
		c.App.HandleError(w, Status(err), "delete cluster error", err)
		return
	}

	Write(w, http.StatusOK, `{"status": "ok"}`)
	log(logger, "Cluster deleted for user %s", username)
}

func (c *ClusterHandler) getRoutes(w http.ResponseWriter, r *http.Request) {
	logger := loggerFromContext(r.Context())
	email := emailFromCtx(r.Context())
	username := usernameFromEmail(email)

	log(logger, "Cluster routes for user %s", username)
	clusterName := GetClusterName(r)

	cluster, err := models.NewCluster(c.App.DB, username, clusterName, nil, nil)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		cluster = &models.Cluster{Username: username}
	} else if err != nil {
		c.App.HandleError(w, Status(err), "retrieve cluster error", err)
		return
	}

	routes, err := cluster.Routes(c.App.AppsRoutesDomain, c.App.Clientset)
	if err != nil {
		c.App.HandleError(w, Status(err), "create cluster error", err)
		return
	}

	routesResponse := make(map[string][]string)
	routesResponse["routes"] = routes
	bts, err := json.Marshal(routesResponse)
	if err != nil {
		c.App.HandleError(w, Status(err), "get routes error", err)
		return
	}

	WriteBytes(w, http.StatusOK, bts)
	log(logger, "Cluster routes built for user %s", username)
}
