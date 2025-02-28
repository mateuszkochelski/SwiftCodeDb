initDb:
	docker run --name postgresDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:17-alpine

startDb:
	docker start postgresDB

initTestDb:
	-docker stop postgres-test && docker rm postgres-test
	docker run --name postgres-test -p 5433:5432 -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=testdb -d postgres:17-alpine
	sleep 2
	docker exec -i postgres-test psql -U test -d testdb < db/schema/up/001_db_up.sql

test:
	-docker stop postgres-test && docker rm postgres-test
	docker run --name postgres-test -p 5433:5432 -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=testdb -d postgres:17-alpine
	sleep 2
	docker exec -i postgres-test psql -U test -d testdb < db/schema/up/001_db_up.sql
	go test -cover ./...
	docker stop postgres-test && docker rm postgres-test
createdb:
	docker exec -it postgresDB createdb --username=root --owner=root swift_codes

migrateup:
	@docker exec -i postgresDB psql -U root -d swift_codes < db/schema/up/001_db_up.sql > out
	
migratedown:
	@docker exec -i postgresDB psql -U root -d swift_codes < db/schema/down/001_db_down.sql > out
	
dropdb:
	docker exec -it postgresDB dropdb swift_codes

seedDatabase:
	docker exec -i go-backend sh -c "cd /app/seeder && go build -o main && ./main"

backendRebuild:
	docker-compose down --volumes
	docker-compose up --build -d
	$(MAKE) migrateup	
	$(MAKE) seedDatabase

backendInit:
	docker-compose down --volumes
	docker-compose up -d
	$(MAKE) migrateup
	$(MAKE) seedDatabase

backendStart:
	docker-compose start

backendStop:
	docker-compose stop