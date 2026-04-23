package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"auth-project/internal/service"
)

// AuthController handles HTTP requests for authentication
type AuthController struct {
	authService *service.AuthService
	cookieName  string
	cookiePath  string
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		cookieName:  "session_id",
		cookiePath:  "/",
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Registration details"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	resp, err := c.authService.Register(req)
	if err != nil {
		switch err {
		case service.ErrUserAlreadyExists:
			BadRequestResponse(ctx, err)
		case service.ErrInvalidEmail, service.ErrWeakPassword:
			ValidationErrorResponse(ctx, err)
		default:
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponseWithMessage(ctx, "User registered successfully", resp)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and create session
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Login credentials"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	resp, err := c.authService.Login(req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			UnauthorizedResponse(ctx, err)
		} else {
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	// Set session cookie
	c.setSessionCookie(ctx, resp.SessionID, resp.ExpiresAt)

	SuccessResponseWithMessage(ctx, "Login successful", resp)
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate user session
// @Tags Auth
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	sessionID, err := ctx.Cookie(c.cookieName)
	if err != nil {
		// No session cookie, still return success
		SuccessResponseWithMessage(ctx, "Logged out successfully", nil)
		return
	}

	if err := c.authService.Logout(sessionID); err != nil {
		// Log error but still return success to client
		// This ensures client cookie is cleared even if server-side cleanup fails
	}

	// Clear session cookie
	c.clearSessionCookie(ctx)

	SuccessResponseWithMessage(ctx, "Logged out successfully", nil)
}

// GetCurrentUser returns the currently authenticated user
// @Summary Get current user
// @Description Get details of the currently authenticated user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /auth/me [get]
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	sessionID, err := ctx.Cookie(c.cookieName)
	if err != nil {
		UnauthorizedResponse(ctx, err)
		return
	}

	user, err := c.authService.ValidateSession(sessionID)
	if err != nil {
		UnauthorizedResponse(ctx, err)
		return
	}

	SuccessResponse(ctx, gin.H{
		"user_id": user.ID.String(),
		"email":   user.Email,
	})
}

// ValidateSession validates the current session
// @Summary Validate session
// @Description Check if the current session is valid
// @Tags Auth
// @Produce json
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /auth/validate [get]
func (c *AuthController) ValidateSession(ctx *gin.Context) {
	sessionID, err := ctx.Cookie(c.cookieName)
	if err != nil {
		UnauthorizedResponse(ctx, err)
		return
	}

	user, err := c.authService.ValidateSession(sessionID)
	if err != nil {
		UnauthorizedResponse(ctx, err)
		return
	}

	SuccessResponse(ctx, gin.H{
		"valid":   true,
		"user_id": user.ID.String(),
		"email":   user.Email,
	})
}

// ListUsers handles GET /users request
// @Summary List all users
// @Description Retrieve a list of all registered users (for assignment dropdown)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /users [get]
func (c *AuthController) ListUsers(ctx *gin.Context) {
	users, err := c.authService.ListUsers()
	if err != nil {
		InternalServerErrorResponse(ctx, err)
		return
	}

	// Strip password hash from response
	var safeUsers []gin.H
	for _, user := range users {
		safeUsers = append(safeUsers, gin.H{
			"id":         user.ID.String(),
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	SuccessResponse(ctx, safeUsers)
}

// Helper methods for cookie management
func (c *AuthController) setSessionCookie(ctx *gin.Context, sessionID string, expiresAt time.Time) {
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		c.cookieName,
		sessionID,
		int(time.Until(expiresAt).Seconds()),
		c.cookiePath,
		"",
		false, // Secure: set to true in production with HTTPS
		true,  // HttpOnly
	)
}

func (c *AuthController) clearSessionCookie(ctx *gin.Context) {
	ctx.SetCookie(
		c.cookieName,
		"",
		-1, // Expire immediately
		c.cookiePath,
		"",
		false,
		true,
	)
}
