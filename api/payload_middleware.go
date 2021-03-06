// mystack-controller api
// https://github.com/topfreegames/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//PayloadMiddleware handles cluster creation and deletion
type PayloadMiddleware struct {
	App  *App
	next http.Handler
}

const configKey = contextKey("clusterConfigKey")

//NewContextWithClusterConfig creates a context with cluster config
func NewContextWithClusterConfig(ctx context.Context, clusterConfig string) context.Context {
	c := context.WithValue(ctx, configKey, clusterConfig)
	return c
}

func clusterConfigFromCtx(ctx context.Context) string {
	clusterConfig := ctx.Value(configKey)
	if clusterConfig == nil {
		return ""
	}
	return clusterConfig.(string)
}

func toLiteral(bts []byte) []byte {
	bts = bytes.Replace(bts, []byte("\n"), []byte(`\n`), -1)
	bts = bytes.Replace(bts, []byte("\t"), []byte("  "), -1)

	return bts
}

func (p *PayloadMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.App.HandleError(w, http.StatusBadRequest, "error reading body", err)
		return
	}

	bts = toLiteral(bts)

	p.App.Logger.Info(fmt.Sprintf("yaml cluster config: %s", string(bts)))

	bodyJSON := make(map[string]string)
	err = json.Unmarshal(bts, &bodyJSON)
	if err != nil {
		msg := fmt.Sprintf("error reading body: %s", string(bts))
		p.App.HandleError(w, http.StatusBadRequest, msg, err)
		return
	}

	ctx := NewContextWithClusterConfig(r.Context(), bodyJSON["yaml"])
	p.next.ServeHTTP(w, r.WithContext(ctx))
}

//SetNext handler
func (p *PayloadMiddleware) SetNext(next http.Handler) {
	p.next = next
}
