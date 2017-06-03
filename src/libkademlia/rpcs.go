package libkademlia

// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"net"
	//"log"
)

type KademliaRPC struct {
	kademlia *Kademlia
}

// Host identification.
type Contact struct {
	NodeID 			ID
	Host   			net.IP
	Port   			uint16
	active			bool
}

///////////////////////////////////////////////////////////////////////////////
// PING
///////////////////////////////////////////////////////////////////////////////
type PingMessage struct {
	Sender Contact
	MsgID  ID
}

type PongMessage struct {
	MsgID  ID
	Sender Contact
}

func (k *KademliaRPC) Ping(ping PingMessage, pong *PongMessage) error {
	// TODO: Finish implementation
	//log.Println("Ping called.")
	pong.MsgID = CopyID(ping.MsgID)

	// TODO: might have problem, race it
	pong.Sender = k.kademlia.SelfContact

	pingCmd := updateCommand{ ping.Sender }
	k.kademlia.updateChan <- pingCmd

	//log.Println("command sent")
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// STORE
///////////////////////////////////////////////////////////////////////////////
type StoreRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
	Value  []byte
}

type StoreResult struct {
	MsgID 	ID
	Err   	error
}

func (k *KademliaRPC) Store(req StoreRequest, res *StoreResult) error {
	// TODO: Implement.

	updateCmd := updateCommand{ req.Sender }
	k.kademlia.updateChan <- updateCmd

	storeCmd := storeCommand{ req.Key, req.Value, make(chan error) }
	k.kademlia.storeChan <- storeCmd

	res.MsgID = CopyID(req.MsgID)
	res.Err = <- storeCmd.ErrChan
	// log.Println("adsfasdfsadfadsfdsfasdf", res.Err)

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// FIND_NODE
///////////////////////////////////////////////////////////////////////////////
type FindNodeRequest struct {
	Sender Contact
	MsgID  ID
	NodeID ID
}

type FindNodeResult struct {
	MsgID ID
	Nodes []Contact
	Err   error
}

func (k *KademliaRPC) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	// TODO: Implement.

	updateCmd := updateCommand{ req.Sender }
	k.kademlia.updateChan <- updateCmd

	findNodeCmd := findNodeCommand{ req.NodeID, make(chan FindNodeResult), kMax }
	k.kademlia.findNodeChan <- findNodeCmd

	*res = <- findNodeCmd.ResChan
	res.MsgID = CopyID(req.MsgID)

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// FIND_VALUE
///////////////////////////////////////////////////////////////////////////////
type FindValueRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
	MsgID ID
	Value []byte
	Nodes []Contact
	Err   error
}

func (k *KademliaRPC) FindValue(req FindValueRequest, res *FindValueResult) error {
	// TODO: Implement.

	// parms(Contact, Key)
	// ret([]byte, []Contact)

	updateCmd := updateCommand{ req.Sender }
	k.kademlia.updateChan <- updateCmd

	findValueCmd := findLocalValueCommand{ req.Key, make(chan findLocalValueResponse) }
	k.kademlia.findLocalValueChan <- findValueCmd

	FindValueRes := <- findValueCmd.LocalValueChan
	if (FindValueRes.Err != nil) {
		findNodeCmd := findNodeCommand{ req.Key, make(chan FindNodeResult), kMax }
		k.kademlia.findNodeChan <- findNodeCmd

		findNodeRes := <- findNodeCmd.ResChan
		res.Nodes = findNodeRes.Nodes
	} else {
		res.Value = FindValueRes.Result
	}

	res.MsgID = CopyID(req.MsgID)

	return nil
}

// For Project 3

type GetVDORequest struct {
	Sender Contact
	VdoID  ID
	MsgID  ID
}

type GetVDOResult struct {
	MsgID ID
	VDO   VanashingDataObject
}

func (k *KademliaRPC) GetVDO(req GetVDORequest, res *GetVDOResult) error {
	ResChan := make(chan getVDOResponse)
	k.kademlia.getVDOChan <- getVDOCommand{req.VdoID, ResChan}

	getVDOResponse := <- ResChan

	res.MsgID = req.MsgID

	if getVDOResponse.Err == nil {
		res.VDO = getVDOResponse.VDO
	}

	return getVDOResponse.Err
}
