package discordgo

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNullString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantBeNull bool
	}{
		{"empty string", args{""}, false},
		{"valid string", args{"test"}, false},
		{"valid nullable string", args{nullStringValue}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name+".Elem()", func(t *testing.T) {
			if got := NullString(tt.args.s); !reflect.DeepEqual(got.Elem(), tt.args.s) {
				t.Errorf("NullString().Elem() = %v, want %v", got.Elem(), tt.args.s)
			}
		})

		t.Run(tt.name+".Ptr()", func(t *testing.T) {
			if got := NullString(tt.args.s); !reflect.DeepEqual(got.Ptr(), &tt.args.s) {
				t.Errorf("NullString().Ptr() = %v, want %v", got.Ptr(), &tt.args.s)
			}
		})

		t.Run(tt.name+".IsNull()", func(t *testing.T) {
			if got := NullString(tt.args.s); !reflect.DeepEqual(got.IsNull(), tt.wantBeNull) {
				t.Errorf("NullString().IsNull() = %v, want %v", got.IsNull(), tt.wantBeNull)
			}
		})
	}
}

func TestNullableString_Null(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty string", args{""}},
		{"valid string", args{"test"}},
		{"valid nullable string", args{nullStringValue}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullString(tt.args.s).Null().Elem()
			if got != nullStringValue {
				t.Errorf("NullableString().Null().Elem() = %v, want %v", got, nullStringValue)
			}
		})
	}
}

func TestNullableString_MarshalJSON(t *testing.T) {
	type fields struct {
		str *string
	}
	tests := []struct {
		name    string
		input   NullableString
		want    []byte
		wantErr bool
	}{
		{"empty string", NullString(""), []byte(`""`), false},
		{"classic string value", NullString("test"), []byte(`"test"`), false},
		{"nullable string value", NullString("").Null(), []byte(`null`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullableString.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullableString.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableString_UnmarshalJSON(t *testing.T) {
	type fields struct {
		str *string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantBeNull bool
		wantErr    bool
	}{
		{"empty string", args{[]byte(`""`)}, "", false, false},
		{"classic string value", args{[]byte(`"test"`)}, "test", false, false},
		{"nullable string value", args{[]byte(`null`)}, "null", true, false},
		{"crash nil", args{nil}, "", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := NullString("")
			if err := ns.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("NullableString.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(ns.Elem(), tt.want) {
				t.Errorf("NullableString.UnmarshalJSON() = %v, want %v", ns.Elem(), tt.want)
			}

			if !reflect.DeepEqual(ns.IsNull(), tt.wantBeNull) {
				fmt.Println("value: ", ns)
				t.Errorf("NullableString.UnmarshalJSON() = %v, want %v", ns.IsNull(), tt.wantBeNull)
			}
		})
	}
}
