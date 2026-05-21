$ErrorActionPreference = "Stop"

# Windows Docker deployment check for VRTrainingServer.
$env:APP_ENV = if ($env:APP_ENV) { $env:APP_ENV } else { "production" }
$env:SERVER_PORT = if ($env:SERVER_PORT) { $env:SERVER_PORT } else { "8080" }
$env:POSTGRES_DB = if ($env:POSTGRES_DB) { $env:POSTGRES_DB } else { "vrtraining" }
$env:POSTGRES_USER = if ($env:POSTGRES_USER) { $env:POSTGRES_USER } else { "vrtraining" }
$env:POSTGRES_PASSWORD = if ($env:POSTGRES_PASSWORD) { $env:POSTGRES_PASSWORD } else { "vrtraining_local_password" }
$env:DATABASE_URL = if ($env:DATABASE_URL) { $env:DATABASE_URL } else { "postgres://$($env:POSTGRES_USER):$($env:POSTGRES_PASSWORD)@postgres:5432/$($env:POSTGRES_DB)?sslmode=disable" }
$env:JWT_SECRET = if ($env:JWT_SECRET) { $env:JWT_SECRET } else { "local-dev-secret-change-before-production" }
$env:CORS_ALLOWED_ORIGINS = if ($env:CORS_ALLOWED_ORIGINS) { $env:CORS_ALLOWED_ORIGINS } else { "http://localhost:3000" }
$env:REPORT_STORAGE_PATH = if ($env:REPORT_STORAGE_PATH) { $env:REPORT_STORAGE_PATH } else { "/app/reports" }
$env:MIGRATIONS_PATH = if ($env:MIGRATIONS_PATH) { $env:MIGRATIONS_PATH } else { "/app/migrations" }
$env:LOG_LEVEL = if ($env:LOG_LEVEL) { $env:LOG_LEVEL } else { "info" }
$env:COMPOSE_PROJECT_NAME = if ($env:COMPOSE_PROJECT_NAME) { $env:COMPOSE_PROJECT_NAME } else { "vrtraining" }

$healthUrl = "http://localhost:$($env:SERVER_PORT)/health"

function Dump-DockerLogs {
    Write-Host "Deployment failed. Dumping container status and recent logs..."
    docker compose ps
    docker compose logs --no-color --tail=240 api
    docker compose logs --no-color --tail=160 postgres
}

try {
    Write-Host "Runner name: $env:RUNNER_NAME"
    Write-Host "Runner OS: $env:RUNNER_OS"
    Write-Host "Workspace: $env:GITHUB_WORKSPACE"

    docker --version
    docker compose version

    if ($env:JWT_SECRET -eq "local-dev-secret-change-before-production") {
        Write-Host "Warning: using fallback JWT_SECRET. Set VRTRAINING_JWT_SECRET for production."
    }

    Write-Host "Stopping old VRTrainingServer containers if they exist..."
    docker compose down --remove-orphans

    Write-Host "Building and starting VRTrainingServer containers..."
    docker compose up -d --build --remove-orphans

    Write-Host "Waiting for API health endpoint: $healthUrl"
    $healthy = $false
    for ($attempt = 1; $attempt -le 90; $attempt++) {
        try {
            $response = Invoke-WebRequest -Uri $healthUrl -UseBasicParsing -TimeoutSec 5
            if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 300) {
                Write-Host "API health check passed."
                Write-Host $response.Content
                $healthy = $true
                break
            }
        } catch {
            Start-Sleep -Seconds 2
        }
    }

    if (-not $healthy) {
        throw "API health check failed after 90 attempts."
    }

    Write-Host "Container status:"
    docker compose ps

    Write-Host "Recent API logs:"
    docker compose logs --no-color --tail=200 api

    Write-Host "Recent Postgres logs:"
    docker compose logs --no-color --tail=120 postgres

    Write-Host "VRTrainingServer Docker deployment completed successfully."
} catch {
    Write-Host "ERROR: $($_.Exception.Message)"
    Dump-DockerLogs
    exit 1
}
