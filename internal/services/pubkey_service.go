package services

import (
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
)

func CreatePubkey(w http.ResponseWriter, r *http.Request) {
	payload := &payloads.PubkeyRequest{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), err))
		return
	}

	pubkeyDao, err := dao.GetPubkeyDao(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err))
		return
	}

	pubkey, err := pubkeyDao.Create(r.Context(), payload.Pubkey)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create pubkey", err))
		return
	}

	if err := render.Render(w, r, payloads.NewPubkeyResponse(pubkey)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "pubkey", err))
	}
}

func ListPubkeys(w http.ResponseWriter, r *http.Request) {
	pubkeyDao, err := dao.GetPubkeyDao(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err))
		return
	}

	pubkeys, err := pubkeyDao.List(r.Context(), 100, 0)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "list pubkeys", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewPubkeyListResponse(pubkeys)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list pubkeys", err))
		return
	}
}

func GetPubkey(w http.ResponseWriter, r *http.Request) {
	id, err := ParseUint64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "ID", err))
		return
	}

	pubkeyDao, err := dao.GetPubkeyDao(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err))
		return
	}

	pubkey, err := pubkeyDao.GetById(r.Context(), id)
	if err != nil {
		var e *dao.NoRowsError
		if errors.As(err, &e) {
			renderError(w, r, payloads.NewNotFoundError(r.Context(), err))
		} else {
			renderError(w, r, payloads.NewDAOError(r.Context(), "get pubkey by id", err))
		}
		return
	}

	if err := render.Render(w, r, payloads.NewPubkeyResponse(pubkey)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "pubkey", err))
	}
}

func DeletePubkey(w http.ResponseWriter, r *http.Request) {
	id, err := ParseUint64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "ID", err))
		return
	}

	pubkeyDao, err := dao.GetPubkeyDao(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err))
		return
	}

	err = pubkeyDao.Delete(r.Context(), id)
	if err != nil {
		var e *dao.MismatchAffectedError
		if errors.As(err, &e) {
			renderError(w, r, payloads.NewNotFoundError(r.Context(), err))
		} else {
			renderError(w, r, payloads.NewDAOError(r.Context(), "delete pubkey", err))
		}
		return
	}

	render.NoContent(w, r)
}
