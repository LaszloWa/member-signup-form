package handlers

import "spondtest/backend/internal/service"

type Handler struct {
	formService       *service.FormService
	submissionService *service.SubmissionService
}

func NewHandler(formService *service.FormService, submissionService *service.SubmissionService) *Handler {
	return &Handler{formService: formService, submissionService: submissionService}
}
