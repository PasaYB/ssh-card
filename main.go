package main

import (
    "fmt"
    "os"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    wish "github.com/charmbracelet/wish"
    bm "github.com/charmbracelet/wish/bubbletea"
    "github.com/charmbracelet/lipgloss/tree"
    "context"
    "os/signal"
    "syscall"
    "time"
    "github.com/charmbracelet/ssh"
    "github.com/charmbracelet/bubbles/paginator"
    // "github.com/charmbracelet/lipgloss/list"
)

const totalPages = 2

// Struct Model
type model struct{
    width       int
    height      int
    currentPage int
    paginator   paginator.Model
    asciiArt    string
}

func section(title, desc, color string) string {
    t := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color(color)).
        MarginBottom(1).
        Underline(true).
        Render("# " + title)

    d := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#ffffff")).
        MarginBottom(2).
        Render(strings.TrimSpace(desc))

    return lipgloss.JoinVertical(lipgloss.Left, t, d)
}

func renderHome(m model) string {
    // Left block style
    leftStyle := lipgloss.NewStyle().
        Width(60).
        Padding(0, 0)
        
    // Right block style
    rightStyle := lipgloss.NewStyle().
        Width(50).
        Padding(2, 1)

    block1 := section("Lanang Gading Pasa | Tech Enthusiast", "I'm an Informatics student who loves exploring technology and building web-based applications. Interested in Fullstack Development, Databases, and learning new frameworks.", "#5bdb00")

    rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00d3c2"))
    itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#bbff00"))

    t := tree.New().
        Root("⁜ What im using").
        Child(
            "PHP",
            tree.New().Child(
                "Laravel",
            ),
            "Javascript",
            tree.New().Child(
                "Vue",
                "Nuxt",
            ),
            "C++",
        ).
        Enumerator(tree.RoundedEnumerator).
        RootStyle(rootStyle).
        ItemStyle(itemStyle)

    tools := lipgloss.NewStyle().Render(t.String())

    footer := lipgloss.NewStyle().
            Foreground(lipgloss.Color("#ffffff")).
            Width(40).
            MarginTop(3).
            MarginLeft(10).
            Render("Thanks For Visiting!")

    rightContent := lipgloss.JoinVertical(lipgloss.Left, block1, tools, footer)
    
    // Render each block
    left  := leftStyle.Render(m.asciiArt)
    right := rightStyle.Render(rightContent)

    main := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

    return lipgloss.NewStyle().
        Render(main)
}

func renderSocials() string {
    return "sup"
}

func renderHints() string {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("#555555")).
        Render("← → navigate   q quit")
}

// Init
func (m model) Init() tea.Cmd {
    return nil
}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tea.KeyMsg:
        switch msg.String() {
    case "right":
        m.currentPage = (m.currentPage + 1) % totalPages 
        m.paginator.Page = m.currentPage
    case "left":
        m.currentPage = (m.currentPage - 1 + totalPages) % totalPages  
        m.paginator.Page = m.currentPage
    case "q":
        return m, tea.Quit
    }
    }
    return m, nil
}

// View
func (m model) View() string {
    var content string

    switch m.currentPage {
    case 0:
        content = renderHome(m)
    case 1:
        content = renderSocials()
    default:
        content = renderHome(m)
    }

    dots := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#874BFD")).
        Render(m.paginator.View())

    hints := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#555555")).
        Render("← → navigate   q quit")

    bottomBar := lipgloss.JoinHorizontal(
        lipgloss.Center,
        dots,
        lipgloss.NewStyle().Width(4).Render(""),
        hints,
    )
    
    outerStyle := lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("#d6a400")).
        Padding(1, 2)

    bordered := outerStyle.Render(content)

    full := lipgloss.JoinVertical(lipgloss.Center, bordered, bottomBar)

    paddingTop := (m.height - lipgloss.Height(bordered)) / 2

    return lipgloss.NewStyle().
        Width(m.width).
        PaddingTop(paddingTop).
        Align(lipgloss.Center).
        Render(full)
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    ascii, err := os.ReadFile("ascii.txt")
    if err != nil {
        ascii = []byte("no ascii art found")  // fallback
    }

    p := paginator.New()
    p.Type = paginator.Dots
    p.ActiveDot = "x"
    p.InactiveDot = "o"
    p.SetTotalPages(totalPages)

    m := model{
        paginator: p,
        asciiArt:  string(ascii),
    }

    return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func main() {
    // m := model{width: 120, height: 40}  // fake terminal size
    // fmt.Println(m.View())
    // os.Exit(0)  // stop before starting SSH server

    s, err := wish.NewServer(
        wish.WithAddress("localhost:2222"),
        wish.WithHostKeyPath(".ssh/id_ed25519"),
        wish.WithMiddleware(
            bm.Middleware(teaHandler),
        ),
    )
    if err != nil {
        fmt.Println("Error starting server:", err)
        os.Exit(1)
    }

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := s.ListenAndServe(); err != nil {
            fmt.Println("Server stopped:", err)
        }
    }()

    fmt.Println("SSH server running on localhost:2222")
    fmt.Println("Connect with: ssh localhost -p 2222")
    <-done

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    s.Shutdown(ctx)
}