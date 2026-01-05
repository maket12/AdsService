package data

type ValidateTokenTestCase struct {
	AccessToken    string
	ExpectedAnswer bool
}

var ValidateTokenTestCases = []ValidateTokenTestCase{
	{
		AccessToken:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiYWNjZXNzIiwidXNlcl9pZCI6MSwiZW1haWwiOiJ2b3ZhMTIzNEBnbWFpbC5jb20iLCJyb2xlIjoidXNlciIsInN1YiI6IjEiLCJleHAiOjE3NTk1MDg0NjQsImlhdCI6MTc1OTUwNzU2NCwianRpIjoiNjc0NWMyZGUtZjFiMC00MjNiLTlmN2YtMmZlODhiOWQxYzE2In0.6pnS5MhBv9-UjOPxVbarHytwvRchyrzdtqasRnVTt7w",
		ExpectedAnswer: true,
	},
	{
		AccessToken:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiYWNjZXNzIiwidXNlcl9pZCI6MiwiZW1haWwiOiJhbGV4NDUxQG1haWwucnUiLCJyb2xlIjoiYWRtaW4iLCJzdWIiOiIyIiwiZXhwIjoxNzU5NTA4NTI5LCJpYXQiOjE3NTk1MDc2MjksImp0aSI6IjViMDUzZTljLWY5NDgtNDhjNS05OTEwLWM2Y2M4ZWE1OTlkYyJ9.hOd5W4UJ-QxBzA2PEeoBRPmFdcLS_MmpIHc49atx21Q",
		ExpectedAnswer: true,
	},
	{
		AccessToken:    "kinda-token",
		ExpectedAnswer: true,
	},
}
