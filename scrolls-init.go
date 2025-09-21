  GNU nano 7.2           /home/dbrock/Desktop/TheAncientScrolls/scrolls-init.go                     
package main

import (
        "bufio"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "os"
        "os/exec"
        "path/filepath"
        "runtime"
        "sort"
        "strconv"
        "strings"
        "time"
)

type Note struct {
        ID          int       `json:"id"`
        Title       string    `json:"title"`
        Content     string    `json:"content"`
        Tags        []string  `json:"tags"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
                                         [ Read 683 lines ]
^G Help       ^O Write Out  ^W Where Is   ^K Cut        ^T Execute    ^C Location   M-U Undo
^X Exit       ^R Read File  ^\ Replace    ^U Paste      ^J Justify    ^/ Go To Line M-E Redo
