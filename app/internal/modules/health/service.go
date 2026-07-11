package health

// Service contains the business logic of the health module.
type Service struct{}

// NewService creates a new health service.
func NewService() *Service {
	return &Service{}
}

// Health returns the current application status.
func (s *Service) Health() Response {

	return Response{
		Status: "ok",
	}
}
