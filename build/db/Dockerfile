FROM postgres:14-alpine

ENV POSTGRES_DB=messenger
ENV POSTGRES_USER=messenger
ENV POSTGRES_PASSWORD=messenger

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
