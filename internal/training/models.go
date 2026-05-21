package training

import "time"

type SessionStatus string

const (
    SessionStatusStarted     SessionStatus = "started"
    SessionStatusCompleted   SessionStatus = "completed"
    SessionStatusPendingSync SessionStatus = "pending_sync"
    SessionStatusFailed      SessionStatus = "failed"
)

type EventType string

const (
    EventModuleStarted       EventType = "module_started"
    EventModuleCompleted     EventType = "module_completed"
    EventHazardFound         EventType = "hazard_found"
    EventHazardMissed        EventType = "hazard_missed"
    EventUnsafeAction        EventType = "unsafe_action"
    EventPPESelected         EventType = "ppe_selected"
    EventDamagedPPESelected  EventType = "damaged_ppe_selected"
    EventChecklistCompleted  EventType = "checklist_item_completed"
    EventHintUsed            EventType = "hint_used"
    EventFeedbackShown       EventType = "feedback_shown"
)

type StartSessionRequest struct {
    ClientSessionID string    `json:"client_session_id"`
    TraineeID       string    `json:"trainee_id"`
    ModuleCode      string    `json:"module_code"`
    ModuleVersion   string    `json:"module_version"`
    DeviceID        string    `json:"device_id"`
    AppVersion      string    `json:"app_version"`
    Language        string    `json:"language"`
    StartedAtClient time.Time `json:"started_at_client"`
}

type StartSessionResponse struct {
    SessionID  string `json:"session_id"`
    ServerTime string `json:"server_time"`
    Status     string `json:"status"`
}

type TrainingEvent struct {
    EventType         EventType       `json:"event_type"`
    ModuleCode        string          `json:"module_code,omitempty"`
    ModuleVersion     string          `json:"module_version,omitempty"`
    HazardCode        string          `json:"hazard_code,omitempty"`
    HazardCategory    string          `json:"hazard_category,omitempty"`
    Severity          string          `json:"severity,omitempty"`
    Action            string          `json:"action,omitempty"`
    IsCorrect         bool            `json:"is_correct"`
    TimeOffsetSeconds int             `json:"time_offset_seconds"`
    Metadata          map[string]any  `json:"metadata,omitempty"`
}

type EventBatchRequest struct {
    Events []TrainingEvent `json:"events"`
}

type CompleteSessionRequest struct {
    ClientSessionID string    `json:"client_session_id"`
    DurationSeconds int       `json:"duration_seconds"`
    Score           int       `json:"score"`
    Status          string    `json:"status"`
    HazardsFound    int       `json:"hazards_found"`
    HazardsMissed   int       `json:"hazards_missed"`
    UnsafeActions   int       `json:"unsafe_actions"`
    HintsUsed       int       `json:"hints_used"`
    EndedAtClient   time.Time `json:"ended_at_client"`
}
