package main

import (
	"bufio"
	"net"
	"strconv"
	"sync"
)

type OrderOnWall struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	FrcId    uint64   `json:"frcid,string"` // return zero to UI, when this order is emtpy
	CameraId uint64   `json:"cameraid,string"`
	Streams  []Stream `json:"stream"`
}

type CameraWall struct {
	ID    int            `json:"id"`
	Name  string         `json:"name"`
	Order []*OrderOnWall `json:"order"`
}

type CameraWalls []CameraWall

type Sequencer struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Interval int            `json:"interval"`
	Order    []*OrderOnWall `json:"order"`
}

type Sequencers []Sequencer

type Site struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Building struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	SiteId    int       `json:"siteid"`
	Coorinate []float64 `json:"coordinate"`
}

type Floor struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	BuildingId int     `json:"buildingid"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	PlanImage  string  `json:"planimage"`
	//	Position []float64 `json:"position"`
}

type Area struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	FloorId  int         `json:"floorid"`
	Position [][]float64 `json:"position"`
}

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

// TODO Used in demo, Should be removed later
type PersonOA2 struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Department string   `json:"department"`
	Role       string   `json:"role"`
	ImageSrc   int      `json:"image_src"` // 0: manually setup, 1: extracted from video stream
	Image      []string `json:"image"`
}

const (
	FaceDetectionTypeNormal  = "normal"
	FaceDetectionTypeCluster = "cluster"
)

type FaceDetection struct {
	ID              int         `json:"id"`
	Type            string      `json:"type"`
	PersonId        int         `json:"personid"`
	PolicyId        int         `json:"policyid"`
	TaskId          int         `json:"taskid"`
	FrcId           uint64      `json:"frcid,string"`
	DevicedId       uint64      `json:"deviceid,string"`
	DeviceType      string      `json:"devicetype"`
	DevicePosition  []float64   `json:"deviceposition"`
	AreaId          int         `json:"areaid"`
	FloorId         int         `json:"floorid"`
	BuildingId      int         `json:"buildingid"`
	SiteId          int         `json:"siteid"`
	AreaName        string      `json:"areaname"`
	FloorName       string      `json:"floorname"`
	BuildingName    string      `json:"buildingname"`
	SiteName        string      `json:"sitename"`
	Image           []string    `json:"image"`
	ImageDeleteFlag []bool      `json:"imagedeleteflag"`
	Time            []uint64    `json:"time"`
	Position        [][]float64 `json:"position"`
	Extra_GUI       []string    `json:"extra_gui"`
}

type FaceImageTransfer struct {
	FaceDetectionId int      `json:"facedetectionid"`
	PersonId        int      `json:"personid"`
	Images          []string `json:"images"`
}

type FaceImage struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
}

type PersonDetection struct {
	ID       int         `json:"id"`
	CameraId int         `json:"cameraid"`
	Time     uint64      `json:"time"`
	Image    string      `json:"image"`
	Position [][]float64 `json:"position"`
}

type CameraEventInfo struct {
	PersonId   int
	PolicyId   int
	FrcId      uint64
	DeviceId   uint64
	AreaId     int
	FloorId    int
	BuildingId int
	SiteId     int
	Event      string
	Time       uint64
}

type AlertLog struct {
	ID         int      `json:"id"`
	RuleId     int      `json:"ruleid"`
	FrcId      uint64   `json:"frcid,string"`
	DeviceId   uint64   `json:"deviceid,string"`
	SiteId     int      `json:"siteid"`
	BuildingId int      `json:"buildingid"`
	FloorId    int      `json:"floorid"`
	AreaId     int      `json:"areaid"`
	PersonId   int      `json:"personid"`
	Event      string   `json:"event"`  //TODO need constrains?
	Action     []string `json:"action"` // TODO need constrians?
	Time       uint64   `json:"time"`   // UnixMillis, uint64 (now.UnixNano() / 1000000)
}

type BuildingCameraStatusS []*BuildingCameraStatus
type BuildingCameraStatus struct {
	BuildingId   int    `json:buildingid`
	BuildingName string `json:buildingname`
	ActiveCamera int    `json:activecamera`
	DownCamera   int    `json:downcamera`
}

type FloorCameraStatusS []*FloorCameraStatus
type FloorCameraStatus struct {
	FloorId      int    `json:floorid`
	FloorName    string `json:floorname`
	ActiveCamera int    `json:activecamera`
	DownCamera   int    `json:downcamera`
}

type PeopleEventNumberS []*PeopleEventNumber
type PeopleEventNumber struct {
	FloorId        int    `json:"floorid"`
	FloorName      string `json:"floorname"`
	DetectedPerson int    `json:"detectedperson"`
	UnknownPerson  int    `json:"unknownperson"`
	alerts         int    `json:"alerts"`
}

const (
	FilterTimeRangeTypeDay         = "day"
	FilterTimeRangeTypeWeek        = "week"
	FilterTimeRangeTypeMonth       = "month"
	FilterTimeRangeTypeYear        = "year"
	FilterTimeRangeDayMilliSeconds = 86400000
)

var FilterTimeRangeMap map[string]uint64 = map[string]uint64{
	FilterTimeRangeTypeDay:   FilterTimeRangeDayMilliSeconds,
	FilterTimeRangeTypeWeek:  FilterTimeRangeDayMilliSeconds * 7,
	FilterTimeRangeTypeMonth: FilterTimeRangeDayMilliSeconds * 30,
	FilterTimeRangeTypeYear:  FilterTimeRangeDayMilliSeconds * 365,
}

type Filter_GUI struct {
	DeviceIds      []string `json:"deviceids"`
	FrcId          uint64   `json:"frcid,string"`
	PersonId       int      `json:"personid"`
	TimeStart      uint64   `json:"timestart"`
	TimeStop       uint64   `json:"timestop"`
	TimeRange      string   `json:"timerange"`
	FacedetectType string   `json:"facedetecttype"`
}

type Filter struct {
	DeviceId       uint64 `json:"deviceid,string"`
	FrcId          uint64 `json:"frcid,string"`
	PersonId       int    `json:"personid"` // -1: known persons, 0: all persons, >0: specific person
	TimeStart      uint64 `json:"timestart"`
	TimeStop       uint64 `json:"timestop"`
	TimeRange      string `json:"timerange"`
	FacedetectType string `json:"facedetecttype"`
}

type HistoryCmd struct {
	TimeStart uint64   `json:"timestart"`
	Interval  uint64   `json:"interval"`
	Count     uint64   `json:"count"`
	DeviceIds []string `json:"deviceids"`
}

type PointHistory struct {
	TimeStart uint64 `json:"timestart"`
	TimeEnd   uint64 `json:"timeend"`
	Known     uint32 `json:"known"`
	Unknown   uint32 `json:"unknown"`
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

type PlaybackCmd struct {
	FrcId     uint64 `json:"frcid,string"`
	DeviceId  uint64 `json:"cameraid,string"`
	TimeStamp uint64 `json:"timestamp"` // 1546555875000 in milliseconds
	EventType uint32 `json:"eventtype"` // no validation yet TODO
	UuId      uint32 `json:"uuid"`
}

const (
	VIDEO_CMD_PLAY         uint32 = 1
	VIDEO_CMD_SEEK         uint32 = 2
	VIDEO_CMD_PAUSE        uint32 = 3
	VIDEO_CMD_STOP         uint32 = 4
	VIDEO_CMD_REPEAT_FRAME uint32 = 5
	VIDEO_CMD_LOOPBACK     uint32 = 8
)

const (
	PLAYBACK_THREAD_IDLE     int64 = 0
	PLAYBACK_THREAD_PLAY     int64 = 1
	PLAYBACK_THREAD_PAUSE    int64 = 2
	PLAYBACK_THREAD_SEEK     int64 = 3
	PLAYBACK_THREAD_TEARDOWN int64 = 4
	PLAYBACK_THREAD_FINISHED int64 = 5
	PLAYBACK_THREAD_LOOP     int64 = 6
)

type VideoCtrlCmd struct {
	FrcId             uint64 `json:"frcid,string"`
	DeviceId          uint64 `json:"cameraid,string"`
	Cmd               uint32 `json:"action"`
	SpeedOrDurationMs uint32 `json:"speed_or_duration"` // 50-300
	TimeStamp         uint64 `json:"timestamp"`
	EventType         uint32 `json:"eventtype"`
	UuId              uint32 `json:"uuid"`
}

type PlaybackInfo struct {
	ID      int    `json:"id"`
	StartTS uint64 `json:"start"`
	Length  uint64 `json:"length"`
}

type PlaybackFileParams struct {
	fileList []string
	tsStart  []uint64
	tsStop   []uint64
	fileCnt  uint32
}

type PlaybackJSONList struct {
	FileURL   []string `json:"file_url"`
	FileStart []uint64 `json:"file_start"`
	FileStop  []uint64 `json:"file_stop"`
}

type PersonPosition struct {
	FloorId  int         `json:"floorid"`
	PersonId int         `json:"personid"`
	Times    []uint64    `josn:"times"`
	Position [][]float64 `json:"position"`
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

const (
	Administrator_level = 1
	Only_For_Read       = 2
	Only_For_Recording  = 3
)

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

type Recorder struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	IpAddr    string `json:"ipaddr"`
	Port      int    `json:"port"`
	SSL       bool   `json:"ssl"`
	httpsPort int    `json:"httpsport"`
	Enabled   bool   `json:"enabled"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Status    int    `json:"status,omitempty"`
	FrcId     uint64 `json:"frcid,omitempty,string"`
	Version   string `json:"version,omitempty"`
}

type Recorders []Recorder

type ChCfg struct {
	frc    Recorder
	action int
}

const (
	CONN_FRC_TCP_CONNECTED           uint32 = 0x00000001
	CONN_FRC_AUTHORIZED              uint32 = 0x00000002
	CONN_FRC_READY                   uint32 = 0x00000003
	CONN_FRC_TCP_FAILED              uint32 = 0x00000004
	CONN_FRC_AUTH_FAILED             uint32 = 0x00000005
	CONN_FRC_RCVED_ENUMERATE         uint32 = 0x00000006
	CONN_FRC_SENT_REQUEST_CAM_DETAIL uint32 = 0x00000007
	CONN_FRC_SSL_HANDSHAKE_FAILED    uint32 = 0x00000008
	CONN_FRC_SSL_CONNECTED           uint32 = 0x00000009
)

const (
	Frc_Op_Video_Data       uint16 = 1
	Frc_Op_PlaybackControl  uint16 = 2
	Frc_Op_MetaDataLiveData uint16 = 3
	Frc_Op_SystemNotify     uint16 = 4
	Frc_Op_PlayerState      uint16 = 5
	Frc_Op_Audio_Data       uint16 = 6
	// in general request and reply in pair : opcode + 1 for reply
	// reserve some range for future one way communication. (1..99)
	Frc_Op_ReqPlayback   uint16 = 100 // 0x64
	Frc_Op_PlaybackReply uint16 = 101 // 0x65
	Frc_Op_ReqPing       uint16 = 102 // 0x66
	Frc_Op_PingReply     uint16 = 103 // 0x67
	Frc_Op_Subscribe     uint16 = 104 // 0x68
	//Frc_Op_SubscribeRepy uint16 =  105
	Frc_Op_UnSubscribe uint16 = 106
	//Frc_Op_UnSubscribeReply uint16 =  107
	Frc_Op_ReqCamEnumerate    uint16 = 108 // 0x6C
	Frc_Op_CamEnumerateResult uint16 = 109 // 0x6D
	//Frc_Op_ReqCamEventList uint16 =  110
	//Frc_Op_CamEventListResult uint16 =  111  //We no longer need the oldest event request query (also it is not supported in the notification branch)
	//Authorization
	Frc_Op_ReqConnect            uint16 = 112 // 0x70
	Frc_Op_ConnectReply          uint16 = 113 // 0x71
	Frc_Op_ReqUserAuthenticate   uint16 = 114 // 0x72
	Frc_Op_UserAuthenticateReply uint16 = 115 // 0x73
	// PTZ
	Frc_Op_ReqCamControl   uint16 = 116
	Frc_Op_CamControlReply uint16 = 117

	// 118
	// 119

	Frc_Op_AuxiliaryDataRequest uint16 = 120 // 0x78
	Frc_Op_AuxiliaryDataReply   uint16 = 121 // 0x79

	// 122
	// 123

	Frc_Op_MetaDataIntervalRequest uint16 = 124
	Frc_Op_MetaDataIntervalReply   uint16 = 125
	Frc_Op_MetaDataQueryRequest    uint16 = 126
	Frc_Op_MetaDataQueryReply      uint16 = 127

	// 128
	// 129

	Frc_Op_DeviceDetailsRequest uint16 = 130 // 0x82
	Frc_Op_DeviceDetailsReply   uint16 = 131 // 0x83

	// 132
	// 133

	// Peer 2 Peer
	Frc_Op_ReqClientID   uint16 = 134
	Frc_Op_ClientIDReply uint16 = 135

	Frc_Op_P2PStatusUpdate uint16 = 136
	// 137

	Frc_Op_ReqP2PNameEtc   uint16 = 138
	Frc_Op_P2PNameEtcReply uint16 = 139

	Frc_Op_ReqP2PViewSubsc   uint16 = 140
	Frc_Op_P2PViewSubscReply uint16 = 141

	Frc_Op_ReqP2PPlaceItemInPane uint16 = 142
	//Frc_Op_P2PPlaceItemInPaneReply uint16 =  143

	Frc_Op_ReqP2PPlaceViewOrLayout uint16 = 144
	//Frc_Op_P2PPlaceViewOrLayoutReply uint16 =  145

	Frc_Op_ReqP2PDoubleClickPane uint16 = 146
	//Frc_Op_P2PDoubleClickPaneReply uint16 =  147

	Frc_Op_ReqWindowControl uint16 = 148
	//Frc_Op_WindowControlReply uint16 =  149

	Frc_Op_P2PChat uint16 = 150
	//Frc_Op_P2PChatReply uint16 =  151

	Frc_Op_ReqDeviceEvents        uint16 = 152 // 0x98
	Frc_Op_DeviceEventsResult     uint16 = 153 // 0x99
	Frc_Op_ReqCreateAnnotation    uint16 = 154 // 0x9a
	Frc_Op_CreateAnnotationResult uint16 = 155
	Frc_Op_ReqPullEdge            uint16 = 156
	Frc_Op_PullEdgeResult         uint16 = 157

	Frc_Op_ReqLockEvents    uint16 = 158 // 0x9e
	Frc_Op_LockEventsResult uint16 = 159

	Frc_Op_GetPtzValuesRequest uint16 = 160
	Frc_Op_GetPtzValuesReply   uint16 = 161
	Frc_Op_SetPtzPresetRequest uint16 = 162
	Frc_Op_SetPtzPresetReply   uint16 = 163
	Frc_Op_DelPtzPresetRequest uint16 = 164
	Frc_Op_DelPtzPresetReply   uint16 = 165

	Frc_Op_BulkDeviceDetailRequest uint16 = 166
	Frc_Op_BulkDeviceDetailReply   uint16 = 167

	// This section as well as the one from 142 till 150 is for P2P
	Frc_Op_ReqP2PScheduleDownload uint16 = 166
	//Frc_Op_P2PScheduleDownloadReply uint16 =  167
	Frc_Op_ReqP2PGetItem   uint16 = 168
	Frc_Op_P2PGetItemReply uint16 = 169
	Frc_Op_ReqP2PGetFile   uint16 = 170
	Frc_Op_P2PGetFileReply uint16 = 171

	//Frc_Op_AuxiliaryDataRequest uint16 =  120     NewRequestWholeVideoFile uint16 =  0x00000006

	//fixme msg 172/173 are not yet supported
	Frc_Op_Req3rdEvents    uint16 = 172
	Frc_Op_3rdEventsResult uint16 = 173

	// IMPORTANT!!! Plase do not use numbers from here till about 250 or 300 (reserved for P2P)
	Frc_Op_ReqAddPersonToTrack uint16 = 300
	Frc_Op_ReqAddCVMarker      uint16 = 302
	Frc_Op_ReqGetPeople        uint16 = 303
	Frc_Op_GetPeopleReply      uint16 = 304

	Frc_Op_ReqGetCVTracks   uint16 = 305
	Frc_Op_GetCVTracksReply uint16 = 306

	Frc_Op_ReqDeletePerson   uint16 = 307
	Frc_Op_DeletePersonReply uint16 = 308

	Frc_Op_ReqCreatePerson   uint16 = 309
	Frc_Op_CreatePersonReply uint16 = 310

	Frc_Op_ReqAttachMugShotToPerson   uint16 = 311
	Frc_Op_AttachMugShotToPersonReply uint16 = 312

	Frc_Op_ReqDeleteTrack   uint16 = 313
	Frc_Op_DeleteTrackReply uint16 = 314

	Frc_Op_ReqRenamePerson   uint16 = 315
	Frc_Op_RenamePersonReply uint16 = 316

	Frc_Op_ReqWhoIsThis   uint16 = 317
	Frc_Op_WhoIsThisReply uint16 = 318

	Frc_Op_ReqGetDBUniqueCameras   uint16 = 319
	Frc_Op_GetDBUniqueCamerasReply uint16 = 320

	//Services
	Frc_Op_ProvideService      uint16 = 340
	Frc_Op_ProviceServiceReply uint16 = 341
	Frc_Op_BackServiceRequest  uint16 = 342
	Frc_Op_BackServiceReply    uint16 = 343

	RESERVED_FOR_SSL_REQUEST uint16 = 998 // used in FrcSocket to init SSL
	RESERVED_FOR_SSL_REPLY   uint16 = 999 //

	// From this point down the operations are local (internal to the app and not meant for Unix side
	Local             uint16 = 1000 // use this to identify the range.
	Lcl_Op_Register   uint16 = 1001
	Lcl_Op_UnRegister uint16 = 1002
)

// enum MsgAuxiliaryInfo
const (
	//Aux Data
	RequestImage             = 0x00000001 // parameter1 : unix Timestamp, parameter2 unix Timestampoffset
	RequestVideoFile         = 0x00000002 // parameter1 (From Time unix) parameter2 (To Time unix)
	RequestEdgeStorageFile   = 0x00000003
	OldRequestWholeVideoFile = 0x00000005 // parameter1 (Time in video) parameter2 (Event State & Event Type)
	NewRequestWholeVideoFile = 0x00000006 // data (i64) Time in video (i32) Event State (u64) Event Type

	AttachChunkData  = 0x00000004 // parameter1 : file size, parameter2 offset, extra data and length : chunk data and size.
	ReqNextChunkData = 0x00000008 // parameter1 : file size, parameter2 : next offset
	AbortChunkData   = 0x00000010

	//Event Lists
	NegativeEvents = 0x00000001 //These are events that should be removed

	// CV QueryManager (and may be other messages)
	ResultSetTruncated = 0x10

	//Error states
	ErrorMark       = 0x80000000
	ACK             = 0x00000001
	NACK            = 0x80000002 // same as MsgAuthFlag.NACK
	VersionMismatch = 0xFFFFFFFD
	ServerBusy      = 0xFFFFFFFE
	RequestError    = 0xFFFFFFFF
)

const (
	Camera_Recorded     = 0x000001
	Camera_Scheduled    = 0x000002
	Camera_Reset        = 0x000004
	Camera_Reboot       = 0x000008
	Camera_Power_Up     = 0x000010
	Camera_Restart      = 0x000020
	Camera_Disable      = 0x000040
	Camera_Enable       = 0x000080
	Camera_Annotate     = 0x000100
	Camera_Notification = 0x000200
	Camera_Interruption = 0x000400
	Camera_Quota        = 0x000800
	Camera_Upgrade      = 0x001000
	Camera_Msg_Event    = 0x002000
	Camera_Suspend      = 0x004000 /* Suspend recording */
	Camera_Resume       = 0x008000 /* Resume recording */
	Camera_SD_Format    = 0x010000 /* SD card format */
	Camera_All          = 0x00ffffff

	Camera_Motion     = 0x01000000
	Camera_Continuous = 0x02000000
	Camera_Temp       = 0x04000000
	Camera_Manual     = 0x08000000
	Camera_Audio      = 0x10000000
	Camera_DI         = 0x20000000
	Camera_PIR        = 0x40000000

	Type_Common_Request = Camera_Recorded | Camera_Motion | Camera_Audio | Camera_DI | Camera_PIR | Camera_Continuous | Camera_Annotate | Camera_Msg_Event
)

const (
	INVALID_TYPE      uint64 = 0x0000000000000000
	RECORDER          uint64 = 0x0001000000000000
	CAMERA            uint64 = 0x0002000000000000
	QUERY             uint64 = 0x0003000000000000
	VIDEO_SEQUENCER   uint64 = 0x0004000000000000
	FRCC_INTERNAL     uint64 = 0x0005000000000000
	MEDIA             uint64 = 0x0006000000000000
	CLIENT            uint64 = 0x0007000000000000 // used by RemoteClient
	MAP3D             uint64 = 0x0008000000000000
	CONVERSATION      uint64 = 0x0009000000000000 // Used in ItemConversation. Can have SubType GROUP or PRIVATE. If Subtype is NONE it indicates an uninitialized conversation (nobody said anything yet).
	WINDOW            uint64 = 0x000A000000000000 // used by RemoteWindow as child of CLIENT (RemoteClient)
	UNUSED1           uint64 = 0x000B000000000000 // NO LONGER Used in Chat.ChatMessage so feel free to use this number.
	CHATGROUP         uint64 = 0x000C000000000000 // used in ItemChatGroup
	FORTIVIEW         uint64 = 0x000D000000000000
	FORTICONNECTION   uint64 = 0x000E000000000000
	ANALYTICS_PROFILE uint64 = 0x000F000000000000
	// NOTE: If you are adding a new type and DeviceId with this type will contain Recorder ID in the parent part, you need to add it in ApplicationState.cs in GetRecorderChildTypes
)

const (
	PTZcapable uint32 = 0x00000002
	Disabled   uint32 = 0x00000001
)

const (
	GLOBAL   uint64 = 0x8000000000000000
	LIVE     uint64 = 0x4000000000000000
	FOREIGN  uint64 = 0x2000000000000000 // unused so far, rename to what you need
	FLAG4    uint64 = 0x1000000000000000 // unused so far, rename to what you need
	FLAG5    uint64 = 0x0800000000000000 // unused so far, rename to what you need
	FLAG6    uint64 = 0x0400000000000000 // unused so far, rename to what you need
	FLAG7    uint64 = 0x0200000000000000 // unused so far, rename to what you need
	INTERNAL uint64 = 0x0100000000000000 // flag indicating that this deviceID represents something internal to the application (this shall or shall not refer to the same internal device in different instances of FRC based on the Global Flag)
)

const (
	NONE    uint64 = 0x0000000000000000
	VIDEO   uint64 = 0x0000010000000000
	AUDIO   uint64 = 0x0000020000000000
	META    uint64 = 0x0000030000000000
	GROUP   uint64 = 0x0000040000000000 // used with Type CONVERSATION for conversations with >3 participants. For these parent part is the own part of CLIENTID of the initiator.
	PRIVATE uint64 = 0x0000050000000000 // used with Type CONVERSATION for conversations with 2 participants. For those, parent part is empty and own part is XOR of 2 participants own parts of CLIENTID.
)

const (
	MASK_FLAGS    uint64 = 0xFF00000000000000
	MASK_TYPE     uint64 = 0x00FF000000000000
	MASK_SUBTYPE  uint64 = 0x0000FF0000000000
	MASK_SUBID    uint64 = 0x000000FF00000000
	MASK_PARENTID uint64 = 0x00000000FFFF0000
	MASK_OWNID    uint64 = 0x000000000000FFFF
	MASK_RECORDER uint64 = MASK_SUBID | MASK_OWNID
)

const (
	CAMERA_MODEL_NAME        int = 0
	CAMERA_TAGS              int = 1
	CAMERA_PTZPRESETS        int = 2
	CAMERA_VIDEO_STREAM_LIST int = 3
	CAMERA_VIDEO_SETTINGS    int = 4
	CAMERA_NAME              int = 5
)

const (
	LOCK_EVENT   uint32 = 1
	UNLOCK_EVENT uint32 = 0
)

func MakeID(flags uint64, Type uint64, subType uint64, subId uint8, parentId uint16, ownId uint16) uint64 {
	return (flags | Type | subType | ((uint64)(subId))<<32 | (uint64)(parentId)<<16 | (uint64)(ownId))
}

type Version struct {
	Major uint16
	Minor uint16
	Build uint16
}

func PrintVersion(ver Version) string {
	return strconv.Itoa(int(ver.Major)) + "." + strconv.Itoa(int(ver.Minor)) + "." + strconv.Itoa(int(ver.Build))
}

const (
	MsgHeaderSize uint32 = 64
	BufferLength  uint32 = 1024 * 1024 * 10
	MaxKeepAlive  uint32 = 10
	MaxStart      uint32 = 3
)

type MsgStatus struct {
	Packet          FrcPacket
	ExtraDataLength int
	HeaderSize      int
}

type FrcPacket struct {
	Header FrcPacketHeader
}

type FrcPacketHeader struct {
	Version         uint16
	OpCode          uint16
	ExtraDataLength uint32
	Flag            uint32
	Index           uint32
	DeviceId        uint64
	Src             uint64
	Dest            uint64
	TimeStamp       int64
	Param1          int64
	Param2          int64
}

type FrcMsg struct {
	Version     uint16
	OpCode      uint16
	Data        []byte
	Flag        uint32
	Index       uint32
	DeviceId    uint64
	Source      uint64
	Destination uint64
	TimeStamp   int64
	Param1      int64
	Param2      int64
	Length      uint32
}
type DbController_Cmd struct {
	Action string `json:"action"` //poolstatus, getID, setLimit, getLimit, getLow, getNum, cleanup, SetByID, DelByID, DelSpecificByID
	Type   string `json:"type"`
	Data   string `json:"data"`
	Key    string `json:"key"`
	Cnt    int    `json:"cnt"` //limit or ID
}
type DhcpController_Cmd struct {
	Ifname  string `json:"ifname"`  //interface name
	Ver     int    `json:"version"` //version: 4 or 6
	Retries int    `json:"retries"` //number of retries
	Debug   bool   `json:"debug"`   //verbose control
	Dryrun  bool   `json:"dryrun"`  //run without setting configuration
	NoIfup  bool   `json:"noifup"`  //no need to bring the interface up
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
type Controller_Cmd struct {
	Action string `json:"action"`
	FrcId  uint64 `json:"frcid,string"`
	Type   int    `json:"type"`
}
type Subscription struct {
	CameraId uint64 `json:"cameraid,string"`
	FrcId    uint64 `json:"frcid,string"`
	VideoId  uint64 `json:"videoid,string"`
	UuId     uint32 `json:"uuid"`
}

// for video download index
var video_download_counter uint32 = 1
var video_download_mutex *sync.Mutex = &sync.Mutex{}

func GetVideoDownloadIndex() uint32 {
	var index uint32
	video_download_mutex.Lock()
	video_download_counter += 1
	index = video_download_counter
	video_download_mutex.Unlock()
	return index
}

const (
	AI_VIDEODOWNLOAD  = 0x00000001
	GUI_VIDEODOWNLOAD = 0x00000002
)

type VideoDownloadCmd struct {
	FrcId    uint64 `json:"frcid,string"`
	DeviceId uint64 `json:"cameraid,string"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Index    uint32 `json:"index"` // let GUI to check the video download status
	Flag     uint32 `json:"flag"`
	FileName string `json:"filename"`
	Type     int    `json:"type"` // type 1 for AI, 2 for GUI
}

const (
	Event_State_Active         = 0x00000001
	Event_State_Inactive       = 0x00000002
	Event_State_Empty_Local    = 0x00000004 //this means that there is an edge recording (posible more or less)
	Event_State_NonEmpty_Local = 0x00000008
	Event_State_Locked         = 0x00000010
	Event_State_UnLocked       = 0x00000020

	Event_State_Common_Request = Event_State_Active | Event_State_Inactive
)

type VideoEvent struct {
	Type      uint32 `json:"type"`
	State     uint32 `json:"state"`
	Start     uint64 `json:"start"`
	End       uint64 `json:"end"`
	ExtraData []byte `json:"data"`
}

type TimelineReq struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}

type DeviceDetail struct {
	typeId uint16
	data   []byte
}

type DeviceDetails []DeviceDetail

type Stream struct {
	Width     int32  `json:"width"`
	Height    int32  `json:"height"`
	FrameRate int32  `json:"framerate"`
	VideoId   uint64 `json:"videoid,string"`
}

type Download struct {
	ThreadId uint64
	FileSize int64
	Offset   int64
	Index    uint32
	Status   uint32
	FileDir  string
	FileName string
	Type     int
}

type DownloadFile struct {
	FileDir  string `json:"filedir"`
	FileName string `json:"filename"`
	DeviceId string `json:"cameraid"`
	Start    string `json:"starttime"`
}

type EventLock struct {
	FrcId    uint64 `json:"frcid,string"`
	DeviceId uint64 `json:"cameraid,string"`
	Action   uint32 `json:"action"` // 1:lock or 0:unlock
	Start    uint64 `json:"start"`
	End      uint64 `json:"end"`
}

type Annotation struct {
	FrcId    uint64 `json:"frcid,string"`
	DeviceId uint64 `json:"cameraid,string"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Action   uint32 `json:"action"` // 0: set, 1: clear
	Name     string `json:"name"`
}

const (
	LIVE_STREAM    = 0x00000001
	HISTORY_STREAM = 0x00000002
	DATA_EXCHANGE  = 0x00000004
)

var LIST_OF_FRC_CLIENT_TYPES = []int{LIVE_STREAM, HISTORY_STREAM, DATA_EXCHANGE}

type Client struct {
	status               int      // 0: working, 1: not working
	frc                  Recorder // recorder config from DB
	ver                  Version
	ID                   int // recorder client id
	Type                 int // Live or History stream
	IpAddr               string
	Port                 int
	conn                 net.Conn
	connStatus           uint32
	preConnStatus        uint32
	keepAliveCnt         uint32
	startCnt             uint32
	OpCode               uint16
	w                    *bufio.Writer
	FrcId                uint64 // fortirecorder id replied from recorder after connection establishes
	uuid_mutex           *sync.RWMutex
	st_mutex             *sync.RWMutex
	Statistics           map[uint64]uint64 // traffic statistics
	readyCamCnt          uint32
	camCnt               uint32
	buf                  [BufferLength]byte
	bufPos               uint32
	dataPos              uint32
	timeResetCnt         uint32 // count time rest
	videoPkgCnt          uint32
	cmd                  chan Cmd
	rcvHack              chan string
	keepAliveHack        chan string
	download             map[uint32]*Download
	download_mutex       *sync.RWMutex
	frcReqHeaderChanMap  map[string](chan *FrcReqHeader)
	frcReqHeaderHookMap  map[string][]*FrcReqHeader
	frcReqHeaderMutexMap map[string]*sync.Mutex
}

// TBD, define actions.
const (
	CmdActionClose   = 0
	CmdActionRestart = 1
)

type Cmd struct {
	act           int
	recorderId    int
	Type          int
	preConnStatus uint32
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

const (
	ClientManagerChanBufferSize   = 1024
	ClientRecoderTimeDiffThrottle = 5 // in Second
)

type ClientManager struct {
	clients      map[int]*Client // fortirecorder client id -- client pointer
	unicast      chan Cmd
	register     chan *Client
	unregister   chan *Client
	uuid_release chan uint32
}

type CameraDescriptor struct {
	CameraId     uint32
	CameraName   [128]byte
	CameraStatus uint32
}

type CamDescs []CameraDescriptor

// FortiRecorder message response header
type FrcReqHeader struct {
	JWTToken  string
	WsId      uint32
	AITaskId  int
	FrcId     uint64
	DeviceId  uint64
	ReqType   string
	TimeStamp uint64
}
