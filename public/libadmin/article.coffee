class Article
  @list: (token, callback, sort='', pageIndex = 0, pageCount = 50) ->
    Util._get(Util.host + "/admin/article/list", {sort:sort, page_index: pageIndex, page_count: pageCount, access_token: token},
      (resp) =>
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          articles = []
          for info in resp.articles
            articles.push(new Article()._update(info))
          callback(articles, resp.page_index, resp.page_total, resp.total_number)
    )

  @lists: (token, sort='', pageIndex = 0, pageCount = 50) ->
    $http.get(Util.host + '/admin/article/list',
      {params: {sort:sort, page_index: pageIndex, page_count: pageCount, access_token: token}}).then(
        (response) =>
          if typeof response.data is 'object'
            response.data
          else
            $q.reject(response.data)
        , (response) ->
            $q.reject(response.data)
      )

  @timeline: (userid, token, callback, pageIndex = 0, pageCount = 50) ->
    Util._get(Util.host + "/admin/article/timeline", {userid: userid, page_index: pageIndex, page_count: pageCount, access_token: token},
      (resp) =>
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          articles = []
          for info in resp.articles
            articles.push(new Article()._update(info))
          callback(articles, resp.page_index, resp.page_total, resp.total_number)
    )

  @search: (token, keyword, callback, sort = '', pageIndex = 0, pageCount = 50) ->
    Util._get(Util.host + "/admin/article/search",
      {keyword: keyword, sort: sort, page_index: pageIndex, page_count: pageCount, access_token: token},
      (resp) =>
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          articles = []
          for info in resp.articles
            articles.push(new Article()._update(info))
          callback(articles, resp.page_index, resp.page_total, resp.total_number)
    )

  getInfo: (article_id, token, callback) ->
    Util._get(Util.host + "/admin/article/info",
      {article_id: article_id, access_token: token},
      (resp) =>
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(@_update(resp))
    )

  post: (token, callback)->
    Util._post(Util.host + "/admin/article/post",
      {article_id: @article_id, author:@author, contents:@contents, tags: [], access_token: token},
      (resp) =>
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(resp)
    )

  delete: (token, callback)->
    Util._post(Util.host + "/admin/article/delete",
      {article_id: @article_id, access_token:token},
      (resp)->
        if Error._hasError(resp)
          callback(new Error(resp.error_id, resp.error_desc))
        else
          callback(resp)
    )

  _update: (data) ->
    @article_id = data.article_id
    @parent = data.parent

    if data.author? && data.author.userid?
      @author = new User(data.author.userid)._update(data.author)
    @cover_image = data.cover_image
    @cover_text = data.cover_text
    if data.cover_text is ''
      @cover_text = '无标题文章'
    @time = Util._formatDate(new Date(data.time * 1000))
    @thumbs_count = data.thumbs_count
    @comments_count = data.comments_count
    @rewards_value = data.rewards_value / 100000000
    @rewards_users = []
    if data.rewards_users isnt null
      @rewards_users = data.rewards_users

    @tags = []
    if data.tags?
      for tag in data.tags
        switch tag
          when 'SPORT_LOG' then @tags.push(new Tag(tag, '运动日志'))
          when 'SPORT_THEORY' then @tags.push(new Tag(tag, '跑步圣经'))
          when 'EQUIP_BLOG' then @tags.push(new Tag(tag, '我爱装备'))
          when 'SPORT_LIFE' then @tags.push(new Tag(tag, '运动生活'))
          when 'PRODUCT_PROPOSAL' then @tags.push(new Tag(tag, '产品建议'))
    if @tags.length is 0
      @tags.push(new Tag('SPORT_LOG', '运动日志'))

    @contents = data.contents
    @_comments(data.comments)
    @

  _comments: (data) ->
    @comments = []
    if data isnt null
      for a in data
        art = new Article()._update(a)
        art._comments(a.comments)
        @comments.push(art)

class Tag
  constructor: (@id, @name)->