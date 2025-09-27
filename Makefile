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
	yarn; \
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

# --------------------------------

.PHONY: design-library
design-library:
	go run ./design-library/main;

.PHONY: game
game:
	cd frontend; \
	go run main.go

.PHONY: bundle-mac
bundle-mac:
	cd frontend; \
	go build -o "Ships Colonies Commerce"; \
	mkdir -p ../build; \
	mkdir -p ../build/"Ships Colonies Commerce.app"/Contents/MacOS; \
	mkdir -p ../build/"Ships Colonies Commerce.app"/Contents/Resources; \
	cp "./Ships Colonies Commerce" ../build/"Ships Colonies Commerce.app"/Contents/MacOS/"Ships Colonies Commerce"; \
	cp ../assets/artwork/another-compass-rose.png ../build/"Ships Colonies Commerce.app"/Contents/Resources/icon.png; \
	echo '<?xml version="1.0" encoding="UTF-8"?>\n\
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">\n\
<plist version="1.0">\n\
<dict>\n\
	<key>CFBundleExecutable</key>\n\
	<string>Ships Colonies Commerce</string>\n\
	<key>CFBundleIconFile</key>\n\
	<string>icon</string>\n\
	<key>CFBundleIdentifier</key>\n\
	<string>com.jhuggett.sea</string>\n\
	<key>CFBundleName</key>\n\
	<string>Ships Colonies Commerce</string>\n\
	<key>CFBundlePackageType</key>\n\
	<string>APPL</string>\n\
	<key>CFBundleVersion</key>\n\
	<string>1.0</string>\n\
	<key>NSHighResolutionCapable</key>\n\
	<true/>\n\
	<key>LSMinimumSystemVersion</key>\n\
	<string>10.13</string>\n\
</dict>\n\
</plist>' > ../build/"Ships Colonies Commerce.app"/Contents/Info.plist; \
	chmod +x ../build/"Ships Colonies Commerce.app"/Contents/MacOS/"Ships Colonies Commerce"; \
	echo "Created Ships Colonies Commerce.app in the build directory"; \
	codesign --force --deep --sign - ../build/"Ships Colonies Commerce.app"; \
	rm "./Ships Colonies Commerce"