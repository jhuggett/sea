package sound

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var audioContext *audio.Context
var loop *audio.Player
var single *audio.Player

func Play(path string) (*audio.Player, error) {

	input, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// defer in.Close()

	steam, err := mp3.DecodeWithoutResampling(input)
	if err != nil {
		return nil, err
	}

	// loop := audio.NewInfiniteLoop(steam, steam.Length())

	singlePlayer, err := audioContext.NewPlayer(steam)
	if err != nil {
		return nil, err
	}

	singlePlayer.SetVolume(0.1) // Set volume to 50%
	singlePlayer.Rewind()
	singlePlayer.Play()

	return singlePlayer, nil
}

func Setup() {
	audioContext = audio.NewContext(48_000)

}

type BackgroundMusic struct {
	songs   []string
	current string

	player *audio.Player
}

func NewBackgroundMusic(songs []string) *BackgroundMusic {
	return &BackgroundMusic{
		songs:   songs,
		current: "",
	}
}

func (bm *BackgroundMusic) AddSong(song string) {
	bm.songs = append(bm.songs, song)
}

func (bm *BackgroundMusic) Run(ctx context.Context) {
	for {
		if bm.player != nil {
			if bm.player.IsPlaying() {
				<-time.After(1 * time.Second)
				continue
			}
		}

		if bm.current != "" {
			bm.songs = append(bm.songs, bm.current)
		}

		songIndex := rand.Intn(len(bm.songs))
		bm.current = bm.songs[songIndex]
		bm.songs = append(bm.songs[:songIndex], bm.songs[songIndex+1:]...)

		player, err := Play(bm.current)
		if err != nil {
			slog.Error("Failed to play song", "song", bm.current, "error", err)
			panic(err)
		}

		bm.player = player
	}
}
