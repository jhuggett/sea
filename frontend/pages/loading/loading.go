package loading

import (
	"design-library/doodad"
	"design-library/label"
	"design-library/position/box"
	"design-library/stack"
)

type LoadingPageController struct {
	Page *LoadingPage
}

type LoadingPage struct {
	doodad.Default

	Controller *LoadingPageController

	LoadingMessage string

	setupFunc func(*LoadingPageController)
}

func New(setupFunc func(*LoadingPageController)) *LoadingPage {
	page := &LoadingPage{
		Controller: &LoadingPageController{},
		setupFunc:  setupFunc,
	}

	page.Controller.Page = page

	return page
}

func (l *LoadingPage) Setup() {

	panel := stack.New(stack.Config{})

	l.AddChild(panel)

	loadingIndicator := label.New(label.Config{
		Message:  "Loading...",
		FontSize: 24,
	})

	panel.AddChild(loadingIndicator)

	loadingMessage := label.New(label.Config{
		Message:  l.LoadingMessage,
		FontSize: 16,
	})

	panel.AddChild(loadingMessage)

	if l.setupFunc != nil {
		go l.setupFunc(l.Controller)
		l.setupFunc = nil // So the page can be re-setup
	}

	l.Children().Setup()

	panel.Box.Computed(func(b *box.Box) {
		b.CenterWithin(l.Box)
	})

	l.Layout().Recalculate()
}
