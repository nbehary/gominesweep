package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
)

type Board struct {
	height int
	width  int

	grid [][]*Cell
}

var infoLog *log.Logger

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
	}
	number := 0
	for i := range aGrid {
		for j := range aGrid[i] {
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
				if i != 0 && j != (width-1) {
					aGrid[i-1][j+1].adjacent += 1
				}
				//Left (row,col-1)
				if j != 0 {
					aGrid[i][j-1].adjacent += 1
				}
				//Right (row,col+1)
				if j != (width - 1) {
					aGrid[i][j+1].adjacent += 1
				}
				//Bottom Right (row+1,col-1)
				if i != (height-1) && j != 0 {
					aGrid[i+1][j-1].adjacent += 1
				}
				//Below (row+1,co1)
				if i != (height - 1) {
					aGrid[i+1][j].adjacent += 1
				}
				//Bottom Left (row+1,col+1)
				if i != (height-1) && j != (width-1) {
					aGrid[i+1][j+1].adjacent += 1
				}

			}
			aGrid[i][j].covered = true
			aGrid[i][j].flagged = false
			aGrid[i][j].mine = mined
			aGrid[i][j].x = i
			aGrid[i][j].y = i
			aGrid[i][j].number = number
			number += 1
		}
	}
	aBoard := Board{
		height: height,
		width:  width,
		grid:   aGrid,
	}

	return &aBoard
}

func drawBoard(s tcell.Screen, b *Board) {
	for i, row := range b.grid {
		for j, cell := range row {
			infoLog.Printf("Cell: %v. \n", cell)
			//Just Draw M for mines, otherwise value of adjacent for now.
			if cell.mine {
				//Draw M
				drawCell(s, i, j, rune('M'))
				infoLog.Printf("Mine at %d,%d.\n", i, j)
			} else {
				//Draw cell.adjacent
				ch := strconv.Itoa(cell.adjacent)
				r := []rune(ch)
				drawCell(s, i, j, r[0])
			}
		}
	}
}

func drawCell(s tcell.Screen, x, y int, ch rune) {
	st := tcell.StyleDefault
	rgb := tcell.NewHexColor(int32(0xb0afac))
	st = st.Background(rgb)
	size := 3
	drawBox(s, x*size, y*size, (x*size)+(size-1), (y*size)+(size-1), st, ch)
}

func whichCell(x, y int) (int, int) {
	outX := int(x / 3)
	outY := int(y / 3)
	return outX, outY
}

func main() {
	infoLogFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	infoLog = log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime)
	if err != nil {
		log.Fatal("Error opening logfile")
	}
	aBoard := NewBoard(10, 50)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.EnableMouse()
	s.Clear()

	drawBoard(s, aBoard)
	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventMouse:
				x, y := ev.Position()
				//button := ev.Buttons()
				switch ev.Buttons() {
				case tcell.Button1:
					x2, y2 := whichCell(x, y)
					infoLog.Printf("Button 1. x:%d, y:%d. Cell x:%d,y%d", x, y, x2, y2)
				case tcell.Button2:
					infoLog.Println("Button 2")
				case tcell.Button3:
					infoLog.Println("Button 3")
				default:
				}
			}
		}
	}()
loop:
	for {
		select {
		case <-quit:
			break loop
		}
	}
	s.Fini()
	os.Exit(0)

}
