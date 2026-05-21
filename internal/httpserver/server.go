package httpserver

import (
    "context"
    "encoding/json"
    "log/slog"
    "net/http"
    "time"

    "github.com/NursultanKoshoev11/VRTrainingServer/internal/config"
)

type Server struct {
    cfg        config.Config
    logger     *slog.Logger
    httpServer *http.Server
}

func New(cfg config.Config, logger *slog.Logger) *Server {
    mux := http.NewServeMux()

    s := &Server{
        cfg:    cfg,
        logger: logger,
    }

    s.registerRoutes(mux)

    s.httpServer = &http.Server{
        Addr:              cfg.HTTPAddress(),
        Handler:           requestLogger(logger, securityHeaders(mux)),
        ReadHeaderTimeout: 10 * time.Second,
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      30 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    return s
}

func (s *Server) Start() error {
    return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.httpServer.Shutdown(ctx)
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
    mux.HandleFunc("GET /ping", s.handlePing)
    mux.HandleFunc("GET /health", s.handleHealth)
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]any{
        "status": "healthy",
        "app_env": s.cfg.AppEnv,
        "time": time.Now().UTC().Format(time.RFC3339),
    })
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}
