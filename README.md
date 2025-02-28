# SwiftCodeDb

![Go](https://img.shields.io/badge/Go-1.23-blue) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-green) ![Docker](https://img.shields.io/badge/Docker-✔️-blue)

Swift Code API is a RESTful service for querying bank details using SWIFT codes.


## Prerequisites
- Docker
- Makefile

```sh
git clone https://github.com/mateuszkochelski/SwiftCodeDb
```
# How to run backend
```sh
make backendInit
```
During the startup of the Docker container, the database is seeded with predefined data parsed using  seeder/seeder.go
# How to run tests
```sh
make test
```

# How to get into docker container to run specific test
```sh
docker exec -i go-backend sh
```
The application is running on localhost:8080

# Instruction without makefile

To be sure that containers are destroyed.
```sh
docker-compose down --volumes
``` 


Build containers
```sh
docker-compose up --build -d
```


Migrate database for testing
```sh
cat db/schema/up/001_db_up.sql | docker exec -i postgresTestDB psql -U test -d testdb 
```


Migrate database
```sh
cat db/schema/up/001_db_up.sql | docker exec -i postgresDB psql -U root -d swift_codes
```


Insert data to database
```sh
docker exec -i go-backend sh -c "cd /app/seeder && go build -o main && ./main"
```

The application is running on localhost:8080




