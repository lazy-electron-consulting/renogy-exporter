// config
package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lazy-electron-consulting/renogy-exporter/internal/util"
	"gopkg.in/yaml.v2"
)

const (
	DefaultAddress  = ":8000"
	DefaultBaudRate = 9600
	DefaultDataBits = 8
	DefaultStopBits = 1
	DefaultParity   = "N"
	DefaultTimeout  = 5 * time.Second
	DefaultUnitID   = 1
)

type Modbus struct {
	Path     string        `json:"path,omitempty" yaml:"path,omitempty"`
	UnitID   byte          `json:"unitId,omitempty" yaml:"unitId,omitempty"`
	BaudRate int           `json:"baudRate,omitempty" yaml:"baudRate,omitempty"`
	DataBits int           `json:"dataBits,omitempty" yaml:"dataBits,omitempty"`
	StopBits int           `json:"stopBits,omitempty" yaml:"stopBits,omitempty"`
	Parity   string        `json:"parity,omitempty" yaml:"parity,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (m *Modbus) Defaults() {
	m.BaudRate = util.Default(m.BaudRate, DefaultBaudRate)
	m.DataBits = util.Default(m.DataBits, DefaultDataBits)
	m.StopBits = util.Default(m.StopBits, DefaultStopBits)
	m.Parity = util.Default(m.Parity, DefaultParity)
	m.Timeout = util.Default(m.Timeout, DefaultTimeout)
	m.UnitID = util.Default(m.UnitID, DefaultUnitID)
}

type State struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value uint16 `json:"value,omitempty" yaml:"value,omitempty"`
}

type Gauge struct {
	Name       string  `json:"name,omitempty" yaml:"name,omitempty"`
	Help       string  `json:"help,omitempty" yaml:"help,omitempty"`
	Address    uint16  `json:"address,omitempty" yaml:"address,omitempty"`
	Byte       *uint8  `json:"byte,omitempty" yaml:"byte,omitempty"`
	Signed     bool    `json:"signed,omitempty" yaml:"signed,omitempty"`
	Multiplier float32 `json:"multiplier,omitempty" yaml:"multiplier,omitempty"`
	States     []State `json:"states,omitempty" yaml:"states,omitempty"`
}

type Config struct {
	Address string  `json:"address,omitempty" yaml:"address,omitempty"`
	Modbus  *Modbus `json:"modbus,omitempty" yaml:"modbus,omitempty"`
	Gauges  []Gauge `json:"gauges,omitempty" yaml:"gauges,omitempty"`
}

func (c *Config) defaults() {
	c.Address = util.Default(c.Address, DefaultAddress)
	c.Modbus = util.Default(c.Modbus, &Modbus{})
	c.Modbus.Defaults()
}

// ParseYaml reads yaml-formatted config in strict mode, filling in any default
// values.
func ParseYaml(r io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(r)
	decoder.SetStrict(true)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}
	config.defaults()
	return &config, nil
}

func ReadYaml(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ParseYaml(f)
}
