package main

const (
	configFileName = "config.yml"
	configLogDebug = "debug"
	configLogInfo  = "info"
	configLogWarn  = "warn"
	configLogError = "error"
	configLogFatal = "fatal"
)

const (
	amiEventInUse    = "AgentConnect"
	amiEventNotInUse = "AgentComplete"
	amiFieldMember   = "MemberName"
	amiCommand       = "devstate change Custom:%s %s"
	amiFlushDevstate = "database deltree CustomDevstate"
	devNotInuse      = "NOT_INUSE"
	devInuse         = "INUSE"
)
