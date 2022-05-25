package request

import (
	"net/http"
	"testing"

	"github.com/qingwave/weave/pkg/utils/set"
	"github.com/stretchr/testify/assert"
)

func TestRequestInfo(t *testing.T) {
	resolver := RequestInfoFactory{set.NewString("api")}
	baseAddr := "http://localhost:8080"

	testCases := []struct {
		Name   string
		Method string
		URI    string

		ExpectedError       bool
		ExpectedRequestInfo *RequestInfo
	}{
		{"non resource", "GET", "/version", false, &RequestInfo{Verb: "get"}},
		{"resource with list", "GET", "/api/v1/containers", false, &RequestInfo{
			IsResourceRequest: true,
			Verb:              "list",
			APIPrefix:         "api",
			APIVersion:        "v1",
			Namespace:         "root",
			Resource:          "containers",
			Parts:             []string{"containers"},
		}},
		{"update sub resource", "PUT", "/api/v1/namespaces/ns1/jobs/1/log", false, &RequestInfo{
			IsResourceRequest: true,
			Verb:              "update",
			APIPrefix:         "api",
			APIVersion:        "v1",
			Namespace:         "ns1",
			Resource:          "jobs",
			Subresource:       "log",
			Name:              "1",
			Parts:             []string{"jobs", "1", "log"},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest(tc.Method, baseAddr+tc.URI, nil)
			assert.NoError(t, err)

			if tc.ExpectedRequestInfo != nil {
				tc.ExpectedRequestInfo.Path = tc.URI
			}

			ri, err := resolver.NewRequestInfo(req)
			if tc.ExpectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.ExpectedRequestInfo, ri)
			}
		})
	}
}
