

## 本程序每次可为客户端提供一个不重复的用户名、密码信息，兼具简易的打点、预警机制

### 1 获取账户信息：
- **测试时的接口地址：** http://192.168.1.169:10001/account
- **请求方法：** HTTP GET方法

#### 1.2 返回结果为以下情况之一：
####  A. 成功接收并处理了用户请求，同时服务器有充足数据时（http头部状态码为200）：
    {
		"username": 用户名
		"password": 密码
	}
#### B. 服务器暂时没有客户端需要的数据时（http头部状态码为410）：
    {
      "desc": "Sorry, there is no data that you need currently! Please try again later~"
    }


### 2 上传帐号使用结果：
- **测试时的接口地址：** http://192.168.1.169:10001/result
- **请求方法：** HTTP POST方法

#### 2.1 Body中需要携带的参数

参数名称	    |类型	    |出现要求	     |描述  
:----		|:---		|:------	 |:---
data		|对象字符串 	|不能为空	     |以下各字段均为该对象中的属性
&emsp;&emsp;result		|int 		|不能为空	     |使用结果，0为失败，1为成功
&emsp;&emsp;account		|string		|不能为空	     |用户名
&emsp;&emsp;password 	|string		|最好不为空	 |密码
&emsp;&emsp;imsi		|string		|最好不为空	 |使用者的SIM卡中的用户标识号
&emsp;&emsp;imei		|string		|最好不为空	 |使用者的设备ID
&emsp;&emsp;installid	|string		|最好不为空	 |安装ID
&emsp;&emsp;ua		    |string		|最好不为空	 |UA
&emsp;&emsp;failreason	|string		|最好不为空	 |错误原因

#### 2.2 返回结果为以下情况之一：
####  A. 成功接收并处理了用户请求，同时服务器有充足数据时（http头部状态码为200）：
    {
		"msg": "Accept successfully."
	}
#### B. 客户端的请求不符合基本要求时（http头部状态码为400）：
    {
      "desc": "The data that you have submitted is not all correct !"
    }
