package service

type Deps struct {
	UserRepo      UserRepo
	SegmentRepo   SegmentRepo
	OperationRepo OperationRepo
}

type Service struct {
	UserService      *UserService
	SegmentService   *SegmentService
	OperationService *OperationService
}

func NewService(deps Deps) *Service {
	return &Service{
		UserService:      NewUserService(deps.UserRepo),
		SegmentService:   NewSegmentService(deps.SegmentRepo),
		OperationService: NewOperationService(deps.OperationRepo),
	}
}
