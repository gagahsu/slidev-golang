package money

import "testing"

func TestString(t *testing.T) {
	tests := []struct {
		name string
		in   Money
		want string
	}{
		{"零元", 0, "NT$0"},
		{"三位數", 999, "NT$999"},
		{"四位數", 1000, "NT$1,000"},
		{"七位數", 1234567, "NT$1,234,567"},
		{"負數三位", -300, "-NT$300"},
		{"負數四位", -1500, "-NT$1,500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.String(); got != tt.want {
				t.Errorf("Money(%d).String() = %q，want %q",
					int(tt.in), got, tt.want)
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	for b.Loop() {
		_ = Money(1234567).String()
	}
}
