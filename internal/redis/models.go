package redis

type ChannelEvent[T any] struct {
	Pattern string `json:"pattern"`
	Id      string `json:"id"`
	Data    T      `json:"data"`
}
