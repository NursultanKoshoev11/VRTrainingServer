package training

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
    return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) StartSession(ctx context.Context, req StartSessionRequest) (string, error) {
    const query = `
        WITH trainee AS (
            SELECT id, company_id
            FROM users
            WHERE id = $1 AND is_active = true
        ), device_row AS (
            SELECT id, company_id
            FROM devices
            WHERE id = $2 AND is_active = true
        ), module_row AS (
            SELECT id, code, current_version
            FROM training_modules
            WHERE code = $3 AND is_active = true
        )
        INSERT INTO training_sessions (
            id,
            client_session_id,
            company_id,
            trainee_id,
            module_id,
            module_version,
            device_id,
            started_at_client,
            started_at_server,
            status,
            language,
            app_version,
            sync_status
        )
        SELECT
            gen_random_uuid(),
            $4,
            trainee.company_id,
            trainee.id,
            module_row.id,
            $5,
            device_row.id,
            $6,
            now(),
            'started',
            $7,
            $8,
            'started'
        FROM trainee, device_row, module_row
        WHERE trainee.company_id = device_row.company_id
        ON CONFLICT (company_id, device_id, client_session_id)
        DO UPDATE SET updated_at = now()
        RETURNING id::text;
    `

    var sessionID string
    err := r.pool.QueryRow(ctx, query,
        req.TraineeID,
        req.DeviceID,
        req.ModuleCode,
        req.ClientSessionID,
        req.ModuleVersion,
        req.StartedAtClient,
        req.Language,
        req.AppVersion,
    ).Scan(&sessionID)
    if err != nil {
        return "", fmt.Errorf("start training session: %w", err)
    }

    return sessionID, nil
}

func (r *PostgresRepository) AppendEvents(ctx context.Context, sessionID string, events []TrainingEvent) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin event transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    var companyID string
    if err := tx.QueryRow(ctx, `SELECT company_id::text FROM training_sessions WHERE id = $1`, sessionID).Scan(&companyID); err != nil {
        return fmt.Errorf("load session company: %w", err)
    }

    const insertEvent = `
        INSERT INTO training_events (
            id,
            company_id,
            session_id,
            event_type,
            module_code,
            module_version,
            hazard_code,
            hazard_category,
            severity,
            action,
            is_correct,
            time_offset_seconds,
            metadata_json
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
        );
    `

    for _, event := range events {
        metadata, err := json.Marshal(event.Metadata)
        if err != nil {
            return fmt.Errorf("marshal event metadata: %w", err)
        }

        if _, err := tx.Exec(ctx, insertEvent,
            companyID,
            sessionID,
            string(event.EventType),
            event.ModuleCode,
            event.ModuleVersion,
            event.HazardCode,
            event.HazardCategory,
            event.Severity,
            event.Action,
            event.IsCorrect,
            event.TimeOffsetSeconds,
            metadata,
        ); err != nil {
            return fmt.Errorf("insert training event: %w", err)
        }
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit event transaction: %w", err)
    }

    return nil
}

func (r *PostgresRepository) CompleteSession(ctx context.Context, sessionID string, req CompleteSessionRequest) error {
    const query = `
        UPDATE training_sessions
        SET
            ended_at_client = $2,
            synced_at_server = now(),
            duration_seconds = $3,
            score = $4,
            status = $5,
            sync_status = 'synced',
            updated_at = now()
        WHERE id = $1
          AND client_session_id = $6;
    `

    result, err := r.pool.Exec(ctx, query,
        sessionID,
        req.EndedAtClient,
        req.DurationSeconds,
        req.Score,
        req.Status,
        req.ClientSessionID,
    )
    if err != nil {
        return fmt.Errorf("complete training session: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("training session not found or client session id mismatch")
    }

    return nil
}
