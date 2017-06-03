package libkademlia

import (
  "container/list"
  "log"
  "errors"
  "sort"
)

type KBucket struct {
  contacts list.List
}

func (k *Kademlia) update(contact Contact) {
  //log.Println("update contact called: ", contact)
  bucket := k.getBucket(k.SelfContact.NodeID.Xor(contact.NodeID))

  for ele := bucket.contacts.Front(); ele != nil; ele = ele.Next() {
    if (ele.Value.(Contact).NodeID == contact.NodeID) {
      bucket.contacts.MoveToBack(ele)
      return
    }
  }

  if (bucket.contacts.Len() < kMax) {
    bucket.contacts.PushBack(contact)
  } else {
    _, err := k.DoPing(bucket.contacts.Front().Value.(Contact).Host, uint16(bucket.contacts.Front().Value.(Contact).Port))
    if err != nil {
      bucket.contacts.Remove(bucket.contacts.Front())
      bucket.contacts.PushBack(contact)
    }
  }
}

func (k *Kademlia) getContact(NodeID ID) (ret findContactResponse) {
  //log.Println("getContact called with NodeID: ", NodeID)

  bucket := k.getBucket(k.SelfContact.NodeID.Xor(NodeID))
  //k.printBucket(bucket)

  for ele := bucket.contacts.Front(); ele != nil; ele = ele.Next() {
    if (ele.Value.(Contact).NodeID == NodeID) {
      ret =  findContactResponse{ele.Value.(Contact), nil}
      return
    }
  }
  //TODO: Find the better implementation
  con := Contact {}
  ret = findContactResponse{ con, errors.New("Contact Not found") }
  return
}

// func (c Contact)

type ContactSlice []Contact

func (contacts ContactSlice) Less(i, j int) bool {
    return contacts[i].NodeID.Less(contacts[j].NodeID)
}

func (contacts ContactSlice) Len() int {
    return len(contacts)
}

func (contacts ContactSlice) Swap(i, j int) {
    contacts[i], contacts[j] = contacts[j], contacts[i]
}

func toContactSlice(l list.List) (contacts ContactSlice) {
  for ele := l.Front(); ele != nil; ele = ele.Next() {
    contacts = append(contacts, ele.Value.(Contact))
  }
  return
}

func (k *Kademlia) getContactsFromBucket(bIndex, currIndex, stillNeed *int) (contacts []Contact) {
  if (k.rt[*currIndex].contacts.Len() <= *stillNeed) {
    // add all
    contacts = toContactSlice(k.rt[*currIndex].contacts)

    *stillNeed -= k.rt[*currIndex].contacts.Len()
    if (*currIndex < *bIndex) {
      *currIndex--
    } else if (*currIndex < b-1) {
      *currIndex++
    } else {
      *currIndex = *bIndex-1
    }
  } else { //
    temp := toContactSlice(k.rt[*currIndex].contacts)
    // find cloest k
    sort.Sort(temp)

    for i := 0; i < *stillNeed; i++ {
      contacts = append(contacts, temp[i])
    }
    *stillNeed = 0
  }
  return
}

func (k *Kademlia) getKContacts(key ID) (ret FindNodeResult) {
  ret = k.getNContacts(key, kMax)
  return
}

func (k *Kademlia) getNContacts(key ID, N int) (ret FindNodeResult) {

  bIndex := k.getBucketIndex(key)
  stillNeed := N
  currIndex := bIndex

  for (stillNeed > 0 && currIndex >= 0) {
    ret.Nodes = append(ret.Nodes, k.getContactsFromBucket(&bIndex, &currIndex, &stillNeed)...)
  }

  ret.Err = nil
  return
}

func (k *Kademlia) getLocalValue(searchKey ID) (ret findLocalValueResponse) {
  //log.Println("getContact called with NodeID: ", NodeID)
  if _, ok := k.hash[searchKey]; ok {
    ret = findLocalValueResponse{ k.hash[searchKey], nil }
  } else {
    ret = findLocalValueResponse{ nil, errors.New("Local Value Not found") }
  }
  return
}

func (k *Kademlia) getBucketIndex(dis ID) (int) {
  return dis.PrefixLen()
}

func (k *Kademlia) getBucket(dis ID) (ret *KBucket) {
  return &k.rt[k.getBucketIndex(dis)]
}

func (k *Kademlia) printBucket(bucket *KBucket) {
  log.Println("Printing bucket.")

  for ele := bucket.contacts.Front(); ele != nil; ele = ele.Next() {
    log.Println("NodeID: ", ele.Value.(Contact).NodeID)
  }

  log.Println("Finish printing.")
  return
}
