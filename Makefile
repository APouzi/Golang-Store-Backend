up:
	docker compose up

down:
	docker compose down

restart:
	docker kill golang-store-backend-app-1; docker-compose -f ./docker-compose.yml build; docker-compose -f ./docker-compose.yml up -d

clean:
	docker compose down -v; docker rm -vf $$(docker ps -aq); docker rmi -f $$(docker images -aq); docker image prune -f; docker volume prune -f; docker system prune -f

clean-start:
	make clean; make bootstrap