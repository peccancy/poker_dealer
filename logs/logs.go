package logs

import (
	logging "github.com/peccancy/chassi/log"
	"log"
	"os"
)

var (
	EntryFactory logging.LogEntryFactory
)

func init() {
	StartService = EntryFactory.MakeEntry(logging.ERROR, 2001, "start service error")
	EntryFactory = logging.NewLogEntryFactory("tickets-client-api")
	Error = EntryFactory.MakeEntry(logging.ERROR, 2000, "error")
	ServicePanicked = EntryFactory.MakeEntry(logging.ERROR, 2001, "service panicked")
	Info = EntryFactory.MakeEntry(logging.INFO, 3000, "info")
	Debug = EntryFactory.MakeEntry(logging.DEBUG, 4000, "debug")
	GrpcServerInitFailure = EntryFactory.MakeEntry(logging.ERROR, 2003, "grpc server init failure")

	MongoDB = &MongoLogger{Logger: *log.New(os.Stdout, "mongodb ", log.LstdFlags)}
}

var (
	Error                 logging.EntryFunc
	ServicePanicked       logging.EntryFunc
	Info                  logging.EntryFunc
	Debug                 logging.EntryFunc
	StartService          logging.EntryFunc
	GrpcServerInitFailure logging.EntryFunc

	MongoDB *MongoLogger
)

type MongoLogger struct {
	log.Logger
}
