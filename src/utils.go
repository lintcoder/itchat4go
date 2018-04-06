package itchat4go

import (
    "regexp"
    "strings"
    "fmt"
)

func emojiFormatter(dict interface{}, key string) {
    s := []byte(strings.Replace(dict.(map[string]interface{})[key].(string), `<span class="emoji emoji1f450"></span`, `<span class="emoji emoji1f450"></span>`, -1))
    reg, _ := regexp.Compile(`<span class="emoji emoji(\S+)"></span>`)

    repdict := map[string]string {"1f63c": "1f601", "1f639": "1f602", "1f63a": "1f603",
                                  "1f4ab": "1f616", "1f64d": "1f614", "1f63b": "1f60d",
                                  "1f63d": "1f618", "1f64e": "1f621", "1f63f": "1f622"}
    replfunc := func(matched []byte) []byte {
        index := reg.FindSubmatchIndex(matched)
        res := []byte{}
        if v, ok := repdict[string(matched[index[2] : index[3]])]; ok {
            res = append(res, matched[:index[2]]...)
            res = append(res, v...)
            res = append(res, matched[index[3]:]...)
        } else {
            res = append(res, matched...)
        }
        return res
    }

    snew := reg.ReplaceAllFunc(s, replfunc)
    fmt.Println(string(snew))
}
