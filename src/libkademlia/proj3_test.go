package libkademlia

import (
	"log"
	"testing"
	"strconv"
)

func TestVanish(t *testing.T) {
	data := []byte("Hello World")

	instance1 := NewKademlia("localhost:9404")
	instance2 := NewKademlia("localhost:9405")

	host2, port2, _ := StringToIpPort("localhost:9405")
	instance1.DoPing(host2, port2)

	tree_node := make([]*Kademlia, 20)
	for i := 0; i < 20; i++ {
		address := "localhost:" + strconv.Itoa(9406+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
	}

	_, vdoid := tree_node[3].Vanish(data, byte(20), byte(10), 0)

	vdo_data := instance2.Unvanish(tree_node[3].NodeID, vdoid)

	log.Println("vdo_data", vdo_data)
	log.Println("data", data)

	if (string(vdo_data) != string(data)) {
		t.Error("test error")
	}


}
