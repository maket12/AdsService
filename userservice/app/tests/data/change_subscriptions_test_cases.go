package data

type ChangeSubscriptionsTestCase struct {
	UserID        uint64
	Subscriptions []string
}

var ChangeSubscriptionsTestCases = []ChangeSubscriptionsTestCase{
	{
		UserID:        1,
		Subscriptions: []string{"andrew223", "alex21"},
	},
	{
		UserID:        2,
		Subscriptions: []string{"shishi", "vladimir", "zaizai1"},
	},
	{
		UserID:        3,
		Subscriptions: []string{"keval34"},
	},
}
