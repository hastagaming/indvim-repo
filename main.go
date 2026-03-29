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
	lineOffset    = 0 // Fitur Scrolling
	filename      = "" 
	
	currentMode   = "VIEW" 
	commandBuffer = ""     
	infoMessage   = "" // Pesan peringatan

	isSelectedAll = false // Fitur Alt+A

	// File Tree Explorer
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
		if err == nil && info.IsDir() {
			currDir = filename
			filename = ""
			showTree, isTreeFocused = true, true
			loadDir(currDir)
		} else {
			loadFromFile()
		}
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
			// LOGIK ALT + I (Siklus Mode)
			if ev.Modifiers()&tcell.ModAlt != 0 && (ev.Rune() == 'i' || ev.Rune() == 'I') {
				if currentMode == "VIEW" { currentMode = "INSERT" } else if currentMode == "INSERT" { currentMode, commandBuffer = "COMMAND", "" } else { currentMode = "VIEW" }
				continue 
			}

			// LOGIK ALT + A (Select All)
			if ev.Modifiers()&tcell.ModAlt != 0 && (ev.Rune() == 'a' || ev.Rune() == 'A') {
				isSelectedAll = !isSelectedAll
				continue
			}

			if ev.Key() == tcell.KeyCtrlE {
				showTree = !showTree
				isTreeFocused = showTree
				if showTree { loadDir(currDir) }
				continue
			}

			if ev.Key() == tcell.KeyCtrlB {
				if filename != "" { saveToFile() ; printToTerminal() ; return }
				continue
			}

			handleInput(ev)
			
			if currentMode == "INSERT" && filename != "" { saveToFile() }
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func loadDir(path string) {
	entries, err := os.ReadDir(path)
	treeNodes = []FileNode{}
	if err != nil { return }
	if path != "." && path != "/" { treeNodes = append(treeNodes, FileNode{Name: "..", IsDir: true}) }
	for _, e := range entries { treeNodes = append(treeNodes, FileNode{Name: e.Name(), IsDir: e.IsDir()}) }
	treeCursor = 0
}

func loadFromFile() {
	data, err := os.ReadFile(filename)
	if err != nil { lines = []string{""} ; return }
	content := string(data)
	if content == "" { lines = []string{""} } else { lines = strings.Split(content, "\n") }
	cursorX, cursorY = 0, 0
}

func saveToFile() {
	if filename == "" { return }
	content := strings.Join(lines, "\n")
	os.WriteFile(filename, []byte(content), 0644)
}

func printToTerminal() {
	fmt.Println("\n--- [ INDVIM Final Output ] ---")
	for _, line := range lines { fmt.Println(line) }
	fmt.Println("--- [ End of Output ] ---\n")
}

func draw(s tcell.Screen) {
	s.Clear()
	w, h := s.Size()
	mainH := h - 1

	// LOGIKA SCROLLING
	if cursorY < lineOffset { lineOffset = cursorY }
	if cursorY >= lineOffset+mainH { lineOffset = cursorY - mainH + 1 }

	styleLineNum := tcell.StyleDefault.Foreground(tcell.ColorDimGray)
	styleText := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	styleSel := tcell.StyleDefault.Background(tcell.ColorDimGray).Foreground(tcell.ColorWhite)

	offsetX := 0
	if showTree {
		offsetX = treeWidth
		for y := 0; y < mainH; y++ { s.SetContent(offsetX-1, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorDimGray)) }
		for i, node := range treeNodes {
			if i >= mainH { break }
			nodeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
			if node.IsDir { nodeStyle = tcell.StyleDefault.Foreground(tcell.ColorDodgerBlue).Bold(true) }
			if isTreeFocused && i == treeCursor { nodeStyle = nodeStyle.Background(tcell.ColorDimGray) }
			prefix := "📁 "
			if !node.IsDir { prefix = "📄 " }
			name := node.Name
			if len(name) > treeWidth-5 { name = name[:treeWidth-7] + ".." }
			row := prefix + name
			for x, char := range row { s.SetContent(x, i, char, nil, nodeStyle) }
		}
	}

	// Tampilkan Logo jika kosong
	if len(lines) == 1 && lines[0] == "" && currentMode != "COMMAND" && !showTree {
		startY := (mainH - len(logo)) / 2
		logoStyle := tcell.StyleDefault.Foreground(tcell.ColorMediumSpringGreen).Bold(true)
		for y, line := range logo {
			startX := (w - len(line)) / 2
			for x, char := range line { s.SetContent(startX+x, startY+y, char, nil, logoStyle) }
		}
	}

	// Render Editor
	for y := 0; y < mainH; y++ {
		lineIdx := y + lineOffset
		if lineIdx >= len(lines) { break }
		line := lines[lineIdx]

		num := fmt.Sprintf(" %2d │ ", lineIdx+1)
		for i, char := range num { s.SetContent(offsetX+i, y, char, nil, styleLineNum) }

		currX, wordStart := offsetX+6, 0
		for i := 0; i <= len(line); i++ {
			if i == len(line) || line[i] == ' ' {
				word := line[wordStart:i]
				color, isKw := keywords[word]
				wordStyle := styleText
				if isSelectedAll { wordStyle = styleSel } else if isKw { wordStyle = tcell.StyleDefault.Foreground(color).Bold(true) }
				
				for j := wordStart; j < i; j++ {
					s.SetContent(currX, y, rune(line[j]), nil, wordStyle)
					currX++
				}
				if i < len(line) {
					spStyle := styleText
					if isSelectedAll { spStyle = styleSel }
					s.SetContent(currX, y, ' ', nil, spStyle)
					currX++
				}
				wordStart = i + 1
			}
		}
	}

	// STATUS BAR
	var statusLeft string
	var barStyle tcell.Style
	if infoMessage != "" {
		statusLeft = " " + infoMessage + " "
		barStyle = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack).Bold(true)
	} else if currentMode == "COMMAND" {
		statusLeft = fmt.Sprintf(" :%s ", commandBuffer)
		barStyle = tcell.StyleDefault.Background(tcell.ColorOrange).Foreground(tcell.ColorBlack).Bold(true)
	} else if currentMode == "INSERT" {
		statusLeft = " [ INSERT ] | name: "
		barStyle = tcell.StyleDefault.Background(tcell.ColorMediumSpringGreen).Foreground(tcell.ColorBlack).Bold(true)
	} else {
		statusLeft = " [ VIEW ] | name: "
		barStyle = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(true)
	}

	for i := 0; i < w; i++ { s.SetContent(i, h-1, ' ', nil, barStyle) }
	currPos := 0
	for _, char := range statusLeft { s.SetContent(currPos, h-1, char, nil, barStyle) ; currPos++ }

	displayFilename := filename
	if displayFilename == "" { displayFilename = "[No Name]" }
	_, bg, _ := barStyle.Decompose()
	redStyle := tcell.StyleDefault.Background(bg).Foreground(tcell.ColorRed).Bold(true)
	for _, char := range displayFilename { s.SetContent(currPos, h-1, char, nil, redStyle) ; currPos++ }

	statusRight := " | ALT+A: Select | Ctrl+E: Tree "
	for _, char := range statusRight { s.SetContent(currPos, h-1, char, nil, barStyle) ; currPos++ }

	if isTreeFocused { s.HideCursor() } else if currentMode == "COMMAND" {
		s.ShowCursor(len(statusLeft)-1+len(commandBuffer), h-1)
	} else { s.ShowCursor(offsetX+6+cursorX, cursorY-lineOffset) }
	s.Show()
}

func handleInput(ev *tcell.EventKey) {
	if isTreeFocused {
		switch ev.Key() {
		case tcell.KeyUp: if treeCursor > 0 { treeCursor-- }
		case tcell.KeyDown: if treeCursor < len(treeNodes)-1 { treeCursor++ }
		case tcell.KeyEnter:
			node := treeNodes[treeCursor]
			if node.IsDir { currDir = filepath.Join(currDir, node.Name) ; loadDir(currDir) } else {
				filename = filepath.Join(currDir, node.Name)
				loadFromFile()
				isTreeFocused, currentMode = false, "VIEW"
			}
		case tcell.KeyLeft: isTreeFocused = false
		}
		return
	}

	if currentMode == "COMMAND" {
		if ev.Key() == tcell.KeyEscape { currentMode, commandBuffer = "VIEW", "" } else if ev.Key() == tcell.KeyEnter { executeCommand() } else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
			if len(commandBuffer) > 0 { commandBuffer = commandBuffer[:len(commandBuffer)-1] }
		} else if ev.Key() == tcell.KeyRune { commandBuffer += string(ev.Rune()) }
		return
	}

	switch ev.Key() {
	case tcell.KeyUp: if cursorY > 0 { cursorY-- }
	case tcell.KeyDown: if cursorY < len(lines)-1 { cursorY++ }
	case tcell.KeyLeft: if cursorX > 0 { cursorX-- } else if showTree && currentMode == "VIEW" { isTreeFocused = true }
	case tcell.KeyRight: if cursorX < len(lines[cursorY]) { cursorX++ }
	case tcell.KeyEscape: currentMode, isSelectedAll, infoMessage = "VIEW", false, ""
	}

	if currentMode == "VIEW" {
		if ev.Key() == tcell.KeyRune {
			if ev.Rune() == 'i' { currentMode = "INSERT" } else if ev.Rune() == ':' { currentMode, commandBuffer = "COMMAND", "" }
		}
		return
	}

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
				cursorY-- ; cursorX = prevLen
			}
		case tcell.KeyRune:
			char := ev.Rune()
			if strings.ContainsRune("{};()", char) { infoMessage = "INFO: Use Ctrl+B to save!" } else { infoMessage = "" }
			lines[cursorY] = lines[cursorY][:cursorX] + string(char) + lines[cursorY][cursorX:]
			cursorX++
			if char == ' ' { checkSnippet() }
		}
	}
}

func executeCommand() {
	parts := strings.Fields(commandBuffer)
	if len(parts) == 0 { return }
	cmd := parts[0]
	switch cmd {
	case "w": if len(parts) > 1 { filename = parts[1] } ; if filename != "" { saveToFile() }
	case "q": printToTerminal() ; os.Exit(0)
	case "wq": if len(parts) > 1 { filename = parts[1] } ; if filename != "" { saveToFile() ; printToTerminal() ; os.Exit(0) }
	}
	commandBuffer, currentMode = "", "VIEW"
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
