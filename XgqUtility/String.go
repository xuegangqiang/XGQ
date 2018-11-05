package XgqUtility

import (
	"fmt"
	"strings"
	"time"
)

func (u *Utility) IP2Str(IP []byte, Port uint32) string {
	var StrIP string

	StrIP = fmt.Sprintf("%d.%d.%d.%d:%d",
		IP[0], IP[1], IP[2], IP[3], Port)

	return StrIP
}

func (u *Utility) IP2StrWithoutPort(IP []byte) string {
	var StrIP string

	StrIP = fmt.Sprintf("%d.%d.%d.%d",
		IP[0], IP[1], IP[2], IP[3])

	return StrIP
}

func (u *Utility) Str2IP(StrIP string) ([]byte, uint32) {
	var IP []byte
	var Port uint32

	IP = make([]byte, 4, 4)
	if strings.Contains(StrIP, ":") {
		fmt.Sscanf(StrIP, "%d.%d.%d.%d:%d", &IP[0], &IP[1], &IP[2], &IP[3], &Port)
	} else {
		fmt.Sscanf(StrIP, "%d.%d.%d.%d", &IP[0], &IP[1], &IP[2], &IP[3])
		Port = 0
	}

	return IP, Port
}

func (u *Utility) GetCurTimeStr() string {
	var CurTimeStr string

	CurTime := time.Now()
	CurTimeStr = fmt.Sprintf(
		"%04d-%02d-%02d %02d:%02d:%02d.%03d",
		CurTime.Year(), CurTime.Month(), CurTime.Day(), CurTime.Hour(),
		CurTime.Minute(), CurTime.Second(), CurTime.Nanosecond()/1e6)

	return CurTimeStr
}
