package metrics

import (
	_ "embed"
	"fmt"
	"log"
	"math"

	"github.com/lazy-electron-consulting/renogy-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const subsystem = "charge_controller"

var logger = log.New(log.Writer(), "[metrics] ", log.Lmsgprefix|log.Flags())

type gaugeReader func(r Reader) (float64, error)

func int8GaugeReader(index uint8, address uint16) gaugeReader {
	return func(r Reader) (float64, error) {
		raw, err := r.Read(address)
		if err != nil {
			return 0, err
		}
		return float64(int8(raw[index])), nil
	}
}

func uint8GaugeReader(index uint8, address uint16) gaugeReader {
	return func(r Reader) (float64, error) {
		raw, err := r.Read(address)
		if err != nil {
			return 0, err
		}
		return float64(uint8(raw[index])), nil
	}
}

func uint16GaugeReader(address uint16) gaugeReader {
	return func(r Reader) (float64, error) {
		val, err := r.ReadUint16(address)
		if err != nil {
			return 0, err
		}
		return float64(val), nil
	}
}

func multipliedGaugeReader(inner gaugeReader, multiplier float64) gaugeReader {
	return func(r Reader) (float64, error) {
		val, err := inner(r)
		if err != nil {
			return 0, err
		}
		return val * multiplier, nil
	}
}

// TODO: tests
func NewGaugeReader(g config.Gauge, r Reader) (func() float64, error) {
	var reader gaugeReader

	switch {
	case g.Byte != nil && g.Signed:
		reader = int8GaugeReader(*g.Byte, g.Address)
	case g.Byte != nil && !g.Signed:
		reader = uint8GaugeReader(*g.Byte, g.Address)
	case g.Byte == nil && !g.Signed:
		reader = uint16GaugeReader(g.Address)
	default:
		return nil, fmt.Errorf("unsupported configuration %+v", g)
	}

	if g.Multiplier != 0 {
		reader = multipliedGaugeReader(reader, float64(g.Multiplier))
	}
	return func() float64 {
		val, err := reader(r)
		if err != nil {
			logger.Printf("error reading %+v: %v", g, err)
			return -1
		}

		return val
	}, nil
}

type Reader interface {
	ReadUint16(address uint16) (uint16, error)
	Read(address uint16) ([]byte, error)
}

func gaugeOpts(g config.Gauge) prometheus.GaugeOpts {
	return prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      g.Name,
		Help:      g.Help,
	}
}

func Register(cc Reader, gauges []config.Gauge) error {
	for _, g := range gauges {

		gr, err := NewGaugeReader(g, cc)
		if err != nil {
			return fmt.Errorf("could not register gauges %w", err)
		}
		if len(g.States) > 0 {
			registerStateset(g, gr)
			continue
		}
		promauto.NewGaugeFunc(gaugeOpts(g), gr)
	}
	return nil
}

func registerStateset(g config.Gauge, r func() float64) {
	c := statesetCollector{
		prometheus.NewGaugeVec(gaugeOpts(g), []string{"state"}),
		r,
		g.States,
	}
	prometheus.MustRegister(c)
}

type statesetCollector struct {
	*prometheus.GaugeVec
	read   func() float64
	states []config.State
}

func (c statesetCollector) Collect(ch chan<- prometheus.Metric) {
	val := c.read()
	for _, s := range c.states {
		state := c.WithLabelValues(s.Name)
		// TODO: skip the float64 steps
		if val > -1 && math.Abs(s.Value-val) < 1 {
			state.Set(1)
		} else {
			state.Set(0)
		}
	}
	c.GaugeVec.Collect(ch)
}
