# 接口说明 #
BaseURL:

## 注册 ##

```
url: /uauth/signup
method: POST
body:
  {
    "Email" : "xxx",
    "Password" : "xxx",
    "Redirect":"xxxx"
  }

response:
  {
    "Status" : code,
    "Message" : "xxx"
  }
  
  curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/signup -d '{"Email":"work_test_a@163.com","Password":"123456","Redirect":"http://www.coderpage.com"}'
```

## 登录 ##
```
url: /uauth/signin
method: POST
body:
{
	"Email": "xxx"
	"Password": "xxx",
	"Web" : "xxx",
	"Duration": "xxx"
}

response:
{
	"Status" : code,
	"Message": "xxx",
	"AuthToken": "xxx"
	"User":{
		"Id":"",
		"UserName": "",
		"Email":"",
		"DisplayName":"",
		"Url":"",
		"Created":"",
		"Activate":"",
		"Logged":"",
		"Group":""
	}
}

curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/signin -d '{"Email":"work_test_a@163.com","Password":"123456","Web":"http://www.coderpage.com","Duration":"1"}'
```

## 获取 USER ##
```
url: /uauth/find/user/withtk
method: POST
body:
{
	"AuthToken": "xxx"
}

response:
{
	"Status": code,
	"Message": "xxx",
	"User":{
		...
	}
}

curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/find/user/withtk -d '{"AuthToken":"fP3zpaBL91U49jnr4qbHLesDpd0dhC8T"}'
```
## 重置密码 ##
```
// 发送重置密码邮件
url: /uauth/user/fpwd/email

curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/user/fpwd/email -d '{"Email":"work_test_a@163.com", "Redirect":"http//www.coderpage.com"}'

// 重置密码
url: /uauth/user/resetpwd
curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/user/resetpwd -d '{"ResetToken":"XVkEtsy1X3AL3LKbVfjqB6AkYW2CLqGA","NewPassword":"1234567"}'
```