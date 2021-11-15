package model

const (
	CounterCmdTotalInc     = "inc_total_counters"
	CounterCmdCursorInc    = "inc_cursor_counters"
	CounterCmdCursorUpdate = "update_cursor_counter"
)

type CounterCmd struct {
	Command     string   `json:"command"`
	Subscribers []string `json:"subscribers"`
}
