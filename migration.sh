export POSTGRESQL_URL="postgres://postgres:fib@localhost:5432/postgres?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path ./migrations up
