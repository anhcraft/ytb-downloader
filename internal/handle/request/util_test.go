package request

import "testing"

func TestExtractPercentage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     Progress
		wantBool bool
	}{
		{
			name:  "valid input with progress prefix",
			input: "[[PROGRESS]]   3.1%,   5.65MiB, 179.84MiB,   7.94MiB/s,00:21",
			want: Progress{
				DownloadProgress:  "3.1%",
				DownloadedSize:    "5.65MiB",
				DownloadTotalSize: "179.84MiB",
				DownloadSpeed:     "7.94MiB/s",
				DownloadEta:       "00:21",
			},
			wantBool: true,
		},
		{
			name:     "missing progress prefix",
			input:    "0.2%, 15.00KiB, 179.84MiB, 791.59KiB/s,03:52",
			want:     Progress{},
			wantBool: false,
		},
		{
			name:     "empty input",
			input:    "",
			want:     Progress{},
			wantBool: false,
		},
		{
			name:     "only progress prefix",
			input:    "[[PROGRESS]]",
			want:     Progress{},
			wantBool: false,
		},
		{
			name:     "wrong prefix",
			input:    "[[DOWNLOAD]] 3.1%,   5.65MiB, 179.84MiB,   7.94MiB/s,00:21",
			want:     Progress{},
			wantBool: false,
		},
		{
			name:     "invalid input - wrong number of parts",
			input:    "[[PROGRESS]] 0.2%, 15.00KiB, 179.84MiB",
			want:     Progress{},
			wantBool: false,
		},
		{
			name:  "input with extra spaces after prefix",
			input: "[[PROGRESS]]    3.1%,   5.65MiB, 179.84MiB,   7.94MiB/s,00:21",
			want: Progress{
				DownloadProgress:  "3.1%",
				DownloadedSize:    "5.65MiB",
				DownloadTotalSize: "179.84MiB",
				DownloadSpeed:     "7.94MiB/s",
				DownloadEta:       "00:21",
			},
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotBool := ExtractProgress(tt.input)
			if gotBool != tt.wantBool {
				t.Errorf("ExtractProgress() bool = %v, want %v", gotBool, tt.wantBool)
			}
			if got != tt.want {
				t.Errorf("ExtractProgress() = %v, want %v", got, tt.want)
			}
		})
	}
}
