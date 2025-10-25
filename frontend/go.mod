module github.com/jhuggett/frontend

go 1.24.0

toolchain go1.24.6

require (
	design-library v0.0.0-20241002000000-000000000000
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/hajimehoshi/ebiten/v2 v2.9.0
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jhuggett/sea v0.0.0
	golang.org/x/sync v0.17.0 // indirect
)

require golang.org/x/sys v0.36.0 // indirect

require (
	github.com/google/uuid v1.6.0
	github.com/guptarohit/asciigraph v0.7.3
	golang.org/x/image v0.31.0
)

require (
	github.com/beefsack/go-astar v0.0.0-20200827232313-4ecf9e304482 // indirect
	github.com/ebitengine/oto/v3 v3.4.0 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20231223183121-56fa3ac82ce7 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	github.com/hajimehoshi/ebiten v1.12.13 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/ojrac/opensimplex-go v1.0.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/exp/shiny v0.0.0-20251002181428-27f1f14c8bb9 // indirect
	golang.org/x/mobile v0.0.0-20250813145510-f12310a0cfd9 // indirect
	golang.org/x/text v0.29.0 // indirect
	gorm.io/driver/sqlite v1.5.6 // indirect
	gorm.io/gorm v1.25.10 // indirect
)

replace github.com/jhuggett/sea => ../backend

replace design-library => ../design-library
