/*
 * INDVIM - Indonesia Text Editor
 * Created by: Nasa (hastagaming) - Dsn. Bangi, Kediri
 * License: MIT License
 * Year: 2026
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type FileNode struct {
	Name  string
	IsDir bool
}

var (
	lines         = []string{""}
	cursorX       = 0
	cursorY       = 0
	filename      = "" // Default kosong jika tidak ada argumen
	
	currentMode   = "VIEW" 
	commandBuffer = ""     

	// Fitur Tree Explorer
	showTree      = false
	isTreeFocused = false
	currDir       = "."
	treeNodes     = []FileNode{}
	treeCursor    = 0
	treeWidth     = 22

	logo = []string{
		"  ___ _   _ ______     _____ __  __ ",
		" |_ _| \\ | |  _ \\ \\   / /_ _|  \\/  |",
		"  | ||  \\| | | | \\ \\ / / | || |\\/| |",
		"  | || |\\  | |_| |\\ V /  | || |  | |",
		" |___|_| \\_|____/  \\_/  |___|_|  |_|",
		"                                    ",
		"      [ Indonesia Text Editor ]     ",
		"        Created by Nasa (2026)      ",
	}

	snippets = map[string]string{
		"!kt":   "fun main() {\n    println(\"Hello Nasa!\")\n}",
		"!java": "public class Main {\n    public static void main(String[] args) {\n        \n    }\n}",
		"!html": "<!DOCTYPE html>\n<html>\n<head><title>INDVIM</title></head>\n<body>\n    \n</body>\n</html>",
		"!py":   "def main():\n    print('INDVIM is the best!')\n\nif __name__ == '__main__':\n    main()",
	}

	keywords = map[string]tcell.Color{
		"func": tcell.ColorOrange, "package": tcell.ColorDeepPink, "import": tcell.ColorYellow,
		"var": tcell.ColorTeal, "if": tcell.ColorPurple, "else": tcell.ColorPurple,
		"return": tcell.ColorRed, "println": tcell.ColorBlue, "val": tcell.ColorTeal,
		"fun": tcell.ColorOrange, "def": tcell.ColorOrange, "class": tcell.ColorYellow,
		"print": tcell.ColorBlue, "main": tcell.ColorGreen,
	}
)

func main() {
	if len(os.Args) > 1 {
		filename = os.Args[1]
		info, err := os.Stat(filename)
		// Jika argumen adalah folder, langsung buka Tree
		if err == nil && info.IsDir() {
			currDir = filename
			filename = ""
			showTree = true
			isTreeFocused = true
			loadDir(currDir)
		} else {
			loadFromFile()
		}
	} else {
		filename = ""
	}

	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer s.Fini()

	loadDir(currDir)

	for {
		draw(s)
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// Toggle Explorer Tree dengan Ctrl + E
			if ev.Key() == tcell.KeyCtrlE {
				showTree = !showTree
				isTreeFocused = showTree
				if showTree { loadDir(currDir) }
				continue
			}

			// Tombol darurat
			if ev.Key() == tcell.KeyCtrlB {
				saveToFile()
				if filename != "" { return } // Hanya keluar jika berhasil disave
				continue
			}

			handleInput(ev)
			
			if currentMode == "INSERT" && filename != "" {
				saveToFile()
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func loadDir(path string) {
	entries, err := os.ReadDir(path)
	treeNodes = []FileNode{}
	if err != nil { return }
	
	// Tambahkan opsi ".." untuk kembali ke folder sebelumnya
	if path != "." && path != "/" {
		treeNodes = append(treeNodes, FileNode{Name: "..", IsDir: true})
	}
	for _, e := range entries {
		treeNodes = append(treeNodes, FileNode{Name: e.Name(), IsDir: e.IsDir()})
	}
	treeCursor = 0
}

func loadFromFile() {
	data, err := os.ReadFile(filename)
	if err != nil {
		lines = []string{""}
		return
	}
	content := string(data)
	if content == "" {
		lines = []string{""}
	} else {
		lines = strings.Split(content, "\n")
	}
	cursorX, cursorY = 0, 0
}

func saveToFile() {
	if filename == "" {
		commandBuffer = "Error: Cannot save without filename. Use :w <name>"
		currentMode = "COMMAND"
		return
	}
	content := strings.Join(lines, "\n")
	os.WriteFile(filename, []byte(content), 0644)
}

func draw(s tcell.Screen) {
	s.Clear()
	w, h := s.Size()

	styleLineNum := tcell.StyleDefault.Foreground(tcell.ColorDimGray)
	styleText := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	offsetX := 0

	// 🌲 GAMBAR FILE TREE EXPLORER
	if showTree {
		offsetX = treeWidth
		for y := 0; y < h-1; y++ {
			s.SetContent(offsetX-1, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorDimGray))
		}

		for i, node := range treeNodes {
			if i >= h-2 { break }
			
			nodeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
			if node.IsDir {
				nodeStyle = tcell.StyleDefault.Foreground(tcell.ColorDodgerBlue).Bold(true)
			}
			
			if isTreeFocused && i == treeCursor {
				nodeStyle = nodeStyle.Background(tcell.ColorDimGray).Foreground(tcell.ColorWhite)
			}

			displayName := node.Name
			if len(displayName) > treeWidth-3 {
				displayName = displayName[:treeWidth-5] + ".."
			}

			prefix := "  "
			if node.IsDir { prefix = "📁 " } else { prefix = "📄 " }

			rowText := prefix + displayName
			for x, char := range rowText {
				s.SetContent(x, i, char, nil, nodeStyle)
			}
		}
	}

	// Logo ASCII (Hanya jika file kosong dan Tree tertutup)
	if len(lines) == 1 && lines[0] == "" && currentMode != "COMMAND" && !showTree {
		startX := (w - len(logo[0])) / 2
		startY := (h - len(logo)) / 2
		logoStyle := tcell.StyleDefault.Foreground(tcell.ColorMediumSpringGreen).Bold(true)
		for y, line := range logo {
			for x, char := range line {
				s.SetContent(startX+x, startY+y, char, nil, logoStyle)
			}
		}
	}

	// Render Teks & Nomor Baris Editor
	for y, line := range lines {
		if y >= h-2 { break }
		
		num := fmt.Sprintf(" %2d │ ", y+1)
		for i, char := range num {
			s.SetContent(offsetX+i, y, char, nil, styleLineNum)
		}

		wordStart := 0
		currX := offsetX + 6
		for i := 0; i <= len(line); i++ {
			if i == len(line) || line[i] == ' ' {
				word := line[wordStart:i]
				color, isKw := keywords[word]
				wordStyle := styleText
				if isKw {
					wordStyle = tcell.StyleDefault.Foreground(color).Bold(true)
				}
				for j := wordStart; j < i; j++ {
					s.SetContent(currX, y, rune(line[j]), nil, wordStyle)
					currX++
				}
				if i < len(line) {
					s.SetContent(currX, y, ' ', nil, styleText)
					currX++
				}
				wordStart = i + 1
			}
		}
	}

	// 🎨 STATUS BAR KUSTOM NASA
	displayFilename := filename
	if displayFilename == "" {
		displayFilename = "[No Name]"
	}

	var statusLeft string
	var barStyle tcell.Style

	if currentMode == "COMMAND" {
		statusLeft = fmt.Sprintf(" :%s ", commandBuffer)
		barStyle = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack).Bold(true)
	} else if currentMode == "INSERT" {
		statusLeft = " [ INSERT ] | name: "
		barStyle = tcell.StyleDefault.Background(tcell.ColorMediumSpringGreen).Foreground(tcell.ColorBlack).Bold(true)
	} else {
		statusLeft = " [ VIEW ] | name: "
		barStyle = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(true)
	}

	// Gambar Background Status Bar
	for i := 0; i < w; i++ { s.SetContent(i, h-1, ' ', nil, barStyle) }

	// Gambar status kiri
	currPos := 0
	for _, char := range statusLeft {
		s.SetContent(currPos, h-1, char, nil, barStyle)
		currPos++
	}

	// Gambar Nama File dengan WARNA MERAH
	// Perbaikan: Ambil warna background dari barStyle secara manual
	_, bg, _ := barStyle.Decompose() 
	redStyle := tcell.StyleDefault.Background(bg).Foreground(tcell.ColorRed).Bold(true)
	
	for _, char := range displayFilename {
		s.SetContent(currPos, h-1, char, nil, redStyle)
		currPos++
	}

	statusRight := " | TREE: Ctrl+E "
	for _, char := range statusRight {
		s.SetContent(currPos, h-1, char, nil, barStyle)
		currPos++
	}

	// Logika Kursor
	if isTreeFocused {
		s.HideCursor()
	} else if currentMode == "COMMAND" {
		s.ShowCursor(len(statusLeft)-2+len(commandBuffer), h-1) // -2 to account for spaces
	} else {
		s.ShowCursor(offsetX+6+cursorX, cursorY)
	}
	s.Show()
}

func handleInput(ev *tcell.EventKey) {
	// LOGIKA TREE EXPLORER
	if isTreeFocused {
		switch ev.Key() {
		case tcell.KeyUp: if treeCursor > 0 { treeCursor-- }
		case tcell.KeyDown: if treeCursor < len(treeNodes)-1 { treeCursor++ }
		case tcell.KeyLeft: isTreeFocused = false // Pindah ke editor
		case tcell.KeyEnter:
			node := treeNodes[treeCursor]
			if node.IsDir {
				currDir = filepath.Join(currDir, node.Name)
				loadDir(currDir)
			} else {
				filename = filepath.Join(currDir, node.Name)
				loadFromFile()
				isTreeFocused = false
				currentMode = "VIEW"
			}
		}
		return
	}

	// KEMBALI KE TREE JIKA TERTUTUP (Navigasi Kiri dari ujung Editor)
	if ev.Key() == tcell.KeyLeft && cursorX == 0 && showTree && currentMode == "VIEW" {
		isTreeFocused = true
		return
	}

	// LOGIKA MODE COMMAND
	if currentMode == "COMMAND" {
		if ev.Key() == tcell.KeyEscape {
			currentMode = "VIEW"
			commandBuffer = ""
		} else if ev.Key() == tcell.KeyEnter {
			executeCommand()
		} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
			if len(commandBuffer) > 0 {
				commandBuffer = commandBuffer[:len(commandBuffer)-1]
			} else {
				currentMode = "VIEW"
			}
		} else if ev.Key() == tcell.KeyRune {
			commandBuffer += string(ev.Rune())
		}
		return
	}

	// NAVIGASI (Bisa di View & Insert)
	switch ev.Key() {
	case tcell.KeyUp: if cursorY > 0 { cursorY-- }
	case tcell.KeyDown: if cursorY < len(lines)-1 { cursorY++ }
	case tcell.KeyLeft: if cursorX > 0 { cursorX-- }
	case tcell.KeyRight: if cursorX < len(lines[cursorY]) { cursorX++ }
	case tcell.KeyEscape: currentMode = "VIEW"
	}

	// NEOVIM BINDINGS
	if currentMode == "VIEW" {
		if ev.Key() == tcell.KeyRune {
			if ev.Rune() == 'i' {
				currentMode = "INSERT"
			} else if ev.Rune() == ':' {
				currentMode = "COMMAND"
				commandBuffer = ""
			}
		}
		return
	}

	// LOGIKA MODE INSERT
	if currentMode == "INSERT" {
		switch ev.Key() {
		case tcell.KeyEnter:
			rem := lines[cursorY][cursorX:]
			lines[cursorY] = lines[cursorY][:cursorX]
			cursorY++
			lines = append(lines[:cursorY], append([]string{rem}, lines[cursorY:]...)...)
			cursorX = 0
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if cursorX > 0 {
				lines[cursorY] = lines[cursorY][:cursorX-1] + lines[cursorY][cursorX:]
				cursorX--
			} else if cursorY > 0 {
				prevLen := len(lines[cursorY-1])
				lines[cursorY-1] += lines[cursorY]
				lines = append(lines[:cursorY], lines[cursorY+1:]...)
				cursorY--
				cursorX = prevLen
			}
		case tcell.KeyRune:
			char := ev.Rune()
			lines[cursorY] = lines[cursorY][:cursorX] + string(char) + lines[cursorY][cursorX:]
			cursorX++
			if char == ' ' { checkSnippet() }
		}
		if cursorX > len(lines[cursorY]) { cursorX = len(lines[cursorY]) }
	}
}

func executeCommand() {
	parts := strings.Split(strings.TrimSpace(commandBuffer), " ")
	if len(parts) == 0 { return }

	cmd := parts[0]
	
	switch cmd {
	case "w":
		if len(parts) > 1 { filename = parts[1] }
		saveToFile()
	case "q":
		os.Exit(0)
	case "wq":
		if len(parts) > 1 { filename = parts[1] }
		saveToFile()
		if filename != "" { os.Exit(0) }
	}
	commandBuffer = ""
	currentMode = "VIEW"
}

func checkSnippet() {
	if cursorX < 2 { return }
	words := strings.Fields(lines[cursorY][:cursorX-1])
	if len(words) == 0 { return }
	lastWord := words[len(words)-1]

	if val, ok := snippets[lastWord]; ok {
		startPos := cursorX - len(lastWord) - 1
		lines[cursorY] = lines[cursorY][:startPos] + lines[cursorY][cursorX:]
		cursorX = startPos
		snipLines := strings.Split(val, "\n")
		lines[cursorY] = lines[cursorY][:cursorX] + snipLines[0]
		for i := 1; i < len(snipLines); i++ {
			cursorY++
			lines = append(lines[:cursorY], append([]string{snipLines[i]}, lines[cursorY:]...)...)
		}
		cursorX = len(lines[cursorY])
	}
}
