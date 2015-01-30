class Task
  @list: (token, userid, week, callback)->
    Util._get('../admin/task/list',
      {userid:userid, week:week,access_token:token},
      (resp)->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          tasks = []
          for task in resp.tasks
            tasks.push(new Task()._update(task))
          callback(tasks)
    )

  @auth: (token, userid, task_id, pass, reason)->
    Util._post('../admin/task/auth',
      {userid:userid, task_id: task_id, pass:pass, reason:reason, access_token: token},
      (resp)->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(resp)
    )

  _update: (data)->
    @task_id = data.task_id
    @type = data.type
    @desc = data.desc
    @images = data.images
    @status = data.status
    @reason = data.reason
    @