package itchat4go

import (
    "fmt"
    "time"
    "testing"
)

func TestLogin(t *testing.T) {
    err := testConnect()
    if err != nil {
        panic(err)
    }

    uuid, err := getQRuuid()
    if err != nil {
        panic(err)
    }

    if !getQRPic(uuid) {
        return
    }

    isloggedin := false
    var status string
    for !isloggedin {
        status = checkLogin(uuid)
        time.Sleep(time.Duration(2)*time.Second)
        switch {
        case status == "200":
            isloggedin = true
        case status == "201":
            fmt.Println("please confirm on your phone")
        case status != "408":
            break
        }
    }

    if isloggedin {
        fmt.Println(status)
    } else {
        fmt.Println("error occurs")
    }
}
