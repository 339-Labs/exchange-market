package signer

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/339-Labs/exchange-market/common"
	"strings"
)

type Signer struct {
	secretKey []byte
}

func (p *Signer) Init(key string) *Signer {
	p.secretKey = []byte(key)
	return p
}

func (p *Signer) Sign(params map[string]string) string {
	str := common.SortParams(params)
	var payload strings.Builder
	payload.WriteString(str)
	hash := hmac.New(sha256.New, p.secretKey)
	hash.Write([]byte(payload.String()))
	result := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return result
}

func (p *Signer) ByBitSign(apiKey string, timestamp string, body string) string {
	var payload strings.Builder
	payload.WriteString(apiKey)
	payload.WriteString(timestamp)
	if body != "" && body != "?" {
		payload.WriteString(body)
	}
	hash := hmac.New(sha256.New, p.secretKey)
	hash.Write([]byte(payload.String()))
	result := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return result
}

func (p *Signer) BnSign(apiKey string, timestamp string) string {
	var payload strings.Builder
	payload.WriteString(apiKey)
	payload.WriteString(timestamp)
	hash := hmac.New(sha256.New, p.secretKey)
	hash.Write([]byte(payload.String()))
	result := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return result
}

func (p *Signer) ParamsSign(method string, requestPath string, body string, timesStamp string) string {
	var payload strings.Builder
	payload.WriteString(timesStamp)
	payload.WriteString(method)
	payload.WriteString(requestPath)
	if body != "" && body != "?" {
		payload.WriteString(body)
	}
	hash := hmac.New(sha256.New, p.secretKey)
	hash.Write([]byte(payload.String()))
	result := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return result
}

func (p *Signer) SignByRSA(method string, requestPath string, body string, timesStamp string) string {
	var payload strings.Builder
	payload.WriteString(timesStamp)
	payload.WriteString(method)
	payload.WriteString(requestPath)
	if body != "" && body != "?" {
		payload.WriteString(body)
	}

	sign, _ := RSASign([]byte(payload.String()), p.secretKey, crypto.SHA256)
	result := base64.StdEncoding.EncodeToString(sign)
	return result
}

func RSASign(src []byte, priKey []byte, hash crypto.Hash) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("key is invalid format")
	}

	var pkixPrivateKey interface{}
	var err error
	if block.Type == "RSA PRIVATE KEY" {
		pkixPrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if block.Type == "PRIVATE KEY" {
		pkixPrivateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	h := hash.New()
	_, err = h.Write(src)
	if err != nil {
		return nil, err
	}

	bytes := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, pkixPrivateKey.(*rsa.PrivateKey), hash, bytes)
	if err != nil {
		return nil, err
	}

	return sign, nil
}
