package data

type AddProfileTestCase struct {
	UserID uint64
	Name   string
	Phone  string
}

var AddProfileTestCases = []AddProfileTestCase{
	{
		UserID: 1,
		Name:   "ShiShi",
		Phone:  "+79137845412",
	},
	{
		UserID: 2,
		Name:   "Vladimir",
		Phone:  "+79137848954",
	},
	{
		UserID: 3,
		Name:   "Andrew",
		Phone:  "+79549857814",
	},
}
