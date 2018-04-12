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
        oldChatroom := searchDictList(chatter.chatroomList, "UserName", chatroom["UserName"].(string))
        if len(oldChatroom) > 0 {
            updateInfoDict(oldChatroom, chatroom)
            var memberList []interface{}
            if m, ok := chatroom["MemberList"]; ok {
                memberList = append(memberList, m)
            }
            oldMemberList := oldChatroom["MemberList"].([]map[string]interface{})
            if len(memberList) > 0 {
                for _, member := range memberList {
                    oldMember := searchDictList(oldMemberList, "UserName", member.(map[string]interface{})["UserName"].(string))
                    if oldMember != nil {
                        updateInfoDict(oldMember, member.(map[string]interface{}))
                    } else {
                        oldMemberList = append(oldMemberList, member.(map[string]interface{}))
                    }
                }
            }
        } else {
            chatter.chatroomList = append(chatter.chatroomList, chatroom)
            oldChatroom = searchDictList(chatter.chatroomList, "UserName", chatroom["UserName"].(string))
        }
        fmt.Println(oldChatroom)

//        if len(chatroom["MemberList"]) != len(oldChatroom["MemberList"]) && chatroom["MemberList"] != nil {
    }
}
