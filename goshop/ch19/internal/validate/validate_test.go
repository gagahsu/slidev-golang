package validate

import (
	"errors"
	"strings"
	"testing"
)

type item struct {
	SKU string `json:"sku" validate:"required"`
	Qty int    `json:"qty" validate:"min=1,max=99"`
}

type order struct {
	Items []item `json:"items" validate:"min=1"`
	Note  string `validate:"max=5"`
	skip  int    // 沒有匯出的欄位會被略過
}

func TestStruct(t *testing.T) {
	good := order{Items: []item{{"A", 1}, {"B", 99}}, Note: "五個中文字"}
	if err := Struct(&good); err != nil {
		t.Errorf("Struct(good) = %v", err)
	}

	bad := order{Items: []item{{"", 0}, {"B", 100}}, Note: "六個中文字喔"}
	err := Struct(bad)
	for _, want := range []string{
		"items[0].sku 不符合規則 required",
		"items[0].qty 不符合規則 min=1",
		"items[1].qty 不符合規則 max=99",
		"Note 不符合規則 max=5",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("錯誤訊息缺少 %q：\n%v", want, err)
		}
	}
	if fe, ok := errors.AsType[*FieldError](err); !ok || fe.Field != "items[0].sku" {
		t.Errorf("errors.AsType 取出 %+v", fe)
	}
	if err := Struct(order{}); err == nil {
		t.Error("空的 items 應該不通過 min=1")
	}
	if err := Struct(42); err == nil {
		t.Error("非結構應該回傳錯誤")
	}
}
