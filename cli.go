package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/emaincourt/codeship-cli/graphics"
	"github.com/emaincourt/codeship-cli/providers"
	"github.com/gizak/termui"
)

func main() {
	organization := flag.String("org", "", "The name of the organization.")

	flag.Parse()

	codeshipUser := os.Getenv("CODESHIP_USERNAME")
	codeshipPass := os.Getenv("CODESHIP_PASSWORD")

	provider, err := providers.NewCodeShipProviderFromCredentials(codeshipUser, codeshipPass, *organization)
	if err != nil {
		panic(err)
	}

	termUI, err := graphics.NewTerm()
	if err != nil {
		panic(err)
	}

	err = termUI.Start()
	if err != nil {
		panic(err)
	}
	defer termUI.Close()

	header, err := provider.GetHeader()
	if err != nil {
		panic(err)
	}
	termUI.SetHeader(header)

	projects, err := provider.GetProjectsList()
	if err != nil {
		panic(err)
	}
	termUI.SetProjectsList(projects)

	latestKBInput := time.Now()
	currentProjectIndex := 0

	termUI.AddTimeLoop(func(e termui.Event) {
		currentProjectID, err := provider.GetProjectIDFromIndex(currentProjectIndex)
		if err != nil {
			panic(err)
		}

		builds, err := provider.GetBuildsList(currentProjectID)
		if err != nil {
			panic(err)
		}

		termUI.SetBuildsList(builds)
		termUI.Clear()
		termUI.Render()
	})

	termUI.AddKeyPressHandler(func(e termui.Event) {
		if time.Now().Sub(latestKBInput) < 1000000000 {
			i, err := strconv.Atoi(strconv.Itoa(currentProjectIndex) + e.Data.(termui.EvtKbd).KeyStr)
			if err == nil {
				currentProjectIndex = i
			}
		} else {
			i, err := strconv.Atoi(e.Data.(termui.EvtKbd).KeyStr)
			if err == nil {
				currentProjectIndex = i
			}
		}
		latestKBInput = time.Now()
	})

	termUI.Loop()
}
