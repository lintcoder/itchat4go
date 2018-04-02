package itchat4go

import (
    "net/http"
)

type chatInfo struct {
    loginInfo map[string]string
    logintime int64
    client *http.Client
}

var chatter *chatInfo

func init() {
    chatter = new(chatInfo)
    chatter.loginInfo = make(map[string]string)
    chatter.client = &http.Client{}
}
