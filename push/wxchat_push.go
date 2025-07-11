package push

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// WeChatBot 微信机器人客户端
type WeChatBot struct {
	BaseURL    string
	HTTPClient *http.Client
	Cookies    []*http.Cookie
	SyncKey    string
	SKey       string
	Sid        string
	Uin        string
	PassTicket string
	DeviceID   string
	UserName   string
	NickName   string
	IsLogin    bool
	Groups     map[string]string // 群名 -> 群ID
	Contacts   map[string]string // 联系人名 -> 联系人ID
}

// LoginResponse 登录响应
type LoginResponse struct {
	Code        int    `json:"code"`
	RedirectURL string `json:"redirect_url"`
	UUID        string `json:"uuid"`
}

// InitResponse 初始化响应
type InitResponse struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	User struct {
		Uin      int64  `json:"Uin"`
		UserName string `json:"UserName"`
		NickName string `json:"NickName"`
	} `json:"User"`
	SyncKey struct {
		Count int `json:"Count"`
		List  []struct {
			Key int `json:"Key"`
			Val int `json:"Val"`
		} `json:"List"`
	} `json:"SyncKey"`
}

// ContactResponse 联系人响应
type ContactResponse struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	MemberList []struct {
		UserName    string `json:"UserName"`
		NickName    string `json:"NickName"`
		RemarkName  string `json:"RemarkName"`
		MemberCount int    `json:"MemberCount"`
		ContactFlag int    `json:"ContactFlag"`
		VerifyFlag  int    `json:"VerifyFlag"`
		HeadImgUrl  string `json:"HeadImgUrl"`
		MemberList  []struct {
			UserName    string `json:"UserName"`
			NickName    string `json:"NickName"`
			DisplayName string `json:"DisplayName"`
		} `json:"MemberList"`
	} `json:"MemberList"`
}

// SendMsgResponse 发送消息响应
type SendMsgResponse struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	MsgID   string `json:"MsgID"`
	LocalID string `json:"LocalID"`
}

// NewWeChatBot 创建新的微信机器人
func NewWeChatBot() *WeChatBot {
	return &WeChatBot{
		BaseURL:    "https://wx.qq.com",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		DeviceID:   "e" + strconv.FormatInt(rand.Int63n(1000000000000000), 10),
		Groups:     make(map[string]string),
		Contacts:   make(map[string]string),
	}
}

// GetQRCode 获取登录二维码
func (bot *WeChatBot) GetQRCode() (string, error) {
	// 1. 获取UUID
	uuid, err := bot.getUUID()
	if err != nil {
		return "", fmt.Errorf("获取UUID失败: %v", err)
	}

	// 2. 生成二维码URL
	qrURL := fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s", uuid)

	fmt.Printf("请用微信扫描二维码登录:\n")
	fmt.Printf("二维码URL: %s\n", qrURL)
	fmt.Printf("或直接访问: https://login.weixin.qq.com/l/%s\n", uuid)

	return uuid, nil
}

// getUUID 获取登录UUID
func (bot *WeChatBot) getUUID() (string, error) {
	apiURL := fmt.Sprintf("https://login.weixin.qq.com/jslogin?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=%d",
		time.Now().UnixNano()/1000000)

	resp, err := bot.HTTPClient.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析UUID
	re := regexp.MustCompile(`window\.QRLogin\.code = 200; window\.QRLogin\.uuid = "([^"]+)";`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("无法获取UUID")
	}

	return matches[1], nil
}

// WaitForLogin 等待登录
func (bot *WeChatBot) WaitForLogin(uuid string) error {
	for {
		status, redirectURL, err := bot.checkLogin(uuid)
		if err != nil {
			return err
		}

		switch status {
		case 200:
			fmt.Println("✓ 登录成功!")
			return bot.processLogin(redirectURL)
		case 201:
			fmt.Println("✓ 已扫描，请在手机上确认登录...")
		case 408:
			return fmt.Errorf("登录超时")
		default:
			fmt.Print(".")
		}

		time.Sleep(2 * time.Second)
	}
}

// checkLogin 检查登录状态
func (bot *WeChatBot) checkLogin(uuid string) (int, string, error) {
	apiURL := fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=%s&tip=0&r=%d&_=%d",
		uuid, time.Now().UnixNano()/1000000, time.Now().UnixNano()/1000000)

	resp, err := bot.HTTPClient.Get(apiURL)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	content := string(body)

	// 解析状态码
	re := regexp.MustCompile(`window\.code=(\d+);`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return 0, "", fmt.Errorf("无法解析登录状态")
	}

	code, _ := strconv.Atoi(matches[1])

	// 如果登录成功，获取跳转URL
	if code == 200 {
		re = regexp.MustCompile(`window\.redirect_uri="([^"]+)";`)
		matches = re.FindStringSubmatch(content)
		if len(matches) < 2 {
			return 0, "", fmt.Errorf("无法获取跳转URL")
		}
		return code, matches[1], nil
	}

	return code, "", nil
}

// processLogin 处理登录
func (bot *WeChatBot) processLogin(redirectURL string) error {
	// 1. 获取登录信息
	resp, err := bot.HTTPClient.Get(redirectURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 保存cookies
	bot.Cookies = resp.Cookies()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 2. 解析登录参数
	if err := bot.parseLoginInfo(string(body)); err != nil {
		return err
	}

	// 3. 更新BaseURL
	u, _ := url.Parse(redirectURL)
	bot.BaseURL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	// 4. 初始化
	if err := bot.webwxinit(); err != nil {
		return err
	}

	// 5. 获取联系人
	if err := bot.getContacts(); err != nil {
		return err
	}

	bot.IsLogin = true
	fmt.Printf("✓ 登录完成! 昵称: %s\n", bot.NickName)
	return nil
}

// parseLoginInfo 解析登录信息
func (bot *WeChatBot) parseLoginInfo(content string) error {
	// 解析各种参数
	patterns := map[string]*string{
		`<skey>([^<]+)</skey>`:               &bot.SKey,
		`<wxsid>([^<]+)</wxsid>`:             &bot.Sid,
		`<wxuin>([^<]+)</wxuin>`:             &bot.Uin,
		`<pass_ticket>([^<]+)</pass_ticket>`: &bot.PassTicket,
	}

	for pattern, target := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) >= 2 {
			*target = matches[1]
		}
	}

	if bot.SKey == "" || bot.Sid == "" || bot.Uin == "" || bot.PassTicket == "" {
		return fmt.Errorf("无法获取登录参数")
	}

	return nil
}

// webwxinit 初始化微信
func (bot *WeChatBot) webwxinit() error {
	apiURL := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxinit?r=%d&pass_ticket=%s",
		bot.BaseURL, time.Now().UnixNano()/1000000, bot.PassTicket)

	initData := map[string]interface{}{
		"BaseRequest": map[string]string{
			"Uin":      bot.Uin,
			"Sid":      bot.Sid,
			"Skey":     bot.SKey,
			"DeviceID": bot.DeviceID,
		},
	}

	jsonData, _ := json.Marshal(initData)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	bot.addCookies(req)

	resp, err := bot.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result InitResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return fmt.Errorf("初始化失败: %s", result.BaseResponse.ErrMsg)
	}

	bot.UserName = result.User.UserName
	bot.NickName = result.User.NickName

	// 构建SyncKey
	var syncKeys []string
	for _, item := range result.SyncKey.List {
		syncKeys = append(syncKeys, fmt.Sprintf("%d_%d", item.Key, item.Val))
	}
	bot.SyncKey = strings.Join(syncKeys, "|")

	return nil
}

// getContacts 获取联系人
func (bot *WeChatBot) getContacts() error {
	apiURL := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxgetcontact?r=%d&skey=%s",
		bot.BaseURL, time.Now().UnixNano()/1000000, bot.SKey)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	bot.addCookies(req)

	resp, err := bot.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result ContactResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return fmt.Errorf("获取联系人失败: %s", result.BaseResponse.ErrMsg)
	}

	// 分类联系人
	groupCount := 0
	contactCount := 0

	for _, member := range result.MemberList {
		if strings.HasPrefix(member.UserName, "@@") {
			// 群聊
			bot.Groups[member.NickName] = member.UserName
			groupCount++
		} else if member.VerifyFlag == 0 && member.ContactFlag != 0 {
			// 普通联系人
			name := member.RemarkName
			if name == "" {
				name = member.NickName
			}
			bot.Contacts[name] = member.UserName
			contactCount++
		}
	}

	fmt.Printf("✓ 获取联系人完成: %d个群聊, %d个联系人\n", groupCount, contactCount)
	return nil
}

// addCookies 添加cookies
func (bot *WeChatBot) addCookies(req *http.Request) {
	for _, cookie := range bot.Cookies {
		req.AddCookie(cookie)
	}
}

// SendMessage 发送消息
func (bot *WeChatBot) SendMessage(toUserName, content string) error {
	if !bot.IsLogin {
		return fmt.Errorf("未登录")
	}

	apiURL := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxsendmsg?pass_ticket=%s",
		bot.BaseURL, bot.PassTicket)

	clientMsgId := fmt.Sprintf("%d%03d", time.Now().UnixNano()/1000000, rand.Intn(1000))

	msgData := map[string]interface{}{
		"BaseRequest": map[string]string{
			"Uin":      bot.Uin,
			"Sid":      bot.Sid,
			"Skey":     bot.SKey,
			"DeviceID": bot.DeviceID,
		},
		"Msg": map[string]interface{}{
			"Type":         1,
			"Content":      content,
			"FromUserName": bot.UserName,
			"ToUserName":   toUserName,
			"LocalID":      clientMsgId,
			"ClientMsgId":  clientMsgId,
		},
	}

	jsonData, _ := json.Marshal(msgData)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	bot.addCookies(req)

	resp, err := bot.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result SendMsgResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return fmt.Errorf("发送消息失败: %s", result.BaseResponse.ErrMsg)
	}

	return nil
}

// SendToGroup 发送消息到群聊
func (bot *WeChatBot) SendToGroup(groupName, message string) error {
	groupID, exists := bot.Groups[groupName]
	if !exists {
		return fmt.Errorf("未找到群聊: %s", groupName)
	}

	return bot.SendMessage(groupID, message)
}

// SendToGroups 发送消息到多个群聊
func (bot *WeChatBot) SendToGroups(groupNames []string, message string) error {
	var errors []string

	for i, groupName := range groupNames {
		if err := bot.SendToGroup(groupName, message); err != nil {
			errors = append(errors, fmt.Sprintf("群聊[%s]: %v", groupName, err))
		} else {
			fmt.Printf("✓ 已发送到群聊: %s\n", groupName)
		}

		// 发送间隔，避免太快
		if i < len(groupNames)-1 {
			time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分发送失败: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ListGroups 列出群聊
func (bot *WeChatBot) ListGroups() []string {
	var groups []string
	for name := range bot.Groups {
		groups = append(groups, name)
	}
	return groups
}

// ======================= 主程序 =======================

func main() {
	fmt.Println("=== 微信个人号群聊推送脚本 ===")
	fmt.Println("注意: 本脚本仅供学习使用，请勿用于商业用途")
	fmt.Println()

	// 创建机器人实例
	bot := NewWeChatBot()

	// 1. 获取二维码并登录
	fmt.Println("步骤1: 获取登录二维码")
	uuid, err := bot.GetQRCode()
	if err != nil {
		log.Fatal("获取二维码失败:", err)
	}

	fmt.Println("\n步骤2: 等待扫码登录")
	if err := bot.WaitForLogin(uuid); err != nil {
		log.Fatal("登录失败:", err)
	}

	// 2. 显示群聊列表
	fmt.Println("\n步骤3: 获取群聊列表")
	groups := bot.ListGroups()
	if len(groups) == 0 {
		fmt.Println("未找到任何群聊")
		return
	}

	fmt.Printf("找到 %d 个群聊:\n", len(groups))
	for i, group := range groups {
		fmt.Printf("%d. %s\n", i+1, group)
	}

	// 3. 交互式发送消息
	fmt.Println("\n步骤4: 开始发送消息")
	fmt.Println("使用说明:")
	fmt.Println("- 输入消息内容直接发送到所有群聊")
	fmt.Println("- 输入 @群名 消息内容 发送到指定群聊")
	fmt.Println("- 输入 'list' 查看群聊列表")
	fmt.Println("- 输入 'quit' 退出程序")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("请输入消息: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "quit" {
			fmt.Println("程序退出")
			break
		}

		if input == "list" {
			fmt.Println("群聊列表:")
			for i, group := range groups {
				fmt.Printf("%d. %s\n", i+1, group)
			}
			continue
		}

		// 解析输入
		if strings.HasPrefix(input, "@") {
			// 发送到指定群聊
			parts := strings.SplitN(input, " ", 2)
			if len(parts) != 2 {
				fmt.Println("格式错误，应为: @群名 消息内容")
				continue
			}

			groupName := parts[0][1:] // 去掉@符号
			message := parts[1]

			if err := bot.SendToGroup(groupName, message); err != nil {
				fmt.Printf("❌ 发送失败: %v\n", err)
			} else {
				fmt.Printf("✓ 已发送到群聊: %s\n", groupName)
			}
		} else {
			// 发送到所有群聊
			fmt.Printf("正在发送到 %d 个群聊...\n", len(groups))
			if err := bot.SendToGroups(groups, input); err != nil {
				fmt.Printf("❌ 部分发送失败: %v\n", err)
			} else {
				fmt.Printf("✓ 已发送到所有群聊\n")
			}
		}

		fmt.Println()
	}
}

// ======================= 工具函数 =======================

// 自动推送任务示例
func autoSendTask(bot *WeChatBot, message string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			groups := bot.ListGroups()
			if len(groups) > 0 {
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				fullMessage := fmt.Sprintf("[%s] %s", timestamp, message)

				if err := bot.SendToGroups(groups, fullMessage); err != nil {
					log.Printf("自动推送失败: %v", err)
				} else {
					log.Printf("自动推送成功: %s", fullMessage)
				}
			}
		}
	}
}

// 定时任务示例
func scheduleTask(bot *WeChatBot, message string, scheduleTime time.Time) {
	duration := time.Until(scheduleTime)
	if duration <= 0 {
		fmt.Println("计划时间已过")
		return
	}

	fmt.Printf("任务已安排，将在 %s 后执行\n", duration)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	<-timer.C

	groups := bot.ListGroups()
	if len(groups) > 0 {
		if err := bot.SendToGroups(groups, message); err != nil {
			log.Printf("定时推送失败: %v", err)
		} else {
			log.Printf("定时推送成功: %s", message)
		}
	}
}
