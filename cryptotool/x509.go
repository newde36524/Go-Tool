package cryptotool

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/issue9/logs"
)

//CreateX509Cer 创建Cer证书数据
//@cerFilePath 证书文件路径
//@key 验证私钥
//@activeTime 证书生效时间
//@duration 证书有效期
//@commonName 证书公开名称
//@organization 证书来源
//@dNSNames 证书解析域名
//@subjectKeyId 证书编号
func CreateX509Cer(cerFilePath string, key *rsa.PrivateKey, activeTime time.Time, duration time.Duration, commonName string, organization, dNSNames []string, subjectKeyId []byte) {
	random := rand.Reader
	now := activeTime
	then := now.Add(duration)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Signature:    x509.MarshalPKCS1PrivateKey(key),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: organization,
		},
		NotBefore:             now,
		NotAfter:              then,
		SubjectKeyId:          subjectKeyId,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              dNSNames,
	}
	derBytes, err := x509.CreateCertificate(random, &template, &template, &key.PublicKey, &key)
	if err != nil {
		logs.Error(err)
	}
	certCerFile, err := os.Create(cerFilePath)
	if err != nil {
		logs.Error(err)
	}
	certCerFile.Write(derBytes)
	certCerFile.Close()

}

//ReadX509Cer .
func ReadX509Cer(cerFilePath string) (*x509.Certificate, error) {
	derBytes, err := ioutil.ReadFile(cerFilePath)
	if err != nil {
	}
	cer, err := x509.ParseCertificate(derBytes)
	return cer, err
}

func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func CreateX509Cer2() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	subject := pkix.Name{                         //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		Organization:       []string{"Manning Publications Co."},
		OrganizationalUnit: []string{"Books"},
		CommonName:         "Go Web Programming",
	}
	template := x509.Certificate{
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, //KeyUsage 与 ExtKeyUsage 用来表明该证书是用来做服务器认证的
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // 密钥扩展用途的序列
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048) //生成一对具有指定字位数的RSA密钥

	//CreateCertificate基于模板创建一个新的证书
	//第二个第三个参数相同，则证书是自签名的
	//返回的切片是DER编码的证书
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk) //DER 格式
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICAET", Bytes: derBytes})
	certOut.Close()
	keyOut, _ := os.Create("key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}
