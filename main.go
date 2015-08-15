package main

import (
	"fmt"
)

const (
	boardstring = "............ooooobbbb..oooooobbb..ooooooobb..oooooooob..ooooooooo..woooooooo..wwooooooo..wwwoooooo..wwwwooooo............"
	NORT        = -11
	SOUT        = 11
	WEST        = -1
	EAST        = 1
	EMPTY       = 'o'
	OFFBD       = '.'
	BLACK       = 'b'
	WHITE       = 'w'
)

type Move struct {
	from, to int
}

type Board struct {
	field      [121]byte
	moves      map[Move]bool
	moverecord [128]Move
	recordcnt  int
	movecolor  byte
}

func InitBoard() (board Board) {
	board.moves = make(map[Move]bool)
	for i := 0; i != len(boardstring); i++ {
		board.field[i] = boardstring[i]
	}
	board.recordcnt = 0
	board.movecolor = WHITE
	return
}

func ToString(board Board) (out string) {
	for i := 0; i != len(boardstring); i++ {
		if i%11 == 0 && i != 0 {
			out += string('\n')
		}
		out += string(board.field[i])
	}
	return
}

/* Returns all legal steps for a figure */
func GenerateSteps(board Board, from int) {
	//TODO MAKE CONST ARRAY
	directions := [4]int{NORT, SOUT, WEST, EAST}

	for _, dir := range directions {
		to := from + dir
		if board.field[to] == EMPTY {
			move := Move{from, to}
			board.moves[move] = true
		}
	}
}

func Jumps(board Board, from int) (jumps []int) {
	directions := [4]int{NORT, SOUT, WEST, EAST}
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

func GenerateJumps(board Board, from int, curr int) {
	for _, to := range Jumps(board, curr) {
		move := Move{from, to}
		if !board.moves[move] {
			board.moves[move] = true
			GenerateJumps(board, from, to)
		}
	}
}

func GenerateMoves(board Board) {
	for i, f := range board.field {
		if f == board.movecolor {
			GenerateSteps(board, i)
			GenerateJumps(board, i, i)
		}
	}
}

func SwitchColor(color byte) byte {
	if color == WHITE {
		return BLACK
	}
	return WHITE
}

func MakeMove(board Board, move Move) {
	board.field[move.to] = board.field[move.from]
	board.field[move.from] = EMPTY
	board.moverecord[board.recordcnt] = move
	board.recordcnt++
	board.movecolor = SwitchColor(board.movecolor)
}

func UnMakeMove(board Board) {
	move := board.moverecord[board.recordcnt]
	board.recordcnt--
	board.field[move.from] = board.field[move.to]
	board.field[move.to] = EMPTY
	board.movecolor = SwitchColor(board.movecolor)
}

func Abs(n int) (ret int) {
	ret = n
	if n < 0 {
		ret = -n
	}
	return
}

func DistanceToGoal(board Board, sq int) (distance int) {
	x, y := sq%11, sq/11
	goalX, goalY := 0, 0
	if board.movecolor == WHITE {
		goalX, goalY = 9, 1
	} else {
		goalX, goalY = 1, 9
	}
	distX, distY := Abs(x-goalX), Abs(y-goalY)
	distance = Abs(distX + distY)
	return
}

func Evaluate(board Board) (score int) {
	for sq, f := range board.field {
		if f == board.movecolor {
			score -= DistanceToGoal(board, sq)
		}
	}
	return
}

func AlphaBeta(board Board, depth int, alpha int, beta int, bestmove *Move) int {
	if depth == 0 {
		return Evaluate(board)
	}

	board.moves = make(map[Move]bool)
	GenerateMoves(board)
	val := 0

	for move, _ := range board.moves {
		MakeMove(board, move)
		val = -AlphaBeta(board, depth-1, -beta, -alpha, bestmove)
		UnMakeMove(board)

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

func main() {
	board := InitBoard()
	fmt.Println(ToString(board))
	bestmove := Move{}
	AlphaBeta(board, 10, -100000, 100000, &bestmove)
	fmt.Println(bestmove)
}
