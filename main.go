package main

import (
    "github.com/syndtr/goleveldb/leveldb"
)

func main() {
    _, err := leveldb.OpenFile("/tmp", nil)
    if err != nil {
        panic(err)
    }
}
