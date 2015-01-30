// task
package models

import (
//"time"
)

const (
	TaskRunning = "PHYSIQUE"
	TaskPost    = "LITERATURE"
	TaskGame    = "MAGIC"
	TaskNormal  = "RUNNING"
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
	Id       int64    `json:"task_id"`
	Type     string   `json:"task_type"`
	Distance int      `json:"-"`
	Duration int64    `json:"-"`
	Desc     string   `json:"task_desc"`
	Status   string   `json:"task_status"`
	Pics     []string `json:"task_pics,omitempty"`
	Result   string   `json:"task_result,omitempty"`
	Tip      string   `json:"task_tip,omitempty"`
}

var Tasks = []Task{
	// week 1
	{Id: 1, Type: TaskPost, Desc: "新手教程"},
	{Id: 2, Type: TaskRunning,
		Distance: 2000, Duration: 1500, Desc: "慢跑1分钟,行走2分钟,重复8次。\n距离：2公里, 时长：25分钟。",
		Tip: `初看这个任务，很多人都觉目标太低，没有多大意义。事实上，我们现在开始需要的是重新思考，再次出发，练习正确的跑步方式，让自己最终能够轻松的完成跑完全程马拉松的目标。跑步需要倾听自己内在的声音，而不是更多地的受到外部目标的影响。你的身体是你的老师和学生，了解身体，这称作“身体感知”。`},
	{Id: 3, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 4, Type: TaskRunning,
		Distance: 2000, Duration: 1200, Desc: "慢跑1分钟,行走2分钟,重复6次。\n距离：2公里, 时长：20分钟。",
		Tip: `现在思考一下正确跑步的方式，我们后面会持续强化以下几个方面：
1. 正确的姿势
2. 放松的四肢
3. 自由活动的关节
4. 核心肌肉的使用
5. 专注的意念
6. 良好的呼吸技术。`},
	{Id: 5, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 6, Type: TaskRunning,
		Distance: 2000, Duration: 1320, Desc: "慢跑1分钟,行走2分钟,重复7次。\n距离：2公里, 时长：22分钟。",
		Tip: `应对一种力量的最佳方式是顺其自然。我们跑步时，身体受到两种力量的影响：向下的重力和向前跑动时道路传来的反向力。跑步时，让身体保持轻微前倾，找到平衡，你就学会了和重力合作，当重力向前拉你走时，你跟着走，让小腿尽可能休息。让我们的脚在身体重心略微靠后一点位置落地，让重心位于落脚点之前，以柔和的全脚掌着地代替脚跟着地，可以化解道路的反向力，而且不会降低速度，也不会对身体造成冲击。这个技术非常关键，需要认真思考，如果有疑问，可以到“跑步圣经”圈子里发帖，和有经验的跑步爱好者一起讨论。`},
	{Id: 7, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	//week 2
	{Id: 8, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 9, Type: TaskRunning,
		Distance: 2500, Duration: 1800, Desc: "慢跑2分钟,行走2分钟,重复7次。\n距离：2.5公里, 时长：30分钟。",
		Tip: `跑步的主要原则：
1. 绵里藏针：让身体成一条直线，四肢要像棉花一样柔软，毫不紧张，学习如何用身体的中心部位带动跑步。想象在你的头部和尾椎骨之间有一条中心线，了解它，感知它，在你的身体找到这个中心，感觉找到这个中心。
2. 循序渐进：在开始阶段一定要慢慢起步，一点点加速让我们身体逐步适应。这是我们做任何事情都必须遵守的原则。
3. 动作中的平衡：在跑步时，我们需要让身体所有部位在一个整体中运行，所有部位都协调工作，以平衡的方式运动。寻找平衡的关键是感知自己的中心线在哪里。`},
	{Id: 10, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 11, Type: TaskRunning,
		Distance: 2500, Duration: 1500, Desc: "慢跑1分钟,行走2分钟,重复7次。\n距离：2.5公里, 时长：25分钟。",
		Tip: `身体感知三个秘诀：
1. 仔细倾听：每当我们对跑姿调整时，都要倾听任何细微的变化：身体是怎样运动的?感觉如何?各部分是否协调？
2. 信息评估：问问自己，身体是否正按着我们希望的方式遇到弄，尽可能辨别我们所作的每次调整是否有效。
3. 逐步调整：进行微调从来都是最好的策略，任何突然的改变都可能对身体造成伤害。`},
	{Id: 12, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 13, Type: TaskRunning,
		Distance: 2500, Duration: 1500, Desc: "慢跑2分钟,行走2分钟,重复6次。\n距离：2.5公里, 时长：25分钟。",
		Tip: `呼气，挖掘你的“气”：呼吸在跑步中起着提供氧气的关键作用，我们要改变日常活动的呼吸方式，要采用腹式呼吸，尽可能吐完空气，放松腹部，保持和步频匹配，一般是三步呼出两步吸入。在进行腹式呼吸时用鼻子呼吸效果会更好，尽可能不用嘴呼吸。`},
	{Id: 14, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 3
	{Id: 15, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 16, Type: TaskRunning,
		Distance: 3000, Duration: 2100, Desc: "慢跑3分钟,行走2分钟,重复7次。\n距离：3公里, 时长：35分钟。",
		Tip: `跑步时，尽量使身体的各部分向同一方向运动。当身体对直时，我们就有一个驱动身体的中心线，这条线笔直和强壮时，它就成了支撑你身体的“针”，你的双臂和双腿就变成了“棉”。当你姿态对直时，身体的重量就是由你的结构而不是肌肉来支撑了，能量或“气”会毫无障碍地流过你的身体，让你不必付出太大努力的情况下跑的更快，更远。`},
	{Id: 17, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 18, Type: TaskRunning,
		Distance: 3000, Duration: 1500, Desc: "慢跑2分钟,行走2分钟,重复6次。\n距离：3公里, 时长：25分钟。",
		Tip: `对直腿和脚：
1. 双脚分开、对直，脚尖指向前方，平行站立。要使双脚对直，不仅仅是让脚尖指向前方，而是要将双腿向内侧扭转，直到腿的正面对直脚尖指向前方，不再是“八”字脚。
2. 平衡双腿的压力，确保你身体的重量是平均的分配在左右脚之间的，然后在双腿间进行平衡，再在双脚间平衡，感觉两侧压力是平均的。
3. 在跑动时，保持脚尖指向前方，假想在地上有一条直线，而你的双脚就平行地在并列在直线两边。`},
	{Id: 19, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 20, Type: TaskRunning,
		Distance: 3000, Duration: 1800, Desc: "慢跑3分钟,行走2分钟,重复6次。\n距离：3公里, 时长：30分钟。",
		Tip: `延长脊柱来对直上半身：向上提后颈部，延长整个脊柱，让身体的中心线又长又直，同时放松膝盖，扩张胸部，以便更充分的呼吸。`},
	{Id: 21, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 4
	{Id: 22, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 23, Type: TaskRunning,
		Distance: 3500, Duration: 2100, Desc: "慢跑3分钟,行走2分钟,重复6次。\n距离：3.5公里, 时长：35分钟。",
		Tip: `保持骨盘水平，利用核心力量，可以强化小腹肌肉，专注自己中心线部分，感受到身体和能量的集中：
1. 在跑步时，保持一条笔直的身体中心线。
2. 在运动中固定住骨盘。
3. 在骨盘和双腿间建立起更加强有力的连接，使整个下半身的运动成为一个整体。`},
	{Id: 24, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 25, Type: TaskRunning,
		Distance: 3500, Duration: 1800, Desc: "慢跑2分钟,行走2分钟,重复5次。\n距离：3.5公里, 时长：30分钟。",
		Tip: `创建身体立柱：让我们的髋部和肩膀对齐，上半身在髋部的正上方，保持平衡。肩膀，髋部和脚部将形成身体的立柱。当身体立柱对直时，身体的重量将有我们的结构支撑，肌肉就可以做它该做的事。`},
	{Id: 26, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 27, Type: TaskRunning,
		Distance: 3500, Duration: 2400, Desc: "慢跑3分钟,行走3分钟,重复6次。\n距离：3.5公里, 时长：40分钟。",
		Tip: `单腿站姿：要通过单独站立训练你的核心肌肉，让我们在跑步的过程中保持笔直的身体姿势，因为在跑步时，我们其实是在进行一系列的单腿站立活动。`},
	{Id: 28, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 5
	{Id: 29, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 30, Type: TaskRunning,
		Distance: 4000, Duration: 2400, Desc: "慢跑3分钟,行走1分钟,重复9次。\n距离：4公里, 时长：40分钟。",
		Tip: `学会如何前倾：身体前倾可以让重力成为我们的好帮手，身体一旦前倾，重心就会落到身体着地点的前面，重力会把我们向下拉使你向前。我们的脚就会像弹跳器发挥作用，让我们前进，要一直在这种状态下保持细微的平衡。`},
	{Id: 31, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 32, Type: TaskRunning,
		Distance: 4000, Duration: 1500, Desc: "慢跑2分钟,行走1分钟,重复8次。\n距离：4公里, 时长：25分钟。",
		Tip: `全脚掌着地：不要一边跑一边制动，不单单是脚跟着地，而是你的整个脚底着地，从前至后，从一侧到另一侧的脚底压力是平均，这样使你的双腿得以充分放松而不用负责推动身体。你的腿只用于在步幅之间短暂地支撑一下，支撑以后，腿部会沿着道路的反向力的方向朝后方摆动。`},
	{Id: 33, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 34, Type: TaskRunning,
		Distance: 3500, Duration: 1800, Desc: "慢跑3分钟,行走1分钟,重复8次。\n距离：4公里, 时长：30分钟。",
		Tip: `摆臂：
1. 屈肘90度：以轻松的方式自肩部摆动，不然胳膊一张一合，保持90度。
2. 向后摆臂，而不是向前。想象着你是用肘击后面的人而不是用拳打前面的，与身体前倾保持相对平衡。
3. 不要让手越过身体中心线：想象双手之间抱着一个篮球，别让双手之间的距离小于这个篮球的宽度。
4. 放松双手：保持手指向内弯曲，拇指在上的动作。`},
	{Id: 35, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 6
	{Id: 36, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 37, Type: TaskRunning,
		Distance: 4500, Duration: 2700, Desc: "慢跑5分钟,行走1分钟,重复7次。\n距离：4.5公里, 时长：45分钟。",
		Tip: `正常摆臂跑.跑步是锻炼骨骼的好方法，所以你有必要补充充足的钙质——每天1000毫克。如果你在50岁以上，则每天需要1500毫克。低脂牛奶，低脂酸奶和深绿色叶片蔬菜都是钙质的重要来源`},
	{Id: 38, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 39, Type: TaskRunning,
		Distance: 4500, Duration: 2100, Desc: "慢跑3分钟,行走1分钟,重复7次。\n距离：4.5公里, 时长：35分钟。"},
	{Id: 40, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 41, Type: TaskRunning,
		Distance: 4500, Duration: 2700, Desc: "慢跑3分钟,行走1分钟,重复10次。\n距离：4.5公里, 时长：45分钟。"},
	{Id: 42, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 7
	{Id: 43, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 44, Type: TaskRunning,
		Distance: 5500, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：5.5公里, 时长：45分钟。"},
	{Id: 45, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 46, Type: TaskRunning,
		Distance: 5000, Duration: 2700, Desc: "慢跑4分钟,行走1分钟,重复6次。\n距离：5公里, 时长：45分钟。"},
	{Id: 47, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 48, Type: TaskRunning,
		Distance: 5000, Duration: 2700, Desc: "慢跑5分钟,行走1分钟,重复7次。\n距离：5公里, 时长：45分钟。"},
	{Id: 49, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 8
	{Id: 50, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 51, Type: TaskRunning,
		Distance: 5500, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：5.5公里, 时长：45分钟。"},
	{Id: 52, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 53, Type: TaskRunning,
		Distance: 5500, Duration: 2400, Desc: "慢跑3分钟,行走1分钟,重复7次。\n距离：5.5公里, 时长：40分钟。"},
	{Id: 54, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 55, Type: TaskRunning,
		Distance: 5500, Duration: 2400, Desc: "慢跑5分钟,行走1分钟,重复6次。\n距离：5.5公里, 时长：40分钟。"},
	{Id: 56, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 9
	{Id: 57, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 58, Type: TaskRunning,
		Distance: 7000, Duration: 3600, Desc: "慢跑10分钟步行1分钟,慢跑15分钟步行1分钟,慢跑20分钟步行1分钟,慢跑10分钟。\n距离：7公里, 时长：60分钟。"},
	{Id: 59, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 60, Type: TaskRunning,
		Distance: 5000, Duration: 2280, Desc: "慢跑5分钟,行走1分钟,重复6次。\n距离：5公里, 时长：38分钟。"},
	{Id: 61, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 62, Type: TaskRunning,
		Distance: 5000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：5公里, 时长：45分钟。"},
	{Id: 63, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 10
	{Id: 64, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 65, Type: TaskRunning,
		Distance: 8000, Duration: 4200, Desc: "慢跑10分钟步行1分钟,慢跑20分钟步行1分钟,慢跑30分钟。\n距离：8公里, 时长：70分钟。"},
	{Id: 66, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 67, Type: TaskRunning,
		Distance: 6000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：6公里, 时长：45分钟。"},
	{Id: 68, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 69, Type: TaskRunning,
		Distance: 6000, Duration: 3000, Desc: "慢跑20分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。\n距离：6公里, 时长：50分钟。"},
	{Id: 70, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 11
	{Id: 71, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 72, Type: TaskRunning,
		Distance: 9000, Duration: 3900, Desc: "慢跑40分钟步行1分钟,慢跑20分钟。\n距离：9公里, 时长：65分钟。"},
	{Id: 73, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 74, Type: TaskRunning,
		Distance: 6000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：6公里, 时长：45分钟。"},
	{Id: 75, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 76, Type: TaskRunning,
		Distance: 8000, Duration: 3000, Desc: "慢跑20分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。\n距离：8公里, 时长：50分钟。"},
	{Id: 77, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 12
	{Id: 78, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 79, Type: TaskRunning,
		Distance: 10000, Duration: 2400, Desc: "慢跑50分钟。\n距离：10公里, 时长：40分钟。"},
	{Id: 80, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 81, Type: TaskRunning,
		Distance: 5000, Duration: 1980, Desc: "慢跑10分钟,行走1分钟,重复3次。\n距离：5公里, 时长：33分钟。"},
	{Id: 82, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 83, Type: TaskRunning,
		Distance: 8000, Duration: 2700, Desc: "慢跑15分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。\n距离：8公里, 时长：45分钟。"},
	{Id: 84, Type: TaskGame, Desc: "玩个游戏放松一下吧"},

	// week 13
	{Id: 85, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 86, Type: TaskRunning,
		Distance: 10000, Duration: 2400, Desc: "慢跑40分钟。\n距离：10公里, 时长：40分钟。"},
	{Id: 87, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 88, Type: TaskRunning,
		Distance: 5000, Duration: 1980, Desc: "慢跑10分钟,行走1分钟,重复3次。\n距离：5公里, 时长：33分钟。"},
	{Id: 89, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 90, Type: TaskRunning,
		Distance: 10000, Duration: 2400, Desc: "慢跑40分钟。\n距离：10公里, 时长：40分钟。"},
	{Id: 91, Type: TaskGame, Desc: "玩个游戏放松一下吧"},
}
