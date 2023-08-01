.PHONY: db cleandb flyway mocks cleanmocks test

all: service

init:
	go get

db:
	if [ "${DATABASE_NAME}" != "" ]; then \
        PGPASSWORD=$${DATABASE_PASSWORD} createdb --host=${DATABASE_HOST} --port=${DATABASE_PORT} --username=${DATABASE_USER} ${DATABASE_NAME} || true; \
    fi;

cleandb:
	if [ "${DATABASE_NAME}" != "" ]; then \
        PGPASSWORD=$${DATABASE_PASSWORD} dropdb --host=${DATABASE_HOST} --port=${DATABASE_PORT} --username=${DATABASE_USER} ${DATABASE_NAME} || true; \
    fi;

migrate-new:
	if [ "${DATABASE_NAME}" != "" ]; then \
		$(eval NAME := $(shell read -p "Enter new file name: " v && echo $$v)) \
		$(eval CMD:= $*) \
		cd db && sql-migrate new ${NAME}; \
	fi;
	
MIG_ENV ?=local
migrate-up:
	if [ "${DATABASE_NAME}" != "" ]; then \
		cd db && sql-migrate up -limit=$${limit:-0} -env=${MIG_ENV} -config=dbconfig.yml;\
	fi;

migrate-down:
	if [ "${DATABASE_NAME}" != "" ]; then \
		cd db && sql-migrate down -limit=$${limit:-1} -env=${MIG_ENV} -config=dbconfig.yml;\
	fi;

mocks:
	mockery --all --inpackage

cleanmocks:
	find . -name "mock_*.go" -exec rm -rf {} \;

test:
	go test ./... -covermode=count -coverprofile=coverage.tmp
	cat coverage.tmp | grep -v ".pb.go" | grep -v "mock_" > coverage.out
	gocover-cobertura < coverage.out > coverage.xml
	go tool cover -html=coverage.out -o coverage.html

cleantest:
	rm -rf coverage.*

service:
	if [ "${SERVICE_NAME}" != "" ]; then go build -o ${SERVICE_NAME}; fi;

cleanservice:
	rm ${SERVICE_NAME} || true

tidy:
	go mod tidy
