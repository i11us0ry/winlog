package src

import (
	"github.com/olekukonko/tablewriter"
	"golang.org/x/sys/windows"
	"os"
	"strconv"
	"syscall"
)

var (
	advapi = syscall.NewLazyDLL("Advapi32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procOpenEventLog                		= advapi.NewProc("OpenEventLogW")
	procReadEventLog 						= advapi.NewProc("ReadEventLogW")
	procCloseEventLog              			= advapi.NewProc("CloseEventLog")
	procGetNumberOfEventLogRecords  		= advapi.NewProc("GetNumberOfEventLogRecords")
	procGetOldestEventLogRecord    			= advapi.NewProc("GetOldestEventLogRecord")
	procLoadLibraryExW            			= kernel32.NewProc("LoadLibraryExW")
	procFormatMessageW             			= kernel32.NewProc("FormatMessageW")
)

type Infos interface {
	run(string)
	checkData(*rInfos)
	output()
}

// mstsc保存的远程连接的相关信息
type rInfo struct {
	host 	string
	port 	string
	user 	string
	pass 	string
}

type rInfos struct {
	rInfo		[]rInfo
}

type PassInfo struct {
	host 	string
	pass 	string
}

type PassInfos struct {
	passInfos 	[]PassInfo
}

type Info struct {
	port 	 				uint64
	mstscKeys 				[]string
	mstscValues 			[][]string
	mstscDefaultValueNames 	[]string
	mstscDefaultValues      []string
}

type idRange struct {
	hi uint32
	lo uint32
}

type SuccessInfo struct {
	time 		string
	sip 		string
	sport 		string
	lName 		string
	lPro 		string
}

type FailInfo struct {
	time 		string
	sip 		string
	sport 		string
	lName 		string
	lPro 		string
	failinfo 	string
}

type EventLog struct {
	eventID		   []idRange
	success 	   []SuccessInfo
	fail 		   []FailInfo
}

func Start(mimi string) {
	var ri = &rInfos{}
	// 获取mstsc信息
	var info = Info{}
	var infos Infos =  &info
	infos.run("")
	infos.checkData(ri)
	//infos.output()

	// 获取登录日志
	var eventLog = EventLog{}
	var event Infos = &eventLog
	eventLog.eventID = append(eventLog.eventID,idRange{4624,4624})
	eventLog.eventID = append(eventLog.eventID,idRange{4625,4625})
	event.run(Log)

	// 抓取密码
	var passInfo = PassInfos{}
	var getPwd Infos = &passInfo
	getPwd.run(mimi)
	getPwd.checkData(ri)
	//getPwd.output()

	if  info.port != 0{
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"local_Rdp_Port"})
		table.Append([]string{strconv.FormatInt(int64(info.port),10)})
		table.Render()
	}

	if len(ri.rInfo) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"r_Host","r_Port","user","pass"})
		for _,v := range(ri.rInfo){
			table.Append([]string{v.host,v.port,v.user,v.pass})
		}
		table.Render()
	}
	event.output()
}