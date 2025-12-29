# How to run the test with the database connection:

## Mac or Linux
REDIS_ADDR=localhost:6379 POSTGRES_URL=postgres://user:password@localhost:5432/db?sslmode=disable go test ./tests/ -v

## Windows PowerShell
$env:REDIS_ADDR="localhost:6379"; $env:POSTGRES_URL="postgres://user:password@localhost:5432/db?sslmode=disable"; go test ./tests/ -v
