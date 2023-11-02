package v1

import (
	"context"
	"errors"
	"net/http"
	"usersegments/internal/service"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=segment.go -destination=mocks/segment_mock.go
type SegmentService interface {
	CreateSegment(ctx context.Context, input service.SegmentCreateInput) (int, error)
	DeleteSegment(ctx context.Context, input service.SegmentDeleteInput) error
}

type SegmentHandler struct {
	segmentService SegmentService
}

func NewSegmentHandler(segmentService SegmentService) *SegmentHandler {
	return &SegmentHandler{segmentService: segmentService}
}

func (h *SegmentHandler) newSegmentRoutes(g *gin.RouterGroup) {
	segments := g.Group("/segments")
	{
		segments.POST("/", h.createSegment)
		segments.DELETE("/", h.deleteSegment)

	}
}

type segmentCreateInput struct {
	Name    string `json:"name" binding:"required"`
	Percent int    `json:"percent"`
}

// @Summary Create Segment
// @Tags segment
// @Description create segment
// @ID create-segment
// @Accept json
// @Produce json
// @Param input body segmentCreateInput true "create segment"
// @Success 201 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments [post]
func (h *SegmentHandler) createSegment(c *gin.Context) {
	var inp segmentCreateInput
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}
	id, err := h.segmentService.CreateSegment(c.Request.Context(), service.SegmentCreateInput{
		Name:    inp.Name,
		Percent: inp.Percent,
	})
	if err != nil {
		if errors.Is(err, service.ErrSegmentAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, dataResponse{map[string]int{
		"ID": id,
	}})
}

type segmentDeleteInput struct {
	Name string `json:"name" binding:"required"`
}

// @Summary Delete Segment
// @Tags segment
// @Description delete segment
// @ID delete-segment
// @Accept json
// @Produce json
// @Param input body segmentDeleteInput true "delete segment"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments [delete]
func (h *SegmentHandler) deleteSegment(c *gin.Context) {
	var inp segmentDeleteInput

	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}

	err := h.segmentService.DeleteSegment(c.Request.Context(), service.SegmentDeleteInput{
		Name: inp.Name,
	})
	if err != nil {
		if errors.Is(err, service.ErrSegmentNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, dataResponse{"Segment deleted"})
}
