package models

/**
 * Msg structure for sending to client.
 */
type Payload struct {
	Cid       string `json:"cid"`
	Ring      bool `json:"ring"`
	Vibrate   bool `json:"vibrate"`
	Cleanable bool `json:"cleanable"`

	Trans     int `json:"trans"`

	Title     string `json:"title"`
	Text      string `json:"text"`
	Logo      string `json:"logo"`
	Url       string `json:"url"`

	MsgId     int64 `json:"msg_id"`
	Ctime     string `json:"ctime"`
}
