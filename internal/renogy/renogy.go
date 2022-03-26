package renogy

import (
	_ "embed"
	"encoding/binary"
	"log"

	"github.com/goburrow/modbus"
	"github.com/lazy-electron-consulting/renogy-exporter/internal/config"
)

var logger = log.New(log.Writer(), "[renogy] ", log.Lmsgprefix|log.Flags())

type Renogy struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

func New(cfg *config.Modbus) (*Renogy, error) {
	handler := modbus.NewRTUClientHandler(cfg.Path)
	handler.BaudRate = cfg.BaudRate
	handler.DataBits = cfg.DataBits
	handler.Parity = cfg.Parity
	handler.StopBits = cfg.StopBits
	handler.SlaveId = byte(cfg.UnitID)
	handler.Timeout = cfg.Timeout
	if err := handler.Connect(); err != nil {
		logger.Printf("failed to connect %v\n", err)
		return nil, err
	}

	return &Renogy{
		handler: handler,
		client:  modbus.NewClient(handler),
	}, nil
}

func (r *Renogy) Close() error { return r.handler.Close() }

func (r *Renogy) ReadUint16(address uint16) (uint16, error) {
	raw, err := r.client.ReadHoldingRegisters(address, 1)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(raw), nil
}

func (r *Renogy) Read(address uint16) ([]byte, error) {
	return r.client.ReadHoldingRegisters(address, 1)
}
