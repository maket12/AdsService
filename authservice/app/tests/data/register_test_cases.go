package data

type RegisterTestCase struct {
	Email                string
	Password             string
	Role                 string
	ExpectedAccessToken  string
	ExpectedRefreshToken string
}

var RegisterTestCases = []RegisterTestCase{
	{
		Email:                "vova1234@gmail.com",
		Password:             "password123",
		Role:                 "user",
		ExpectedAccessToken:  "access-xxx",
		ExpectedRefreshToken: "refresh-yyy",
	},
	{
		Email:                "alex451@mail.ru",
		Password:             "password178",
		Role:                 "user",
		ExpectedAccessToken:  "access-zzz",
		ExpectedRefreshToken: "refresh-zzz",
	},
	{
		Email:                "newemail@yahoo.com",
		Password:             "newpass",
		Role:                 "user",
		ExpectedAccessToken:  "access-nnn",
		ExpectedRefreshToken: "refresh-nnn",
	},
}
