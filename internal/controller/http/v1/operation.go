package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"usersegments/internal/entity"
	"usersegments/internal/service"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=operation.go -destination=mocks/operation_mock.go
type OperationService interface {
	GetOperationHistory(ctx context.Context, user_id int, input service.GetOperationHistoryInput) ([]entity.Operation, error)
}
type OperationHandler struct {
	operationService OperationService
}

func NewOperationHandler(operationService OperationService) *OperationHandler {
	return &OperationHandler{operationService: operationService}
}

func (h *OperationHandler) newOperationRoutes(g *gin.RouterGroup) {
	operations := g.Group("/operations")
	{
		operations.POST("/:id", h.getOperationHistory)
	}
}

type getOperationHistoryInput struct {
	Month string `json:"month" binding:"required"`
	Year  string `json:"year" binding:"required"`
}

// @Summary Get Operation History
// @Tags operation
// @Description get operaion history
// @ID get-operaion-history
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Param input body getOperationHistoryInput true "get operaion history"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /operations/{id} [post]
func (o *OperationHandler) getOperationHistory(c *gin.Context) {
	var inp getOperationHistoryInput
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidUserId)
		return
	}
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}

	operations, err := o.operationService.GetOperationHistory(c.Request.Context(), id, service.GetOperationHistoryInput{
		Month: inp.Month,
		Year:  inp.Year,
	})
	if err != nil {
		if errors.Is(err, service.ErrOperationHistoryNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{operations})
}
