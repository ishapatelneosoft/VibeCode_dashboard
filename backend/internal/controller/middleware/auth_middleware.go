package middleware

import (
	"github.com/gin-gonic/gin"

	"auth-project/internal/controller"
	"auth-project/internal/domain"
	"auth-project/internal/service"
)

// AuthMiddleware provides authentication middleware for protected routes
type AuthMiddleware struct {
	authService *service.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware ensures the request is authenticated
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			controller.UnauthorizedResponse(ctx, err)
			ctx.Abort()
			return
		}

		user, err := m.authService.ValidateSession(sessionID)
		if err != nil {
			controller.UnauthorizedResponse(ctx, err)
			ctx.Abort()
			return
		}

		// Store user in context for downstream handlers
		ctx.Set("user_id", user.ID.String())
		ctx.Set("user_email", user.Email)
		ctx.Set("user", user)

		ctx.Next()
	}
}

// OptionalAuth middleware attaches user context if authenticated, but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			// No session cookie, continue without authentication
			ctx.Next()
			return
		}

		user, err := m.authService.ValidateSession(sessionID)
		if err != nil {
			// Invalid session, continue without authentication
			ctx.Next()
			return
		}

		// Store user in context for downstream handlers
		ctx.Set("user_id", user.ID.String())
		ctx.Set("user_email", user.Email)
		ctx.Set("user", user)

		ctx.Next()
	}
}

// RequireRoles middleware ensures the authenticated user has required roles
// Note: This is a placeholder for role-based authorization
func (m *AuthMiddleware) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// First, require authentication
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			controller.UnauthorizedResponse(ctx, err)
			ctx.Abort()
			return
		}

		user, err := m.authService.ValidateSession(sessionID)
		if err != nil {
			controller.UnauthorizedResponse(ctx, err)
			ctx.Abort()
			return
		}

		// Store user in context
		ctx.Set("user_id", user.ID.String())
		ctx.Set("user_email", user.Email)
		ctx.Set("user", user)

		// TODO: Implement role checking when role system is added
		// For now, just continue
		ctx.Next()
	}
}

// GetUserFromContext retrieves the authenticated user from context
func GetUserFromContext(ctx *gin.Context) (*domain.User, bool) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, false
	}

	domainUser, ok := user.(*domain.User)
	return domainUser, ok
}

// GetUserIDFromContext retrieves the authenticated user ID from context
func GetUserIDFromContext(ctx *gin.Context) (string, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetUserEmailFromContext retrieves the authenticated user email from context
func GetUserEmailFromContext(ctx *gin.Context) (string, bool) {
	email, exists := ctx.Get("user_email")
	if !exists {
		return "", false
	}

	emailStr, ok := email.(string)
	return emailStr, ok
}
