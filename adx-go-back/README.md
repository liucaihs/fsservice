注：本服务为中转控制系统！！

# 1. 与客户端(手机端)通讯数据格式定义

## 1.1 客户端注册（Signin）
- **接口地址：**  ws:// IP :9090/rs
- **请求方法：**  WebSoket send，采用JSON格式

### 1.1.1 消息内容格式
  
参数名称						|类型		|出现要求	|描述  
:----						|:---		|:------	|:---
sid			        |int		|	--			|消息编号（自增）
cmd			        |string		|	--			|消息标示	
data				|对象数组	|不能为空！		|&nbsp;
&emsp;phone		    |string		|不能为空！		|手机号
&emsp;iccid		    |string		|不能为空！		|iccid
&emsp;imei			|string		|	--			|imei
&emsp;imsi			|string		|	--			|imsi
 

接受消息内容示例：

```
{
    "sid":123,
    "cmd":"Signin",
    "data":{
        "iccid":"2D906981-874E-4052-91A4-755B8A081B54",
        "imei":"13450211224",
		"imsi":"",
		"phone":"13450211224",
    }
} 

返回消息示例：
{
    "sid":0,
    "cmd":"Signin",
    "data":{
         "status":true
    }
} 
``` 

## 1.2 获取上传设备信息接口信息（Device）
- **接口地址：**  ws:// IP :9090/rs
- **请求方法：**  WebSoket send，采用JSON格式

### 1.2.1 消息内容格式
  
参数名称				|类型		            |出现要求	        |描述  
:----				|:---		            |:------	    |:---
sid			        |int		            |	--			|消息编号（自增）
cmd			        |string("Device")		|	--			|消息标示	
data				|string	                |   --		    |&nbsp;
 
 

接受消息内容示例：

```
{
    "sid":123,
    "cmd":"Device",
    "data":""
} 
```
### 1.2.2 返回消息内容格式
  
参数名称						|类型		|出现要求	|描述  
:----						|:---		|:------	|:---
sid			        |int		|	--			|消息编号（自增）
cmd		 （"Device"）|string		|	--			|消息标示	
data				|对象数组	|不能为空！		|&nbsp;
&emsp;url		    |string		|不能为空！		|上传设备信息接口地址
 
 

消息内容示例：

```
{
    "sid":123,
    "cmd":"Device",
    "data":"{
		"url":"http://www.abc.com"
	}"
} 
```


## 1.3 通知获取验证码消息数据格式（Vcode）
- **接口地址：**  ws:// IP :9090/rs
- **请求方法：**  WebSoket send，采用JSON格式

### 1.3.1 消息内容格式
  
参数名称				|类型		            |出现要求	        |描述  
:----				|:---		            |:------	    |:---
sid			        |int		            |	--			|消息编号（自增）
cmd			        |string("Vcode")		|	--			|消息标示	
data				|string	                |   --		    |&nbsp;
 
 

发送消息内容示例：

```
{
    "sid":0,
    "cmd":"Vcode",
    "data":"{}"
} 
```
### 1.3.2 返回消息内容格式 
消息内容示例：

```
{
   
} 
```