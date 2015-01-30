class Util
  @host: "http://172.24.222.54:8080"

  @_get: (url, data, callback) ->
    $.getJSON url, data, callback

  @_post: (url, data, callback) ->
    $.ajax url,
      type: "POST"
      url: url
      data: JSON.stringify(data)
      dataType: "json"
      success: callback

  @_formatDate: (date) ->
    d = [date.getFullYear(), date.getMonth() + 1, date.getDate()].join("-")
    t = [date.getHours(), date.getMinutes(), date.getSeconds()].join(":")
    [d, t].join(" ")

  @_birth2Age: (birth) ->
    new Date().getFullYear() - birth.getFullYear()

class Error
  constructor: (@error_id, @error_desc)->

  @_hasError: (data) ->
    if data.error_id? and data.error_id > 0
      yes
    else
      no
  String: ->
    "#{@error_id}: #{@error_desc}"