package common

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWrapFunc(t *testing.T) {
	testCases := []struct {
		Name   string
		Func   interface{}
		Inputs []interface{}

		ExpectedPanic    bool
		ExpectedRespCode int
		ExpectedRespBody string
	}{
		{"invalid input parameters", func(int) string { return "" }, []interface{}{1, 2}, true, 0, ""},
		{"invalid output parameters", func() {}, nil, true, 0, ""},
		{"error response", func() (string, error) { return "", assert.AnError }, nil, false, 500, fmt.Sprintf(`{"code":500,"msg":"%s","data":null}`, assert.AnError)},
		{"function panic", func() (string, error) { panic("some error"); return "", assert.AnError }, nil, false, 500, fmt.Sprintf(`{"code":500,"msg":"some error","data":null}`)},
		{"success response with one outputs", func() string { return "some msg" }, nil, false, 200, `"some msg"`},
		{"success response with two outputs", func() (string, error) { return "some msg", nil }, nil, false, 200, `"some msg"`},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.ExpectedPanic {
				defer func() {
					err := recover()
					assert.NotNil(t, err)
				}()
			}

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			h := WrapFunc(tc.Func, tc.Inputs...)
			h(c)

			assert.Equal(t, tc.ExpectedRespCode, w.Code)
			assert.Equal(t, tc.ExpectedRespBody, w.Body.String())
		})
	}
}
