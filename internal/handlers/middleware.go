package handlers

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/maxsnegir/url-shortener/internal/auth"
)

const UserID = "userID"

func (h *URLHandler) CookieAuthenticationMiddleware(next http.Handler) http.Handler {

	setTokenToCookie := func(w http.ResponseWriter, token string) {
		http.SetCookie(w, &http.Cookie{
			Name:  auth.AuthorizationCookieName,
			Value: token,
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIP := strings.Split(r.RemoteAddr, ":")[0]
		token, err := r.Cookie(auth.AuthorizationCookieName)
		if err != nil {
			setTokenToCookie(w, h.authentication.CreateToken())
		} else {
			if _, err := h.authentication.ParseToken(token.Value); err != nil {
				setTokenToCookie(w, h.authentication.CreateToken())
			}
		}
		ctx := context.WithValue(r.Context(), UserID, userIP)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *BaseHandler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		h.logger.Infof("%s :: %s", r.RequestURI, r.Method)
	})
}

func (h *BaseHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		defer func() {
			if err := gz.Close(); err != nil {
				h.logger.Error(err)
			}
		}()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (h *BaseHandler) UnzipMiddleware(next http.Handler) http.Handler {
	var reader io.ReadCloser
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				if err := gz.Close(); err != nil {
					h.logger.Error(err)
				}
			}()
			reader = gz
		} else {
			reader = r.Body
		}
		r.Body = reader
		next.ServeHTTP(w, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
