package remote

import (
	"encoding/json"
	"log"
	"os"
)

type BaseMessage struct {
	Sid  int             `json:"sid"`
	Cmd  string          `json:"cmd"`
	Data json.RawMessage `json:"data"`
}

func GetEnvDef(name string, def_str string) string {
	value := os.Getenv(name)
	if len(value) > 0 {
		return value
	}
	return def_str
}

func (c *Client) sendMessage(m *BaseMessage) {
	buffer, err := json.Marshal(m)
	if err != nil {
		log.Println("marshal err:", err)
		return
	}
	c.send(buffer)
}

func (c *Client) sendBaseMsg(cmd string, data map[string]interface{}) {
	if len(cmd) == 0 {
		log.Println("client sendBaseMsg: cmd cannot empty")
		return
	}
	var sendmsg BaseMessage
	datastr, err := json.Marshal(data)
	if err != nil {
		log.Println("client sendBaseMsg Marshal: ", err.Error())
	}
	sendmsg.Cmd = cmd
	sendmsg.Data = datastr
	c.sendMessage(&sendmsg)

}

func Json2map(req []byte) (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal(req, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type CHandler func(c *Client, m *BaseMessage)

var Router map[string]CHandler = map[string]CHandler{
	"Signin": onSignin,
	"Device": onDevice,
	"Test":   onTest,
}

func handleMessage(c *Client, buffer []byte) {

	var message BaseMessage
	err := json.Unmarshal(buffer, &message)
	if err != nil {
		log.Println("Unmarsha1 err:", err)
		return
	}
	handle := Router[message.Cmd]
	if handle == nil {
		log.Println("message content:", string(buffer), "client id:", c.id)
		log.Println("not exist message type:", message.Cmd)
		c.sendBaseMsg("DEBUG", map[string]interface{}{"msg": "server receive" + string(buffer)})
		return
	}
	handle(c, &message)
}

func onSignin(c *Client, message *BaseMessage) {
	inputData, err := Json2map(message.Data)
	if err != nil {
		log.Println("onSignin Unmarsha1 err:", err)
	}
	onlineMsg := inputData
	onlineMsg["type"] = "online"
	datastr, err := json.Marshal(onlineMsg)
	if err != nil {
		log.Println("onSignin Marshal err:", err)
	}
	c.hub.Srv.mq.PutMessage("onoff", "", string(datastr))
	log.Println(inputData)
	c.hub.Srv.model.SaveClientInfo(c.id, inputData["phone"], inputData["imei"], inputData["iccid"], inputData["imsi"], c.ip)

	c.hub.Srv.model.ShowAllRec()
}

func onDevice(c *Client, message *BaseMessage) {
	uploadDev(c)
}

func onTest(c *Client, message *BaseMessage) {
	inputData, err := Json2map(message.Data)
	if err != nil {
		log.Println("Unmarsha1 err:", err)
	}

	log.Println(inputData)
	cliExt := c.hub.Srv.model.GetClintInfo(inputData["phone"].(string), "")
	log.Println("cliExt info:", cliExt)
	if sendid, ok := cliExt["id"]; ok {
		// ...
		selClient, cok := c.hub.clients[uint32(sendid.(int))]
		log.Println("selClient info:", selClient)
		if cok {
			selClient.sendBaseMsg("SEND2ONE", map[string]interface{}{"msg": "server receive" + inputData["phone"].(string)})
		}
	}

	return
}

func uploadDev(c *Client) {
	var device = make(map[string]interface{})
	device["url"] = GetEnvDef("DEVICE_URL", "http://192.168.1.169:8090") ///deviceinfo
	c.sendBaseMsg("Device", device)

}
