// Package validate 用反射讀取結構的 validate 標籤，檢查欄位的值。
//
// 支援的規則（用逗號分隔）：
//
//	required  不可以是零值
//	min=N     數字 ≥ N；字串的字數、切片的長度 ≥ N
//	max=N     數字 ≤ N；字串的字數、切片的長度 ≤ N
//
// 切片中的結構會逐一檢查。
package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

// FieldError 描述一個欄位沒有通過哪條規則。
type FieldError struct {
	Field string // 欄位名稱，優先使用 json 標籤的名稱
	Rule  string // 例如 required、min=1
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("欄位 %s 不符合規則 %s", e.Field, e.Rule)
}

// Struct 檢查 v（結構或結構的指標），把所有問題用 errors.Join 一次回報。
func Struct(v any) error {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("validate: 需要結構，收到 %T", v)
	}
	return errors.Join(check(rv, "")...)
}

func check(rv reflect.Value, prefix string) []error {
	var errs []error
	for f, fv := range rv.Fields() { // Go 1.26：用迭代器走訪欄位
		if !f.IsExported() {
			continue
		}
		name := prefix + fieldName(f)
		for rule := range strings.SplitSeq(f.Tag.Get("validate"), ",") {
			if rule != "" && !ok(fv, rule) {
				errs = append(errs, &FieldError{Field: name, Rule: rule})
			}
		}
		if fv.Kind() == reflect.Slice { // 切片裡的結構也要檢查
			for i := range fv.Len() {
				if elem := reflect.Indirect(fv.Index(i)); elem.Kind() == reflect.Struct {
					sub := fmt.Sprintf("%s[%d].", name, i)
					errs = append(errs, check(elem, sub)...)
				}
			}
		}
	}
	return errs
}

// fieldName 取 json 標籤的名稱，例如 `json:"sku"` → sku。
func fieldName(f reflect.StructField) string {
	if name, _, _ := strings.Cut(f.Tag.Get("json"), ","); name != "" && name != "-" {
		return name
	}
	return f.Name
}

// ok 判斷欄位值 v 是否符合一條規則。
func ok(v reflect.Value, rule string) bool {
	if rule == "required" {
		return !v.IsZero()
	}
	key, arg, _ := strings.Cut(rule, "=")
	n, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		panic("validate: 規則寫錯了：" + rule) // 這是程式設計師的錯誤
	}
	var size int64
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		size = v.Int()
	case reflect.String:
		size = int64(utf8.RuneCountInString(v.String()))
	case reflect.Slice, reflect.Map:
		size = int64(v.Len())
	default:
		panic("validate: 不支援的型別 " + v.Type().String())
	}
	switch key {
	case "min":
		return size >= n
	case "max":
		return size <= n
	}
	panic("validate: 不認識的規則 " + rule)
}
