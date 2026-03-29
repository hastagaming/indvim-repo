package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var (
	lines       = []string{""}
	cursorX     = 0
	cursorY     = 0
	filename    = "nasa_project.txt"
	showWarning = false

	// Snippets otomatis (Ketik lalu Spasi)
	snippets = map[string]string{
		"!kt":   "fun main() {\n    println(\"Hello Nasa!\")\n}",
		"!java": "public class Main {\n    public static void main(String[] args) {\n        \n    }\n}",
		"!html": "<!DOCTYPE html>\n<html>\n<head><title>INDVIM</title></head>\n<body>\n    \n</body>\n</html>",
		"!py":   "def main():\n    print('INDVIM is the best!')\n\nif __name__ == '__main__':\n    main()",
	}

	// Syntax Highlighting dengan warna yang didukung tcell
	keywords = map[string]tcell.Color{
		"func": tcell.ColorOrange, "package": tcell.ColorDeepPink, "import": tcell.ColorYellow,
		"var": tcell.ColorTeal, "if": tcell.ColorPurple, "else": tcell.ColorPurple,
		"return": tcell.ColorRed, "println": tcell.ColorBlue, "val": tcell.ColorTeal,
		"fun": tcell.ColorOrange, "def": tcell.ColorOrange, "class": tcell.ColorYellow,
		"print": tcell.ColorBlue,
	}
)

func main() {
	s, _ := tcell.NewScreen()
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
			// Tekan ESC -> Munculkan Peringatan
			if ev.Key() == tcell.KeyEscape {
				showWarning = true
				continue
			}
			// Tekan CTRL + B -> Save & Exit
			if ev.Key() == tcell.KeyCtrlB {
				saveToFile()
				return
			}
			
			// Jika menekan tombol lain, hilangkan peringatan
			showWarning = false 
			handleInput(ev)
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func draw(s tcell.Screen) {
	s.Clear()
	w, h := s.Size()

	// Style Definitions (NvChad Vibe)
	styleLineNum := tcell.StyleDefault.Foreground(tcell.ColorDimGray)
	styleText := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	// Gambar Editor & Syntax Highlighting
	for y, line := range lines {
		if y >= h-2 { break }
		
		// Line Numbers (Nomor Baris)
		num := fmt.Sprintf(" %2d │ ", y+1)
		for i, char := range num {
			s.SetContent(i, y, char, nil, styleLineNum)
		}

		// Syntax Highlighting Engine Sederhana
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

	// Peringatan CTRL+B (Hanya muncul jika ESC ditekan)
	if showWarning {
		msg1 := " Are you want to leave without saving? "
		msg2 := " type : ctrl + b for exit and save the file "
		warnStyle := tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite).Bold(true)
		
		// Taruh di tengah layar
		startX := (w - len(msg2)) / 2
		if startX < 0 { startX = 0 }
		startY := h / 2
		
		for i, char := range msg1 { s.SetContent(startX+i+((len(msg2)-len(msg1))/2), startY, char, nil, warnStyle) }
		for i, char := range msg2 { s.SetContent(startX+i, startY+1, char, nil, warnStyle) }
	}

	// Gambar Status Bar (Di bagian bawah)
	status := fmt.Sprintf(" [INDVIM] | FILE: %s | SAVE & EXIT: CTRL + B ", filename)
	barStyle := tcell.StyleDefault.Background(tcell.ColorMediumSpringGreen).Foreground(tcell.ColorBlack).Bold(true)
	for i, char := range status { s.SetContent(i, h-1, char, nil, barStyle) }
	for i := len(status); i < w; i++ { s.SetContent(i, h-1, ' ', nil, barStyle) }

	s.ShowCursor(cursorX+6, cursorY)
	s.Show()
}

func handleInput(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyUp:
		if cursorY > 0 { cursorY-- }
	case tcell.KeyDown:
		if cursorY < len(lines)-1 { cursorY++ }
	case tcell.KeyLeft:
		if cursorX > 0 { cursorX-- }
	case tcell.KeyRight:
		if cursorX < len(lines[cursorY]) { cursorX++ }
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
	
	// Cegah kursor lompat keluar batas baris
	if cursorX > len(lines[cursorY]) { cursorX = len(lines[cursorY]) }
}

func checkSnippet() {
	// Cek kata terakhir sebelum kursor (tidak termasuk spasi yang baru diketik)
	if cursorX < 2 { return }
	words := strings.Fields(lines[cursorY][:cursorX-1])
	if len(words) == 0 { return }
	lastWord := words[len(words)-1]

	if val, ok := snippets[lastWord]; ok {
		// Hapus trigger snippet dan spasi
		startPos := cursorX - len(lastWord) - 1
		lines[cursorY] = lines[cursorY][:startPos] + lines[cursorY][cursorX:]
		cursorX = startPos

		// Masukkan isi snippet multi-line
		snipLines := strings.Split(val, "\n")
		lines[cursorY] = lines[cursorY][:cursorX] + snipLines[0]
		for i := 1; i < len(snipLines); i++ {
			cursorY++
			lines = append(lines[:cursorY], append([]string{snipLines[i]}, lines[cursorY:]...)...)
		}
		cursorX = len(lines[cursorY])
	}
}

func saveToFile() {
	content := strings.Join(lines, "\n")
	os.WriteFile(filename, []byte(content), 0644)
}
