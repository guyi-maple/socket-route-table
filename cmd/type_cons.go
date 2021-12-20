package cmd

type CommandType int8

const (
	PING         CommandType = 0
	REPORT_ROUTE CommandType = 1
)
