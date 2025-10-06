package data

type BanUserTestCase struct {
	AdminUserID     uint64
	RequestedUserID uint64
	ExpectedAnswer  bool
}

var BanUserTestCases = []BanUserTestCase{
	{
		AdminUserID:     1,
		RequestedUserID: 5,
		ExpectedAnswer:  true,
	},
	{
		AdminUserID:     2,
		RequestedUserID: 6,
		ExpectedAnswer:  true,
	},
	{
		AdminUserID:     3,
		RequestedUserID: 7,
		ExpectedAnswer:  true,
	},
}
