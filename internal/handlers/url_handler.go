package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/maxsnegir/url-shortener/internal/auth"
	"github.com/maxsnegir/url-shortener/internal/services"
)

type URLHandler struct {
	BaseHandler
	shortener      services.URLService
	authentication auth.CookieAuthentication
}

func (h *URLHandler) SetURLTextHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Context().Value(UserIDKey).(string)
		url, err := io.ReadAll(r.Body)
		if err != nil || len(url) == 0 {
			h.TextResponse(w, http.StatusUnprocessableEntity, "URL in request body is missing")
			return
		}
		shortURL, err := h.shortener.SetShortURL(userToken, string(url))
		if err != nil {
			errMsg, statusCode := h.processSetURLError(err)
			h.TextResponse(w, statusCode, errMsg)
		}
		h.TextResponse(w, http.StatusCreated, shortURL)
	}
}

func (h *URLHandler) SetURLJSONHandler() http.HandlerFunc {

	type RequestData struct {
		URL string `json:"url"`
	}
	type ResponseData struct {
		Result   string `json:"result"`
		ErrorMsg string `json:"error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Context().Value(UserIDKey).(string)
		requestData := &RequestData{}
		responseData := &ResponseData{}

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil || requestData.URL == "" {
			// ToDo Сделать позже через errors wrap/unwrap
			responseData.ErrorMsg = "Wrong URL type or URL in body is missing"
			h.JSONResponse(w, http.StatusBadRequest, responseData)
			return
		}
		shortURL, err := h.shortener.SetShortURL(userToken, requestData.URL)
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

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		urlID := vars["urlID"]
		shortURL := fmt.Sprintf("%s/%s/", h.shortener.GetHostURL(), urlID)
		originalURL, err := h.shortener.GetOriginalURLByShort(shortURL)
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

func (h *URLHandler) GetAllUserURLs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Context().Value(UserIDKey).(string)
		userURLs, err := h.shortener.GetAllUserURLs(userToken)
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
