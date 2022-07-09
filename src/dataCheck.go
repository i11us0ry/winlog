package src

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"regexp"
	"strings"
	"syscall"
	"unicode/utf16"
)

func bytesToString(b []byte) (string, uint32) {
	var i int
	s := make([]uint16, len(b)/2)
	for i = range s {
		s[i] = uint16(b[i*2]) + uint16(b[(i*2)+1])<<8
		if s[i] == 0 {
			s = s[0:i]
			break
		}
	}
	return string(utf16.Decode(s)), uint32(i * 2)
}

func getResourceMessage(providerName, sourceName string, eventID uint32, argsptr uintptr) (string, error) {
	regkey := fmt.Sprintf(
		"SYSTEM\\CurrentControlSet\\Services\\EventLog\\%s\\%s",
		providerName, sourceName)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regkey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	val, _, err := key.GetStringValue("EventMessageFile")
	if err != nil {
		return "", err
	}
	val, err = registry.ExpandString(val)
	if err != nil {
		return "", err
	}

	handle, err := loadLibraryEx(syscall.StringToUTF16Ptr(val), 0,
		DONT_RESOLVE_DLL_REFERENCES|LOAD_LIBRARY_AS_DATAFILE)
	if err != nil {
		return "", err
	}
	defer syscall.CloseHandle(handle)

	msgbuf := make([]byte, 1<<16)
	numChars, err := formatMessage(
		syscall.FORMAT_MESSAGE_FROM_SYSTEM|
			syscall.FORMAT_MESSAGE_FROM_HMODULE|
			syscall.FORMAT_MESSAGE_ARGUMENT_ARRAY,
		handle,
		eventID,
		0,
		&msgbuf[0],
		uint32(len(msgbuf)),
		argsptr)
	if err != nil {
		return "", err
	}
	message, _ := bytesToString(msgbuf[:numChars*2])
	message = strings.Replace(message, "\r", "", -1)
	message = strings.TrimSuffix(message, "\n")
	return message, nil
}

func getFail(message string) string{
	// 失败原因:	%%2313
	reg := regexp.MustCompile(`失败原因:[\t\n\f\r]+(.*)`)
	result := reg.FindAllStringSubmatch(message,-1)
	if len(result)>0{
		return result[0][1]
	} else {
		return ""
	}
}
func getSrcIP(message string) string{
	// 源网络地址:	192.168.43.251
	reg := regexp.MustCompile(`源网络地址:[\t\n\f\r]+(.*)`)
	result := reg.FindAllStringSubmatch(message,-1)
	if len(result)>0{
		sip := result[0][1]
		if strings.Contains(sip,"-"){
			return ""
		}
		return result[0][1]
	} else {
		return ""
	}
}

func getSrcPort(message string) string{
	// 源端口:		0
	reg := regexp.MustCompile(`源端口:[\t\n\f\r]+(.*)`)
	result := reg.FindAllStringSubmatch(message,-1)
	if len(result)>0{
		return result[0][1]
	} else {
		return ""
	}
}

// 失败 帐户名
func getFailLoginName(message string) string{
	// 帐户名:		2
	reg := regexp.MustCompile(`帐户名:[\t\n\f\r]+(.*)`)
	result := reg.FindAllStringSubmatch(message,-1)
	if len(result)>0{
		return result[1][1]
	} else {
		return ""
	}
}

// 成功 帐户名称
func getSucLoginName(message string) string{
	// 帐户名称:		2
	reg := regexp.MustCompile(`帐户名称:[\t\n\f\r]+(.*)`)
	result := reg.FindAllStringSubmatch(message,-1)
	if len(result)>0{
		return result[1][1]
	} else {
		return ""
	}
}

func getLoginPro(message string) string{
	// 登录进程:	NtLmSsp
	reg := regexp.MustCompile(`登录进程:[\t\n\f\r]+(.*)`)
	result := reg.FindAllStringSubmatch(message,-1)
	if len(result)>0{
		return result[0][1]
	} else {
		return ""
	}
}