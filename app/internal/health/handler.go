package health

import (
	"net/http"

	"github.com/mzeahmed/go-booking/internal/response"
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
