package kubecontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/qingwave/weave/pkg/common"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
)

type errorResponder struct{}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	logrus.Errorf("proxy k8s err: %v", err)
	resp := &common.Response{
		Code: http.StatusInternalServerError,
		Msg:  err.Error(),
	}
	json.NewEncoder(w).Encode(resp)
}

func ProxyKubeAPIServer(config *rest.Config) gin.HandlerFunc {
	kubernetes, _ := url.Parse(config.Host)
	defaultTransport, err := rest.TransportFor(config)
	if err != nil {
		logrus.Errorf("Unable to create transport from rest.Config: %v", err)
		return func(*gin.Context) {}
	}

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err == http.ErrAbortHandler {
				return
			} else if err != nil {
				common.ResponseFailed(c, http.StatusInternalServerError, fmt.Errorf("%v", err))
			}
		}()

		s := *c.Request.URL
		s.Host = kubernetes.Host
		s.Scheme = kubernetes.Scheme

		// make sure we don't override kubernetes's authorization
		c.Request.Header.Del("Authorization")
		httpProxy := proxy.NewUpgradeAwareHandler(&s, defaultTransport, true, false, &errorResponder{})
		httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(defaultTransport, defaultTransport)
		httpProxy.ServeHTTP(c.Writer, c.Request)
	}
}
