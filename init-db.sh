#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER world_news WITH PASSWORD 'world_news';
	CREATE DATABASE world_news;
	\c world_news
	GRANT ALL ON SCHEMA public TO world_news;
EOSQL
