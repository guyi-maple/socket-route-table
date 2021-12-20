package route

import (
	"strconv"
	"strings"
)

func isBelong(ip string, cidr string) bool {
	ipAddr := strings.Split(ip, `.`)
	if len(ipAddr) < 4 {
		return false
	}
	cidrArr := strings.Split(cidr, `/`)
	if len(cidrArr) < 2 {
		return false
	}
	var tmp = make([]string, 0)
	for key, value := range strings.Split(`255.255.255.0`, `.`) {
		iint, _ := strconv.Atoi(value)
		iint2, _ := strconv.Atoi(ipAddr[key])
		tmp = append(tmp, strconv.Itoa(iint&iint2))
	}
	return strings.Join(tmp, `.`) == cidrArr[0]
}
