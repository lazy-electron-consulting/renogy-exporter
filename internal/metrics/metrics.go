package metrics

import (
	_ "embed"
	"errors"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/yaml.v2"
)

const subsystem = "charge_controller"

var logger = log.New(log.Writer(), "[metrics] ", log.Lmsgprefix|log.Flags())

//go:embed renogy.yaml
var rawConfig []byte

type register struct {
	Address    uint16  `yaml:"address"`
	Multiplier float64 `yaml:"multiplier"`
	Byte       *uint   `yaml:"byte"`
	Signed     bool    `yaml:"signed"`
}

func (r register) Read(cc ChargeController) (float64, error) {
	if r.Byte != nil && r.Signed {
		raw, err := cc.Read(r.Address)
		if err != nil {
			return 0, err
		}
		return float64(int8(raw[*r.Byte])), nil
	} else if r.Byte == nil && !r.Signed {
		val, err := cc.ReadUint16(r.Address)
		if err != nil {
			return 0, err
		}
		return float64(val), nil
	}
	return 0, errors.New("unsupported configuration")
}

type metric struct {
	Name string `yaml:"name"`
	Help string `yaml:"help"`
}

type gauge struct {
	metric   `yaml:",inline"`
	register `yaml:",inline"`
}

type config struct {
	Gauges []gauge `yaml:"gauges,flow"`
}

type ChargeController interface {
	ReadUint16(address uint16) (uint16, error)
	Read(address uint16) ([]byte, error)
}

func Register(cc ChargeController) error {
	var config config
	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		logger.Printf("failed parse config %v\n", err)
		return err
	}

	for _, g := range config.Gauges {
		promauto.NewGaugeFunc(prometheus.GaugeOpts{
			Subsystem: subsystem,
			Name:      g.Name,
			Help:      g.Help,
		}, makeReader(g.register, cc))
	}

	return nil
}

func makeReader(r register, cc ChargeController) func() float64 {
	return func() float64 {
		val, err := r.Read(cc)
		if err != nil {
			logger.Printf("error reading %+v: %v", r, err)
			return -1
		}
		if r.Multiplier != 0 {
			return val * r.Multiplier
		}
		return val
	}
}
