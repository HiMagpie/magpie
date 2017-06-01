package models

/**
 * msg ack after pushed
 */
type MsgAck struct {
	MsgId int64 `json:"msg_id"`
	Cid string `json:"cid"`
}

func NewMsgAck() *MsgAck {
	return new(MsgAck)
}