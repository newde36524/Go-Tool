package ali

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttP2PClient struct {
	opt          *MqttClientOption
	clientID     string
	userName     string
	passWord     string
	selfP2pTopic string

	client mqtt.Client
}

type MqttClientOption struct {
	instanceID string
	brokerURL  string
	accessKey  string
	secretKey  string
	topic      string
	deviceID   string
	groupID    string

	qos      byte
	retained bool
}

func (opt *MqttClientOption) SetInstanceID(instanceID string) *MqttClientOption {
	opt.instanceID = instanceID
	return opt
}
func (opt *MqttClientOption) SetBrokerURL(brokerURL string) *MqttClientOption {
	opt.brokerURL = brokerURL
	return opt
}
func (opt *MqttClientOption) SetAccessKey(accessKey string) *MqttClientOption {
	opt.accessKey = accessKey
	return opt
}
func (opt *MqttClientOption) SetSecretKey(secretKey string) *MqttClientOption {
	opt.secretKey = secretKey
	return opt
}
func (opt *MqttClientOption) SetTopic(topic string) *MqttClientOption {
	opt.topic = topic
	return opt
}
func (opt *MqttClientOption) SetDeviceID(deviceID string) *MqttClientOption {
	opt.deviceID = deviceID
	return opt
}
func (opt *MqttClientOption) SetGroupID(groupID string) *MqttClientOption {
	opt.groupID = groupID
	return opt
}

func (opt *MqttClientOption) SetQos(qos byte) *MqttClientOption {
	opt.qos = qos
	return opt
}
func (opt *MqttClientOption) SetRetained(retained bool) *MqttClientOption {
	opt.retained = retained
	return opt
}
func NewMqttP2PClient(optFunc func(*MqttClientOption)) *MqttP2PClient {
	opt := &MqttClientOption{
		qos:      0x01,
		retained: false,
	}
	optFunc(opt)
	var (
		clientID     = opt.groupID + "@@@" + opt.deviceID
		userName     = "Signature|" + opt.accessKey + "|" + opt.instanceID
		passWord     = HMACSHA1(opt.secretKey, clientID)
		selfP2pTopic = opt.topic + "/p2p/" + clientID
	)
	return &MqttP2PClient{
		opt:          opt,
		clientID:     clientID,
		userName:     userName,
		passWord:     passWord,
		selfP2pTopic: selfP2pTopic,
	}
}

func (mqttP2PClient *MqttP2PClient) Regist(fn func(b []byte)) error {
	opts := mqtt.NewClientOptions().
		AddBroker(mqttP2PClient.opt.brokerURL).
		SetClientID(mqttP2PClient.clientID).
		SetUsername(mqttP2PClient.userName).
		SetPassword(mqttP2PClient.passWord).
		SetAutoReconnect(true). //自动重连
		SetDefaultPublishHandler(mqtt.MessageHandler(func(c mqtt.Client, msg mqtt.Message) {
			fn(msg.Payload())
		}))
	mqttP2PClient.client = mqtt.NewClient(opts)
	if token := mqttP2PClient.client.Connect(); token.Wait() && token.Error() != nil {
		return errors.New("未成功建立连接")
	}
	if mqttP2PClient.client.IsConnected() && mqttP2PClient.client.IsConnectionOpen() {
		return errors.New("连接建立成功")
	} else {
		return errors.New("连接未建立")
	}
}

func (mqttP2PClient *MqttP2PClient) Send(data []byte, deviceID string) (bool, error) {
	opt := mqttP2PClient.opt
	token := mqttP2PClient.client.Publish(opt.topic+"/p2p/"+opt.groupID+"@@@"+deviceID, opt.qos, opt.retained, data)
	return token.Wait(), token.Error()
}

func HMACSHA1(key, dataToSign string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(dataToSign))
	result := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(result)
}

func (mqttP2PClient *MqttP2PClient) Close() {
	mqttP2PClient.client.Disconnect(0)
}
