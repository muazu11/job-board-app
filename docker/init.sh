#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- First, terminate any existing connections
    SELECT pg_terminate_backend(pid) 
    FROM pg_stat_activity 
    WHERE datname = '$POSTGRES_DB' AND pid <> pg_backend_pid();
    
    -- Drop the database if it exists
    DROP DATABASE IF EXISTS $POSTGRES_DB;
    
    -- Create fresh database
    CREATE DATABASE $POSTGRES_DB;
EOSQL

# Run schema creation
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/schema.sql

# Run user creation
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/admin_user.sql

# Run fake data population
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/fake_data.sql
