package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/xurl/internal/store"
	"github.com/xurl/util"
	"net/http"
	"strings"
)

func (app *application) createShortenedUrlHandler(w http.ResponseWriter, r *http.Request) {
	var payload store.URLMapper

	if err := util.ParseJSON(w, r, &payload); err != nil {
		_ = util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := util.Validate.Struct(&payload); err != nil {

		var ValidationErrors validator.ValidationErrors

		if errors.As(err, &ValidationErrors) {
			for _, e := range ValidationErrors {
				switch e.Tag() {
				case "required":
					_ = util.WriteError(w, http.StatusBadRequest, errors.New(strings.ToUpper(e.Field())+" is required"))
					return
				case "url":
					_ = util.WriteError(w, http.StatusBadRequest, errors.New("invalid "+strings.ToUpper(e.Field())))
					return
				default:
					_ = util.WriteError(w, http.StatusBadRequest, e)
					return
				}
			}
		}

		_ = util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !strings.HasPrefix(payload.Url, "http://") && !strings.HasPrefix(payload.Url, "https://") {
		payload.Url = "https://" + payload.Url
	}

	if payload.GeneratedUrl == "" {
		payload.GeneratedUrl = util.RandStringRunes(6)
	}

	_, err := app.store.Get(r.Context(), payload)

	if err != nil && !errors.Is(err, store.ErrURLNotFound) {
		_ = util.WriteError(w, http.StatusInternalServerError, err)
		return
	} else {
		if err == nil {
			_ = util.WriteError(w, http.StatusBadRequest, errors.New("URL already exists"))
			return
		}
	}

	if generatedUrl, err := app.store.Create(r.Context(), payload); err != nil {
		_ = util.WriteError(w, http.StatusInternalServerError, err)
	} else {
		generatedUrl = r.Host + "/" + generatedUrl
		_ = util.WriteJSON(w, http.StatusCreated, util.APIResponseBody{Message: "Shortened URL generated successfully", Data: generatedUrl})
	}
	return
}

func (app *application) getOriginalUrlHandler(w http.ResponseWriter, r *http.Request) {
	var payload store.URLMapper

	slug := chi.URLParam(r, "slug")

	payload.GeneratedUrl = slug

	if originalUrl, err := app.store.Get(r.Context(), payload); err != nil {
		_ = util.WriteError(w, http.StatusInternalServerError, err)
	} else {
		_ = util.WriteJSON(w, http.StatusOK, util.APIResponseBody{Message: "Original URL fetched successfully", Data: originalUrl})
	}
}
