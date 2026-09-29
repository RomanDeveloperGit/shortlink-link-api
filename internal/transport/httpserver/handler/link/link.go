package link

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/apperror"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/middleware/observability"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/requestvalidation"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/response"
)

type LinkService interface {
	Create(ctx context.Context, fullURL string, ttlDays int) (*model.Link, error)
	GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error)
	GetByID(ctx context.Context, id int) (*model.Link, error)
	Visit(ctx context.Context, shortCode string) (*model.Link, error)
}

type handler struct {
	linkService LinkService
	validator   *validator.Validate
	logger      *slog.Logger
}

func NewHandler(
	linkService LinkService,
	validator *validator.Validate,
	logger *slog.Logger,
) *handler {
	return &handler{
		linkService: linkService,
		validator:   validator,
		logger:      logger,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		observability.NewRequestLogger(r.Context(), h.logger).Error(
			"failed to decode request body",
			slog.Any("error", err),
		)

		response.RespondStatus(w, http.StatusBadRequest)

		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.RespondError(
			w,
			http.StatusUnprocessableEntity,
			&response.Error{
				Code:    apperror.ErrValidation.Code,
				Message: apperror.ErrValidation.Message,
				Details: requestvalidation.ParseErrors(err),
			},
		)

		return
	}

	link, err := h.linkService.Create(r.Context(), req.FullURL, req.TTLDays)

	if err != nil {
		observability.NewRequestLogger(r.Context(), h.logger).Error(
			"failed to create link",
			slog.Any("error", err),
		)

		response.RespondStatus(w, http.StatusInternalServerError)

		return
	}

	response.RespondJSON(w, http.StatusCreated, NewLinkResponse(link))
}

func (h *handler) GetByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Query().Get("short_code")

	if shortCode == "" {
		response.RespondStatus(w, http.StatusBadRequest)

		return
	}

	link, err := h.linkService.GetByShortCode(r.Context(), shortCode)

	if err != nil {
		observability.NewRequestLogger(r.Context(), h.logger).Error(
			"failed to get link by short code",
			slog.Any("error", err),
		)

		if errors.Is(err, apperror.ErrLinkNotFound) {
			response.RespondError(w, http.StatusNotFound, &response.Error{
				Code:    apperror.ErrLinkNotFound.Code,
				Message: apperror.ErrLinkNotFound.Message,
			})

			return
		}

		response.RespondStatus(w, http.StatusInternalServerError)

		return
	}

	response.RespondJSON(w, http.StatusOK, NewLinkResponse(link))
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		response.RespondStatus(w, http.StatusBadRequest)

		return
	}

	link, err := h.linkService.GetByID(r.Context(), id)

	if err != nil {
		observability.NewRequestLogger(r.Context(), h.logger).Error(
			"failed to get link by id",
			slog.Any("error", err),
		)

		if errors.Is(err, apperror.ErrLinkNotFound) {
			response.RespondError(w, http.StatusNotFound, &response.Error{
				Code:    apperror.ErrLinkNotFound.Code,
				Message: apperror.ErrLinkNotFound.Message,
			})

			return
		}

		response.RespondStatus(w, http.StatusInternalServerError)

		return
	}

	response.RespondJSON(w, http.StatusOK, NewLinkResponse(link))
}

func (h *handler) Visit(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	link, err := h.linkService.Visit(r.Context(), shortCode)

	if err != nil {
		observability.NewRequestLogger(r.Context(), h.logger).Error(
			"failed to visit link",
			slog.Any("error", err),
		)

		if errors.Is(err, apperror.ErrLinkNotFound) {
			response.RespondError(w, http.StatusNotFound, &response.Error{
				Code:    apperror.ErrLinkNotFound.Code,
				Message: apperror.ErrLinkNotFound.Message,
			})

			return
		}

		response.RespondStatus(w, http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, link.FullURL, http.StatusFound)
}
