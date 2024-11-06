docker run \
    --name msa-messanger-db \
    -p 5433:5432 \
    -e POSTGRES_USER=messenger \
    -e POSTGRES_PASSWORD=messenger \
    -e POSTGRES_DB=messenger \
    -e POSTGRES_INITDB_ARGS=--encoding=UTF8 \
    -e POSTGRES_HOST_AUTH_METHOD=password \
    -e PGDATA=/data/ \
    -d postgres:12.4