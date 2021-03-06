package itchat4go

import (
	"fmt"
	"regexp"
	"strings"
)

func emojiFormatter(dict interface{}, key string) {
	s := []byte(strings.Replace(dict.(map[string]interface{})[key].(string), `<span class="emoji emoji1f450"></span`, `<span class="emoji emoji1f450"></span>`, -1))
	reg, _ := regexp.Compile(`<span class="emoji emoji(\S+)"></span>`)

	repdict := map[string]string{"1f63c": "1f601", "1f639": "1f602", "1f63a": "1f603",
		"1f4ab": "1f616", "1f64d": "1f614", "1f63b": "1f60d",
		"1f63d": "1f618", "1f64e": "1f621", "1f63f": "1f622"}
	replfunc := func(matched []byte) []byte {
		index := reg.FindSubmatchIndex(matched)
		res := []byte{}
		if v, ok := repdict[string(matched[index[2]:index[3]])]; ok {
			res = append(res, matched[:index[2]]...)
			res = append(res, v...)
			res = append(res, matched[index[3]:]...)
		} else {
			res = append(res, matched...)
		}
		return res
	}

	snew := reg.ReplaceAllFunc(s, replfunc)

	formatfunc := func(matched []byte) []byte {
		s := reg.FindSubmatch(matched)[1]
		arr := [8]byte{'0', '0', '0', '0', '0', '0', '0', '0'}
		res1 := []byte{}
		res2 := []byte{}
		sz := len(s)
		var fmtstr string
		if sz == 6 {
			res1 = append(res1, arr[:6]...)
			res1 = append(res1, s[:2]...)
			res2 = append(res2, arr[:4]...)
			res2 = append(res2, s[2:]...)
			fmtstr = fmt.Sprintf("\\U%s\\U%s", string(res1), string(res2))
		} else if sz == 10 {
			res1 = append(res1, arr[:3]...)
			res1 = append(res1, s[:5]...)
			res2 = append(res2, arr[:3]...)
			res2 = append(res2, s[5:]...)
			fmtstr = fmt.Sprintf("\\U%s\\U%s", string(res1), string(res2))
		} else {
			if sz < 8 {
				res1 = append(res1, arr[:8-sz]...)
			}
			res1 = append(res1, s...)
			fmtstr = fmt.Sprintf("\\U%s", string(res1))
		}
		return []byte(fmtstr)
	}
	dict.(map[string]interface{})[key] = string(reg.ReplaceAllFunc(snew, formatfunc))
}

func wrapUserDict(d map[string]interface{}) userDictWrapper {
	var r userDictWrapper
	vflag := d["VerifyFlag"].(int)
	if vflag&8 == 0 {
		r = new(userClass)
		r.selfInit(d)
	}
	return r
}

func searchDictList(l []map[string]interface{}, key string, value string) map[string]interface{} {
	for _, i := range l {
		if v, ok := i[key]; ok && v.(string) == value {
			return i
		}
	}
	return nil
}

func updateInfoDict(oldInfoDict map[string]interface{}, newInfoDict map[string]interface{}) {
	for k, v := range newInfoDict {
		switch v.(type) {
		case []interface{}:
			return
		case map[string]interface{}:
			return
		default:
			if _, ok := oldInfoDict[k]; !ok || (v != nil && v != "" && v != '0' && v != 0) {
				oldInfoDict[k] = v
			}
		}
	}
}

func deepCopy(src interface{}) interface{} {
	if value, ok := src.(map[string]interface{}); ok {
		dst := make(map[string]interface{})
		for k, v := range value {
			dst[k] = deepCopy(v)
		}
		return dst
	} else if value, ok := src.([]interface{}); ok {
		dst := make([]interface{}, len(value))
		for k, v := range value {
			dst[k] = deepCopy(v)
		}
		return dst
	}
	return src
}
