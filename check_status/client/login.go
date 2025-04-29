package client

import (
	"io/ioutil"
	"learn/check_status/config"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

type Client struct {
	account *config.AccountConfig
	Client  *http.Client
}

func NewClient(account *config.AccountConfig) *Client {
	client := &http.Client{}
	return &Client{
		account: account,
		Client:  client,
	}
}

func (cl *Client) Login() {
	loginUrl := "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	lt, execution, client := find(loginUrl)
	cl.Client = client //保存client以复用

	//创建表单
	formData := url.Values{
		"username":  {cl.account.Username},
		"password":  {cl.account.Password},
		"lt":        {lt},
		"execution": {execution},
		"_eventId":  {"submit"},
		"submit":    {"%E7%99%BB%E5%BD%95"},
	}

	// 创建 POST 请求
	//strings.NewReader将其变为io.Reader格式
	//formData.Encode() 将表单改为http可阅读的格式 application/x-www-form-urlencoded
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("origin", "https://account.ccnu.edu.cn")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	// 发送登录请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
}

// 函数寻找 lt和 execution
func find(loginUrl string) (lt string, execution string, client *http.Client) {

	//创建带cookie的HTTP客户端
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}

	res, err := client.Get(loginUrl)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	//使用正则表达式找到It和execution
	ret := regexp.MustCompile(` <input type="hidden" name="lt" value="(.*?)" />`)
	lts := ret.FindAllStringSubmatch(string(body), -1) //查找所有的lt
	//fmt.Println(lts)
	lt = lts[0][1]
	//fmt.Print(lt)

	//获取execution
	ret1 := regexp.MustCompile(` <input type="hidden" name="execution" value="(.*?)" />`)
	executions := ret1.FindAllStringSubmatch(string(body), -1)
	execution = executions[0][1]
	return lt, execution, client
}
