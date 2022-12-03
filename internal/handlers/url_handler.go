package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/maxsnegir/url-shortener/internal/auth"
	"github.com/maxsnegir/url-shortener/internal/services"
)

type URLHandler struct {
	BaseHandler
	shortener      services.URLService
	authentication auth.CookieAuthentication
}

func (h *URLHandler) getUserToken(ctx context.Context) string {
	userToken := ctx.Value(UserIDKey)
	if userToken == nil {
		return ""
	}
	return userToken.(string)
}

func (h *URLHandler) SetURLTextHandler() http.HandlerFunc {
	const timeout = 3 * time.Second

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		userToken := h.getUserToken(r.Context())
		url, err := io.ReadAll(r.Body)
		if err != nil || len(url) == 0 {
			h.TextResponse(w, http.StatusUnprocessableEntity, "URL in request body is missing")
			return
		}
		shortURL, err := h.shortener.SaveData(ctx, userToken, string(url))
		if err != nil {
			errMsg, statusCode := h.processSetURLError(err)
			h.TextResponse(w, statusCode, errMsg)
		}
		h.TextResponse(w, http.StatusCreated, shortURL)
	}
}

func (h *URLHandler) SetURLJSONHandler() http.HandlerFunc {
	const timeout = 3 * time.Second
	type RequestData struct {
		URL string `json:"url"`
	}
	type ResponseData struct {
		Result   string `json:"result"`
		ErrorMsg string `json:"error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		userToken := h.getUserToken(r.Context())
		requestData := &RequestData{}
		responseData := &ResponseData{}

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil || requestData.URL == "" {
			// ToDo Сделать позже через errors wrap/unwrap
			responseData.ErrorMsg = "Wrong URL type or URL in body is missing"
			h.JSONResponse(w, http.StatusBadRequest, responseData)
			return
		}
		shortURL, err := h.shortener.SaveData(ctx, userToken, requestData.URL)
		if err != nil {
			errMsg, statusCode := h.processSetURLError(err)
			responseData.ErrorMsg = errMsg
			h.JSONResponse(w, statusCode, responseData)
			return
		}
		responseData.Result = shortURL
		h.JSONResponse(w, http.StatusCreated, responseData)
	}
}

func (h *URLHandler) GetURLByIDHandler() http.HandlerFunc {
	const timeout = 3 * time.Second

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		vars := mux.Vars(r)
		urlID := vars["urlID"]
		shortURL := fmt.Sprintf("%s/%s/", h.shortener.GetHostURL(), urlID)
		originalURL, err := h.shortener.GetOriginalURL(ctx, shortURL)
		if err != nil {
			switch err.(type) {
			case services.OriginalURLNotFound:
				h.TextResponse(w, http.StatusNotFound, err.Error())
			default:
				h.logger.Error(err)
				h.TextResponse(w, http.StatusInternalServerError, InternalServerError.Error())
			}
			return
		}
		w.Header().Add("Location", originalURL)
		h.TextResponse(w, http.StatusTemporaryRedirect, "")
	}
}

func (h *URLHandler) GetUserURLs() http.HandlerFunc {
	const timeout = 3 * time.Second

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		userToken := h.getUserToken(r.Context())
		userURLs, err := h.shortener.GetUserURLs(ctx, userToken)
		if err != nil {
			h.TextResponse(w, http.StatusInternalServerError, InternalServerError.Error())
			return
		}
		if len(userURLs) == 0 {
			h.JSONResponse(w, http.StatusNoContent, userURLs)
			return
		}
		h.JSONResponse(w, http.StatusOK, userURLs)
	}
}

func (h *URLHandler) SaveDataBatch() http.HandlerFunc {
	const timeout = 3 * time.Second

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		userToken := h.getUserToken(r.Context())
		var requestData []services.URLDataBatchRequest

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			h.JSONResponse(w, http.StatusBadRequest, errors.New("wrong request"))
			return
		}
		responseData, err := h.shortener.SaveDataBatch(ctx, userToken, requestData)
		if err != nil {
			h.processSetURLError(err)
			errMsg, statusCode := h.processSetURLError(err)
			h.JSONResponse(w, statusCode, errMsg)
			return
		}
		h.JSONResponse(w, http.StatusCreated, responseData)
	}
}

func (h *URLHandler) Ping() http.HandlerFunc {
	const timeout = 3 * time.Second

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		if err := h.shortener.Ping(ctx); err != nil {
			h.TextResponse(w, http.StatusInternalServerError, "")
			return
		}
		h.TextResponse(w, http.StatusOK, "")
	}
}

func (h *URLHandler) processSetURLError(err error) (string, int) {
	var errMsg string
	var statusCode int
	switch err.(type) {
	case services.URLIsNotValidError:
		errMsg = err.Error()
		statusCode = http.StatusBadRequest
	default:
		errMsg = InternalServerError.Error()
		statusCode = http.StatusInternalServerError
		h.logger.Error(err)
	}
	return errMsg, statusCode
}

func NewURLHandler(shortener services.URLService, auth auth.CookieAuthentication, logger *logrus.Logger) URLHandler {
	return URLHandler{
		BaseHandler:    BaseHandler{logger: logger},
		shortener:      shortener,
		authentication: auth,
	}
}
