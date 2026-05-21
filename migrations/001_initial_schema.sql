CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    industry TEXT,
    country TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id),
    role TEXT NOT NULL,
    full_name TEXT NOT NULL,
    email TEXT,
    employee_id TEXT,
    password_hash TEXT,
    language TEXT NOT NULL DEFAULT 'en',
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(company_id, email),
    UNIQUE(company_id, employee_id)
);

CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id),
    device_name TEXT NOT NULL,
    platform TEXT NOT NULL,
    serial_number TEXT,
    app_version TEXT,
    last_seen_at TIMESTAMPTZ,
    last_sync_status TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(company_id, serial_number)
);

CREATE TABLE IF NOT EXISTS training_modules (
    id UUID PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    category TEXT,
    current_version TEXT NOT NULL,
    pass_score INTEGER NOT NULL DEFAULT 80,
    supported_languages TEXT[] NOT NULL DEFAULT ARRAY['en'],
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS training_assignments (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id),
    trainee_id UUID NOT NULL REFERENCES users(id),
    module_id UUID NOT NULL REFERENCES training_modules(id),
    assigned_by UUID REFERENCES users(id),
    status TEXT NOT NULL DEFAULT 'assigned',
    due_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS training_sessions (
    id UUID PRIMARY KEY,
    client_session_id TEXT NOT NULL,
    company_id UUID NOT NULL REFERENCES companies(id),
    trainee_id UUID NOT NULL REFERENCES users(id),
    assignment_id UUID REFERENCES training_assignments(id),
    module_id UUID NOT NULL REFERENCES training_modules(id),
    module_version TEXT NOT NULL,
    device_id UUID NOT NULL REFERENCES devices(id),
    started_at_client TIMESTAMPTZ,
    started_at_server TIMESTAMPTZ NOT NULL DEFAULT now(),
    ended_at_client TIMESTAMPTZ,
    synced_at_server TIMESTAMPTZ,
    duration_seconds INTEGER,
    score INTEGER,
    status TEXT NOT NULL DEFAULT 'started',
    language TEXT NOT NULL DEFAULT 'en',
    app_version TEXT,
    sync_status TEXT NOT NULL DEFAULT 'started',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(company_id, device_id, client_session_id),
    CHECK(score IS NULL OR (score >= 0 AND score <= 100)),
    CHECK(duration_seconds IS NULL OR duration_seconds >= 0)
);

CREATE TABLE IF NOT EXISTS training_events (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id),
    session_id UUID NOT NULL REFERENCES training_sessions(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    module_code TEXT,
    module_version TEXT,
    hazard_code TEXT,
    hazard_category TEXT,
    severity TEXT,
    action TEXT,
    is_correct BOOLEAN,
    time_offset_seconds INTEGER NOT NULL DEFAULT 0,
    metadata_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK(time_offset_seconds >= 0)
);

CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id),
    session_id UUID NOT NULL REFERENCES training_sessions(id),
    report_type TEXT NOT NULL,
    pdf_url TEXT,
    csv_url TEXT,
    generated_by UUID REFERENCES users(id),
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY,
    company_id UUID,
    user_id UUID,
    action TEXT NOT NULL,
    resource_type TEXT,
    resource_id TEXT,
    success BOOLEAN NOT NULL DEFAULT true,
    ip_address TEXT,
    user_agent TEXT,
    metadata_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_devices_company_id ON devices(company_id);
CREATE INDEX IF NOT EXISTS idx_training_sessions_company_id ON training_sessions(company_id);
CREATE INDEX IF NOT EXISTS idx_training_sessions_trainee_id ON training_sessions(trainee_id);
CREATE INDEX IF NOT EXISTS idx_training_sessions_module_id ON training_sessions(module_id);
CREATE INDEX IF NOT EXISTS idx_training_events_session_id ON training_events(session_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_company_id ON audit_logs(company_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
