package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCurrentUserIDUsesAuthenticatedIdentityInsteadOfURLID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Params = gin.Params{{Key: "id", Value: "99"}}
	ctx.Set("user_id", int64(42))

	userID, ok := currentUserID(ctx)
	if !ok {
		t.Fatal("currentUserID() did not accept the authenticated user")
	}
	if userID != 42 {
		t.Fatalf("currentUserID() = %d, want authenticated user ID 42", userID)
	}
}

func TestCurrentUserIDRejectsMissingAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	if _, ok := currentUserID(ctx); ok {
		t.Fatal("currentUserID() accepted a request without an authenticated user")
	}
	if recorder.Code != 401 {
		t.Fatalf("currentUserID() status = %d, want 401", recorder.Code)
	}
}
