package itchat4go

import (
    "net/http"
    "sync"
)

type chatInfo struct {
    loginInfo map[string]interface{}
    loginBaseRequest map[string]string
    client *http.Client
    wg sync.WaitGroup
}

var chatter *chatInfo

func init() {
    chatter = new(chatInfo)
    chatter.loginInfo = make(map[string]interface{})
    chatter.client = &http.Client{}
}

