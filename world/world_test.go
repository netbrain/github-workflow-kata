package world

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewGame(t *testing.T) {
	w := NewWorld(8,4)

	require.EqualValues(t, w.Grid,[][]byte{
		{'.','.','.','.','.','.','.','.'},
		{'.','.','.','.','.','.','.','.'},
		{'.','.','.','.','.','.','.','.'},
		{'.','.','.','.','.','.','.','.'},
	})
}

func TestNewGameFromReader(t *testing.T) {
	input :=
	"Generation 1:\n" +
	"4 8\n" +
	"........\n" +
	"....*...\n" +
	"...**...\n" +
	"........"

	w,err := NewWorldFromReader(bytes.NewBufferString(input))
	require.NoError(t,err)

	require.Equal(t, 1, w.Generation)
	require.Equal(t, 4, len(w.Grid))
	require.Equal(t, 8, len(w.Grid[0]))
	require.EqualValues(t,[][]byte{
		{'.','.','.','.','.','.','.','.'},
		{'.','.','.','.','*','.','.','.'},
		{'.','.','.','*','*','.','.','.'},
		{'.','.','.','.','.','.','.','.'},
	}, w.Grid)
	require.Equal(t,input, w.String())
}

func Test_IsAlive(t *testing.T) {
	tests := []struct{
		name string
		input string
		x int
		y int
		expected bool
	}{
		{"dead",".", 0, 0, false},
		{"alive","*",0, 0,true},
		{"out of bounds -1 0","*", -1, 0, false},
		{"out of bounds 0 -1","*", 0, -1, false},
		{"out of bounds -1 -1","*", -1, -1, false},
		{"out of bounds 1 0","*", 1, 0, false},
		{"out of bounds 0 1","*", 0, 1, false},
		{"out of bounds 1 1","*", 1, 1, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w,err := NewWorldFromReader(bytes.NewBufferString("Generation 1:\n1 1\n"+test.input))
			require.NoError(t,err)

			require.Equal(t,test.expected, w.IsAlive(test.x,test.y))
		})
	}
}

func TestGame_LiveNeighbours(t *testing.T) {
	tests := []struct{
		name string
		input string
		x int
		y int
		expected int
	}{
		{"all dead","...\n...\n...", 1, 1, 0},
		{"one","*..\n...\n...", 1, 1, 1},
		{"two","**.\n...\n...", 1, 1, 2},
		{"three","***\n...\n...", 1, 1, 3},
		{"four","***\n*..\n...", 1, 1, 4},
		{"five","***\n*.*\n...", 1, 1, 5},
		{"six","***\n*.*\n*..", 1, 1, 6},
		{"seven","***\n*.*\n**.", 1, 1, 7},
		{"eight","***\n*.*\n***", 1, 1, 8},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g,err := NewWorldFromReader(bytes.NewBufferString("Generation 1:\n3 3\n"+test.input))
			require.NoError(t,err)

			require.Equal(t,test.expected,g.LiveNeighbours(test.x,test.y))
		})
	}
}

func TestGame_Tick(t *testing.T) {
	tests := []struct{
		input  string
		output string
	}{
		{
			"Generation 1:\n" +
				"4 8\n" +
				"........\n" +
				"....*...\n" +
				"...**...\n" +
				"........",
			"Generation 2:\n" +
				"4 8\n" +
				"........\n" +
				"...**...\n" +
				"...**...\n" +
				"........",
		},
		{
			"Generation 99:\n" +
				"3 3\n" +
				".*.\n" +
				".*.\n" +
				".*.",
			"Generation 100:\n" +
				"3 3\n" +
				"...\n" +
				"***\n" +
				"...",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("tick test #%d",i), func(t *testing.T) {
			w,err := NewWorldFromReader(bytes.NewBufferString(test.input))
			require.NoError(t,err)

			w = w.NextGeneration()
			fmt.Println(w)

			require.Equal(t,test.output, w.String())
		})
	}


}