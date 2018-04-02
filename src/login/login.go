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

    pattern := "window.QRLogin.code = \\d+; window.QRLogin.uuid = \"\\S+?\";"
    reg, _ := regexp.Compile(pattern)
    if matchstr := reg.FindSubmatch(data); matchstr == nil {
        return uuid, errors.New("No match pattern in response body")
    } else {
        var ct int
        for i, ch := range matchstr[0] {
            if ch == '=' && ct == 0 && string(matchstr[0][i+2:i+5]) == "200" {
                ct++
            } else if ch == '=' && ct == 1 {
                j := i + 3
                for ; j < len(matchstr[0]); j++ {
                    if matchstr[0][j] == '"' {
                        break
                    }
                }
                uuid = string(matchstr[0][i+3:j])
                break
            }
        }
    }
    if uuid == "" {
        err = errors.New("get QR uuid failed")
    }

    return uuid, err
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
    pattern := "window.redirect_uri=\"\\S+\";"
    reg, _ := regexp.Compile(pattern)
    loginInfoURL := reg.Find(loginContent)
    chatter.loginInfo["url"] = string(loginInfoURL[bytes.IndexByte(loginInfoURL, '"')+1 :
    bytes.LastIndex(loginInfoURL, []byte{'/'})])

    fmt.Println(chatter.loginInfo["url"])

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
    chatter.logintime = time.Now().Unix() * 1e3
    chatter.loginInfo["BaseRequest"] = ""
    fmt.Println(chatter.loginInfo)
    fmt.Println(chatter.logintime)

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

    pattern := "window.code=\\d+"
    reg, _ := regexp.Compile(pattern)
    if matchstr := reg.FindSubmatch(data); matchstr != nil {
        for i, ch := range matchstr[0] {
           if ch == '=' {
               status := string(matchstr[0][i+1:i+4])
               if status != "200" {
                   return status
               } else {
                   if processLoginInfo(data) {
                       return "200"
                   } else {
                       return "400"
                   }
               }
           }
       }
   }

   return "400"
}
