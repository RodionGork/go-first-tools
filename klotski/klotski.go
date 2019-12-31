package main

import "fmt"
import "runtime"
import "sort"
import "strconv"
import "strings"

var nei = [][]int{{1, 0}, {0, -1}, {-1, 0}, {0, 1}}

var board [][]byte
var boardsSeen = make(map[string]string)
var boardsToCheck = make([]string, 0)
var congr map[byte]byte
var width int
var height int
var goalId byte
var goalX, goalY int

func main() {
    readInput()
    congruency()
    addPosition("x")
    for len(boardsToCheck) > 0 {
        b := boardsToCheck[0]
        setPosition(b)
        if checkGoal() {
            traceBack()
            break
        }
        boardsToCheck = boardsToCheck[1:]
        processPosition()
    }
}

func readInput() {
    res := make([][]byte, 0)
    var line string
    for true {
        fmt.Scanf("%s", &line)
        if strings.Contains(line, "-") { break }
        res = append(res, []byte(line))
    }
    board = res
    height = len(board)
    width = len(board[0])
    goal := strings.Split(line, "-")
    goalId = goal[0][0]
    goalX, _ = strconv.Atoi(goal[1])
    goalY, _ = strconv.Atoi(goal[2])
}

func congruency() {
    congr = make(map[byte]byte)
    backMap := make(map[string]byte)
    for _, row := range board {
        for _, v := range row {
            if v == '.' || congr[v] != 0 { continue }
            sq := findSquaresSorted(v)
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

func describePosition() (string, string) {
    exact := ""
    for _, row := range board {
        exact += string(row)
    }
    b := []byte(exact)
    for i, v := range b {
        if v != '.' {
            b[i] = congr[v]
        } else {
            b[i] = '.'
        }
    }
    plain := string(b)
    return exact, plain
}

func addPosition(parent string) (bool) {
    exact, plain := describePosition()
    if boardsSeen[plain] != "" {
        return false
    }
    boardsSeen[plain] = parent + " " + exact
    boardsToCheck = append(boardsToCheck, exact)
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

func processPosition() {
    //fmt.Println("Trying position:")
    //printPosition()
    _, plain := describePosition()
    moves := findMoves()
    for _, m := range moves {
        id, dx, dy := byte(m[0]), m[1], m[2]
        makeMove(id, dx, dy)
        addPosition(plain)
        makeMove(id, -dx, -dy)
    }
    seen := len(boardsSeen)
    if seen < 200000 || seen % 1000 == 0 {
        fmt.Println("SEEN", seen)
        fmt.Println("TOCHECK", len(boardsToCheck))
        var ms runtime.MemStats
        runtime.ReadMemStats(&ms)
        fmt.Println("MEM", ms.Alloc / (1024 * 1024))
    }
}

func findMoves() ([][]int) {
    empties := findSquares('.')
    moveables := make(map[byte]bool)
    for _, e := range empties {
        for _, n := range nei {
            x, y := e[0] + n[0], e[1] + n[1]
            if x >= 0 && y >= 0 && x < width && y < height {
                id := board[y][x]
                if id > '.' {
                    moveables[id] = true
                }
            }
        }
    }
    moves := make([][]int, 0)
    for m := range moveables {
        mList := findMovesFor(m)
        moves = append(moves, mList...)
    }
    return moves
}

func findMovesFor(id byte) ([][]int) {
    sq := findSquares(id)
    lst := [][]int{{int(id), 0, 0}}
    for i := 0; i < len(lst); i++ {
        for _, dir := range nei {
            dx, dy := lst[i][1] + dir[0], lst[i][2] + dir[1]
            seen := false
            for _, m := range lst {
                if m[1] == dx && m[2] == dy {
                    seen = true
                    break
                }
            }
            if seen { continue }
            allowed := true
            for _, s := range sq {
                x, y := s[0] + dx, s[1] + dy
                if x < 0 || x >= width || y < 0 || y >= height {
                    allowed = false
                    break
                }
                if board[y][x] != '.' && board[y][x] != id {
                    allowed = false
                    break
                }
            }
            if allowed {
                lst = append(lst, []int{int(id), dx, dy})
            }
        }
    }
    return lst[1:]
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

func findSquaresSorted(id byte) ([][]int) {
    sq := findSquares(id)
    sort.Slice(sq, func(i, j int) (bool) {
        return sq[i][1] * 256 + sq[i][0] < sq[j][1] * 256 + sq[j][0]
    })
    return sq
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

func makeMove(id byte, dx int, dy int) {
    sq := findSquares(id)
    for _, s := range sq {
        board[s[1]][s[0]] = '.'
    }
    for _, s := range sq {
        board[s[1] + dy][s[0] + dx] = id
    }
}

func checkGoal() (bool) {
    sq := findSquaresSorted(goalId)
    return sq[0][0] == goalX && sq[0][1] == goalY
}

func traceBack() {
    fmt.Println("TRACEBACK!!!")
    res := make([]string, 0)
    for true {
        exact, plain := describePosition()
        res = append(res, exact)
        s := strings.Split(boardsSeen[plain], " ")
        if s[0] == "x" { break }
        s2 := strings.Split(boardsSeen[s[0]], " ")
        setPosition(s2[1])
    }
    for i := range res {
        setPosition(res[len(res) - i - 1])
        printPosition()
        fmt.Println("-", i)
    }
}
