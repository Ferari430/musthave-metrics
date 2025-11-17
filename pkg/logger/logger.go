package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger // классический zap
	ReqLogger   // zap для http
}

func InitLogger(level string) (*Logger, error) {

	var log *zap.Logger = zap.NewNop()

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	// singltone
	log = zl

	Logger := &Logger{
		Logger: log,
		ReqLogger: ReqLogger{
			logger: log,
		},
	}

	return Logger, nil
}

type responseData struct {
	statuscode int
	size       int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size += size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(statuscode int) {
	l.responseData.statuscode = statuscode

	l.ResponseWriter.WriteHeader(statuscode)

}

type ReqLogger struct {
	logger *zap.Logger
}

func (Rl *ReqLogger) RequestLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		responseData := &responseData{
			statuscode: 200,
			size:       0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		Rl.logger.Info("got incoming http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", duration),
			zap.Int("statuscode", lw.responseData.statuscode),
			zap.Int("size", lw.responseData.size),
		)

	})
}
