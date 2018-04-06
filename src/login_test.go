package itchat4go

import (
    "fmt"
    "time"
    "runtime"
    "testing"
)

func TestLogin(t *testing.T) {
    runtime.GOMAXPROCS(runtime.NumCPU())

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

//    loginStatus := make(chan struct{}, 1)
//    var filename string = "QRcode.png"
//    chatter.wg.Add(1)
//    go displayQRPic(loginStatus, filename)

    isloggedin := false
    var status string
    for !isloggedin {
        status = checkLogin(uuid)
        time.Sleep(time.Duration(1)*time.Second)
        switch {
        case status == "200":
//            loginStatus <- struct{}{}
            isloggedin = true
        case status == "201":
            fmt.Println("please confirm on your phone")
        case status != "408":
//            close(loginStatus)
            break
        }
    }
    chatter.wg.Wait()

    if isloggedin {
        fmt.Println(status)
    } else {
        fmt.Println("error occurs")
    }
    webInit()
}
