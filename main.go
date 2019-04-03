package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
	"olympos.io/encoding/edn"
)

const (
	LANG_VIEW   = "languages"
	MAIN_VIEW   = "main"
	PROMPT_VIEW = "prompt"
)

type Config struct {
	Languages []string `edn:"languages"`
}

var (
	viewArr = []string{LANG_VIEW, MAIN_VIEW}
	active  = 0
	c       Config
)

func relativeSize(g *gocui.Gui) (int, int) {
	tw, th := g.Size()
	return (tw * 3) / 10, (th * 70) / 100
}

func setKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(MAIN_VIEW, gocui.KeyArrowDown, gocui.ModNone, goDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(MAIN_VIEW, gocui.KeyArrowUp, gocui.ModNone, goUp); err != nil {
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

	if err := g.SetKeybinding(LANG_VIEW, gocui.KeyEnter, gocui.ModNone, fetchLangRepos); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, addLang); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(PROMPT_VIEW, gocui.KeyEnter, gocui.ModNone, addNewLang); err != nil {
		log.Panicln(err)
	}

	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
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

	dat, err := ioutil.ReadFile("config.edn")
	if err != nil {
		log.Fatal(err)
	}

	//var c Config // temporarily will be a global variable
	err = edn.Unmarshal(dat, &c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(c.Languages)

	setKeyBindings(g)

	c.Languages = append(c.Languages, "foobar")
	//bs, _ := edn.Marshal(c.Languages)
	bs, _ := edn.MarshalPPrint(c, nil)
	log.Println(string(bs))

	configFile, err := os.OpenFile("config.edn", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	_, err = configFile.WriteAt(bs, 0)
	if err != nil {
		log.Fatal(err)
	}

	// important to set focus on the languates panel right on the start
	go g.Update(func(g *gocui.Gui) error {
		v, _ := g.SetCurrentView(LANG_VIEW)
		for _, lang := range c.Languages {
			fmt.Fprintln(v, lang)
		}
		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func createPromptView(g *gocui.Gui) error {
	widthT, heightT := g.Size()
	v, err := g.SetView(PROMPT_VIEW, widthT/6, (heightT/2)-1, (widthT*5)/6, (heightT/2)+1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Editable = true

	g.Cursor = true
	//g.Highlight = true

	//t, err := g.SetCurrentView(PROMPT_VIEW)
	_, err = g.SetCurrentView(PROMPT_VIEW)
	if err != nil {
		log.Println("error setting the current view to prompt view")
		return err
	}
	//log.Println(t)

	g.SetViewOnTop(PROMPT_VIEW)

	return nil
}

func addLang(g *gocui.Gui, v *gocui.View) error {
	createPromptView(g)
	return nil
}

func addNewLang(g *gocui.Gui, v *gocui.View) error {
	newLang := strings.TrimSpace(v.Buffer())
	log.Println("Adding new language: ", newLang)
	langView, err := g.View(LANG_VIEW)
	if err != nil {
		return err
	}
	fmt.Fprintln(langView, newLang)
	//saveLang(newLang)
	g.DeleteView(PROMPT_VIEW)
	g.SetCurrentView(LANG_VIEW)
	return nil
}

func createLangView(g *gocui.Gui) error {
	_, heightTerm := g.Size()
	relWidth, _ := relativeSize(g)

	if langView, err := g.SetView(LANG_VIEW, 0, 0, relWidth, heightTerm-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		langView.Highlight = true
		langView.SelBgColor = gocui.ColorGreen
		langView.SelFgColor = gocui.ColorBlack
	}

	//	_, err := g.SetCurrentView(LANG_VIEW)
	//	if err != nil {
	//		log.Println("erro lang view")
	//	}
	return nil
}

func createMainView(g *gocui.Gui) error {
	widthTerm, heightTerm := g.Size()
	relWidth, _ := relativeSize(g)

	if mainView, err := g.SetView(MAIN_VIEW, relWidth+1, 0, widthTerm-1, heightTerm-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		mainView.Wrap = true
	}
	return nil
}

func layout(g *gocui.Gui) error {
	createMainView(g)
	createLangView(g)
	//log.Println("CUrrent view: ", g.CurrentView)
	//v := g.CurrentView()
	//if v != nil {
	//log.Println(v.Name())
	//} else {
	//log.Println("nenhum foco")
	//}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// need to refactor
func goDown(g *gocui.Gui, v *gocui.View) error {
	//mainView, _ := g.View(MAIN_VIEW)
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
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

//it has a bug when it reaches the top of the panel, check out
func goUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	log.Println("cy: ", cy)
	if err := v.SetCursor(cx, cy-1); err != nil {
		ox, oy := v.Origin()
		log.Println("oy: ", oy)
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

func getCurrLine(v *gocui.View) string {
	_, cy := v.Cursor()
	l, err := v.Line(cy)
	if err != nil {
		return ""
	}
	return l
}

func fetchLangRepos(g *gocui.Gui, v *gocui.View) error {
	currLang := getCurrLine(v)
	log.Println("Searching for language: ", currLang)

	reposChan := make(chan string)
	go GetTrendingRepos(currLang, "daily", reposChan)

	mainView, err := g.View(MAIN_VIEW)
	if err != nil {
		return err
	}
	go updateView(g, mainView, <-reposChan)
	return nil
}

func updateView(g *gocui.Gui, v *gocui.View, content string) error {
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		fmt.Fprintln(v, content)
		return nil
	})
	return nil
}
