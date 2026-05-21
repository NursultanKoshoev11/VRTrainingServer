# VRTrainingServer Docker Deployment

This document explains how to run the VRTrainingServer backend with Docker.

## 1. Requirements

Install:

- Docker
- Docker Compose v2

## 2. Environment Setup

Copy the example environment file:

```bash
cp .env.example .env
```

Then edit `.env` and change all placeholder values before production use.

Required production changes:

- `POSTGRES_PASSWORD`
- `JWT_SECRET`
- `DATABASE_URL`
- `CORS_ALLOWED_ORIGINS`

Never commit `.env` to GitHub.

## 3. Start the Stack

```bash
docker compose up -d --build
```

This starts:

- PostgreSQL 16
- VRTrainingServer API

## 4. Check Containers

```bash
docker compose ps
```

## 5. Check API Health

```bash
curl http://localhost:8080/ping
curl http://localhost:8080/health
```

Expected `/ping` response:

```json
{"status":"ok"}
```

## 6. View Logs

```bash
docker compose logs -f api
```

## 7. Stop the Stack

```bash
docker compose down
```

To remove volumes also:

```bash
docker compose down -v
```

## 8. Important Production Notes

The current Docker stack is enough to run the API and database. Before a real customer deployment, add:

- HTTPS reverse proxy or WAF in front of the API;
- managed database backups;
- migration runner;
- secret manager;
- monitoring and alerts;
- production domain and CORS origin;
- log retention policy.

The API image runs as a non-root user in the final container image.
