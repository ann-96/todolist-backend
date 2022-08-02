export HTTP_HOST=localhost
export HTTP_PORT=8080
export SQL_HOST=localhost
export SQL_PORT=5432
export SQL_USER=postgres
export SQL_PASS=postgres
export SQL_DBNAME=postgres

go build -o service-binary
./service-binary
