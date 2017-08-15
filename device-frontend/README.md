
本程序主要提供设备信息的存储服务。

（1）上传设备信息需要发Post请求，url格式为：
        http://服务器IP:8090/deviceinfo
    采用"application/json"方式传递Body内容，其数据结构为：
      {   
          字段名                 说明或者客户端可能的获取方法
          imei                  设备ID，不能为空！!！
        	imsi                  getSubscriberId()
        	phone                 手机号
        	simserial             getSimSerialNumber()
        	simcountryiso         getSimCountryIso()
        	simoperator           getSimOperator()
        	simoperatorname       getSimOperatorName()
        	networkcountryiso     getNetworkCountryIso()
        	networkoperator       getNetworkOperator()
        	networkoperatorname   getNetworkOperatorName()
        	wifimac               getMacAddress()或getHardwareAddress()
        	bssid                 getBSSID()
        	ssid                  getSSID()
        	bluemac               蓝牙的地址信息
        	bluename              蓝牙的名称
        	model                 Build.MODEL
        	manufacturer          Build.MANUFACTURER
        	brand                 Build.BRAND
        	hardware              Build.HARDWARE
        	board                 Build.BOARD
        	serial                Build.SERIAL
        	device                Build.DEVICE
        	buildid               Build.ID
        	product               Build.PRODUCT
        	display               Build.DISPLAY
        	fingerprint           Build.FINGERPRINT
        	nowrelease            Build.RELEASE
        	sdk                   Build.SDK
        	radioversion          getRadioVersion()
        	androidid             Secure.getString()
        	mnc                 	运营商的mnc信息
        	mcc                   运营商的mcc信息
        	latitude              纬度
        	longitude             经度
        	gsmlac                getLac(), gsm基站信息
        	gsmcid                getCid(), gsm基站信息
        	cdmalatitude          getBaseStationLatitude(), cdma基站信息
        	Cdmalongitude         getBaseStationLongitude(), cdma基站信息
        	cdmabid               getBaseStationId(), cdma基站信息
        	cdmasid               getSystemId(), cdma基站信息
        	cdmanid               getNetworkId(), cdma基站信息
        	wh                    表示屏幕分辨率，以“width`heigth”格式拼接，比如768`1184, 获取方
                                   法可能为display.getWidth()、display.getHeight()
          pkginfos              表示当前设备已安装的APP列表
    }
注：各字段的值必须为字符串类型。

(2)返回内容：
  A. 保存成功时（http头部的状态码为200），结构为：
    {
       "msg": "Save successfully~"
    }
  B. 服务器程序运行出错（http头部的状态码为500），返回结构为：
   {
     "desc": "The server is busy, please try again later.",
   }

  C. 客户端所发请求中某些字段的数据类型不正确或者该字段未出现（http头部的状态码为400），返回结构为：
   {
     "desc": "The data that you have submitted is not all correct.",
   }

  D. 客户端提交了重复的设备信息时（http头部的状态码为406），返回结构为：
   {
     "desc": "The data that you have submitted seems to be a copy of former data.",
   }
  E. 客户端请求的路由不正确（http头部的状态码为404），返回结构为：“404 page not found”
