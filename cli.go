package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/emaincourt/codeship-cli/graphics"
	"github.com/emaincourt/codeship-cli/providers"
	"github.com/gizak/termui"
)

func main() {
	organization := flag.String("org", "", "The name of the organization.")
	username := flag.String("username", os.Getenv("CODESHIP_USERNAME"), "Codeship username.")
	password := flag.String("password", os.Getenv("CODESHIP_PASSWORD"), "Codeship password.")

	flag.Parse()

	if *username == "" {
		panic(fmt.Errorf("username must be provided either from env var CODESHIP_USERNAME or --username"))
	}

	if *password == "" {
		panic(fmt.Errorf("password must be provided either from env var CODESHIP_PASSWORD or --password"))
	}

	provider, err := providers.NewCodeShipProviderFromCredentials(*username, *password, *organization)
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
