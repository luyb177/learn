package tool

import (
	"crypto/tls"
	"fmt"
	"learn/check_status/config"
	"net/smtp"
)

type Mail struct {
	cfg *config.QQConfig
}

func NewMail(cfg *config.QQConfig) *Mail {
	return &Mail{
		cfg: cfg,
	}
}

// SendEmailByQQEmail 发送邮件函数
func (mail *Mail) SendEmailByQQEmail(to, roomName, seat, name, start, end, count string) error {
	from := mail.cfg.Email
	password := mail.cfg.Key // 邮箱授权码
	smtpServer := "smtp.qq.com:465"

	// 设置 PlainAuth
	auth := smtp.PlainAuth("", from, password, "smtp.qq.com")

	// 创建 tls 配置
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.qq.com",
	}

	// 连接到 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, "smtp.qq.com")
	if err != nil {
		return fmt.Errorf("SMTP 客户端创建失败: %v", err)
	}
	defer client.Quit()

	// 使用 auth 进行认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %v", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("发件人设置失败: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("收件人设置失败: %v", err)
	}

	// 写入邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("数据写入失败: %v", err)
	}
	defer wc.Close()

	subject := "【反内卷通知】您已被连续识别为高频卷王"
	body := `
<div style="font-family:Arial, sans-serif; max-width:600px; margin:0 auto; border:1px solid #eee; padding:20px; border-radius:10px;">
	<h2 style="color:#d9534f;">🚨 反内卷系统警告通知</h2>
	<p>Hi <strong>` + name + `</strong> 👀，你又被我们盯上啦！</p>

	<table border="1" cellpadding="8" cellspacing="0" style="border-collapse: collapse; width:100%; margin-top:10px;">
		<tr style="background-color:#f2f2f2;">
			<th>🚪 场地</th>
			<th>🪑 座位</th>
			<th>⏰ 时间</th>
			<th>📸 被抓次数</th>
		</tr>
		<tr>
			<td>` + roomName + `</td>
			<td>` + seat + `</td>
			<td>` + start + ` ~ ` + end + `</td>
			<td>` + count + ` 次</td>
		</tr>
	</table>
	<p style="margin-top:10px; font-size:16px;">
	📢 <strong>警告：</strong>你已经被系统 <strong style="color:red;">连续侦测</strong> 到卷王行为。
	<br>建议你适当摸鱼，卷出健康，卷出风采 ✨
	<br>别再让我们👮天天盯着你咯~
	<br><br>
	—— 来自反内卷总部 ❤️
	</p>
</div>
`

	msg := []byte("From: Sender Name <" + from + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body)

	_, err = wc.Write(msg)
	if err != nil {
		return fmt.Errorf("消息发送失败: %v", err)
	}
	return nil
}
