// task
package models

const (
	TaskRunning = "PHYSIQUE"
	TaskPost    = "LITERATURE"
	TaskGame    = "MAGIC"
)

const (
	TaskCompleted = iota
	TaskUncompleted
)

type Task struct {
	Id     int      `json:"task_id"`
	Type   string   `json:"task_type"`
	Desc   string   `json:"task_desc"`
	Status string   `json:"task_status"`
	Pics   []string `json:"task_pics,omitempty"`
	Result string   `json:"task_result,omitempty"`
}

var Tasks = []Task{
	// week 1
	{Id: 1, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复8次"},
	{Id: 2, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复6次"},
	{Id: 3, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复7次"},
	{Id: 4, Type: TaskPost, Desc: "发表一篇运动日志"},
	{Id: 5, Type: TaskPost, Desc: "发表一篇运动日志"},
	{Id: 6, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 7, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	//week 2
}
