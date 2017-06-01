package httptool

import (
	"net/url"
	"net/http"
	"io/ioutil"
"magpie/internal/com/logger"
)

// 普通get/post参数类型
type Params map[string]string

/**
 * 发起get请求
 */
func Get(urlStr string, params Params) ([]byte, error) {
	// 解析url并拼接参数
	u, err := url.Parse(urlStr)
	q := u.Query()
	var res *http.Response
	var resBytes []byte
	if err != nil {
		goto deal_err
	}

	for k, p := range params {
		q.Set(k, p)
	}
	u.RawQuery = q.Encode()
	logger.Debug("http.tool", u.String())
	res, err = http.Get(u.String())
	if err != nil {
		goto deal_err
	}

	// 读取返回结果
	resBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		goto deal_err
	}

	//  处理错误
	deal_err:
	if err != nil {
		logger.Error("http.tool", logger.Format(
			"err", err.Error(),
			"url", urlStr,
			"params", params,
		))
		return nil, err
	}

	return resBytes, nil
}