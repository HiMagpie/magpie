package models

/**
 * msg entity
 * means: content of the msg.  
 */
type MsgEntity struct {
	MsgId int64 `json:"msg_id"`
	Ring      bool `json:"ring"`
	Vibrate   bool `json:"vibrate"`
	Cleanable bool `json:"cleanable"`
	Trans     int `json:"trans"`
	Begin     string `json:"begin"`
	End       string `json:"end"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Logo      string `json:"logo"`
	Url       string `json:"url"`
	Ctime     string `json:"ctime"`
}

func NewMsgEntity() *MsgEntity {
	return new(MsgEntity)
}
