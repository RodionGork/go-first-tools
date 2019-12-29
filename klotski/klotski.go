package main

import "fmt"

var nei = [][]int{{1, 0}, {0, -1}, {-1, 0}, {0, 1}}

var board [][]byte
var width int
var height int

func main() {
    readBoard()
    findMoves()
}

func readBoard() {
    res := make([][]byte, 0)
    var line string
    for true {
        _, e := fmt.Scanf("%s", &line)
        if e != nil { break }
        res = append(res, []byte(line))
    }
    board = res
    height = len(board)
    width = len(board[0])
}

func findMoves() {
    empties := findSquares('.')
    for _, e := range empties {
        for dir, n := range nei {
            x := e[0] + n[0]
            y := e[1] + n[1]
            if x >= 0 && y >= 0 && x < width && y < height {
                id := board[y][x]
                dirMove := (dir + 2) % 4
                if id > '.' && movePossible(id, dirMove) {
                    fmt.Println("MOVE:", string(rune(id)), x, y, dirMove)
                }
            }
        }
    }
}

func findSquares(id byte) ([][]int) {
    res := make([][]int, 0)
    for y, row := range board {
        for x, v := range row {
            if v == id {
                res = append(res, []int{x, y})
            }
        }
    }
    return res
}

func movePossible(id byte, dir int) (bool) {
    dx := nei[dir][0]
    dy := nei[dir][1]
    sq := findSquares(id)
    for _, s := range sq {
        x := s[0] + dx
        y := s[1] + dy
        if x < 0 || y < 0 || x >= width || y >= height {
            return false
        }
        v := board[y][x]
        if v != '.' && v != id {
            return false
        }
    }
    return true
}
