package main

import "fmt"
import "sort"
import "strconv"

var nei = [][]int{{1, 0}, {0, -1}, {-1, 0}, {0, 1}}

var board [][]byte
var boardsSeen = make(map[string]bool)
var boardsToCheck = make([]string, 0)
var congr map[byte]byte
var width int
var height int

func main() {
    readBoard()
    congruency()
    addPosition()
    for len(boardsToCheck) > 0 {
        setPosition(boardsToCheck[0])
        boardsToCheck = boardsToCheck[1:]
        findMoves()
    }
}

func readBoard() {
    res := make([][]byte, 0)
    var line string
    for true {
        _, e := fmt.Scanf("%s", &line)
        if e != nil || line == "-" { break }
        res = append(res, []byte(line))
    }
    board = res
    height = len(board)
    width = len(board[0])
}

func congruency() {
    congr = make(map[byte]byte)
    backMap := make(map[string]byte)
    for _, row := range board {
        for _, v := range row {
            if v == '.' || congr[v] != 0 { continue }
            sq := findSquares(v)
            sort.Slice(sq, func(i, j int) (bool) {
                return sq[i][0] * 256 + sq[i][1] < sq[j][0] * 256 + sq[j][1]
            })
            x0 := sq[0][0]
            y0 := sq[0][1]
            s := ""
            for _, xy := range sq {
                s += strconv.Itoa(xy[0] - x0) + strconv.Itoa(xy[1] - y0)
            }
            if backMap[s] == 0 {
                backMap[s] = v
            }
            congr[v] = backMap[s]
        }
    }
}

func addPosition() (bool) {
    s := ""
    for _, row := range board {
        for _, v := range row {
            if v != '.' {
                s += string(congr[v])
            } else {
                s += "."
            }
        }
    }
    if boardsSeen[s] {
        return false
    }
    boardsSeen[s] = true
    s = ""
    for _, row := range board {
        s += string(row)
    }
    boardsToCheck = append(boardsToCheck, s)
    return true
}

func setPosition(s string) {
    for y, row := range board {
        for x := range row {
            row[x] = s[y * width + x]
        }
    }
}

func printPosition() {
    for _, row := range board {
        fmt.Println(string(row))
    }
}

func findMoves() {
    fmt.Println("Trying position:")
    printPosition()
    empties := findSquares('.')
    for _, e := range empties {
        for dir, n := range nei {
            x := e[0] + n[0]
            y := e[1] + n[1]
            if x >= 0 && y >= 0 && x < width && y < height {
                id := board[y][x]
                dirMove := (dir + 2) % 4
                if id > '.' && movePossible(id, dirMove) {
                    fmt.Println("MOVE:", string(id), dirMove)
                    makeMove(id, dirMove)
                    if addPosition() {
                        printPosition()
                    } else {
                        fmt.Println("skip...")
                    }
                    makeMove(id, dir)
                }
            }
        }
    }
    fmt.Println("SEEN", len(boardsSeen))
    fmt.Println("TOCHECK", len(boardsToCheck))
    ur := ""
    fmt.Scanf("%s", &ur)
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

func makeMove(id byte, dir int) {
    sq := findSquares(id)
    for _, s := range sq {
        board[s[1]][s[0]] = '.'
    }
    dx := nei[dir][0]
    dy := nei[dir][1]
    for _, s := range sq {
        board[s[1] + dy][s[0] + dx] = id
    }
}
