package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpx "spondtest/backend/internal/http"
	"spondtest/backend/internal/http/handlers"
	"spondtest/backend/internal/repository/memory"
	"spondtest/backend/internal/service"
)

func setupRouter() http.Handler {
	formService, err := service.NewFormService(service.DefaultFormDetails())
	if err != nil {
		panic(err)
	}

	repo := memory.NewSubmissionRepository()
	submissionService := service.NewSubmissionService(formService, repo)
	h := handlers.NewHandler(formService, submissionService)
	return httpx.NewRouter(h, httpx.RouterOptions{AllowedOrigins: []string{"http://localhost:5173"}})
}

func TestGetFormDetails(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/forms/public", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if body["clubId"] != "britsport" {
		t.Fatalf("expected clubId britsport, got %v", body["clubId"])
	}
}

func TestCreateSubmission_Success(t *testing.T) {
	router := setupRouter()
	payload := map[string]string{
		"name":         "Test User",
		"email":        "test@example.com",
		"phoneNumber":  "+47 12345678",
		"birthDate":    "1990-04-21",
		"memberTypeId": "8FE4113D4E4020E0DCF887803A886981",
		"clubId":       "britsport",
		"formId":       "B171388180BC457D9887AD92B6CCFC86",
	}
	encoded, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/forms/public/submissions", bytes.NewReader(encoded))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateSubmission_DuplicateReturnsConflict(t *testing.T) {
	router := setupRouter()
	payload := map[string]string{
		"name":         "Test User",
		"email":        "test@example.com",
		"phoneNumber":  "+47 12345678",
		"birthDate":    "1990-04-21",
		"memberTypeId": "8FE4113D4E4020E0DCF887803A886981",
		"clubId":       "britsport",
		"formId":       "B171388180BC457D9887AD92B6CCFC86",
	}
	encoded, _ := json.Marshal(payload)

	firstReq := httptest.NewRequest(http.MethodPost, "/api/v1/forms/public/submissions", bytes.NewReader(encoded))
	firstReq.Header.Set("Content-Type", "application/json")
	firstRR := httptest.NewRecorder()
	router.ServeHTTP(firstRR, firstReq)

	if firstRR.Code != http.StatusCreated {
		t.Fatalf("expected first submission 201, got %d body=%s", firstRR.Code, firstRR.Body.String())
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/api/v1/forms/public/submissions", bytes.NewReader(encoded))
	secondReq.Header.Set("Content-Type", "application/json")
	secondRR := httptest.NewRecorder()
	router.ServeHTTP(secondRR, secondReq)

	if secondRR.Code != http.StatusConflict {
		t.Fatalf("expected duplicate submission 409, got %d body=%s", secondRR.Code, secondRR.Body.String())
	}

	var body map[string]map[string]string
	if err := json.Unmarshal(secondRR.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON error envelope: %v", err)
	}

	if body["errors"]["submission"] == "" {
		t.Fatalf("expected submission duplicate error, got %v", body)
	}
}

func TestCreateSubmission_UnknownField(t *testing.T) {
	router := setupRouter()
	payload := []byte(`{"name":"Test User","email":"test@example.com","phoneNumber":"+47 12345678","birthDate":"1990-04-21","memberTypeId":"8FE4113D4E4020E0DCF887803A886981","clubId":"britsport","formId":"B171388180BC457D9887AD92B6CCFC86","unexpected":"nope"}`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/forms/public/submissions", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHealth(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "ok" {
		t.Fatalf("expected ok, got %q", rr.Body.String())
	}
}

func TestMethodNotAllowed_UsesEnvelope(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/forms/public", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}

	var body map[string]map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON error envelope: %v", err)
	}

	if body["errors"]["method"] == "" {
		t.Fatalf("expected method error in envelope, got %v", body)
	}
}

func TestNotFound_UsesEnvelope(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/not-a-route", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}

	var body map[string]map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON error envelope: %v", err)
	}

	if body["errors"]["route"] == "" {
		t.Fatalf("expected route error in envelope, got %v", body)
	}
}

func TestPreflight_AllowsConfiguredOrigin(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/forms/public/submissions", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}

	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("expected allow-origin header to echo request origin, got %q", got)
	}
}

func TestPreflight_RejectsUnknownOrigin(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/forms/public/submissions", nil)
	req.Header.Set("Origin", "https://evil.example")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}

	var body map[string]map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON error envelope: %v", err)
	}

	if body["errors"]["origin"] == "" {
		t.Fatalf("expected origin error in envelope, got %v", body)
	}
}
