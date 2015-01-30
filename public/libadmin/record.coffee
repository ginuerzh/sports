class Record
  @timeline: (token, userid, type, callback, pageIndex = 0, pageCount = 50) ->
    Util._get('../admin/record/timeline',
      {userid: userid, type: type, page_index: pageIndex, page_count: pageCount, access_token: token},
      (resp)->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          records = []
          for rec in resp.records
            records.push(new Record()._update(rec))
          callback(records, resp.page_index, resp.page_total, resp.total_number)
    )

  delete: (token, userid, type, callback) ->
    Util._post('../admin/record/delete',
      {userid: userid, type: type},
      (resp)->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(resp)
    )

  _update: (data)->
    @record_id = data.record_id
    @type = data.type
    @duration = data.duration
    @distance = data.distance
    @iamges = data.iamges
    @game_name = data.game_name
    @game_score = data.game_score
    @time = Util._formatDate(new Date(data.time * 1000))
    @pub_time = Util._formatDate(new Date(data.pub_time * 1000))
    @
