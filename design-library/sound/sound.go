package sound

import (
	"os"

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

	loop := audio.NewInfiniteLoop(steam, steam.Length())

	singlePlayer, err := audioContext.NewPlayer(loop)
	if err != nil {
		return nil, err
	}

	singlePlayer.SetVolume(0.5) // Set volume to 50%
	singlePlayer.Rewind()
	singlePlayer.Play()

	return singlePlayer, nil
}

func Setup() {
	audioContext = audio.NewContext(48_000)

}
