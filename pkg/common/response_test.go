package common

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	testCases := []struct {
		name    string
		code    int
		err     error
		success bool
		resp    Response
	}{
		{
			name: "with code 400",
			code: 400,
			err:  errors.New("some err"),
			resp: Response{
				Code: 400,
				Msg:  "some err",
			},
		},
		{
			name: "with code 401",
			code: 401,
			err:  errors.New("some err"),
			resp: Response{
				Code: 401,
				Msg:  "some err",
			},
		},
		// {
		// 	name: "with code 400 and err is nil",
		// 	code: 400,
		// 	resp: Response{
		// 		Code: 400,
		// 	},
		// },
		{
			name: "with error",
			err:  errors.New("some err"),
			resp: Response{
				Code: 500,
				Msg:  "some err",
			},
		},
		{
			name:    "response success",
			success: true,
			resp: Response{
				Code: 200,
				Msg:  "success",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "http://localhost/api/v1/apps", nil)

			if tc.success {
				ResponseSuccess(c, nil)
			} else {
				ResponseFailed(c, tc.code, tc.err)
			}

			resp := Response{}
			err := json.NewDecoder(w.Body).Decode(&resp)
			assert.Empty(t, err)
			assert.Equal(t, tc.resp, resp)
		})
	}
}
