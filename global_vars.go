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
	seDefaultIndex        = -iota
	seDependencyIssue     = -iota
	seAPINotSupport       = -iota
)

var serverErrorMessages = map[int]string{
	seNoError:             "No error",
	seInputParamNotValid:  "The input parameter is not valid",
	seInputSchemaNotValid: "The input schema is not valid",
	seInputJSONNotValid:   "The input json format is not valid",
	seDBResourceQuery:     "Internal DB query error occurs",
	seResourceNotFound:    "Could not find resource by PID",
	seResourceDuplicated:  "Could not create duplicated resource (Please check PID)",
	seResourceNotChange:   "Could not edit resource since no field has been changed",
	seDefaultIndex:        "Could not create new resource with default index",
	seDependencyIssue:     "Could not remove resource with resource dependency issue",
	seAPINotSupport:       "API is not supported yet.",
}
