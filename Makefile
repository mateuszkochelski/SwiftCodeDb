migrateup:
	cat db/schema/up/001_db_up.sql | docker exec -i postgresDB psql -U root -d swift_codes
	
migratedown:
	cat db/schema/down/001_db_down.sql | docker exec -i postgresDB psql -U root -d swift_codes
	
migrateupTest:
	cat db/schema/up/001_db_up.sql | docker exec -i postgresTestDB psql -U test -d testdb

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
