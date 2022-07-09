package src

import (
	"golang.org/x/sys/windows/registry"
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