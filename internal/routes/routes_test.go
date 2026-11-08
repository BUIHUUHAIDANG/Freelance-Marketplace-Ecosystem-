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
		"GET /user/:id",
		"PUT /user/:id",
		"DELETE /user/:id",
		"GET /user/:id/freelancer",
		"GET /user/:id/hirer",
		"GET /swagger/*any",
	} {
		if !registeredRoutes[route] {
			t.Errorf("route %q is not registered", route)
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
				"/user/{id}",
				"/user/{id}/freelancer",
				"/user/{id}/hirer",
				"UserAccountProfileRespond",
			} {
				if !strings.Contains(response.Body.String(), expected) {
					t.Errorf("generated Swagger document does not contain %q", expected)
				}
			}
		}
	}
}
