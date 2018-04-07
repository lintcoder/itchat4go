package itchat4go

import (
    "net/http"
    "sync"
)

type chatInfo struct {
    loginInfo map[string]string
    loginBaseRequest map[string]string
    loginTime int64
    inviteStartCount int
    client *http.Client
    wg sync.WaitGroup
}

var chatter *chatInfo

func init() {
    chatter = new(chatInfo)
    chatter.loginInfo = make(map[string]string)
    chatter.client = &http.Client{}
}

