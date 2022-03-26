package renogy

import (
	_ "embed"
	"encoding/binary"
	"log"
	"time"

	"github.com/goburrow/modbus"
)

var logger = log.New(log.Writer(), "[renogy] ", log.Lmsgprefix|log.Flags())

type Renogy struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

func New(path string) (*Renogy, error) {
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
