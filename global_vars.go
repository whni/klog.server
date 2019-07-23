package main

const (
	serrClientNotReady    = iota
	serrInputNotValid     = iota
	serrResourceNotFound  = iota
	serrDefaultIndex      = iota
	serrDuplicateResource = iota
	serResourceValidation = iota
	serrResourceNoChange  = iota
	serrDependencyIssue   = iota
	serrAPINotSupport     = iota
)

var serverErrorMessages = map[int]string{
	serrClientNotReady:    "The record client is not ready yet, please wait and retry again later",
	serrInputNotValid:     "The input parameter is not valid.",
	serrResourceNotFound:  "Could not find %s by PID!",
	serrDefaultIndex:      "Could not create new %s with default index",
	serrDuplicateResource: "Could not create duplicated %s",
	serResourceValidation: "Could not create %s with unvalidated input schema",
	serrResourceNoChange:  "Could not edit resource %s since no field has been changed",
	serrDependencyIssue:   "Could not remove %s with resource dependency issue",
	serrAPINotSupport:     "API is not supported yet.",
}
