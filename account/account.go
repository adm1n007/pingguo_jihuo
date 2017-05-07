package account

import (
	"fmt"
	. "active_apple/ml/strings"

	"active_apple/utility"
)

const (
	ACCOUNT_UNREGISTERED    = 1
	ACCOUNT_BANNED          = 2
	ACCOUNT_EXISTS          = 3
	ACCOUNT_ACTIVATED       = 4
	ACCOUNT_UNACTIVATED     = 5
	ACCOUNT_EMAIL_NOT_FOUND = 6
	ACCOUNT_LOGIN_FAILED    = 7

	ACCOUNT_CANT_ACTIVATE       = 80
	ACCOUNT_ACTIVATE_FAILED     = 81
	ACCOUNT_NEED_RESET_PASSWORD = 82
	ACCOUNT_PASSWORD_ERROR      = 83
	ACCOUNT_SECURITY_BANNED     = 84

	ACCOUNT_BUY_LOCKED   = 90
	ACCOUNT_BUY_UNLOCKED = 91

	ACCOUNT_UPLOAD_DSID = 1000
)

const (
	DEFAULT_COUNTRY = 111
)

type AppleAccount struct {
	Id            int64
	UserName      string
	MailPassword  string
	ApplePassword string
	PopAddress    string
	Answer1       string
	Answer2       string
	Answer3       string
	QuestionText1 string
	QuestionText2 string
	QuestionText3 string
	Question1     string
	Question2     string
	Question3     string
	Verify        int64
	Birth         string
	RecoveryEmail string
	CreationTime  int64
	Allocation    string
	Status        int64
	OrigStatus    int64
	Dsid          int64

	Country			int64

	FailureCount            int
	FromActivator           bool
	InitFromDatabase        bool
	VerificationEmailResent bool
	CantVerify              bool
	PasswordVerified        bool
}

func (self AppleAccount) String2() string {
	return fmt.Sprintf("id = %v, user = %v", self.Id, self.UserName)
	return fmt.Sprintf(
		String("\n").Join([]String{
			"id             = %d",
			"userName       = %s",
			"mailPassword   = %s",
			"popAddress     = %s",
			"applePassword  = %s",
			"creationTime   = %d",
		}).String(),
		self.Id,
		self.UserName,
		self.MailPassword,
		self.PopAddress,
		self.ApplePassword,
		self.CreationTime,
	)
}

func (self *AppleAccount) CreateRandomPassword() string {
	if self.InitFromDatabase {
		return self.ApplePassword
	}

	self.ApplePassword = string(utility.GeneratePassword())
	return self.ApplePassword
}
