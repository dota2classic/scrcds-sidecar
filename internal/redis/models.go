package redis

type ChannelEvent[T any] struct {
	Pattern string `json:"pattern"`
	Data    T      `json:"data"`
}
