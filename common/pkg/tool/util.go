package tool

import (
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cast"
	"github.com/syyongx/php2go"
	"study/common/pkg/timeutil"
)

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

func ToUrlParams(data map[string]interface{}, keys ...[]string) string {
	var params []string
	if len(keys) > 0 {
		for _, k := range keys[0] {
			if "sign" == k {
				continue
			}
			v := data[k]
			switch v.(type) {
			case string, int, int64, int32, int16, int8, float64, float32:
				params = append(params, k+"="+fmt.Sprintf("%v", v))
			}
		}
	} else {
		for k, v := range data {
			if "sign" == k {
				continue
			}
			switch v.(type) {
			case string, int, int64, int16, int8, float64, float32:
				params = append(params, k+"="+fmt.Sprintf("%v", v))
			}
		}
	}
	return strings.Join(params, "&")
}

func TransMap(data map[string]interface{}) map[string]string {
	var params map[string]string
	params = make(map[string]string)
	for k, v := range data {
		params[k] = fmt.Sprintf("%v", v)
	}
	return params
}

// 判断元素是否存在数组中
func IsContainInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func IsContainString(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 三目运算的函数
func Ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// 字符串转int64切片
func StringToSliceInt64(data string) []int64 {
	var slice []int64

	if php2go.Empty(data) {
		return slice
	}

	stringSlice := php2go.Explode(",", data)
	for _, v := range stringSlice {
		slice = append(slice, cast.ToInt64(v))
	}

	return slice
}

// 字符串截取
func SubStrRange(s string, length int) string {
	var size, n int
	for i := 0; i < length && n < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[n:])
		n += size
	}

	return s[:n]
}

// 获取腾讯云图片（Picture）存储文件路径（工单）
func GetQcloudCosFilePathPictureGd(env, fileName string) string {
	return fmt.Sprintf("/gd/%s/%s/%s", env, timeutil.GetTodayStr(), fileName)
}

// 格式化
func GetTimeFormatText(createTime int64) string {
	if createTime == 0 {
		return ""
	}
	createTimeUnix := time.Unix(createTime, 0)
	return createTimeUnix.Format("2006-01-02 15:04:05")
}

// 获取指定时间单位的时间戳（单位：秒）
func GetAppointTimeUnitTimestamp(number, unit int64) int64 {
	timestamp := cast.ToInt64(0)
	switch unit {
	case 1: // 秒
		timestamp = number * 1
	case 2: // 分
		timestamp = number * 60
	case 3: // 时
		timestamp = number * 60 * 60
	case 4: // 天
		timestamp = number * 60 * 60 * 24
	}

	return timestamp
}

// 对字符串进行md5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 获取文件后缀，如："123.text" 返回 ".text"
func GetFileExt(filePath string) string {
	return path.Ext(filePath)
}

// ToString 强转interface类型到string类型
func ToString(i interface{}) string {
	v, _ := ToStringE(i)
	return v
}

// ToStringE 强转interface类型到string, 支持错误返回值
func ToStringE(i interface{}) (string, error) {
	i = indirectToStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatInt(int64(s), 10), nil
	case uint64:
		return strconv.FormatInt(int64(s), 10), nil
	case uint32:
		return strconv.FormatInt(int64(s), 10), nil
	case uint16:
		return strconv.FormatInt(int64(s), 10), nil
	case uint8:
		return strconv.FormatInt(int64(s), 10), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		v := reflect.ValueOf(i)
		if method, ok := reflect.TypeOf(i).MethodByName("String"); ok && method.Type.NumIn() == 0 &&
			method.Type.NumOut() == 1 && method.Type.Out(0).Kind() == reflect.String {
			return method.Func.Call([]reflect.Value{v})[0].String(), nil
		}
		switch v.Kind() {
		case reflect.Func:
			fullName := runtime.FuncForPC(v.Pointer()).Name()
			ss := strings.Split(fullName, ".")
			if len(ss) > 0 {
				return ss[len(ss)-1], nil
			} else {
				return fullName, nil
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return ToStringE(v.Uint())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return ToStringE(v.Int())
		case reflect.Float64, reflect.Float32:
			return ToStringE(v.Float())
		case reflect.Bool:
			return ToStringE(v.Bool())
		case reflect.String:
			return v.String(), nil
		}
		return "", fmt.Errorf("unable to cast %#v of types %T to string", i, i)
	}
}

func indirectToStringerOrError(a any) any {
	if a == nil {
		return nil
	}
	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func StructToMap(obj interface{}) map[string]interface{} {
	data, _ := json.Marshal(&obj)
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return m
}

// 复制list
func CopyList(originalList list.List) *list.List {
	copiedList := list.New()
	for e := originalList.Front(); e != nil; e = e.Next() {
		copiedList.PushBack(e.Value)
	}
	return copiedList
}
