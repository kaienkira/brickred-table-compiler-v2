BUILD_RELEASE_FLAGS = -trimpath -ldflags="-s -w"

build:
	@cd src/cmd/brickred-table-compiler && \
		go build -o ../../../bin/brickred-table-compiler
	@cd src/cmd/brickred-table-cutter && \
		go build -o ../../../bin/brickred-table-cutter

fmt:
	@cd src && go fmt ./...

clean:
	@rm -rf bin/
