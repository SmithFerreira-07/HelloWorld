package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	message         = "Hello, World!"
	framesPerSecond = 30
	framesPerLetter = 2
	sidePadding     = 12
)

var boxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#666666")).
	Padding(3, sidePadding).
	Width(len(message) + 2*sidePadding)

type tickMsg time.Time

type model struct {
	frame  int
	width  int
	height int
}

// Init function
func (m model) Init() tea.Cmd {
	return waitForNextFrame()
}

// wait for next frame
func waitForNextFrame() tea.Cmd {
	return tea.Tick(time.Second/framesPerSecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update animation frame and handle messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.frame++
		return m, waitForNextFrame()
	}

	return m, nil
}

// View render current frame
func (m model) View() string {
	lettersToShow := m.frame / framesPerLetter
	if lettersToShow > len(message) {
		lettersToShow = len(message)
	}

	var text strings.Builder
	for position, letter := range message[:lettersToShow] {
		letterStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(shadeOfGray(position, m.frame))
		text.WriteString(letterStyle.Render(string(letter)))
	}

	box := boxStyle.Render(text.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// gray shimmer effect
func shadeOfGray(position, frame int) lipgloss.Color {
	wave := math.Sin(float64(position)*0.5 - float64(frame)*0.2)
	brightness := 128 + int(127*(wave+1)/2)
	channel := fmt.Sprintf("%02x", brightness)
	return lipgloss.Color("#" + channel + channel + channel)
}

func main() {
	program := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Println("Something went wrong:", err)
		os.Exit(1)
	}
}
