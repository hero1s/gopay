package gopay

import (
	"encoding/json"
	"fmt"
	"time"
)

// toney 扩展接口
// --------------------------------------------------------------------------------
type CertifyBizCode string

const (
	K_CERTIFY_BIZ_CODE_FACE            CertifyBizCode = "FACE"            // 多因子人脸认证
	K_CERTIFY_BIZ_CODE_CERT_PHOTO      CertifyBizCode = "CERT_PHOTO"      // 多因子证照认证
	K_CERTIFY_BIZ_CODE_CERT_PHOTO_FACE CertifyBizCode = "CERT_PHOTO_FACE" // 多因子证照和人脸认证
	K_CERTIFY_BIZ_CODE_SMART_FACE      CertifyBizCode = "SMART_FACE"      // 多因子快捷认证
)

// https://docs.open.alipay.com/api_2/alipay.user.certify.open.initialize
type UserCertifyOpenInitialize struct {
	AppAuthToken        string         `json:"-"`                               // 可选
	OuterOrderNo        string         `json:"outer_order_no"`                  // 必选  商户请求的唯一标识，商户要保证其唯一性，值为32位长度的字母数字组合。建议：前面几位字符是商户自定义的简称，中间可以使用一段时间，后段可以使用一个随机或递增序列
	BizCode             CertifyBizCode `json:"biz_code"`                        // 必选 认证场景码。入参支持的认证场景码和商户签约的认证场景相关，取值如下: FACE：多因子人脸认证 CERT_PHOTO：多因子证照认证 CERT_PHOTO_FACE ：多因子证照和人脸认证 SMART_FACE：多因子快捷认证
	IdentityParam       IdentityParam  `json:"identity_param"`                  // 必选
	MerchantConfig      MerchantConfig `json:"merchant_config"`                 // 必选 商户个性化配置，格式为json，详细支持的字段说明为： return_url：需要回跳的目标地址，必填，一般指定为商户业务页面
	FaceContrastPicture string         `json:"face_contrast_picture,omitempty"` // 可选 自定义人脸比对图片的base64编码格式的string字符串
}

type IdentityParam struct {
	IdentityType string `json:"identity_type"` // 身份信息参数类型，必填，必须传入CERT_INFO
	CertType     string `json:"cert_type"`     // 证件类型，必填，当前支持身份证，必须传入IDENTITY_CARD
	CertName     string `json:"cert_name"`     // 真实姓名，必填，填写需要验证的真实姓名
	CertNo       string `json:"cert_no"`       // 证件号码，必填，填写需要验证的证件号码
}

type MerchantConfig struct {
	ReturnURL string `json:"return_url"`
}

type UserCertifyOpenInitializeRsp struct {
	Content struct {
		Code      string `json:"code"`
		Msg       string `json:"msg"`
		SubCode   string `json:"sub_code"`
		SubMsg    string `json:"sub_msg"`
		CertifyId string `json:"certify_id"`
	} `json:"alipay_user_certify_open_initialize_response"`
	Sign string `json:"sign"`
}

// --------------------------------------------------------------------------------
// https://docs.open.alipay.com/api_2/alipay.user.certify.open.certify
type UserCertifyOpenCertify struct {
	AppAuthToken string `json:"-"`          // 可选
	CertifyId    string `json:"certify_id"` // 必选 本次申请操作的唯一标识，由开放认证初始化接口调用后生成，后续的操作都需要用到
}
type UserCertifyOpenCertifyRsp struct {
	Content struct {
		Code    string `json:"code"`
		Msg     string `json:"msg"`
		SubCode string `json:"sub_code"`
		SubMsg  string `json:"sub_msg"`
	} `json:"alipay_user_certify_open_certify_response"`
	Sign string `json:"sign"`
}

// --------------------------------------------------------------------------------
// https://docs.open.alipay.com/api_2/alipay.user.certify.open.query/
type UserCertifyOpenQuery struct {
	AppAuthToken string `json:"-"`          // 可选
	CertifyId    string `json:"certify_id"` // 必选 本次申请操作的唯一标识，由开放认证初始化接口调用后生成，后续的操作都需要用到
}
type UserCertifyOpenQueryRsp struct {
	Content struct {
		Code         string `json:"code"`
		Msg          string `json:"msg"`
		SubCode      string `json:"sub_code"`
		SubMsg       string `json:"sub_msg"`
		Passed       string `json:"passed"`
		IdentityInfo string `json:"identity_info"`
		MaterialInfo string `json:"material_info"`
	} `json:"alipay_user_certify_open_query_response"`
	Sign string `json:"sign"`
}

// UserCertifyOpenInitialize 身份认证初始化服务 https://docs.open.alipay.com/api_2/alipay.user.certify.open.initialize
func (a *AliPayClient) UserCertifyOpenInitialize(body BodyMap) (resp UserCertifyOpenInitializeRsp, err error) {
	var bs []byte
	if bs, err = a.doAliPay(body, "alipay.user.certify.open.initialize"); err != nil {
		return
	}
	err = json.Unmarshal(bs, &resp)
	return
}

// UserCertifyOpenCertify 身份认证开始认证 https://docs.open.alipay.com/api_2/alipay.user.certify.open.certify
func (a *AliPayClient) UserCertifyOpenCertify(body BodyMap) (result []byte, err error) {
	var (
		bodyStr, sign, urlParam string
		bodyBs                  []byte
	)

	pubBody := make(BodyMap)
	pubBody.Set("app_id", a.AppId)
	pubBody.Set("method", "alipay.user.certify.open.certify")
	pubBody.Set("format", "JSON")
	if body != nil {
		if bodyBs, err = json.Marshal(body); err != nil {
			return nil, fmt.Errorf("json.Marshal：%v", err.Error())
		}
		bodyStr = string(bodyBs)
	}
	if a.AppCertSN != null {
		pubBody.Set("app_cert_sn", a.AppCertSN)
	}
	if a.AlipayRootCertSN != null {
		pubBody.Set("alipay_root_cert_sn", a.AlipayRootCertSN)
	}
	if a.ReturnUrl != null {
		pubBody.Set("return_url", a.ReturnUrl)
	}
	if a.Charset == null {
		pubBody.Set("charset", "utf-8")
	} else {
		pubBody.Set("charset", a.Charset)
	}
	if a.SignType == null {
		pubBody.Set("sign_type", "RSA2")
	} else {
		pubBody.Set("sign_type", a.SignType)
	}
	pubBody.Set("timestamp", time.Now().Format(TimeLayout))
	pubBody.Set("version", "1.0")
	if a.AppAuthToken != null {
		pubBody.Set("app_auth_token", a.AppAuthToken)
	}
	if a.AuthToken != null {
		pubBody.Set("auth_token", a.AuthToken)
	}
	if bodyStr != null {
		pubBody.Set("biz_content", bodyStr)
	}
	if sign, err = getRsaSign(pubBody, pubBody.Get("sign_type"), FormatPrivateKey(a.PrivateKey)); err != nil {
		return
	}
	pubBody.Set("sign", sign)
	urlParam = FormatAliPayURLParam(pubBody)
	if !a.IsProd {
		result,err = []byte(zfbSandboxBaseUrl + "?" + urlParam), nil
	} else {
		result,err = []byte(zfbBaseUrl + "?" + urlParam), nil
	}
	return
}

// UserCertifyOpenQuery 身份认证记录查询 https://docs.open.alipay.com/api_2/alipay.user.certify.open.query/
func (a *AliPayClient) UserCertifyOpenQuery(body BodyMap) (resp UserCertifyOpenQueryRsp, err error) {
	var bs []byte
	if bs, err = a.doAliPay(body, "alipay.user.certify.open.query"); err != nil {
		return
	}
	err = json.Unmarshal(bs, &resp)
	return
}
