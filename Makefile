# Go compiler
GO = go

# Directory containing Go source files
FILE_DIR = cmd

# Directory containing compiled executables
OUT_DIR = bin

# File suffix for Go source files
FILE_SUFFIX = go

all: .clean .build .generate

.clean:
	rm -rf bin/

.build:
	for filename in $$(find ${FILE_DIR} -name "*${FILE_SUFFIX}"); do\
		path=$${filename%/*.*};\
		path=$${path##${FILE_DIR}/};\
		echo Compiling $${filename} to $${path} ...;\
		${GO} build -o ${OUT_DIR}/$${path} $${filename};\
	done

.generate:
	for filename in $$(find ${OUT_DIR} -type f -executable); do\
		$${filename} --generate-description;\
	done

.PHONY: build
build: .build .generate

.PHONY: clean
clean: .clean

.PHONY: generate
generate: .generate

.PHONY: install
install:
	@echo "There's no target 'install'. Use 'build', 'clean' or 'generate' instead."
