# 接口说明 #

BaseURL: localhost:9090

## 注册账号 ##

**URL:** `/uauth/signup`

**Method:** `POST`

**POST BODY: **

```
  {
    "Email" : "xxx",
    "Password" : "xxx",
    "Redirect":"xxxx"
  }
```
> 目前只支持**邮箱**注册<br>
> *Email*: 注册邮箱<br>
> *Password*: 账号密码<br>
> *Redirect*: 回调地址，注册成功后会向邮箱中发送激活链接，激活成功后，会重定向此链接<br>

**Response**

```
  {
    "Status" : code,
    "Message" : "xxx"
  }
```
**Status 说明：**

|Status|说明|
|---|---|
|-1|出现未知错误|
|200|注册成功|
|422|服务器读取解析 Post Body 失败|
|430|用户名不合法|
|431|密码不合法|
|432|用户已经存在|

**回调**

用户点击邮箱激活链接后，uAuth 完成激活会跳转到指定的回到链接，链接携带 `status` 和 `msg` 参数，例：
*http://www.coderpage.com?status=200&msg=success*

**激活回调 Status 说明**

|Status|说明|
|---|---|
|200|激活成功|
|403|Token 过期|
|404|Token 错误，或因为过期已经没服务器删除|
|500|服务器出错|

**curl 测试**

```
curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/signup -d '{"Email":"820631390@qq.com","Password":"123456","Redirect":"http://www.coderpage.com"}'
```

## 重新发送激活邮件 ##
**URL:** `/uauth/user/active/sendemail`

**Method:** `POST`

**POST BODY: **

```
  {
    "Email" : "xxx",
    "Redirect" : "xxx"
  }
```

> *Email*: 注册邮箱<br>
> *Redirect*: 回调地址，注册成功后会向邮箱中发送激活链接，激活成功后，会重定向此链接<br>

**Response**

```
  {
    "Status" : code,
    "Message" : "xxx"
  }
```
**Status 说明：**

|Status|说明|
|---|---|
|200|成功|
|400|请求参数错误|
|404|用户不存在|
|436|发送邮件失败|

**curl 测试**

```
curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/user/active/sendemail -d '{"Email":"820631390@qq.com","Redirect":"http://www.coderpage.com"}'
```


## 登录 ##

**URL:** `/uauth/signin`

**Method:** `POST`

**POST BODY: **

```
  {
    "Email" : "xxx",
    "Password" : "xxx",
	"Web" : "xxx",
	"Duration": "xxx"
  }
```

> *Email*: 邮箱<br>
> *Password*: 账号密码<br>
> *Web*: 请求登录 Web 的地址<br>
> *Duration*: 本次登录的有效市场，小时为单位<br>

**Response**

```
{
	"Status" : code,
	"Message": "xxx",
	"AuthToken": "xxx",
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
```
**Status 说明：**

|Status|说明|
|---|---|
|-1|出现未知错误|
|200|登录成功|
|400|请求参数有错误|
|422|服务器读取解析 Post Body 失败|
|434|用户未激活|
|435|用户名或密码错误|
|500|服务器出错|

**curl 测试**

```
curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/signin -d '{"Email":"820631390@qq.com","Password":"123456","Web":"http://www.coderpage.com","Duration":"1"}'
```

## 获取 USER ##

**URL:** `/uauth/find/user/withtk`

**Method:** `POST`

**POST BODY: **

```
  {
    "AuthToken" : "xxx"
  }
```
> *AuthToken*: 登录成功后返回的令牌<br>

**Response**

```
{
	"Status" : code,
	"Message": "xxx",
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
```
**Status 说明：**

|Status|说明|
|---|---|
|-1|出现未知错误|
|200|成功|
|403|Token 不可用|
|404|此 Token 找不到对应的用户信息|
|422|服务器读取解析 Post Body 失败|
|434|用户未激活|

**curl 测试**

```
curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/find/user/withtk -d '{"AuthToken":"jbE3vU5fu6oYgGCHHMhsnnY1PWo9Fbe8"}'
```

## 找回密码 ##
### 通过注册邮箱找回密码 ###

#### 发送邮件 ####

**URL:** `/uauth/user/fpwd/email`

**Method:** `POST`

**POST BODY: **

```
  {
    "Email" : "xxx",
    "Redirect": ""
  }
```
> *Email*: 注册邮箱<br>
> *Redirect*: 回调地址，用户点击重置密码链接后，uAuth 完成验证会跳转到指定的回到链接<br>

**Response**

```
{
	"Status" : code,
	"Message": "xxx",
}
```
**Status 说明：**

|Status|说明|
|---|---|
|-1|出现未知错误|
|200|成功|
|400|请求参数错误|
|404|用户不存在|
|422|服务器读取解析 Post Body 失败|

**回调**

用户重置密码链接后，uAuth 完成验证会跳转到指定的回到链接，链接携带 `status` 和 `reset` 参数，例：
*http://www.coderpage.com?status=200&reset=xxxx* ，`reset` 对应值为重置密码的令牌

**激活回调 Status 说明**

|Status|说明|
|---|---|
|200|成功|
|403|Token 不可用|
|404|Token 不存在|
|500|服务器出错|

**curl 测试**

```
curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/user/fpwd/email -d '{"Email":"work_test_a@163.com", "Redirect":"http://www.coderpage.com"}'
```
#### 重置密码 ####

**URL:** `/uauth/user/resetpwd`

**Method:** `POST`

**POST BODY: **

```
  {
    "ResetToken" : "xxx",
    "NewPassword": ""
  }
```
> *ResetToken*: 重置密码令牌<br>
> *NewPassword*: 新密码<br>

**Response**

```
{
	"Status" : code,
	"Message": "xxx",
}
```
**Status 说明：**

|Status|说明|
|---|---|
|-1|出现未知错误|
|200|成功|
|400|请求参数错误|
|403|Token 不可用|
|404|此 Token 找不到对应的用户信息|
|422|服务器读取解析 Post Body 失败|
|500|服务器出错|

**curl 测试**

```
curl -i -X POST -H "Content-Type: application/json" http://localhost:9090/uauth/user/resetpwd -d '{"ResetToken":"XVkEtsy1X3AL3LKbVfjqB6AkYW2CLqGA","NewPassword":"1234567"}
```