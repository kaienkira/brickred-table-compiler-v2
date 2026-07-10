BUILD_DEBUG_FLAGS = -gcflags="all=-N -l"
BUILD_RELEASE_FLAGS = -buildvcs=false -trimpath -ldflags="-s -w"

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

build-linux:
	@cd src/cmd/brickred-table-compiler && \
		GOOS=linux arch=amd64 go build $(BUILD_RELEASE_FLAGS) \
			-o ../../../bin/brickred-table-compiler
	@cd src/cmd/brickred-table-cutter && \
		GOOS=linux arch=amd64 go build $(BUILD_RELEASE_FLAGS) \
			-o ../../../bin/brickred-table-cutter

build-windows:
	@cd src/cmd/brickred-table-compiler && \
		GOOS=windows arch=amd64 go build $(BUILD_RELEASE_FLAGS) \
			-o ../../../bin/brickred-table-compiler.exe
	@cd src/cmd/brickred-table-cutter && \
		GOOS=windows arch=amd64 go build $(BUILD_RELEASE_FLAGS) \
			-o ../../../bin/brickred-table-cutter.exe

fmt:
	@cd src && go fmt ./...

clean:
	@rm -rf bin/
