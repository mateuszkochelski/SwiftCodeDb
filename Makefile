postgres:
	docker run --name postgresDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:17-alpine

createdb:
	docker exec -it postgresDB createdb --username=root --owner=root swift_codes

migrateup:
	docker exec -i postgresDB psql -U root -d swift_codes < db/schema/up/001_banks_up.sql
	
migratedown:
	docker exec -i postgresDB psql -U root -d swift_codes < db/schema/down/001_banks_down.sql
	
dropdb:
	docker exec -it postgresDB dropdb swift_codes

.PHONY: initDocker postgres