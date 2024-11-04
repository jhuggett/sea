.PHONY: run-backend
run-backend: 
	cd ./backend/main; \
	go run main.go; \

.PHONY: backend-dev
backend-dev:
	$(MAKE) backend; \
	$(MAKE) run-backend;

.PHONY: frontend-setup
frontend-setup:
	cd ./frontend/web-react; \
	yarn install; \

.PHONY: frontend-dev
frontend-dev:
	cd ./frontend/web-react; \
	yarn dev; \

.PHONY: generate-types
generate-types:
	cd ./backend; \
	tygo generate; \

.PHONY: backend
backend:
	make generate-types; \
	cd ./frontend/web-react; \
	yarn upgrade @shared/sea;