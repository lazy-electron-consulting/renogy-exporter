package renology

import (
	_ "embed"
	"encoding/binary"
	"log"
	"time"

	"github.com/goburrow/modbus"
)

var logger = log.New(log.Writer(), "[renology] ", log.Lmsgprefix|log.Flags())

type Renology struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

func New(path string) (*Renology, error) {
	handler := modbus.NewRTUClientHandler(path)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second
	if err := handler.Connect(); err != nil {
		logger.Printf("failed to connect %v\n", err)
		return nil, err
	}

	return &Renology{
		handler: handler,
		client:  modbus.NewClient(handler),
	}, nil
}

func (r *Renology) Close() error { return r.handler.Close() }

func (r *Renology) ReadUint16(address uint16) (uint16, error) {
	raw, err := r.client.ReadHoldingRegisters(address, 1)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(raw), nil
}

func (r *Renology) Read(address uint16) ([]byte, error) {
	return r.client.ReadHoldingRegisters(address, 1)
}
