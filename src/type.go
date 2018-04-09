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


