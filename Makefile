run:
	@docker-compose --env-file .env.dev up

build:
	@docker-compose --env-file .env.dev build

run-prod:
	@docker-compose -f docker-compose.prod.yml up

build-prod:
	@docker-compose -f docker-compose.prod.yml build