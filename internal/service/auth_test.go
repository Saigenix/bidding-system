package service

import (
	"context"
	"testing"

	"github.com/saigenix/bidding-system/internal/mocks"
)

func newTestAuthService() (*AuthService, *mocks.MockUserRepository) {
	repo := mocks.NewMockUserRepository()
	svc := NewAuthService(repo, "test-secret-key-for-testing", 24)
	return svc, repo
}

// ============================================================================
// Register
// ============================================================================

func TestAuthService_Register_Success(t *testing.T) {
	svc, _ := newTestAuthService()

	user, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Register() unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("Register() returned nil user")
	}
	if user.Email != "test@example.com" {
		t.Errorf("Register() email = %q, want %q", user.Email, "test@example.com")
	}
	if user.ID == "" {
		t.Error("Register() user ID is empty")
	}
	if user.PasswordHash == "" {
		t.Error("Register() password hash is empty")
	}
	if user.PasswordHash == "password123" {
		t.Error("Register() password was not hashed")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc, _ := newTestAuthService()

	_, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("First Register() unexpected error: %v", err)
	}

	_, err = svc.Register(context.Background(), "test@example.com", "password456")
	if err == nil {
		t.Error("Second Register() expected error for duplicate email, got nil")
	}
}

// ============================================================================
// Login
// ============================================================================

func TestAuthService_Login_Success(t *testing.T) {
	svc, _ := newTestAuthService()

	_, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Register() unexpected error: %v", err)
	}

	token, err := svc.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Login() unexpected error: %v", err)
	}
	if token == "" {
		t.Error("Login() returned empty token")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc, _ := newTestAuthService()

	_, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Register() unexpected error: %v", err)
	}

	_, err = svc.Login(context.Background(), "test@example.com", "wrongpassword")
	if err == nil {
		t.Error("Login() expected error for wrong password, got nil")
	}
}

func TestAuthService_Login_NonExistentUser(t *testing.T) {
	svc, _ := newTestAuthService()

	_, err := svc.Login(context.Background(), "nobody@example.com", "password123")
	if err == nil {
		t.Error("Login() expected error for non-existent user, got nil")
	}
}

// ============================================================================
// ValidateToken
// ============================================================================

func TestAuthService_ValidateToken_Valid(t *testing.T) {
	svc, _ := newTestAuthService()

	user, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Register() unexpected error: %v", err)
	}

	token, err := svc.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Login() unexpected error: %v", err)
	}

	userID, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() unexpected error: %v", err)
	}
	if userID != user.ID {
		t.Errorf("ValidateToken() userID = %q, want %q", userID, user.ID)
	}
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	svc, _ := newTestAuthService()

	_, err := svc.ValidateToken("invalid.token.string")
	if err == nil {
		t.Error("ValidateToken() expected error for invalid token, got nil")
	}
}

func TestAuthService_ValidateToken_WrongSecret(t *testing.T) {
	svc1, _ := newTestAuthService()
	_, err := svc1.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Register() unexpected error: %v", err)
	}

	token, err := svc1.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Login() unexpected error: %v", err)
	}

	// Validate with a different secret
	svc2 := NewAuthService(mocks.NewMockUserRepository(), "different-secret-key", 24)
	_, err = svc2.ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken() expected error for token signed with different secret, got nil")
	}
}
