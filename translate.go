package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var TRANSLATE_URL = "http://api.fanyi.baidu.com/api/trans/vip/translate"

type TranslateResultElement struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type TranslateResult struct {
	From    string                   `json:"from"`
	To      string                   `json:"to"`
	Results []TranslateResultElement `json:"trans_result"`
}

func calculateSign(appid string, key string, query string) (string, string) {
	salt := fmt.Sprintf("%d", time.Now().Unix())
	str1 := appid + query + salt + key
	hash := md5.Sum([]byte(str1))
	return hex.EncodeToString(hash[:]), salt
}

func translate(content string, appid string, key string, to string) (*TranslateResult, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, TRANSLATE_URL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	sign, salt := calculateSign(appid, key, content)
	q.Add("q", content)
	q.Add("appid", appid)
	q.Add("salt", salt)
	q.Add("sign", sign)
	q.Add("from", "auto")
	q.Add("to", to)
	req.URL.RawQuery = q.Encode()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result TranslateResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
