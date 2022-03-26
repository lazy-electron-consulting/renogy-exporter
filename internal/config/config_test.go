package config_test

import (
	"embed"
	"testing"
	"time"

	"github.com/lazy-electron-consulting/renogy-exporter/internal/config"
	"github.com/lazy-electron-consulting/renogy-exporter/internal/util"
	"github.com/stretchr/testify/require"
)

//go:embed testdata
var testdata embed.FS

func TestParseYaml(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		path    string
		want    *config.Config
		wantErr bool
	}{
		{
			name: "full",
			path: "testdata/full.yaml",
			want: &config.Config{
				Address: ":8080",
				Modbus: &config.Modbus{
					Path:     "/dev/ttyS0",
					BaudRate: 19200,
					DataBits: 8,
					StopBits: 1,
					Parity:   "N",
					Timeout:  5 * time.Minute,
					UnitID:   2,
				},
				Gauges: []config.Gauge{
					{Name: "16bit", Help: "16bit int", Address: 0x100},
					{Name: "8bit", Help: "8bit int", Address: 0x101, Byte: util.Ptr(uint8(1)), Signed: true},
					{Name: "16bit with multiplier", Help: "16bit int", Address: 0x102, Multiplier: 0.1},
				},
			},
		},
		{
			name:    "invalidYaml",
			path:    "testdata/test.json",
			wantErr: true,
		},
		{
			name: "emptyUsesDefaults",
			path: "testdata/empty.yaml",
			want: &config.Config{
				Address: config.DefaultAddress,
				Modbus: &config.Modbus{
					BaudRate: config.DefaultBaudRate,
					DataBits: config.DefaultDataBits,
					StopBits: config.DefaultStopBits,
					Parity:   config.DefaultParity,
					Timeout:  config.DefaultTimeout,
					UnitID:   config.DefaultUnitID,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			file, err := testdata.Open(tt.path)
			require.NoError(t, err)

			got, err := config.ParseYaml(file)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, tt.want, got)
		})
	}
}
