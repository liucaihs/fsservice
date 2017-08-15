
本程序主要提供设备信息的查阅服务，并保证在读完数据库中所有设备信息之前每次获取的设备信息是不同的。

（1）获取设备信息需要发GET请求，url格式为：
        http://服务器IP:8091/deviceinfo/:applicationame
    其中，“:applicationame”需要用请求信息的应用程序名代替。
(2)返回内容：
    A. 服务器处理成功并返回一条设备信息时（http头部的状态码为200），结构为：
        {
           "data": 数据结构为DeviceInfo
        }
    B. 服务器程序运行出错（http头部的状态码为500），返回结构为：
       {
         "desc": "The server is busy, please try again later."
       }
    C. 暂无新数据给客户端时（http头部的状态码为208），返回结构为：
       {
         "desc": "Nowdays, you have obtained whole data that we have collected. If you want to get them again, please try later~"
       }
    D. 客户端请求的路由不正确（http头部的状态码为404），返回结构为：“404 page not found”

附： DeviceInfo的结构为：
      {   
          字段名                 说明或者原始数据来源处
          imei                  设备ID
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
        	gsmlac                gsm基站信息
        	gsmcid                gsm基站信息
        	cdmalatitude          cdma基站信息
        	Cdmalongitude         cdma基站信息
        	cdmabid               cdma基站信息
        	cdmasid               getSystemId(), cdma基站信息
        	cdmanid               getNetworkId(), cdma基站信息
        	wh                    表示屏幕分辨率
          pkginfos              表示当前设备已安装的APP列表
    }
