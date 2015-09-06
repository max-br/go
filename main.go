package main

import (
	"fmt"
)

const (
	boardstring = "############.....BBBB##......BBB##.......BB##........B##.........##W........##WW.......##WWW......##WWWW.....############"
	debugstring = "############.....BBBB##......B.B##.......BB##.........##.....B.B.##...wB.B..##.........##.........##.........############"
	NORT        = -11
	SOUT        = 11
	WEST        = -1
	EAST        = 1
	EMPTY       = '.'
	OFFBD       = '#'
	BLACK       = 'B'
	WHITE       = 'W'
)

var directions = [...]int{NORT, SOUT, WEST, EAST}

type Move struct {
	from, to int
}

type Board struct {
	field      [121]byte
	moverecord [128]Move
	recordcnt  int
	us         byte
	them       byte
}

func (board *Board) InitBoard() {
	for i := 0; i != len(boardstring); i++ {
		board.field[i] = boardstring[i]
	}
	board.recordcnt = 0
	board.us = WHITE
	board.them = BLACK
	return
}

func (board *Board) ToString() (out string) {
	for i := 0; i != len(boardstring); i++ {
		if i%11 == 0 && i != 0 {
			out += string('\n')
		}
		out += string(board.field[i])
		out += string(' ')
	}
	out += string('\n')
	return
}

/* Returns all legal steps for a figure */
func (board *Board) GenerateSteps(from int, moves map[Move]bool) {
	for _, dir := range directions {
		to := from + dir
		if board.field[to] == EMPTY {
			move := Move{from, to}
			moves[move] = true
		}
	}
}

func (board *Board) Jumps(from int) (jumps []int) {
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

func (board *Board) GenerateJumps(from int, curr int, moves map[Move]bool) {
	for _, to := range board.Jumps(curr) {
		move := Move{from, to}
		if !moves[move] {
			moves[move] = true
			board.GenerateJumps(from, to, moves)
		}
	}
}

func (board *Board) GenerateMoves(moves map[Move]bool) {
	for i, f := range board.field {
		if f == board.us {
			board.GenerateJumps(i, i, moves)
			board.GenerateSteps(i, moves)
		}
	}
}

func (board *Board) MakeMove(move Move) {
	board.field[move.to] = board.field[move.from]
	board.field[move.from] = EMPTY
	board.moverecord[board.recordcnt] = move
	board.recordcnt++
	board.us, board.them = board.them, board.us
}

func (board *Board) UnMakeMove() {
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

func Goal(color byte) (goalX int, goalY int) {
	if color == WHITE {
		goalX, goalY = 9, 1
	} else {
		goalX, goalY = 1, 9
	}
	return
}

func DistanceTo(index int, goalX int, goalY int) (distance int) {
	x, y := IndexToCoord(index)
	distance = Abs(x-goalX) + Abs(y-goalY)
	return
}

func (board *Board) Evaluate() (score int) {
	for index, f := range board.field {
		if f == board.us {
			goalX, goalY := Goal(board.us)
			score -= DistanceTo(index, goalX, goalY)
		}
		if f == board.them {
			goalX, goalY := Goal(board.them)
			score += DistanceTo(index, goalX, goalY)
		}
	}
	return -score
}

var cnt int = 0

func (board *Board) AlphaBeta(depth int, alpha int, beta int) int {
	cnt++
	if depth == 0 {
		return board.Evaluate()
	}
	moves := make(map[Move]bool)
	board.GenerateMoves(moves)
	val := 0

	for move, _ := range moves {
		board.MakeMove(move)
		val = -board.AlphaBeta(depth-1, -beta, -alpha)
		board.UnMakeMove()
		if val >= beta {
			return beta
		}
		if val > alpha {
			alpha = val
		}
	}
	return alpha
}

func (board *Board) SearchBestMove(depth int) (bestmove Move) {
	moves := make(map[Move]bool)
	board.recordcnt = 0
	board.GenerateMoves(moves)
	bestscore := -50000
	for move, _ := range moves {
		board.MakeMove(move)
		score := board.AlphaBeta(depth-1, -10000, 10000)
		if score > bestscore {
			bestmove = move
			bestscore = score
			fmt.Println(move, score)
		}
		board.UnMakeMove()
	}
	return
}

func (board *Board) Perft(depth int) (nodes int) {
	if depth == 0 {
		return 1
	}
	moves := make(map[Move]bool)
	board.GenerateMoves(moves)
	for move, _ := range moves {
		board.MakeMove(move)
		nodes += board.Perft(depth - 1)
		board.UnMakeMove()
	}
	return nodes
}

func (board *Board) Divide(depth int) {
	moves := make(map[Move]bool)
	board.GenerateMoves(moves)
	fmt.Println(len(moves))
	for move, _ := range moves {
		board.MakeMove(move)
		x, y := IndexToCoord(move.from)
		x2, y2 := IndexToCoord(move.to)
		fmt.Println(board.Perft(depth-1), move, "from: ", x, y, "to: ", x2, y2)
		board.UnMakeMove()
	}
}

func main() {
	var board Board
	board.InitBoard()
	//fmt.Println(board.Perft(5)) // should be 1381888
	for {
		bm := board.SearchBestMove(5)
		board.MakeMove(bm)
		fmt.Println(board.ToString())
	}
}
