package health

// Response represents the health endpoint response.
type Response struct {
	Status string `json:"status"`
}

// ProtectedResponse is returned by the protected health endpoint. It
// echoes back the identity extracted from the caller's access token, as
// a usage example of the Authenticate middleware.
type ProtectedResponse struct {
	Status string   `json:"status"`
	UserID int      `json:"user_id"`
	Roles  []string `json:"roles"`
}
