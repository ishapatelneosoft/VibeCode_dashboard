package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"auth-project/internal/service"
)

// TaskController handles HTTP requests for task operations
type TaskController struct {
	taskService *service.TaskService
}

// NewTaskController creates a new task controller
func NewTaskController(taskService *service.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// CreateTask handles POST /tasks request
// @Summary Create a new task
// @Description Create a new task in a column
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateTaskRequest true "Task details"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /tasks [post]
func (c *TaskController) CreateTask(ctx *gin.Context) {
	var req service.CreateTaskRequest
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

	// Parse user ID
	createdBy, err := uuid.Parse(userID.(string))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	// Create task
	resp, err := c.taskService.CreateTask(req, createdBy)
	if err != nil {
		switch err {
		case service.ErrInvalidTitle:
			ValidationErrorResponse(ctx, err)
		case service.ErrColumnNotFound:
			NotFoundResponse(ctx, err)
		default:
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponseWithMessage(ctx, "Task created successfully", resp)
}

// GetTask handles GET /tasks/:id request
// @Summary Get a task by ID
// @Description Get details of a specific task
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /tasks/{id} [get]
func (c *TaskController) GetTask(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	task, err := c.taskService.GetTask(taskID)
	if err != nil {
		if err == service.ErrTaskNotFound {
			NotFoundResponse(ctx, err)
		} else {
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponse(ctx, task)
}

// GetTasksByColumn handles GET /columns/:columnId/tasks request
func (c *TaskController) GetTasksByColumn(ctx *gin.Context) {
	columnID, err := uuid.Parse(ctx.Param("columnId"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	tasks, err := c.taskService.GetTasksByColumn(columnID)
	if err != nil {
		InternalServerErrorResponse(ctx, err)
		return
	}

	SuccessResponse(ctx, tasks)
}

// UpdateTask handles PUT /tasks/:id request
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	if err := c.taskService.UpdateTask(taskID, updates); err != nil {
		switch err {
		case service.ErrInvalidTitle:
			ValidationErrorResponse(ctx, err)
		case service.ErrTaskNotFound:
			NotFoundResponse(ctx, err)
		default:
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponseWithMessage(ctx, "Task updated successfully", nil)
}

// MoveTask handles PATCH /tasks/:id/move request
// @Summary Move a task to a different column or position
// @Description Move a task using fractional indexing for ordering
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body service.MoveTaskRequest true "Movement details"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /tasks/{id}/move [patch]
func (c *TaskController) MoveTask(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	var req service.MoveTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	resp, err := c.taskService.MoveTask(taskID, req)
	if err != nil {
		switch err {
		case service.ErrTaskNotFound:
			NotFoundResponse(ctx, err)
		case service.ErrColumnNotFound:
			NotFoundResponse(ctx, err)
		default:
			// Check if it's a validation error (string matching)
			if err.Error() == "before task not found" ||
				err.Error() == "after task not found" ||
				err.Error() == "before task is not in destination column" ||
				err.Error() == "after task is not in destination column" {
				ValidationErrorResponse(ctx, err)
			} else {
				InternalServerErrorResponse(ctx, err)
			}
		}
		return
	}

	SuccessResponseWithMessage(ctx, "Task moved successfully", resp)
}

// UpdateAssignee handles PATCH /tasks/{id}/assignee request
// @Summary Update task assignee
// @Description Assign or unassign a user to a task, recording history
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body service.UpdateAssigneeRequest true "Assignee details"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /tasks/{id}/assignee [patch]
func (c *TaskController) UpdateAssignee(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	var req service.UpdateAssigneeRequest
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
	changedBy, err := uuid.Parse(userID.(string))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	if err := c.taskService.UpdateAssignee(taskID, req.AssigneeID, changedBy); err != nil {
		switch err {
		case service.ErrTaskNotFound:
			NotFoundResponse(ctx, err)
		case service.ErrAssigneeNotFound:
			ValidationErrorResponse(ctx, err)
		default:
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponseWithMessage(ctx, "Assignee updated successfully", nil)
}

// GetAssignmentHistory handles GET /tasks/{id}/history request
// @Summary Get assignment history for a task
// @Description Retrieve all assignee changes for a task, most recent first
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /tasks/{id}/history [get]
func (c *TaskController) GetAssignmentHistory(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	histories, err := c.taskService.GetAssignmentHistory(taskID)
	if err != nil {
		if err == service.ErrTaskNotFound {
			NotFoundResponse(ctx, err)
		} else {
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponse(ctx, histories)
}

// DeleteTask handles DELETE /tasks/:id request
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		BadRequestResponse(ctx, err)
		return
	}

	if err := c.taskService.DeleteTask(taskID); err != nil {
		if err == service.ErrTaskNotFound {
			NotFoundResponse(ctx, err)
		} else {
			InternalServerErrorResponse(ctx, err)
		}
		return
	}

	SuccessResponseWithMessage(ctx, "Task deleted successfully", nil)
}
