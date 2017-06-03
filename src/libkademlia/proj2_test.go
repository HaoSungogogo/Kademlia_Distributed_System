package libkademlia

import (
	//"net"
	"strconv"
	"testing"
	//"time"
	"log"
	//"fmt"
)

func TestIterativeFindNode(t *testing.T) {
	instance1 := NewKademlia("localhost:9104")
	instance2 := NewKademlia("localhost:9105")
	host2, port2, _ := StringToIpPort("localhost:9105")
	instance1.DoPing(host2, port2)

	tree_node := make([]*Kademlia, 20)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(9106+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
	}

	contacts, err := tree_node[0].DoIterativeFindNode(tree_node[3].NodeID)

	if err == nil {
		for i := 0; i < len(contacts); i++ {
			log.Println(contacts[i].NodeID.AsString())
		}
	}
	log.Println(tree_node[3].NodeID.AsString())

	log.Println("\n -----\n", "Test TestIterativeFindNode passed!\n", "-----\n")
	// t.Error("test error")
}

func TestIterativeStore(t *testing.T) {
	instance1 := NewKademlia("localhost:9204")
	instance2 := NewKademlia("localhost:9205")
	host2, port2, _ := StringToIpPort("localhost:9205")
	instance1.DoPing(host2, port2)

	tree_node := make([]*Kademlia, 20)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(9206+i)

		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
	}

	key := tree_node[3].NodeID
	value := []byte("Hello World")

	tree_node[0].DoIterativeStore(key, value)

	for i := 0; i < 10; i++ {
		storedValue, err := tree_node[i].LocalFindValue(key)
		if err == nil {
			log.Println(string(storedValue))
		}
	}

	log.Println("\n -----\n", "Test TestIterativeStore passed!\n", "-----\n")
}

func TestIterativeFindValue(t *testing.T) {
	instance1 := NewKademlia("localhost:9304")
	instance2 := NewKademlia("localhost:9305")
	host2, port2, _ := StringToIpPort("localhost:9305")
	instance1.DoPing(host2, port2)

	tree_node := make([]*Kademlia, 20)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(9306+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
	}
	key := tree_node[3].NodeID
	value := []byte("Hello World @ DoIterativeFindValue")

	instance2.DoStore(&tree_node[3].SelfContact, key, value)

	foundValue, err := tree_node[0].DoIterativeFindValue(key)

	if err == nil {
		log.Println(string(foundValue))
	}

	log.Println("\n -----\n", "Test TestIterativeFindValue passed!\n", "-----\n")
	// t.Error("test error")

}
