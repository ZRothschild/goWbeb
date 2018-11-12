# OAuth2.0认证授权
## 认证授权过程
- 服务提供方(provider)，用户使用服务提供方来存储受保护的资源，如照片，视频，联系人列表
- 用户，存放在服务提供方的受保护的资源的拥有者
- 客户端，要访问服务提供方资源的第三方应用，通常是网站，如提供照片打印服务的网站。在认证过程之前，客户端要向服务提供者申请客户端标识

**使用OAuth进行认证和授权的过程如下所示:**
> 用户想操作存放在服务提供方的资源。
> 用户登录客户端向服务提供方请求一个临时令牌。
> 服务提供方验证客户端的身份后，授予一个临时令牌。
> 客户端获得临时令牌后，将用户引导至服务提供方的授权页面请求用户授权。在这个过程中将临时令牌和客户端的回调连接发送给服务提供方。
> 用户在服务提供方的网页上输入用户名和密码，然后授权该客户端访问所请求的资源。
> 授权成功后，服务提供方引导用户返回客户端的网页。
> 客户端根据临时令牌从服务提供方那里获取访问令牌。
> 服务提供方根据临时令牌和用户的授权情况授予客户端访问令牌。
> 客户端使用获取的访问令牌访问存放在服务提供方上的受保护的资源。
## 目录结构
> 主目录`basicauth`

```html
    —— main.go
    —— templates
        —— index.html
        —— user.html
```
## 代码示例
> `main.go`

```go
package main

// 任何OAuth2（甚至是纯golang/x/net/oauth2）包
// 可以与iris一起使用，但在这个例子中我们将看到markbates' goth：
// 获取包  go get github.com/markbates/goth/...

// 这个OAuth2示例适用于会话，因此我们需要
// 附加会话管理器
// 可选：为了更安全的会话值，
// 开发人员可以使用任何第三方包添加自定义cookie编码器/解码器。
// 在这个例子中，我们将使用gorilla的securecookie：
// 获取包 go get github.com/gorilla/securecookie
// securecookie的示例可以在“sessions / securecookie”示例文件夹中找到。

// 注意:
// 整个示例由markbates/goth/example/main.go转换。
// 它使用我自己的TWITTER应用程序进行测试，即使对于localhost也可以使用它。
// 我猜其他一切都按预期工作，goth库社区报告了所有错误
// 在我写这个例子的时候修好了，玩得开心！
import (
	"errors"
	"os"
	"sort"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/gorilla/securecookie"  //可选，用于seesion的编码器/解码器
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/dailymotion"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/fitbit"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/meetup"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
)

var sessionsManager *sessions.Sessions

func init() {
	//附加session管理器
	cookieName := "mycustomsessionid"
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	// AES仅支持16,24或32字节的密钥大小。
	// 您需要准确提供该字节数，或者从您键入的内容中获取密钥。
	hashKey := []byte("the-big-and-secret-fash-key-here")
	blockKey := []byte("lot-secret-of-characters-big-too")
	secureCookie := securecookie.New(hashKey, blockKey)

	sessionsManager = sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Decode: secureCookie.Decode,
	})
}
//下面是一些辅助函数

// GetProviderName函数是用于获取提供者名称（授权应用名称）通过请求。
// 默认情况下，将从URL查询字符串中提取此授权应用名称。
// 如果您以不同的方式提供， 将自己的函数分配给返回提供者的变量用您的请求的名称。
var GetProviderName = func(ctx iris.Context) (string, error) {
	//尝试从的url参数中获取provider
	if p := ctx.URLParam("provider"); p != "" {
		return p, nil
	}
	//尝试从url PATH参数“{provider}或：provider或{provider：string}或{provider：alphabetical}”获取它
	if p := ctx.Params().Get("provider"); p != "" {
		return p, nil
	}
	//尝试从上下文的每个请求存储中获取它
	if p := ctx.Values().GetString("provider"); p != "" {
		return p, nil
	}
	//如果没有找到，则返回一个带有相应错误的空字符串
	return "", errors.New("you must select a provider")
}
/*
BeginAuthHandler是用于启动身份验证过程的便捷处理程序。
它希望能够从查询参数中获取提供程序的名称授权应用名称）
作为“provider”或“：provider”。

BeginAuthHandler会将用户重定向到相应的身份验证端点对于请求的provider。
请参阅https://github.com/markbates/goth/examples/main.go以查看此操作。
*/
func BeginAuthHandler(ctx iris.Context) {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("%v", err)
		return
	}
	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}
/*
GetAuthURL使用provided的请求启动身份验证过程。
它将返回一个应该用于向用户发送的URL。

它希望能够从查询参数中获取provider的名称
作为“provider”或“：provider”或来自“provider”键的上下文值。

我建议使用BeginAuthHandler而不是执行所有这些步骤，但那完全取决于你。
*/
func GetAuthURL(ctx iris.Context) (string, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}
	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	session := sessionsManager.Start(ctx)
	session.Set(providerName, sess.Marshal())
	return url, nil
}

// SetState设置与给定请求关联的状态字符串。
// 如果没有状态字符串与请求相关联，则会生成一个。
// 此状态发送给provider，可以在提取期间检索回调
var SetState = func(ctx iris.Context) string {
	state := ctx.URLParam("state")
	if len(state) > 0 {
		return state
	}
	return "state"
}

// GetState获取回调期间provider的返回的状态。
// 这用于防止CSRF攻击，请参阅 http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx iris.Context) string {
	return ctx.URLParam("state")
}
/*
CompleteUserAuth在锡上做了它所说的。 它完成了身份验证处理并从provider处获取有关用户的所有基本信息。

它希望能够从查询参数中获取provider的名称作为“provider”或“：provider”。

请参阅https://github.com/markbates/goth/examples/main.go以查看此操作。
*/
var CompleteUserAuth = func(ctx iris.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}
	session := sessionsManager.Start(ctx)
	value := session.GetString(providerName)
	if value == "" {
		return goth.User{}, errors.New("session value for " + providerName + " not found")
	}
	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}
	user, err := provider.FetchUser(sess)
	if err == nil {
		//可以找到现有session数据的用户
		return user, err
	}
	//获取新令牌并重试获取
	_, err = sess.Authorize(provider, ctx.Request().URL.Query())
	if err != nil {
		return goth.User{}, err
	}
	session.Set(providerName, sess.Marshal())
	return provider.FetchUser(sess)
}

// 注销使用户session
func Logout(ctx iris.Context) error {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return err
	}
	session := sessionsManager.Start(ctx)
	session.Delete(providerName)
	return nil
}

//一些函数助手的结尾 设置key secret 与回调方法
func main() {
	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),
		//如果您想在Twitter提供程序中使用authenticate而不是authorize，请改用它。
		// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),

		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
		fitbit.New(os.Getenv("FITBIT_KEY"), os.Getenv("FITBIT_SECRET"), "http://localhost:3000/auth/fitbit/callback"),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), "http://localhost:3000/auth/gplus/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
		spotify.New(os.Getenv("SPOTIFY_KEY"), os.Getenv("SPOTIFY_SECRET"), "http://localhost:3000/auth/spotify/callback"),
		linkedin.New(os.Getenv("LINKEDIN_KEY"), os.Getenv("LINKEDIN_SECRET"), "http://localhost:3000/auth/linkedin/callback"),
		lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), "http://localhost:3000/auth/lastfm/callback"),
		twitch.New(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), "http://localhost:3000/auth/twitch/callback"),
		dropbox.New(os.Getenv("DROPBOX_KEY"), os.Getenv("DROPBOX_SECRET"), "http://localhost:3000/auth/dropbox/callback"),
		digitalocean.New(os.Getenv("DIGITALOCEAN_KEY"), os.Getenv("DIGITALOCEAN_SECRET"), "http://localhost:3000/auth/digitalocean/callback", "read"),
		bitbucket.New(os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"), "http://localhost:3000/auth/bitbucket/callback"),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), "http://localhost:3000/auth/instagram/callback"),
		intercom.New(os.Getenv("INTERCOM_KEY"), os.Getenv("INTERCOM_SECRET"), "http://localhost:3000/auth/intercom/callback"),
		box.New(os.Getenv("BOX_KEY"), os.Getenv("BOX_SECRET"), "http://localhost:3000/auth/box/callback"),
		salesforce.New(os.Getenv("SALESFORCE_KEY"), os.Getenv("SALESFORCE_SECRET"), "http://localhost:3000/auth/salesforce/callback"),
		amazon.New(os.Getenv("AMAZON_KEY"), os.Getenv("AMAZON_SECRET"), "http://localhost:3000/auth/amazon/callback"),
		yammer.New(os.Getenv("YAMMER_KEY"), os.Getenv("YAMMER_SECRET"), "http://localhost:3000/auth/yammer/callback"),
		onedrive.New(os.Getenv("ONEDRIVE_KEY"), os.Getenv("ONEDRIVE_SECRET"), "http://localhost:3000/auth/onedrive/callback"),

		//将localhost.com指向http：//localhost：3000/auth/yahoo/通过代理回调作为雅虎
		//不允许在重定向uri中放置自定义端口
		yahoo.New(os.Getenv("YAHOO_KEY"), os.Getenv("YAHOO_SECRET"), "http://localhost.com"),
		slack.New(os.Getenv("SLACK_KEY"), os.Getenv("SLACK_SECRET"), "http://localhost:3000/auth/slack/callback"),
		stripe.New(os.Getenv("STRIPE_KEY"), os.Getenv("STRIPE_SECRET"), "http://localhost:3000/auth/stripe/callback"),
		wepay.New(os.Getenv("WEPAY_KEY"), os.Getenv("WEPAY_SECRET"), "http://localhost:3000/auth/wepay/callback", "view_user"),
		//默认使用paypal production auth urls，请将PAYPAL_ENV = sandbox设置为环境变量进行测试
		//在沙盒环境中
		paypal.New(os.Getenv("PAYPAL_KEY"), os.Getenv("PAYPAL_SECRET"), "http://localhost:3000/auth/paypal/callback"),
		steam.New(os.Getenv("STEAM_KEY"), "http://localhost:3000/auth/steam/callback"),
		heroku.New(os.Getenv("HEROKU_KEY"), os.Getenv("HEROKU_SECRET"), "http://localhost:3000/auth/heroku/callback"),
		uber.New(os.Getenv("UBER_KEY"), os.Getenv("UBER_SECRET"), "http://localhost:3000/auth/uber/callback"),
		soundcloud.New(os.Getenv("SOUNDCLOUD_KEY"), os.Getenv("SOUNDCLOUD_SECRET"), "http://localhost:3000/auth/soundcloud/callback"),
		gitlab.New(os.Getenv("GITLAB_KEY"), os.Getenv("GITLAB_SECRET"), "http://localhost:3000/auth/gitlab/callback"),
		dailymotion.New(os.Getenv("DAILYMOTION_KEY"), os.Getenv("DAILYMOTION_SECRET"), "http://localhost:3000/auth/dailymotion/callback", "email"),
		deezer.New(os.Getenv("DEEZER_KEY"), os.Getenv("DEEZER_SECRET"), "http://localhost:3000/auth/deezer/callback", "email"),
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), "http://localhost:3000/auth/discord/callback", discord.ScopeIdentify, discord.ScopeEmail),
		meetup.New(os.Getenv("MEETUP_KEY"), os.Getenv("MEETUP_SECRET"), "http://localhost:3000/auth/meetup/callback"),

		// Auth0为每个客户分配域，必须提供域以使auth0正常工作
		auth0.New(os.Getenv("AUTH0_KEY"), os.Getenv("AUTH0_SECRET"), "http://localhost:3000/auth/auth0/callback", os.Getenv("AUTH0_DOMAIN")),
		xero.New(os.Getenv("XERO_KEY"), os.Getenv("XERO_SECRET"), "http://localhost:3000/auth/xero/callback"),
	)
	// OpenID Connect基于OpenID Connect自动发现URL（https://openid.net/specs/openid-connect-discovery-1_0-17.html）
	//因为OpenID Connect提供程序在New（）中初始化它，它可以返回应该处理或忽略的错误
	//暂时忽略错误
	openidConnect, _ := openidConnect.New(os.Getenv("OPENID_CONNECT_KEY"), os.Getenv("OPENID_CONNECT_SECRET"), "http://localhost:3000/auth/openid-connect/callback", os.Getenv("OPENID_CONNECT_DISCOVERY_URL"))
	if openidConnect != nil {
		goth.UseProviders(openidConnect)
	}
	m := make(map[string]string)
	m["amazon"] = "Amazon"
	m["bitbucket"] = "Bitbucket"
	m["box"] = "Box"
	m["dailymotion"] = "Dailymotion"
	m["deezer"] = "Deezer"
	m["digitalocean"] = "Digital Ocean"
	m["discord"] = "Discord"
	m["dropbox"] = "Dropbox"
	m["facebook"] = "Facebook"
	m["fitbit"] = "Fitbit"
	m["github"] = "Github"
	m["gitlab"] = "Gitlab"
	m["soundcloud"] = "SoundCloud"
	m["spotify"] = "Spotify"
	m["steam"] = "Steam"
	m["stripe"] = "Stripe"
	m["twitch"] = "Twitch"
	m["uber"] = "Uber"
	m["wepay"] = "Wepay"
	m["yahoo"] = "Yahoo"
	m["yammer"] = "Yammer"
	m["gplus"] = "Google Plus"
	m["heroku"] = "Heroku"
	m["instagram"] = "Instagram"
	m["intercom"] = "Intercom"
	m["lastfm"] = "Last FM"
	m["linkedin"] = "Linkedin"
	m["onedrive"] = "Onedrive"
	m["paypal"] = "Paypal"
	m["twitter"] = "Twitter"
	m["salesforce"] = "Salesforce"
	m["slack"] = "Slack"
	m["meetup"] = "Meetup.com"
	m["auth0"] = "Auth0"
	m["openid-connect"] = "OpenID Connect"
	m["xero"] = "Xero"
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}
	//创建我们的应用，
	//设置一个视图
	//设置sessions
	//并为展示设置路由器
	app := iris.New()
	//附加并构建我们的模板
	app.RegisterView(iris.HTML("./templates", ".html"))
	//启动路由器
	app.Get("/auth/{provider}/callback", func(ctx iris.Context) {
		user, err := CompleteUserAuth(ctx)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Writef("%v", err)
			return
		}
		ctx.ViewData("", user)
		if err := ctx.View("user.html"); err != nil {
			ctx.Writef("%v", err)
		}
	})
	app.Get("/logout/{provider}", func(ctx iris.Context) {
		Logout(ctx)
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
	})
	app.Get("/auth/{provider}", func(ctx iris.Context) {
		//尝试让用户无需重新进行身份验证
		if gothUser, err := CompleteUserAuth(ctx); err == nil {
			ctx.ViewData("", gothUser)
			if err := ctx.View("user.html"); err != nil {
				ctx.Writef("%v", err)
			}
		} else {
			BeginAuthHandler(ctx)
		}
	})
	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("", providerIndex)
		if err := ctx.View("index.html"); err != nil {
			ctx.Writef("%v", err)
		}
	})
	// http://localhost:3000
	app.Run(iris.Addr("localhost:3000"))
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}
```
> `index.html`

```html
{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}
```
> `user.html`

```html
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
```
### 提示
1. 去Github设置第三方授权登录，设置回调路径，把key,secret填入上面代码指定位置
2. 执行程序选择github 即可看到想要的效果