migrate-up:
	migrate -path ./migrations/migrations -database "$(DB_DSN)" up