package graphics

import (
	"github.com/emaincourt/codeship-cli/providers"
	"github.com/gizak/termui"
)

// Graphics references the interface of a graphical environment to monitor builds
type Graphics interface {
	SetHeader(content string) error
	SetProjectsList(projects []string) error
	SetBuildsList(builds []providers.Build) error
	Render() error

	AddTimeLoop(func(e termui.Event)) error
	AddKeyPressHandler(func(e termui.Event)) error
	Loop()
	Close()
}
