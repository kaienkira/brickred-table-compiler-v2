BUILD_DEBUG_FLAGS = -gcflags="all=-N -l"
BUILD_RELEASE_FLAGS = -trimpath -ldflags="-s -w"

build:
	@cd src/cmd/brickred-table-compiler && \
		go build -o ../../../bin/brickred-table-compiler
	@cd src/cmd/brickred-table-cutter && \
		go build -o ../../../bin/brickred-table-cutter

build-debug:
	@cd src/cmd/brickred-table-compiler && \
		go build ${BUILD_DEBUG_FLAGS} -o ../../../bin/brickred-table-compiler
	@cd src/cmd/brickred-table-cutter && \
		go build ${BUILD_DEBUG_FLAGS} -o ../../../bin/brickred-table-cutter

fmt:
	@cd src && go fmt ./...

clean:
	@rm -rf bin/
