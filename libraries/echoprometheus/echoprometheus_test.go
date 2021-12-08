package echoprometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
)

func makeRequest(e *echo.Echo, path string, rec http.ResponseWriter) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e.ServeHTTP(rec, req)
}

func TestPrometheusMiddleware(t *testing.T) {
	e := echo.New()
	NewPrometheus("").Use(e)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	e.GET("/test_echo_error", func(c echo.Context) error {
		return c.String(http.StatusInternalServerError, "test")
	})

	errorHandler := func(c echo.Context) error {
		return fmt.Errorf("internal user error")
	}

	e.GET("/test_user_error_1", errorHandler)
	e.GET("/test_user_error_2", errorHandler)

	rec := httptest.NewRecorder()
	for i := 0; i < 100; i++ {
		makeRequest(e, "/test", rec)            // 100
		makeRequest(e, "/test_echo_error", rec) // 100
	}
	for i := 0; i < 96; i++ {
		// new: 96 per each request, old: 500,GET_/404 96*2
		makeRequest(e, "/test_user_error_1", rec)
		makeRequest(e, "/test_user_error_2", rec)
	}
	for i := 0; i < 69; i++ {
		makeRequest(e, "/test_get_notfound", rec) // new 69 old 404,GET_/404
	}

	// request not found
	req := httptest.NewRequest(http.MethodPost, "/test_post_notfound", nil)
	e.ServeHTTP(rec, req)

	makeRequest(e, "/metrics", rec)
	bodyString := rec.Body.String()
	if !strings.Contains(bodyString, `request_duration_seconds_count{code="200",path="GET_/test"} 100`) {
		t.Error("GET_/test doesnt show")
	}
	if !strings.Contains(bodyString, `request_duration_seconds_count{code="500",path="GET_/test_echo_error"} 100`) {
		t.Error("GET_/test_echo_error doesnt show")
	}

	// // old assert
	// if !strings.Contains(bodyString, `request_duration_seconds_count{code="500",path="GET_/404"} 192`) {
	// 	t.Error("GET_/test_user_error doesnt show")
	// }

	if !strings.Contains(bodyString, `request_duration_seconds_count{code="500",path="GET_/test_user_error_1"} 96`) {
		t.Error("GET_/test_user_error doesnt show")
	}

	if !strings.Contains(bodyString, `request_duration_seconds_count{code="500",path="GET_/test_user_error_2"} 96`) {
		t.Error("GET_/test_user_error doesnt show")
	}

	if !strings.Contains(bodyString, `request_duration_seconds_count{code="404",path="GET_/404"} 69`) {
		t.Error("GET_/404 doesnt show")
	}
	if !strings.Contains(bodyString, `request_duration_seconds_count{code="404",path="POST_/404"} 1`) {
		t.Error("POST_/404 doesnt show")
	}

	// ioutil.WriteFile("a.txt", []byte(bodyString), 0644)
	// t.Logf("response %v '%v'", rec.Code, bodyString)
}
