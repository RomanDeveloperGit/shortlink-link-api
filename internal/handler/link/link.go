package link

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/errors"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/handler/link/dto"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/httpjson"
	"github.com/go-playground/validator/v10"
)

type LinkService interface {
	Create(ctx context.Context, fullURL string, ttlDays int) (*model.Link, error)
}

type handler struct {
	linkService LinkService
	validator   *validator.Validate
}

func NewHandler(linkService LinkService, validator *validator.Validate) *handler {
	return &handler{
		linkService: linkService,
		validator:   validator,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(
			w,
			http.StatusBadRequest,
			errors.ErrInvalidRequest,
			err.Error(),
		)

		return
	}

	if err := h.validator.Struct(req); err != nil {
		httpjson.WriteError(
			w,
			http.StatusBadRequest,
			errors.ErrInvalidValidation,
			err.Error(),
		)

		return
	}

	link, err := h.linkService.Create(r.Context(), req.FullURL, req.TTLDays)

	if err != nil {
		httpjson.WriteError(
			w,
			http.StatusInternalServerError,
			errors.ErrInternalServer,
			http.StatusText(http.StatusInternalServerError),
		)

		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, dto.NewCreateLinkResponse(link))
}
