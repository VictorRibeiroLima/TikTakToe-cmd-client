package tiktaktoe

import (
	"errors"
	"fmt"
)

type Game struct {
	table [3][3]int8
	turn  int8
	moves int8
}

func (g Game) Draw() {
	for i := 0; i < len(g.table); i++ {
		for j := 0; j < len(g.table); j++ {
			var symbol string
			square := g.table[i][j]
			if square == 0 {
				symbol = " "
			} else if square == 1 {
				symbol = "x"
			} else {
				symbol = "o"
			}
			fmt.Print(symbol)
			if j != 2 {
				fmt.Print(" | ")
			}
		}
		fmt.Println()
		if i != 2 {
			fmt.Println("----------")
		}
	}
}

func (g *Game) MakePlay(row int8, column int8) error {
	if g.turn == 0 {
		g.turn = 1
	}
	if row < 0 || row > 2 {
		return errors.New("Invalid row")
	}
	if column < 0 || column > 2 {
		return errors.New("Invalid column")
	}
	g.moves++
	if g.table[row][column] != 0 {
		return errors.New("square already marked")
	}
	g.table[row][column] = g.turn

	if g.turn == 1 {
		g.turn = 2
	} else {
		g.turn = 1
	}
	return nil
}
