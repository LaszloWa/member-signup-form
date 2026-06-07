package handlers

import (
	"errors"
	"net/http"

	"spondtest/backend/internal/domain"
	"spondtest/backend/internal/http/httputil"
	"spondtest/backend/internal/service"
)

func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var input domain.SubmissionInput
	if err := httputil.DecodeStrictJSON(w, r, &input); err != nil {
		httputil.WriteErrors(w, http.StatusBadRequest, httputil.DecodeErrorFields(err))
		return
	}

	submission, validationErrors, err := h.submissionService.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateSubmission) {
			httputil.WriteErrors(w, http.StatusConflict, map[string]string{
				"submission": "a submission with the same form, club, email, phone number, and birth date already exists",
			})
			return
		}

		httputil.WriteErrors(w, http.StatusInternalServerError, map[string]string{"server": "internal server error"})
		return
	}

	if len(validationErrors) > 0 {
		httputil.WriteErrors(w, http.StatusBadRequest, validationErrors)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, submission)
}
