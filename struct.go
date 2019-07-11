package main

// Description of registered person in OA system
const (
	PERSON_TYPE_EMPLOYEE = "Employee"
	PERSON_TYPE_GUEST    = "Guest"
	PERSON_TYPE_UNKNOWN  = "Unknown"
)

type Person struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Department      string   `json:"department"`
	Role            string   `json:"role"`
	PersonType      string   `json:"persontype"`
	AIClusterEnable bool     `json:"aiclusterenable"`
	Images          []string `json:"images"`
	AIImages        []string `json:"aiimages"`
}

type CamLocationParams struct {
	AreaId       int
	FloorId      int
	Buildingd    int
	SiteId       int
	AreaName     string
	FloorName    string
	BuildingName string
	SiteName     string
	CamPosition  []float64
}

type MailServer struct {
	ID           int    `json:"id"`
	Description  string `json:"description"`
	ServerName   string `json:"server_name"`
	ServerPort   int    `json:"server_port"`
	ConnSecurity string `json:"conn_security"`
	AuthMethod   string `json:"auth_method"`
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
}

type MailServers []MailServer

type MailSender struct {
	ID           int    `json:"id"`
	DisplayName  string `json:"display_name"`
	EmailAddress string `json:"email_address"`
}

type MailSenders []MailSender

type SmsProvider struct {
	ID           int    `json:"id"`
	ProviderName string `json:"provider_name"`
	UserId       string `json:"user_id"`
	AuthToken    string `json:"auth_token"`
	FromPhoneNum string `json:"from_phone_num"`
}

type SmsProviders []SmsProvider

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Administrator struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	EmailAddress string `json:"email_address"`
	PhoneNum     string `json:"phone_num"`
	MsgMethod    string `json:"msg_method"` // email or sms or both or none
	Permission   int    `json:"permission"` // Administrator_level, Only_for_read, only_for_recording
}

type Administrators []Administrator

type Version struct {
	Major uint16
	Minor uint16
	Build uint16
}

type LogController_Cmd struct {
	Action string `json:"action"`    //level, module, format, output, default
	Level  string `json:"level"`     //trace, debug, info, warn, error, fatal, panic
	Module string `json:"module"`    //all or individual module name
	Format string `json:"format"`    //json or text
	Output string `json:"output"`    //STDOUT, FILE, BOTH
	Op     string `json:"operation"` //ONLY used for module: set or clear
	Rate   int    `json:"rate"`      //for Profiling block or cpu rate
	//BELOW are for TEXT format ONLY
	ForceColor     bool `json:"forcecolor"`     //true or false
	DisableColor   bool `json:"disablecolor"`   //true or false
	FullTimestamp  bool `json:"fulltimestamp"`  //true or false
	DisableSorting bool `json:"disablesorting"` //true or false
}

// response OK 200, response Error 409
type Response struct {
	Status  int
	payload Body
}

type Body struct {
	Message string      `json:"message"`
	Percent int64       `json:"percent"`
	Data    interface{} `json:"data"`
}

type CameraDescriptor struct {
	CameraId     uint32
	CameraName   [128]byte
	CameraStatus uint32
}

type CamDescs []CameraDescriptor
