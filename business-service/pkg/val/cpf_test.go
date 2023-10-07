package val_test

import (
	"testing"

	"github.com/lucasd-coder/business-service/pkg/val"
)

func TestIsCPF(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "InvalidData_ShouldReturnFalse",
			arg:  "3467875434578764345789654",
			want: false,
		},
		{"InvalidData_ShouldReturnFalse", "", false},
		{"InvalidData_ShouldReturnFalse", "AAAAAAAAAAA", false},
		{"InvalidPattern_ShouldReturnFalse", "000.000.000-00", false},
		{"InvalidPattern_ShouldReturnFalse", "222.222.222-22", false},
		{"InvalidPattern_ShouldReturnFalse", "333.333.333-33", false},
		{"InvalidPattern_ShouldReturnFalse", "444.444.444-44", false},
		{"InvalidPattern_ShouldReturnFalse", "555.555.555-55", false},
		{"InvalidPattern_ShouldReturnFalse", "666.666.666-66", false},
		{"InvalidPattern_ShouldReturnFalse", "777.777.777-77", false},
		{"InvalidPattern_ShouldReturnFalse", "888.888.888-88", false},
		{"InvalidPattern_ShouldReturnFalse", "999.999.999-99", false},
		{"InvalidPattern_ShouldReturnFalse", "248.438.034-08", false},
		{"InvalidPattern_ShouldReturnFalse", "248 438 034 80", false},
		{"InvalidPattern_ShouldReturnFalse", "099-075-865.60", false},
		{"Valid_ShouldReturnTrue", "042.618.600-15", true},
		{"Valid_ShouldReturnTrue", "04261860015", true},
		{"Valid_ShouldReturnTrue", "099.075.865-60", true},
		{"Valid_ShouldReturnTrue", "09907586560", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := val.IsCPF(tt.arg); got != tt.want {
				t.Errorf("IsCPF() = %v, want %v", got, tt.want)
			}
		})
	}
}
