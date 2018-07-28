package graphics

import (
	"github.com/emaincourt/codeship-cli/providers"
	"github.com/gizak/termui"
)

// Term refers to a type of monitoring graphical interface that leverages the terminal
type Term struct {
	Graphics

	Header   *termui.Par
	Projects *termui.List
	Builds   *termui.Table
}

// NewTerm instantiates a new empty Term
func NewTerm() (*Term, error) {
	err := termui.Init()
	if err != nil {
		return nil, err
	}

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	return &Term{
		Header:   &termui.Par{},
		Projects: &termui.List{},
		Builds:   &termui.Table{},
	}, nil
}

// SetHeader sets the header of the Term as a simple string
func (t *Term) SetHeader(content string) error {
	block := *termui.NewBlock()
	block.Height = 4
	block.Width = 17
	block.BorderLabel = "Codeship CLI 0.0.1"

	t.Header = &termui.Par{
		Block: block,
		Text:  content,
	}

	return nil
}

// SetProjectsList attaches the projects list to the given Term instance
func (t *Term) SetProjectsList(projects []string) error {
	block := *termui.NewBlock()
	block.BorderLabel = "Projects"
	block.Height = len(projects)

	t.Projects = &termui.List{
		Block:       block,
		Items:       projects,
		ItemFgColor: termui.ColorYellow,
		Overflow:    "hidden",
	}

	return nil
}

// SetBuildsList attaches the builds list to the given Term instance
func (t *Term) SetBuildsList(builds []providers.Build) error {
	rows := [][]string{
		[]string{"Started At", "Finished At", "Commit Message", "Status", "Username"},
	}

	for _, build := range builds {
		rows = append(rows, []string{
			build.StartedAt,
			build.FinishedAt,
			build.CommitMessage,
			build.Status,
			build.Username,
		})
	}

	table := termui.NewTable()
	table.Rows = rows
	table.FgColor = termui.ColorWhite
	table.BgColor = termui.ColorDefault
	table.Separator = false
	table.CellWidth = []int{50, 50, 100, 10, 20}
	table.Analysis()
	table.SetSize()
	table.BgColors[0] = termui.ColorRed

	t.Builds = table

	return nil
}

// Render renders the underlying termui
func (t *Term) Render() error {
	termui.Body.Rows = []*termui.Row{
		termui.NewRow(
			termui.NewCol(3, 0, t.Header),
		),
		termui.NewRow(
			termui.NewCol(3, 0, t.Projects),
			termui.NewCol(9, 0, t.Builds),
		),
	}

	termui.Body.Align()
	termui.Render(termui.Body)

	return nil
}

// AddTimeLoop starts a time loop in the underlying termui
func (t *Term) AddTimeLoop(fn func(e termui.Event)) error {
	termui.Handle("/timer/1s", fn)

	return nil
}

// AddKeyPressHandler add a key press handler to the underlying termui
func (t *Term) AddKeyPressHandler(fn func(e termui.Event)) error {
	termui.Handle("/sys/kbd", fn)

	return nil
}

// Loop starts the loop of the underlying termui
func (t *Term) Loop() {
	termui.Loop()
}

// Close closes the loop of the underlying termui
func (t *Term) Close() {
	termui.Close()
}
