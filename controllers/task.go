package controllers

import (
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"math/rand"
	"net/http"
	"time"
)

type Task struct {
	Id     int      `json:"task_id"`
	Type   string   `json:"task_type"`
	Desc   string   `json:"task_desc"`
	Status string   `json:"task_status"`
	Pics   []string `json:"task_pics,omitempty"`
	Result string   `json:"task_result,omitempty"`
}

var tasks = []Task{
	// week 1
	{Id: 1, Type: models.TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复8次"},
	{Id: 2, Type: models.TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复6次"},
	{Id: 3, Type: models.TaskRunning, Desc: "慢跑1分钟,行走2分钟,重复7次"},
	{Id: 4, Type: models.TaskPost, Desc: "发表一篇运动日志"},
	{Id: 5, Type: models.TaskPost, Desc: "发表一篇运动日志"},
	{Id: 6, Type: models.TaskGame, Desc: "玩个游戏放松一下吧"},
	{Id: 7, Type: models.TaskGame, Desc: "玩个游戏放松一下吧"},

	//week 2
}

var tips = []string{
	"慢跑为小碎步跑.为了给你的训练增加能量，你可以在出门前的两小时吃一点水果或者巧克力，然后在出门前一小时喝适量（约240克）的运动饮料，这样既能够保证你有充足的水分，也能补充钠和钾。",
	"慢跑为小碎步跑.开始训练前可先慢走2-3分钟热身，训练结束后再慢走2-3分钟放松。不要在跑步前舒展关节，而应该在训练后或晚上看电视的时候进行舒展。",
	"正常摆臂跑.跑步过程中双臂一定要保持放松。跑步时手肘弯曲约90度，在腰间前后摆臂。手指弯曲成放松的拳头，不要让手在上身中部胡乱地摇摆。",
	"正常摆臂跑.如果天气炎热，太阳猛烈，一定要涂防晒霜，戴上太阳眼镜和鸭嘴帽，防止阳光直射脸部。如果天气特别炎热潮湿，一定要注意多行走休息。尽可能在清早或者傍晚的时候跑步。",
	"慢跑为小碎步跑.有时你可以跳过行走和跑步的训练，做一些交替运动，如骑30-40分钟单车，上健身房或者参加一些举重训练课程。跑步训练期间的间歇能让你更快地恢复精力，同时还能够锻炼到新的肌肉群。",
	"正常摆臂跑.跑步是锻炼骨骼的好方法，所以你有必要补充充足的钙质——每天1000毫克。如果你在50岁以上，则每天需要1500毫克。低脂牛奶，低脂酸奶和深绿色叶片蔬菜都是钙质的重要来源。",
	"新手跑者通常会觉得胫骨、肋骨或者膝盖酸痛，如果你在训练后能够及时进行冰敷，这些痛感很快就会消失，你还可以把豆子装进袋子冷藏后敷在膝盖上15分钟。如果疼痛还持续的话，就需要停止几天的训练。",
	"要想呼吸新鲜的空气让肺部健康的话，尽量不要到繁忙的街道或者在交通高峰时跑步。找一个车辆比较少的地方，这样废气就可以很快驱散。最好就是能够找一些绿化带或公园等。",
}

func BindTaskApi(m *martini.ClassicMartini) {
	m.Get("/1/tasks/getList", binding.Form(getTasksForm{}), ErrorHandler, getTasksHandler)
	m.Get("/1/tasks/getInfo", binding.Form(getTaskInfoForm{}), ErrorHandler, getTaskInfoHandler)
	m.Post("/1/tasks/execute", binding.Json(completeTaskForm{}), ErrorHandler, completeTaskHandler)
}

type getTasksForm struct {
	Token string `form:"access_token" binding:"required"`
}

func getTaskStatus(tid int, tasklist *models.TaskList) (status string) {
	for i := len(tasklist.Completed) - 1; i >= 0; i-- {
		if tid == tasklist.Completed[i] {
			return "FINISH"
		}
	}
	for _, id := range tasklist.Waited {
		if tid == id {
			return "AUTHENTICATION"
		}
	}
	for _, id := range tasklist.Uncompleted {
		if tid == id {
			return "UNFINISH"
		}
	}

	return "NORMAL"
}

func getTasksHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getTasksForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.User{Id: user.Id}
	tasklist, err := u.GetTasks()
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	week := len(tasklist.Completed) / 7
	list := make([]Task, 7)
	copy(list, tasks[week*7:week*7+7])
	for i, _ := range list {
		list[i].Status = getTaskStatus(list[i].Id, &tasklist)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	respData := map[string]interface{}{
		"week_id":   week + 1,
		"task_list": list,
		"week_desc": tips[r.Int()%len(tips)],
	}

	writeResponse(request.RequestURI, resp, respData, nil)
}

type getTaskInfoForm struct {
	Token string `form:"access_token" binding:"required"`
	Tid   int    `form:"task_id" binding:"required"`
}

func getTaskInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form getTaskInfoForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	u := &models.User{Id: user.Id}
	tasklist, err := u.GetTasks()
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}
	var task Task
	if form.Tid > 0 && form.Tid <= len(tasks) {
		task = tasks[form.Tid-1]
	}

	task.Status = getTaskStatus(task.Id, &tasklist)

	for _, proof := range tasklist.Proofs {
		if proof.Tid == task.Id {
			task.Pics = proof.Pics
			task.Result = proof.Result
			break
		}
	}

	writeResponse(request.RequestURI, resp, map[string]interface{}{"task_info": task}, nil)
}

type completeTaskForm struct {
	Token  string   `json:"access_token" binding:"required"`
	Tid    int      `json:"task_id" binding:"required"`
	Proofs []string `json:"task_pics"`
}

func completeTaskHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form completeTaskForm) {
	user := redis.OnlineUser(form.Token)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	if form.Tid < 1 || form.Tid > len(tasks) {
		writeResponse(request.RequestURI, resp, nil, nil)
		return
	}

	u := &models.User{Id: user.Id}
	err := u.AddTask(tasks[form.Tid-1].Type, form.Tid, form.Proofs)

	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, err)
}
