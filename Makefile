.PHONY: run-frontend docker-up docker-down

# Frontend operations
run-frontend:
	@echo "Starting frontend development server..."
	@cd frontend && npm run dev

# Docker operations
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

