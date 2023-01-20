package token

import "testing"

func TestToken_String(t1 *testing.T) {
	type fields struct {
		TT TokenType
		S  string
		E  *any
	}
	var data any = "extra-data"
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"with S no E", fields{TTnl, "\n", nil}, `TT.NL("\n")`},
		{"no S no E", fields{TTcustom, "", nil}, `TT.Custom`},
		{"no S with E", fields{TTcustom, "", &data}, `TT.Custom["extra-data"]`},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Token{
				TT: tt.fields.TT,
				S:  tt.fields.S,
				E:  tt.fields.E,
			}
			if got := t.String(); got != tt.want {
				t1.Errorf("\n got: `%s`\nwant: `%s`", got, tt.want)
			}
		})
	}
}
