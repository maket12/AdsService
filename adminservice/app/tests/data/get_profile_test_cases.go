package data

type GetProfileTestCase struct {
	AdminUserID     uint64
	RequestedUserID uint64
}

var GetProfileTestCases = []GetProfileTestCase{
	{
		AdminUserID:     1,
		RequestedUserID: 5,
	},
	{
		AdminUserID:     2,
		RequestedUserID: 6,
	},
	{
		AdminUserID:     3,
		RequestedUserID: 7,
	},
}
