package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	codeship "github.com/codeship/codeship-go"
	"github.com/gizak/termui"
)

func makeInterface(projectsList *termui.List, buildsTable *termui.Table) {
	termui.Clear()
	termui.Body.Rows = termui.Body.Rows[:0]
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(3, 0, projectsList),
			termui.NewCol(9, 0, buildsTable),
		),
	)
	termui.Body.Align()
	termui.Render(termui.Body)
}

func makeProjectsList(projects []codeship.Project) *termui.List {
	var strs []string

	for index, project := range projects {
		strs = append(strs, fmt.Sprintf("[%d] %s", index, project.Name))
	}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Projects"
	ls.Height = len(strs)

	return ls
}

func makeBuildsTable(builds []codeship.Build) *termui.Table {
	rows := [][]string{
		[]string{"Started At", "Finished At", "Commit Message", "Status", "Username"},
	}

	for _, build := range builds {
		commitMessage := build.CommitMessage
		if len(commitMessage) > 70 {
			commitMessage = build.CommitMessage[:70]
		}

		status := fmt.Sprintf("[%s](fg-cyan)", build.Status)
		if strings.Contains(build.Status, "success") {
			status = fmt.Sprintf("[%s](fg-green)", build.Status)
		}
		if strings.Contains(build.Status, "error") {
			status = fmt.Sprintf("[%s](fg-red)", build.Status)
		}

		rows = append(rows, []string{build.AllocatedAt.Format("02/01/06 03:04:05"), build.FinishedAt.Format("02/01/06 03:04:05"), commitMessage, status, build.Username})
	}

	table := termui.NewTable()
	table.Rows = rows
	table.FgColor = termui.ColorWhite
	table.BgColor = termui.ColorDefault
	table.Separator = false
	table.CellWidth = []int{50, 50, 200, 10, 20}
	table.Analysis()
	table.SetSize()
	table.BgColors[0] = termui.ColorRed

	return table
}

func main() {
	organization := flag.String("org", "", "The name of the organization.")

	flag.Parse()

	codeshipUser := os.Getenv("CODESHIP_USERNAME")
	codeshipPass := os.Getenv("CODESHIP_PASSWORD")

	ctx := context.Background()
	auth := codeship.NewBasicAuth(codeshipUser, codeshipPass)
	currentProjectIndex := 0

	client, err := codeship.New(auth)
	if err != nil {
		panic(err)
	}

	org, err := client.Organization(ctx, *organization)
	if err != nil {
		panic(err)
	}

	projects, _, err := org.ListProjects(ctx)
	if err != nil {
		panic(err)
	}

	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	projectsList := makeProjectsList(projects.Projects)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		builds, _, err := org.ListBuilds(ctx, projects.Projects[currentProjectIndex].UUID)
		if err != nil {
			panic(err)
		}

		buildsList := makeBuildsTable(builds.Builds)

		makeInterface(projectsList, buildsList)
	})

	termui.Handle("/sys/kbd", func(e termui.Event) {
		i, err := strconv.Atoi(e.Data.(termui.EvtKbd).KeyStr)
		if err == nil {
			currentProjectIndex = i
		}
	})

	termui.Loop()
}
