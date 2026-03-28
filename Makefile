# ============================================================
#  Daily Hello — Root Makefile
#  Backend port : 8282
#  DB port      : 5461 (PostgreSQL)
# ============================================================

SERVICE_DIR  := ./daily-hello-service
APP_DIR      := ./daily-hello-app

.DEFAULT_GOAL := help

.PHONY: help \
        service.start service.stop service.restart service.logs \
        app.start app.run app.analyze app.build \
        db.start db.stop db.migrate \
        dev all stop

# ─────────────────────────────────────────────
#  Help
# ─────────────────────────────────────────────
help:
	@echo ""
	@echo "  Daily Hello — Available commands"
	@echo ""
	@echo "  ── Backend ──────────────────────────────────────"
	@echo "  make service.start     Start backend (go run)"
	@echo "  make service.stop      Kill backend process on port 8282"
	@echo "  make service.restart   Restart backend"
	@echo "  make service.logs      Tail backend log file"
	@echo ""
	@echo "  ── Flutter App ──────────────────────────────────"
	@echo "  make app.start         Run Flutter app (default device)"
	@echo "  make app.run           Run with verbose logs"
	@echo "  make app.analyze       Analyze Dart code"
	@echo "  make app.build         Build APK release"
	@echo ""
	@echo "  ── Database ─────────────────────────────────────"
	@echo "  make db.start          Start PostgreSQL via docker"
	@echo "  make db.stop           Stop PostgreSQL container"
	@echo "  make db.migrate        Run auto migration (via backend)"
	@echo ""
	@echo "  ── Combined ─────────────────────────────────────"
	@echo "  make dev               Start DB + Backend together"
	@echo "  make all               Start DB + Backend + Flutter app"
	@echo "  make stop              Stop all running services"
	@echo ""

# ─────────────────────────────────────────────
#  Backend
# ─────────────────────────────────────────────
service.start:
	@echo ">> Starting backend service..."
	@cd $(SERVICE_DIR) && go run cmd/server/main.go

service.stop:
	@echo ">> Stopping backend service on port 8282..."
	@lsof -ti :8282 | xargs kill -9 2>/dev/null && echo "   Stopped." || echo "   Not running."

service.restart: service.stop
	@sleep 1
	@$(MAKE) service.start

service.logs:
	@tail -f $(SERVICE_DIR)/app.log 2>/dev/null || echo "No log file found at $(SERVICE_DIR)/app.log"

# ─────────────────────────────────────────────
#  Flutter App
# ─────────────────────────────────────────────
app.start:
	@echo ">> Running Flutter app..."
	@cd $(APP_DIR) && flutter run

app.run:
	@echo ">> Running Flutter app (verbose)..."
	@cd $(APP_DIR) && flutter run -v

app.analyze:
	@echo ">> Analyzing Dart code..."
	@cd $(APP_DIR) && flutter analyze lib/

app.build:
	@echo ">> Building APK release..."
	@cd $(APP_DIR) && flutter build apk --release
	@echo "   Output: $(APP_DIR)/build/app/outputs/flutter-apk/app-release.apk"

# ─────────────────────────────────────────────
#  Database (Docker)
# ─────────────────────────────────────────────
DB_CONTAINER := daily-hello-postgres
DB_PORT      := 5461
DB_USER      := ams
DB_PASS      := ams123
DB_NAME      := dhdb

db.start:
	@echo ">> Starting PostgreSQL container..."
	@docker run -d \
		--name $(DB_CONTAINER) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		--restart unless-stopped \
		postgres:16-alpine 2>/dev/null \
	&& echo "   Container started." \
	|| docker start $(DB_CONTAINER) 2>/dev/null \
	&& echo "   Container resumed." \
	|| echo "   Already running."

db.stop:
	@echo ">> Stopping PostgreSQL container..."
	@docker stop $(DB_CONTAINER) 2>/dev/null && echo "   Stopped." || echo "   Not running."

db.migrate:
	@echo ">> Running auto migration via backend (one-shot)..."
	@cd $(SERVICE_DIR) && go run cmd/server/main.go & \
		sleep 3 && kill %% 2>/dev/null; echo "   Migration done."

# ─────────────────────────────────────────────
#  Combined
# ─────────────────────────────────────────────
dev: db.start
	@echo ">> Starting backend (background)..."
	@cd $(SERVICE_DIR) && go run cmd/server/main.go & echo $$! > /tmp/dhs.pid
	@echo "   Backend PID saved to /tmp/dhs.pid"
	@echo "   Backend running at http://localhost:8282"

all: dev
	@echo ">> Starting Flutter app..."
	@$(MAKE) app.start

stop:
	@echo ">> Stopping all services..."
	@[ -f /tmp/dhs.pid ] && kill $$(cat /tmp/dhs.pid) 2>/dev/null && rm /tmp/dhs.pid && echo "   Backend stopped." || true
	@$(MAKE) service.stop
	@$(MAKE) db.stop
	@echo "   Done."
