package src

import (
	"github.com/olekukonko/tablewriter"
	"golang.org/x/sys/windows/registry"
	"os"
	"strconv"
	"strings"
)

// SYSTEM\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp  本机rdp端口
// 计算机\HKEY_CURRENT_USER\Software\Microsoft\Terminal Server Client\Default mstsc缓存记录，包括端口
// 计算机\HKEY_CURRENT_USER\SOFTWARE\Microsoft\Terminal Server Client\Servers 当前用户登陆历史，包括用户名
// 计算机\HKEY_USERS\{SID}\SOFTWARE\Microsoft\Terminal Server Client\Servers 所有用户登录历史  暂不考虑

func (info *Info)getRdp() {
	// PortNumber REG_dword uint64
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp`, registry.ALL_ACCESS)
	if err != nil {
		//fmt.Println("getRdp:",err)
		return
	}
	defer k.Close()
	port, _, err := k.GetIntegerValue("PortNumber")
	if err != nil {
		//fmt.Println("getRdp:",err)
		return
	}
	info.port = port
}

func (info *Info)getMstscKeys() {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Terminal Server Client\Servers`, registry.ALL_ACCESS)
	if err != nil {
		//fmt.Println("getMstscKeys:",err)
		return
	}
	defer k.Close()
	hosts, err := k.ReadSubKeyNames(0)
	if err != nil {
		//fmt.Println("getMstscKeys:",err)
		return
	}
	info.mstscKeys = hosts
}

func (info *Info)getMstscValue(key string) {
	// UsernameHint REG_SZ string
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Terminal Server Client\Servers\`+key, registry.ALL_ACCESS)
	if err != nil {
		//fmt.Println("getMstscValue:",err)
		return
	}
	defer k.Close()
	user, _, err := k.GetStringValue("UsernameHint")
	if err != nil {
		//fmt.Println("getMstscValue:",err)
		return
	}
	info.mstscValues = append(info.mstscValues, []string{key,user})
}

func (info *Info)getDefaultMstscValues(){
	// UsernameHint REG_SZ string
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Terminal Server Client\Default`, registry.ALL_ACCESS)
	if err != nil {
		//fmt.Println("getDefaultMstscValues:",err)
		return
	}
	defer k.Close()
	ValueNames, err := k.ReadValueNames(0)
	if err != nil {
		//fmt.Println("getDefaultMstscValues:",err)
		return
	}
	info.mstscDefaultValueNames = ValueNames
}

func (info *Info)getDefaultMstsc(name string){
	// UsernameHint REG_SZ string
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Terminal Server Client\Default`, registry.ALL_ACCESS)
	if err != nil {
		//fmt.Println("getDefaultMstsc:",err)
		return
	}
	defer k.Close()
	Value, _, err := k.GetStringValue(name)
	if err != nil {
		//fmt.Println("getDefaultMstsc:",err)
		return
	}
	info.mstscDefaultValues = append(info.mstscDefaultValues,Value)
}

// 数据归纳
func (info *Info)checkData(info2 *rInfos){
	// 先归纳host和port
	if len(info.mstscDefaultValues) > 0{
		for _,hosts := range(info.mstscDefaultValues){
			if strings.Contains(hosts,":") {
				host_port := strings.Split(hosts,":")
				info2.rInfo = append(info2.rInfo,rInfo{host_port[0], host_port[1],"",""})
			} else {
				info2.rInfo = append(info2.rInfo,rInfo{hosts, "3389","",""})
			}
		}
	}
	// 归纳host和username
	if len(info.mstscValues) > 0 {
		for _, host := range(info.mstscValues) {
			index := findHost(host[0],info2.rInfo)
			if index == -1 {
				info2.rInfo = append(info2.rInfo,rInfo{host[0], "",host[1],""})
			} else {
				info2.rInfo[index].user = host[1]
			}
		}
	}
}

func findHost(host string, rinfo []rInfo) int{
	for k, v := range(rinfo){
		if v.host == host{
			return k
		}
	}
	return -1
}

func (info *Info)output(){
	//fmt.Println(info.port,"\n",info.mstscValues,"\n",info.mstscDefaultValues,"\n")
	if  info.port != 0{
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"local_Rdp_Port"})
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
}

func (info *Info)run(path string){
	info.getRdp()
	info.getMstscKeys()
	for _, keyName := range(info.mstscKeys){
		info.getMstscValue(keyName)
	}
	info.getDefaultMstscValues()
	for _, name := range(info.mstscDefaultValueNames){
		info.getDefaultMstsc(name)
	}
}