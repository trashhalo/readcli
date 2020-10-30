package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/go-shiori/go-readability"
	"github.com/mattn/go-runewidth"

	"github.com/urfave/cli/v2"
)

func main() {
	ctx := context.Background()
	app := cli.App{
		Name:  "readcli",
		Usage: "read website content on the command line",
		Action: func(c *cli.Context) error {
			url := c.Args().First()
			if url == "" {
				return fmt.Errorf("url is required input")
			}

			article, err := readability.FromURL(url, 30*time.Second)
			if err != nil {
				return err
			}

			converter := md.NewConverter("", true, nil)
			in, err := converter.ConvertString(article.Content)
			if err != nil {
				return err
			}

			out, err := glamour.Render(in, "dark")
			if err != nil {
				return err
			}

			p := tea.NewProgram(model{title: article.Title, content: out})

			p.EnterAltScreen()
			defer p.ExitAltScreen()

			p.EnableMouseCellMotion()
			defer p.DisableMouseCellMotion()

			return p.Start()
		},
	}
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		panic(err)
	}

}

type model struct {
	content  string
	title    string
	ready    bool
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

const (
	headerHeight = 3
	footerHeight = 3
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Ctrl+c exits
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		verticalMargins := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we need
			// to wait until we've received the window dimensions before we
			// can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.Model{Width: msg.Width, Height: msg.Height - verticalMargins}
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargins
		}
	}

	m.viewport, _ = viewport.Update(msg, m.viewport)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initalizing..."
	}

	headerTop := "╭─"
	for i := 0; i < len(m.title); i++ {
		headerTop += "─"
	}
	headerTop += "─╮"
	headerMid := fmt.Sprintf("│ %s ├", m.title)
	headerBot := "╰─"
	for i := 0; i < len(m.title); i++ {
		headerBot += "─"
	}
	headerBot += "─╯"
	headerMid += strings.Repeat("─", m.viewport.Width-runewidth.StringWidth(headerMid))
	header := fmt.Sprintf("%s\n%s\n%s", headerTop, headerMid, headerBot)

	footerTop := "╭──────╮"
	footerMid := fmt.Sprintf("┤ %3.f%% │", m.viewport.ScrollPercent()*100)
	footerBot := "╰──────╯"
	gapSize := m.viewport.Width - runewidth.StringWidth(footerMid)
	footerTop = strings.Repeat(" ", gapSize) + footerTop
	footerMid = strings.Repeat("─", gapSize) + footerMid
	footerBot = strings.Repeat(" ", gapSize) + footerBot
	footer := fmt.Sprintf("%s\n%s\n%s", footerTop, footerMid, footerBot)

	return fmt.Sprintf("%s\n%s\n%s", header, viewport.View(m.viewport), footer)
}
