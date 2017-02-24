package main

import (
    "io"
    "os"
    "strings"
)

type rot13Reader struct {
    r io.Reader
}

func (r rot13Reader) Read(b []byte) (i int, e error) {
    i, e = r.r.Read(b)
    if e != nil {
        return 0, e
    }
    for x := range b {
        if 'A' <= b[x] && b[x] <= 'Z' {
            b[x] += 13
            if b[x] > 'Z' {
                b[x] -= 26
            }
        } else if 'a' <= b[x] && b[x] <= 'z' {
            b[x] += 13
            if b[x] > 'z' {
                b[x] -= 26
            }
        }
    }
    return
}

func main() {
    s := strings.NewReader("Lbh penpxrq gur pbqr!")
    r := rot13Reader{s}
    io.Copy(os.Stdout, &r)
}