CMDS=cleanmocks db cleandb flyway cleantest service cleanservice tidy migrate-new migrate-up migrate-down
SERVICES=question

deps:
	docker compose -f docker-compose.deps.yml up -d

tool:
	go install github.com/rubenv/sql-migrate/...@latest

lint:
	golangci-lint run --issues-exit-code 0 --out-format code-climate \
		./services/question/... \
		| jq -r '.[] | "\(.location.path):\(.location.lines.begin) \(.description)"'

cleanlint:
	rm -rf gl-code-quality-report.json

mocks:
	cd services/common/domain && mockery --all --inpackage

$(CMDS):
	make -C services $@

$(SERVICES):
	make -C services/$@ service

test:
	gotestsum --junitfile report.xml --format testname -- -covermode=count -coverprofile=coverage.tmp smartkid/services/.../domain/...
	cat coverage.tmp | grep -v ".pb.go" | grep -v "mock_" > coverage.out
	gocover-cobertura < coverage.out > coverage.xml
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out | grep total | awk '{print "total coverage: "$$3" of statements"}'

clean:
	git clean -fX