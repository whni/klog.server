package main

const (
	seNoError             = iota // 0, reserverd for no error
	seInputParamNotValid  = iota
	seInputSchemaNotValid = iota
	seDBResourceQuery     = iota
	seResourceNotFound    = iota
	seResourceDuplicated  = iota
	seResourceNotChange   = iota
	seDefaultIndex        = iota
	seDependencyIssue     = iota
	seAPINotSupport       = iota
)

var serverErrorMessages = map[int]string{
	seNoError:             "No error",
	seInputParamNotValid:  "The input parameter is not valid",
	seInputSchemaNotValid: "The input schema is not valid",
	seDBResourceQuery:     "Internal DB query error occurs",
	seResourceNotFound:    "Could not find <%s> by PID",
	seResourceDuplicated:  "Could not create duplicated <%s>",
	seResourceNotChange:   "Could not edit resource <%s> since no field has been changed",
	seDefaultIndex:        "Could not create new <%s> with default index",
	seDependencyIssue:     "Could not remove <%s> with resource dependency issue",
	seAPINotSupport:       "API is not supported yet.",
}
