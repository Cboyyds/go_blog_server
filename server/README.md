├── server
├── api               (api层)
├── assets            (静态资源包)
├── config            (配置包)
├── core              (核心文件)
├── flag              (flag命令)
├── global            (全局对象)
├── initialize        (初始化)
├── log               (日志文件)
├── middleware        (中间件层)
├── model             (模型层)
│   ├── appTypes      (自定义类型)
│   ├── database      (mysql结构体)
│   ├── elasticsearch (es结构体)
│   ├── other         (其他结构体)
│   ├── request       (入参结构体)
│   └── response      (出参结构体)
├── router            (路由层)
├── service           (service层)
├── task              (定时任务包)
├── uploads           (文件上传目录)
└── utils             (工具包)
├── hotSearch    (热搜接口封装)
└── upload        (oss接口封装)
[Go春绝迹-绝迹之春的个人博客](https://www.scc749.cn/)

问题：

core

server_other.go

### 1，server_win.go 里面存在相同的函数名，但是小写，报错了，但是视频里面还没讲，用了下面这个就解决了

[Windows环境下github.com\fvbock\endless库报错：undefined: syscall.SIGUSR1_go undefined: syscall.sigusr1-CSDN博客](https://blog.csdn.net/hblzong/article/details/140820067)

### 2，数据库表的创建：

![image-20250424220452731](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250424220452731.png)

(
`id` bigint unsigned NOT NULL AUTO_INCREMENT,
`created_at` datetime(3) DEFAULT NULL,
`updated_at` datetime(3) DEFAULT NULL,
`deleted_at` datetime(3) DEFAULT NULL,
`article_id` longtext,
`p_id` bigint unsigned DEFAULT NULL,
`user_uuid` char(36) DEFAULT NULL,
`content` longtext,
PRIMARY KEY (`id`),
KEY `idx_comments_deleted_at` (`deleted_at`),
KEY `fk_comments_children` (`p_id`),
CONSTRAINT `fk_comments_children` FOREIGN KEY (`p_id`) REFERENCES `comments` (`id`)
)

![image-20250425090946981](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250425090946981.png)

### 3，binding:"required,email" 表示在数据绑定时，该字段是必填项，并且必须符合电子邮件格式。

![image-20250514170721988](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250514170721988.png)

### 4，密钥生成

```
go get "github.com/gin-contrib/sessions/cookie"
```

新建一个go文件

```
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomKey(length int) (string, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}
func main() {
	kenLength := 32
	randomKey, err := GenerateRandomKey(kenLength)
	if err != nil {
		panic(err)
	}
	fmt.Println("Randomly generated key:", randomKey)
}

```

![image-20250515101922888](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250515101922888.png)

### 5 ，gin的默认日志会打印到1，如何将它打印到日志文件

![image-20250515190508558](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250515190508558.png)

```
package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"server/global"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 是一个 Gin 中间件，用于记录请求日志。
// 该中间件会在每次请求结束后，使用 Zap 日志记录请求信息。
// 通过此中间件，可以方便地追踪每个请求的情况以及性能。
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 获取请求的路径和查询参数
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 继续执行后续的处理
		c.Next()

		// 计算请求处理的耗时
		cost := time.Since(start)

		// 使用 Zap 记录请求日志
		global.Log.Info(path,
			// 记录响应状态码
			zap.Int("status", c.Writer.Status()),
			// 记录请求方法
			zap.String("method", c.Request.Method),
			// 记录请求路径
			zap.String("path", path),
			// 记录查询参数
			zap.String("query", query),
			// 记录客户端 IP
			zap.String("ip", c.ClientIP()),
			// 记录 User-Agent 信息
			zap.String("user-agent", c.Request.UserAgent()),
			// 记录错误信息（如果有）
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			// 记录请求耗时
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 是一个 Gin 中间件，用于捕获和处理请求中的 panic 错误。
// 该中间件的主要作用是确保服务在遇到未处理的异常时不会崩溃，并通过日志系统提供详细的错误追踪。
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 defer 确保 panic 被捕获，并且处理函数会在 panic 后执行
		defer func() {
			// 检查是否发生了 panic 错误
			if err := recover(); err != nil {
				// 检查是否是连接被断开的问题（如 broken pipe），这些错误不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求信息，包括请求体等
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 如果是 broken pipe 错误，则只记录错误信息，不记录堆栈信息
				if brokenPipe {
					global.Log.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 由于连接断开，不能再向客户端写入状态码
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()                // 中止请求处理
					return
				}

				// 如果是其他类型的 panic，根据 `stack` 参数决定是否记录堆栈信息
				if stack {
					// 记录详细的错误和堆栈信息
					global.Log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					// 只记录错误信息，不记录堆栈
					global.Log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 返回 500 错误状态码，表示服务器内部错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 继续执行后续的请求处理
		c.Next()
	}
}
```

![image-20250515190642700](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250515190642700.png)

### 6，双token

![image-20250515191206213](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250515191206213.png)

```
以下是双Token机制的工作过程：
1. 用户登录：
• 用户通过输入用户名和密码进行身份验证。
• 服务器验证用户凭据，若验证通过，生成两个令牌：Access Token和Refresh Token。
2. Access Token：
• 这是一个短期有效的令牌，用于访问受保护的资源。
• 它通常包含用户的身份信息和权限声明。
• 由于有效期短，即使被窃取，攻击者利用它进行恶意操作的时间窗口也有限。
3. Refresh Token：
• 这是一个长期有效的令牌，用于在Access Token过期后获取新的Access Token，而无需用户重新登录。
• 它通常存储在更安全的位置，如HTTP-Only Cookie中，以降低被窃取的风险。
4. 访问受保护资源：
• 客户端在后续请求中携带Access Token，服务器验证该令牌的有效性。
• 如果令牌有效，服务器处理请求并返回响应。
5. Access Token过期：
• 当Access Token过期后，客户端在下次请求时会收到服务器返回的“令牌过期”错误。
6. 刷新Access Token：
• 客户端使用Refresh Token向服务器发起刷新请求。
• 服务器验证Refresh Token的有效性，如果有效，会生成一个新的Access Token，并返回给客户端。
7. 新的访问资源：
• 客户端使用新的Access Token继续访问受保护资源。
8. Refresh Token过期：
• 当Refresh Token也过期时，用户需要重新登录以获取新的访问令牌和刷新令牌。这种机制的优势在于它结合了安全性和用户体验。Access Token的短期有效性减少了令牌被盗用的风险，而Refresh Token的存在则避免了用户频繁重新登录的不便。
此外，Refresh Token可以设计为在每次成功刷新Access Token时更新，或者维持其不变，具体策略根据业务需求而定。在实际应用中，双Token机制还可以包括一些额外的安全措施，例如令牌续期、令牌撤销等，以进一步增强系统的安全性。
```

先学习理解双token验证：

[gin框架使用jwt, 双token刷新,续期 - dogRuning - 博客园](https://www.cnblogs.com/dogHuang/p/16621331.html)

#### 双token的刷新 access_token和refresh_token

第一次用账号密码登录服务器会返回两个 token : access_token 和 refresh_token，时效长短不一样。短的access_token 时效过了之后，发送时效长的 refresh_token 重新获取一个短时效token，如果都过期，就需要重新登录了。

***双token的使用步骤要清晰：***

如何生成对应的token，

如何将token发送给对应的前端，（如刚注册时，刚登录时，还有过期时等一些情况）

### 7.验证码部分

api/base.go

```
package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type BaseApi struct {
}

var store = base64Captcha.DefaultMemStore

// Captcha 生成数字验证码
func (baseApi *BaseApi) Captcha(c *gin.Context) {
	// 创建数字验证码的驱动
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)
	
	// 创建验证码对象
	captcha := base64Captcha.NewCaptcha(driver, store)
	
	// 生成验证码
	id, b64s, _, err := captcha.Generate()
	
	if err != nil {
		global.Log.Error("Failed to generate captcha:", zap.Error(err))
		response.FailWithMessage("Failed to generate captcha", c)
		return
	}
	response.OkWithData(response.Captcha{
		CaptchaID: id,
		PicPath:   b64s,
	}, c)
}

// SendEmailVerificationCode 处理发送邮件验证码的请求。
// 该函数首先从请求中解析出验证码和邮箱信息，然后验证验证码是否正确。
// 如果验证码正确，则调用基础服务发送邮件验证码，并根据发送结果返回相应的响应。
// 如果验证码错误或发送邮件失败，则返回错误信息。
//
// 参数:
//   - c: *gin.Context, Gin框架的上下文对象，用于处理HTTP请求和响应。
//
// 返回值:
//   无返回值，但通过Gin的上下文对象返回HTTP响应。
func (baseApi *BaseApi) SendEmailVerificationCode(c *gin.Context) {
	// 解析请求体中的JSON数据到SendEmailVerificationCode结构体
	var req request.SendEmailVerificationCode
	err := c.ShouldBindJSON(&req)
	if err != nil {
		// 如果解析失败，返回错误信息
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	// 验证验证码是否正确
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		// 如果验证码正确，调用基础服务发送邮件验证码
		err = baseService.SendEmailVerificationCode(c, req.Email)
		if err != nil {
			// 如果发送邮件失败，记录错误日志并返回错误信息
			global.Log.Error("Failed to send email:", zap.Error(err))
			response.FailWithMessage("Failed to send email", c)
			return
		}
		// 发送成功，返回成功信息
		response.OkWithMessage("Successfully sent email", c)
		return
	}
	
	// 验证码错误，返回错误信息
	response.FailWithMessage("Incorrect verification code", c)
}

// QQLoginURL 返回 QQ 登录链接
func (baseApi *BaseApi) QQLoginURL(c *gin.Context) {
	url := global.Config.QQ.QQLoginURL()
	response.OkWithData(url, c)
}

```



### 8.注册时

initial/router.go

```
var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	Router.Use(sessions.Sessions("session", store))
```

service/base.go

```
// 将验证码、验证邮箱、过期时间存入会话中
	session := sessions.Default(c)
	session.Set("verification_code", verificationCode)
	session.Set("email", to)
	session.Set("expire_time", expireTime)
	_ = session.Save()
```

service\base.go，里面需要保存session信息，然后发送邮件

```
package service

import (
	"server/global"
	"server/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type BaseService struct {
}

func (baseService *BaseService) SendEmailVerificationCode(c *gin.Context, to string) error {
	verificationCode := utils.GenerateVerificationCode(6)
	expireTime := time.Now().Add(5 * time.Minute).Unix()

	// 将验证码、验证邮箱、过期时间存入会话中
	session := sessions.Default(c)
	session.Set("verification_code", verificationCode)
	session.Set("email", to)
	session.Set("expire_time", expireTime)
	_ = session.Save()

	subject := "您的邮箱验证码"
	body := `亲爱的用户[` + to + `]，<br/>
<br/>
感谢您注册` + global.Config.Website.Name + `的个人博客！为了确保您的邮箱安全，请使用以下验证码进行验证：<br/>
<br/>
验证码：[<font color="blue"><u>` + verificationCode + `</u></font>]<br/>
该验证码在 5 分钟内有效，请尽快使用。<br/>
<br/>
如果您没有请求此验证码，请忽略此邮件。
<br/>
如有任何疑问，请联系我们的支持团队：<br/>
邮箱：` + global.Config.Email.From + `<br/>
<br/>
祝好，<br/>` +
		global.Config.Website.Title + `<br/>
<br/>`

	_ = utils.Email(to, subject, body)

	return nil
}
```

api/user.go 的register里面

```
//"github.com/gin-contrib/sessions"
session := sessions.Default(c) // session.Session
	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}
```

发送邮件utils/email.go

```
package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"server/global"
	"strings"

	"github.com/jordan-wright/email"
)

// Email 发送电子邮件
func Email(To, subject string, body string) error {
	to := strings.Split(To, ",") // 将收件人邮箱地址按逗号拆分成多个地址
	return send(to, subject, body)
}

// send 执行邮件发送操作
func send(to []string, subject string, body string) error {
	emailCfg := global.Config.Email // 获取全局配置中的邮件设置

	from := emailCfg.From
	nickname := emailCfg.Nickname
	secret := emailCfg.Secret
	host := emailCfg.Host
	port := emailCfg.Port
	isSSL := emailCfg.IsSSL

	// 使用 PlainAuth 创建认证信息，用到的是net/stmp包
	auth := smtp.PlainAuth("", from, secret, host)

	// 创建新的电子邮件对象
	e := email.NewEmail()
	if nickname != "" {
		// 如果设置了昵称，则格式化发件人地址为 "昵称 <邮箱>"
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		// 否则直接使用发件人邮箱
		e.From = from
	}

	// 设置收件人、主题和邮件内容
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)

	// 定义错误变量
	var err error
	// 构建邮件服务器的地址，格式为 host:port
	hostAddr := fmt.Sprintf("%s:%d", host, port)

	// 根据配置的是否使用 SSL 来选择邮件发送方法
	if isSSL {
		// 使用带 TLS 的邮件发送
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		// 使用普通的邮件发送
		err = e.Send(hostAddr, auth)
	}

	return err
}

```

### 9，在utils/hotSearch/zhihu.go里面有一句

```
	reg := regexp.MustCompile(`(?s)<script id="js-initialData" type="text/json">(.*?)</script>`)  
	//原来是正则表达式库啊regexp
```

功能如下：

```html
正则表达式：(?s)<script id="js-initialData" type="text/json">(.*?)</script>
(?s)：启用单行模式，使 . 匹配包括换行符在内的所有字符。
```



<script id="js-initialData" type="text/json">：匹配具有特定 id 和 type 的 <script> 标签。
(.*?)：非贪婪匹配，捕获标签之间的所有内容。
功能：该正则表达式用于从HTML中提取 id="js-initialData" 的 <script> 标签内的JSON数据。


### 10.新闻获取

如百度的用正则

这种获取方式如果不给url就难受了

```
reg := regexp.MustCompile(`<!--s-data:({.*?})-->`) // 正则表达式
```

![image-20250522172451737](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250522172451737.png)

提取的数据的大致格式是这样的：

```
{
    "data":
    {
        "cards":
        [

        {
                "component":"hotList",
                "content":
                [

                     {
                        "appUrl":"https://www.baidu.com/s?wd=%E6%80%BB%E4%B9%A6%E8%AE%B0%E7%9A%84%E2%80%9C%E5%BE%85%E5%AE%A2%E8%8C%B6%E2%80%9D&sa=fyb_news&rsv_dl=fyb_news",
                        "desc":"在国际外交舞台上，习近平主席频频以茶会友、以茶论道，在弘扬中华优秀传统文化的同时，阐发出“和合共生”等价值理念。",
                        "hotChange":"same",
                        "hotScore":"7904434",
                        "hotTag":"0",
                        "img":"https://fyb-2.cdn.bcebos.com/hotboard_image/e594fb69f99568ed58a8afc7415fb238",
                        "index":0,
                        "indexUrl":"",
                        "query":"总书记的“待客茶”",
                        "rawUrl":"https://www.baidu.com/s?wd=%E6%80%BB%E4%B9%A6%E8%AE%B0%E7%9A%84%E2%80%9C%E5%BE%85%E5%AE%A2%E8%8C%B6%E2%80%9D",
                        "show":
                        [

                        ],
                        "url":"https://www.baidu.com/s?wd=%E6%80%BB%E4%B9%A6%E8%AE%B0%E7%9A%84%E2%80%9C%E5%BE%85%E5%AE%A2%E8%8C%B6%E2%80%9D&sa=fyb_news&rsv_dl=fyb_news",
                        "word":"总书记的“待客茶”",
                        "isTop":true
                    },

```

### 11，获取主机ip

```
IP:=c.ClientIP()
```

![image-20250523095625587](C:\Users\c博\AppData\Roaming\Typora\typora-user-images\image-20250523095625587.png)

在这里出现了一个小bug，因为本机跑在本地，ip地址为127.0.0.1，去进行请求定位信息请求不到，所以出现天气查询不到的情况了，其实高德的接口还可以只写一个key值也能请求到定位，