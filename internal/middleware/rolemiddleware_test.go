package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/utils"
	"github.com/gin-gonic/gin"
)

func TestIsConfiguredAdmin(t *testing.T) {
	t.Setenv("ADMIN_USER_IDS", "2, 15")

	for _, test := range []struct {
		userID int64
		want   bool
	}{
		{userID: 2, want: true},
		{userID: 15, want: true},
		{userID: 7, want: false},
	} {
		got, err := isConfiguredAdmin(test.userID)
		if err != nil {
			t.Fatalf("isConfiguredAdmin(%d) returned error: %v", test.userID, err)
		}
		if got != test.want {
			t.Errorf("isConfiguredAdmin(%d) = %t, want %t", test.userID, got, test.want)
		}
	}
}

func TestIsConfiguredAdminRejectsInvalidConfiguration(t *testing.T) {
	t.Setenv("ADMIN_USER_IDS", "2,not-an-id")

	if _, err := isConfiguredAdmin(2); err == nil {
		t.Fatal("isConfiguredAdmin accepted an invalid ADMIN_USER_IDS value")
	}
}

func TestConfiguredAdminCanPassJobRoleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("SECRET_KEY", "test-secret")
	t.Setenv("TOKEN_LIFESPAN", "1")
	t.Setenv("ADMIN_USER_IDS", "42")

	token, err := utils.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	router := gin.New()
	router.POST("/jobs", JwtAuthMiddleware(), RequireJobManagerOrAdmin(nil), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodPost, "/jobs", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("admin request returned %d, want %d", response.Code, http.StatusNoContent)
	}
}
