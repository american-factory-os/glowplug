package service

import "testing"

func TestValidateMqttURL(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "valid mqtt url",
			args:    args{"tcp://localhost:1883"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid mqtt url",
			args:    args{"mqtt://localhost:1883"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid ws url",
			args:    args{"ws://localhost:1883"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid wss url",
			args:    args{"wss://localhost:1883"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			args:    args{"http://localhost:1883"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "port too short",
			args:    args{"mqtt://localhost:/"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "host too short",
			args:    args{"mqtt://s:1883"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "empty url",
			args:    args{""},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateBrokerURI(tt.args.rawURL)

			if err != nil && !tt.wantErr {
				t.Errorf("validateBrokerURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("validateBrokerURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
