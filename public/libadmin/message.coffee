class Message
  @timeline: (token, callback, from, to = '', pageIndex = 0, pageCount = 50)->
    Util._get('../admin/chat/timeline',
      {from: from, to: to, page_index: pageIndex, page_count: pageCount, access_token: token},
      (resp)->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          messages = []
          for msg in resp.messages
            messages.push(new Message()._update(msg))
          callback(messages, resp.page_index, resp.page_total, resp.total_number)
    )

  send: (token, callback)->
    Util._post('../admin/chat/send',
      {from:@from, to:@to, contents: @contents, time: @time, access_token: token},
      (resp)->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(resp)

    )
  delete: (token, callback) ->
    Util._post('../admin/chat/delete',
      {message_id: @message_id, access_token: token},
      (resp)->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(resp)
    )

  _update:(data)->
    @message_id = data.message_id
    @from = data.from
    @to = data.to
    @time = Util._formatDate(new Date(data.time * 1000))
    @contents = data.contents
    @