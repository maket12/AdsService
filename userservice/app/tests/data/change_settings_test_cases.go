package data

type ChangeSettingsTestCase struct {
	UserID               uint64
	NotificationsEnabled bool
}

var ChangeSettingsTestCases = []ChangeSettingsTestCase{
	{
		UserID:               1,
		NotificationsEnabled: true,
	},
	{
		UserID:               2,
		NotificationsEnabled: false,
	},
	{
		UserID:               3,
		NotificationsEnabled: false,
	},
}
