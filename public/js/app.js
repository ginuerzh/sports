var Error, Util;

Util = (function() {
  function Util() {}

  Util.host = "http://172.24.222.54:8080";

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
    return Util._post(Util.host + "/admin/login", {
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
    return Util._post(Util.host + "/admin/logout", {
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
    return Util._get(Util.host + "/admin/user/list", {
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
    return Util._get(Util.host + "/admin/user/search", {
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
    return Util._get(Util.host + "/admin/user/info", {
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
    return Util._post(Util.host + "/admin/user/ban", {
      userid: this.userid,
      duration: duration,
      access_token: token
    }, function(resp) {
      if (Error._hasError(resp)) {
        return callb00ack(new Error(resp.error_id, resp.error_desc));
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
    this.equips = {
      shoes: "",
      hardwares: "",
      softwares: ""
    };
    if (data.equips != null) {
      if (data.equips.shoes !== null && data.equips.shoes.length > 0) {
        this.equips.shoes = data.equips.shoes.join(",");
      }
      if (data.equips.hardwares !== null && data.equips.hardwares.length > 0) {
        this.equips.hardwares = data.equips.hardwares.join(",");
      }
      if (data.equips.softwares !== null && data.equips.softwares.length > 0) {
        this.equips.softwares = data.equips.softwares.join(",");
      }
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
    return Util._get(Util.host + "/admin/article/list", {
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
    return Util._get(Util.host + "/admin/article/timeline", {
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

  Article.search = function(token, keyword, tag, callback, sort, pageIndex, pageCount) {
    if (sort == null) {
      sort = '';
    }
    if (pageIndex == null) {
      pageIndex = 0;
    }
    if (pageCount == null) {
      pageCount = 50;
    }
    return Util._get(Util.host + "/admin/article/search", {
      keyword: keyword,
      tag: tag,
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
    return Util._get(Util.host + "/admin/article/info", {
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
    return Util._post(Util.host + "/admin/article/post", {
      article_id: this.article_id,
      author: this.author,
      contents: this.contents,
      tags: [],
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

  Article.prototype["delete"] = function(token, callback) {
    return Util._post(Util.host + "/admin/article/delete", {
      article_id: this.article_id,
      access_token: token
    }, function(resp) {
      if (Error._hasError(resp)) {
        return callback(new Error(resp.error_id, resp.error_desc));
      } else {
        return callback(resp);
      }
    });
  };

  Article.prototype.getComments = function(token, callback, pageIndex, pageCount) {
    if (pageIndex == null) {
      pageIndex = 0;
    }
    if (pageCount == null) {
      pageCount = 10;
    }
    return Util._get(Util.host + "/admin/article/comments", {
      article_id: this.article_id,
      access_token: token,
      page_index: pageIndex,
      page_count: pageCount
    }, function(resp) {
      if (Error.hasError(resp)) {
        return callback(new Error(resp.error_id, resp.error_desc));
      } else {
        this._comments(resp.comments);
        return callback(this.comments, resp.page_index, resp.page_total, resp.total_number);
      }
    });
  };

  Article.prototype._update = function(data) {
    var tag, _i, _len, _ref;
    this.article_id = data.article_id;
    this.parent = data.parent;
    if ((data.author != null) && (data.author.userid != null)) {
      this.author = new User(data.author.userid)._update(data.author);
    }
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

var app, articleObj, checkRequest, urlPath, userId, userObj, userToken;

app = angular.module('app', ['ngRoute', 'akoenig.deckgrid', 'smart-table', 'ng.ueditor']);

app.constant('app', {
  version: Date.now()
});

app.filter('length', [
  'utf8', 'JSONKit', function(utf8, JSONKit) {
    var text;
    text = JSONKit.toStr(text);
    return utf8.stringToBytes(text).length;
  }
]);

app.filter('cutText', [
  'utf8', 'JSONKit', function(utf8, JSONKit) {
    return function(text, len) {
      var bytes, _ref;
      text = JSONKit.toStr(text).trim();
      bytes = utf8.stringToBytes(text);
      len = (_ref = len > 0) != null ? _ref : {
        len: 0
      };
      if (bytes.length > len) {
        bytes.length = len;
        text = utf8.bytesToString(bytes);
        text = text.slice(0, -2) + '…';
      }
      return text;
    };
  }
]);

userObj = new User();

articleObj = new Article();

userToken = "";

userId = "";

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
  var aricledetail, ariclelist, articleimport, login, tasklist, tasklisthistory, userdetail, userlist;
  login = {
    templateUrl: 'html/user-login.html',
    controller: 'loginController'
  };
  userlist = {
    templateUrl: 'html/user-list.html',
    controller: 'userlistController'
  };
  userdetail = {
    templateUrl: 'html/user-details.html',
    controller: 'userdetailController'
  };
  ariclelist = {
    templateUrl: 'html/article-list.html',
    controller: 'articleListController'
  };
  aricledetail = {
    templateUrl: 'html/article-detail.html',
    controller: 'articleDetailController'
  };
  tasklist = {
    templateUrl: 'html/task-list.html',
    controller: 'tasklistController'
  };
  articleimport = {
    templateUrl: 'html/article-import.html',
    controller: 'articleimportController'
  };
  tasklisthistory = {
    templateUrl: 'html/task-listhistory.html',
    controller: 'tasklisthistoryController'
  };
  return $routeProvider.when('/', login).when('/userlist', userlist).when('/detail/:id', userdetail).when('/articledetail/:artid', aricledetail).when('/1', userlist).when('/2', ariclelist).when('/3', tasklist).when('/4', articleimport).when('/tag/:tagid', ariclelist).when('/tasklisthistory', tasklisthistory);
});

app.run([
  'app', '$rootScope', 'utils', '$filter', function(app, $rootScope, utils, $filter) {
    $rootScope.isLogin = utils.getItem('isLogin');
    $rootScope.profile = utils.getItem('profile');
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
    $rootScope.global = {
      "ContentMinLen": 0,
      "ContentMaxLen": 1024
    };
    app.filter = $filter;
    if ($rootScope.isLogin == null) {
      $rootScope.isLogin = false;
    } else {
      userToken = utils.getItem("access_token");
      userId = utils.getItem("id");
    }
    app.checkUser = function(data) {
      var tmp;
      if (data.isLogin != null) {
        tmp = utils.setItem('isLogin', data.isLogin);
        $rootScope.isLogin = data.isLogin;
      }
      if (data.userid != null) {
        utils.setItem('id', data.userid);
        $rootScope.id = data.userid;
        userId = data.userid;
      }
      if (data.access_token != null) {
        utils.setItem('access_token', data.access_token);
        $rootScope.access_token = data.access_token;
        userToken = data.access_token;
      }
      if (data.profile != null) {
        utils.setItem('profile', data.profile);
        return $rootScope.profile = data.profile;
      }
    };
    app.getCookie = function(key) {
      console.log(utils.getItem(key));
      return utils.getItem(key);
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
    $rootScope.navBarItems = ["首页", "用户管理", "博文管理", "任务管理", "文章导入"];
    return app.rootScope = $rootScope;
  }
]);

app.factory('ArticleFac', [
  '$q', '$http', function($q, $http) {
    var getlist;
    getlist = function() {
      return $http.get('http://172.24.222.54:8080/admin/article/list', {
        params: {
          sort: '',
          page_index: 0,
          page_count: 50,
          access_token: userToken
        }
      }).then(function(response) {
        if (typeof response.data === 'object') {
          return response.data;
        } else {
          return $q.reject(response.data);
        }
      }, function(response) {
        return $q.reject(response.data);
      });
    };
    return {
      getlist: getlist
    };
  }
]).factory('utils', function() {
  var factory;
  return factory = {
    getItem: function(item) {
      var data;
      data = window.localStorage.getItem("app_local_data");
      if (!data) {
        data = {};
      } else {
        data = JSON.parse(data);
      }
      return data[item];
    },
    setItem: function(key, value) {
      var data;
      data = window.localStorage.getItem("app_local_data");
      if (!data) {
        data = {};
        data[key] = value;
      } else {
        data = JSON.parse(data);
        data[key] = value;
      }
      window.localStorage.setItem("app_local_data", JSON.stringify(data));
      return true;
    },
    removeItem: function(key) {
      var data;
      data = window.localStorage.getItem("app_local_data");
      if (!data) {
        data = {};
      } else {
        data = JSON.parse(data);
        delete data[key];
      }
      window.localStorage.setItem("app_local_data", JSON.stringify(data));
      return true;
    },
    logout: function() {
      return window.localStorage.removeItem("app_local_data");
    }
  };
});

app.filter("statusname", function() {
  return function(data) {
    var status;
    status = "其他状态";
    switch (data) {
      case "FINISH":
        status = "审核通过";
        break;
      case "UNFINISH":
        status = "被拒绝";
        break;
      case "AUTHENTICATION":
        status = "待审核";
    }
    return status;
  };
});

app.filter("typename", function() {
  return function(data) {
    var typename;
    typename = "其他任务";
    switch (data) {
      case "PHYSIQUE":
        typename = "跑步任务";
    }
    return typename;
  };
});

app.directive('genParseMd', [
  'mdParse', 'sanitize', 'pretty', 'isVisible', '$timeout', function(mdParse, sanitize, pretty, isVisible, $timeout) {
    return function(scope, element, attr) {
      var parseDoc;
      $scope.$watch(attr.genParseMd, function(value) {
        if (isVisible(element)) {
          return parseDoc(value);
        } else {
          return $timeout(function() {
            return parseDoc(value);
          }, 500);
        }
      });
      return parseDoc = function(value) {
        if (angular.isDefined(value)) {
          value = mdParse(value);
          value = sanitize(value);
          element.html(value);
          angular.forEach(element.find('code'), function(value) {
            value = angular.element(value);
            if (!value.parent().is('pre')) {
              return value.addClass('prettyline');
            }
          });
          element.find('pre').addClass('prettyprint');
          element.find('a').attr('target', function() {
            if (this.host !== location.host) {
              return '_blank';
            }
          });
          return pretty();
        }
      };
    };
  }
]);

app.directive('zjcustomize', function() {
  return {
    restrict: 'E',
    template: '<div>Hissss3333<span ng-transclude></span> there</div>',
    transclude: true,
    link: function(scope, element, attrs) {
      element.css('background-color', 'white');
      element.bind('mouseover', function() {
        if (element.css("color") !== "#C10066") {
          return element.find('span').css({
            "color": "#C10066"
          });
        } else {
          return element.find('span').css({
            "color": 'white'
          });
        }
      });
      return console.log("enter zjcustomize");
    }
  };
});

var articleDetailController;

articleDetailController = app.controller('articleDetailController', [
  'app', '$scope', '$routeParams', '$rootScope', '$filter', function(app, $scope, $routeParams, $rootScope, $filter) {
    var articleID, converter, getArticleByUser, initReply;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    converter = new Markdown.Converter();
    articleID = $routeParams.artid;
    $scope.article = {
      "authorInfo": {
        "articlesList": {}
      }
    };
    $scope.comment = {
      title: '',
      content: '',
      refer: '',
      replyToComment: ''
    };
    $scope.parent = {
      wmdPreview: false,
      contentBytes: 0
    };
    $scope.replyMoving = {};
    $scope.commentMoving = {};
    $scope.markdownModal = {
      title: "Markdown简明语法",
      cancelBtn: "返回"
    };
    initReply = function() {
      var article, comment;
      comment = $scope.comment;
      article = $scope.article;
      comment.replyToComment = '';
      comment.title = '评论：' + article.cover_text;
      comment.content = '';
      return comment.refer = article.article_id;
    };
    getArticleByUser = function(userId) {
      return Article.timeline(userId, userToken, function(retData) {
        if (checkRequest(retData)) {
          $scope.article.author.articlesList = retData;
          return $scope.$apply();
        }
      });
    };
    $scope.initArtDetail = function() {
      return articleObj.getInfo(articleID, userToken, function(retData) {
        if (checkRequest(retData)) {
          $scope.article = retData;
          document.getElementById('markdownstring').innerHTML = converter.makeHtml(retData.contents);
          return getArticleByUser(retData.author.userid);
        }
      });
    };
    $scope.reply = function(articleID) {
      var comment;
      comment = $scope.comment;
      comment.refer = articleID;
      $scope.parent.wmdPreview = false;
      if (articleID === $scope.article.article_id) {
        initReply();
      } else {
        comment.replyToComment = articleID;
        $scope.replyMoving.appendTo('#' + articleID);
      }
      return $scope.replyMoving.scrollIntoView();
    };
    $scope.wmdPreview = function() {
      $scope.parent.wmdPreview = !$scope.parent.wmdPreview;
      return $scope.replyMoving.scrollIntoView(true);
    };
    $scope.submit = function() {
      if ($scope.comment.content.length > 0) {
        articleObj.article_id = $scope.article.article_id;
        articleObj.author = userId;
        articleObj.contents = $scope.comment.content;
        return articleObj.post(userToken, function(retData) {
          if (checkRequest(retData)) {
            return console.log("postSuccess");
          }
        });
      }
    };
    $scope.deleteArticle = function(articleId) {
      $rootScope.sel.nCount = 1;
      $rootScope.note.successState = false;
      $scope.note.errState = false;
      articleObj.article_id = articleId;
      return articleObj["delete"](userToken, function(retData) {
        if (checkRequest(retData)) {
          $rootScope.sel.nSuccess = 1;
          $rootScope.note.successState = true;
          console.log("delete article success");
        } else {
          $rootScope.sel.nSuccess = 1;
          $rootScope.note.errState = true;
        }
        $scope.$apply();
        window.location.href = "#/2";
        return setTimeout(app.hideNote, 1500);
      });
    };
    $scope.checkContentMin = function(scope, model) {
      var length;
      length = $filter.lengthFn(model.$value);
      $scope.parent.contentBytes = length;
      return length >= $rootScope.global.ContentMinLen;
    };
    $scope.checkContentMax = function(scope, model) {
      return $filter.lengthFn(model.$value) <= $rootScope.global.ContentMaxLen;
    };
    return $scope.wmdPreview = function() {
      $scope.parent.wmdPreview = !$scope.parent.wmdPreview;
      return $scope.replyMoving.scrollIntoView(true);
    };
  }
]);

var articleListController;

articleListController = app.controller('articleListController', [
  'app', '$scope', '$routeParams', '$rootScope', 'ArticleFac', function(app, $scope, $routeParams, $rootScope, ArticleFac) {
    var articlePageIndex, articleSearchPageIndex, pageCount, searchMode, searchStr, tagID;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
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
      return ArticleFac.getlist().then(function(data) {
        $scope.articleList = [];
        $scope.articleList = data.articles;
        $scope.arrPage = [0, data.page_total];
        articlePageIndex = $scope.currentPage = data.page_index;
        $scope.$apply();
        return console.log(data);
      }, function(error) {
        return console.log(error);
      });
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

var articleimportController;

articleimportController = app.controller('articleimportController', [
  'app', '$scope', '$rootScope', 'articleService', function(app, $scope, $rootScope, articleService) {
    var refreshinput;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    console.log($scope.content);
    $scope.tags = [
      {
        id: "SPORT_THEORY",
        tagname: "跑步圣经"
      }, {
        id: "SPORT_LOG",
        tagname: "运动日志"
      }, {
        id: "EQUIP_BLOG",
        tagname: "我爱装备"
      }, {
        id: "SPORT_LIFE",
        tagname: "运动生活"
      }, {
        id: "PRODUCT_PROPOSAL",
        tagname: "产品建议"
      }
    ];
    refreshinput = function() {
      $scope.title = "";
      return $scope.content = "";
    };
    refreshinput();
    $scope.tag = $scope.tags[0];
    return $scope.importarticle = function() {
      var imglist;
      console.log("enter importarticle");
      if ($scope.title === "") {
        console.log("before get imagelist");
        alert("please input the title");
        return imglist = articleService.getimagelist($scope.content);
      } else {
        imglist = articleService.getimagelist($scope.content);
        console.log(imglist);
        articleService.articlepost($scope.title, "1419305934614", imglist, $scope.content, $scope.tag.id);
        return refreshinput();
      }
    };
  }
]);

var loginController;

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
      console.log("sdfasdf");
      console.log(app.getCookie("isLogin"));
      if ($routeParams.index == null) {
        return $rootScope.isLogin = false;
      }
    };
  }
]);

var tasklistController;

tasklistController = app.controller('tasklistController', [
  'app', '$scope', '$rootScope', 'taskService', function(app, $scope, $rootScope, taskService) {
    var refreshtable;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    refreshtable = function() {
      $scope.rowCollection = taskService.gettasklist("Auditting");
      return $scope.displayedCollection = [].concat($scope.rowCollection);
    };
    refreshtable();
    console.log("before enter the taskservice");
    $scope.itemsByPage = 20;
    $scope.Approve = function(row) {
      this.reason = row.reason.trim();
      taskService.taskapprove(row.userid, row.taskid, this.reason);
      return refreshtable();
    };
    return $scope.Reject = function(row) {
      this.reason = row.reason.trim();
      if (this.reason === "") {
        return alert("please input the reason for the rejection");
      } else {
        taskService.taskreject(row.userid, row.taskid, this.reason);
        return refreshtable();
      }
    };
  }
]);

var tasklisthistoryController;

tasklisthistoryController = app.controller('tasklisthistoryController', [
  'app', '$scope', '$rootScope', 'taskService', function(app, $scope, $rootScope, taskService) {
    var refreshtable;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    refreshtable = function() {
      $scope.rowCollection = taskService.gettasklist("Audited");
      return $scope.displayedCollection = [].concat($scope.rowCollection);
    };
    refreshtable();
    console.log("before enter the taskservice");
    $scope.itemsByPage = 20;
    $scope.Approve = function(row) {
      this.reason = row.reason.trim();
      taskService.taskapprove(row.userid, row.taskid, this.reason);
      return refreshtable();
    };
    return $scope.Reject = function(row) {
      this.reason = row.reason.trim();
      if (this.reason === "") {
        return alert("please input the reason for the rejection");
      } else {
        taskService.taskreject(row.userid, row.taskid, this.reason);
        return refreshtable();
      }
    };
  }
]);

var userdetailController;

userdetailController = app.controller('userdetailController', [
  'app', '$scope', '$routeParams', '$rootScope', function(app, $scope, $routeParams, $rootScope) {
    var getArrayString, userinfoId;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    $scope.userDetail = {};
    $scope.userDetailBtnState = [1, 1, 1];
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
            $scope.userDetail.ban_statusTmp = "正常1";
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

var userlistController,
  __indexOf = [].indexOf || function(item) { for (var i = 0, l = this.length; i < l; i++) { if (i in this && this[i] === item) return i; } return -1; };

userlistController = app.controller('userlistController', [
  'app', '$scope', '$rootScope', function(app, $scope, $rootScope) {
    var banuserFinish, pageCount, resetCheckState, searchMode, sortStr, timer, userPageIndex, userSearchPageIndex;
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
    $scope.dropdowmItems = ["20222332项", "50项", "100项", "200项"];
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

app.factory('taskq', [
  '$q', '$http', function($q, $http) {
    return {
      gettasklist: function() {
        return $http.get('http://172.24.222.54:8080/admin/task/list', {
          params: {
            page_index: 0,
            page_count: 50,
            access_token: userToken
          }
        });
      },
      taskaudit: function(userid, taskid, pass, reason) {
        return $http.post('http://172.24.222.54:8080/admin/task/auth', {
          userid: userid,
          task_id: taskid,
          pass: pass,
          reason: reason,
          access_token: userToken
        });
      }
    };
  }
]);

app.factory('articleq', [
  '$q', '$http', function($q, $http) {
    return {
      articlepost: function(article_id, title, author, imglist, contents, tag) {
        return $http.post('http://172.24.222.54:8080/admin/article/post', {
          article_id: article_id,
          title: title,
          author: author,
          image: imglist,
          contents: contents,
          tags: tag,
          access_token: userToken
        });
      }
    };
  }
]);

app.factory('taskService', [
  '$q', 'taskq', function($q, $taskq) {
    return {
      gettasklist: function(tasktype) {
        var tasklist;
        console.log("enter the taskservice");
        tasklist = [];
        $taskq.gettasklist().then(function(response) {
          var task, taskitem, taskjson, _i, _len, _ref, _results;
          if (typeof response.data === 'object') {
            _ref = response.data.users;
            _results = [];
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              taskitem = _ref[_i];
              console.log(taskitem);
              _results.push((function() {
                var _j, _len1, _ref1, _results1;
                _ref1 = taskitem.tasks;
                _results1 = [];
                for (_j = 0, _len1 = _ref1.length; _j < _len1; _j++) {
                  task = _ref1[_j];
                  if (task.status === null) {
                    continue;
                  }
                  taskjson = {
                    taskid: task.task_id,
                    type: task.type,
                    desc: task.desc,
                    status: task.status,
                    reason: task.reason,
                    images: task.images,
                    userid: taskitem.userid,
                    nickname: taskitem.nickname,
                    profile: taskitem.profile
                  };
                  switch (tasktype) {
                    case "Auditting":
                      if (taskjson.status === "AUTHENTICATION") {
                        _results1.push(tasklist.push(taskjson));
                      } else {
                        _results1.push(void 0);
                      }
                      break;
                    case "Audited":
                      if (taskjson.status === "FINISH" || taskjson.status === "UNFINISH") {
                        _results1.push(tasklist.push(taskjson));
                      } else {
                        _results1.push(void 0);
                      }
                      break;
                  }
                }
                return _results1;
              })());
            }
            return _results;
          } else {
            return $q.reject(response.data);
          }
        }, function(response) {
          return $q.reject(response.data);
        });
        console.log(tasklist);
        return tasklist;
      },
      taskapprove: function(userid, taskid, reason) {
        console.log("enter taskapprove");
        return $taskq.taskaudit(userid, taskid, true, reason).success(function(response) {
          return console.log(response);
        });
      },
      taskreject: function(userid, taskid, reason) {
        return $taskq.taskaudit(userid, taskid, false, reason).success(function(response) {
          return console.log("reject");
        });
      }
    };
  }
]);

app.factory('articleService', [
  '$q', 'articleq', function($q, $articleq) {
    return {
      articlepost: function(title, author, imglist, contents, tag) {
        console.log("enter articlepost");
        return $articleq.articlepost("", title, author, imglist, contents, tag).success(function(response) {
          return console.log(response);
        });
      },
      getimagelist: function(contents) {
        var elem, imglist;
        console.log("enter getimagelist");
        imglist = (function() {
          var _i, _len, _ref, _results;
          _ref = $(contents).find("img");
          _results = [];
          for (_i = 0, _len = _ref.length; _i < _len; _i++) {
            elem = _ref[_i];
            _results.push(elem.src);
          }
          return _results;
        })();
        return imglist;
      }
    };
  }
]);
