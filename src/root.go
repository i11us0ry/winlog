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
	getRdp()
	getMstscKeys()
	getMstscValue(key string)
	getDefaultMstscValues()
	getDefaultMstsc(name string)
}

type Info struct {
	port 	 				uint64
	mstscKeys 				[]string
	mstscValues 			[][]string
	mstscDefaultValueNames 	[]string
	mstscDefaultValues      []string
}

type Event interface {
	Run(logName string) (err error)
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

func Start() {
	var info = Info{}
	var infos Infos =  &info
	infos.getRdp()
	infos.getMstscKeys()
	for _, keyName := range(info.mstscKeys){
		infos.getMstscValue(keyName)
	}
	infos.getDefaultMstscValues()
	for _, name := range(info.mstscDefaultValueNames){
		infos.getDefaultMstsc(name)
	}

	//fmt.Println(info.port,"\n",info.mstscValues,"\n",info.mstscDefaultValues,"\n")
	if  info.port != 0{
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"RdpPort"})
		table.Append([]string{strconv.FormatInt(int64(info.port),10)})
		table.Render()
	}

	if len(info.mstscValues) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Hosts","Name"})
		for _,hn := range(info.mstscValues){
			table.Append([]string{hn[0],hn[1]})
		}
		table.Render()
	}

	if len(info.mstscDefaultValues) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Hosts"})
		for _,h := range(info.mstscDefaultValues){
			table.Append([]string{h})
		}
		table.Render()
	}

	var eventLog = EventLog{}
	var event Event = &eventLog

	eventLog.eventID = append(eventLog.eventID,idRange{4624,4624})
	eventLog.eventID = append(eventLog.eventID,idRange{4625,4625})
	// 启动
	event.Run(Log)
	if len(eventLog.success) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"id","Time", "Source Ip", "Source Port","Account Name","Logon Process"})
		for _, s := range(eventLog.success){
			//fmt.Println(fmt.Sprintf("time:%v,sip:%v,sport:%v,loginName:%v,loginprocess:%v",s.time,s.sip,s.sport,s.lName,s.lPro))
			table.Append([]string{"4624",s.time,s.sip,s.sport,s.lName,s.lPro})
		}
		table.Render()
	}
	if len(eventLog.fail) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Id","Time", "Source Ip", "Source Port","Account Name","Logon Process","Failure Reason"})
		for _, s := range(eventLog.fail){
			//fmt.Println(fmt.Sprintf("time:%v,sip:%v,sport:%v,loginName:%v,loginprocess:%v,Reasons for failure%v",s.time,s.sip,s.sport,s.lName,s.lPro,s.failinfo))
			table.Append([]string{"4625",s.time,s.sip,s.sport,s.lName,s.lPro,s.failinfo})
		}
		table.Render()
	}
}