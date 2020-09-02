package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Board struct {
	height int
	width  int

	grid [][]Cell
}

type Cell struct {
	mine     bool
	covered  bool
	flagged  bool
	adjacent int
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func NewBoard(width, height int) *Board {
	rand.Seed(time.Now().UnixNano())
	matcher := randomInt(1, 10)
	aGrid := make([][]Cell, height)
	for i := range aGrid {
		aRow := make([]Cell, width)
		aGrid[i] = aRow
		for j, row := range aRow {
			fmt.Printf("i: %d, j: %d.\n", i, j)
			mined := false
			matchee := randomInt(1, 10)
			if matchee == matcher {
				mined = true
				//Increment adjacent field of adjacent cells since we are a mine.
				//Top Right (row-1,col-1)
				if i != 0 && j != 0 {
					aGrid[i-1][j-1].adjacent += 1
				}
				//Above (row-1,col)
				if i != 0 {
					aGrid[i-1][j].adjacent += 1
				}
				//Top Left (row-1,col+1)
				if i != 0 && j != width {
					aGrid[i-1][j+1].adjacent += 1
				}
				//Left (row,col-1)
				if j != 0 {
					aGrid[i][j-1].adjacent += 1
				}
				//Right (row,col+1)
				if j != width {
					aGrid[i][j+1].adjacent += 1
				}
				//Bottom Right (row+1,col-1)
				if i != height && j != 0 {
					aGrid[i+1][j-1].adjacent += 1
				}
				//Below (row+1,co1)
				if i != height {
					aGrid[i+1][j].adjacent += 1
				}
				//Bottom Left (row+1,col+1)
				if i != height && j != width {
					aGrid[i+1][j+1].adjacent += 1
				}

			}
			row.covered = true
			row.flagged = false
			row.mine = mined
		}
	}
	aBoard := Board{
		height: height,
		width:  width,
		grid:   aGrid,
	}

	return &aBoard
}

func main() {
	aBoard := NewBoard(25, 25)
	fmt.Printf("The Board: %v", aBoard)

}
