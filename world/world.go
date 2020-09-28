package world

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

var ErrFormat = fmt.Errorf("incorrect input format")

type World struct {
	Generation int
	Grid       [][]byte
}

func NewWorld(width,height int) *World {
	g := &World{}
	for y := 0; y < height; y++ {
		g.Grid = append(g.Grid,[]byte{})
		for x := 0; x < width; x++{
			g.Grid[y] = append(g.Grid[y],'.')
		}
	}
	return g
}

func NewWorldFromReader(reader io.Reader)  (*World, error) {
	g := &World{}
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil,err
	}
	buf = bytes.Trim(buf,"\r\n\t ")
	lines := bytes.Split(buf,[]byte{'\n'})
	if !strings.HasPrefix(strings.ToLower(string(lines[0])),"generation") {
		return nil, ErrFormat
	}
	generation := lines[0][len("generation "):len(lines[0])-1]
	wh := bytes.SplitN(lines[1],[]byte{' '},3)
	if len(wh) != 2 {
		return nil, ErrFormat
	}

	height,err := strconv.Atoi(string(wh[0]))
	if err != nil {
		return nil, ErrFormat
	}

	width,err := strconv.Atoi(string(wh[1]))
	if err != nil {
		return nil, ErrFormat
	}

	g.Generation,err = strconv.Atoi(string(generation))
	if err != nil {
		return nil, ErrFormat
	}

	g.Grid = lines[2:]

	if len(g.Grid) != height {
		return nil, ErrFormat
	}

	for _, r := range g.Grid {
		if len(r) != width {
			return nil, ErrFormat
		}
	}

	return g,nil
}

func (w *World) String() string {
	return fmt.Sprintf("Generation %d:\n%d %d\n%s",
		w.Generation,
		w.Height(),
		w.Width(),
		string(bytes.Join(w.Grid,[]byte{'\n'})),
	)
}

func (w *World) Height() int {
	return len(w.Grid)
}

func (w *World) Width() int {
	if w.Height() == 0 {
		return 0
	}
	return len(w.Grid[0])
}

func (w *World) IsAlive(x,y int) bool {
	height := len(w.Grid)
	if height == 0 {
		return false
	}
	width := len(w.Grid[0])

	if y >= height || y < 0 || x >= width || x < 0 {
		return false
	}

	return w.Grid[y][x] == '*'
}

func (w *World) LiveNeighbours(x,y int) int {
	i := 0
	for _,xy := range [][]int{
		{x-1,y-1},
		{x,y-1},
		{x+1,y-1},
		{x-1,y},
		{x+1,y},
		{x-1,y+1},
		{x,y+1},
		{x+1,y+1},
	} {
		if w.IsAlive(xy[0],xy[1]){
			i++
		}
	}
	return i
}

func (w *World) Perish(x,y int) {
	w.Grid[y][x] = '.'
}

func (w *World) Resuscitate(x,y int) {
	w.Grid[y][x] = '*'
}

func (w *World) NextGeneration() *World {
	newGame := NewWorld(w.Width(), w.Height())
	newGame.Generation = w.Generation+1
	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++{
			liveNeigh := w.LiveNeighbours(x,y)
			alive := w.IsAlive(x,y)
			if alive && (liveNeigh < 2 || liveNeigh > 3) {
				newGame.Perish(x,y)
			}else if alive && (liveNeigh == 2 || liveNeigh == 3){
				newGame.Resuscitate(x,y)
			}else if !alive && liveNeigh == 3 {
				newGame.Resuscitate(x,y)
			}
		}
	}
	return newGame
}