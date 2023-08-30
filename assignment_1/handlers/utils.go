package handlers

import (
	"net/http"
	"strconv"
)

func CheckId(id string) (int, *HTTPError) {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return 0, &HTTPError{
			Status:  http.StatusBadRequest,
			Message: "Invalid Type ID",
		}
	}

	if idInt == 0 || idInt > len(data) {
		return 0, &HTTPError{
			Status:  http.StatusNotFound,
			Message: "Invalid ID",
		}
	}

	return idInt, nil
}
