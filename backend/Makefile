
atlas-schema-diff:
	atlas schema diff \
		--to "file://schema.sql" \
		--from "postgres://dev:dev@localhost:5432/dforget?sslmode=disable" \
		--dev-url "docker://postgres"

atlas-schema-apply:
	atlas schema apply \
		--to "file://schema.sql" \
		--url "postgres://dev:dev@localhost:5432/dforget?sslmode=disable" \
		--dev-url "docker://postgres"

atlas-schema-apply-test:
	atlas schema apply \
		--to "file://schema.sql" \
		--url "postgres://dev:dev@localhost:5432/test?sslmode=disable" \
		--dev-url "docker://postgres"

fmt:
	go fmt ./...

test: fmt
	go test ./...

gen-wire:
	wire gen ./server

gen-db-schema:
	go run cmd/schema/main.go
