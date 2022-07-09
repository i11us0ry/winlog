package src

import (
	"regexp"
	"syscall"
)

type EventLogHandle 		= uintptr
type EventSourceHandle 		= uintptr
type WIN32_ERROR uint32
type SIDType uint32
type EventType uint16

var stateRe = regexp.MustCompile(`^([A-Z]):[/\\]`)
var Language = ""

const Log = "Security"
const errorInvalidParameter = syscall.Errno(87)

// ReadEventLog Flags
const (
	EVENTLOG_SEEK_READ       = 0x0002
	EVENTLOG_SEQUENTIAL_READ = 0x0001
	EVENTLOG_FORWARDS_READ   = 0x0004
	EVENTLOG_BACKWARDS_READ  = 0x0008
	MAX_BUFFER_SIZE                = 0x7ffff
	MAX_DEFAULT_BUFFER_SIZE        = 0x10000
)

const (
	DONT_RESOLVE_DLL_REFERENCES         uint32 = 0x0001
	LOAD_LIBRARY_AS_DATAFILE            uint32 = 0x0002
	LOAD_WITH_ALTERED_SEARCH_PATH       uint32 = 0x0008
	LOAD_IGNORE_CODE_AUTHZ_LEVEL        uint32 = 0x0010
	LOAD_LIBRARY_AS_IMAGE_RESOURCE      uint32 = 0x0020
	LOAD_LIBRARY_AS_DATAFILE_EXCLUSIVE  uint32 = 0x0040
	LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR    uint32 = 0x0100
	LOAD_LIBRARY_SEARCH_APPLICATION_DIR uint32 = 0x0200
	LOAD_LIBRARY_SEARCH_USER_DIRS       uint32 = 0x0400
	LOAD_LIBRARY_SEARCH_SYSTEM32        uint32 = 0x0800
	LOAD_LIBRARY_SEARCH_DEFAULT_DIRS    uint32 = 0x1000
)

// 日志字段
type EVENTLOGRECORD struct {
	Length              uint32
	Reserved            uint32
	RecordNumber        uint32
	TimeGenerated       uint32
	TimeWritten         uint32
	EventID             uint32
	EventType           uint16
	NumStrings          uint16
	EventCategory       uint16
	ReservedFlags       uint16
	ClosingRecordNumber uint32
	StringOffset        uint32
	UserSidLength       uint32
	UserSidOffset       uint32
	DataLength          uint32
	DataOffset          uint32
}

const (
	SidTypeUser SIDType = 1 + iota
	SidTypeGroup
	SidTypeDomain
	SidTypeAlias
	SidTypeWellKnownGroup
	SidTypeDeletedAccount
	SidTypeInvalid
	SidTypeUnknown
	SidTypeComputer
	SidTypeLabel
)

const (
	EVENTLOG_SUCCESS    EventType = 0
	EVENTLOG_ERROR_TYPE           = 1 << (iota - 1)
	EVENTLOG_WARNING_TYPE
	EVENTLOG_INFORMATION_TYPE
	EVENTLOG_AUDIT_SUCCESS
	EVENTLOG_AUDIT_FAILURE
)

func (et EventType) String() string {
	switch et {
	case EVENTLOG_SUCCESS:
		return "Success"
	case EVENTLOG_ERROR_TYPE:
		return "Error"
	case EVENTLOG_AUDIT_FAILURE:
		return "Audit Failure"
	case EVENTLOG_AUDIT_SUCCESS:
		return "Audit Success"
	case EVENTLOG_INFORMATION_TYPE:
		return "Information"
	case EVENTLOG_WARNING_TYPE:
		return "Warning"
	default:
		return "Unknown"
	}
}

func (st SIDType) String() string {
	switch st {
	case SidTypeUser:
		return "User"
	case SidTypeGroup:
		return "Group"
	case SidTypeDomain:
		return "Domain"
	case SidTypeAlias:
		return "Alias"
	case SidTypeWellKnownGroup:
		return "Well Known Group"
	case SidTypeDeletedAccount:
		return "Deleted Account"
	case SidTypeInvalid:
		return "Invalid"
	case SidTypeUnknown:
		return "Unknown"
	case SidTypeComputer:
		return "Unknown"
	case SidTypeLabel:
		return "Label"
	default:
		return "Unknown"
	}
}