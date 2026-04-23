package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"auth-project/internal/service"
)

// WorklogController handles HTTP requests for worklog operations
type WorklogController struct {
	worklogService *service.WorklogService
}

// NewWorklogController creates a new worklog controller
func NewWorklogController(worklogService *service.WorklogService) *WorklogController {
	return &WorklogController{
		worklogService: worklogService,
	}
}

// LogWork handles POST /tasks/{id}/worklogs request
// @Summary Log time spent on a task
// @Description Create an immutable worklog entry for a task
// @Tags Worklogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body service.LogWorkRequest true "Worklog details"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /tasks/{id}/worklogs [post]
func (c *WorklogController) LogWork(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	var req service.LogWorkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		UnauthorizedResponse(ctx, gin.Error{})
		return
	}
	loggedBy, err := uuid.Parse(userID.(string))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	resp, err := c.worklogService.LogWork(taskID, loggedBy, req)
	if err != nil {
		switch err {
		case service.ErrInvalidTimeSpent:
			ValidationErrorResponse(ctx, err)
		case service.ErrTaskNotFound:
			NotFoundResponse(ctx, err)
		default:
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponseWithMessage(ctx, "Worklog created successfully", resp)
}

// GetWorklogs handles GET /tasks/{id}/worklogs request
// @Summary Get worklogs for a task
// @Description Retrieve all worklog entries for a task, most recent first
// @Tags Worklogs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /tasks/{id}/worklogs [get]
func (c *WorklogController) GetWorklogs(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	worklogs, err := c.worklogService.GetWorklogs(taskID)
	if err != nil {
		if err == service.ErrTaskNotFound {
			NotFoundResponse(ctx, err)
		} else {
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponse(ctx, worklogs)
}

// GetTimeReport handles GET /reports/time request
// @Summary Get time report across all tasks
// @Description Retrieve aggregated time spent per task and grand total
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /reports/time [get]
func (c *WorklogController) GetTimeReport(ctx *gin.Context) {
	report, err := c.worklogService.GetTimeReport()
	if err != nil {
		InternalServerErrorResponse(ctx, err)
		return
	}

	SuccessResponse(ctx, report)
}
