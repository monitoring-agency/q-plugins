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
	rm -rf ${OUT_DIR}/

.build:
	find $$(find ${FILE_DIR} -name "*${FILE_SUFFIX}") | while read FILE ; do IFS="/" read -ra P <<< "$$FILE"; [[ $${P[-1]} = $${P[-2]}* ]] && \
		echo $$(source=$$(echo $${P[@]:0:$$(echo $${#P[@]} - 1)} | tr " " "/"); target=$$(echo ${OUT_DIR}/$$(echo $${P[@]:1:$$(echo $${#P[@]} - 2)} | tr " " "/")); echo Compiling $${source} to $${target} ...; CGO_ENABLED=0 ${GO} build -ldflags=-w -o $${target} ./$${source}; echo "done");\
	done; exit 0

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
