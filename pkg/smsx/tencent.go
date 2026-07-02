package smsx

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Config 描述腾讯云短信发送接口所需配置。
type Config struct {
	SecretID   string
	SecretKey  string
	SDKAppID   string
	SignName   string
	TemplateID string
	Region     string
	Endpoint   string
	HTTPClient *http.Client
}

// TencentSender 调用腾讯云短信发送接口。
type TencentSender struct {
	cfg        Config
	httpClient *http.Client
}

type sendSMSRequest struct {
	PhoneNumberSet   []string `json:"PhoneNumberSet"`
	SmsSdkAppId      string   `json:"SmsSdkAppId"`
	SignName         string   `json:"SignName"`
	TemplateId       string   `json:"TemplateId"`
	TemplateParamSet []string `json:"TemplateParamSet"`
}

type sendSMSResponse struct {
	Response struct {
		SendStatusSet []struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"SendStatusSet"`
		Error *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
	} `json:"Response"`
}

// NewTencentSender 创建腾讯云短信发送器。
func NewTencentSender(cfg Config) *TencentSender {
	if strings.TrimSpace(cfg.Endpoint) == "" {
		cfg.Endpoint = "https://sms.tencentcloudapi.com"
	}
	if strings.TrimSpace(cfg.Region) == "" {
		cfg.Region = "ap-guangzhou"
	}
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &TencentSender{cfg: cfg, httpClient: client}
}

// SendRegisterCode 发送注册验证码短信。
func (s *TencentSender) SendRegisterCode(ctx context.Context, phone string, code string, ttl time.Duration) error {
	if err := s.validate(); err != nil {
		return err
	}
	payload := sendSMSRequest{
		PhoneNumberSet:   []string{normalizeMainlandPhone(phone)},
		SmsSdkAppId:      s.cfg.SDKAppID,
		SignName:         s.cfg.SignName,
		TemplateId:       s.cfg.TemplateID,
		TemplateParamSet: []string{code, strconv.Itoa(int(ttl.Minutes()))},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("编码短信请求失败: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.Endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("创建短信请求失败: %w", err)
	}
	s.sign(req, body, time.Now())
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("调用短信接口失败: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取短信响应失败: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("短信接口状态异常: %d", resp.StatusCode)
	}
	var result sendSMSResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析短信响应失败: %w", err)
	}
	if result.Response.Error != nil {
		return fmt.Errorf("短信接口返回错误: %s", result.Response.Error.Message)
	}
	if len(result.Response.SendStatusSet) == 0 {
		return errors.New("短信接口未返回发送状态")
	}
	if result.Response.SendStatusSet[0].Code != "Ok" {
		return fmt.Errorf("短信发送失败: %s", result.Response.SendStatusSet[0].Message)
	}
	return nil
}

func (s *TencentSender) validate() error {
	if strings.TrimSpace(s.cfg.SecretID) == "" || strings.TrimSpace(s.cfg.SecretKey) == "" {
		return errors.New("腾讯云短信密钥未配置")
	}
	if strings.TrimSpace(s.cfg.SDKAppID) == "" || strings.TrimSpace(s.cfg.SignName) == "" || strings.TrimSpace(s.cfg.TemplateID) == "" {
		return errors.New("腾讯云短信应用、签名或模板未配置")
	}
	return nil
}

func (s *TencentSender) sign(req *http.Request, body []byte, now time.Time) {
	const service = "sms"
	date := now.UTC().Format("2006-01-02")
	timestamp := strconv.FormatInt(now.Unix(), 10)
	host := req.URL.Host
	if parsed, err := url.Parse(s.cfg.Endpoint); err == nil && parsed.Host != "" {
		host = parsed.Host
	}
	canonicalHeaders := "content-type:application/json; charset=utf-8\nhost:" + host + "\n"
	signedHeaders := "content-type;host"
	hashedRequestPayload := sha256Hex(body)
	canonicalRequest := "POST\n/\n\n" + canonicalHeaders + "\n" + signedHeaders + "\n" + hashedRequestPayload
	credentialScope := date + "/" + service + "/tc3_request"
	stringToSign := "TC3-HMAC-SHA256\n" + timestamp + "\n" + credentialScope + "\n" + sha256Hex([]byte(canonicalRequest))
	secretDate := hmacSHA256([]byte("TC3"+s.cfg.SecretKey), date)
	secretService := hmacSHA256(secretDate, service)
	secretSigning := hmacSHA256(secretService, "tc3_request")
	signature := hex.EncodeToString(hmacSHA256(secretSigning, stringToSign))
	authorization := "TC3-HMAC-SHA256 Credential=" + s.cfg.SecretID + "/" + credentialScope + ", SignedHeaders=" + signedHeaders + ", Signature=" + signature

	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Host", host)
	req.Header.Set("X-TC-Action", "SendSms")
	req.Header.Set("X-TC-Version", "2021-01-11")
	req.Header.Set("X-TC-Timestamp", timestamp)
	req.Header.Set("X-TC-Region", s.cfg.Region)
}

func normalizeMainlandPhone(phone string) string {
	trimmed := strings.TrimSpace(phone)
	if strings.HasPrefix(trimmed, "+") {
		return trimmed
	}
	return "+86" + trimmed
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(data))
	return mac.Sum(nil)
}
