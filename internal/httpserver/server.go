package httpserver

import (
    "context"
    "encoding/json"
    "log/slog"
    "net/http"
    "time"

    "github.com/NursultanKoshoev11/VRTrainingServer/internal/config"
    "github.com/NursultanKoshoev11/VRTrainingServer/internal/training"
)

type Server struct {
    cfg             config.Config
    logger          *slog.Logger
    trainingService *training.Service
    httpServer      *http.Server
}

func New(cfg config.Config, logger *slog.Logger, trainingService *training.Service) *Server {
    mux := http.NewServeMux()

    s := &Server{
        cfg:             cfg,
        logger:          logger,
        trainingService: trainingService,
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
    mux.HandleFunc("POST /training/sessions/start", s.handleStartTrainingSession)
    mux.HandleFunc("POST /training/sessions/{id}/events", s.handleAppendTrainingEvents)
    mux.HandleFunc("POST /training/sessions/{id}/complete", s.handleCompleteTrainingSession)
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

func (s *Server) handleStartTrainingSession(w http.ResponseWriter, r *http.Request) {
    var req training.StartSessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid_json", "Invalid request body.")
        return
    }

    resp, err := s.trainingService.StartSession(r.Context(), req)
    if err != nil {
        s.logger.Warn("start training session failed", "error", err)
        writeError(w, http.StatusBadRequest, "training_session_start_failed", err.Error())
        return
    }

    writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) handleAppendTrainingEvents(w http.ResponseWriter, r *http.Request) {
    sessionID := r.PathValue("id")
    var req training.EventBatchRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid_json", "Invalid request body.")
        return
    }

    if err := s.trainingService.AppendEvents(r.Context(), sessionID, req.Events); err != nil {
        s.logger.Warn("append training events failed", "session_id", sessionID, "error", err)
        writeError(w, http.StatusBadRequest, "training_event_upload_failed", err.Error())
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{"status": "events_saved"})
}

func (s *Server) handleCompleteTrainingSession(w http.ResponseWriter, r *http.Request) {
    sessionID := r.PathValue("id")
    var req training.CompleteSessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid_json", "Invalid request body.")
        return
    }

    if err := s.trainingService.CompleteSession(r.Context(), sessionID, req); err != nil {
        s.logger.Warn("complete training session failed", "session_id", sessionID, "error", err)
        writeError(w, http.StatusBadRequest, "training_session_complete_failed", err.Error())
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}
