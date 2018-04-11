package itchat4go

import (
    "fmt"
)

func updateLocalChatrooms(l []map[string]interface{}) {
    for _, chatroom := range l {
        emojiFormatter(chatroom, "NickName")
        fmt.Println(chatroom["NickName"])
        if chatroom["MemberList"] != nil && len(chatroom["MemberList"].([]interface{})) > 0 {
            for _, member := range chatroom["MemberList"].([]interface{}) {
                if _, ok := member.(map[string]interface{})["NickName"]; ok {
                    emojiFormatter(member, "NickName")
                }
                if _, ok := member.(map[string]interface{})["DisplayName"]; ok {
                    emojiFormatter(member, "DisplayName")
                }
                if _, ok := member.(map[string]interface{})["RemarkName"]; ok {
                    emojiFormatter(member, "RemarkName")
                }
            }
        }
    }
}
