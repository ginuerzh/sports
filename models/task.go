// task
package models

import (
//"time"
)

const (
	TaskRunning = "PHYSIQUE"
	TaskPost    = "LITERATURE"
	TaskGame    = "MAGIC"
)

/*
const (
	TaskCompleted = iota
	TaskUncompleted
)

type Proof struct {
	Tid    int
	Pics   []string
	Result string `bson:",omitempty"`
}

type TaskList struct {
	Completed   []int
	Uncompleted []int
	Waited      []int
	Proofs      []Proof
	Last        time.Time
}

func (tl *TaskList) TaskStatus(tid int) (status string) {
	for i := len(tl.Completed) - 1; i >= 0; i-- {
		if tid == tl.Completed[i] {
			return "FINISH"
		}
	}
	for _, id := range tl.Waited {
		if tid == id {
			return "AUTHENTICATION"
		}
	}
	for _, id := range tl.Uncompleted {
		if tid == id {
			return "UNFINISH"
		}
	}

	return "NORMAL"
}

func (tl *TaskList) GetProof(tid int) Proof {
	for _, proof := range tl.Proofs {
		if proof.Tid == tid {
			return proof
		}
	}
	return Proof{}
}
*/
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
	{Id: 8, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复8次"},
	{Id: 9, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复6次"},
	{Id: 10, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复7次"},
	{Id: 11, Type: TaskPost, Desc: "发表一篇运动日志"},
	{Id: 12, Type: TaskPost, Desc: "发表一篇运动日志"},
	{Id: 13, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 14, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 3
	{Id: 15, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复8次"},
	{Id: 16, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复6次"},
	{Id: 17, Type: TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复7次"},
	{Id: 18, Type: TaskPost, Desc: "发表一篇运动日志"},
	{Id: 19, Type: TaskPost, Desc: "发表一篇运动日志"},
	{Id: 20, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 21, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
}
