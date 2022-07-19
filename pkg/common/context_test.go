package common

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/utils/request"
	"github.com/qingwave/weave/pkg/utils/trace"
	"github.com/stretchr/testify/assert"
)

func TestTraceContext(t *testing.T) {

	trace := trace.New("test", logr.Discard())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SetTrace(c, nil)
	assert.Nil(t, GetTrace(c))

	SetTrace(c, trace)

	TraceStep(c, "msg")

	assert.NotNil(t, GetTrace(c))
}

func TestUserContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SetUser(c, nil)
	assert.Nil(t, GetUser(c))

	user := &model.User{ID: 1, Name: "some"}
	SetUser(c, user)

	assert.Equal(t, user, GetUser(c))
}

func TestRequestInfoContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SetRequestInfo(c, nil)
	assert.Nil(t, GetRequestInfo(c))

	ri := &request.RequestInfo{Verb: "get", Resource: "apps"}
	SetRequestInfo(c, ri)

	assert.Equal(t, ri, GetRequestInfo(c))
}
