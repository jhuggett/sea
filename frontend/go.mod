module github.com/jhuggett/frontend

go 1.22.4

require (
	design-library v0.0.0-20241002000000-000000000000
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/hajimehoshi/ebiten/v2 v2.8.8
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jhuggett/sea v0.0.0
	golang.org/x/sync v0.8.0 // indirect
)

require golang.org/x/sys v0.30.0 // indirect

require (
	github.com/google/uuid v1.6.0
	github.com/guptarohit/asciigraph v0.7.3
	golang.org/x/image v0.20.0
)

require (
	github.com/beefsack/go-astar v0.0.0-20200827232313-4ecf9e304482 // indirect
	github.com/ebitengine/oto/v3 v3.3.3 // indirect
	github.com/go-text/typesetting v0.2.0 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/ojrac/opensimplex-go v1.0.2 // indirect
	golang.org/x/text v0.18.0 // indirect
	gorm.io/driver/sqlite v1.5.6 // indirect
	gorm.io/gorm v1.25.10 // indirect
)

replace github.com/jhuggett/sea => ../backend

replace design-library => ../design-library
