run:
	git pull origin main;
	go run caddy/main.go;

test:
	git pull origin main;
	k6 run --vus 1 --iterations 1 index.js;

load-test:
	git pull origin main;
	LOAD_TEST=true k6 run index.js;

load-test-debug:
	git pull origin main;
	LOAD_TEST=true DEBUG_ALL=true k6 run --vus 300 --iterations 300 index.js;

test-debug:
	git pull origin main;
	DEBUG_ALL=true k6 run --vus 1 --iterations 1 index.js;

test-no-pull:
	DEBUG_ALL=true k6 run --vus 1 --iterations 1 index.js;

tidy:
	go mod tidy

compile-proto:
	protoc --go_out=./caddy/entity/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./caddy/entity/pb \
		--go-grpc_opt=paths=source_relative *.proto