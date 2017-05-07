package globals

import (
	"active_apple/ml/encoding/json"
	"active_apple/ml/os2"
	. "active_apple/ml/strings"
	. "active_apple/ml/trace"
	"path/filepath"
	"time"
)

const Debugging = false
const UseFiddler = false

type pref struct {
	MaxRegisterWorkers  int `json:"MaxRegisterWorkers,omitempty"`
	MaxActivatorWorkers int `json:"MaxActivatorWorkers,omitempty"`

	MaxProxys int `json:"MaxProxys,omitempty"`

	TotalProcess   int `json:"TotalProcess,omitempty"`
	CurrentProcess int

	Mail               String `json:"Mail,omitempty"`
	UseSocks5Proxy     bool   `json:"UseSocks5Proxy,omitempty"`
	MaxProxyReuseTimes int    `json:"MaxProxyReuseTimes,omitempty"`
	SaveLogToFile      bool   `json:"SaveLogToFile,omitempty"`

	Register struct {
		VerifyApppleIdDelay    time.Duration `json:"VerifyApppleIdDelay,omitempty"`
		VerifyApppleIdDelayiOS time.Duration `json:"VerifyApppleIdDelayiOS,omitempty"`
		RegisterAlive          time.Duration `json:"RegisterAlive,omitempty"`
		RetryOnFailure         int           `json:"RetryOnFailure,omitempty"`
	}

	Activator struct {
		ActivateDelay  time.Duration `json:"ActivateDelay,omitempty"`
		ActivatorAlive time.Duration `json:"ActivatorAlive,omitempty"`
	}

	Proxy struct {
		Port         int    `json:"Port,omitempty"`
		User         String `json:"User,omitempty"`
		Password     String `json:"Password,omitempty"`
		MinimumAlive int    `json:"MinimumAlive,omitempty"`
		ProxyUrl     String `json:"ProxyUrl,omitempty"`
	}

	ProxyServer struct {
		Host  String `json:"Host,omitempty"`
		Url   String `json:"Url,emitempty"`
		Port  int    `json:"Port,omitempty"`
		Port2 int    `json:"Port2,omitempty"`
	}

	Database struct {
		Host     String `json:"Host,omitempty"`
		Port     int    `json:"Port,omitempty"`
		User     String `json:"User,omitempty"`
		Password String `json:"Password,omitempty"`
		Db       String `json:"db,omitempty"`
	}
}

func (self *pref) LoadFile() {
	if Try(func() { json.LoadFile(filepath.Join(os2.ExecutablePath(), "Preferences.json"), self) }) != nil {
		json.LoadFile(`D:\Desktop\Source\GoProject\src\AppleIdRegister\Preferences.json`, self)
	}

	switch {
	case self.CurrentProcess <= 0,
		self.CurrentProcess > self.TotalProcess:
		Raisef("incorrect Preference: CurrentProcess = %d TotalProcess = %d", self.CurrentProcess, self.TotalProcess)
	}

	self.Register.RegisterAlive *= time.Second
	self.Register.VerifyApppleIdDelay *= time.Second
	self.Register.VerifyApppleIdDelayiOS *= time.Second
	self.Activator.ActivatorAlive *= time.Second
	self.Activator.ActivateDelay *= time.Minute

	if self.MaxProxyReuseTimes == 0 {
		self.MaxProxyReuseTimes = 3
	}
}

var Preferences pref

func init() {
	Preferences.LoadFile()
}
