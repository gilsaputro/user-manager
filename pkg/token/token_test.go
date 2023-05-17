package token

import "testing"

func TestTokenConfig_GenerateToken(t *testing.T) {
	type args struct {
		bodyGenerate TokenBody
		bodyValidate TokenBody
	}
	tests := []struct {
		name            string
		tr              TokenConfig
		args            args
		wantErrValidate bool
		wantErr         bool
	}{
		{
			name: "success flow",
			tr: TokenConfig{
				Secret:        "my_secret_key",
				ExpTimeInHour: 1,
			},
			args: args{
				bodyGenerate: TokenBody{
					UserID:   1,
					Username: "user1",
				},
				bodyValidate: TokenBody{
					UserID:   1,
					Username: "user1",
				},
			},
			wantErrValidate: false,
			wantErr:         false,
		},
		{
			name: "error validate flow",
			tr: TokenConfig{
				Secret:        "my_secret_key",
				ExpTimeInHour: 1,
			},
			args: args{
				bodyGenerate: TokenBody{
					UserID:   1,
					Username: "user1",
				},
				bodyValidate: TokenBody{
					UserID:   2,
					Username: "user1",
				},
			},
			wantErrValidate: true,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.GenerateToken(tt.args.bodyGenerate)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenConfig.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				err = tt.tr.ValidateToken(got, tt.args.bodyValidate)
				if (err != nil) != tt.wantErrValidate {
					t.Errorf("TokenConfig.ValidateToken() error = %v, wantErr %v", err, tt.wantErrValidate)
					return
				}
			}
		})
	}
}
