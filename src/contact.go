package itchat4go

import (
    "fmt"
    "strconv"
)

func updateLocalChatrooms(l []map[string]interface{}) map[string]interface{} {
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

        if chatroom["MemberList"] != nil && len(chatroom["MemberList"].([]interface{})) != len(oldChatroom["MemberList"].([]interface{})) {
            var existUserNames map[string]struct{}
            var delList []int
            for _, member := range chatroom["MemberList"].([]interface{}) {
                existUserNames[member.(map[string]interface{})["UserName"].(string)] = struct{}{}
            }
            for i, member := range oldChatroom["MemberList"].([]interface{}) {
                if _, ok := existUserNames[member.(map[string]interface{})["UserName"].(string)]; ok {
                    delList = append(delList, i)
                }
            }
            var res []interface{}
            index := 0
            for i := 0; i < len(delList); i++ {
                res = append(res, oldChatroom["MemberList"].([]interface{})[index:delList[i]]...)
                index = delList[i]+1
            }
            oldChatroom["MemberList"] = append(res, oldChatroom["MemberList"].([]interface{})[index:]...)
        }

        if _, ok := oldChatroom["ChatRoomOwner"]; ok {
            if _, ok = oldChatroom["MemberList"]; ok {
                owner := searchDictList(oldChatroom["MemberList"].([]map[string]interface{}), "UserName", oldChatroom["ChatRoomOwner"].(string))
                if owner == nil {
                    oldChatroom["OwnerUin"] = 0
                } else if uin, ok := owner["Uin"]; ok {
                    oldChatroom["OwnerUin"] = uin
                } else {
                    oldChatroom["OwnerUin"] = 0
                }
            }
        }

        if uin, ok := oldChatroom["OwnerUin"]; ok && int(uin.(float64)) != 0 {
            wxuin, _ := strconv.Atoi(chatter.loginInfo["wxuin"].(string))
            oldChatroom["IsAdmin"] = int(uin.(float64)) == wxuin
        } else {
            oldChatroom["IsAdmin"] = nil
        }

        if ml := oldChatroom["MemberList"]; ml != nil {
            newSelf := searchDictList(ml.([]map[string]interface{}), "UserName", chatter.userName)
            if newSelf != nil {
                oldChatroom["Self"] = newSelf
            } else {
                oldChatroom["Self"] = deepCopy(chatter.loginInfo["User"])
            }
        } else {
            oldChatroom["Self"] = deepCopy(chatter.loginInfo["User"])
        }
    }

    res := make(map[string]interface{})
    res["Type"] = "System"
    res["SystemInfo"] = "chatrooms"
    res["FromUserName"] = chatter.userName
    res["ToUserName"] = chatter.userName
    res["Text"] = make([]interface{}, len(l))
    for i := 0; i < len(l); i++ {
        res["Text"].([]interface{})[i] = l[i]["UserName"]
    }
    return res
}

func updateLocalFriends(l []map[string]interface{}) {
    fullList := append(chatter.memberList, chatter.mpList...)
    for _, friend := range l {
        if _, ok := friend["NickName"]; ok {
            emojiFormatter(friend, "NickName")
        }
        if _, ok := friend["DisplayName"]; ok {
            emojiFormatter(friend, "DisplayName")
        }
        if _, ok := friend["RemarkName"]; ok {
            emojiFormatter(friend, "RemarkName")
        }
        oldInfoDict := searchDictList(fullList, "UserName", friend["UserName"].(string))
        if oldInfoDict == nil {
            oldInfoDict = deepCopy(friend).(map[string]interface{})
            if int(oldInfoDict["VerifyFlag"].(float64)) & 8 == 0 {
                chatter.memberList = append(chatter.memberList, oldInfoDict)
            } else {
                chatter.mpList = append(chatter.mpList, oldInfoDict)
            }
        } else {
            updateInfoDict(oldInfoDict, friend)
        }
    }
}
