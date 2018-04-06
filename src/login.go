package itchat4go

import (
    "net/http"
    "fmt"
    "io/ioutil"
    "regexp"
    "errors"
    "time"
    "bytes"
    "strings"
    "math/rand"
    "strconv"
    "encoding/json"
    qrcode "github.com/skip2/go-qrcode"
)

const (
    BASE_URL string = "https://login.weixin.qq.com"
    USER_AGENT string = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0"
)

func testConnect() error {
    var retry_time int = 5
    var login_url string = BASE_URL
    var err error

    for i := 0; i < retry_time; i++ {
        _, err := http.NewRequest("GET", login_url, nil)
        if err == nil {
            break
        }
    }
    return err
}

func getQRuuid() (string, error) {
    var uuid string
    url := BASE_URL + "/jslogin?appid=wx782c26e4c19acffb&fun=new"
    req, err := http.NewRequest("GET", url, nil)

    if err != nil {
        return uuid, err
    }

    req.Header.Add("User-Agent", USER_AGENT)
    response, err := chatter.client.Do(req)

    if err != nil {
        return uuid, err
    }

    data, _ := ioutil.ReadAll(response.Body)
    response.Body.Close()

    pattern := "window.QRLogin.code = (\\d+); window.QRLogin.uuid = \"(\\S+?)\";"
    reg, _ := regexp.Compile(pattern)
    if matchindex := reg.FindSubmatchIndex(data); matchindex == nil {
        return uuid, errors.New("No match pattern in response body")
    } else if string(data[matchindex[2]:matchindex[3]]) != "200" {
        return uuid, errors.New("Response not OK")
    } else {
        uuid = string(data[matchindex[4]:matchindex[5]])
        return uuid, nil
    }
}

func getQRPic(uuid string) bool {
    err := qrcode.WriteFile("https://login.weixin.qq.com/l/" + uuid, qrcode.Medium, 512, "qrpic.png")
    if err != nil {
        fmt.Println("write QR pic error")
        return false
    }
    fmt.Println("Generate QR pic qrpic.png")
    return true
}

func processLoginInfo(loginContent []byte) bool {
    pattern := "window.redirect_uri=\"(\\S+)\";"
    reg, _ := regexp.Compile(pattern)
    loginInfoURL := reg.FindSubmatch(loginContent)
    chatter.loginInfo["url"] = string(loginInfoURL[1][: bytes.LastIndex(loginInfoURL[1], []byte{'/'})])

    indexUrlGrp := [5]string {"wx2.qq.com", "wx8.qq.com", "qq.com", "web2.wechat.com", "wechat.com"}
    detailedUrlGrp := [5][2]string {{"file.wx2.qq.com", "webpush.wx2.qq.com"},
                                    {"file.wx8.qq.com", "webpush.wx8.qq.com"},
                                    {"file.wx.qq.com", "webpush.wx.qq.com"},
                                    {"file.web2.wechat.com", "webpush.web2.wechat.com"},
                                    {"file.web.wechat.com", "webpush.web.wechat.com"}}

    flag := false
    for i := 0; i < 5; i++ {
        if strings.Contains(chatter.loginInfo["url"], indexUrlGrp[i]) {
            chatter.loginInfo["fileUrl"] = fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin", detailedUrlGrp[i][0])
            chatter.loginInfo["syncUrl"] = fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin", detailedUrlGrp[i][1])
            flag = true
            break
        }
    }
    if !flag {
        chatter.loginInfo["fileUrl"] = chatter.loginInfo["url"]
        chatter.loginInfo["syncUrl"] = chatter.loginInfo["url"]
    }

    rand.Seed(time.Now().UnixNano())
    chatter.loginInfo["deviceid"] = "e" + strconv.FormatFloat(rand.Float64(), 'f', 6, 64)[2:] + strconv.FormatFloat(rand.Float64(), 'f', 6, 64)[2:] + strconv.FormatFloat(rand.Float64(), 'f', 3, 64)[2:]
    chatter.loginTime = time.Now().Unix() * 1e3
    chatter.loginBaseRequest = make(map[string]string)

    chatter.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }

    url := string(loginInfoURL[1])
    req, _ := http.NewRequest("GET", url, nil)

    req.Header.Add("User-Agent", USER_AGENT)
    resp, _ := chatter.client.Do(req)
    data, _ := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    chatter.client.CheckRedirect = nil
    return getChildNodes(data)
}

func getChildNodes(xmltext []byte) bool {
    text := string(xmltext)
    targetnodes := [8]string {"skey", "wxsid", "wxuin", "pass_ticket",
                               "</skey>", "</wxsid>", "</wxuin>", "</pass_ticket>"}
    baserequest := [4]string {"Skey", "Sid", "Uin", "DeviceID"}

    for i := 0; i < 4; i++ {
        if begin, end := strings.Index(text, targetnodes[i]), strings.Index(text, targetnodes[i+4]); begin != -1 && end != -1 {
            chatter.loginInfo[targetnodes[i]] = text[begin+len(targetnodes[i])+1 : end]
            chatter.loginBaseRequest[baserequest[i]] = text[begin+len(targetnodes[i])+1 : end]
        } else {
           fmt.Println("Your wechat account may be LIMITED to log in WEB wechat, error info:")
           fmt.Println(text)
           return false
        }
    }
    return true
}

func checkLogin(uuid string) string {
    localtime := time.Now()
    totalsecs := localtime.Unix()
    url := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=%s&tip=1&r=%d&_=%d", BASE_URL, uuid,
    -totalsecs/1579, totalsecs)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("User-Agent", USER_AGENT)
    response, _ := chatter.client.Do(req)
    data, _ := ioutil.ReadAll(response.Body)

    pattern := "window.code=(\\d+)"
    reg, _ := regexp.Compile(pattern)
    if matchstr := reg.FindSubmatchIndex(data); matchstr != nil {
        if status := string(data[matchstr[2]:matchstr[3]]); status == "200" {
            if processLoginInfo(data) {
                return "200"
            } else {
                return "400"
            }
        } else {
            return status
        }
    } else {
        return "400"
    }
}

func webInit() {
    localtime := time.Now()
    totalsecs := localtime.Unix()
    url := fmt.Sprintf("%s/webwxinit?r=%d&pass_ticket=%s", chatter.loginInfo["url"], -totalsecs/1579, chatter.loginInfo["pass_ticket"])

    oridata := map[string] map[string]string {
        "BaseRequest": chatter.loginBaseRequest,
    }
    reqdata, _ := json.Marshal(oridata)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqdata))
    req.Header.Add("Content-Type", "application/json;charset=UTF-8")
    req.Header.Add("User-Agent", USER_AGENT)
    resp, _ := chatter.client.Do(req)
    fmt.Println(resp.StatusCode)

    respdata, _ := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    var dict map[string]interface{}
    json.Unmarshal(respdata, &dict)
    emojiFormatter(dict["User"], "NickName")
}
