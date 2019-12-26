#!/usr/bin/env bash


export AGENT_TRACING_HOST=13.250.21.165:5775


#main.go
export OTTOPOINT_OP="0.0.0.0:8002"


#db/otmdbs/ottomart_postgres.go
export DB_TYPE="POSTGRES"
export DB_POSTGRES_USER="ottopay"
export DB_POSTGRES_PASS="Ottopay!23"
export DB_POSTGRES_NAME="ottopay-api_development"
export DB_POSTGRES_HOST="159.65.139.167"
export DB_POSTGRES_PORT="5432"
export DB_POSTGRES_SSL_MODE="disable"
export DB_POSTGRES_DEBUG="true"

# go run main.go
# go build
# nohup ./ottopoint-purchase > nohup.out 2>&1 &