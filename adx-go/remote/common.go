package remote

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func urlQueryTrans(inPut map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, val := range inPut {
		for _, va := range val {
			result[key] = va
		}
	}
	return result
}
func showSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func showErro(c *gin.Context, code string) {
	erroMsg := map[string]string{
		"-201": "未找到广告标示符",
		"-301": "缺失设备参数",
		"-302": "投放已暂停",
		"-303": "投放时间错误",
		"-304": "投放已超量",
		"-316": "参数无法解析",
		"-401": "IP地址无效",
		"-501": "mac地址无效",
	}
	msg, ok := erroMsg[code]
	if !ok {
		code = "-316"
		msg, ok = erroMsg[code]
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func showNoFound(c *gin.Context, code string) {
	c.String(http.StatusNotFound, code, "")
}

func base64_encode(input string) string {
	if len(input) > 0 {
		inputBy := []byte(input)
		// 演示base64编码
		encodeString := base64.StdEncoding.EncodeToString(inputBy)
		return encodeString

	}
	return ""
}

func base64_decode(encodeString string) string {
	if len(encodeString) > 0 {
		// 对上面的编码结果进行base64解码
		decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
		if err != nil {
			log.Println("base64解码失败", err)
		}
		return string(decodeBytes)

	}
	return ""
}

func calcmd5(data string) string {
	databy := []byte(data)
	hash := md5.New()
	hash.Write(databy)
	cipherText := hash.Sum(nil)
	return hex.EncodeToString(cipherText)
}

func checkIp(ipstr string) bool {
	var b bool
	b, _ = regexp.MatchString("^([1-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])(\\.([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])){3}$", ipstr)
	return b
}

func checkrep_url(urlstr string) bool {
	if len(urlstr) == 0 {
		return false
	}
	matched, err := regexp.MatchString("\\{.*\\}", urlstr)
	if err != nil || !matched {
		log.Println("func checkrep_url erro", urlstr, err)
	}
	return matched
}
