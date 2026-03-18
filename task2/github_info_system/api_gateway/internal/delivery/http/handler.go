package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"task2/github_info_system/api_gateway/internal/delivery/grpc"
	"task2/github_info_system/proto/collector"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	collectorClient *grpc.Client
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewHandler(collectorClient *grpc.Client) *Handler {
	return &Handler{
		collectorClient: collectorClient,
	}
}

func (h *Handler) GetRepositoryInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repo := vars["repo"]

	if owner == "" || repo == "" {
		h.sendError(w, "Invalid request", "owner and repo are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	resp, err := h.collectorClient.GetRepositoryInfo(ctx, owner, repo)
	if err != nil {
		h.handleGRPCError(w, err)
		return
	}

	h.sendSuccess(w, resp)
}

func (h *Handler) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger-ui/index.html")
}

func (h *Handler) SwaggerJSON(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.json")
}

func (h *Handler) handleGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		h.sendError(w, "Internal error", err.Error(), http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		h.sendError(w, "Not found", st.Message(), http.StatusNotFound)
	case codes.InvalidArgument:
		h.sendError(w, "Bad request", st.Message(), http.StatusBadRequest)
	case codes.Internal:
		h.sendError(w, "Internal error", st.Message(), http.StatusInternalServerError)
	default:
		h.sendError(w, "Unknown error", st.Message(), http.StatusInternalServerError)
	}
}

func (h *Handler) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) sendError(w http.ResponseWriter, error, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   error,
		Message: message,
		Code:    code,
	})
}
