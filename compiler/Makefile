BUILD_RELEASE_FLAGS = -trimpath -ldflags="-s -w"

build:
	@cd src/cmd/brickred-table-compiler && \
		go build -o ../../../bin/brickred-table-compiler

fmt:
	@cd src && go fmt ./...

clean:
	@rm -rf bin/
