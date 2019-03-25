package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"os"
)

const (
	LANG_VIEW = "languages"
	MAIN_VIEW = "main"
)

var (
	viewArr = []string{LANG_VIEW, MAIN_VIEW}
	active  = 0
)

func relativeSize(g *gocui.Gui) (int, int) {
	tw, th := g.Size()
	return (tw * 3) / 10, (th * 70) / 100
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		//	log.Panicln(err)
		log.Fatal("Failed to initialize GUI", err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	file, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(MAIN_VIEW, gocui.KeyArrowDown, gocui.ModNone, goDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(LANG_VIEW, gocui.KeyArrowDown, gocui.ModNone, goDownLang); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(LANG_VIEW, gocui.KeyArrowUp, gocui.ModNone, goUpLang); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(LANG_VIEW, gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {

		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	widthTerm, heightTerm := g.Size()

	relWidth, _ := relativeSize(g)
	if langView, err := g.SetView(LANG_VIEW, 0, 0, relWidth, heightTerm-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		langView.Highlight = true
		langView.SelBgColor = gocui.ColorGreen
		langView.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(langView, "clojure")
		fmt.Fprintln(langView, "go")
		fmt.Fprintln(langView, "elixir")
	}

	if mainView, err := g.SetView(MAIN_VIEW, relWidth+1, 0, widthTerm-1, heightTerm-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		mainView.Wrap = true
	}
	if _, err := g.SetCurrentView(LANG_VIEW); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// needs refactoring
func goDown(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View(MAIN_VIEW)
	cx, cy := mainView.Cursor()
	if err := mainView.SetCursor(cx, cy+1); err != nil {
		ox, oy := mainView.Origin()
		if err := mainView.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func goDownLang(g *gocui.Gui, v *gocui.View) error {
	//log.Println("GODOWNLANG: ", v)
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func goUpLang(g *gocui.Gui, v *gocui.View) error {
	//log.Println("GOUPLANG: ", v)
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}
	active = nextIndex
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	log.Println("Searching for language: ", l)

	var reposChan chan string
	reposChan = make(chan string)
	go GetTrendingRepos(l, "daily", reposChan)

	g.Update(func(g *gocui.Gui) error {
		mainView, _ := g.View(MAIN_VIEW)
		mainView.Clear()
		fmt.Fprintln(mainView, <-reposChan)
		return nil
	})

	return nil
}
