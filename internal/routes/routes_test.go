package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSwaggerDocumentationRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router, nil)

	registeredRoutes := make(map[string]bool)
	for _, route := range router.Routes() {
		registeredRoutes[route.Method+" "+route.Path] = true
	}
	for _, route := range []string{
		"GET /health",
		"POST /auth/login",
		"POST /auth/register",
		"GET /user/me",
		"PUT /user/me",
		"DELETE /user/me",
		"GET /user/me/freelancer",
		"GET /user/me/hirer",
		"GET /job",
		"GET /job/:jobid",
		"POST /job",
		"PUT /job/:jobid",
		"DELETE /job/:jobid",
		"POST /job/:jobid/apply",
		"GET /swagger/*any",
	} {
		if !registeredRoutes[route] {
			t.Errorf("route %q is not registered", route)
		}
	}
	for _, removedRoute := range []string{
		"GET /user/:id",
		"PUT /user/:id",
		"DELETE /user/:id",
		"GET /user/:id/freelancer",
		"GET /user/:id/hirer",
	} {
		if registeredRoutes[removedRoute] {
			t.Errorf("user route %q should not accept a user ID from the URL", removedRoute)
		}
	}

	for _, path := range []string{"/swagger/index.html", "/swagger/doc.json"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Errorf("GET %s returned %d", path, response.Code)
		}
		if path == "/swagger/doc.json" {
			for _, expected := range []string{
				"/health",
				"/auth/login",
				"/auth/register",
				"/user/me",
				"/user/me/freelancer",
				"/user/me/hirer",
				"Get my user account profile",
				"Update my user account profile",
				"Delete my user account",
				"/job",
				"/job/{jobid}",
				"/job/{jobid}/apply",
				"JobListRespond",
				"JobRespond",
				"AddJobRequest",
				"ApplyJobRequest",
				"UserAccountProfileRespond",
			} {
				if !strings.Contains(response.Body.String(), expected) {
					t.Errorf("generated Swagger document does not contain %q", expected)
				}
			}
		}
	}
}

func TestProtectedRoutesRequireJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router, nil)

	for _, test := range []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/user/me"},
		{method: http.MethodPut, path: "/user/me"},
		{method: http.MethodDelete, path: "/user/me"},
		{method: http.MethodGet, path: "/user/me/freelancer"},
		{method: http.MethodGet, path: "/user/me/hirer"},
		{method: http.MethodPost, path: "/job"},
		{method: http.MethodPut, path: "/job/1"},
		{method: http.MethodDelete, path: "/job/1"},
		{method: http.MethodPost, path: "/job/1/apply"},
	} {
		request := httptest.NewRequest(test.method, test.path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Errorf("%s %s returned %d, want %d", test.method, test.path, response.Code, http.StatusUnauthorized)
		}
	}
}
