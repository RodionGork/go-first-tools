package main

import "fmt"
import "io"
import "os"
import "strings"
import "encoding/base64"

func main() {
    defer failureCheck()
    room, user, text := parseInput()
    fmt.Printf("Content-Type: text/plain; charset=utf-8\r\n\r\n")
    fmt.Printf("ROOM=%v, USER=%v, DATA='%v'\n", room, user, text)
}

func failureCheck() {
    if r := recover(); r != nil {
        fmt.Printf("Status: 400 Bad Request\r\n\r\n")
        fmt.Println("Problem: " + r.(string))
        os.Exit(0)
    }
}

func parseInput() (string, string, string) {
    
    data := make([]byte, 16000)
    n, errRead := io.ReadFull(os.Stdin, data)
    if errRead != nil && errRead != io.ErrUnexpectedEOF {
        panic("Error on reading input: " + errRead.Error())
    }
    
    chunks := strings.SplitN(string(data[:n]), " ", 3)
    if len(chunks) < 3 {
        panic("Malformed input")
    }
    
    room, user := chunks[0], chunks[1]
    if chunks[2] == "===" {
        return room, user, "" // means retrieve messages, rather than send
    }
    
    text, errUnpack := base64.StdEncoding.DecodeString(chunks[2])
    if errUnpack != nil {
        panic("Payload unpacking failed")
    }
    
    return room, user, string(text)
}
