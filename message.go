package queue

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
)

type Message struct {
	Id string `json:"id"`
	CreateTime time.Time `json:"createTime""`
	ConsumeTime time.Time `json:"consumeTime"`
	Body interface{} `json:"body"`
}

// NewMessage 创建消息实体
func NewMessage(id string, consumeTime time.Time, body interface{}) *Message {
	if id == "" {
		id = uuid.NewV4().String()
	}
	return &Message{
		Id: id,
		CreateTime: time.Now(),
		ConsumeTime: consumeTime,
		Body: body,
	}
}

func (m *Message) GetScore() float64 {
	return float64(m.ConsumeTime.Unix())
}

func (m *Message) GetId() string {
	return m.Id
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
