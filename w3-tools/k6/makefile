
run-debug:
	# git pull origin main;
	DEBUG_ALL=true k6 run script.js > output.txt 2>&1;

run:
	# git pull origin main;
	k6 run script.js;

run-load-test:
	# git pull origin main;
	LOAD_TEST=true k6 run script.js;

compile-pb:
	protoc --go_out=./srv/entity/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./srv/entity/pb \
		--go-grpc_opt=paths=source_relative *.proto
tidy-go:
	@cd ./srv && go mod tidy

compile-go:
	@cd ./srv && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ./cmd/main ./main.go

run-go-binary:
	./cmd/main

run-go:
	@cd ./srv && go run main.go