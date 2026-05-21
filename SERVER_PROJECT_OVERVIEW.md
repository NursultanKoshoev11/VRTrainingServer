# VRTrainingServer — Production Backend Documentation

## 1. Project Summary

**VRTrainingServer** is the production backend for **BuildSafe VR**, a B2B construction safety training platform powered by an Unreal Engine VR client. The backend is responsible for receiving training data from VR devices, storing it reliably, enforcing company-level data isolation, managing users and devices, supporting dashboard APIs, generating reports, and providing a secure foundation for commercial customer use.

This backend is not a temporary prototype server. It must be designed from the beginning as a production service that can support real companies, real trainees, real training records, and real reporting workflows.

The server turns the VR application into a complete business product. The VR headset creates the training experience, but the backend makes that experience measurable, reviewable, auditable, and reportable.

Production expectations:

- secure authentication;
- role-based authorization;
- multi-tenant company isolation;
- reliable training session storage;
- idempotent offline sync handling;
- structured event ingestion;
- dashboard-ready APIs;
- PDF/CSV reports;
- audit logging;
- database migrations;
- Docker-based deployment;
- environment-based configuration;
- monitoring-friendly health checks;
- safe error handling;
- production HTTPS deployment.

---

## 2. Production Server Vision

The backend must become the operational system of record for VR construction safety training.

For every company, the server should answer:

- who completed training;
- which module was completed;
- which version of the module was used;
- when training started and ended;
- how long the session lasted;
- what score was achieved;
- whether the trainee passed or failed;
- which hazards were found;
- which hazards were missed;
- which unsafe actions were performed;
- which users need retraining;
- which modules create the most mistakes;
- which devices are actively syncing;
- whether reports can be exported for internal records.

The backend must support real operational usage, not only development testing.

---

## 3. Role of This Repository

This repository contains the server-side application for BuildSafe VR.

The server handles:

1. Authentication.
2. Authorization.
3. Company management.
4. Tenant isolation.
5. User and trainee management.
6. Trainer and safety manager workflows.
7. Device registration.
8. Module metadata and versioning.
9. Training assignments.
10. Training session lifecycle.
11. Training event ingestion.
12. Session completion.
13. Score and result storage.
14. Dashboard API.
15. PDF report generation.
16. CSV export.
17. Audit logging.
18. Offline sync support.
19. Duplicate upload prevention.
20. Health checks and production monitoring support.

The backend must expose APIs to the VR client and the web dashboard. The database must never be accessed directly by the VR client.

---

## 4. System Architecture

Production architecture:

```text
Meta Quest / PC VR Client
        ↓ HTTPS REST API
VRTrainingServer
        ↓
PostgreSQL Database
        ↓
Dashboard API / Reports / Exports
```

Recommended deployment architecture:

```text
Internet
  ↓
HTTPS Reverse Proxy / WAF
  ↓
VRTrainingServer Container
  ↓
PostgreSQL
  ↓
Backup Storage / Report Storage
```

All production traffic must use HTTPS. Secrets must be managed through environment variables or a secure secret management system, not committed into the repository.

---

## 5. Production Data Flow

## 5.1 Session Access

1. Trainer assigns a module to a trainee or group.
2. Trainee starts the VR application.
3. Trainee authenticates through login, PIN, QR code, or assigned training code.
4. VR client requests available assignments or validates the training code.
5. Server confirms whether the trainee can start the module.

## 5.2 Session Start

1. VR client starts the module.
2. Client creates a local `client_session_id`.
3. Client sends session start request.
4. Server creates a server-side session record.
5. Server returns `session_id`.
6. Client continues training.

## 5.3 Event Recording

1. VR client records training events locally.
2. Client uploads event batches during the session or after completion.
3. Server validates the session, company, module version, and event structure.
4. Server stores events in `training_events`.

## 5.4 Session Completion

1. VR client completes the module.
2. Client saves final result locally.
3. Client sends completion request to server.
4. Server stores final score, pass/fail status, duration, and summary.
5. Server marks session as completed.
6. Dashboard can display the session result.
7. Manager can export PDF or CSV report.

## 5.5 Offline Sync

If the VR client is offline:

1. Client completes the session locally.
2. Client keeps data in pending sync queue.
3. Client retries upload when internet returns.
4. Server uses idempotency rules to prevent duplicate records.
5. Server confirms successful sync.
6. Client marks local session as uploaded.

---

## 6. Production Authentication

The server must support secure authentication for dashboard users and controlled access for VR training devices.

Production authentication methods may include:

- email/password login for company admins and managers;
- JWT or secure session tokens;
- refresh token or session renewal flow;
- PIN or QR-based trainee access;
- device registration token;
- future SSO/SAML for enterprise customers.

Password rules:

- passwords must never be stored in plain text;
- use a strong password hashing algorithm;
- protect login endpoints with rate limiting;
- return safe error messages;
- log failed login attempts in audit logs.

---

## 7. Production Authorization and Roles

The server must implement role-based access control.

Recommended roles:

## 7.1 Platform Admin

Can manage platform-level settings if needed.

Permissions:

- manage companies;
- view platform health;
- support customer onboarding;
- access operational tools only when authorized.

## 7.2 Company Admin

Manages one company account.

Permissions:

- manage company profile;
- manage users;
- manage devices;
- manage module assignments;
- view all company training data;
- export company reports.

## 7.3 Safety Manager

Manages training results and safety review.

Permissions:

- assign modules;
- view trainee results;
- view analytics;
- export reports;
- identify retraining needs.

## 7.4 Trainer

Runs training sessions and reviews assigned groups.

Permissions:

- view assigned trainees;
- start or support training sessions;
- view assigned session results.

## 7.5 Trainee

Completes training.

Permissions:

- access assigned modules;
- view own completion result if enabled.

Every API endpoint must check role permissions and company boundaries.

---

## 8. Multi-Tenant Company Isolation

BuildSafe VR is a B2B product. The backend must support multiple companies safely.

Production rule:

**A user from Company A must never access data from Company B.**

Important tables should include `company_id` where relevant:

- users;
- devices;
- training assignments;
- training sessions;
- reports;
- audit logs.

All queries must be scoped by company context unless a platform admin action explicitly requires broader access.

---

## 9. Device Management

The server must track VR devices used by companies.

Device records should include:

- device ID;
- company ID;
- device name;
- platform;
- serial number or assigned hardware ID if available;
- app version;
- last seen timestamp;
- last sync status;
- active/inactive status.

Device management is important because training centers may use multiple headsets. Managers need to know whether devices are syncing correctly and whether they are running the correct app version.

---

## 10. Training Modules and Versioning

The backend must store module metadata and module versions.

Each module should include:

- module ID;
- module code;
- title;
- description;
- category;
- current version;
- pass score;
- supported languages;
- active/inactive status.

Module versioning is required because reports must remain understandable after module content changes.

Example module codes:

- `ppe_site_entry_001`
- `fall_hazard_recognition_001`
- `construction_hazard_hunt_001`
- `electrical_hazard_awareness_001`
- `trench_safety_awareness_001`

---

## 11. Training Assignments

The server should support assigning modules to trainees or groups.

Assignment data may include:

- assignment ID;
- company ID;
- trainee ID or group ID;
- module ID;
- required completion date;
- assigned by;
- status;
- created at;
- completed at.

Assignment support helps companies manage training operations instead of relying only on manual headset usage.

---

## 12. Training Sessions

A training session is one attempt by one trainee to complete one training module.

Each session should store:

- server session ID;
- client session ID;
- company ID;
- trainee ID;
- assignment ID if available;
- module ID;
- module version;
- device ID;
- started at client;
- started at server;
- ended at client;
- synced at server;
- duration seconds;
- score;
- pass/fail status;
- language;
- app version;
- sync status;
- created at;
- updated at.

The training session is the main record for reports and dashboard results.

---

## 13. Training Events

Training events explain what happened inside the VR module.

Examples:

- hazard found;
- hazard missed;
- PPE selected;
- damaged PPE selected;
- damaged PPE rejected;
- hint used;
- unsafe area entered;
- checklist item completed;
- wrong object selected;
- critical error triggered.

Each event should store:

- event ID;
- company ID;
- session ID;
- event type;
- module code;
- module version;
- hazard code;
- hazard category;
- severity;
- action;
- correctness;
- time offset seconds;
- metadata JSON;
- created at.

The event system makes reports more useful than a simple final score.

---

## 14. Scoring and Validation

The VR client may calculate the score locally for immediate feedback, but the server must store and validate the result.

Production requirements:

- score must be tied to module version;
- pass threshold must be stored;
- server must reject impossible or malformed results;
- server must prevent duplicate completion;
- score data must remain readable in historical reports.

Initial scoring model:

- start score: 100;
- missed critical hazard: -15;
- missed medium hazard: -8;
- missed minor hazard: -3;
- unsafe action: -20;
- hint used: -5;
- timeout: -10.

Default pass threshold: 80.

---

## 15. Production API Endpoints

## 15.1 Health

```http
GET /ping
GET /health
```

## 15.2 Authentication

```http
POST /auth/register
POST /auth/login
POST /auth/logout
POST /auth/refresh
GET  /auth/me
```

## 15.3 Companies

```http
GET   /companies/current
PATCH /companies/current
```

## 15.4 Users

```http
GET    /users
POST   /users
GET    /users/{id}
PATCH  /users/{id}
DELETE /users/{id}
```

## 15.5 Devices

```http
POST  /devices/register
GET   /devices
GET   /devices/{id}
PATCH /devices/{id}
POST  /devices/{id}/heartbeat
```

## 15.6 Training Modules

```http
GET  /training/modules
GET  /training/modules/{id}
GET  /training/modules/by-code/{code}
```

## 15.7 Assignments

```http
GET  /training/assignments
POST /training/assignments
GET  /training/assignments/{id}
POST /training/assignments/{id}/cancel
```

## 15.8 Training Sessions

```http
POST /training/sessions/start
POST /training/sessions/{id}/events
POST /training/sessions/{id}/complete
GET  /training/sessions
GET  /training/sessions/{id}
```

## 15.9 Offline Sync

```http
POST /sync/training-session
POST /sync/training-session/batch
```

## 15.10 Reports

```http
GET /reports/sessions/{id}/pdf
GET /reports/sessions/export.csv
GET /reports/company/summary
```

## 15.11 Audit Logs

```http
GET /audit-logs
```

Access to audit logs must be restricted.

---

## 16. Database Schema Draft

## 16.1 companies

```sql
id UUID PRIMARY KEY
name TEXT NOT NULL
industry TEXT
country TEXT
status TEXT NOT NULL
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

## 16.2 users

```sql
id UUID PRIMARY KEY
company_id UUID REFERENCES companies(id)
role TEXT NOT NULL
full_name TEXT NOT NULL
email TEXT
employee_id TEXT
password_hash TEXT
language TEXT
is_active BOOLEAN NOT NULL DEFAULT true
last_login_at TIMESTAMP
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

## 16.3 devices

```sql
id UUID PRIMARY KEY
company_id UUID REFERENCES companies(id)
device_name TEXT
platform TEXT
serial_number TEXT
app_version TEXT
last_seen_at TIMESTAMP
last_sync_status TEXT
is_active BOOLEAN NOT NULL DEFAULT true
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

## 16.4 training_modules

```sql
id UUID PRIMARY KEY
code TEXT UNIQUE NOT NULL
title TEXT NOT NULL
description TEXT
category TEXT
current_version TEXT NOT NULL
pass_score INTEGER NOT NULL DEFAULT 80
supported_languages TEXT[]
is_active BOOLEAN NOT NULL DEFAULT true
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

## 16.5 training_assignments

```sql
id UUID PRIMARY KEY
company_id UUID REFERENCES companies(id)
trainee_id UUID REFERENCES users(id)
module_id UUID REFERENCES training_modules(id)
assigned_by UUID REFERENCES users(id)
status TEXT NOT NULL
due_at TIMESTAMP
completed_at TIMESTAMP
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

## 16.6 training_sessions

```sql
id UUID PRIMARY KEY
client_session_id TEXT NOT NULL
company_id UUID REFERENCES companies(id)
trainee_id UUID REFERENCES users(id)
assignment_id UUID REFERENCES training_assignments(id)
module_id UUID REFERENCES training_modules(id)
module_version TEXT NOT NULL
device_id UUID REFERENCES devices(id)
started_at_client TIMESTAMP
started_at_server TIMESTAMP NOT NULL
ended_at_client TIMESTAMP
synced_at_server TIMESTAMP
duration_seconds INTEGER
score INTEGER
status TEXT
language TEXT
app_version TEXT
sync_status TEXT
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
UNIQUE(company_id, device_id, client_session_id)
```

## 16.7 training_events

```sql
id UUID PRIMARY KEY
company_id UUID REFERENCES companies(id)
session_id UUID REFERENCES training_sessions(id)
event_type TEXT NOT NULL
module_code TEXT
module_version TEXT
hazard_code TEXT
hazard_category TEXT
severity TEXT
action TEXT
is_correct BOOLEAN
time_offset_seconds INTEGER
metadata_json JSONB
created_at TIMESTAMP NOT NULL
```

## 16.8 reports

```sql
id UUID PRIMARY KEY
company_id UUID REFERENCES companies(id)
session_id UUID REFERENCES training_sessions(id)
report_type TEXT NOT NULL
pdf_url TEXT
csv_url TEXT
generated_by UUID REFERENCES users(id)
generated_at TIMESTAMP NOT NULL
created_at TIMESTAMP NOT NULL
```

## 16.9 audit_logs

```sql
id UUID PRIMARY KEY
company_id UUID
user_id UUID
action TEXT NOT NULL
resource_type TEXT
resource_id TEXT
success BOOLEAN
ip_address TEXT
user_agent TEXT
metadata_json JSONB
created_at TIMESTAMP NOT NULL
```

---

## 17. API Payload Examples

## 17.1 Start Session Request

```json
{
  "client_session_id": "8d1d1f4b-4c7b-4c90-9f2f-01a3f2a5c111",
  "trainee_id": "user_123",
  "module_code": "fall_hazard_recognition_001",
  "module_version": "1.0.0",
  "device_id": "quest_001",
  "app_version": "1.0.0",
  "language": "en",
  "started_at_client": "2026-05-21T00:00:00Z"
}
```

## 17.2 Start Session Response

```json
{
  "session_id": "session_abc123",
  "server_time": "2026-05-21T00:00:01Z",
  "status": "started"
}
```

## 17.3 Event Batch Request

```json
{
  "events": [
    {
      "event_type": "hazard_found",
      "hazard_code": "open_edge_missing_guardrail",
      "hazard_category": "fall",
      "severity": "critical",
      "action": "marked_hazard",
      "is_correct": true,
      "time_offset_seconds": 125,
      "metadata": {
        "zone": "scaffold_area"
      }
    }
  ]
}
```

## 17.4 Complete Session Request

```json
{
  "client_session_id": "8d1d1f4b-4c7b-4c90-9f2f-01a3f2a5c111",
  "duration_seconds": 480,
  "score": 86,
  "status": "passed",
  "hazards_found": 18,
  "hazards_missed": 3,
  "unsafe_actions": 1,
  "hints_used": 0,
  "ended_at_client": "2026-05-21T00:08:00Z"
}
```

---

## 18. Report Requirements

## 18.1 PDF Session Report

The server must generate a professional PDF report for each completed training session.

Report title:

**Training Completion Report — Supplemental VR Safety Training**

The report should include:

- company name;
- trainee name or employee ID;
- module title;
- module version;
- training date;
- duration;
- final score;
- pass/fail status;
- hazards found;
- hazards missed;
- critical hazards missed;
- unsafe actions;
- recommendation;
- app version;
- device ID;
- report generation timestamp.

The report must not claim official OSHA certification.

## 18.2 CSV Export

CSV export should support manager workflows and external recordkeeping.

Recommended columns:

- company_name;
- trainee_name;
- employee_id;
- module_code;
- module_title;
- module_version;
- started_at;
- ended_at;
- duration_seconds;
- score;
- status;
- hazards_found;
- hazards_missed;
- unsafe_actions;
- hints_used;
- app_version;
- device_id.

---

## 19. Audit Logging Requirements

The backend must record important actions.

Audit log examples:

- successful login;
- failed login;
- user created;
- user updated;
- user deactivated;
- device registered;
- training assignment created;
- training session started;
- training session completed;
- report generated;
- report downloaded;
- CSV export generated;
- permission denied event.

Audit logs should include:

- actor user ID;
- company ID;
- action;
- resource type;
- resource ID;
- timestamp;
- IP address where available;
- user agent where available;
- success/failure;
- metadata.

---

## 20. Security Requirements

Production security requirements:

- HTTPS only in production;
- strong password hashing;
- rate limiting for auth endpoints;
- role-based authorization;
- company-level query scoping;
- secure environment variables;
- no committed secrets;
- safe CORS configuration;
- input validation;
- request size limits;
- structured error responses;
- audit logging;
- database backups;
- least-privilege database user;
- secure Docker configuration;
- controlled access to reports.

---

## 21. Privacy Requirements

The server should collect only data required for training operations and reporting.

Recommended data:

- name or employee ID;
- company/group;
- assigned module;
- completion status;
- score;
- mistakes;
- training timestamps;
- device ID;
- app version.

Avoid by default:

- biometric data;
- raw body movement history;
- unnecessary health data;
- unrelated personal information.

If future versions collect sensitive data, privacy requirements must be reviewed before implementation.

---

## 22. Idempotency and Duplicate Prevention

Offline sync makes duplicate prevention mandatory.

The backend must support idempotent session uploads.

Recommended strategy:

- client generates `client_session_id`;
- server stores unique `(company_id, device_id, client_session_id)`;
- completion endpoint can be safely retried;
- event batch upload should avoid duplicate event insertion where possible;
- server returns existing session if the same upload is repeated.

This prevents duplicate training records when a headset retries after connection loss.

---

## 23. Error Handling

The API must return clear and safe errors.

Example:

```json
{
  "error": {
    "code": "invalid_module",
    "message": "The requested training module is not available for this company."
  }
}
```

Rules:

- do not expose internal stack traces;
- use stable error codes;
- return helpful messages;
- log internal error details server-side;
- keep client-facing errors safe.

---

## 24. Monitoring and Health Checks

Production endpoints:

```http
GET /ping
GET /health
```

`/ping` should confirm the server process is running.

`/health` should check:

- application status;
- database connection;
- migration state if implemented;
- storage access if reports are stored externally.

The server should produce structured logs suitable for production monitoring.

---

## 25. Deployment Requirements

The backend should be deployable with Docker.

Required production files later:

- `Dockerfile`;
- `docker-compose.yml` for local development;
- production environment example file;
- database migrations;
- startup command;
- health check configuration;
- README deployment instructions.

Environment variables may include:

```env
APP_ENV=production
SERVER_PORT=8080
DATABASE_URL=postgres://...
JWT_SECRET=...
CORS_ALLOWED_ORIGINS=https://...
REPORT_STORAGE_PATH=...
LOG_LEVEL=info
```

Secrets must never be committed.

---

## 26. Backup and Data Retention Direction

Production training records may be important to customers. The backend should support reliable backups and retention planning.

Recommended direction:

- regular PostgreSQL backups;
- restore testing;
- report file backup if reports are stored outside DB;
- configurable retention policy;
- clear deletion workflow if customer requests removal;
- audit logs for deletion/export actions.

Exact retention periods should be decided based on customer contract and legal requirements.

---

## 27. Production Development Roadmap

## Release 1 — Production Foundation

Goal: deliver an end-to-end production-ready backend foundation.

Deliverables:

- Go project initialization;
- HTTP router;
- PostgreSQL connection;
- database migrations;
- config management;
- structured logging;
- `/ping` and `/health`;
- companies table;
- users table;
- auth endpoints;
- role middleware;
- company isolation middleware;
- devices table;
- training modules table;
- training sessions table;
- training events table;
- start session endpoint;
- event upload endpoint;
- complete session endpoint;
- idempotency support.

## Release 2 — Dashboard and Reports

Goal: support manager workflows.

Deliverables:

- dashboard API;
- session list;
- session details;
- trainee list;
- module analytics;
- PDF report generation;
- CSV export;
- audit logs for report access.

## Release 3 — Training Operations

Goal: support real customer operations.

Deliverables:

- training assignments;
- group support;
- device heartbeat;
- device status;
- improved offline sync handling;
- company settings;
- admin tools.

## Release 4 — Enterprise Readiness

Goal: support larger customers.

Deliverables:

- SSO/SAML option;
- advanced audit logs;
- LMS integration direction;
- advanced analytics;
- custom module configuration;
- subscription/billing integration if needed.

---

## 28. Production Acceptance Criteria

The backend is production-ready when:

1. It can authenticate users securely.
2. It enforces role permissions.
3. It isolates company data.
4. It registers VR devices.
5. It stores module metadata and versions.
6. It starts training sessions.
7. It receives training events.
8. It completes sessions idempotently.
9. It prevents duplicate offline sync records.
10. It exposes dashboard APIs.
11. It generates PDF reports.
12. It exports CSV records.
13. It writes audit logs.
14. It runs through Docker.
15. It passes health checks.
16. It does not require secrets in the repository.

---

## 29. Final Summary

This repository is the production backend for BuildSafe VR.

The VR client creates the immersive training experience. This server makes the training operational, measurable, secure, and reportable.

The first commercial backend goal is clear:

**Receive training results from the Unreal Engine VR app, store them reliably, protect company data, support dashboard review, and generate professional training reports.**

This server must be built as a real production system from the beginning: secure, reliable, maintainable, auditable, and ready for customer use.