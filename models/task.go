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
	Id        int64    `json:"task_id"`
	Type      string   `json:"task_type"`
	Index     int      `json:"-"`
	BeginTime int64    `json:"begin_time,omitempty"`
	EndTime   int64    `json:"end_time,omitempty"`
	Distance  int      `json:"distance,omitempty"`
	Duration  int64    `json:"duration,omitempty"`
	Source    string   `json:"source,omitempty"`
	Desc      string   `json:"task_desc"`
	Video     string   `json:"task_video"`
	Status    string   `json:"task_status"`
	Pics      []string `json:"task_pics,omitempty"`
	Result    string   `json:"task_result,omitempty"`
	Tip       string   `json:"task_tip,omitempty"`
}

var Tasks = []Task{
	// week 1
	{Id: 1, Type: TaskPost, Desc: "新手教程"},
	{Id: 2, Type: TaskRunning,
		Distance: 2000, Duration: 1500, Desc: "慢跑1分钟,行走2分钟,重复8次。\n距离：2公里, 时长：25分钟。",
		Tip: tips[0]},
	{Id: 3, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 4, Type: TaskRunning,
		Distance: 2000, Duration: 1200, Desc: "慢跑1分钟,行走2分钟,重复6次。\n距离：2公里, 时长：20分钟。",
		Tip: tips[1]},
	{Id: 5, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 6, Type: TaskRunning,
		Distance: 2000, Duration: 1320, Desc: "慢跑1分钟,行走2分钟,重复7次。\n距离：2公里, 时长：22分钟。",
		Tip: tips[2]},
	{Id: 7, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	//week 2
	{Id: 8, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 9, Type: TaskRunning,
		Distance: 2500, Duration: 1800, Desc: "慢跑2分钟,行走2分钟,重复7次。\n距离：2.5公里, 时长：30分钟。",
		Tip: tips[3]},
	{Id: 11, Type: TaskRunning,
		Distance: 2500, Duration: 1500, Desc: "慢跑1分钟,行走2分钟,重复7次。\n距离：2.5公里, 时长：25分钟。",
		Tip: tips[4]},
	{Id: 12, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 13, Type: TaskRunning,
		Distance: 2500, Duration: 1500, Desc: "慢跑2分钟,行走2分钟,重复6次。\n距离：2.5公里, 时长：25分钟。",
		Tip: tips[5]},
	{Id: 14, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 3
	{Id: 15, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 16, Type: TaskRunning,
		Distance: 3000, Duration: 2100, Desc: "慢跑3分钟,行走2分钟,重复7次。\n距离：3公里, 时长：35分钟。",
		Tip: tips[6]},
	{Id: 17, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 18, Type: TaskRunning,
		Distance: 3000, Duration: 1500, Desc: "慢跑2分钟,行走2分钟,重复6次。\n距离：3公里, 时长：25分钟。",
		Tip: tips[7]},
	{Id: 19, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 20, Type: TaskRunning,
		Distance: 3000, Duration: 1800, Desc: "慢跑3分钟,行走2分钟,重复6次。\n距离：3公里, 时长：30分钟。",
		Tip: tips[8]},
	{Id: 21, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 4
	{Id: 22, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 23, Type: TaskRunning,
		Distance: 3500, Duration: 2100, Desc: "慢跑3分钟,行走2分钟,重复6次。\n距离：3.5公里, 时长：35分钟。",
		Tip: tips[9]},
	{Id: 24, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 25, Type: TaskRunning,
		Distance: 3500, Duration: 1800, Desc: "慢跑2分钟,行走2分钟,重复5次。\n距离：3.5公里, 时长：30分钟。",
		Tip: tips[10]},
	{Id: 26, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 27, Type: TaskRunning,
		Distance: 3500, Duration: 2400, Desc: "慢跑3分钟,行走3分钟,重复6次。\n距离：3.5公里, 时长：40分钟。",
		Tip: tips[11]},
	{Id: 28, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 5
	{Id: 29, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 30, Type: TaskRunning,
		Distance: 4000, Duration: 2400, Desc: "慢跑3分钟,行走1分钟,重复9次。\n距离：4公里, 时长：40分钟。",
		Tip: tips[12]},
	{Id: 31, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 32, Type: TaskRunning,
		Distance: 4000, Duration: 1500, Desc: "慢跑2分钟,行走1分钟,重复8次。\n距离：4公里, 时长：25分钟。",
		Tip: tips[13]},
	{Id: 33, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 34, Type: TaskRunning,
		Distance: 3500, Duration: 1800, Desc: "慢跑3分钟,行走1分钟,重复8次。\n距离：4公里, 时长：30分钟。",
		Tip: tips[14]},
	{Id: 35, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 6
	{Id: 36, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 37, Type: TaskRunning,
		Distance: 4500, Duration: 2700, Desc: "慢跑5分钟,行走1分钟,重复7次。\n距离：4.5公里, 时长：45分钟。",
		Tip: tips[15]},
	{Id: 38, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 39, Type: TaskRunning,
		Distance: 4500, Duration: 2100, Desc: "慢跑3分钟,行走1分钟,重复7次。\n距离：4.5公里, 时长：35分钟。"},
	{Id: 40, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 41, Type: TaskRunning,
		Distance: 4500, Duration: 2700, Desc: "慢跑3分钟,行走1分钟,重复10次。\n距离：4.5公里, 时长：45分钟。"},
	{Id: 42, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 7
	{Id: 43, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 44, Type: TaskRunning,
		Distance: 5500, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：5.5公里, 时长：45分钟。"},
	{Id: 45, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 46, Type: TaskRunning,
		Distance: 5000, Duration: 2700, Desc: "慢跑4分钟,行走1分钟,重复6次。\n距离：5公里, 时长：45分钟。"},
	{Id: 47, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 48, Type: TaskRunning,
		Distance: 5000, Duration: 2700, Desc: "慢跑5分钟,行走1分钟,重复7次。\n距离：5公里, 时长：45分钟。"},
	{Id: 49, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 8
	{Id: 50, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 51, Type: TaskRunning,
		Distance: 5500, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：5.5公里, 时长：45分钟。"},
	{Id: 52, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 53, Type: TaskRunning,
		Distance: 5500, Duration: 2400, Desc: "慢跑3分钟,行走1分钟,重复7次。\n距离：5.5公里, 时长：40分钟。"},
	{Id: 54, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 55, Type: TaskRunning,
		Distance: 5500, Duration: 2400, Desc: "慢跑5分钟,行走1分钟,重复6次。\n距离：5.5公里, 时长：40分钟。"},
	{Id: 56, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 9
	{Id: 57, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 58, Type: TaskRunning,
		Distance: 7000, Duration: 3600, Desc: "慢跑10分钟步行1分钟,慢跑15分钟步行1分钟,慢跑20分钟步行1分钟,慢跑10分钟。\n距离：7公里, 时长：60分钟。"},
	{Id: 59, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 60, Type: TaskRunning,
		Distance: 5000, Duration: 2280, Desc: "慢跑5分钟,行走1分钟,重复6次。\n距离：5公里, 时长：38分钟。"},
	{Id: 61, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 62, Type: TaskRunning,
		Distance: 5000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：5公里, 时长：45分钟。"},
	{Id: 63, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 10
	{Id: 64, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 65, Type: TaskRunning,
		Distance: 8000, Duration: 4200, Desc: "慢跑10分钟步行1分钟,慢跑20分钟步行1分钟,慢跑30分钟。\n距离：8公里, 时长：70分钟。"},
	{Id: 66, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 67, Type: TaskRunning,
		Distance: 6000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：6公里, 时长：45分钟。"},
	{Id: 68, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 69, Type: TaskRunning,
		Distance: 6000, Duration: 3000, Desc: "慢跑20分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。\n距离：6公里, 时长：50分钟。"},
	{Id: 70, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 11
	{Id: 71, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 72, Type: TaskRunning,
		Distance: 9000, Duration: 3900, Desc: "慢跑40分钟步行1分钟,慢跑20分钟。\n距离：9公里, 时长：65分钟。"},
	{Id: 73, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 74, Type: TaskRunning,
		Distance: 6000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。\n距离：6公里, 时长：45分钟。"},
	{Id: 75, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 76, Type: TaskRunning,
		Distance: 8000, Duration: 3000, Desc: "慢跑20分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。\n距离：8公里, 时长：50分钟。"},
	{Id: 77, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 12
	{Id: 78, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 79, Type: TaskRunning,
		Distance: 10000, Duration: 2400, Desc: "慢跑50分钟。\n距离：10公里, 时长：40分钟。"},
	{Id: 80, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 81, Type: TaskRunning,
		Distance: 5000, Duration: 1980, Desc: "慢跑10分钟,行走1分钟,重复3次。\n距离：5公里, 时长：33分钟。"},
	{Id: 82, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 83, Type: TaskRunning,
		Distance: 8000, Duration: 2700, Desc: "慢跑15分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。\n距离：8公里, 时长：45分钟。"},
	{Id: 84, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},

	// week 13
	{Id: 85, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 86, Type: TaskRunning,
		Distance: 10000, Duration: 2400, Desc: "慢跑40分钟。\n距离：10公里, 时长：40分钟。"},
	{Id: 87, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
	{Id: 88, Type: TaskRunning,
		Distance: 5000, Duration: 1980, Desc: "慢跑10分钟,行走1分钟,重复3次。\n距离：5公里, 时长：33分钟。"},
	{Id: 89, Type: TaskPost, Desc: "发表一篇运动日志,分享一下你的运动心情吧"},
	{Id: 90, Type: TaskRunning,
		Distance: 10000, Duration: 2400, Desc: "慢跑40分钟。\n距离：10公里, 时长：40分钟。"},
	{Id: 91, Type: TaskGame, Desc: "让我们开始魔法冒险，提升魔法值吧"},
}

var tips = []string{
	`初看这个任务，很多人都觉目标太低，没有多大意义。事实上，我们现在开始需要的是重新思考，再次出发，练习正确的跑步方式，让自己最终能够轻松的完成跑完全程马拉松的目标。
跑步需要倾听自己内在的声音，而不是更多地的受到外部目标的影响。你的身体是你的老师和学生，了解身体，这称作“身体感知”。`,
	`现在思考一下正确跑步的方式，我们后面会持续强化以下几个方面：
1. 正确的姿势
2. 放松的四肢
3. 自由活动的关节
4. 核心肌肉的使用
5. 专注的意念
6. 良好的呼吸技术。`,
	`应对一种力量的最佳方式是顺其自然。我们跑步时，身体受到两种力量的影响：向下的重力和向前跑动时道路传来的反向力。
跑步时，让身体保持轻微前倾，找到平衡，你就学会了和重力合作，当重力向前拉你走时，你跟着走，让小腿尽可能休息。
让我们的脚在身体重心略微靠后一点位置落地，让重心位于落脚点之前，以柔和的全脚掌着地代替脚跟着地，可以化解道路的反向力，而且不会降低速度，也不会对身体造成冲击。
这个技术非常关键，需要认真思考，如果有疑问，可以到“跑步圣经”圈子里发帖，和有经验的跑步爱好者一起讨论。`,
	`跑步的主要原则：
1. 绵里藏针：让身体成一条直线，四肢要像棉花一样柔软，毫不紧张，学习如何用身体的中心部位带动跑步。想象在你的头部和尾椎骨之间有一条中心线，了解它，感知它，在你的身体找到这个中心，感觉找到这个中心。
2. 循序渐进：在开始阶段一定要慢慢起步，一点点加速让我们身体逐步适应。这是我们做任何事情都必须遵守的原则。
3. 动作中的平衡：在跑步时，我们需要让身体所有部位在一个整体中运行，所有部位都协调工作，以平衡的方式运动。寻找平衡的关键是感知自己的中心线在哪里。`,
	`身体感知三个秘诀：
1. 仔细倾听：每当我们对跑姿调整时，都要倾听任何细微的变化：身体是怎样运动的?感觉如何?各部分是否协调？
2. 信息评估：问问自己，身体是否正按着我们希望的方式遇到弄，尽可能辨别我们所作的每次调整是否有效。
3. 逐步调整：进行微调从来都是最好的策略，任何突然的改变都可能对身体造成伤害。`,
	`呼气，挖掘你的“气”：呼吸在跑步中起着提供氧气的关键作用，我们要改变日常活动的呼吸方式，要采用腹式呼吸，尽可能吐完空气，放松腹部，保持和步频匹配，一般是三步呼出两步吸入。
在进行腹式呼吸时用鼻子呼吸效果会更好，尽可能不用嘴呼吸。`,
	`跑步时，尽量使身体的各部分向同一方向运动。当身体对直时，我们就有一个驱动身体的中心线，这条线笔直和强壮时，它就成了支撑你身体的“针”，你的双臂和双腿就变成了“棉”。
当你姿态对直时，身体的重量就是由你的结构而不是肌肉来支撑了，能量或“气”会毫无障碍地流过你的身体，让你不必付出太大努力的情况下跑的更快，更远。`,
	`对直腿和脚：
1. 双脚分开、对直，脚尖指向前方，平行站立。要使双脚对直，不仅仅是让脚尖指向前方，而是要将双腿向内侧扭转，直到腿的正面对直脚尖指向前方，不再是“八”字脚。
2. 平衡双腿的压力，确保你身体的重量是平均的分配在左右脚之间的，然后在双腿间进行平衡，再在双脚间平衡，感觉两侧压力是平均的。
3. 在跑动时，保持脚尖指向前方，假想在地上有一条直线，而你的双脚就平行地在并列在直线两边。`,
	`延长脊柱来对直上半身：向上提后颈部，延长整个脊柱，让身体的中心线又长又直，同时放松膝盖，扩张胸部，以便更充分的呼吸。`,
	`保持骨盘水平，利用核心力量，可以强化小腹肌肉，专注自己中心线部分，感受到身体和能量的集中：
1. 在跑步时，保持一条笔直的身体中心线。
2. 在运动中固定住骨盘。
3. 在骨盘和双腿间建立起更加强有力的连接，使整个下半身的运动成为一个整体。`,
	`创建身体立柱：让我们的髋部和肩膀对齐，上半身在髋部的正上方，保持平衡。肩膀，髋部和脚部将形成身体的立柱。当身体立柱对直时，身体的重量将有我们的结构支撑，肌肉就可以做它该做的事。`,
	`单腿站姿：要通过单独站立训练你的核心肌肉，让我们在跑步的过程中保持笔直的身体姿势，因为在跑步时，我们其实是在进行一系列的单腿站立活动。`,
	`学会如何前倾：身体前倾可以让重力成为我们的好帮手，身体一旦前倾，重心就会落到身体着地点的前面，重力会把我们向下拉使你向前。
我们的脚就会像弹跳器发挥作用，让我们前进，要一直在这种状态下保持细微的平衡。`,
	`全脚掌着地：不要一边跑一边制动，不单单是脚跟着地，而是你的整个脚底着地，从前至后，从一侧到另一侧的脚底压力是平均，这样使你的双腿得以充分放松而不用负责推动身体。
你的腿只用于在步幅之间短暂地支撑一下，支撑以后，腿部会沿着道路的反向力的方向朝后方摆动。`,
	`摆臂：
1. 屈肘90度：以轻松的方式自肩部摆动，不然胳膊一张一合，保持90度。
2. 向后摆臂，而不是向前。想象着你是用肘击后面的人而不是用拳打前面的，与身体前倾保持相对平衡。
3. 不要让手越过身体中心线：想象双手之间抱着一个篮球，别让双手之间的距离小于这个篮球的宽度。
4. 放松双手：保持手指向内弯曲，拇指在上的动作。`,
	`正常摆臂跑.跑步是锻炼骨骼的好方法，所以你有必要补充充足的钙质——每天1000毫克。如果你在50岁以上，则每天需要1500毫克。低脂牛奶，低脂酸奶和深绿色叶片蔬菜都是钙质的重要来源`,
}

const (
	PostDesc = "发表文章，分享你的运动经验吧"
	GameDesc = "开始冒险，提升自己的魔法力"
)

var NewTasks = []Task{
	{Id: 1, Type: TaskRunning, Index: 0,
		Distance: 2000, Duration: 1500, Desc: "慢跑1分钟,行走2分钟,重复8次。", Tip: tips[0]},
	{Id: 2, Type: TaskPost, Index: 0, Desc: PostDesc, Tip: tips[0]},
	{Id: 3, Type: TaskGame, Index: 0, Desc: GameDesc, Tip: tips[0]},
	{Id: 4, Type: TaskRunning, Index: 1,
		Distance: 2000, Duration: 1200, Desc: "慢跑1分钟,行走2分钟,重复6次。", Tip: tips[1]},
	{Id: 5, Type: TaskPost, Index: 1, Desc: PostDesc, Tip: tips[1]},
	{Id: 6, Type: TaskGame, Index: 1, Desc: GameDesc, Tip: tips[1]},
	{Id: 7, Type: TaskRunning, Index: 2,
		Distance: 2000, Duration: 1320, Desc: "慢跑1分钟,行走2分钟,重复7次。", Tip: tips[2]},
	{Id: 8, Type: TaskPost, Index: 2, Desc: PostDesc, Tip: tips[2]},
	{Id: 9, Type: TaskGame, Index: 2, Desc: GameDesc, Tip: tips[2]},
	{Id: 10, Type: TaskRunning, Index: 3,
		Distance: 2500, Duration: 1800, Desc: "慢跑2分钟,行走2分钟,重复7次。", Tip: tips[3]},
	{Id: 11, Type: TaskPost, Index: 3, Desc: PostDesc, Tip: tips[3]},
	{Id: 12, Type: TaskGame, Index: 3, Desc: GameDesc, Tip: tips[3]},
	{Id: 13, Type: TaskRunning, Index: 4,
		Distance: 2500, Duration: 1500, Desc: "慢跑1分钟,行走2分钟,重复7次。", Tip: tips[4]},
	{Id: 14, Type: TaskPost, Index: 4, Desc: PostDesc, Tip: tips[4]},
	{Id: 15, Type: TaskGame, Index: 4, Desc: GameDesc, Tip: tips[4]},
	{Id: 16, Type: TaskRunning, Index: 5,
		Distance: 2500, Duration: 1500, Desc: "慢跑2分钟,行走2分钟,重复6次。", Tip: tips[5]},
	{Id: 17, Type: TaskPost, Index: 5, Desc: PostDesc, Tip: tips[5]},
	{Id: 18, Type: TaskGame, Index: 5, Desc: GameDesc, Tip: tips[5]},
	{Id: 19, Type: TaskRunning, Index: 6,
		Distance: 3000, Duration: 2100, Desc: "慢跑3分钟,行走2分钟,重复7次。", Tip: tips[6]},
	{Id: 20, Type: TaskPost, Index: 6, Desc: PostDesc, Tip: tips[6]},
	{Id: 21, Type: TaskGame, Index: 6, Desc: GameDesc, Tip: tips[6]},
	{Id: 22, Type: TaskRunning, Index: 7,
		Distance: 3000, Duration: 1500, Desc: "慢跑2分钟,行走2分钟,重复6次。", Tip: tips[7]},
	{Id: 23, Type: TaskPost, Index: 7, Desc: PostDesc, Tip: tips[7]},
	{Id: 24, Type: TaskGame, Index: 7, Desc: GameDesc, Tip: tips[7]},
	{Id: 25, Type: TaskRunning, Index: 8,
		Distance: 3000, Duration: 1800, Desc: "慢跑3分钟,行走2分钟,重复6次。", Tip: tips[8]},
	{Id: 26, Type: TaskPost, Index: 8, Desc: PostDesc, Tip: tips[8]},
	{Id: 27, Type: TaskGame, Index: 8, Desc: GameDesc, Tip: tips[8]},
	{Id: 28, Type: TaskRunning, Index: 9,
		Distance: 3500, Duration: 2100, Desc: "慢跑3分钟,行走2分钟,重复6次。", Tip: tips[9]},
	{Id: 29, Type: TaskPost, Index: 9, Desc: PostDesc, Tip: tips[9]},
	{Id: 30, Type: TaskGame, Index: 9, Desc: GameDesc, Tip: tips[9]},
	{Id: 31, Type: TaskRunning, Index: 10,
		Distance: 3500, Duration: 1800, Desc: "慢跑2分钟,行走2分钟,重复5次。", Tip: tips[10]},
	{Id: 32, Type: TaskPost, Index: 10, Desc: PostDesc, Tip: tips[10]},
	{Id: 33, Type: TaskGame, Index: 10, Desc: GameDesc, Tip: tips[10]},
	{Id: 34, Type: TaskRunning, Index: 11,
		Distance: 3500, Duration: 2400, Desc: "慢跑3分钟,行走3分钟,重复6次。", Tip: tips[11]},
	{Id: 35, Type: TaskPost, Index: 11, Desc: PostDesc, Tip: tips[11]},
	{Id: 36, Type: TaskGame, Index: 11, Desc: GameDesc, Tip: tips[11]},
	{Id: 37, Type: TaskRunning, Index: 12,
		Distance: 4000, Duration: 2400, Desc: "慢跑3分钟,行走1分钟,重复9次。", Tip: tips[12]},
	{Id: 38, Type: TaskPost, Index: 12, Desc: PostDesc, Tip: tips[12]},
	{Id: 39, Type: TaskGame, Index: 12, Desc: GameDesc, Tip: tips[12]},
	{Id: 40, Type: TaskRunning, Index: 13,
		Distance: 4000, Duration: 1500, Desc: "慢跑2分钟,行走1分钟,重复8次。", Tip: tips[13]},
	{Id: 41, Type: TaskPost, Index: 13, Desc: PostDesc, Tip: tips[13]},
	{Id: 42, Type: TaskGame, Index: 13, Desc: GameDesc, Tip: tips[13]},
	{Id: 43, Type: TaskRunning, Index: 14,
		Distance: 3500, Duration: 1800, Desc: "慢跑3分钟,行走1分钟,重复8次。", Tip: tips[14]},
	{Id: 44, Type: TaskPost, Index: 14, Desc: PostDesc, Tip: tips[14]},
	{Id: 45, Type: TaskGame, Index: 14, Desc: GameDesc, Tip: tips[14]},
	{Id: 46, Type: TaskRunning, Index: 15,
		Distance: 4500, Duration: 2700, Desc: "慢跑5分钟,行走1分钟,重复7次。", Tip: tips[15]},
	{Id: 47, Type: TaskPost, Index: 15, Desc: PostDesc, Tip: tips[15]},
	{Id: 48, Type: TaskGame, Index: 15, Desc: GameDesc, Tip: tips[15]},
	{Id: 49, Type: TaskRunning, Index: 16,
		Distance: 4500, Duration: 2100, Desc: "慢跑3分钟,行走1分钟,重复7次。", Tip: tips[0]},
	{Id: 50, Type: TaskPost, Index: 16, Desc: PostDesc, Tip: tips[0]},
	{Id: 51, Type: TaskGame, Index: 16, Desc: GameDesc, Tip: tips[0]},
	{Id: 52, Type: TaskRunning, Index: 17,
		Distance: 4500, Duration: 2700, Desc: "慢跑3分钟,行走1分钟,重复10次。", Tip: tips[1]},
	{Id: 53, Type: TaskPost, Index: 17, Desc: PostDesc, Tip: tips[1]},
	{Id: 54, Type: TaskGame, Index: 17, Desc: GameDesc, Tip: tips[1]},
	{Id: 55, Type: TaskRunning, Index: 18,
		Distance: 5500, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。", Tip: tips[2]},
	{Id: 56, Type: TaskPost, Index: 18, Desc: PostDesc, Tip: tips[2]},
	{Id: 57, Type: TaskGame, Index: 18, Desc: GameDesc, Tip: tips[2]},
	{Id: 58, Type: TaskRunning, Index: 19,
		Distance: 5000, Duration: 2700, Desc: "慢跑4分钟,行走1分钟,重复6次。", Tip: tips[3]},
	{Id: 59, Type: TaskPost, Index: 19, Desc: PostDesc, Tip: tips[3]},
	{Id: 60, Type: TaskGame, Index: 19, Desc: GameDesc, Tip: tips[3]},
	{Id: 61, Type: TaskRunning, Index: 20,
		Distance: 5000, Duration: 2700, Desc: "慢跑5分钟,行走1分钟,重复7次。", Tip: tips[4]},
	{Id: 62, Type: TaskPost, Index: 20, Desc: PostDesc, Tip: tips[4]},
	{Id: 63, Type: TaskGame, Index: 20, Desc: GameDesc, Tip: tips[4]},
	{Id: 64, Type: TaskRunning, Index: 21,
		Distance: 5500, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。", Tip: tips[5]},
	{Id: 65, Type: TaskPost, Desc: PostDesc, Tip: tips[5]},
	{Id: 66, Type: TaskGame, Desc: GameDesc, Tip: tips[5]},
	{Id: 67, Type: TaskRunning, Index: 22,
		Distance: 5500, Duration: 2400, Desc: "慢跑3分钟,行走1分钟,重复7次。", Tip: tips[6]},
	{Id: 68, Type: TaskPost, Desc: PostDesc, Tip: tips[6]},
	{Id: 69, Type: TaskGame, Desc: GameDesc, Tip: tips[6]},
	{Id: 70, Type: TaskRunning, Index: 23,
		Distance: 5500, Duration: 2400, Desc: "慢跑5分钟,行走1分钟,重复6次。", Tip: tips[7]},
	{Id: 71, Type: TaskPost, Desc: PostDesc, Tip: tips[7]},
	{Id: 72, Type: TaskGame, Desc: GameDesc, Tip: tips[7]},
	{Id: 73, Type: TaskRunning, Index: 24,
		Distance: 7000, Duration: 3600, Desc: "慢跑10分钟步行1分钟,慢跑15分钟步行1分钟,慢跑20分钟步行1分钟,慢跑10分钟。", Tip: tips[8]},
	{Id: 74, Type: TaskPost, Desc: PostDesc, Tip: tips[8]},
	{Id: 75, Type: TaskGame, Desc: GameDesc, Tip: tips[8]},
	{Id: 76, Type: TaskRunning, Index: 25,
		Distance: 5000, Duration: 2280, Desc: "慢跑5分钟,行走1分钟,重复6次。", Tip: tips[9]},
	{Id: 77, Type: TaskPost, Desc: PostDesc, Tip: tips[9]},
	{Id: 78, Type: TaskGame, Desc: GameDesc, Tip: tips[9]},
	{Id: 79, Type: TaskRunning, Index: 26,
		Distance: 5000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。", Tip: tips[10]},
	{Id: 80, Type: TaskPost, Desc: PostDesc, Tip: tips[10]},
	{Id: 81, Type: TaskGame, Desc: GameDesc, Tip: tips[10]},
	{Id: 82, Type: TaskRunning, Index: 27,
		Distance: 8000, Duration: 4200, Desc: "慢跑10分钟步行1分钟,慢跑20分钟步行1分钟,慢跑30分钟。", Tip: tips[11]},
	{Id: 83, Type: TaskPost, Desc: PostDesc, Tip: tips[11]},
	{Id: 84, Type: TaskGame, Desc: GameDesc, Tip: tips[11]},
	{Id: 85, Type: TaskRunning, Index: 28,
		Distance: 6000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。", Tip: tips[12]},
	{Id: 86, Type: TaskPost, Desc: PostDesc, Tip: tips[12]},
	{Id: 87, Type: TaskGame, Desc: GameDesc, Tip: tips[12]},
	{Id: 88, Type: TaskRunning, Index: 29,
		Distance: 6000, Duration: 3000, Desc: "慢跑20分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。", Tip: tips[13]},
	{Id: 89, Type: TaskPost, Desc: PostDesc, Tip: tips[13]},
	{Id: 90, Type: TaskGame, Desc: GameDesc, Tip: tips[13]},
	{Id: 91, Type: TaskRunning, Index: 30,
		Distance: 9000, Duration: 3900, Desc: "慢跑40分钟步行1分钟,慢跑20分钟。", Tip: tips[14]},
	{Id: 92, Type: TaskPost, Desc: PostDesc, Tip: tips[14]},
	{Id: 93, Type: TaskGame, Desc: GameDesc, Tip: tips[14]},
	{Id: 94, Type: TaskRunning, Index: 31,
		Distance: 6000, Duration: 2700, Desc: "慢跑10分钟,行走1分钟,重复4次。", Tip: tips[15]},
	{Id: 95, Type: TaskPost, Desc: PostDesc, Tip: tips[15]},
	{Id: 96, Type: TaskGame, Desc: GameDesc, Tip: tips[15]},
	{Id: 97, Type: TaskRunning, Index: 32,
		Distance: 8000, Duration: 3000, Desc: "慢跑20分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。"},
	{Id: 98, Type: TaskPost, Desc: PostDesc},
	{Id: 99, Type: TaskGame, Desc: GameDesc},
	{Id: 100, Type: TaskRunning, Index: 33,
		Distance: 10000, Duration: 2400, Desc: "慢跑50分钟。"},
	{Id: 101, Type: TaskPost, Desc: PostDesc},
	{Id: 102, Type: TaskGame, Desc: GameDesc},
	{Id: 103, Type: TaskRunning, Index: 34,
		Distance: 5000, Duration: 1980, Desc: "慢跑10分钟,行走1分钟,重复3次。"},
	{Id: 104, Type: TaskPost, Desc: PostDesc},
	{Id: 105, Type: TaskGame, Desc: GameDesc},
	{Id: 106, Type: TaskRunning, Index: 35,
		Distance: 8000, Duration: 2700, Desc: "慢跑15分钟步行1分钟,慢跑15分钟步行1分钟,慢跑10分钟。"},
	{Id: 107, Type: TaskPost, Desc: PostDesc},
	{Id: 108, Type: TaskGame, Desc: GameDesc},
	{Id: 109, Type: TaskRunning, Index: 36,
		Distance: 10000, Duration: 2400, Desc: "慢跑40分钟。"},
	{Id: 110, Type: TaskPost, Desc: PostDesc},
	{Id: 111, Type: TaskGame, Desc: GameDesc},
	{Id: 112, Type: TaskRunning, Index: 37,
		Distance: 5000, Duration: 1980, Desc: "慢跑10分钟,行走1分钟,重复3次。"},
	{Id: 113, Type: TaskPost, Desc: PostDesc},
	{Id: 114, Type: TaskGame, Desc: GameDesc},
	{Id: 115, Type: TaskRunning, Index: 38,
		Distance: 10000, Duration: 2400, Desc: "慢跑40分钟。"},
	{Id: 116, Type: TaskPost, Desc: PostDesc},
	{Id: 117, Type: TaskGame, Desc: GameDesc},
}
