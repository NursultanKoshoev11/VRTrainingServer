package training

import (
    "context"
    "errors"
    "time"
)

var ErrInvalidSessionPayload = errors.New("invalid training session payload")

type Repository interface {
    StartSession(ctx context.Context, req StartSessionRequest) (string, error)
    AppendEvents(ctx context.Context, sessionID string, events []TrainingEvent) error
    CompleteSession(ctx context.Context, sessionID string, req CompleteSessionRequest) error
}

type Service struct {
    repo Repository
    clock func() time.Time
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo, clock: func() time.Time { return time.Now().UTC() }}
}

func (s *Service) StartSession(ctx context.Context, req StartSessionRequest) (StartSessionResponse, error) {
    if req.ClientSessionID == "" || req.TraineeID == "" || req.ModuleCode == "" || req.DeviceID == "" {
        return StartSessionResponse{}, ErrInvalidSessionPayload
    }

    sessionID, err := s.repo.StartSession(ctx, req)
    if err != nil {
        return StartSessionResponse{}, err
    }

    return StartSessionResponse{
        SessionID: sessionID,
        ServerTime: s.clock().Format(time.RFC3339),
        Status: string(SessionStatusStarted),
    }, nil
}

func (s *Service) AppendEvents(ctx context.Context, sessionID string, events []TrainingEvent) error {
    if sessionID == "" || len(events) == 0 {
        return ErrInvalidSessionPayload
    }
    return s.repo.AppendEvents(ctx, sessionID, events)
}

func (s *Service) CompleteSession(ctx context.Context, sessionID string, req CompleteSessionRequest) error {
    if sessionID == "" || req.ClientSessionID == "" || req.DurationSeconds < 0 || req.Score < 0 || req.Score > 100 {
        return ErrInvalidSessionPayload
    }
    return s.repo.CompleteSession(ctx, sessionID, req)
}
