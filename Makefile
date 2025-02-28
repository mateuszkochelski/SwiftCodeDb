migrateup:
	@docker exec -i postgresDB psql -U root -d swift_codes < db/schema/up/001_db_up.sql > out
	
migratedown:
	@docker exec -i postgresDB psql -U root -d swift_codes < db/schema/down/001_db_down.sql > out
	
migrateupTest:
	@docker exec -i postgresTestDB psql -U test -d testdb < db/schema/up/001_db_up.sql > out

dropdb:
	docker exec -it postgresDB dropdb swift_codes

seedDatabase:
	docker exec -i go-backend sh -c "cd /app/seeder && go build -o main && ./main"

backendRebuild:
	-docker-compose down --volumes
	docker-compose up --build -d
	$(MAKE) migrateupTest
	$(MAKE) migrateup	
	$(MAKE) seedDatabase

backendInit:
	-docker-compose down --volumes
	docker-compose up -d
	$(MAKE) migrateupTest
	$(MAKE) migrateup
	$(MAKE) seedDatabase

backendStart:
	docker-compose start

backendStop:
	docker-compose stop

test:
	docker exec -i go-backend sh -c "cd ../.. && go test -cover ./..."
