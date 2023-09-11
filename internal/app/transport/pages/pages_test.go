package pages

import (
	//"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleShurlPage(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		// TODO: Add test cases.
		{
			name: "something to test GET",
			want: want{
				contentType: "text/plain",
				statusCode:  307,
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleShurlPage)
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			err := result.Body.Close()
			if err != nil {
				return
			}
			//HandleShurlPage(tt.args.res, tt.args.req)
		})
	}
}
