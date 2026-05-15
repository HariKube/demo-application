package apiserver

import (
	"context"
	"net/http"

	kaf "github.com/HariKube/kubernetes-aggregator-framework/pkg/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Group   = "custom.report.example.example.com"
	Version = "v1"
)

func New(port, certFile, keyFile string) *customAPIServer {
	sas := customAPIServer{
		Server: *kaf.NewServer(kaf.ServerConfig{
			Port:     port,
			CertFile: certFile,
			KeyFile:  keyFile,
			Group:    Group,
			Version:  Version,
			APIKinds: []kaf.APIKind{
				{
					ApiResource: metav1.APIResource{
						Name:       "customreports",
						Namespaced: true,
						Kind:       "CustomReport",
						Verbs:      []string{"get", "list", "watch", "create", "update", "delete"},
					},
					CustomResource: &kaf.CustomResource{
						CreateHandler: func(namespace, name string, w http.ResponseWriter, r *http.Request) {
							http.Error(w, "I'm a teapot", http.StatusTeapot)
						},
						GetHandler: func(namespace, name string, w http.ResponseWriter, r *http.Request) {
							http.Error(w, "I'm a teapot", http.StatusTeapot)
						},
						ListHandler: func(namespace, name string, w http.ResponseWriter, r *http.Request) {
							http.Error(w, "I'm a teapot", http.StatusTeapot)
						},
						ReplaceHandler: func(namespace, name string, w http.ResponseWriter, r *http.Request) {
							http.Error(w, "I'm a teapot", http.StatusTeapot)
						},
						DeleteHandler: func(namespace, name string, w http.ResponseWriter, r *http.Request) {
							http.Error(w, "I'm a teapot", http.StatusTeapot)
						},
						WatchHandler: func(namespace, name string, w http.ResponseWriter, r *http.Request) {
							http.Error(w, "I'm a teapot", http.StatusTeapot)
						},
					},
				},
			},
		}),
	}

	return &sas
}

type customAPIServer struct {
	kaf.Server
}

func (s *customAPIServer) Start(ctx context.Context) (err error) {
	return s.Server.Start(ctx)
}
