package controllers

import (
	//"encoding/json"
	//"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/jinzhu/now"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"math/rand"
	"net/http"
	"time"
)

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

func init() {
	now.FirstDayMonday = true
}

func BindTaskApi(m *martini.ClassicMartini) {
	m.Get("/1/tasks/getList",
		binding.Form(getTasksForm{}),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		getTasksHandler)
	m.Get("/1/tasks/getInfo",
		binding.Form(getTaskInfoForm{}),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		getTaskInfoHandler)
	m.Post("/1/tasks/execute",
		binding.Json(completeTaskForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		completeTaskHandler)
}

type getTasksForm struct {
	parameter
}

func getTasksHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	tasklist := user.Tasks

	week := len(tasklist.Completed) / 7
	if week > 0 && len(tasklist.Completed)%7 == 0 &&
		tasklist.Last.After(now.BeginningOfWeek()) {
		week -= 1
	}
	list := make([]models.Task, 7)
	if week*7+7 <= len(models.Tasks)+1 {
		copy(list, models.Tasks[week*7:week*7+7])
	}
	for i, _ := range list {
		list[i].Status = tasklist.TaskStatus(list[i].Id)
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
	Tid int `form:"task_id" binding:"required"`
	parameter
}

func getTaskInfoHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(getTaskInfoForm)
	tasklist := user.Tasks
	var task models.Task
	if form.Tid > 0 && form.Tid <= len(models.Tasks) {
		task = models.Tasks[form.Tid-1]
	}

	task.Status = tasklist.TaskStatus(task.Id)
	proof := tasklist.GetProof(task.Id)
	task.Pics = proof.Pics
	task.Result = proof.Result

	writeResponse(request.RequestURI, resp, map[string]interface{}{"task_info": task}, nil)
}

type completeTaskForm struct {
	Tid    int      `json:"task_id" binding:"required"`
	Proofs []string `json:"task_pics"`
	parameter
}

func completeTaskHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(completeTaskForm)

	if form.Tid < 1 || form.Tid > len(models.Tasks) {
		writeResponse(request.RequestURI, resp, nil, nil)
		return
	}

	//u := &models.User{Id: user.Id}
	err := user.AddTask(models.Tasks[form.Tid-1].Type, form.Tid, form.Proofs)

	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, err)
}
