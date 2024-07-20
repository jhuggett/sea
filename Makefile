.PHONY: dev
dev:
	cd ./backend/main; \
	go build -o ../../dist/dev/server; \
	cd ../../frontend; \
	bun build ./index.ts --compile --outfile ../dist/dev/game; \
	cd .. && ./dist/dev/game; \

.PHONY: generate-types
generate-types:
	cd ./backend; \
	tygo generate; \

.PHONY: backend
backend:
	make generate-types; \
	cd ./frontend/web; \
	yarn upgrade @shared/sea;