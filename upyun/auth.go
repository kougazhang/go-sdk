package upyun

import (
	"fmt"
	"sort"
	"strings"
)

type RESTAuthConfig struct {
	Method    string
	Uri       string
	DateStr   string
	LengthStr string
}

type PurgeAuthConfig struct {
	PurgeList string
	DateStr   string
}

type UnifiedAuthConfig struct {
	Method     string
	Uri        string
	DateStr    string
	Policy     string
	ContentMD5 string
}

func (up *UpYun) MakeRESTAuth(config *RESTAuthConfig) string {
	sign := []string{
		config.Method,
		config.Uri,
		config.DateStr,
		config.LengthStr,
		up.Password,
	}
	return "UpYun " + up.Operator + ":" + md5Str(strings.Join(sign, "&"))
}

func (up *UpYun) MakePurgeAuth(config *PurgeAuthConfig) string {
	sign := []string{
		config.PurgeList,
		up.Bucket,
		config.DateStr,
		up.Password,
	}
	return "UpYun " + up.Bucket + ":" + up.Operator + ":" + md5Str(strings.Join(sign, "&"))
}

func (up *UpYun) MakeFormAuth(policy string) string {
	return md5Str(base64ToStr([]byte(policy)) + "&" + up.Secret)
}

func (up *UpYun) MakeProcessAuth(kwargs map[string]string) string {
	keys := []string{}
	for k := range kwargs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	auth := ""
	for _, k := range keys {
		auth += k + kwargs[k]
	}
	return fmt.Sprintf("UpYun %s:%s", up.Operator, md5Str(up.Operator+auth+up.Password))
}

func (up *UpYun) MakeUnifiedAuth(config *UnifiedAuthConfig) string {
	sign := []string{
		config.Method,
		config.Uri,
		config.DateStr,
		config.Policy,
		config.ContentMD5,
	}
	signNoEmpty := []string{}
	for _, v := range sign {
		if v != "" {
			signNoEmpty = append(signNoEmpty, v)
		}
	}
	signStr := base64ToStr(hmacSha1(up.Password, []byte(strings.Join(signNoEmpty, "&"))))
	return "UpYun " + up.Operator + ":" + signStr
}
