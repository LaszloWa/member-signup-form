package handlers

import (
	"net/http"

	"spondtest/backend/internal/http/httputil"
)

func (h *Handler) GetFormDetails(w http.ResponseWriter, r *http.Request) {
	form := h.formService.GetFormDetails()

	httputil.WriteJSON(w, http.StatusOK, form)
}
