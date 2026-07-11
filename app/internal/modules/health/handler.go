package health

import (
	"net/http"

	"github.com/mzeahmed/gobooking/internal/response"
)

type Response struct {
	Status string `json:"status"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	response.JSON(
		w,
		http.StatusOK,
		Response{
			Status: "ok",
		},
	)
}
