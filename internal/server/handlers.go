package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/maxsnegir/url-shortener/internal/services"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type BaseHandler struct {
	logger *logrus.Logger
}

func (h *BaseHandler) TextResponse(w http.ResponseWriter, code int, data string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	if _, err := w.Write([]byte(data)); err != nil {
		h.logger.Error(err)
	}
}

func (h *BaseHandler) JSONResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error(err)
	}
}

type URLHandler struct {
	BaseHandler
	shortener services.URLService
}

func (h *URLHandler) SetURLTextHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if err != nil || len(url) == 0 {
			h.TextResponse(w, http.StatusUnprocessableEntity, "URL in request body is missing")
			return
		}
		shortURL, err := h.shortener.SetURL(string(url))
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
		requestData := &RequestData{}
		responseData := &ResponseData{}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil || requestData.URL == "" {
			responseData.ErrorMsg = "Wrong URL type or URL in body is missing"
			h.JSONResponse(w, http.StatusBadRequest, responseData)
			return
		}
		shortURL, err := h.shortener.SetURL(requestData.URL)
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

func (h *URLHandler) GetURLByIDHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		urlID := vars["urlID"]
		originalURL, err := h.shortener.GetURLByID(urlID)
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

func NewURLHandler(shortener services.URLService, logger *logrus.Logger) URLHandler {
	return URLHandler{
		BaseHandler: BaseHandler{logger: logger},
		shortener:   shortener,
	}
}
