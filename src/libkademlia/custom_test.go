package libkademlia

import (
	"bytes"
	//"net"
	"strconv"
	"testing"
	//"time"
	"log"
	//"fmt"
)

func TestFindKNodes(t *testing.T) {
	// tree structure;
	// A->B->tree
	/*
	       C
	      /
	  A-B -- D
	      \
	       E
	*/
	instance1 := NewKademlia("localhost:9000")
	instance2 := NewKademlia("localhost:9001")
	host2, port2, _ := StringToIpPort("localhost:9001")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	tree_node := make([]*Kademlia, 30)
	for i := 0; i < 30; i++ {
		address := "localhost:" + strconv.Itoa(9002+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
	}
	key := NewRandomID()
	contacts, err := instance1.DoFindNode(contact2, key)
	if err != nil {
		t.Error("Error doing FindNode")
	}

	if contacts == nil || len(contacts) != 20 {
		t.Error("Did not Found K (20) nodes")
	}

	// TODO: Check that the correct contacts were stored
	//       (and no other contacts)
	log.Println("\n -----\n", "Test FindKNodes passed!\n", "-----\n")
	return
}

func TestStoreConflict(t *testing.T) {
	// test Dostore() function and LocalFindValue() function
	instance1 := NewKademlia("localhost:9050")
	instance2 := NewKademlia("localhost:9051")
	host2, port2, _ := StringToIpPort("localhost:9051")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	key := NewRandomID()
	value := []byte("Hello")
	err = instance1.DoStore(contact2, key, value)
	if err != nil {
		t.Error("Can not store this value")
	}

  // key2 := CopyID(key)
  // key = NewRandomID()
  // report:
  // rpc: gob error encoding body: gob: type not registered for interface: errors.errorString
  // see more in comment
	// value = []byte("World")
	// err = instance1.DoStore(contact2, key, value)
	// if err == nil {
	// 	t.Error("This error is introduced by incomplete rpc interface, See more in our document")
  //   t.Error("Call: ", err)
	// }

	storedValue, err := instance2.LocalFindValue(key)
	if err != nil {
		t.Error("Stored value not found!")
	}
	if !bytes.Equal(storedValue, value) {
		t.Error("Stored value did not match found value")
	}
	log.Println("\n -----\n", "Test Store passed!\n", "-----\n")
	return
}
