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
	"strings"

	"github.com/gdamore/tcell/v2"
)

var (
	lines         = []string{""}
	cursorX       = 0
	cursorY       = 0
	filename      = "new_file.txt"
	
	// Variabel Mode Baru
	currentMode   = "VIEW" // Default saat aplikasi dibuka
	commandBuffer = ""     // Menyimpan ketikan saat di mode command

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
		loadFromFile()
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

	for {
		draw(s)
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// FITUR SAKTI: Tekan Alt + i untuk siklus ganti mode
			if ev.Modifiers()&tcell.ModAlt != 0 && (ev.Rune() == 'i' || ev.Rune() == 'I') {
				if currentMode == "VIEW" {
					currentMode = "INSERT"
				} else if currentMode == "INSERT" {
					currentMode = "COMMAND"
					commandBuffer = "" // Bersihkan command sebelumnya
				} else {
					currentMode = "VIEW"
				}
				continue
			}

			// Standar Vim: Tekan ESC untuk kembali ke VIEW mode dengan cepat
			if ev.Key() == tcell.KeyEscape {
				currentMode = "VIEW"
				continue
			}

			// Tombol darurat
			if ev.Key() == tcell.KeyCtrlB {
				saveToFile()
				return
			}

			handleInput(ev)
			
			// Auto-save HANYA terjadi saat kita di Insert mode
			if currentMode == "INSERT" {
				saveToFile()
			}

		case *tcell.EventResize:
			s.Sync()
		}
	}
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
}

func saveToFile() {
	content := strings.Join(lines, "\n")
	os.WriteFile(filename, []byte(content), 0644)
}

func draw(s tcell.Screen) {
	s.Clear()
	w, h := s.Size()

	styleLineNum := tcell.StyleDefault.Foreground(tcell.ColorDimGray)
	styleText := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	// Logo ASCII
	if len(lines) == 1 && lines[0] == "" && currentMode != "COMMAND" {
		startX := (w - len(logo[0])) / 2
		startY := (h - len(logo)) / 2
		logoStyle := tcell.StyleDefault.Foreground(tcell.ColorMediumSpringGreen).Bold(true)
		for y, line := range logo {
			for x, char := range line {
				s.SetContent(startX+x, startY+y, char, nil, logoStyle)
			}
		}
	}

	for y, line := range lines {
		if y >= h-2 { break }
		
		num := fmt.Sprintf(" %2d │ ", y+1)
		for i, char := range num {
			s.SetContent(i, y, char, nil, styleLineNum)
		}

		wordStart := 0
		currX := 6
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

	// 🎨 WARNA STATUS BAR BERDASARKAN MODE
	var status string
	var barStyle tcell.Style

	if currentMode == "COMMAND" {
		status = fmt.Sprintf(" :%s ", commandBuffer)
		barStyle = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack).Bold(true)
	} else if currentMode == "INSERT" {
		status = fmt.Sprintf(" [ INSERT ] | FILE: %s | MODE: Alt+i ", filename)
		barStyle = tcell.StyleDefault.Background(tcell.ColorMediumSpringGreen).Foreground(tcell.ColorBlack).Bold(true)
	} else {
		status = fmt.Sprintf(" [ VIEW ] | FILE: %s | MODE: Alt+i ", filename)
		barStyle = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(true)
	}

	for i, char := range status { s.SetContent(i, h-1, char, nil, barStyle) }
	for i := len(status); i < w; i++ { s.SetContent(i, h-1, ' ', nil, barStyle) }

	// Pindah Kursor ke bawah jika di mode Command
	if currentMode == "COMMAND" {
		s.ShowCursor(len(status), h-1)
	} else {
		s.ShowCursor(cursorX+6, cursorY)
	}
	s.Show()
}

func handleInput(ev *tcell.EventKey) {
	// LOGIKA MODE COMMAND
	if currentMode == "COMMAND" {
		if ev.Key() == tcell.KeyEnter {
			executeCommand()
		} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
			if len(commandBuffer) > 0 {
				commandBuffer = commandBuffer[:len(commandBuffer)-1]
			} else {
				currentMode = "VIEW" // Keluar command jika kosong
			}
		} else if ev.Key() == tcell.KeyRune {
			commandBuffer += string(ev.Rune())
		}
		return
	}

	// Navigasi aktif di VIEW dan INSERT
	switch ev.Key() {
	case tcell.KeyUp: if cursorY > 0 { cursorY-- }
	case tcell.KeyDown: if cursorY < len(lines)-1 { cursorY++ }
	case tcell.KeyLeft: if cursorX > 0 { cursorX-- }
	case tcell.KeyRight: if cursorX < len(lines[cursorY]) { cursorX++ }
	}

	// Jika mode VIEW, blokir input ketikan
	if currentMode == "VIEW" {
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
	cmd := strings.TrimSpace(commandBuffer)
	switch cmd {
	case "w":
		saveToFile() // Save file
	case "q":
		os.Exit(0)   // Keluar
	case "wq":
		saveToFile()
		os.Exit(0)   // Save & Keluar
	}
	commandBuffer = ""
	currentMode = "VIEW" // Setelah enter, otomatis kembali ke View
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
