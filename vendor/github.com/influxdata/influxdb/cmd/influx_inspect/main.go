// The influx_inspect command displays detailed information about InfluxDB data files.
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/influxdata/influxdb/cmd"
	"github.com/influxdata/influxdb/cmd/influx_inspect/buildtsi"
	"github.com/influxdata/influxdb/cmd/influx_inspect/deletetsm"
	"github.com/influxdata/influxdb/cmd/influx_inspect/dumptsi"
	"github.com/influxdata/influxdb/cmd/influx_inspect/dumptsm"
	"github.com/influxdata/influxdb/cmd/influx_inspect/dumptsmwal"
	"github.com/influxdata/influxdb/cmd/influx_inspect/export"
	"github.com/influxdata/influxdb/cmd/influx_inspect/help"
	"github.com/influxdata/influxdb/cmd/influx_inspect/report"
	"github.com/influxdata/influxdb/cmd/influx_inspect/reporttsi"
	"github.com/influxdata/influxdb/cmd/influx_inspect/verify/seriesfile"
	"github.com/influxdata/influxdb/cmd/influx_inspect/verify/tsm"
	_ "github.com/influxdata/influxdb/tsdb/engine"
)

func main() {
	m := NewMain()
	if err := m.Run(os.Args[1:]...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main represents the program execution.
type Main struct {
	Logger *log.Logger

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		Logger: log.New(os.Stderr, "[influx_inspect] ", log.LstdFlags),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run determines and runs the command specified by the CLI args.
func (m *Main) Run(args ...string) error {
	name, args := cmd.ParseCommandName(args)

	// Extract name from args.
	switch name {
	case "", "help":
		if err := help.NewCommand().Run(args...); err != nil {
			return fmt.Errorf("help: %s", err)
		}
	case "deletetsm":
		name := deletetsm.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("deletetsm: %s", err)
		}
	case "dumptsi":
		name := dumptsi.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("dumptsi: %s", err)
		}
	case "dumptsmdev":
		fmt.Fprintf(m.Stderr, "warning: dumptsmdev is deprecated, use dumptsm instead.\n")
		fallthrough
	case "dumptsm":
		name := dumptsm.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("dumptsm: %s", err)
		}
	case "dumptsmwal":
		name := dumptsmwal.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("dumptsmwal: %s", err)
		}
	case "export":
		name := export.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("export: %s", err)
		}
	case "buildtsi":
		name := buildtsi.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("buildtsi: %s", err)
		}
	case "report":
		name := report.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("report: %s", err)
		}
	case "reporttsi":
		name := reporttsi.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("reporttsi: %s", err)
		}
	case "verify":
		name := tsm.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("verify: %s", err)
		}
	case "verify-seriesfile":
		name := seriesfile.NewCommand()
		if err := name.Run(args...); err != nil {
			return fmt.Errorf("verify-seriesfile: %s", err)
		}
	default:
		return fmt.Errorf(`unknown command "%s"`+"\n"+`Run 'influx_inspect help' for usage`+"\n\n", name)
	}

	return nil
}
