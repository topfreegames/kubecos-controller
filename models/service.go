// mystack-controller api
// https://github.com/topfreegames/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"bytes"
	"strings"
	"text/template"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

const serviceYaml = `
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    mystack/routable: "true"
spec:
  selector:
    app: {{.Name}}
  ports:
    - protocol: TCP
      port: {{.Port}}
      targetPort: {{.TargetPort}}
  type: ClusterIP
`

//Service represents a service
type Service struct {
	Name       string
	Namespace  string
	Port       int
	TargetPort int
}

//NewService is the service ctor
func NewService(name, username string, port, targetPort int) *Service {
	username = strings.Replace(username, ".", "-", -1)
	namespace := usernameToNamespace(username)
	return &Service{
		Name:       name,
		Namespace:  namespace,
		Port:       port,
		TargetPort: targetPort,
	}
}

//Expose exposes a deployment
func (s *Service) Expose(clientset kubernetes.Interface) (*v1.Service, error) {
	tmpl, err := template.New("expose").Parse(serviceYaml)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, s)
	if err != nil {
		return nil, err
	}

	decoder := api.Codecs.UniversalDecoder()
	obj, _, err := decoder.Decode(buf.Bytes(), nil, nil)
	if err != nil {
		return nil, err
	}

	src := obj.(*api.Service)
	dst := &v1.Service{}

	err = api.Scheme.Convert(src, dst, 0)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Services(s.Namespace).Create(dst)
}

//Delete deletes service
func (s *Service) Delete(clientset kubernetes.Interface) error {
	deleteOptions := &v1.DeleteOptions{}
	return clientset.CoreV1().Services(s.Namespace).Delete(s.Name, deleteOptions)
}
