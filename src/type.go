package itchat4go

type userDictWrapper interface {
    selfInit(map[string]interface{})
}

type userClass struct {
    dict map[string]interface{}
}

func (u *userClass) selfInit(d map[string]interface{}) {
    u.dict = d
}

type returnValue struct {
    retDict map[string]interface{}
}

func (r *returnValue) setValue(value map[string]interface{}) {
    if r.retDict == nil {
        r.retDict = make(map[string]interface{})
    }

    for k, v := range value {
        r.retDict[k] = v
    }
    if _, ok := r.retDict["BaseResponse"]; !ok {
        r.retDict["BaseResponse"] = map[string]interface{} {
            "ErrMsg": "no BaseResponse in raw response",
            "Ret": -1000,
        }
    }
}
