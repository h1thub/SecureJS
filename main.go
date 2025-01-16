package main

import (
    "fmt"
    "os"
    "runtime/debug"

    "SecureJS/cmd"
)

func main() {
    debug.SetTraceback("none")
    
    // 2) 在最顶层拦截 panic
    defer func() {
        if r := recover(); r != nil {
            fmt.Fprintf(os.Stderr, "[!] An error occurred: %v\n", r)
            os.Exit(1)
        }
    }()

    cmd.Execute()
}