
$("#login").click =>
  User.login "admin", "123456"
    , (resp) -> console.log resp

$("#user_info").click =>
  new User("2467309932").getInfo '2f614864-7262-47bd-55c1-50bf2a7de97c'
    , (resp) -> console.log resp

$("#user_list").click =>
  User.list '2f614864-7262-47bd-55c1-50bf2a7de97c', '',
    (resp) -> console.log resp

$("#user_search").click =>
  User.search '2f614864-7262-47bd-55c1-50bf2a7de97c', "yuan", '', '', '', '',
    (resp, index, pages, total)->
      console.log(resp)
      console.log(index, pages, total)

$("#user_ban").click =>
  new User("yuan.li3@tcl.com").ban('2f614864-7262-47bd-55c1-50bf2a7de97c', 0,
    (resp) ->
      console.log(resp)
  )

$("#article_info").click =>
  new Article().getInfo '541941d0763d947b5c000005', '2f614864-7262-47bd-55c1-50bf2a7de97c',
    (resp) -> console.log resp

$("#article_list").click =>
  Article.list '2f614864-7262-47bd-55c1-50bf2a7de97c',
    (resp, index, pages, total) ->
      console.log resp
      console.log(index, pages, total)

$("#article_search").click =>
  Article.search '2f614864-7262-47bd-55c1-50bf2a7de97c', "hello",
    (resp, index, pages, total)->
      console.log(resp)
      console.log(index, pages, total)