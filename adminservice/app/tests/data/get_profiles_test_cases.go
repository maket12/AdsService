package data

type GetProfilesTestCase struct {
	AdminUserID uint64
	Limit       uint32
	Offset      uint32
}

var GetProfilesTestCases = []GetProfilesTestCase{
	{
		AdminUserID: 1,
		Limit:       10,
		Offset:      0,
	},
	{
		AdminUserID: 2,
		Limit:       10,
		Offset:      5,
	},
	{
		AdminUserID: 3,
		Limit:       5,
		Offset:      10,
	},
}
