package main

const (
	seNoError             = -iota // 0, reserverd for no error
	seInputParamNotValid  = -iota
	seInputSchemaNotValid = -iota
	seInputJSONNotValid   = -iota
	seDBResourceQuery     = -iota
	seResourceNotFound    = -iota
	seResourceDuplicated  = -iota
	seResourceNotChange   = -iota
	seDependencyIssue     = -iota
	seUnresolvedError     = -iota
	seAPINotSupport       = -iota
)

var serverErrorMessages = map[int]string{
	seNoError:             "No error",
	seInputParamNotValid:  "Invalid input parameter",
	seInputSchemaNotValid: "Invalid input schema",
	seInputJSONNotValid:   "Invalid input JSON format",
	seDBResourceQuery:     "DB resource query error",
	seResourceNotFound:    "Resource not found",
	seResourceDuplicated:  "Resource duplicated",
	seResourceNotChange:   "Resource not changed",
	seDependencyIssue:     "Dependency not resolved",
	seUnresolvedError:     "Unresolved server error",
	seAPINotSupport:       "API not supported",
}
