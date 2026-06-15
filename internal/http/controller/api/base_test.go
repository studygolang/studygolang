package api

import "testing"

func TestIsValidObjType(t *testing.T) {
	tests := []struct {
		name     string
		objtype  int
		expected bool
	}{
		// 合法枚举值
		{"topic=0 is valid (regression for C3)", 0, true},
		{"article", 1, true},
		{"resource", 2, true},
		{"wiki", 3, true},
		{"project", 4, true},
		{"book", 5, true},
		{"interview", 6, true},
		// 非法值
		{"-1 rejected", -1, false},
		{"7 rejected (out of enum)", 7, false},
		{"100 rejected (out of enum)", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidObjType(tt.objtype); got != tt.expected {
				t.Errorf("isValidObjType(%d) = %v, want %v", tt.objtype, got, tt.expected)
			}
		})
	}
}
