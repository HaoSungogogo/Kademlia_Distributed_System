package libkademlia

import (
  "sort"
  "fmt"
  //"log"
)

type ParallelDoFindNodeRes struct {
  Target    *Contact
  Nodes     []Contact
  Err       bool
}

func (k *Kademlia) ParallelDoFindNode (target *Contact, key ID, resChan chan ParallelDoFindNodeRes) {
  contacts, err := k.DoFindNode(target, key)
  if err != nil {
    resChan <- ParallelDoFindNodeRes{ target, contacts, true }
  } else {
    resChan <- ParallelDoFindNodeRes{ target, contacts, false }
  }
}

type ParallelDoFindValueRes struct {
  Target    *Contact
  Value     []byte
  Nodes     []Contact
  Err       bool
}

func (k *Kademlia) ParallelDoFindValue (target *Contact, key ID, resChan chan ParallelDoFindValueRes) {
  value, contacts, err := k.DoFindValue(target, key)
  resChan <- ParallelDoFindValueRes{ target, value, contacts, err != nil }
}

type ParallelDoPingRes struct {
  Target    *Contact
  Err       bool
}

func (k *Kademlia) ParallelDoPing (target *Contact, resChan chan ParallelDoPingRes) {
  _, err := k.DoPing(target.Host, target.Port)
  if err != nil {
    resChan <- ParallelDoPingRes{ target, true }
  } else {
    resChan <- ParallelDoPingRes{ target, false }
  }
}

type ContactTargetSlice struct {
  Contacts  []Contact
  Target    ID
}

func (self ContactTargetSlice) Less(i, j int) bool {
  return self.Contacts[i].NodeID.Xor(self.Target).Less(self.Contacts[j].NodeID.Xor(self.Target))
}

func (self ContactTargetSlice) Len() int {
  return len(self.Contacts)
}

func (self ContactTargetSlice) Swap(i, j int) {
  self.Contacts[i], self.Contacts[j] = self.Contacts[j], self.Contacts[i]
}

func toContactTargetSlice(l []Contact, target ID) (contacts ContactTargetSlice) {
  contacts.Contacts = l
  contacts.Target = target
  return
}



// For project 2!
func (k *Kademlia) DoIterativeFindNode(id ID) (shortlist []Contact, err error) {
  findNodeCmd := findNodeCommand{ id, make(chan FindNodeResult), alpha }
	k.findNodeChan <- findNodeCmd

  findNodeRes := <- findNodeCmd.ResChan
  alphaNodes := findNodeRes.Nodes

  var bufferContacts []Contact

  stop := false
  for (!stop) {

    ParallelDoFindNodeResChan := make(chan ParallelDoFindNodeRes)

    for i := 0; i < len(alphaNodes); i++ {
      go k.ParallelDoFindNode(&alphaNodes[i], id, ParallelDoFindNodeResChan)
    }

    for i := 0; i < len(alphaNodes); i++ {
      currDoFindNodeRes := <- ParallelDoFindNodeResChan

      if !currDoFindNodeRes.Err {
        bufferContacts = append(bufferContacts, currDoFindNodeRes.Nodes...)

        shortlist = append(shortlist, *currDoFindNodeRes.Target)
      }
    }
    temp := toContactTargetSlice(shortlist, id)
    sort.Sort(temp)
    shortlist = temp.Contacts

    temp = toContactTargetSlice(bufferContacts, id)
    sort.Sort(temp)
    bufferContacts = temp.Contacts

    shortlist = removeDup(shortlist)
    bufferContacts = removeDup(bufferContacts)
    bufferContacts = minus(bufferContacts, shortlist)

    bufferContacts = firstKEle(bufferContacts, kMax - len(shortlist))

    if len(shortlist) >= kMax {
      stop = true

    } else if len(shortlist) > 0 && id.Xor(bufferContacts[0].NodeID).Less(id.Xor(shortlist[0].NodeID)) {

      alphaNodes = firstKEle(bufferContacts, alpha)

    } else {
      ParallelDoPingResChan := make(chan ParallelDoPingRes)

      for i := 0; i < len(bufferContacts); i++ {
        go k.ParallelDoPing(&bufferContacts[i], ParallelDoPingResChan)
      }

      for i := 0; i < len(bufferContacts); i++ {
        currDoPingRes := <- ParallelDoPingResChan

        if !currDoPingRes.Err {
          shortlist = append(shortlist, *currDoPingRes.Target)
        }
      }
      stop = true
    }
  }

  shortlist = firstKEle(shortlist, kMax)
  return
}

type DoStoreRes struct {
  Err     bool
  Target  *Contact
}

func (k *Kademlia) DoIterativeStore(key ID, value []byte) (res []Contact, err error) {
  contacts, err := k.DoIterativeFindNode(key)

  if err == nil {
    DoStoreResChan := make(chan DoStoreRes)

    for i := 0; i < len(contacts); i++ {

      currContact := &contacts[i]
      go func () {
        // log.Println(contacts[i].NodeID);
        err = k.DoStore(currContact, key, value)
        if err == nil {
          DoStoreResChan <- DoStoreRes{ false, currContact }
        } else {
          DoStoreResChan <- DoStoreRes{ true, nil }
        }
      } ()
    }

    for i := 0; i < len(contacts); i++ {
      doStoreRes := <- DoStoreResChan
      if !doStoreRes.Err {
        res = append(res, *doStoreRes.Target)
      }
    }
    return res, nil
  }
  return nil, &CommandFailed{ "DoIterativeStore Failed" }
}



func (k *Kademlia) DoIterativeFindValue(key ID) (value []byte, err error) {
  //should be a return value but the assignment doesn't
  var shortlist []Contact
  //check if the key is stored locally
  localvalue, err := k.LocalFindValue(key)

  if err == nil {
    return localvalue, nil
  }

  //find node closest to the key
  findNodeCmd := findNodeCommand{ key, make(chan FindNodeResult), alpha }
  k.findNodeChan <- findNodeCmd

  findNodeRes := <- findNodeCmd.ResChan
  alphaNodes := findNodeRes.Nodes

  var bufferContacts []Contact

  stop := false
  for (!stop) {

    ParallelDoFindValueResChan := make(chan ParallelDoFindValueRes)

    for i := 0; i < len(alphaNodes); i++ {
      go k.ParallelDoFindValue(&alphaNodes[i], key, ParallelDoFindValueResChan)
    }

    for i := 0; i < len(alphaNodes); i++ {
      currDoFindValueRes := <- ParallelDoFindValueResChan

      //We don't have error return here due to the incomplete interface

      if currDoFindValueRes.Value != nil {

        if len(shortlist) > 0 {
          k.DoStore(&shortlist[0], key, currDoFindValueRes.Value)
        }

        return currDoFindValueRes.Value, nil
      }

      if !currDoFindValueRes.Err {
        bufferContacts = append(bufferContacts, currDoFindValueRes.Nodes...)

        shortlist = append(shortlist, *currDoFindValueRes.Target)
      }
    }
    temp := toContactTargetSlice(shortlist, key)
    sort.Sort(temp)
    shortlist = temp.Contacts

    temp = toContactTargetSlice(bufferContacts, key)
    sort.Sort(temp)
    bufferContacts = temp.Contacts

    shortlist = removeDup(shortlist)
    bufferContacts = removeDup(bufferContacts)
    bufferContacts = minus(bufferContacts, shortlist)

    bufferContacts = firstKEle(bufferContacts, kMax - len(shortlist))

    if len(shortlist) >= kMax {
      stop = true

    } else if len(shortlist) > 0 && key.Xor(bufferContacts[0].NodeID).Less(key.Xor(shortlist[0].NodeID)) {

      alphaNodes = firstKEle(bufferContacts, alpha)

    } else {
      ParallelDoPingResChan := make(chan ParallelDoPingRes)

      for i := 0; i < len(bufferContacts); i++ {
        go k.ParallelDoPing(&bufferContacts[i], ParallelDoPingResChan)
      }

      for i := 0; i < len(bufferContacts); i++ {
        currDoPingRes := <- ParallelDoPingResChan

        if !currDoPingRes.Err {
          shortlist = append(shortlist, *currDoPingRes.Target)
        }
      }
      stop = true
    }
  }

  shortlist = firstKEle(shortlist, kMax)
  return nil, &CommandFailed {
    "Unable to find value, closest node queried: " + fmt.Sprintf("%s", shortlist[0].NodeID.AsString()) }
}
