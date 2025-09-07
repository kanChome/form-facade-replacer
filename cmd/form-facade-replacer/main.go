package main

import (
    "os"

    "form-facade-replacer/internal/ffr"
)

func main() {
    os.Exit(ffr.Run(os.Args))
}

