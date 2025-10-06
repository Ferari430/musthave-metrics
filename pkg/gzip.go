package pkg

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(p []byte) (int, error) {
	return g.Writer.Write(p)
}

func GzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ---------- ВХОДЯЩИЙ ЗАПРОС ----------
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "invalid gzip body", http.StatusBadRequest)
				return
			}
			defer gzr.Close()
			r.Body = gzr
		}

		// ---------- ИСХОДЯЩИЙ ОТВЕТ ----------
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				http.Error(w, "failed to init gzip", http.StatusInternalServerError)
				return
			}

			w = &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzw,
			}

			defer gzw.Close()
		}

		// ---------- ВЫЗОВ ХЕНДЛЕРА ----------
		next.ServeHTTP(w, r)
	})
}
