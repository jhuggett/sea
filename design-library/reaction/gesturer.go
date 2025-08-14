package reaction

import (
	"log/slog"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Press struct {
	StartX, StartY int
	X, Y           int
	TimeStart      time.Time
	Button         ebiten.MouseButton
}

type gesturer struct {
	events map[ReactionType][]Reaction
	MouseX int
	MouseY int
	Press  *Press
}

func (g *gesturer) Register(reaction Reaction) func() {
	if g.events == nil {
		g.events = make(map[ReactionType][]Reaction)
	}
	if g.events[reaction.ReactionType()] == nil {
		g.events[reaction.ReactionType()] = []Reaction{}
	}
	g.events[reaction.ReactionType()] = append(g.events[reaction.ReactionType()], reaction)

	return func() {
		for i, r := range g.events[reaction.ReactionType()] {
			if r == reaction {
				g.events[reaction.ReactionType()] = append(g.events[reaction.ReactionType()][:i], g.events[reaction.ReactionType()][i+1:]...)
				break
			}
		}
	}
}

type Event struct {
	stopPropagation bool
}

func (e *Event) StopPropagation() {
	e.stopPropagation = true
}

func (g *gesturer) trigger(reactionType ReactionType, data any) {
	event := &Event{}

	if reactions, ok := g.events[reactionType]; ok {
		for _, reaction := range reactions {
			if event.stopPropagation {
				return
			}
			err := reaction.TryPerform(event, data)
			if err != nil {
				slog.Error("Error performing reaction", "reaction", reaction, "error", err)
			}
		}
	}
}

type ReactionType string

type Gesturer interface {
	Update()
	Register(reaction Reaction) func()
}

func NewGesturer() *gesturer {
	return &gesturer{}
}

// type MouseMoved struct {
// 	X, Y int
// }

// type MouseUpEvent struct {
// 	X, Y   int
// 	Button ebiten.MouseButton
// }

func (g *gesturer) Update() {
	x, y := ebiten.CursorPosition()

	// Keydown events
	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		g.trigger(KeyDown, KeyDownEvent{
			Key: key,
		})
	}

	if x != g.MouseX || y != g.MouseY {
		g.trigger(MouseMoved, MouseMovedEvent{X: x, Y: y})
	}

	g.MouseX = x
	g.MouseY = y

	_, yoff := ebiten.Wheel()
	if yoff != 0 {
		g.trigger(MouseWheel, MouseWheelEvent{YOffset: yoff})
	}

	var pressedMouseButton ebiten.MouseButton = -1
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		pressedMouseButton = ebiten.MouseButtonLeft
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		pressedMouseButton = ebiten.MouseButtonRight
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		pressedMouseButton = ebiten.MouseButtonMiddle
	}

	if pressedMouseButton != -1 {
		if g.Press == nil {
			g.Press = &Press{
				StartX:    x,
				StartY:    y,
				X:         x,
				Y:         y,
				TimeStart: time.Now(),
				Button:    pressedMouseButton,
			}
		}

		if time.Since(g.Press.TimeStart) > 100*time.Millisecond || (math.Abs(float64(g.Press.StartX-x)) > 25 || math.Abs(float64(g.Press.StartY-y)) > 25) {
			if g.Press.X != x || g.Press.Y != y {
				g.trigger(MouseDrag, OnMouseDragEvent{
					StartX:    g.Press.StartX,
					StartY:    g.Press.StartY,
					X:         x,
					Y:         y,
					TimeStart: g.Press.TimeStart,
					Button:    g.Press.Button,
				})
			}
		}

		g.Press.X = x
		g.Press.Y = y

	} else {
		if g.Press != nil {
			if time.Since(g.Press.TimeStart) < 100*time.Millisecond || (math.Abs(float64(g.Press.StartX-g.Press.X)) < 8 && math.Abs(float64(g.Press.StartY-g.Press.Y)) < 8) {
				g.trigger(MouseUp, MouseUpEvent{
					X:      g.Press.X,
					Y:      g.Press.Y,
					Button: g.Press.Button,
				})
			}
			g.Press = nil
		}
	}
}

func (g *gesturer) Teardown() {
}
