package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/process"
)

type ghostProc struct {
	PID    int32
	Name   string
	X      int // horizontal position
	V      int
	CPU    float64
	Spooky bool
}

type model struct {
	procs []ghostProc
	err   error
}

var spookyNames = []string{
	"🧟 Zombie", "👻 Phantom", "💀 Skull", "🕷 Spider", "🦇 Bat", "🔥 Cursed",
}

func getProcesses() ([]ghostProc, error) {
	ps, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var gprocs []ghostProc
	for _, p := range ps {
		name, err := p.Name()
		if err != nil {
			name = "unknown"
		}
		cpu, _ := p.CPUPercent()
		gprocs = append(gprocs, ghostProc{PID: p.Pid, Name: name, X: rand.Intn(20), CPU: cpu})
	}
	return gprocs, nil
}

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg { return "tick" })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case string:
		if msg == "tick" {
			for i := range m.procs {
				// move horizontally
				m.procs[i].X++
				if m.procs[i].X > 50 {
					m.procs[i].X = 0
				}
				// randomly rename some processes
				if rand.Float32() < 0.05 {
					m.procs[i].Name = spookyNames[rand.Intn(len(spookyNames))]
					m.procs[i].Spooky = true
				} else {
					m.procs[i].Spooky = false
				}
			}
			return m, tick()
		}
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
	}

	output := ""
	for _, p := range m.procs {
		spaces := ""
		for i := 0; i < p.X; i++ {
			spaces += " "
		}
		red := "\033[31m"
		reset := "\033[0m"
		name := p.Name
		if p.CPU > 20 || p.Spooky { // high CPU or spooky flash
			name = red + name + reset
		}
		output += fmt.Sprintf("%s%d: %s\n", spaces, p.PID, name)
	}
	output += "\nPress q to quit."
	return output
}

func main() {
	rand.Seed(time.Now().UnixNano())
	procs, err := getProcesses()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	p := tea.NewProgram(model{procs: procs, err: nil})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
