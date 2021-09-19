package storage

// sort.Sort(ByMsgID([]ChatMessage{}))

type ByMsgID []ChatMessage

func (msg ByMsgID) Len() int {
	return len(msg)
}

func (msg ByMsgID) Swap(i, j int) {
	msg[i], msg[j] = msg[j], msg[i]
}

func (msg ByMsgID) Less(i, j int) bool {
	return msg[i].ID() < msg[j].ID()
}
