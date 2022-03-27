# renogy-exporter

[Prometheus exporter](https://prometheus.io/docs/instrumenting/exporters/) for
Renogy charge controllers.

## Usage

- download a binary from the github releases or build from source using `make
  build`
- connect your computer to a renogy device, I'm using a renogy brand [RS485 to USB serial cable](https://www.renogy.com/rs485-to-usb-serial-cable/)
- run `./renogy-exporter config.yaml`
- see your metrics via http

See `renogy.yaml` for an example of the config syntax that works for me.

## Releasing

- merging to `main` will make new "latest" binaries
- push a tag w/ `v$SEMVER` to make a versioned binaries

## References

- [Modbus protocol](https://en.m.wikipedia.org/wiki/Modbus)
- [serial ports](https://en.wikipedia.org/wiki/Serial_port)
- [KyleJamesWalker/renogy_rover](https://github.com/KyleJamesWalker/renogy_rover)
  - python modbus integration
- [modbus_exporter](https://github.com/RichiH/modbus_exporter)
