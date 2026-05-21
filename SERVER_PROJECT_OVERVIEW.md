# VRTrainingServer — BuildSafe VR Backend Project Overview

## 1. Project Summary

**VRTrainingServer** is the backend server for the BuildSafe VR construction safety training platform. The server is responsible for receiving training data from the Unreal Engine VR client, storing the data safely, calculating and organizing training results, providing APIs for the dashboard, generating reports, and supporting company-level training management.

The VR client is responsible for the immersive training experience. This server is responsible for everything that must be stored, managed, analyzed, exported, and shown to supervisors or training managers.

The product should be understood as a B2B training system, not only as a VR application. The server is what turns a VR simulation into a measurable training platform. Without the server, the VR app can show a scenario, but the company cannot easily prove who completed training, what mistakes were made, what scores were achieved, and who needs retraining.

---

## 2. Server Vision

The vision of this backend is to become the central system for construction VR safety training data.

The server should answer these questions for a company:

- Who completed training?
- Which module did they complete?
- When did they start and finish?
- How long did the training take?
- What score did they receive?
- Did they pass or fail?
- What hazards did they find?
- What hazards did they miss?
- What unsafe actions did they perform?
- Which workers need retraining?
- Which modules have the most common mistakes?
- Can the manager download a report?

The backend should make the VR training useful for real business operations.

---

## 3. Role of This Repository

This repository contains the server-side system for the VR construction training platform.

The server will handle:

1. Authentication and authorization.
2. Company and tenant management.
3. User and trainee management.
4. Device registration and tracking.
5. Training module metadata.
6. Training session creation.
7. Training event ingestion.
8. Training session completion.
9. Score storage and result history.
10. Report generation.
11. Dashboard APIs.
12. CSV export.
13. PDF report generation.
14. Audit logging.
15. Security and privacy controls.
16. Future integrations with LMS or enterprise systems.

---

## 4. Relationship Between VR Client and Server

The system has two main parts:

```text
Unreal Engine VR Client
        ↓ HTTPS API
VRTrainingServer Backend
        ↓
PostgreSQL Database
        ↓
Web Dashboard / Reports
```

The VR client sends training data to the server. The server stores and organizes that data. A manager or trainer uses the dashboard to view results and export reports.

The VR client should not connect directly to the database. All communication must happen through secure backend API endpoints.

---

## 5. Main Product Data Flow

### 5.1 Training Start

1. Trainee opens VR application.
2. Trainee logs in or enters a training code.
3. VR client starts a module.
4. VR client sends `POST /training/sessions/start` to the server.
5. Server creates a new training session.
6. Server returns a `session_id` to the VR client.

### 5.2 During Training

1. Trainee interacts with objects in VR.
2. VR client records events locally.
3. VR client may send event batches during the session or at the end.
4. Server validates the session and stores events.

### 5.3 Training Completion

1. Trainee finishes the module.
2. VR client calculates initial local result.
3. VR client sends completion data to the server.
4. Server stores final score, duration, pass/fail status, and event summary.
5. Server makes the result available in dashboard.
6. Server can generate PDF/CSV report.

---

## 6. Core Server Responsibilities

## 6.1 Authentication

The server should support secure login for managers, trainers, and admins.

Possible authentication options:

- email and password
- magic code
- training PIN for trainee session start
- QR code for headset login
- JWT access tokens
- refresh tokens or session-based auth

For the MVP, simple email/password for managers and PIN/QR-based training access for trainees may be enough.

## 6.2 Authorization and Roles

The server should use role-based access control.

Recommended roles:

1. **Super Admin**
   - platform-level administration
   - can manage all companies if needed

2. **Company Admin**
   - manages company settings
   - manages users
   - manages devices
   - sees company reports

3. **Safety Manager**
   - assigns modules
   - views training results
   - exports reports

4. **Trainer**
   - starts group training
   - views trainee results for assigned groups

5. **Trainee**
   - completes assigned training
   - may view own result if product requires it

The first MVP can implement fewer roles, but the database should be designed so roles can grow later.

---

## 7. Companies and Multi-Tenant Structure

The server should support multiple companies. Each company must only see its own users, sessions, devices, and reports.

This is very important for B2B SaaS.

Every important table should include `company_id` where appropriate.

Examples:

- users belong to company
- devices belong to company
- trainees belong to company
- sessions belong to company
- reports belong to company

The server must prevent one company from accessing another company’s data.

---

## 8. Training Modules

The server stores metadata about training modules.

The VR client contains the actual interactive scenes, but the server should know which modules exist and what results belong to which module.

Example modules:

- `ppe_site_entry_001`
- `fall_hazard_recognition_001`
- `construction_hazard_hunt_001`
- `electrical_hazard_awareness_001`
- `trench_safety_awareness_001`

Module metadata should include:

- module ID
- code
- title
- category
- description
- version
- pass score
- active/inactive status
- supported languages

---

## 9. Training Sessions

A training session represents one attempt by one trainee to complete one module.

Each session should store:

- session ID
- company ID
- trainee ID
- module ID
- device ID
- started at
- ended at
- duration
- score
- pass/fail status
- language
- app version
- sync status
- created at
- updated at

The session is the main record used for reports.

---

## 10. Training Events

Training events are detailed actions that happened during the VR session.

Examples:

- trainee found a hazard
- trainee missed a hazard
- trainee selected PPE
- trainee selected damaged PPE
- trainee used a hint
- trainee entered unsafe area
- trainee completed checklist item
- trainee clicked wrong object

Each event should store:

- event ID
- session ID
- event type
- hazard code
- hazard category
- severity
- action
- whether action was correct
- time offset in seconds
- metadata JSON
- created at

The event system gives the manager more than just a final score. It explains what actually happened.

---

## 11. Scoring and Result Storage

The first version can let the VR client calculate the initial score, but the server should still validate and store the final result.

In future versions, the server can calculate or verify score using module rules.

Example score model:

- start score: 100
- missed critical hazard: -15
- missed medium hazard: -8
- missed minor hazard: -3
- unsafe action: -20
- hint used: -5
- timeout: -10

Default pass threshold: 80.

The pass threshold should be configurable per module.

---

## 12. Dashboard API

The dashboard will use the backend API to show training data to managers.

Important dashboard pages:

1. **Overview**
   - total trainees
   - completed sessions
   - average score
   - pass rate
   - high-risk trainees
   - recent training activity

2. **Trainees**
   - list of trainees
   - assigned modules
   - completion status
   - latest score
   - training history

3. **Modules**
   - list of modules
   - average score by module
   - common missed hazards
   - completion count

4. **Sessions**
   - session detail
   - score
   - events
   - missed hazards
   - unsafe actions

5. **Reports**
   - PDF report export
   - CSV export
   - date filters
   - group filters

6. **Devices**
   - registered headsets
   - last sync time
   - app version
   - company assignment

---

## 13. MVP API Endpoints

The first backend API can include these endpoints.

### Auth

```http
POST /auth/register
POST /auth/login
POST /auth/logout
GET  /auth/me
```

### Companies

```http
GET  /companies/current
PATCH /companies/current
```

### Users and Trainees

```http
GET  /users
POST /users
GET  /users/{id}
PATCH /users/{id}
DELETE /users/{id}
```

### Devices

```http
POST /devices/register
GET  /devices
GET  /devices/{id}
PATCH /devices/{id}
```

### Training Modules

```http
GET /training/modules
GET /training/modules/{id}
POST /training/modules/assign
```

### Training Sessions

```http
POST /training/sessions/start
POST /training/sessions/{id}/events
POST /training/sessions/{id}/complete
GET  /training/sessions
GET  /training/sessions/{id}
```

### Reports

```http
GET /reports/sessions/{id}/pdf
GET /reports/sessions/export.csv
GET /reports/company/summary
```

### Health

```http
GET /ping
GET /health
```

---

## 14. Recommended Technology Stack

Recommended stack based on project needs:

- Go backend
- PostgreSQL database
- Docker Compose for local development
- HTTPS in production
- JWT or session-based authentication
- REST API first
- structured logging
- PDF generation library
- CSV export
- migration tool for database schema

Possible Go libraries can be chosen later depending on preference. The important point is that the server should be simple, reliable, and easy to deploy.

---

## 15. Database Schema Draft

### companies

```sql
id UUID PRIMARY KEY
name TEXT NOT NULL
industry TEXT
country TEXT
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

### users

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
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

### devices

```sql
id UUID PRIMARY KEY
company_id UUID REFERENCES companies(id)
device_name TEXT
platform TEXT
serial_number TEXT
app_version TEXT
last_seen_at TIMESTAMP
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

### training_modules

```sql
id UUID PRIMARY KEY
code TEXT UNIQUE NOT NULL
title TEXT NOT NULL
category TEXT
version TEXT
pass_score INTEGER NOT NULL DEFAULT 80
is_active BOOLEAN NOT NULL DEFAULT true
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

### training_sessions

```sql
id UUID PRIMARY KEY
company_id UUID REFERENCES companies(id)
trainee_id UUID REFERENCES users(id)
module_id UUID REFERENCES training_modules(id)
device_id UUID REFERENCES devices(id)
started_at TIMESTAMP NOT NULL
ended_at TIMESTAMP
duration_seconds INTEGER
score INTEGER
status TEXT
language TEXT
app_version TEXT
sync_status TEXT
created_at TIMESTAMP NOT NULL
updated_at TIMESTAMP NOT NULL
```

### training_events

```sql
id UUID PRIMARY KEY
session_id UUID REFERENCES training_sessions(id)
event_type TEXT NOT NULL
hazard_code TEXT
hazard_category TEXT
severity TEXT
action TEXT
is_correct BOOLEAN
time_offset_seconds INTEGER
metadata_json JSONB
created_at TIMESTAMP NOT NULL
```

### reports

```sql
id UUID PRIMARY KEY
company_id UUID REFERENCES companies(id)
session_id UUID REFERENCES training_sessions(id)
report_type TEXT
pdf_url TEXT
csv_url TEXT
generated_at TIMESTAMP NOT NULL
created_at TIMESTAMP NOT NULL
```

### audit_logs

```sql
id UUID PRIMARY KEY
company_id UUID
user_id UUID
action TEXT NOT NULL
resource_type TEXT
resource_id TEXT
ip_address TEXT
user_agent TEXT
metadata_json JSONB
created_at TIMESTAMP NOT NULL
```

---

## 16. Example JSON Payloads

### Start Session Request

```json
{
  "trainee_id": "user_123",
  "module_code": "fall_hazard_recognition_001",
  "device_id": "quest_001",
  "app_version": "0.1.0",
  "language": "en"
}
```

### Start Session Response

```json
{
  "session_id": "session_abc123",
  "server_time": "2026-05-21T00:00:00Z"
}
```

### Event Batch Request

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

### Complete Session Request

```json
{
  "duration_seconds": 480,
  "score": 86,
  "status": "passed",
  "hazards_found": 18,
  "hazards_missed": 3,
  "unsafe_actions": 1,
  "hints_used": 0
}
```

---

## 17. PDF Report Requirements

The server should generate a clean and professional PDF report for each training session.

The report should include:

- company name
- trainee name or employee ID
- module title
- training date
- duration
- score
- pass/fail status
- hazards found
- hazards missed
- unsafe actions
- recommendation
- app version
- report generation date

The report title should be:

**Training Completion Report — Supplemental VR Safety Training**

Avoid language such as:

- OSHA certified
- OSHA approved
- official certification

---

## 18. CSV Export Requirements

The server should allow managers to export training results as CSV.

Useful CSV columns:

- company_name
- trainee_name
- employee_id
- module_code
- module_title
- started_at
- ended_at
- duration_seconds
- score
- status
- hazards_found
- hazards_missed
- unsafe_actions
- hints_used
- app_version
- device_id

CSV export is important because companies may want to upload results into their own internal systems.

---

## 19. Audit Logging

The backend should record important system actions.

Examples:

- user login
- failed login
- user created
- user updated
- report downloaded
- training session created
- training session completed
- device registered
- export generated

Audit logs should include:

- who did it
- what action was performed
- when it happened
- what resource was affected
- IP address if available
- success/failure if relevant

This is useful for security, troubleshooting, and future compliance needs.

---

## 20. Security Principles

The backend should follow basic security rules from the beginning.

Required principles:

- HTTPS in production
- password hashing, never plain text passwords
- role-based access control
- company-level data isolation
- input validation
- rate limiting for login endpoints
- audit logs for important actions
- secure environment variables
- no secrets committed to GitHub
- database backups
- safe error messages that do not leak internals

---

## 21. Privacy Principles

The system should collect only the data needed for training and reports.

The MVP should avoid collecting unnecessary sensitive data.

Recommended data to collect:

- name or employee ID
- company/group
- module completed
- score
- mistakes
- date/time
- device ID

Avoid collecting by default:

- biometric data
- continuous body tracking history
- unnecessary health data
- unrelated personal information

If future versions collect more sensitive data, the privacy and legal requirements must be reviewed carefully.

---

## 22. Offline Sync Support

The VR client may complete training while offline. The server should support later sync.

Important behavior:

1. Client creates local session.
2. Client sends it when internet is available.
3. Server accepts session with client timestamps.
4. Server prevents duplicate uploads using idempotency keys or client session IDs.
5. Server responds clearly whether sync succeeded.

Possible fields:

- `client_session_id`
- `device_id`
- `started_at_client`
- `ended_at_client`
- `synced_at_server`
- `idempotency_key`

This prevents duplicate training records if the headset retries upload.

---

## 23. Idempotency

The backend should be prepared for repeated requests from the VR client.

For example, if the Quest loses internet during upload, it may retry the same session completion request. The server should not create duplicate completed sessions.

Use one of these approaches:

- client-generated session UUID
- idempotency key header
- unique constraint on client session ID and device ID

---

## 24. Error Handling

The API should return clear errors.

Example error format:

```json
{
  "error": {
    "code": "invalid_module",
    "message": "The requested training module is not available for this company."
  }
}
```

The server should avoid exposing internal stack traces to clients.

---

## 25. Health Checks

The server should include health endpoints for deployment and monitoring.

Recommended endpoints:

```http
GET /ping
GET /health
```

`/ping` can return simple response.

`/health` can check:

- server running
- database connection
- migration state if needed

---

## 26. Deployment Direction

The server should be easy to deploy with Docker.

Recommended deployment path:

- Dockerfile
- docker-compose.yml for local development
- PostgreSQL container
- environment variables
- migration command
- production build command
- HTTPS behind reverse proxy or WAF

Important environment variables may include:

```env
APP_ENV=development
SERVER_PORT=8080
DATABASE_URL=postgres://...
JWT_SECRET=...
PDF_STORAGE_PATH=...
CORS_ALLOWED_ORIGINS=...
```

Secrets must not be committed to the repository.

---

## 27. MVP Development Plan

### Phase 1 — Backend Foundation

Goal: create the base backend.

Tasks:

- initialize Go project
- configure HTTP router
- connect PostgreSQL
- add migrations
- add `/ping` and `/health`
- create basic config system
- add structured logging

### Phase 2 — Auth and Company Structure

Goal: support managers and company data isolation.

Tasks:

- create companies table
- create users table
- implement password hashing
- implement login
- implement JWT/session auth
- implement roles
- add middleware for company context

### Phase 3 — Training Session API

Goal: receive data from VR client.

Tasks:

- create devices table
- create modules table
- create sessions table
- create events table
- implement start session endpoint
- implement event upload endpoint
- implement complete session endpoint
- handle duplicate sync attempts

### Phase 4 — Dashboard API

Goal: provide data for web dashboard.

Tasks:

- list trainees
- list sessions
- session detail
- module analytics
- overview stats
- filters by date/module/user

### Phase 5 — Reports

Goal: generate proof of training activity.

Tasks:

- PDF session report
- CSV export
- report download endpoint
- audit log for report download

---

## 28. Definition of Done for First Server Demo

The first backend demo is ready when:

1. The server starts locally.
2. `/ping` works.
3. PostgreSQL connection works.
4. A trainee or training code can be created.
5. VR client can start a session.
6. VR client can upload training events.
7. VR client can complete a session.
8. Dashboard API can show the completed session.
9. CSV or PDF report can be generated.
10. Duplicate uploads do not create duplicate records.

---

## 29. Future Features

Future server features may include:

- full web dashboard
- company subscription billing
- module assignment workflows
- group management
- trainer-led class sessions
- LMS integration
- SCORM/xAPI support
- email reports
- scheduled reports
- advanced analytics
- AI-generated retraining recommendations
- custom client site configuration
- multi-region deployment
- SSO/SAML for enterprise customers

---

## 30. Final Summary

This server is the operational brain of the BuildSafe VR platform.

The VR client creates the training experience. The server makes the training measurable, reportable, manageable, and useful for companies.

The first goal is simple:

**Receive training results from the Unreal Engine VR app, store them reliably, show them to managers, and generate clear training reports.**

Once this foundation works, the product can grow into a full B2B VR safety training platform.