package main

const (
	seNoError             = -iota // 0, reserverd for no error
	seInputParamNotValid  = -iota
	seInputSchemaNotValid = -iota
	seInputJSONNotValid   = -iota
	seInputBSONNotValid   = -iota
	seDBResourceQuery     = -iota
	seResourceNotFound    = -iota
	seResourceNotMatched  = -iota
	seResourceDuplicated  = -iota
	seResourceConflict    = -iota
	seResourceExpired     = -iota
	seResourceNotChange   = -iota
	seDependencyIssue     = -iota
	seCloudOpsError       = -iota
	seUnresolvedError     = -iota
	seAPINotSupport       = -iota
)

var serverErrorMessages = map[int]string{
	seNoError:             "NO_ERROR",
	seInputParamNotValid:  "INVALID_INPUT_PARAMS",
	seInputSchemaNotValid: "INVALID_INPUT_SCHEMA",
	seInputJSONNotValid:   "INVALID_INPUT_JSON",
	seInputBSONNotValid:   "INVALID_INPUT_BSON",
	seDBResourceQuery:     "DB_QUERY_ERROR",
	seResourceNotFound:    "RESOURCE_NOT_FOUND",
	seResourceNotMatched:  "RESOURCE_NOT_MATCHED",
	seResourceDuplicated:  "RESOURCE_DUPLICATED",
	seResourceConflict:    "RESOURCE_CONFLICT",
	seResourceExpired:     "RESOURCE_EXPIRED",
	seResourceNotChange:   "RESOURCE_NOT_CHANGED",
	seDependencyIssue:     "DEPENDENCY_UNRESOLVED",
	seCloudOpsError:       "CLOUD_OPS_ERROR",
	seUnresolvedError:     "UNRESOLVED_SERVER_ERROR",
	seAPINotSupport:       "API_NOT_SUPPORTED",
}
