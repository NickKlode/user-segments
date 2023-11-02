package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"usersegments/internal/service"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=user.go -destination=mocks/user_mock.go
type UserService interface {
	CreateUser(ctx context.Context, input service.UserCreateInput) (int, error)
	AddUserToSegment(ctx context.Context, input service.AddUserToSegmentsInput) error
	DeleteUserFromSegment(ctx context.Context, input service.DeleteUserFromSegmentsInput) error
	GetAllUserSegments(ctx context.Context, user_id int) ([]service.AllUserSegmentsOutput, error)
}

type UserHandler struct {
	UserService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) newUserRoutes(g *gin.RouterGroup) {
	users := g.Group("/users")
	{
		users.POST("/", h.createUser)
		users.POST("/:id", h.addUserToSegment)
		users.DELETE("/:id", h.deleteUserFromSegment)
		users.GET("/:id", h.getAllUserSegments)

	}
}

type userCreateInput struct {
	Name string `json:"name" binding:"required"`
}

// @Summary Create User
// @Tags user
// @Description create user
// @ID create-user
// @Accept json
// @Produce json
// @Param input body userCreateInput true "user creation"
// @Success 201 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users [post]
func (h *UserHandler) createUser(c *gin.Context) {
	var inp userCreateInput
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}
	id, err := h.UserService.CreateUser(c.Request.Context(), service.UserCreateInput{
		Name: inp.Name,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
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

type addUserToSegmentInput struct {
	Segments []string `json:"segments" binding:"required"`
	Timeout  int      `json:"timeout"`
}

// @Summary Add User To Segment
// @Tags user
// @Description add user to segment
// @ID add-user-to-segment
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Param input body addUserToSegmentInput true "add user to segment"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/{id} [post]
func (h *UserHandler) addUserToSegment(c *gin.Context) {
	var inp addUserToSegmentInput
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidUserId)
		return
	}
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}

	err = h.UserService.AddUserToSegment(c.Request.Context(), service.AddUserToSegmentsInput{
		UserID:   id,
		Segments: inp.Segments,
		Timeout:  inp.Timeout,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyInSegments) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{"User Added"})
}

type deleteUserFromSegmentInput struct {
	Segments []string `json:"segments" binding:"required"`
}

// @Summary Delete User To Segment
// @Tags user
// @Description delete user to segment
// @ID delete-user-to-segment
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Param input body deleteUserFromSegmentInput true "delete user to segment"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) deleteUserFromSegment(c *gin.Context) {
	var inp deleteUserFromSegmentInput
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidUserId)
		return
	}
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidInput)
		return
	}
	err = h.UserService.DeleteUserFromSegment(c.Request.Context(), service.DeleteUserFromSegmentsInput{
		UserID:   id,
		Segments: inp.Segments,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotInSegment) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, dataResponse{"User Deleted"})
}

// @Summary Get All User Segments
// @Tags user
// @Description get all user segments
// @ID get-all-user-segments
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/{id} [get]
func (h *UserHandler) getAllUserSegments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errInvalidUserId)
		return
	}
	segments, err := h.UserService.GetAllUserSegments(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNoSegments) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, dataResponse{segments})
}
