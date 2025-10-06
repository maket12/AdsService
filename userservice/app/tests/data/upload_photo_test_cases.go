package data

type UploadPhotoTestCase struct {
	UserID         uint64
	Data           []byte
	FileName       string
	ContentType    string
	ExpectedAnswer string
}

var UploadPhotoTestCases = []UploadPhotoTestCase{
	{
		UserID:         1,
		Data:           []byte{0, 1, 1, 2, 34, 56},
		FileName:       "image.png",
		ContentType:    "image",
		ExpectedAnswer: "hex-string",
	},
	{
		UserID:         2,
		Data:           []byte{0, 8, 1, 2, 240, 52},
		FileName:       "new.jpeg",
		ContentType:    "image",
		ExpectedAnswer: "hex-string1",
	},
	{
		UserID:         3,
		Data:           []byte{0, 1, 1, 2, 34, 56},
		FileName:       "test123.webm",
		ContentType:    "image",
		ExpectedAnswer: "hex-string222",
	},
}

var UploadPhotoErrTestCases = []UploadPhotoTestCase{
	{
		UserID:      1,
		Data:        []byte{},
		FileName:    "image.png",
		ContentType: "image",
	},
	{
		UserID:      2,
		Data:        []byte{0, 8, 1, 2, 240, 52},
		FileName:    "",
		ContentType: "image",
	},
	{
		UserID:      3,
		Data:        []byte{0, 1, 1, 2, 34, 56},
		FileName:    "test123.webm",
		ContentType: "",
	},
}
