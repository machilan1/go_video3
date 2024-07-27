RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m

default: 
	@echo " ${RED}Please specify the target to run"


setup: 
ifeq ($(OS),wsl)
	@echo "Setup in wsl"
	make tailwindcss-install-wsl
else ifeq ($(OS),mac)
	@echo "Setup in mac"
	make tailwindcss-install-mac
else
	@echo "Please specify the OS with OS flag (ex: make setup OS=mac)"
endif


run:
	air

tidy:
	go mod tidy

resetDB:
	docker-compose down -v
	docker-compose up 

psql:
	psql -h localhost -p 5432 -U postgres -d maya

migrate:
	migrate create -seq -ext=.sql -dir=sql/migrations $(name)

migrate-up:
	migrate -path sql/migrations -database="postgresql://postgres:123456@localhost:5432/maya?sslmode=disable" -verbose up



tailwindcss-install-wsl:
	@echo "${YELLOW}Installing taiwind for wsl environment...${NC}"
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
	chmod +x tailwindcss-linux-x64
	mv tailwindcss-linux-x64 bin/tailwindcss 
	@echo "${GREEN}Done!${NC}"

tailwindcss-install-mac:
	@echo "${YELLOW}Installing taiwind for mac environment...${NC}"
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
	chmod +x tailwindcss-macos-arm64
	mv tailwindcss-macos-arm64 bin/tailwindcss
	@echo "${GREEN}Done!${NC}"

tailwindcss-watch:
	./bin/tailwindcss -i statics/ui/main.css -o statics/ui/output.css --watch

tailwindcss-compile:
	./bin/tailwindcss -i statics/ui/main.css -o statics/ui/output.css --minify

.PHONY: default