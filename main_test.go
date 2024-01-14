package main

import "testing"

func Test_eval(t *testing.T) {
	type args struct {
		jsonIn     string
		expression string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
		wantErr    bool
	}{
		{
			name: "string equals",
			args: args{
				jsonIn:     `{"name": "abc"}`,
				expression: "i.name == 'abc'",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "startsWith",
			args: args{
				jsonIn:     `{"name": "abc"}`,
				expression: "i.name.startsWith('x')",
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name: "logical and",
			args: args{
				jsonIn:     `{"name": "abc", "val": 9001}`,
				expression: "i.name == 'abc' && i.val > 9000",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "logical or",
			args: args{
				jsonIn:     `{"name": "abc", "val": 9000}`,
				expression: "i.name == 'abc' || i.val < 10",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "macro has true",
			args: args{
				jsonIn:     `{"name": "abc"}`,
				expression: "has(i.name)",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "macro has false",
			args: args{
				jsonIn:     `{"name": "abc"}`,
				expression: "has(i.val)",
			},
			wantResult: false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := eval(tt.args.jsonIn, tt.args.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("eval() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
