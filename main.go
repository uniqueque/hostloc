package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
)

type Config struct {
	TgToken string `yaml:"tg_token"`
	TgId    int64  `yaml:"tg_id"`
	Users   []User `yaml:"users"`
}

type User struct {
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	FastLoginField string `yaml:"fastloginfield"`
	TgId           int64  `yaml:"tg_id"`
}

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"

var (
	conf Config
	bot  *tgbotapi.BotAPI
)

func main() {
	//设置时区
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
	data, err := os.ReadFile("./conf.yaml")
	if err != nil {
		log.Println("读取配置文件失败", err)
		panic(err)
	} else {
		//解析配置文件
		err := yaml.Unmarshal(data, &conf)
		if err != nil {
			log.Println("解析配置文件失败", err)
			panic(err)
		}
	}
	if conf.TgToken != "" {
		bot, err = tgbotapi.NewBotAPI(conf.TgToken)
		if err != nil {
			log.Panic(err)
		}
	}
	Task()
}

func Task() {
	for _, user := range conf.Users {
		getCredit(user)
		//随机休眠
		time.Sleep(time.Duration(rand.Intn(600)+60) * time.Second)
	}
	//随机休眠到第二天8-20点
	s := (24-time.Now().Hour()+8)*3600 + rand.Intn(12*3600)
	log.Println("下次任务将在", time.Now().Add(time.Duration(s)*time.Second).Format("2006-01-02 15:04:05"), "执行")
	time.Sleep(time.Duration(s) * time.Second)
	Task()
}

func getCredit(user User) {
	if user.FastLoginField == "" {
		if strings.Contains(user.Username, "@") {
			user.FastLoginField = "email"
		} else {
			user.FastLoginField = "username"
		}
	}
	// 创建一个CookieJar
	cookieJar, _ := cookiejar.New(nil)

	// 创建一个带有CookieJar的HTTP客户端
	client := &http.Client{
		Jar: cookieJar,
	}
	loginData := url.Values{
		"fastloginfield": {user.FastLoginField},
		"username":       {user.Username},
		"password":       {user.Password},
		"quickforward":   {"yes"},
		"handlekey":      {"ls"},
	}
	loginBody := strings.NewReader(loginData.Encode())
	// 创建POST请求
	loginReq, err := http.NewRequest("POST", "https://hostloc.com/member.php?mod=logging&action=login&loginsubmit=yes&infloat=yes&lssubmit=yes&inajax=1", loginBody)
	if err != nil {
		log.Println("Creating POST request failed:", err)
		return
	}
	// 设置Content-Type头字段
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 设置User-Agent头字段
	loginReq.Header.Set("User-Agent", userAgent)
	// 发送POST请求
	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Println("POST request failed:", err)
		return
	}
	defer loginResp.Body.Close()
	// 检查响应状态码
	if loginResp.StatusCode != http.StatusOK {
		log.Println("POST request failed with status:", loginResp.StatusCode)
		return
	}
	for i, e := 0, 0; i < 20 && e < 5; i++ {
		//随机休眠
		time.Sleep(time.Duration(rand.Intn(5)+3) * time.Second)
		// 使用获取到的cookie发送GET请求
		uid := rand.Intn(60000) + 1
		spaceReq, err := http.NewRequest("GET", "https://hostloc.com/space-uid-"+strconv.Itoa(uid)+".html", nil)
		if err != nil {
			log.Println("Creating GET request failed:", err)
			return
		}
		// 设置User-Agent
		spaceReq.Header.Set("User-Agent", userAgent)
		// 发送GET请求
		spaceResp, err := client.Do(spaceReq)
		if err != nil {
			log.Println("GET request failed:", err)
			return
		}
		defer spaceResp.Body.Close()

		// 处理GET请求的响应
		spaceBody, err := io.ReadAll(spaceResp.Body)
		if err != nil {
			fmt.Println("Reading GET response failed:", err)
			return
		}
		if strings.Contains(string(spaceBody), "您指定的用户空间不存在") {
			e++
		} else {
			i++
		}
	}
	// 休眠
	time.Sleep(time.Second)
	homeReq, err := http.NewRequest("GET", "https://hostloc.com/home.php?mod=spacecp&ac=credit", nil)
	if err != nil {
		fmt.Println("Creating GET request failed:", err)
		return
	}
	// 设置User-Agent
	homeReq.Header.Set("User-Agent", userAgent)
	// 发送GET请求
	homeResp, err := client.Do(homeReq)
	if err != nil {
		fmt.Println("GET request failed:", err)
		return
	}
	defer homeResp.Body.Close()

	// 处理GET请求的响应
	homeBody, err := io.ReadAll(homeResp.Body)
	if err != nil {
		fmt.Println("Reading GET response failed:", err)
		return
	}
	homeHtml := string(homeBody)
	str := "用户: " + user.Username
	group := ""
	groupReg := regexp.MustCompile(`用户组: ([^<]+)`)
	groupMatch := groupReg.FindStringSubmatch(homeHtml)
	if len(groupMatch) > 1 {
		group = groupMatch[1]
	}
	str += "\n用户组: " + group
	credit := ""
	creditReg := regexp.MustCompile(`积分: (\d+)`)
	creditMatch := creditReg.FindStringSubmatch(homeHtml)
	if len(creditMatch) > 1 {
		credit = creditMatch[1]
	}
	str += "\n积分: " + credit
	money := ""
	moneyReg := regexp.MustCompile(`金钱: </em>(\d+)`)
	moneyMatch := moneyReg.FindStringSubmatch(homeHtml)
	if len(moneyMatch) > 1 {
		money = moneyMatch[1]
	}
	str += "\n金钱: " + money
	prestige := ""
	prestigeReg := regexp.MustCompile(`威望: </em>(\d+)`)
	prestigeMatch := prestigeReg.FindStringSubmatch(homeHtml)
	if len(prestigeMatch) > 1 {
		prestige = prestigeMatch[1]
	}
	str += "\n威望: " + prestige
	// 发送消息
	if bot != nil {
		if user.TgId != 0 {
			msg := tgbotapi.NewMessage(user.TgId, str)
			msg.ParseMode = "HTML"
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(conf.TgId, str)
			msg.ParseMode = "HTML"
			bot.Send(msg)
		}
	} else {
		log.Println(str)
	}
}
