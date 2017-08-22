
注：本服务建议在docker中运行，推荐为每个渠道专门开设一个服务器，运行前请做好相关分流配置！！

# 1. 接口定义

## 1.1 上传手机号信息到本服务器
- **接口地址：** /phones
- **请求方法：** HTTP POST方法，采用JSON格式

### 1.1.1 Body中的请求参数
  
参数名称						|类型		|出现要求	|描述  
:----						|:---		|:------	|:---	
data				|对象数组	|不能为空！		|&nbsp;
&emsp;mobile		|string		|不能为空！		|手机号
&emsp;cid		    |string		|	--			|分配手机号的渠道号
&emsp;imsi			|string		|	--			|存于SIM卡中的用户标识号
&emsp;imei			|string		|	--			|设备ID
&emsp;iccid			|string		|	--			|SIM卡号
&emsp;ip			|string		|	--			|客户端网络地址

Body中参数内容示例：

```
{
    "data":[
			{
			"mobile":"18250794865",
			"cid":"hello",
			"imsi":"zjgcghsdv",
			"imei":"56416",
			"iccid":"s5d45s", 
			"ip":"8.88.888.886"
			},
			{
			"mobile":"17342060825",
			"cid":"world",
			"imsi":"ghsdv",
			"imei":"12356416",
			"iccid":"ws5dfg45s", 
			"ip":"6.66.666.686"
			}
		   ]
}

```
### 1.1.2 返回结果
#### A. 服务器内部运行出错时
	{
		"code": 500,
		"desc": "The server is busy, please try again later.",
	}		
#### B. 暂时不需要客户端上传的数据时
	{
		"code": 406,
		"desc": "Thanks for your contribution. However, as we have enough data to be used, thus we don not need your data just now. Please devote your data later~",
	}
#### C. 全部接收或者仅使用了部分客户端的数据时，服务端会通知客户端使用了哪些手机号
	{
		"code":       200,
		"msg":        "Thanks for your contribution. And as follows, we have used some datas that you provided, please offer their verify codes later.",
		"phonesUsed": 字符串数组, 例如：["18250794865","17342060825"]
	}
#### D. 客户端传递的某个参数的数据类型或数值不符合要求时，
	{
	    "code": 400,
		"desc": "The data that you have submitted is not all correct !",
	}



## 1.2 从本服务器获取手机号信息
- **接口地址：** /mobiles/此处用需要手机号的那个应用包名填充/此处用需要的手机号数量填充
- **请求方法：** HTTP GET方法

### 1.2.1 请求示例：
```
http://服务器IP:端口号/mobiles/com.yy.huanju/8
```

### 1.2.2 返回结果
#### A. 服务器内部运行出错时, 同上文；
#### B. 服务器暂时数据给客户端时，
	{
		"code": 208,
		"desc": "Sorry, there is no data that you need currently! Please try again later~",
	}
#### C. 客户端传递的某个参数的数据类型或数值不符合要求时，同上文
#### D. 接收请求并返回客户端所需数据
	{
		"code": 200,
		"data": 对象数组, 例如：[{"mobile": "17342060825", "cid": "world", "imsi": "ghsdv", "imei": "12356416", 	"iccid": "ws5dfg45s", "ip": "6.66.666.686"}, {"mobile": "18250794865", "cid": "hello", "imsi": "zjgcghsdv", "imei": "56416", "iccid": "s5d45s", "ip": "8.88.888.886"}]
	}


	
## 1.3 上传验证码信息到本服务器
- **接口地址：** /vcodes
- **请求方法：** HTTP POST方法，采用JSON格式

### 1.3.1 Body中的请求参数
  
参数名称						|类型		|出现要求	|描述  
:----						|:---		|:------	|:---	
data				|对象数组	|不能为空！		|&nbsp;
&emsp;mobile		|string		|不能为空！		|手机号
&emsp;vcode			|string		|不能为空！  	|对应所传手机号的那个验证码
&emsp;appname		|string		|   --  	    |需要验证码的那个应用包名
&emsp;chanlid		|string		|	--			|分配验证码的渠道号

Body中参数内容示例：

```
{
	"data":[{
		"mobile":"18250794865",
		"vcode":"201708",
		"appname":"com.yy.huanju",
		"chanlid":"hello"
	},{
		"mobile":"17342060825",
		"vcode":"171331",
		"appname":"com.yy.huanju",
		"chanlid":"world"
	}
		]
}

```

### 1.3.2 返回结果
#### A. 服务器内部运行出错时，同上文；
#### B. 客户端传递的某个参数的数据类型或数值不符合要求时，同上文；
#### C. 成功接收并处理了客户端的请求；
	{
		"code": 200,
		"msg":  "Save successfully~",
	}


	
## 1.4 从本服务器获取验证码信息
- **接口地址：** /vcodes/此处用需要验证码的那个应用包名填充（要与请求手机号的应用包名一致）
- **请求方法：** HTTP GET方法

### 1.4.1 请求示例：
```
http://服务器IP:端口号/vcodes/com.yy.huanju
```

### 1.4.2 返回结果
#### A. 服务器内部运行出错时, 同上文；
#### B. 服务器暂时数据给客户端时，同上文；	
#### C. 成功接收请求并返回客户端所需数据时，
	{
		"code": 200,
		"data": 对象数组,例如：[{"mobile": "17342060825", "vcode": "171331", "appname": "com.yy.huanju", "chanlid": "world"},
        {"mobile": "18250794865", "vcode": "201708", "appname": "com.yy.huanju","chanlid": "hello"}]
	}

	
	
## 1.5 控制客户端上传数据的速率
- **接口地址：** /setsize
- **接口参数：** newsize, 整数
- **请求方法：** HTTP PUT方法

### 1.5.1 请求示例：
```
http://服务器IP:端口号/setsize?newsize=2
```

### 1.5.2 返回结果
#### A. 客户端传递的某个参数的数据类型或数值不符合要求时，同上文；
#### B. 成功接收并处理请求时，
	{
		"code": 200,
		"data": "Update Successfully!"
	}





