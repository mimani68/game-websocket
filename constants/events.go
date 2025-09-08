package constants

type EventType string

const (
	EventTypeBroadcast EventType = "core.game.state.broadcast"
	EventTypeSendMsg   EventType = "core.game.state.send_message"
	EventTypeBuyGood   EventType = "core.game.state.buy_good"
)
