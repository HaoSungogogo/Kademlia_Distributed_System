package libkademlia

import (
  "net"
  "strconv"
  "net/rpc"
  "log"
)


func StringToIpPort(laddr string) (ip net.IP, port uint16, err error) {
	hostString, portString, err := net.SplitHostPort(laddr)
	if err != nil {
		return
	}
	ipStr, err := net.LookupHost(hostString)
	if err != nil {
		return
	}
	for i := 0; i < len(ipStr); i++ {
		ip = net.ParseIP(ipStr[i])
		if ip.To4() != nil {
			break
		}
	}
	portInt, err := strconv.Atoi(portString)
	port = uint16(portInt)
	return
}

func IpPortToString(ip net.IP, port uint16) (ipStr string, portStr string) {
  ipStr = ip.String()
  portStr = strconv.Itoa(int(port))
  return
}

func getClient(host net.IP, port uint16) (client *rpc.Client) {
  hostStr, portStr := IpPortToString(host, port)

	client, err := rpc.DialHTTPPath("tcp", net.JoinHostPort(hostStr, portStr),
		rpc.DefaultRPCPath+hostStr+portStr)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	}
  return
}

func removeDup(l []Contact) (res []Contact) {
  res = append(res, l[0])

  for i := 1; i < len(l); i++ {
    if !l[i].NodeID.Equals(l[i-1].NodeID) {
      res = append(res, l[i])
    }
  }
  return
}

func minus(le, ri []Contact) (res []Contact) {
  for i := 0; i < len(le); i++ {
    dup := false
    for j := 0; j < len(ri); j++ {
      if le[i].NodeID.Equals(ri[j].NodeID) {
        dup = true
        break
      }
    }

    if !dup {
      res = append(res, le[i])
    }
  }
  return
}


func firstKEle(l []Contact, k int) (res []Contact) {
  if k > len(l) {
    res = l[:]
  } else if k < 0 {
    res = l[:0]
  } else {
    res = l[:k]
  }
  return
}
