package v1

type Deps struct {
	UserService      UserService
	SegmentService   SegmentService
	OperationService OperationService
}

type Handler struct {
	UserHandler      *UserHandler
	SegmentHandler   *SegmentHandler
	OperationHandler *OperationHandler
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		UserHandler:      NewUserHandler(deps.UserService),
		SegmentHandler:   NewSegmentHandler(deps.SegmentService),
		OperationHandler: NewOperationHandler(deps.OperationService),
	}
}
