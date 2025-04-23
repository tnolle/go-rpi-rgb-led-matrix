package rgbmatrix

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type BDFGlyph struct {
	Encoding    rune
	Width       int
	Height      int
	XOffset     int
	YOffset     int
	DeviceWidth int
	Bitmap      []string
}

type BDFFont struct {
	Glyphs map[rune]*BDFGlyph
	Ascent int
	width  int
	height int
}

func (f *BDFFont) Height() int {
	if f.height != 0 {
		return f.height
	}
	for _, glyph := range f.Glyphs {
		if glyph.Height > f.height {
			f.height = glyph.Height
		}
	}
	return f.height
}

func (f *BDFFont) Width() int {
	if f.width != 0 {
		return f.width
	}
	for _, glyph := range f.Glyphs {
		if glyph.Width > f.width {
			f.width = glyph.Width
		}
	}
	return f.width
}

func LoadBDF(path string) (*BDFFont, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	font := &BDFFont{Glyphs: make(map[rune]*BDFGlyph)}
	var current *BDFGlyph
	inBitmap := false

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "FONT_ASCENT":
			font.Ascent, _ = strconv.Atoi(fields[1])
		case "STARTCHAR":
			current = &BDFGlyph{}
		case "ENCODING":
			code, _ := strconv.Atoi(fields[1])
			current.Encoding = rune(code)
		case "DWIDTH":
			current.DeviceWidth, _ = strconv.Atoi(fields[1])
		case "BBX":
			current.Width, _ = strconv.Atoi(fields[1])
			current.Height, _ = strconv.Atoi(fields[2])
			current.XOffset, _ = strconv.Atoi(fields[3])
			current.YOffset, _ = strconv.Atoi(fields[4])
		case "BITMAP":
			inBitmap = true
		case "ENDCHAR":
			font.Glyphs[current.Encoding] = current
			current = nil
			inBitmap = false
		default:
			if inBitmap && current != nil {
				current.Bitmap = append(current.Bitmap, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return font, nil
}
