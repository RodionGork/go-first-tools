package main

import "fmt"

func main() {
    n := 3000000
    var sieve [3000000]byte
    sieve[0] = 1
    sieve[1] = 1
    c := 2
    cnt := 0
    for true {
        for sieve[c] != 0 {
            c++
        }
        sieve[c] = 2
        cnt++
        if cnt >= 199999 {
            break
        }
        for k := c * 2; k < n; k += c {
            sieve[k] = 1
        }
    }
    fmt.Println("hi ", c)
}
