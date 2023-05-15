PORT = 80

restart:
	docker kill golang-store-backend-app-1; docker-compose -f ./docker-compose.yml build; docker-compose -f ./docker-compose.yml up -d