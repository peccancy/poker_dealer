package logs

import (
	"github.com/peccancy/chassi/log"
)

var (
	EntryFactory log.LogEntryFactory
)

func init() {
	EntryFactory = log.NewLogEntryFactory("tickets-client-api")
	Error = EntryFactory.MakeEntry(log.ERROR, 2000, "error")
	ServicePanicked = EntryFactory.MakeEntry(log.ERROR, 2001, "service panicked")
	Info = EntryFactory.MakeEntry(log.INFO, 3000, "info")
	Debug = EntryFactory.MakeEntry(log.DEBUG, 4000, "debug")
}

var (
	Error           log.EntryFunc
	ServicePanicked log.EntryFunc
	Info            log.EntryFunc
	Debug           log.EntryFunc
)
