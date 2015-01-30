var Error, Util;

Util = (function() {
  function Util() {}

  Util._get = function(url, data, callback) {
    return $.getJSON(url, data, callback);
  };

  Util._post = function(url, data, callback) {
    return $.ajax(url, {
      type: "POST",
      url: url,
      data: JSON.stringify(data),
      dataType: "json",
      success: callback
    });
  };

  Util._formatDate = function(date) {
    var d, t;
    d = [date.getFullYear(), date.getMonth() + 1, date.getDate()].join("-");
    t = [date.getHours(), date.getMinutes(), date.getSeconds()].join(":");
    return [d, t].join(" ");
  };

  Util._birth2Age = function(birth) {
    return new Date().getFullYear() - birth.getFullYear();
  };

  return Util;

})();

Error = (function() {
  function Error(error_id, error_desc) {
    this.error_id = error_id;
    this.error_desc = error_desc;
  }

  Error._hasError = function(data) {
    if ((data.error_id != null) && data.error_id > 0) {
      return true;
    } else {
      return false;
    }
  };

  Error.prototype.String = function() {
    return "" + this.error_id + ": " + this.error_desc;
  };

  return Error;

})();

var User;

User = (function() {
  function User(userid) {
    this.userid = userid;
  }

  User.login = function(userid, password, callback) {
    return Util._post("../admin/login", {
      username: userid,
      password: password
    }, function(resp) {
      if (Error._hasError(resp)) {
        return callback(new Error(resp.error_id, resp.error_desc));
      } else {
        return callback(resp);
      }
    });
  };

  User.logout = function(token, callback) {
    return Util._post("../admin/logout", {
      access_token: token
    }, function(resp) {
      if (Error._hasError(resp)) {
        return callback(new Error(resp.error_id, resp.error_desc));
      } else {
        return callback(resp);
      }
    });
  };

  User.list = function(token, sort, callback, pageIndex, pageCount) {
    if (pageIndex == null) {
      pageIndex = 0;
    }
    if (pageCount == null) {
      pageCount = 50;
    }
    return Util._get("../admin/user/list", {
      sort: sort,
      page_index: pageIndex,
      page_count: pageCount,
      access_token: token
    }, (function(_this) {
      return function(resp) {
        var info, users, _i, _len, _ref;
        if (Error._hasError(resp)) {
          return callback(new Error(resp.error_id, resp.error_desc));
        } else {
          users = [];
          _ref = resp.users;
          for (_i = 0, _len = _ref.length; _i < _len; _i++) {
            info = _ref[_i];
            users.push(new User(info.userid)._update(info));
          }
          return callback(users, resp.page_index, resp.page_total, resp.total_number);
        }
      };
    })(this));
  };

  User.search = function(token, keyword, gender, age, ban_status, sort, callback, pageIndex, pageCount) {
    if (pageIndex == null) {
      pageIndex = 0;
    }
    if (pageCount == null) {
      pageCount = 50;
    }
    return Util._get("../admin/user/search", {
      keyword: keyword,
      gender: gender,
      age: age,
      ban_status: ban_status,
      sort: sort,
      page_index: pageIndex,
      page_count: pageCount,
      access_token: token
    }, (function(_this) {
      return function(resp) {
        var info, users, _i, _len, _ref;
        if (Error._hasError(resp)) {
          return callback(new Error(resp.error_id, resp.error_desc));
        } else {
          users = [];
          _ref = resp.users;
          for (_i = 0, _len = _ref.length; _i < _len; _i++) {
            info = _ref[_i];
            users.push(new User(info.userid)._update(info));
          }
          return callback(users, resp.page_index, resp.page_total, resp.total_number);
        }
      };
    })(this));
  };

  User.prototype.getInfo = function(token, callback) {
    return Util._get("../admin/user/info", {
      userid: this.userid,
      access_token: token
    }, (function(_this) {
      return function(resp) {
        if (Error._hasError(resp)) {
          return callback(new Error(resp.error_id, resp.error_desc));
        } else {
          return callback(_this._update(resp));
        }
      };
    })(this));
  };

  User.prototype.ban = function(token, duration, callback) {
    return Util._post("../admin/user/ban", {
      userid: this.userid,
      duration: duration,
      access_token: token
    }, function(resp) {
      if (Error._hasError(resp)) {
        return callback(new Error(resp.error_id, resp.error_desc));
      } else {
        return callback(resp);
      }
    });
  };

  User.prototype._update = function(data) {
    var birth;
    this.nickname = data.nickname;
    this.role = data.role;
    this.profile = "";
    if (data.profile.search("http://") === 0) {
      this.profile = data.profile;
    }
    this.gender = "未知";
    if (data.gender != null) {
      if (data.gender.search("f") === 0) {
        this.gender = "女";
      }
      if (data.gender.search("m") === 0) {
        this.gender = "男";
      }
    }
    this.phone = data.phone;
    this.about = data.about;
    this.address = data.address;
    if (data.photos != null) {
      this.photos = data.photos;
    }
    this.hobby = data.hobby;
    this.birthday = "";
    this.age = "未知";
    if ((data.birthday != null) && data.birthday !== 0) {
      birth = new Date(data.birthday * 1000);
      this.birthday = Util._formatDate(birth);
      this.age = Util._birth2Age(birth);
    }
    this.reg_time = "";
    if ((data.reg_time != null) && data.reg_time > 0) {
      this.reg_time = Util._formatDate(new Date(data.reg_time * 1000));
    }
    if ((data.last_login_time != null) && data.last_login_time > 0) {
      this.last_login_time = Util._formatDate(new Date(data.last_login_time * 1000));
    } else {
      this.last_login_time = "未知";
    }
    this.height = data.height;
    this.weight = data.weight;
    this.loc_latitude = data.loc_latitude;
    this.loc_longitude = data.loc_longitude;
    this.equips = {};
    if (data.equips != null) {
      this.equips = data.equips;
    }
    this.physique_value = data.physique_value;
    this.literature_value = data.literature_value;
    this.magic_value = data.magic_value;
    this.coin_value = data.coin_value / 100000000;
    this.score = data.score;
    this.level = data.level;
    this.wallet = data.wallet;
    this.articles_count = data.articles_count;
    this.follows_count = data.follows_count;
    this.followers_count = data.followers_count;
    this.friends_count = data.friends_count;
    this.blacklist_count = data.blacklist_count;
    this.ban_time = "";
    this.ban_status = "normal";
    if (data.ban_time != null) {
      if (data.ban_time < 0) {
        this.ban_status = "ban";
      }
      if (data.ban_time > 0) {
        this.ban_time = Util._formatDate(new Date(data.ban_time * 1000));
        this.ban_status = "lock";
      }
    }
    return this;
  };

  return User;

})();

var Article, Tag;

Article = (function() {
  function Article() {}

  Article.list = function(token, callback, sort, pageIndex, pageCount) {
    if (sort == null) {
      sort = '';
    }
    if (pageIndex == null) {
      pageIndex = 0;
    }
    if (pageCount == null) {
      pageCount = 50;
    }
    return Util._get("../admin/article/list", {
      sort: sort,
      page_index: pageIndex,
      page_count: pageCount,
      access_token: token
    }, (function(_this) {
      return function(resp) {
        var articles, info, _i, _len, _ref;
        if (Error._hasError(resp)) {
          return callback(new Error(resp.error_id, resp.error_desc));
        } else {
          articles = [];
          _ref = resp.articles;
          for (_i = 0, _len = _ref.length; _i < _len; _i++) {
            info = _ref[_i];
            articles.push(new Article()._update(info));
          }
          return callback(articles, resp.page_index, resp.page_total, resp.total_number);
        }
      };
    })(this));
  };

  Article.timeline = function(userid, token, callback, pageIndex, pageCount) {
    if (pageIndex == null) {
      pageIndex = 0;
    }
    if (pageCount == null) {
      pageCount = 50;
    }
    return Util._get("../admin/article/timeline", {
      userid: userid,
      page_index: pageIndex,
      page_count: pageCount,
      access_token: token
    }, (function(_this) {
      return function(resp) {
        var articles, info, _i, _len, _ref;
        if (Error._hasError(resp)) {
          return callback(new Error(resp.error_id, resp.error_desc));
        } else {
          articles = [];
          _ref = resp.articles;
          for (_i = 0, _len = _ref.length; _i < _len; _i++) {
            info = _ref[_i];
            articles.push(new Article()._update(info));
          }
          return callback(articles, resp.page_index, resp.page_total, resp.total_number);
        }
      };
    })(this));
  };

  Article.search = function(token, keyword, callback, sort, pageIndex, pageCount) {
    if (sort == null) {
      sort = '';
    }
    if (pageIndex == null) {
      pageIndex = 0;
    }
    if (pageCount == null) {
      pageCount = 50;
    }
    return Util._get("../admin/article/search", {
      keyword: keyword,
      sort: sort,
      page_index: pageIndex,
      page_count: pageCount,
      access_token: token
    }, (function(_this) {
      return function(resp) {
        var articles, info, _i, _len, _ref;
        if (Error._hasError(resp)) {
          return callback(new Error(resp.error_id, resp.error_desc));
        } else {
          articles = [];
          _ref = resp.articles;
          for (_i = 0, _len = _ref.length; _i < _len; _i++) {
            info = _ref[_i];
            articles.push(new Article()._update(info));
          }
          return callback(articles, resp.page_index, resp.page_total, resp.total_number);
        }
      };
    })(this));
  };

  Article.prototype.getInfo = function(article_id, token, callback) {
    return Util._get("../admin/article/info", {
      article_id: article_id,
      access_token: token
    }, (function(_this) {
      return function(resp) {
        if (Error._hasError(resp)) {
          return callback(new Error(resp.error_id, resp.error_desc));
        } else {
          return callback(_this._update(resp));
        }
      };
    })(this));
  };

  Article.prototype.post = function(token, callback) {
    return Util._post("../admin/article/post", {
      article_id: this.article_id,
      author: this.author,
      contents: this.contents,
      tags: this.tags,
      access_token: token
    }, (function(_this) {
      return function(resp) {
        if (Error._hasError(resp)) {
          return callback(new Error(resp.error_id, resp.error_desc));
        } else {
          return callback(resp);
        }
      };
    })(this));
  };

  Article.prototype._update = function(data) {
    var tag, _i, _len, _ref;
    this.article_id = data.article_id;
    this.parent = data.parent;
    this.author = data.author;
    this.cover_image = data.cover_image;
    this.cover_text = data.cover_text;
    if (data.cover_text === '') {
      this.cover_text = '无标题文章';
    }
    this.time = Util._formatDate(new Date(data.time * 1000));
    this.thumbs_count = data.thumbs_count;
    this.comments_count = data.comments_count;
    this.rewards_value = data.rewards_value / 100000000;
    this.rewards_users = [];
    if (data.rewards_users !== null) {
      this.rewards_users = data.rewards_users;
    }
    this.tags = [];
    if (data.tags != null) {
      _ref = data.tags;
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        tag = _ref[_i];
        switch (tag) {
          case 'SPORT_LOG':
            this.tags.push(new Tag(tag, '运动日志'));
            break;
          case 'SPORT_THEORY':
            this.tags.push(new Tag(tag, '跑步圣经'));
            break;
          case 'EQUIP_BLOG':
            this.tags.push(new Tag(tag, '我爱装备'));
            break;
          case 'SPORT_LIFE':
            this.tags.push(new Tag(tag, '运动生活'));
            break;
          case 'PRODUCT_PROPOSAL':
            this.tags.push(new Tag(tag, '产品建议'));
        }
      }
    }
    if (this.tags.length === 0) {
      this.tags.push(new Tag('SPORT_LOG', '运动日志'));
    }
    this.contents = data.contents;
    this._comments(data.comments);
    return this;
  };

  Article.prototype._comments = function(data) {
    var a, art, _i, _len, _results;
    this.comments = [];
    if (data !== null) {
      _results = [];
      for (_i = 0, _len = data.length; _i < _len; _i++) {
        a = data[_i];
        art = new Article()._update(a);
        art._comments(a.comments);
        _results.push(this.comments.push(art));
      }
      return _results;
    }
  };

  return Article;

})();

Tag = (function() {
  function Tag(id, name) {
    this.id = id;
    this.name = name;
  }

  return Tag;

})();

var app, articleObj, checkRequest, urlPath, userObj;

app = angular.module('app', ['ngRoute', 'ngCookies']);

app.constant('app', {
  version: Date.now()
});

userObj = new User();

articleObj = new Article();

checkRequest = function(reqData) {
  if (reqData === null || (reqData.error_id == null)) {
    return true;
  }
  return false;
};

urlPath = function(url) {
  return url;
};

app.config(function($routeProvider) {
  var aricledetail, ariclelist, login, userdetail, userlist;
  login = {
    templateUrl: '../html/user-login.html',
    controller: 'loginController'
  };
  userlist = {
    templateUrl: '../html/user-list.html',
    controller: 'userlistController'
  };
  userdetail = {
    templateUrl: '../html/user-details.html',
    controller: 'userdetailController'
  };
  ariclelist = {
    templateUrl: '../html/article-list.html',
    controller: 'articleListController'
  };
  aricledetail = {
    templateUrl: '../html/article-detail.html',
    controller: 'articleDetailController'
  };
  return $routeProvider.when('/', login).when('/userlist', userlist).when('/detail/:id', userdetail).when('/articledetail/:artid', aricledetail).when('/1', userlist).when('/2', ariclelist).when('/tag/:tagid', ariclelist);
});

app.run([
  'app', '$rootScope', '$cookieStore', function(app, $rootScope, $cookieStore) {
    $rootScope.isLogin = $cookieStore.get('isLogin');
    $rootScope.profile = $cookieStore.get('profile');
    $rootScope.note = {
      'successState': false,
      'errState': false
    };
    $rootScope.sel = {
      "nCount": 0,
      "nSuccess": 0
    };
    $rootScope.dialog = {};
    $rootScope.banStyleList = [
      {
        "background-color": 'rgb(0,199,0)'
      }, {
        "background-color": 'rgb(184,184,184)'
      }, {
        "background-color": 'rgb(139,139,139)'
      }
    ];
    if ($rootScope.isLogin == null) {
      $rootScope.isLogin = false;
    }
    app.checkUser = function(data) {
      if (data.isLogin != null) {
        $cookieStore.put('isLogin', data.isLogin);
        $rootScope.isLogin = data.isLogin;
      }
      if (data.userid != null) {
        $cookieStore.put('id', data.userid);
        $rootScope.id = data.userid;
      }
      if (data.access_token != null) {
        $cookieStore.put('access_token', data.access_token);
        $rootScope.access_token = data.access_token;
      }
      if (data.profile != null) {
        $cookieStore.put('profile', data.profile);
        return $rootScope.profile = data.profile;
      }
    };
    app.getCookie = function(key) {
      return $cookieStore.get(key);
    };
    app.hideNote = function() {
      $rootScope.note.successState = false;
      $rootScope.note.errState = false;
      return $rootScope.$apply();
    };
    $rootScope.showDialog = function(configData) {
      $rootScope.dialog.content = configData.content;
      if (configData.DialBtnHidden != null) {
        $rootScope.dialog.DialBtnHidden = configData.DialBtnHidden;
      } else {
        $rootScope.dialog.DialBtnHidden = 1;
      }
      if (configData.okBtn != null) {
        $rootScope.dialog.okBtn = configData.okBtn;
      } else {
        $rootScope.dialog.okBtn = null;
      }
      $('#myModalDialog').modal('show');
      return true;
    };
    $rootScope.logout = function() {
      return User.logout($rootScope.userid, function(retData) {
        var data;
        console.log(retData);
        data = {
          isLogin: false,
          userid: '',
          access_token: '',
          profile: ''
        };
        app.checkUser(data);
        return window.location.href = "#/";
      });
    };
    $rootScope.logOutApp = function() {
      var dialogInfo;
      dialogInfo = {
        content: "你确认要退出吗？",
        DialBtnHidden: 0,
        okBtn: function() {
          return $rootScope.logout();
        }
      };
      return $rootScope.showDialog(dialogInfo);
    };
    $rootScope.navBarItems = ["首页", "用户管理", "博文管理", "权限管理", "文章导入", "消息管理"];
    return app.rootScope = $rootScope;
  }
]);

var articleDetailController, articleListController, loginController, userdetailController, userlistController,
  __indexOf = [].indexOf || function(item) { for (var i = 0, l = this.length; i < l; i++) { if (i in this && this[i] === item) return i; } return -1; };

loginController = app.controller('loginController', [
  'app', '$scope', '$routeParams', '$rootScope', function(app, $scope, $routeParams, $rootScope) {
    $scope.loginAlert = false;
    $scope.onLogin = function() {
      if (($scope.username != null) && ($scope.pwd != null)) {
        $scope.loginAlert = false;
        return User.login($scope.username, $scope.pwd, function(retData) {
          if (checkRequest(retData)) {
            retData.isLogin = true;
            app.checkUser(retData);
            window.location.href = "#/userlist";
            $scope.username = "";
            return $scope.pwd = "";
          } else {
            return $scope.loginAlert = true;
          }
        });
      } else {
        return $scope.loginAlert = true;
      }
    };
    $scope.enterLogin = function() {
      var event;
      event = window.event || arguments.callee.caller["arguments"][0];
      if (event.keyCode === 13) {
        return $scope.onLogin();
      }
    };
    return $scope.checkLogin = function() {
      $rootScope.isLogin = app.getCookie("isLogin");
      if ($routeParams.index == null) {
        return $rootScope.isLogin = false;
      }
    };
  }
]);

userlistController = app.controller('userlistController', [
  'app', '$scope', '$rootScope', function(app, $scope, $rootScope) {
    var banuserFinish, pageCount, resetCheckState, searchMode, sortStr, timer, userId, userPageIndex, userSearchPageIndex, userToken;
    pageCount = 50;
    searchMode = false;
    userPageIndex = 0;
    userSearchPageIndex = 0;
    timer = void 0;
    sortStr = "-regtime";
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    userToken = app.getCookie("access_token");
    userId = app.getCookie("id");
    $scope.dropdowmItems = ["20项", "50项", "100项", "200项"];
    $scope.selectType = ["选择", "性别", "年龄", "状态"];
    $scope.selectItemList = [["选择"], ["男", "女"], ["< 20岁", "20～40岁", "> 40岁"], ["正常", "禁言", "拉黑"]];
    $scope.selectItem = $scope.selectItemList[0];
    $scope.filtStr = ["启用过滤", "取消过滤"];
    $scope.filtState = 0;
    $scope.userlist = [];
    $scope.selectindex = 1;
    $scope.typeIndex = 0;
    $scope.filtItemIndex = 0;
    $scope.currentPage = 0;
    $scope.arrPage = [1];
    $scope.checkAllBool = false;
    $scope.searchData = {
      "data": ""
    };
    $scope.selectedList = [];
    $scope.currentIndex = 0;
    $scope.sortImg = ["../images/nosort.png", "../images/sort_asc.png", "../images/sort_des.png"];
    $scope.sortState = {
      "userid": 0,
      "nickname": 0,
      "age": 0,
      "regtime": 2,
      "logintime": 0,
      "score": 0,
      "ban": 0,
      "gender": 0
    };
    resetCheckState = function(state) {
      var item, _i, _len, _ref, _results;
      $scope.selectedList = [];
      _ref = $scope.userlist;
      _results = [];
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        item = _ref[_i];
        _results.push($scope.selectedList.push(state));
      }
      return _results;
    };
    $scope.$watch("selectedList", function(newData, oldData) {
      var item, _i, _len, _ref, _results;
      $scope.checkAllBool = $scope.selectedList.length > 0 && !(__indexOf.call($scope.selectedList, false) >= 0);
      $rootScope.sel.nCount = 0;
      _ref = $scope.selectedList;
      _results = [];
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        item = _ref[_i];
        if (item) {
          _results.push($rootScope.sel.nCount++);
        } else {
          _results.push(void 0);
        }
      }
      return _results;
    }, true);
    $scope.menuClass = function(page) {
      var current, _ref;
      current = $location.path().substring(1);
      return (_ref = page === current) != null ? _ref : {
        "active": ""
      };
    };
    $scope.loginOK = function() {
      if ($rootScope.profile == null) {
        userObj.userid = userId;
        userObj.getInfo(userToken, function(retData) {
          var data;
          if (checkRequest(retData)) {
            if ((retData.profile == null) || retData.profile.length === 0) {
              retData.profile = "../images/lanhan.png";
            }
            data = {
              "profile": retData.profile
            };
            return app.checkUser(data);
          }
        });
      }
      return $scope.getUserList(0);
    };
    $scope.search = function(pageIndex, reset) {
      var searchDetail;
      if (reset == null) {
        reset = true;
      }
      searchDetail = {
        gender: "",
        age: "",
        ban_status: "",
        keyword: $scope.searchData.data
      };
      if ($scope.filtState) {
        if ($scope.typeIndex === 1) {
          if ($scope.filtItemIndex === 0) {
            searchDetail.gender = "male";
          } else {
            searchDetail.gender = "female";
          }
        } else if ($scope.typeIndex === 2) {
          if ($scope.filtItemIndex === 0) {
            searchDetail.age = "0-19";
          } else if ($scope.filtItemIndex === 1) {
            searchDetail.age = "20-40";
          } else {
            searchDetail.age = "41-100";
          }
        } else if ($scope.typeIndex === 3) {
          if ($scope.filtItemIndex === 0) {
            searchDetail.ban_status = "normal";
          } else if ($scope.filtItemIndex === 1) {
            searchDetail.ban_status = "lock";
          } else {
            searchDetail.ban_status = "ban";
          }
        }
      } else {
        searchDetail.gender = "";
        searchDetail.age = "";
        searchDetail.ban_status = "";
      }
      return User.search(userToken, searchDetail.keyword, searchDetail.gender, searchDetail.age, searchDetail.ban_status, sortStr, function(retData, page_index, page_total, total_count) {
        var useritem, _i, _j, _len, _results;
        if (checkRequest(retData)) {
          $scope.arrPage = (function() {
            _results = [];
            for (var _i = 0; 0 <= page_total ? _i < page_total : _i > page_total; 0 <= page_total ? _i++ : _i--){ _results.push(_i); }
            return _results;
          }).apply(this);
          $scope.currentPage = page_index;
          userSearchPageIndex = page_index;
          if ($scope.userlist.length > 0) {
            $scope.userlist = [];
          }
          for (_j = 0, _len = retData.length; _j < _len; _j++) {
            useritem = retData[_j];
            if (useritem.ban_status === "normal") {
              useritem.ban_statusTmp = "正常";
              useritem.banStyle = $rootScope.banStyleList[0];
            } else if (useritem.ban_status === "lock") {
              useritem.ban_statusTmp = "禁言";
              useritem.banStyle = $rootScope.banStyleList[1];
            } else {
              useritem.ban_statusTmp = "拉黑";
              useritem.banStyle = $rootScope.banStyleList[2];
            }
            if ((useritem.profile == null) || useritem.profile.length === 0) {
              useritem.profile = "../images/lanhan.png";
            }
            $scope.userlist.push(useritem);
          }
          if (reset) {
            resetCheckState(false);
          }
          searchMode = true;
          return $scope.$apply();
        }
      }, pageIndex, pageCount);
    };
    $scope.searchChange = function() {
      if ((($scope.searchData.data != null) && $scope.searchData.data.length > 0) || $scope.filtState) {
        return $scope.search(0);
      } else {
        return $scope.getUserList(userPageIndex);
      }
    };
    $scope.getUserList = function(pageIndex, reset) {
      if (reset == null) {
        reset = true;
      }
      return User.list(userToken, sortStr, function(retData, page_index, page_total, total_count) {
        var useritem, _i, _j, _len, _results;
        if (checkRequest(retData)) {
          $scope.arrPage = (function() {
            _results = [];
            for (var _i = 0; 0 <= page_total ? _i < page_total : _i > page_total; 0 <= page_total ? _i++ : _i--){ _results.push(_i); }
            return _results;
          }).apply(this);
          $scope.currentPage = page_index;
          userPageIndex = page_index;
          if ($scope.userlist.length > 0) {
            $scope.userlist = [];
          }
          for (_j = 0, _len = retData.length; _j < _len; _j++) {
            useritem = retData[_j];
            if (useritem.ban_status === "normal") {
              useritem.ban_statusTmp = "正常";
              useritem.banStyle = $rootScope.banStyleList[0];
            } else if (useritem.ban_status === "lock") {
              useritem.ban_statusTmp = "禁言";
              useritem.banStyle = $rootScope.banStyleList[1];
            } else {
              useritem.ban_statusTmp = "拉黑";
              useritem.banStyle = $rootScope.banStyleList[2];
            }
            if ((useritem.profile == null) || useritem.profile.length === 0) {
              useritem.profile = "../images/lanhan.png";
            }
            $scope.userlist.push(useritem);
            $scope.$apply();
          }
          if (reset) {
            resetCheckState(false);
          }
          searchMode = false;
          return $scope.$apply();
        } else {
          return setTimeout($scope.getUserList(pageIndex), 500);
        }
      }, pageIndex, pageCount);
    };
    $scope.selected = function(index) {
      return $scope.selectedList[index] = !$scope.selectedList[index];
    };
    $scope.selectAll = function() {
      var item, _i, _len, _ref, _results;
      $scope.checkAllBool = !$scope.checkAllBool;
      _ref = $scope.userlist;
      _results = [];
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        item = _ref[_i];
        if ($scope.checkAllBool) {
          _results.push(resetCheckState(true));
        } else {
          _results.push(resetCheckState(false));
        }
      }
      return _results;
    };
    $scope.countPageChange = function(index) {
      $scope.selectindex = index;
      if ($scope.selectindex === 0) {
        pageCount = 20;
      } else if ($scope.selectindex === 2) {
        pageCount = 100;
      } else if ($scope.selectindex === 3) {
        pageCount = 200;
      } else {
        pageCount = 50;
      }
      if (searchMode) {
        return $scope.search(0);
      } else {
        $scope.getUserList(0);
        return console.log("get user list from countPageChange");
      }
    };
    $scope.changePage = function(index) {
      if (index >= 0 && index !== $scope.arrPage.length && index !== $scope.currentPage) {
        $scope.currentPage = index;
        if (searchMode) {
          return $scope.search(index);
        } else {
          return $scope.getUserList(index);
        }
      }
    };
    banuserFinish = function(refresh) {
      if ($rootScope.sel.nSuccess === $rootScope.sel.nCount) {
        $rootScope.note.successState = true;
      } else {
        $scope.note.errState = true;
      }
      timer = setTimeout(app.hideNote, 1500);
      if (searchMode && refresh) {
        return $scope.search(userSearchPageIndex, false);
      } else {
        return $scope.getUserList(userPageIndex, false);
      }
    };
    $scope.banUser = function(nState, blDetail) {
      var index, nBanState, nStateTmp, refresh, selectCount, _results;
      if (blDetail == null) {
        blDetail = false;
      }
      refresh = false;
      nStateTmp = nState;
      $rootScope.sel.nSuccess = 0;
      if (nState > 0) {
        nState = 30 * 24 * 60 * 60;
      }
      selectCount = 0;
      index = 0;
      _results = [];
      while (index < $scope.selectedList.length) {
        if ($scope.selectedList[index]) {
          userObj.userid = $scope.userlist[index].userid;
          nBanState = 0;
          if ($scope.userlist[index].ban_status === "lock") {
            nBanState = 1;
          } else if ($scope.userlist[index].ban_status === "ban") {
            nBanState = -1;
          }
          if (nBanState !== nStateTmp) {
            index++;
            _results.push(userObj.ban(userToken, nState, function(retData) {
              selectCount++;
              if (checkRequest(retData)) {
                $rootScope.sel.nSuccess++;
                refresh = true;
              }
              if (selectCount === $rootScope.sel.nCount) {
                return banuserFinish(refresh);
              }
            }));
          } else {
            selectCount++;
            index++;
            $rootScope.sel.nSuccess++;
            if (selectCount === $rootScope.sel.nCount) {
              _results.push(banuserFinish(refresh));
            } else {
              _results.push(void 0);
            }
          }
        } else {
          _results.push(index++);
        }
      }
      return _results;
    };
    $scope.changeFilt = function() {
      if ($scope.filtState) {
        $scope.filtState = 0;
        $scope.selectItem = $scope.selectItemList[0];
        $scope.typeIndex = 0;
        $scope.filtItemIndex = 0;
        if ($scope.searchData.data.length > 0) {
          return $scope.search(userSearchPageIndex);
        } else {
          return $scope.getUserList(userPageIndex);
        }
      } else {
        $scope.filtState = 1;
        if ($scope.typeIndex !== 0) {
          return $scope.search(0);
        }
      }
    };
    $scope.filtChange = function(index, type) {
      if (type === 0) {
        $scope.selectItem = $scope.selectItemList[index];
        $scope.typeIndex = index;
        $scope.filtItemIndex = 0;
        if ($scope.filtState) {
          return $scope.search(0);
        }
      } else {
        $scope.filtItemIndex = index;
        if ($scope.filtState) {
          return $scope.search(0);
        }
      }
    };
    return $scope.sort = function(str) {
      sortStr = str;
      if ($scope.sortState[str] === 1) {
        sortStr = "-" + sortStr;
        $scope.sortState[str] = 2;
      } else if ($scope.sortState[str] === 2) {
        $scope.sortState[str] = 1;
      } else {
        $scope.sortState = {
          "userid": 0,
          "nickname": 0,
          "age": 0,
          "regtime": 0,
          "logintime": 0,
          "score": 0,
          "ban": 0,
          "gender": 0
        };
        $scope.sortState[str] = 2;
        sortStr = "-" + sortStr;
      }
      if (searchMode) {
        return $scope.search(0);
      } else {
        return $scope.getUserList(0);
      }
    };
  }
]);

userdetailController = app.controller('userdetailController', [
  'app', '$scope', '$routeParams', '$rootScope', function(app, $scope, $routeParams, $rootScope) {
    var getArrayString, userId, userToken, userinfoId;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    $scope.userDetail = {};
    $scope.userDetailBtnState = [1, 1, 1];
    userToken = app.getCookie("access_token");
    userId = app.getCookie("id");
    userinfoId = $routeParams.id;
    getArrayString = function(arrData) {
      if ((arrData != null) && arrData.length > 0) {
        return String(arrData);
      }
      return "未知";
    };
    $scope.getuserInfo = function() {
      userObj.userid = userinfoId;
      return userObj.getInfo(userToken, function(retData) {
        if (checkRequest(retData)) {
          if (retData.ban_status === "normal") {
            retData.ban_statusTmp = "正常";
            retData.banStyle = $rootScope.banStyleList[0];
          } else if (retData.ban_status === "lock") {
            retData.ban_statusTmp = "禁言";
            retData.banStyle = $rootScope.banStyleList[1];
          } else {
            retData.ban_statusTmp = "拉黑";
            retData.banStyle = $rootScope.banStyleList[2];
          }
          if ((retData.profile == null) || retData.profile.length === 0) {
            retData.profile = "../images/lanhan.png";
          }
          $scope.userDetail = retData;
          $scope.userDetailBtnState = [1, 1, 1];
          if (retData.ban_status === "normal") {
            $scope.userDetailBtnState[0] = 0;
          } else if (retData.ban_status === "lock") {
            $scope.userDetailBtnState[1] = 0;
          } else {
            $scope.userDetailBtnState[2] = 0;
          }
          $scope.userDetail.equips.hardwares = getArrayString(retData.equips.hardwares);
          $scope.userDetail.equips.shoes = getArrayString(retData.equips.shoes);
          $scope.userDetail.equips.softwares = getArrayString(retData.equips.softwares);
          return $scope.$apply();
        }
      });
    };
    return $scope.banUser = function(nState) {
      $rootScope.sel.nSuccess = 0;
      if (nState > 0) {
        nState = 30 * 24 * 60 * 60;
      }
      userObj.userid = userinfoId;
      return userObj.ban(userToken, nState, function(retData) {
        var refresh;
        $rootScope.sel.nSuccess++;
        if (checkRequest(retData)) {
          $scope.userDetailBtnState = [1, 1, 1];
          refresh = true;
          $rootScope.note.successState = true;
          if (nState === 0) {
            $scope.userDetail.ban_statusTmp = "正常";
            $scope.userDetailBtnState[0] = 0;
            $scope.userDetail.banStyle = $rootScope.banStyleList[0];
          } else if (nState > 0) {
            $scope.userDetail.ban_statusTmp = "禁言";
            $scope.userDetailBtnState[1] = 0;
            $scope.userDetail.banStyle = $rootScope.banStyleList[1];
          } else {
            $scope.userDetail.ban_statusTmp = "拉黑";
            $scope.userDetailBtnState[2] = 0;
            $scope.userDetail.banStyle = $rootScope.banStyleList[2];
          }
        } else {
          $rootScope.note.errState = true;
        }
        $rootScope.$apply();
        return setTimeout(app.hideNote, 1500);
      });
    };
  }
]);

articleListController = app.controller('articleListController', [
  'app', '$scope', '$routeParams', '$rootScope', function(app, $scope, $routeParams, $rootScope) {
    var articlePageIndex, articleSearchPageIndex, pageCount, searchMode, searchStr, tagID, userId, userToken;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    userToken = app.getCookie("access_token");
    userId = app.getCookie("id");
    pageCount = 50;
    searchMode = false;
    articlePageIndex = 0;
    articleSearchPageIndex = 0;
    searchStr = "";
    $scope.articleList = {};
    $scope.parent = {
      'sumModel': 0,
      'title': ''
    };
    $scope.arrPage = [1];
    $scope.currentPage = 0;
    $scope.dropdowmItems = ["20项", "50项", "100项", "200项"];
    $scope.selectindex = 1;
    $scope.searchData = {
      "data": ""
    };
    $scope.getArticleList = function(page_index) {
      return Article.list(userToken, function(retData, pageindex, page_total, total_count) {
        var _i, _results;
        if (checkRequest(retData)) {
          $scope.arrPage = (function() {
            _results = [];
            for (var _i = 0; 0 <= page_total ? _i < page_total : _i > page_total; 0 <= page_total ? _i++ : _i--){ _results.push(_i); }
            return _results;
          }).apply(this);
          $scope.currentPage = pageindex;
          articlePageIndex = pageindex;
          if ($scope.articleList.length > 0) {
            $scope.articleList = [];
          }
          $scope.articleList = retData;
          searchMode = false;
          return $scope.$apply();
        }
      }, '', page_index, pageCount);
    };
    $scope.search = function(pageIndex) {
      return Article.search(userToken, searchStr, function(retData, page_index, page_total, total_count) {
        var _i, _results;
        if (checkRequest(retData)) {
          $scope.arrPage = (function() {
            _results = [];
            for (var _i = 0; 0 <= page_total ? _i < page_total : _i > page_total; 0 <= page_total ? _i++ : _i--){ _results.push(_i); }
            return _results;
          }).apply(this);
          $scope.currentPage = page_index;
          articleSearchPageIndex = page_index;
          if ($scope.articleList.length > 0) {
            $scope.articleList = [];
          }
          $scope.articleList = retData;
          searchMode = true;
          return $scope.$apply();
        }
      }, '', pageIndex, pageCount);
    };
    $scope.countPageChange = function(index) {
      $scope.selectindex = index;
      if ($scope.selectindex === 0) {
        pageCount = 20;
      } else if ($scope.selectindex === 2) {
        pageCount = 100;
      } else if ($scope.selectindex === 3) {
        pageCount = 200;
      } else {
        pageCount = 50;
      }
      if (searchMode) {
        return $scope.search(0);
      } else {
        return $scope.getArticleList(0);
      }
    };
    $scope.changePage = function(index) {
      if (index >= 0 && index !== $scope.arrPage.length && index !== $scope.currentPage) {
        $scope.currentPage = index;
        if (searchMode) {
          return $scope.search(index);
        } else {
          return $scope.getArticleList(index);
        }
      }
    };
    $scope.searchChange = function() {
      if ((($scope.searchData.data != null) && $scope.searchData.data.length > 0) || $scope.filtState) {
        searchStr = $scope.searchData.data;
        return $scope.search(0);
      } else {
        return $scope.getArticleList(articlePageIndex);
      }
    };
    tagID = $routeParams.tagid;
    if (tagID != null) {
      searchStr = tagID;
      return $scope.search(0);
    } else {
      return $scope.getArticleList(0);
    }
  }
]);

articleDetailController = app.controller('articleDetailController', [
  'app', '$scope', '$routeParams', '$rootScope', function(app, $scope, $routeParams, $rootScope) {
    var articleID, userId, userToken;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    articleID = $routeParams.artid;
    userToken = app.getCookie("access_token");
    userId = app.getCookie("id");
    $scope.article = {};
    return $scope.initArtDetail = function() {
      return articleObj.getInfo(articleID, userToken, function(retData) {
        if (checkRequest(retData)) {
          $scope.article = retData;
          $scope.$apply();
          return console.log(retData);
        }
      });
    };
  }
]);
