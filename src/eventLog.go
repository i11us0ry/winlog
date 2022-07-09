package src

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

// OpenRemoteEventLog does the same as Open, but on different computer host.
//func (el *EventLog)OpenRemoteEventLog(host, source string) {
//	if source == "" {
//		fmt.Println("Specify event log source")
//	}
//	var s *uint16
//	if host != "" {
//		s = syscall.StringToUTF16Ptr(host)
//	}
//	h, err := openEventLog(s, syscall.StringToUTF16Ptr(source))
//	if err != nil {
//		fmt.Println(err)
//	}
//	buf := make([]byte, MAX_BUFFER_SIZE+1)
//	el.Handle = h
//	el.bufferSize = uint32(MAX_DEFAULT_BUFFER_SIZE)
//	el.buffer = buf
//	el.minRead = uint32(0)
//}

func openEventLog(uncServerName *uint16, sourceName *uint16) (handle syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procOpenEventLog.Addr(), 2, uintptr(unsafe.Pointer(uncServerName)), uintptr(unsafe.Pointer(sourceName)), 0)
	handle = syscall.Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getNumberOfEventLogRecords(eventLog syscall.Handle, numberOfRecords *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetNumberOfEventLogRecords.Addr(), 2, uintptr(eventLog), uintptr(unsafe.Pointer(numberOfRecords)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func closeEventLog(eventLog syscall.Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procCloseEventLog.Addr(), 1, uintptr(eventLog), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func readEventLog(eventLog syscall.Handle, readFlags uint32, recordOffset uint32, buffer *byte, numberOfBytesToRead uint32, bytesRead *uint32, minNumberOfBytesNeeded *uint32) (err error) {
	r1, _, e1 := syscall.Syscall9(procReadEventLog.Addr(), 7, uintptr(eventLog), uintptr(readFlags), uintptr(recordOffset), uintptr(unsafe.Pointer(buffer)), uintptr(numberOfBytesToRead), uintptr(unsafe.Pointer(bytesRead)), uintptr(unsafe.Pointer(minNumberOfBytesNeeded)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getOldestEventLogRecord(eventLog syscall.Handle, oldestRecord *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetOldestEventLogRecord.Addr(), 2, uintptr(eventLog), uintptr(unsafe.Pointer(oldestRecord)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func loadLibraryEx(filename *uint16, file syscall.Handle, flags uint32) (handle syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procLoadLibraryExW.Addr(), 3, uintptr(unsafe.Pointer(filename)), uintptr(file), uintptr(flags))
	handle = syscall.Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func formatMessage(flags uint32, source syscall.Handle, messageID uint32, languageID uint32, buffer *byte, bufferSize uint32, arguments uintptr) (numChars uint32, err error) {
	r0, _, e1 := syscall.Syscall9(procFormatMessageW.Addr(), 7, uintptr(flags), uintptr(source), uintptr(messageID), uintptr(languageID), uintptr(unsafe.Pointer(buffer)), uintptr(bufferSize), uintptr(arguments), 0, 0)
	numChars = uint32(r0)
	if numChars == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// 搜索日志
func (el *EventLog) Run(logName string) (err error) {
	// 获取statefile
	recordNumber := uint32(0)

	ptr := syscall.StringToUTF16Ptr(logName)
	// 打开日志
	h, err := openEventLog(nil,ptr)
	if err != nil {
		fmt.Println("系统权限可能不够!Administrator privileges are required!")
		return
	}
	// 关闭日志
	defer closeEventLog(h)

	var num, oldnum uint32

	// 获取日志数量
	getNumberOfEventLogRecords(h, &num)
	if err != nil {
		fmt.Println("系统权限可能不够!Administrator privileges are required!")
		return
	}

	getOldestEventLogRecord(h, &oldnum)
	if err != nil {
		fmt.Println("系统权限可能不够!Administrator privileges are required!")
		return
	}

	// 比较日志数量
	if oldnum <= recordNumber {
		if recordNumber == oldnum+num-1 {
			return err
		}
		recordNumber++
	} else {
		recordNumber = oldnum
	}

	size := uint32(1)
	buf := []byte{0}

	var readBytes uint32
	var nextSize uint32

	// 循环读取日志
loop_events:
	for i := recordNumber; i < oldnum+num; i++ {
		flags := EVENTLOG_FORWARDS_READ | EVENTLOG_SEEK_READ
		if i == 0 {
			flags = EVENTLOG_FORWARDS_READ | EVENTLOG_SEQUENTIAL_READ
		}
		err = readEventLog(
			h,
			uint32(flags),
			i,
			&buf[0],
			size,
			&readBytes,
			&nextSize)
		if err != nil {
			if err != syscall.ERROR_INSUFFICIENT_BUFFER {
				if err != errorInvalidParameter {
					return err
				}
				break
			}
			buf = make([]byte, nextSize)
			size = nextSize
			err = readEventLog(
				h,
				uint32(flags),
				i,
				&buf[0],
				size,
				&readBytes,
				&nextSize)
			if err != nil {
				log.Printf("eventlog.ReadEventLog: %v", err)
				break
			}
		}
		r := *(*EVENTLOGRECORD)(unsafe.Pointer(&buf[0]))
		// 4624
		if r.EventID == 4624 {
			si := SuccessInfo{}
			si.time = time.Unix(int64(r.TimeGenerated), 0).String()
			// even code takes last 4 byte
			eventID := r.EventID & 0x0000FFFF
			if len(el.eventID) > 0 {
				accepted := false
				for _, idr := range el.eventID {
					if idr.lo <= eventID && eventID <= idr.hi {
						accepted = true
						break
					}
				}
				if !accepted {
					continue loop_events
				}
			}

			// 获取字段
			sourceName, _ := bytesToString(buf[unsafe.Sizeof(EVENTLOGRECORD{}):])
			off := uint32(0)
			args := make([]*byte, uintptr(r.NumStrings)*unsafe.Sizeof((*uint16)(nil)))
			for n := 0; n < int(r.NumStrings); n++ {
				args[n] = &buf[r.StringOffset+off]
				_, boff := bytesToString(buf[r.StringOffset+off:])
				off += boff + 2
			}
			var argsptr uintptr
			if r.NumStrings > 0 {
				argsptr = uintptr(unsafe.Pointer(&args[0]))
			}
			// 信息, 帐户登录失败。 已成功登录帐户。
			message, _ := getResourceMessage(logName, sourceName, r.EventID, argsptr)
			//fmt.Println(message)
			//log.Printf("Message=%v", message)
			if srcIp := getSrcIP(message);srcIp!=""{
				//fmt.Println("sip:",srcIp)
				si.sip = srcIp
				if srcPort := getSrcPort(message);srcPort!=""{
					//fmt.Println("srcPort:",srcPort)
					si.sport = srcPort
				}
				if ln := getSucLoginName(message);len(ln)!=0{
					//fmt.Println("loginName:",ln)
					si.lName = ln
				}
				if lp := getLoginPro(message); lp !=""{
					//fmt.Println("loginProccess:",lp)
					si.lPro = lp
				}
				el.success = append(el.success,si)
			}
		}
		if r.EventID == 4625{
			//fmt.Println("TimeGenerated:", time.Unix(int64(r.TimeGenerated), 0).String())
			fi := FailInfo{}
			fi.time = time.Unix(int64(r.TimeGenerated), 0).String()
			eventID := r.EventID & 0x0000FFFF
			if len(el.eventID) > 0 {
				accepted := false
				for _, idr := range el.eventID {
					if idr.lo <= eventID && eventID <= idr.hi {
						accepted = true
						break
					}
				}
				if !accepted {
					continue loop_events
				}
			}

			// 获取字段
			sourceName, _ := bytesToString(buf[unsafe.Sizeof(EVENTLOGRECORD{}):])
			off := uint32(0)
			args := make([]*byte, uintptr(r.NumStrings)*unsafe.Sizeof((*uint16)(nil)))
			for n := 0; n < int(r.NumStrings); n++ {
				args[n] = &buf[r.StringOffset+off]
				_, boff := bytesToString(buf[r.StringOffset+off:])
				off += boff + 2
			}

			var argsptr uintptr
			if r.NumStrings > 0 {
				argsptr = uintptr(unsafe.Pointer(&args[0]))
			}
			// 信息, 帐户登录失败。 已成功登录帐户。
			message, _ := getResourceMessage(logName, sourceName, r.EventID, argsptr)
			//fmt.Println(message)
			//log.Printf("Message=%v", message)
			if srcIp := getSrcIP(message);srcIp!=""{
				//fmt.Println("sip:",srcIp)
				fi.sip = srcIp
				if srcPort := getSrcPort(message);srcPort!=""{
					//fmt.Println("srcPort:",srcPort)
					fi.sport = srcPort
				}
				if failm := getFail(message);failm!=""{
					//fmt.Println("failMessage:",failm)
					fi.failinfo = failm
				}
				if ln := getFailLoginName(message);len(ln)!=0{
					//fmt.Println("loginName:",ln)
					fi.lName = ln
				}
				if lp := getLoginPro(message); lp !=""{
					//fmt.Println("loginProccess:",lp)
					fi.lPro = lp
				}
				el.fail = append(el.fail,fi)
			}
		}
	}
	return nil
}

