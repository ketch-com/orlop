package env

import (
	"go.ketch.com/lib/orlop/v2/service"
	"testing"
)

func Test_environImpl_GetPrefix(t *testing.T) {
	prefix := service.Name("test")
	type fields struct {
		prefix service.Name
	}
	tests := []struct {
		name   string
		fields fields
		want   service.Name
	}{
		{
			name: "",
			fields: fields{
				prefix: prefix,
			},
			want: prefix,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := environImpl{
				prefix: tt.fields.prefix,
			}
			if got := e.GetPrefix(); got != tt.want {
				t.Errorf("GetPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
