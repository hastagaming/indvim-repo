/*
 * INDVIM - Indonesia Text Editor
 * Created by: Nasa (hastagaming) - Dsn. Bangi, Kediri
 * Versi: Global Explorer Edition
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
	lineOffset    = 0
	filename      = "" 
	
	currentMode   = "VIEW" 
	commandBuffer = ""     
	infoMessage   = "" 

	isSelectedAll = false

	// File Tree Explorer (Global)
	showTree      = false
	isTreeFocused = false
	currDir       = "" // Akan diisi Absolute Path di main()
	treeNodes     = []FileNode{}
	treeCursor    = 0
	treeWidth     = 25

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
	// Ambil direktori kerja saat ini secara absolut
	wd, _ := os.Getwd()
	currDir = wd

	if len(os.Args) > 1 {
		argPath, _ := filepath.Abs(os.Args[1])
		info, err := os.Stat(argPath)
		if err == nil && info.IsDir() {
			currDir = argPath
			showTree, isTreeFocused = true, true
		} else {
			filename = argPath
			loadFromFile()
		}
	}

	s, err := tcell.NewScreen()
	if err != nil { fmt.Fprintf(os.Stderr, "%v\n", err) ; os.Exit(1) }
	if err := s.Init(); err != nil { fmt.Fprintf(os.Stderr, "%v\n", err) ; os.Exit(1) }
	defer s.Fini()

	loadDir(currDir)

	for {
		draw(s)
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Modifiers()&tcell.ModAlt != 0 && (ev.Rune() == 'i' || ev.Rune() == 'I') {
				if currentMode == "VIEW" { currentMode = "INSERT" } else if currentMode == "INSERT" { currentMode, commandBuffer = "COMMAND", "" } else { currentMode = "VIEW" }
				continue 
			}
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
	absPath, _ := filepath.Abs(path)
	currDir = absPath
	entries, err := os.ReadDir(currDir)
	treeNodes = []FileNode{}
	
	// Tambahkan opsi ".." untuk naik level kecuali di ROOT
	if currDir != "/" && currDir != filepath.VolumeName(currDir)+"\\" {
		treeNodes = append(treeNodes, FileNode{Name: "..", IsDir: true})
	}

	if err != nil { 
		infoMessage = "Error: Access Denied"
		return 
	}

	for _, e := range entries {
		treeNodes = append(treeNodes, FileNode{Name: e.Name(), IsDir: e.IsDir()})
	}
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

	if cursorY < lineOffset { lineOffset = cursorY }
	if cursorY >= lineOffset+mainH { lineOffset = cursorY - mainH + 1 }

	offsetX := 0
	if showTree {
		offsetX = treeWidth
		// Header Folder di Tree
		folderName := filepath.Base(currDir)
		if folderName == "." || folderName == "/" { folderName = "ROOT" }
		header := " 📂 " + folderName
		for x, char := range header {
			if x < treeWidth-1 { s.SetContent(x, 0, char, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)) }
		}

		for y := 1; y < mainH; y++ {
			s.SetContent(offsetX-1, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorDimGray))
			nodeIdx := y - 1
			if nodeIdx >= len(treeNodes) { continue }
			node := treeNodes[nodeIdx]
			
			nodeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
			if node.IsDir { nodeStyle = tcell.StyleDefault.Foreground(tcell.ColorDodgerBlue).Bold(true) }
			if isTreeFocused && nodeIdx == treeCursor { nodeStyle = nodeStyle.Background(tcell.ColorDimGray) }
			
			prefix := "📁 "
			if !node.IsDir { prefix = "📄 " }
			name := node.Name
			row := prefix + name
			for x, char := range row {
				if x < treeWidth-2 { s.SetContent(x, y, char, nil, nodeStyle) }
			}
		}
	}

	// Render Teks
	for y := 0; y < mainH; y++ {
		lineIdx := y + lineOffset
		if lineIdx >= len(lines) { break }
		line := lines[lineIdx]
		num := fmt.Sprintf(" %2d │ ", lineIdx+1)
		for i, char := range num { s.SetContent(offsetX+i, y, char, nil, tcell.StyleDefault.Foreground(tcell.ColorDimGray)) }

		currX, wordStart := offsetX+6, 0
		for i := 0; i <= len(line); i++ {
			if i == len(line) || line[i] == ' ' {
				word := line[wordStart:i]
				color, isKw := keywords[word]
				wordStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
				if isSelectedAll { wordStyle = wordStyle.Background(tcell.ColorDimGray) } else if isKw { wordStyle = tcell.StyleDefault.Foreground(color).Bold(true) }
				for j := wordStart; j < i; j++ {
					s.SetContent(currX, y, rune(line[j]), nil, wordStyle)
					currX++
				}
				if i < len(line) {
					s.SetContent(currX, y, ' ', nil, wordStyle)
					currX++
				}
				wordStart = i + 1
			}
		}
	}

	// Status Bar
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

	fn := filename
	if fn == "" { fn = "[No Name]" } else { fn = filepath.Base(fn) }
	_, bg, _ := barStyle.Decompose()
	for _, char := range fn { s.SetContent(currPos, h-1, char, nil, tcell.StyleDefault.Background(bg).Foreground(tcell.ColorRed).Bold(true)) ; currPos++ }

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
			newPath := filepath.Join(currDir, node.Name)
			if node.Name == ".." { newPath = filepath.Dir(currDir) }
			
			info, _ := os.Stat(newPath)
			if info != nil && info.IsDir() {
				loadDir(newPath)
			} else {
				filename = newPath
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
	case "w": if len(parts) > 1 { filename, _ = filepath.Abs(parts[1]) } ; if filename != "" { saveToFile() }
	case "q": printToTerminal() ; os.Exit(0)
	case "wq": if len(parts) > 1 { filename, _ = filepath.Abs(parts[1]) } ; if filename != "" { saveToFile() ; printToTerminal() ; os.Exit(0) }
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
