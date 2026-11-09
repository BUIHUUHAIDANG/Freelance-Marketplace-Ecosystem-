package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGenerateAndExtractTokenUserID(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret")
	t.Setenv("TOKEN_LIFESPAN", "1")

	token, err := GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/user/42", nil)
	ctx.Request.Header.Set("Authorization", "Bearer "+token)

	userID, err := ExtractTokenUserID(ctx)
	if err != nil {
		t.Fatalf("ExtractTokenUserID() error = %v", err)
	}
	if userID != 42 {
		t.Fatalf("ExtractTokenUserID() = %d, want 42", userID)
	}
}

func TestGenerateTokenRequiresJWTConfiguration(t *testing.T) {
	t.Setenv("SECRET_KEY", "")
	t.Setenv("TOKEN_LIFESPAN", "1")
	if _, err := GenerateToken(42); err == nil {
		t.Fatal("GenerateToken() succeeded without SECRET_KEY")
	}

	t.Setenv("SECRET_KEY", "test-secret")
	t.Setenv("TOKEN_LIFESPAN", "0")
	if _, err := GenerateToken(42); err == nil {
		t.Fatal("GenerateToken() succeeded with non-positive TOKEN_LIFESPAN")
	}
}

func TestExtractTokenRequiresBearerScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/user/42", nil)
	ctx.Request.Header.Set("Authorization", "Basic abc")
	if token := ExtractToken(ctx); token != "" {
		t.Fatalf("ExtractToken() = %q, want empty token", token)
	}

	ctx.Request.Header.Set("Authorization", "bearer valid-token")
	if token := ExtractToken(ctx); token != "valid-token" {
		t.Fatalf("ExtractToken() = %q, want valid-token", token)
	}
}
