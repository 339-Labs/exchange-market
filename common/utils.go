package common

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

func TimesStamp() string {
	timesStamp := time.Now().Unix() * 1000
	return strconv.FormatInt(timesStamp, 10)
}

func TimesStampSec() string {
	timesStamp := time.Now().Unix()
	return strconv.FormatInt(timesStamp, 10)
}

func BuildJsonParams(params map[string]string) (string, error) {
	if params == nil {
		return "", errors.New("illegal parameter")
	}
	data, err := json.Marshal(params)
	if err != nil {
		return "", errors.New("json convert string error")
	}
	jsonBody := string(data)
	return jsonBody, nil
}

func BuildGetParams(params map[string]string) string {
	//urlParams := url.Values{}
	//if params != nil && len(params) > 0 {
	//	for k := range params {
	//		urlParams.Add(k, params[k])
	//	}
	//}
	//return "?" + urlParams.Encode()
	if len(params) == 0 {
		return ""
	}
	return "?" + SortParams(params)
}

func SortParams(params map[string]string) string {
	keys := make([]string, len(params))
	i := 0
	for k, _ := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	sorted := make([]string, len(params))
	i = 0
	for _, k := range keys {
		//sorted[i] = k + "=" + url.QueryEscape(params[k])
		sorted[i] = k + "=" + params[k]
		i++
	}
	return strings.Join(sorted, "&")
}

func JSONToMap(str string) map[string]interface{} {

	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}

func BytesToMap(bytes []byte) map[string]interface{} {

	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(bytes), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}

func JSONToArrMap(str string) []map[string]interface{} {
	var tempMap []map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}

func BytesToArrMap(bytes []byte) []map[string]interface{} {

	var tempMap []map[string]interface{}
	err := json.Unmarshal([]byte(bytes), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}

func NewParams() map[string]string {
	return make(map[string]string)
}

func ToJson(v interface{}) (string, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
