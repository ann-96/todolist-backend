export TODO_HTTP_PORT=8080
export TODO_HTTP_HOST=localhost
export USER_HTTP_PORT=9080
export USER_HTTP_HOST=localhost
export SQL_HOST=localhost
export SQL_PORT=5432
export SQL_USER=postgres
export SQL_PASS=postgres
export SQL_DBNAME=postgres
export REDIS_HOST=localhost
export REDIS_PORT=6379

go build -o service-binary
./service-binary
