package itchat4go

import (
    "net/http"
    "sync"
)

type storageClass struct {
    userName string
    nickName string
}

type chatInfo struct {
    loginInfo map[string]interface{}
    loginBaseRequest map[string]string
    memberList []interface{}
    chatroomList []map[string]interface{}
    storageClass
    client *http.Client
    wg sync.WaitGroup
}

var chatter *chatInfo

var friendInfoArr1 = []string {"UserName", "City", "DisplayName", "PYQuanPin", "RemarkPYInitial", "Province",
                        "KeyWord", "RemarkName", "PYInitial", "EncryChatRoomId", "Alias", "Signature",
                        "NickName", "RemarkPYQuanPin", "HeadImgUrl"}
var friendInfoArr2 = []string {"UniFriend", "Sex", "AppAccountFlag", "VerifyFlag", "ChatRoomId", "HideInputBarFlag",
                        "AttrStatus", "SnsFlag", "MemberCount", "OwnerUin", "ContactFlag", "Uin",
                        "StarFriend", "Statues"}

func init() {
    chatter = new(chatInfo)
    chatter.loginInfo = make(map[string]interface{})
    chatter.chatroomList = make([]map[string]interface{}, 0)
    chatter.client = &http.Client{}
}

