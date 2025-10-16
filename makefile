# Makefile

# Путь к файлу конфигурации
CONFIG_FILE := config.yaml

# Используем yq для чтения YAML
DB_USER     := $(shell yq e '.database.user' $(CONFIG_FILE))
DB_PASSWORD := $(shell yq e '.database.password' $(CONFIG_FILE))
DB_NAME     := $(shell yq e '.database.name' $(CONFIG_FILE))
DB_HOST     := $(shell yq e '.database.host' $(CONFIG_FILE))
DB_PORT     := $(shell yq e '.database.port' $(CONFIG_FILE))
DB_SSLMODE  := $(shell yq e '.database.sslmode' $(CONFIG_FILE))
MIG_PATH    := $(shell yq e '.migrations.path' $(CONFIG_FILE))

# Формируем DSN
DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Применить все миграции (up)
migrate-up:
	migrate -database "$(DSN)" -path "$(MIG_PATH)" up

# Откатить все миграции (down)
migrate-down:
	migrate -database "$(DSN)" -path "$(MIG_PATH)" down

# Посмотреть версию миграций
migrate-version:
	migrate -database "$(DSN)" -path "$(MIG_PATH)" version
