NAME=jat
CURRENT_RELEASE=v0.1.0
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
CC=go
FLAGS=-X github.com/filz0r/jat/internal/version.Version=$(CURRENT_RELEASE) -X github.com/filz0r/jat/internal/version.Commit=$(GIT_COMMIT) -X github.com/filz0r/jat/internal/version.BuildDate=$(BUILD_DATE)

all: $(NAME)

$(NAME):
	@echo "Compiling new server binary"
	@echo "Version: $(CURRENT_RELEASE)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@$(CC) build -ldflags "$(FLAGS)" -o $(NAME)
	@echo "Done!"

release:
	@if git rev-parse "$(CURRENT_RELEASE)" >/dev/null 2>&1; then \
              echo "Error: tag $(CURRENT_RELEASE) already exists"; \
              exit 1; \
	fi
	@echo "Creating release $(CURRENT_RELEASE)"
	@git tag -a "$(CURRENT_RELEASE)" -m "Release $(CURRENT_RELEASE)"

clean_server:
	@echo "Removing old server binary"
	@rm -rf $(NAME)

test:
	@echo $(GIT_COMMIT)
	@echo $(CC) $(FLAGS)

build_server: $(NAME)

dev_server: clean_server build_server
	./jat server

re: clean_server build_server

.PHONY: all re test clean_server build_server