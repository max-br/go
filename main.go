package main

import (
	"fmt"
)

const (
	boardstring = "............ooooobbbb..oooooobbb..ooooooobb..oooooooob..ooooooooo..woooooooo..wwooooooo..wwwoooooo..wwwwooooo............"
	debugstring = "............ooooobbbb..oooooobob..ooooooobb..ooooooooo..ooooobobo..ooowboboo..ooooooooo..ooooooooo..ooooooooo............"
	NORT        = -11
	SOUT        = 11
	WEST        = -1
	EAST        = 1
	EMPTY       = 'o'
	OFFBD       = '.'
	BLACK       = 'b'
	WHITE       = 'w'
)

var directions = [...]int{NORT, SOUT, WEST, EAST}
var board Board

type Move struct {
	from, to int
}

type Board struct {
	field      [121]byte
	moves      map[Move]bool
	moverecord [128]Move
	recordcnt  int
	us         byte
	them       byte
}

func InitBoard() {
	board.moves = make(map[Move]bool)
	for i := 0; i != len(boardstring); i++ {
		board.field[i] = boardstring[i]
	}
	board.recordcnt = 0
	board.us = WHITE
	board.them = BLACK
	return
}

func ToString() (out string) {
	for i := 0; i != len(boardstring); i++ {
		if i%11 == 0 && i != 0 {
			out += string('\n')
		}
		out += string(board.field[i])
	}
	return
}

/* Returns all legal steps for a figure */
func GenerateSteps(from int) {
	for _, dir := range directions {
		to := from + dir
		if board.field[to] == EMPTY {
			move := Move{from, to}
			board.moves[move] = true
		}
	}
}

func Jumps(from int) (jumps []int) {
	for _, dir := range directions {
		onestep := from + dir
		if board.field[onestep] == BLACK || board.field[onestep] == WHITE {
			to := onestep + dir // Jump over figure
			if board.field[to] == EMPTY {
				jumps = append(jumps, to)
			}
		}
	}
	return
}

func GenerateJumps(from int, curr int) {
	for _, to := range Jumps(curr) {
		move := Move{from, to}
		if !board.moves[move] {
			board.moves[move] = true
			GenerateJumps(from, to)
		}
	}
}

func GenerateMoves() {
	for i, f := range board.field {
		if f == board.us {
			GenerateSteps(i)
			GenerateJumps(i, i)
		}
	}
}

func MakeMove(move Move) {
	board.field[move.to] = board.field[move.from]
	board.field[move.from] = EMPTY
	board.moverecord[board.recordcnt] = move
	board.recordcnt++
	board.us, board.them = board.them, board.us
}

func UnMakeMove() {
	board.recordcnt--
	move := board.moverecord[board.recordcnt]
	board.field[move.from] = board.field[move.to]
	board.field[move.to] = EMPTY
	board.us, board.them = board.them, board.us
}

func Abs(n int) (ret int) {
	ret = n
	if n < 0 {
		ret = -n
	}
	return
}

func IndexToCoord(index int) (x int, y int) {
	x = index % 11
	y = index / 11
	return
}

func DistanceToGoal(sq int) (distance int) {
	x, y := IndexToCoord(sq)
	goalX, goalY := 0, 0
	if board.us == WHITE {
		goalX, goalY = 9, 1
	} else {
		goalX, goalY = 1, 9
	}
	distX, distY := Abs(x-goalX), Abs(y-goalY)
	distance = Abs(distX + distY)
	return
}

func Evaluate() (score int) {
	for sq, f := range board.field {
		if f == board.us {
			score -= DistanceToGoal(sq)
		}
		if f == board.them {
			score += DistanceToGoal(sq)
		}
	}
	return
}

var cnt int = 0

func AlphaBeta(depth int, alpha int, beta int, bestmove *Move) int {
	cnt++
	if depth == 0 {
		return Evaluate()
	}

	board.moves = make(map[Move]bool)
	GenerateMoves()
	val := 0

	for move, _ := range board.moves {
		MakeMove(move)
		val = -AlphaBeta(depth-1, -beta, -alpha, bestmove)
		UnMakeMove()
		if val >= beta {
			return beta
		}
		if val > alpha {
			alpha = val
			bestmove = &move
		}
	}
	return alpha
}

func Perft(depth int) (nodes int) {
	if depth == 0 {
		return 1
	}
	board.moves = make(map[Move]bool)
	GenerateMoves()
	for move, _ := range board.moves {
		MakeMove(move)
		nodes += Perft(depth - 1)
		UnMakeMove()
	}
	return nodes
}

func Divide(depth int) {
	board.moves = make(map[Move]bool)
	GenerateMoves()
	fmt.Println(len(board.moves))
	for move, _ := range board.moves {
		MakeMove(move)
		x, y := IndexToCoord(move.from)
		x2, y2 := IndexToCoord(move.to)
		fmt.Println(Perft(depth-1), move, "from: ", x, y, "to: ", x2, y2)
		UnMakeMove()
	}
}

func main() {
	InitBoard()
	fmt.Println(Perft(2))
}
