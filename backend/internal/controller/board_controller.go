package controller

import (
	"github.com/gin-gonic/gin"

	"auth-project/internal/service"
)

// BoardController handles HTTP requests for board operations
type BoardController struct {
	boardService *service.BoardService
}

// NewBoardController creates a new board controller
func NewBoardController(boardService *service.BoardService) *BoardController {
	return &BoardController{
		boardService: boardService,
	}
}

// GetBoard handles GET /board request
// @Summary Get board
// @Description Get the board with columns and tasks
// @Tags Board
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /board [get]
func (c *BoardController) GetBoard(ctx *gin.Context) {
	board, err := c.boardService.GetBoard()
	if err != nil {
		InternalServerErrorResponse(ctx, err)
		return
	}

	SuccessResponse(ctx, board)
}
