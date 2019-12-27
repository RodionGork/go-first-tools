package main

import (
    "fmt"
    "os"
    "bufio"
    "unicode"
)

var keywords = []string {
    "break", "case", "chan", "const", "continue", "default",
    "defer", "else", "fallthrough", "for" , "func", "go", "if",
    "goto", "import", "interface", "map", "package", "range",
    "return", "select", "struct", "switch", "type", "var",
}

var names = []string {
    "true", "false", "iota", "nil", "int", "int8", "int16", "int32", "int64",
    "uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
    "float32", "float64", "complex128", "complex64",
    "bool", "byte", "rune", "string", "error",
}

var keywordsMap map[string]bool = strSetFromArray(keywords)
var namesMap map[string]bool = strSetFromArray(names)

var advance bool = true
var state string = "space"
var cur rune
var prev rune
var keep string

func main() {
    var e error
    in := bufio.NewReader(os.Stdin)
    for {
        if advance {
            prev = cur
            cur, _, e = in.ReadRune()
            if e != nil {
                break
            }
        } else {
            advance = true;
        }
        switch state {
            case "space":
                doSpace(cur)
            case "mb-cmn":
                doMaybeComment(cur)
            case "cmn-line", "cmn-blk":
                doComment(cur)
            case "num":
                doNumber(cur)
            case "word":
                doWord(cur)
            case "str":
                doString(cur)
            default:
                panic("Unknown state: " + state)
        }
    }
}

func doSpace(c rune) {
    if c == '/' {
        changeState("mb-cmn")
    } else if unicode.IsDigit(c) {
        changeState("num")
        fmt.Print(`<i class="num">`)
        advance = false
    } else if unicode.IsLetter(c) {
        changeState("word")
        keep = string(c)
    } else if c == '"' || c == '`' || c == '\'' {
        changeState("str")
        keep = string(c)
        cls := "str"
        if c == '\'' { cls = "chr" }
        fmt.Print(`<i class="`, cls,`">`, keep)
    } else {
        putc(c)
    }
}

func doMaybeComment(c rune) {
    if c == '/' {
        changeState("cmn-line")
        fmt.Print(`<i class="cmn">//`)
    } else if c == '*' {
        changeState("cmn-blk")
        fmt.Print(`<i class="cmn">/*`)
    } else {
        fmt.Print("/")
        doSpace(c)
    }
}

func doComment(c rune) {
    if state == "cmn-line" && (c == '\n' || c == '\r') {
        changeState("space")
    }
    putc(c)
    if state == "cmn-blk" && c == '/' && prev == '*' {
        changeState("space")
    } 
}

func doNumber(c rune) {
    if unicode.IsDigit(c) {
        putc(c)
    } else {
        changeState("space")
        advance = false
    }
}

func doWord(c rune) {
    if unicode.IsLetter(c) {
        keep += string(c)
    } else {
        cls := ""
        if keywordsMap[keep] { cls = "kw" }
        if namesMap[keep] { cls = "name" }
        if cls != "" {
            fmt.Print(`<i class="`, cls, `">`)
        } else {
            state = "space"
        }
        fmt.Print(keep)
        changeState("space")
        advance = false
    }
}

func doString(c rune) {
    putc(c)
    if c == '\\' && prev == '\\' {
        cur = 0
    }
    if c == []rune(keep)[0] && prev != '\\' {
        changeState("space")
    }
}

func changeState(newState string) {
    if newState == "space" && state != "space" {
        fmt.Print("</i>")
    }
    state = newState
}

func strSetFromArray(a []string) map[string]bool {
    m := make(map[string]bool)
    for _, e := range(a) {
        m[e] = true
    }
    return m
}

func putc(c rune) {
    fmt.Print(htmlChar(c))
}

func htmlChar(c rune) string {
    switch c {
        case '<':
            return "&lt;"
        case '&':
            return "&amp;"
        default:
            return string(c)
    }
}
