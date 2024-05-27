package tool

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/pkg/errors"

	"golang.org/x/exp/constraints"
)

// reflectStruct 获取结构体反射
func reflectStruct(obj interface{}) *reflect.Value {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		return &v
	}
	return nil
}

// GetColumnFromStruct 从结构体切片中所有结构体获取属性值
// array 为传入结构体切片
// field 为结构体对应的字段
// 如果字段名不存在、或者泛型类型和字段不一致则不获取
// 另外，字段需要导出才有效
func GetColumnFromStruct[S interface{}, T comparable](array []S, field string) []T {
	if len(array) == 0 {
		return []T{}
	}

	fields := make([]T, 0, len(array))
	for _, s := range array {
		v := reflectStruct(s)
		if v == nil {
			continue
		}

		f := v.FieldByName(field)
		if !f.CanInterface() {
			continue
		}

		e, ok := f.Interface().(T)
		if ok {
			fields = append(fields, e)
		}
	}

	return fields
}

// GetColumnFromMap 从map切片中所有map获取key对应值
// array 为传入map切片
// key 为map对应的key
// 如果字段名不存在、或者泛型类型和字段不一致则不获取
func GetColumnFromMap[S any, T comparable](array []map[string]S, key string) []T {
	if len(array) == 0 {
		return []T{}
	}

	fields := make([]T, 0, len(array))
	for _, m := range array {
		if sv, ok := m[key]; ok {
			if v, ok := reflect.ValueOf(sv).Interface().(T); ok {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

// GenMapFromMap 从map切片中选取两个值生成map
// array 为传入map切片
// keyCol 为map对应的key的列名
// valCol 为map对应的val的列名
// 如果字段名不存在、或者泛型类型和字段不一致则不获取
func GenMapFromMap[B any, S comparable, T any](array []map[string]B, keyCol string, valCol string) map[S]T {
	if len(array) == 0 {
		return map[S]T{}
	}

	cols := make(map[S]T, len(array))
	for _, m := range array {
		if k, ok := reflect.ValueOf(m[keyCol]).Interface().(S); ok {
			if v, ok := reflect.ValueOf(m[valCol]).Interface().(T); ok {
				cols[k] = v
			}
		}
	}

	return cols
}

// GenMapFromStruct 从结构体切片中选取两个值生成map
// array 为传入map切片
// keyCol 为结构体对应的key的列名
// valCol 为结构体对应的val的列名
// 如果字段名不存在、或者泛型类型和字段不一致则不获取
// 另外，字段需要导出才有效
func GenMapFromStruct[B any, S comparable, T any](array []B, keyCol string, valCol string) map[S]T {
	if len(array) == 0 {
		return map[S]T{}
	}

	cols := make(map[S]T, len(array))
	for _, obj := range array {
		rObj := reflectStruct(obj)
		if nil == rObj {
			continue
		}

		// 检查是否能转为 Interface
		f := rObj.FieldByName(keyCol)
		if !f.CanInterface() {
			continue
		}

		key, ok := f.Interface().(S)
		if !ok {
			continue
		}

		f = rObj.FieldByName(valCol)
		if !f.CanInterface() {
			continue
		}

		val, ok := f.Interface().(T)
		if !ok {
			continue
		}

		cols[key] = val
	}

	return cols
}

// Filter 过滤器
type Filter interface {
	Where(string, string, any) Filter
	In(string, []any) Filter
}

type Number interface {
	constraints.Float | constraints.Integer
}

func errorWrap(err *error, msg string) {
	if *err == nil {
		*err = fmt.Errorf("%s\n%s", msg, debug.Stack())
	} else {
		*err = fmt.Errorf("%s\n%s; %w", msg, debug.Stack(), *err)
	}
}

// StructsFilter 切片过滤器
type StructsFilter[T any] struct {
	data      []T
	errorWrap error
	mtx       sync.RWMutex
}

// SourceStruct 数据来源是 Struct
func SourceStruct[T any](data []T) *StructsFilter[T] {
	return &StructsFilter[T]{
		data: data,
	}
}

// Data 返回数据
func (f *StructsFilter[T]) Data() []T {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	return f.data
}

// Error 返回过程错误
func (f *StructsFilter[T]) Error() error {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	return f.errorWrap
}

// Filter 过滤器
func (f *StructsFilter[T]) Filter(cond func(t *T) bool) *StructsFilter[T] {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	// 复制过滤
	fData := make([]T, 0, len(f.data))
	for _, v := range f.data {
		if cond(&v) {
			fData = append(fData, v)
		}
	}
	f.data = fData
	return f
}

// compareNumber 数字对比函数
func compareNumber[T constraints.Float | constraints.Integer](a T, op string, b T) bool {
	switch op {
	case "=":
		fallthrough
	case "==":
		return a-b == 0
	case "!=":
		return a-b != 0
	case "<":
		return a-b < 0
	case "<=":
		return a-b <= 0
	case ">":
		return a-b > 0
	case ">=":
		return a-b >= 0
	}
	return false
}

// Where 通过比较过滤数字
// val 为比较值，支持整型、浮点型和字符串
// 使用常量时，当 val 为整型时不可带小数点，而当 val 为浮点型时必须带小数点
// val 为字符串时仅支持等于、不等于操作
func (f *StructsFilter[T]) Where(field string, op string, val any) *StructsFilter[T] {

	var errWhere error
	defer func() {
		if errWhere != nil {
			f.mtx.Lock()
			errorWrap(&f.errorWrap, errWhere.Error())
			f.mtx.Unlock()
		}
	}()

	switch val.(type) {
	case float32, float64: // 浮点型
		target := reflect.ValueOf(val).Float()

		f.Filter(func(t *T) bool {
			tVal := reflectStruct(*t)
			if tVal == nil {
				errWhere = errors.New("无效的结构体")
				return false
			}

			fVal := tVal.FieldByName(field)
			if !fVal.CanFloat() {
				errWhere = fmt.Errorf("结构体属性%s比较值非浮点型", field)
				return false
			}

			origin := fVal.Float()
			return compareNumber(origin, op, target)
		})

	case int, int8, int16, int32, int64: // 整型
		target := reflect.ValueOf(val).Int()

		f.Filter(func(t *T) bool {
			tVal := reflectStruct(*t)
			if tVal == nil {
				errWhere = errors.New("无效的结构体")
				return false
			}

			fVal := tVal.FieldByName(field)
			if !fVal.CanInt() {
				errWhere = fmt.Errorf("结构体属性%s比较值非整型", field)
				return false
			}

			origin := fVal.Int()
			return compareNumber(origin, op, target)
		})

	case string: // 字符串
		target := reflect.ValueOf(val).String()

		f.Filter(func(t *T) bool {
			tVal := reflectStruct(*t)
			if tVal == nil {
				errWhere = errors.New("无效的结构体")
				return false
			}

			fVal := tVal.FieldByName(field)
			if !fVal.CanInterface() {
				errWhere = fmt.Errorf("结构体属性%s比较值无法导出", field)
				return false
			}

			// 反射结构体成员
			origin, ok := fVal.Interface().(string)
			if !ok {
				errWhere = fmt.Errorf("结构体属性%s比较值非字符串", field)
				return false
			}

			switch op {
			case "=":
				fallthrough
			case "==":
				return origin == target
			case "!=":
				return origin != target
			default:
				errWhere = fmt.Errorf("不支持的操作:%s", op)
			}

			return false
		})

	default:
		errWhere = fmt.Errorf("不支持的类型:%s", reflect.TypeOf(val).Kind().String())
	}

	return f
}

// In 过滤元素
// pattern 为对应过滤的切片类型
func (f *StructsFilter[T]) In(field string, patterns any) *StructsFilter[T] {
	patVal := reflect.ValueOf(patterns)
	if patVal.Kind() != reflect.Slice {
		f.mtx.Lock()
		errorWrap(&f.errorWrap, fmt.Sprintf("无效的 patterns 类型 %s，请传入切片类型", patVal.Kind().String()))
		f.mtx.Unlock()
		return f
	}

	var values []*reflect.Value
	for i := 0; i < patVal.Len(); i++ {
		v := patVal.Index(i)
		if v.Type().Comparable() && v.CanInterface() {
			values = append(values, &v)
		}
	}

	f.Filter(func(t *T) bool {
		st := reflectStruct(*t)
		if st == nil {
			return false
		}

		fVal := st.FieldByName(field)
		for _, val := range values {
			if fVal.Type().Comparable() && fVal.CanInterface() && reflect.DeepEqual((*val).Interface(), fVal.Interface()) {
				return true
			}
		}

		return false
	})
	return f
}

// MapFilter map 过滤器
// s := SourceMap(users).Where("weight", ">", 13.0).In("name", "xx")
// s.Data() 数据
// s.Error() 错误
type MapFilter[T comparable] struct {
	data      []map[T]any
	errorWrap error
	mtx       sync.RWMutex
}

// SourceMap 数据来源是 map
func SourceMap[T comparable](data []map[T]any) *MapFilter[T] {
	return &MapFilter[T]{
		data: data,
	}
}

func (f *MapFilter[T]) Data() []map[T]any {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	return f.data
}

func (f *MapFilter[T]) Error() error {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	return f.errorWrap
}

// Filter 过滤函数
func (f *MapFilter[T]) Filter(cond func(t map[T]any) bool) *MapFilter[T] {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	// 复制过滤
	fData := make([]map[T]any, 0, len(f.data))
	for _, v := range f.data {
		if cond(v) {
			fData = append(fData, v)
		}
	}
	f.data = fData
	return f
}

// Where 通过比较过滤数字
// val 为比较值，支持整型、浮点型和字符串
// 使用常量时，当 val 为整型时不可带小数点，而当 val 为浮点型时必须带小数点
// val 为字符串时仅支持等于、不等于操作
func (f *MapFilter[T]) Where(key T, op string, val any) *MapFilter[T] {

	var errWhere error
	defer func() {
		if errWhere != nil {
			f.mtx.Lock()
			errorWrap(&f.errorWrap, errWhere.Error())
			f.mtx.Unlock()
		}
	}()

	switch val.(type) {
	case float32, float64: // 浮点型
		target := reflect.ValueOf(val).Float()

		f.Filter(func(t map[T]any) bool {
			obj, ok := t[key]
			if !ok {
				return false
			}

			vObj := reflect.ValueOf(obj)
			if !vObj.CanFloat() {
				return false
			}

			origin := vObj.Float()

			return compareNumber(origin, op, target)
		})

	case int, int8, int16, int32, int64: // 整型
		target := reflect.ValueOf(val).Int()

		f.Filter(func(t map[T]any) bool {
			obj, ok := t[key]
			if !ok {
				return false
			}

			vObj := reflect.ValueOf(obj)
			if !vObj.CanInt() {
				return false
			}

			origin := vObj.Int()

			return compareNumber(origin, op, target)
		})

	case string: // 字符串
		target := reflect.ValueOf(val).String()

		f.Filter(func(t map[T]any) bool {
			obj, ok := t[key]
			if !ok {
				return false
			}

			vObj := reflect.ValueOf(obj)
			if !vObj.CanInt() {
				return false
			}
			// 反射结构体成员
			origin, ok := vObj.Interface().(string)
			if !ok {
				return false
			}

			switch op {
			case "=":
				fallthrough
			case "==":
				return origin == target
			case "!=":
				return origin != target
			default:
				errWhere = fmt.Errorf("不支持的操作:%s", op)
			}

			return false
		})

	default:
		errWhere = fmt.Errorf("不支持的类型:%s", reflect.TypeOf(val).Kind().String())
	}

	return f
}

// In 过滤元素
func (f *MapFilter[T]) In(key T, patterns any) *MapFilter[T] {
	patVal := reflect.ValueOf(patterns)
	if patVal.Kind() != reflect.Slice {
		f.mtx.Lock()
		errorWrap(&f.errorWrap, fmt.Sprintf("无效的 patterns 类型 %s，请传入切片类型", patVal.Kind().String()))
		f.mtx.Unlock()
		return f
	}

	var values []*reflect.Value
	for i := 0; i < patVal.Len(); i++ {
		v := patVal.Index(i)
		if v.Type().Comparable() && v.CanInterface() {
			values = append(values, &v)
		}
	}

	f.Filter(func(t map[T]any) bool {
		obj, ok := t[key]
		if !ok {
			return false
		}

		fVal := reflect.ValueOf(obj)
		for _, val := range values {
			if fVal.Type().Comparable() && fVal.CanInterface() && reflect.DeepEqual((*val).Interface(), fVal.Interface()) {
				return true
			}
		}

		return false
	})
	return f
}
