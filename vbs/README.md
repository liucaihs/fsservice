
## 本程序主要提供手机号及其验证码信息的下发服务。


### 1 从本服务器获取手机号
- **接口地址：** /1/mobile
- **请求方法：** HTTP GET方法

#### 1.1 返回结果
##### A. 服务器暂时没有手机号码时（http头部状态码为404），
	{
		"desc": "Sorry, there is no data that you need currently! Please try again later~",
	}
##### B. 成功接收并处理了客户端的请求时（http头部状态码为200）：
	{
		"mobile":    手机号,
		"imsi":      用户的SIM卡信息,
		"imei":      设备ID,
		"iccid":     SIM卡编号信息,
		"ip":        提供此号码方的IP地址  
	}


### 2 从本服务器获取验证码信息
- **接口地址：** /vcode/:pkg/:phone
- **请求方法：** HTTP GET方法

#### 2.1 url地址中相关字段说明
变量名     |出现要求	     |描述  
:----	  |:---		     |:------	
pkg		  |不能为空	     |需要使用此验证码的应用包名
phone	  |不能为空	     |使用此验证码的手机号

#### 请求示例：
```
若在“咪咕音乐”中进行注册操作时使用的手机号码为“18811620823”, 则可发如下请求来索取其验证码：

        http://服务器IP：8093/vcode/cmccwm.mobilemusic/18811620823    
```

#### 2.2 返回结果
##### A. 服务器暂时没有收到该手机号的验证码时（http头部状态码为404），
	{
		"desc": "Sorry, there is no data that you need currently! Please try again later~",
	}
##### B. 成功接收并处理了客户端的请求时（http头部状态码为200）：
	{
		"phone":    手机号,
		"vcode":    验证码，
		"pkg":      应用包名
	}
