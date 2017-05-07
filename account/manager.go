package account

import (
	. "fmt"
	. "active_apple/ml/strings"
	. "active_apple/ml/trace"

	_ "github.com/go-sql-driver/mysql"

	"active_apple/ml/channels"
	"active_apple/ml/io2"
	"active_apple/ml/logging/logger"
	"active_apple/ml/random"
	"active_apple/ml/sync2"

	"active_apple/globals"
	"active_apple/utility"

	"database/sql"
	"sync"
	"time"
)

const (
	UnactivatedManager   = 1
	UnregisteredManager  = 2
	EmptyDsidManager     = 3
	EmailNotFoundManager = 4
	LockedManager        = 5
)

type Manager struct {
	accountMap   map[int64]bool
	processedMap map[int64]bool
	accountList  *channels.InfiniteChannel
	commitTasks  *channels.InfiniteChannel

	accountProcessed int
	accountAvailable int
	managerType      int

	stop    bool
	running bool
	lock    *sync.Mutex
	event   *sync2.Event

	db        *sql.DB
	executor  func(db *sql.DB) (*sql.Rows, error)
	processor func(rows *sql.Rows) error
}

func NewManager(t int) *Manager {
	mgr := &Manager{
		accountMap:   map[int64]bool{},
		processedMap: map[int64]bool{},
		accountList:  channels.NewInfiniteChannel(),
		commitTasks:  channels.NewInfiniteChannel(),
		lock:         &sync.Mutex{},
		event:        sync2.NewEvent(),
		managerType:  t,
	}

	switch t {
	case UnactivatedManager:
		//logger.Debug("case UnactivatedManager")
		mgr.executor = mgr.getUnactivatedExecutor
		mgr.processor = mgr.getUnactivatedProcessor

	case EmailNotFoundManager:
		//logger.Debug("case EmailNotFoundManager")
		mgr.executor = mgr.getEmailNotFoundExecutor
		mgr.processor = mgr.getUnactivatedProcessor
/*
	case UnregisteredManager:
		//logger.Debug("case UnregisteredManager, debugging: %v", debugging)
		if debugging {
			mgr.executor = mgr.debugExecutor
			mgr.processor = mgr.debugProcessor

		} else {
			mgr.executor = mgr.getUnregisteredExecutor
			mgr.processor = mgr.getUnregisteredProcessor
		}

	case LockedManager:
		//logger.Debug("case LockedManager")
		mgr.executor = mgr.getLockedExecutor
		mgr.processor = mgr.getLockedProcessor
*/
	default:
		Raise(Sprintf("unknown manager type: %d", t))
	}

	return mgr
}

func (self *Manager) Close() {
	self.Stop()

	logger.Debug("stop exit")

	self.accountList.Close()
	logger.Debug("close1")

	self.commitTasks.Close()
	logger.Debug("close2")
}

func (self *Manager) isDebugging() bool {
	return self.managerType == UnregisteredManager && debugging
}

func (self *Manager) isAccountExists(acc *AppleAccount) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.accountMap[acc.Id]
}

func (self *Manager) removeAccount(acc *AppleAccount) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.accountAvailable--
	self.accountProcessed--
	delete(self.accountMap, acc.Id)
	delete(self.processedMap, acc.Id)
}

func (self *Manager) addAccount(acc *AppleAccount) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.accountMap[acc.Id] = true
	self.accountAvailable++
	self.accountList.In() <- acc
}

func (self *Manager) accountCount() int {
	self.lock.Lock()
	defer self.lock.Unlock()

	return len(self.accountMap)
}

func (self *Manager) setAccountProcessed(acc *AppleAccount) {
	self.lock.Lock()

	if self.processedMap[acc.Id] == false {
		self.processedMap[acc.Id] = true
		self.accountProcessed++
	}

	self.lock.Unlock()
}

var accid int64 = 0

func (self *Manager) debugExecutor(db *sql.DB) (*sql.Rows, error) {
	for _ = range [100]int{} {
		acc := self.getAccountLocal(false)
		if acc == nil {
			return nil, nil
		}

		acc.Status = ACCOUNT_UNREGISTERED
		acc.Id = accid

		accid++

		self.addAccount(acc)
	}

	return nil, nil
}

func (self *Manager) debugProcessor(rows *sql.Rows) error {
	return nil
}

func (self *Manager) getUnactivatedExecutor(db *sql.DB) (*sql.Rows, error) {
	return db.Query(
		`SELECT id, username, mailpass, popaddr, applepass, createtime, status FROM appstore_buydata_appleid
                WHERE id = 176321
                LIMIT 5000`,
	)
}

func (self *Manager) getEmailNotFoundExecutor(db *sql.DB) (*sql.Rows, error) {
	return db.Query(
		`SELECT id, username, mailpass, popaddr, applepass, createtime, status FROM appstore_buydata_appleid
                WHERE createtime is not NULL AND status = 6
                LIMIT 500`,
	)
}

func (self *Manager) getUnactivatedProcessor(rows *sql.Rows) error {
	for rows.Next() {
		acc := &AppleAccount{
			Country: DEFAULT_COUNTRY,
		}

		var userName sql.NullString
		var mailPassword sql.NullString
		var popAddress sql.NullString
		var applePassword sql.NullString

		err := rows.Scan(
			&acc.Id,
			&userName,
			&mailPassword,
			&popAddress,
			&applePassword,
			&acc.CreationTime,
			&acc.Status,
		)
		if err != nil {
			logger.Debug("err while scan: %v", err)
			continue
		}

		acc.UserName = userName.String
		acc.MailPassword = mailPassword.String
		acc.PopAddress = popAddress.String
		acc.ApplePassword = applePassword.String

		if self.isAccountExists(acc) {
			logger.Debug("id in accountMap(%d %d): %v", self.accountAvailable, self.accountProcessed, acc.UserName)
			continue
		}

		// acc.Status = ACCOUNT_UNACTIVATED
		acc.FromActivator = true

		logger.Debug("getUnactivatedProcessor, acc: %+v", acc)

		self.addAccount(acc)
	}

	return nil
}


func (self *Manager) handleCommit(request *CommitRequest) {
	var err error

	account := request.account

	logger.Debug("process: %v", account)

	exist := self.isAccountExists(account)

	if exist == false {
		logger.Debug("%v does not exist", account)
		return
	}

	self.removeAccount(account)

	if self.isDebugging() {
		return
	}

	status := account.Status
	query := ""

	switch status {
	case ACCOUNT_UPLOAD_DSID:
		query = Sprintf("UPDATE appstore_buydata_appleid SET dsid = %d WHERE id = %d", account.Dsid, account.Id)

	case account.OrigStatus:
		logger.Debug("invalid new status, ignore")
		return

	case ACCOUNT_BANNED,
		ACCOUNT_ACTIVATED,
		ACCOUNT_EMAIL_NOT_FOUND,
		ACCOUNT_CANT_ACTIVATE,
		ACCOUNT_NEED_RESET_PASSWORD,
		ACCOUNT_ACTIVATE_FAILED,
		ACCOUNT_PASSWORD_ERROR,
		ACCOUNT_SECURITY_BANNED,
		ACCOUNT_LOGIN_FAILED:

		query = Sprintf("UPDATE appstore_buydata_appleid SET status = %d WHERE id = %d", status, account.Id)

	case ACCOUNT_EXISTS:
		// from unlocker?
		if account.FromActivator {
			query = Sprintf("UPDATE appstore_buydata_appleid SET status = %d WHERE id = %d", status, account.Id)
			break
		}
		status = ACCOUNT_UNACTIVATED
		fallthrough

	case ACCOUNT_UNACTIVATED:
		if account.InitFromDatabase {
			query = string(String("\n").Join([]string{
				Sprintf(`UPDATE appstore_buydata_appleid SET`),
				Sprintf(`answerNO1_title    = "%s",`, account.QuestionText1),
				Sprintf(`answerNO2_title    = "%s",`, account.QuestionText2),
				Sprintf(`answerNO3_title    = "%s",`, account.QuestionText3),
				Sprintf(`createtime         = %d,`, account.CreationTime),
				Sprintf(`dsid               = %d,`, account.Dsid),
				Sprintf(`rescue             = "%s",`, account.RecoveryEmail),
				Sprintf(`country            = %d,`, account.Country),
				Sprintf(`status             = %d`, status),
				Sprintf(`WHERE id           = %d`, account.Id),
			}))

		} else {
			query = string(String("\n").Join([]string{
				Sprintf(`UPDATE appstore_buydata_appleid SET`),
				Sprintf(`applepass         = "%s",`, account.ApplePassword),
				Sprintf(`answer1            = "%s",`, account.Answer1),
				Sprintf(`answer2            = "%s",`, account.Answer2),
				Sprintf(`answer3            = "%s",`, account.Answer3),
				Sprintf(`answerNO1_title    = "%s",`, account.QuestionText1),
				Sprintf(`answerNO2_title    = "%s",`, account.QuestionText2),
				Sprintf(`answerNO3_title    = "%s",`, account.QuestionText3),
				Sprintf(`answerNO1          = "%s",`, account.Question1),
				Sprintf(`answerNO2          = "%s",`, account.Question2),
				Sprintf(`answerNO3          = "%s",`, account.Question3),
				Sprintf(`birth              = "%s",`, account.Birth),
				Sprintf(`createtime         = %d,`, account.CreationTime),
				Sprintf(`dsid               = %d,`, account.Dsid),
				Sprintf(`rescue             = "%s",`, account.RecoveryEmail),
				Sprintf(`country            = %d,`, account.Country),
				Sprintf(`status             = %d`, status),
				Sprintf(`WHERE id           = %d`, account.Id),
			}))
		}

	default:
		logger.Debug("unknown status")
		level := logger.Level()
		logger.SetLevel(99999)
		utility.Pause("debug me")
		logger.SetLevel(level)
	}

	logger.Debug("query = %s", query)
	_, err = self.db.Exec(query)
	if err != nil {
		logger.Debug("occur error while commit [%v]: %v", account, err)
	}
}

func (self *Manager) mainLoop2() {
	var err error

	if self.isDebugging() == false {
		dataSourceName := Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			globals.Preferences.Database.User,
			globals.Preferences.Database.Password,
			globals.Preferences.Database.Host,
			globals.Preferences.Database.Port,
			globals.Preferences.Database.Db,
		)

		self.db, err = sql.Open("mysql", dataSourceName)
		logger.Debug("connect db, err: %v", err)
		if err != nil {
			return
		}

		defer self.db.Close()
	}

	commitTaskHandlerExited := sync.WaitGroup{}
	commitTaskHandlerExited.Add(1)
	dbLock := &sync.Mutex{}

	go func() {
		defer commitTaskHandlerExited.Done()

		for self.stop == false {
			select {
			case req := <-self.commitTasks.Out():
				r := req.(*CommitRequest)
				dbLock.Lock()
				exp := Try(func() {
					self.handleCommit(r)
				})
				dbLock.Unlock()

				if exp != nil {
					logger.Debug("handleCommit exception: %v", exp)
				}

				r.Done()

			default:
				time.Sleep(time.Millisecond)
			}
		}
	}()

	logger.Debug("enter account loop")
	for self.stop == false {
		var maxWorkers int
		switch self.managerType {
		case UnregisteredManager:
			maxWorkers = globals.Preferences.MaxRegisterWorkers
		case UnactivatedManager, EmailNotFoundManager:
			maxWorkers = globals.Preferences.MaxActivatorWorkers
		case LockedManager:
			maxWorkers = globals.Preferences.MaxActivatorWorkers
		}

		logger.Debug("account loop, self.accountCount(): %v", self.accountCount())
		// if self.accountAvailable - self.accountProcessed > maxWorkers {
		if self.accountAvailable >= maxWorkers {
			logger.Debug("[%d] accountList not empty. accountAvailable: %d, accountProcessed: %d", self.managerType, self.accountAvailable, self.accountProcessed)
			if self.isDebugging() == false {
				self.db.Ping()
			}

			time.Sleep(time.Second * 10)
			continue
		}

		numberOfAccounts := self.accountCount()

		exp := Try(func() {
			dbLock.Lock()
			defer dbLock.Unlock()

			rows, err := self.executor(self.db)

			RaiseIf(err)
			//万恶的起源
			err = self.processor(rows)
			RaiseIf(err)
		})

		if exp != nil {
			logger.Debug("query account error: %v", exp)
		}

		// logger.Debug("executing")
		// rows, err := self.executor(self.db)
		// logger.Debug("execute: %v", err)
		// if err != nil {
		//     time.Sleep(time.Second)
		//     continue
		// }

		// origLength := self.accountCount()

		// logger.Debug("process")
		// err = self.processor(rows)
		// logger.Debug("process end: %v", err)

		// logger.Debug("%d", self.accountCount())

		if self.accountCount() == numberOfAccounts {
			time.Sleep(10 * time.Second)
		}
	}

	commitTaskHandlerExited.Wait()
}

func (self *Manager) mainLoop() {
	exp := Try(self.mainLoop2)

	if exp != nil {
		logger.Debug("mainLoop2 error: %v", exp)
	}

	logger.Debug("mainLoop2 exit")
	self.event.Signal()
}

func (self *Manager) Start() {
	if self.running {
		return
	}

	self.running = true
	self.stop = false
	go self.mainLoop()
}

func (self *Manager) Stop() {
	if self.running == false {
		return
	}

	self.stop = true
	self.event.Wait()

	self.stop = false
	self.running = false
}

var emailIndex = 0
var localaccs = []*AppleAccount{}

func init() {
	return
	if debugging == false {
		return
	}

	var id int64 = 1

	for _, mail := range io2.ReadLines("D:\\Desktop\\mails.txt") {
		if mail.IsEmpty() {
			continue
		}

		r := mail.Split(" ", 1)

		account := &AppleAccount{
			Id:            id,
			UserName:      r[0].String(),
			MailPassword:  r[1].String(),
			ApplePassword: string(utility.GeneratePassword()),
			Status:        ACCOUNT_UNREGISTERED,
			Country:       DEFAULT_COUNTRY,
		}

		id++

		localaccs = append(localaccs, account)
	}
}

func (self *Manager) getAccountLocal(block bool) *AppleAccount {
	emails := []String{
		`@outlook.com`,
		`@hotmail.com`,
		`@162.com`,
		`@qq.com`,
		`@foxmail.com`,
		`@gmail.com`,
		`@microsoft.com`,
		`@google.com`,
		`@sina.com`,
		// `@163.com`,
	}

	// email := Sprintf("%s%d%s", utility.GeneratePinyin(), random.IntRange(1000, 100000), random.Choice(emails).(String))
	email := Sprintf("%s%d%s", utility.GeneratePinyin(), random.IntRange(1000, 100000), emails[emailIndex%len(emails)])
	emailIndex++

	// logger.Debug("email = %v", email)
	// email = utility.GeneratePassword() + "@" + utility.GeneratePassword()[:random.IntRange(4, 8)] + ".com"

	account := &AppleAccount{
		UserName:  string(email),
		Status:    ACCOUNT_UNREGISTERED,
		Country:   DEFAULT_COUNTRY,
		Question1: Sprintf("%d", random.IntRange(130, 136)),
		Question2: Sprintf("%d", random.IntRange(136, 142)),
		Question3: Sprintf("%d", random.IntRange(142, 148)),
		Answer1:   utility.GenerateAnswer().String(),
		Answer2:   utility.GenerateAnswer().String(),
		Answer3:   utility.GenerateAnswer().String(),
		Birth:     Sprintf("%04d-%d-%d", random.IntRange(1970, 2001), random.IntRange(1, 13), random.IntRange(1, 27)),
	}

	// account.Answer1 = utility.GeneratePinyin().String()
	// account.Answer2 = utility.GeneratePinyin().String()
	// account.Answer3 = utility.GeneratePinyin().String()

	account.CreateRandomPassword()
	account.InitFromDatabase = true

	// logger.Debug("local account: %v %v", account.UserName, account.ApplePassword)

	return account
}

func (self *Manager) GetAccount(block bool) *AppleAccount {
	for globals.Exiting == false {
		select {
		case v := <-self.accountList.Out():
			acc := v.(*AppleAccount)
			//logger.Debug("GetAccount, got account, acc: %+v", acc)
			self.setAccountProcessed(acc)
			return acc

		default:
			if block == false {
				return nil
			}
		}

		//logger.Debug("GetAccount, no account, waiting")
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}

func (self *Manager) ReleaseAccount(account *AppleAccount) {
	self.accountList.In() <- account
}

func (self *Manager) Commit(account *AppleAccount) {
	request := NewRequest(account)
	self.commitTasks.In() <- request
	request.Wait()
}

const debugging = globals.Debugging
