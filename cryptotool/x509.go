package cryptotool

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
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

// //your_private_key和your_public_key都是上述《go加载rsa公、私钥》部分创建的结构体类型
// func Sign(src []byte, hash crypto.Hash) ([]byte, error) {
// 	h := hash.New()
// 	h.Write(src)
// 	hashed := h.Sum(nil)
// 	return rsa.SignPKCS1v15(rand.Reader, your_private_key, hash, hashed)
// }

// func Verify(src []byte, sign []byte, hash crypto.Hash) error {
// 	h := hash.New()
// 	h.Write(src)
// 	hashed := h.Sum(nil)
// 	return rsa.VerifyPKCS1v15(your_public_key, hash, hashoed, sign)
// }
