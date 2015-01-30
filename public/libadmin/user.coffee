
class User
  constructor: (@userid) ->

  @login: (userid, password, callback) ->
    Util._post(Util.host + "/admin/login", {username:userid, password: password},
      (resp) ->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
         callback(resp)
    )

  @logout: (token, callback) ->
    Util._post(Util.host + "/admin/logout", {access_token: token},
      (resp) ->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(resp)
    )

  @list: (token, sort, callback, pageIndex = 0, pageCount = 50) ->
    Util._get(Util.host + "/admin/user/list", {sort: sort, page_index: pageIndex, page_count: pageCount, access_token: token},
      (resp) =>
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          users = []
          for info in resp.users
            users.push(new User(info.userid)._update(info))
          callback(users, resp.page_index, resp.page_total, resp.total_number)
    )
  #sort:  regtime(注册时间), logintime(登录时间), userid(用户名), nickname(昵称), score(总分), age(年龄), ban(状态).
  @search: (token, keyword, gender, age, ban_status, sort, callback, pageIndex = 0, pageCount = 50) ->
    Util._get(Util.host + "/admin/user/search",
      {keyword: keyword, gender: gender, age: age, ban_status: ban_status, sort: sort, \
        page_index: pageIndex, page_count: pageCount, access_token: token},
      (resp) =>
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          users = []
          for info in resp.users
            users.push(new User(info.userid)._update(info))
          callback(users, resp.page_index, resp.page_total, resp.total_number)
    )

  getInfo: (token, callback) ->
    Util._get(Util.host + "/admin/user/info",
      {userid: @userid, access_token: token},
      (resp) =>
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(@_update(resp))
    )

  ban: (token, duration, callback) ->
    Util._post(Util.host + "/admin/user/ban", {userid: @userid, duration: duration, access_token: token},
      (resp) ->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(resp)
    )


  _update: (data) ->
    @nickname = data.nickname
    @role = data.role
    @profile = ""
    if data.profile.search("http://") is 0 then @profile = data.profile

    @gender = "未知"
    if data.gender?
      if data.gender.search("f") is 0 then @gender = "女"
      if data.gender.search("m") is 0 then @gender = "男"

    @phone = data.phone
    @about = data.about
    @address = data.address

    if data.photos? then @photos = data.photos

    @hobby = data.hobby

    @birthday = ""
    @age = "未知"
    if data.birthday? and data.birthday isnt 0
      birth = new Date(data.birthday * 1000)
      @birthday = Util._formatDate(birth)
      @age = Util._birth2Age(birth)
    @reg_time = ""
    if data.reg_time? and data.reg_time > 0 then @reg_time = Util._formatDate(new Date(data.reg_time * 1000))

    if data.last_login_time? and data.last_login_time > 0
      @last_login_time = Util._formatDate(new Date(data.last_login_time * 1000))
    else
      @last_login_time = "未知"

    @height = data.height
    @weight = data.weight

    @loc_latitude = data.loc_latitude
    @loc_longitude = data.loc_longitude

    @equips = {shoes: "", hardwares: "", softwares:""}
    if data.equips?
      if data.equips.shoes isnt null and data.equips.shoes.length > 0
        @equips.shoes = data.equips.shoes.join ","
      if data.equips.hardwares isnt null and data.equips.hardwares.length > 0
        @equips.hardwares = data.equips.hardwares.join ","
      if data.equips.softwares isnt null and data.equips.softwares.length > 0
        @equips.softwares = data.equips.softwares.join ","

    @physique_value = data.physique_value
    @literature_value = data.literature_value
    @magic_value = data.magic_value
    @coin_value = data.coin_value / 100000000
    @score = data.score
    @level = data.level

    @wallet = data.wallet

    @articles_count = data.articles_count
    @follows_count = data.follows_count
    @followers_count = data.followers_count
    @friends_count = data.friends_count
    @blacklist_count = data.blacklist_count

    @ban_time = ""
    @ban_status = "normal"
    if data.ban_time?
      if data.ban_time < 0 then @ban_status = "ban"
      if data.ban_time > 0
        @ban_time = Util._formatDate(new Date(data.ban_time * 1000))
        @ban_status = "lock"
    @





