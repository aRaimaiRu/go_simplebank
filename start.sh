#!/bin/sh

# exit immediately if a command exits with a non-zero status
set -e

echo "test"
ls
pwd
echo "run db migrations"
ls /app/db/migration
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
# takes all parameters pass to the script and run it
exec "$@"