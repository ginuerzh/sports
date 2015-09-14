var Error, Util, app;

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

  Util._checkDate = function(data) {
    var retData;
    retData = data;
    retData += "";
    return retData.replace(/^(\d)$/, "0$1");
  };

  Util._formatDate = function(date) {
    var d, t;
    d = [date.getFullYear(), Util._checkDate(date.getMonth() + 1), Util._checkDate(date.getDate())].join("-");
    t = [Util._checkDate(date.getHours()), Util._checkDate(date.getMinutes()), Util._checkDate(date.getSeconds())].join(":");
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

app = angular.module('app', ['ngRoute', 'akoenig.deckgrid', 'smart-table', 'ng.ueditor']);

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

  User.search = function(token, keyword, gender, age, ban_status, role, actor, sort, callback, pageIndex, pageCount) {
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
      role: role,
      actor: actor,
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
    this.online = data.online;
    this.onlinetime = data.onlinetime;
    this.actor = data.actor;
    this.admin = data.admin;
    this.sign = data.sign;
    this.emotion = data.emotion;
    this.profession = data.profession;
    this.fond = data.fond;
    this.hometown = data.hometown;
    this.oftenAppear = data.oftenAppear;
    this.action = data.action;
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
    this.email = data.email;
    this.auth = data.auth;
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

app.factory('articleq', [
  '$http', function($http) {
    return {
      articlepost: function(article_id, title, author, imglist, contents, tag) {
        return $http.post(Util.host + "/admin/article/post", {
          article_id: article_id,
          title: title,
          author: author,
          image: imglist,
          contents: contents,
          tags: tag,
          access_token: userToken
        });
      },
      getarticlelist: function(sort, tag, page_index, page_count) {
        return $http.get(Util.host + "/admin/article/list", {
          params: {
            sort: sort,
            page_index: page_index,
            page_count: page_count,
            tag: tag,
            access_token: userToken
          }
        });
      },
      searcharticle: function(keyword, tag, sort, pageIndex, pageCount) {
        return $http.get(Util.host + "/admin/article/search", {
          params: {
            keyword: keyword,
            tag: tag,
            sort: sort,
            page_index: pageIndex,
            page_count: pageCount,
            access_token: userToken
          }
        });
      },
      articlemark: function(article_id, type) {
        return $http.post(Util.host + "/admin/article/mark", {
          article_id: article_id,
          type: type,
          access_token: userToken
        });
      },
      topicpost: function(topic, rec, image) {
        return $http.post(Util.host + "/admin/article/topic/post", {
          topic: topic,
          rec: rec,
          image: image,
          access_token: userToken
        });
      },
      gettopiclist: function(sort, tag, page_index, page_count) {
        return $http.get(Util.host + "/admin/article/topic/list", {
          params: {
            sort: sort,
            page_index: page_index,
            page_count: page_count,
            tag: tag,
            access_token: userToken
          }
        });
      },
      deletearticle: function(article_id) {
        return $http.post(Util.host + "/admin/article/delete", {
          article_id: article_id,
          access_token: userToken
        });
      },
      getarticleinfo: function(article_id) {
        return $http.get(Util.host + "/admin/article/info", {
          params: {
            article_id: article_id,
            access_token: userToken
          }
        });
      }
    };
  }
]);

var Task;

Task = (function() {
  function Task() {}

  Task.list = function(token, callback, pageIndex, pageCount) {
    if (pageIndex == null) {
      pageIndex = 0;
    }
    if (pageCount == null) {
      pageCount = 50;
    }
    return Util._get(Util.host + '/admin/task/list', {
      access_token: token
    }, function(resp) {
      if (Error._hasError(resp)) {
        return callback(new Error(resp.error_id, resp.error_desc));
      } else {
        return callback(resp.users, resp.page_index, resp.page_total, resp.total_number);
      }
    });
  };

  Task.timeline = function(token, userid, week, callback) {
    return Util._get(Util.host + '/admin/task/timeline', {
      userid: userid,
      week: week,
      access_token: token
    }, function(resp) {
      var task, tasks, _i, _len, _ref;
      if (Error._hasError(resp)) {
        return callback(new Error(resp.error_id, resp.error_desc));
      } else {
        tasks = [];
        _ref = resp.tasks;
        for (_i = 0, _len = _ref.length; _i < _len; _i++) {
          task = _ref[_i];
          tasks.push(new Task()._update(task));
        }
        return callback(tasks);
      }
    });
  };

  Task.auth = function(token, userid, task_id, pass, reason) {
    return Util._post(Util.host + '/admin/task/auth', {
      userid: userid,
      task_id: task_id,
      pass: pass,
      reason: reason,
      access_token: token
    }, function(resp) {
      if (Error._hasError(resp)) {
        return callback(new Error(resp.error_id, resp.error_desc));
      } else {
        return callback(resp);
      }
    });
  };

  Task.prototype._update = function(data) {
    this.task_id = data.task_id;
    this.type = data.type;
    this.desc = data.desc;
    this.images = data.images;
    this.status = data.status;
    this.reason = data.reason;
    return this;
  };

  return Task;

})();

app.factory('taskq', [
  '$http', function($http) {
    return {
      gettasklist: function(page_index, finished) {
        return $http.get(Util.host + '/admin/task/list', {
          params: {
            finished: finished,
            page_index: page_index,
            page_count: 50,
            access_token: userToken
          }
        });
      },
      taskaudit: function(userid, taskid, pass, reason) {
        return $http.post(Util.host + '/admin/task/auth', {
          userid: userid,
          task_id: taskid,
          pass: pass,
          reason: reason,
          access_token: userToken
        });
      },
      taskauditall: function(auditlist) {
        return $http.post(Util.host + '/admin/task/auth_list', {
          auths: auditlist,
          access_token: userToken
        });
      },
      search: function(nickname, finish, page_count, page_index) {
        return $http.get(Util.host + '/admin/task/timeline', {
          params: {
            nickname: nickname,
            finish: finish,
            page_index: page_index,
            page_count: page_count,
            access_token: userToken
          }
        });
      }
    };
  }
]);

var articleObj, checkDate, checkRequest, converter, urlPath, userId, userObj, userToken;

app.constant('app', {
  version: Date.now()
});

converter = new Markdown.Converter();

userObj = new User();

articleObj = new Article();

userToken = "";

userId = "";

checkRequest = function(reqData) {
  var data;
  if (reqData === null || (reqData.error_id == null)) {
    return true;
  } else if (reqData.error_id === 1003) {
    data = {
      isLogin: false,
      userid: '',
      access_token: '',
      profile: ''
    };
    app.checkUser(data);
    window.location.href = "#/";
  }
  return false;
};

urlPath = function(url) {
  return url;
};

checkDate = function(data) {
  var retData;
  retData = data;
  retData += "";
  return retData.replace(/^(\d)$/, "0$1");
};

app.config(function($routeProvider) {
  var aricledetail, ariclelist, articleimport, authdetail, authenticationlist, configsetting, coversetting, dashboard, login, tasklist, userdetail, userlist;
  login = {
    templateUrl: 'html/user-login.html',
    controller: 'loginController'
  };
  userlist = {
    templateUrl: 'html/user-list.html',
    controller: 'userlistController'
  };
  dashboard = {
    templateUrl: 'html/dashboard-info.html',
    controller: 'dashboardController'
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
  authenticationlist = {
    templateUrl: 'html/authentication-list.html',
    controller: 'authenticationController'
  };
  authdetail = {
    templateUrl: 'html/authentication-detail.html',
    controller: 'authenticationController'
  };
  configsetting = {
    templateUrl: 'html/system-setting.html',
    controller: 'configController'
  };
  coversetting = {
    templateUrl: 'html/cover-setting.html',
    controller: 'coverController'
  };
  return $routeProvider.when('/', login).when('/userlist', userlist).when('/detail/:id', userdetail).when('/articledetail/:artid', aricledetail).when('/0', dashboard).when('/1', userlist).when('/2', ariclelist).when('/3', tasklist).when('/4', articleimport).when('/5', authenticationlist).when('/6', configsetting).when('/7', coversetting).when('/authdetail/:authid', authdetail).when('/tag/:tagid', ariclelist).when('/tasklisthistory', tasklist);
});

app.run([
  'app', '$rootScope', 'utils', '$filter', 'authService', 'taskService', function(app, $rootScope, utils, $filter, authService, taskService) {
    var checkTaskAndAuth;
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
      $rootScope.profile = utils.getItem("profile");
      $rootScope.access_token = userToken;
      $rootScope.id = userId;
      $rootScope.isLogin = utils.getItem("isLogin");
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
    $rootScope.navBarItems = [
      {
        title: "Dashboard",
        shownew: false
      }, {
        title: "用户管理",
        shownew: false
      }, {
        title: "博文管理",
        shownew: false
      }, {
        title: "任务管理",
        shownew: false
      }, {
        title: "文章导入",
        shownew: false
      }, {
        title: "认证管理",
        shownew: false
      }, {
        title: "系统设置",
        shownew: false
      }, {
        title: "封面推荐设置",
        shownew: false
      }
    ];
    app.rootScope = $rootScope;
    $rootScope.checkAuthList = function(data) {
      if (data.authlist.length > 0) {
        return $rootScope.navBarItems[5].shownew = true;
      } else {
        return $rootScope.navBarItems[5].shownew = false;
      }
    };
    $rootScope.checkTaskList = function(data) {
      if (data.tasklist.length > 0) {
        return $rootScope.navBarItems[3].shownew = true;
      } else {
        return $rootScope.navBarItems[3].shownew = false;
      }
    };
    checkTaskAndAuth = function() {
      var authlist;
      authlist = authService.getauthlist('', 0, 20).then($rootScope.checkAuthList);
      taskService.gettasklist(false, 0).then($rootScope.checkTaskList);
      return window.setTimeout(checkTaskAndAuth, 360000);
    };
    if ($rootScope.isLogin) {
      return checkTaskAndAuth();
    }
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
        status = "拒绝";
        break;
      case "AUTHENTICATION":
        status = "待审核";
    }
    return status;
  };
}).filter("typename", function() {
  return function(data) {
    var typename;
    typename = "跑步纪录";
    switch (data) {
      case "PHYSIQUE":
        typename = "跑步任务";
    }
    return typename;
  };
}).filter("profilefilter", function() {
  return function(data) {
    var imagePath;
    if (angular.isString(data)) {
      imagePath = "../images/lanhan.png";
      if (data.length > 0) {
        imagePath = data;
      }
    }
    return imagePath;
  };
}).filter("timeconvert", function() {
  return function(data) {
    var d, date, t;
    date = new Date(data * 1000);
    d = [date.getFullYear(), checkDate(date.getMonth() + 1), checkDate(date.getDate())].join("-");
    t = [checkDate(date.getHours()), checkDate(date.getMinutes()), checkDate(date.getSeconds())].join(":");
    return [d, t].join(" ");
  };
}).filter("timechange", function() {
  return function(data) {
    var hour, min, sec;
    if (data >= 3600) {
      hour = Math.floor(data / 3600);
      min = Math.floor((data % 3600) / 60);
      sec = (data % 3600) % 60;
      return hour + " 小时 " + min + " 分 " + sec + " 秒";
    } else if (data >= 60) {
      min = Math.floor(data / 60);
      sec = data % 60;
      return min + " 分 " + sec + " 秒";
    } else {
      sec = data;
      return sec + " 秒";
    }
  };
}).filter("actorfilter", function() {
  return function(data) {
    var actor;
    actor = "普通用户";
    if (data.length > 0) {
      actor = "教练";
    }
    return actor;
  };
}).filter("yesornofilter", function() {
  return function(data) {
    var yesorno;
    yesorno = "否";
    if (data) {
      yesorno = "是";
    }
    return yesorno;
  };
}).filter("articletitle", function() {
  return function(data) {
    var titlestr;
    titlestr = data;
    if (angular.isString(data)) {
      if (data.length > 50) {
        titlestr = data.substr(0, 50) + "......";
      } else if (data.length === 0) {
        titlestr = "无标题文章";
      }
    }
    return titlestr;
  };
}).filter("articletag", function() {
  return function(data) {
    var tagstr;
    tagstr = data;
    if (angular.isString(data)) {
      switch (data) {
        case "SPORT_LOG":
          tagstr = "运动日志";
          break;
        case "SPORT_THEORY":
          tagstr = "跑步圣经";
          break;
        case "EQUIP_BLOG":
          tagstr = "我爱装备";
          break;
        case "SPORT_LIFE":
          tagstr = "运动生活";
          break;
        case "PRODUCT_PROPOSAL":
          tagstr = "产品建议";
      }
    }
    return tagstr;
  };
}).filter("taskSource", function() {
  return function(data) {
    var taskSource;
    taskSource = data;
    if (data.length === 0) {
      taskSource = "手动";
    }
    return taskSource;
  };
}).filter("authstatus", function() {
  return function(data) {
    var authstr;
    authstr = "未认证";
    if (angular.isString(data)) {
      switch (data) {
        case "verifying":
          authstr = "认证中";
          break;
        case "verified":
          authstr = "已认证";
          break;
        case "refused":
          authstr = "认证拒绝";
      }
    }
    return authstr;
  };
}).filter("authclass", function() {
  return function(data) {
    var authclass;
    authclass = {
      "background-color": "#999999"
    };
    if (angular.isString(data)) {
      switch (data) {
        case "verifying":
          authclass = {
            "background-color": "#f0ad4e"
          };
          break;
        case "verified":
          authclass = {
            "background-color": "#5cb85c"
          };
          break;
        case "refused":
          authclass = {
            "background-color": "#d9534f"
          };
      }
    }
    return authclass;
  };
}).filter("rolefilter", function() {
  return function(data) {
    var rolestr;
    rolestr = "手机";
    if (angular.isString(data)) {
      switch (data) {
        case "weibo":
          rolestr = "微博";
          break;
        case "email":
          rolestr = "邮箱";
      }
    }
    return rolestr;
  };
}).filter("gender", function() {
  return function(data) {
    var gender;
    gender = "男";
    if (angular.isString(data) && data === "female") {
      gender = "女";
    }
    return gender;
  };
}).filter("age", function() {
  return function(data) {
    var age, birth;
    age = 0;
    if ((data != null) && data !== 0) {
      birth = new Date(data * 1000);
      age = new Date().getFullYear() - birth.getFullYear();
    }
    return age;
  };
});

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
}).directive('genPagination', function() {
  return {
    scope: true,
    templateUrl: '../html/mbb-pagination.html',
    link: function(scope, element, attrs) {
      return scope.$watchCollection(attrs.genPagination, function(value) {
        var lastPage, pageIndex, showPages, _ref, _ref1;
        showPages = [];
        lastPage = value.pagetotal;
        pageIndex = value.pageIndex;
        showPages[0] = lastPage;
        while (showPages[0] > 1) {
          showPages.unshift(showPages[0] - 1);
        }
        scope.prev = (_ref = pageIndex <= 1) != null ? _ref : {
          0: pageIndex - 1
        };
        scope.next = (_ref1 = pageIndex >= lastPage) != null ? _ref1 : {
          0: pageIndex + 1
        };
        scope.total = value.total;
        scope.pageIndex = pageIndex;
        scope.showPages = showPages;
        scope.pagetotal = value.pagetotal;
        return scope.paginationTo = function(p) {
          if (p > 0 && p <= scope.pagetotal) {
            return scope.$emit('genPagination', p, element.context.id);
          }
        };
      });
    }
  };
}).directive('genParseMd', function() {
  return {
    link: function(scope, element, attrs) {
      return scope.$watchCollection(attrs.genParseMd, function(value) {
        if (angular.isDefined(value)) {
          value = converter.makeHtml(value);
          element.html(value);
          angular.forEach(element.find('code'), function(value) {
            value = angular.element(value);
            if (!value.parent().is('pre')) {
              return value.addClass('prettyline');
            }
          });
          angular.forEach(element.find('p'), function(value) {
            value = angular.element(value);
            return value.addClass('content-p-show');
          });
          return element.find('a').attr('target', function() {
            if (this.host !== location.host) {
              return '_blank';
            }
          });
        }
      });
    }
  };
});

var articleDetailController;

articleDetailController = app.controller('articleDetailController', [
  'app', '$scope', '$routeParams', '$rootScope', 'articleService', function(app, $scope, $routeParams, $rootScope, articleService) {
    var articleID, getArticleByUser;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    articleID = $routeParams.artid;
    $scope.article = {
      "authorInfo": {
        "articlesList": {}
      }
    };
    $scope.comment = {
      refer: '',
      title: '',
      content: ''
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
          $scope.comment.title = '评论: ' + $scope.article.cover_text;
          return getArticleByUser(retData.author.userid);
        }
      });
    };
    $scope.reback = function() {
      return $scope.comment = {
        refer: '',
        content: ''
      };
    };
    $scope.submit = function() {
      var imglist;
      imglist = articleService.getimagelist($scope.comment.content);
      return articleService.articlepost($scope.article.article_id, $scope.comment.title, userId, imglist, $scope.comment.content, $scope.article.tags).then($scope.initArtDetail);
    };
    return $scope.deleteArticle = function(articleId) {
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
  }
]);

var articleListController;

articleListController = app.controller('articleListController', [
  'app', '$scope', '$routeParams', '$rootScope', 'articleService', function(app, $scope, $routeParams, $rootScope, articleService) {
    var articleMark, articleMarkFailed, articlePageIndex, getTopicList, pageCount, searchMode, searchStr, tagID;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    pageCount = 50;
    searchMode = false;
    articlePageIndex = 0;
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
    $scope.pagination = {};
    articleMarkFailed = function(reason) {
      return alert("设置不成功，" + reason);
    };
    articleMark = function(artcile_id, type) {
      return articleService.articlemark(artcile_id, type).then('', articleMarkFailed);
    };
    getTopicList = function(index) {
      var topiclist;
      topiclist = articleService.gettopiclist('', '', index, pageCount);
      $scope.topiclist = topiclist.articlelist;
      return $scope.topicpagination = topiclist.pagination;
    };
    $scope.getArticleList = function(page_index) {
      var articleinfo;
      articleinfo = articleService.getarticlelist('', '', page_index, pageCount);
      $scope.articleList = articleinfo.articlelist;
      return $scope.pagination = articleinfo.pagination;
    };
    $scope.search = function(pageIndex) {
      var articleinfo;
      articleinfo = articleService.searcharticle(searchStr, tagID, '', pageIndex, pageCount);
      $scope.articleList = articleinfo.articlelist;
      return $scope.pagination = articleinfo.pagination;
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
    $scope.article = function(index) {
      if ($scope.articleList[index].type !== 0) {
        return articleMark($scope.articleList[index].article_id, '');
      }
    };
    $scope.interview = function(index) {
      if ($scope.articleList[index].type !== 1) {
        return articleMark($scope.articleList[index].article_id, 'topic');
      }
    };
    $scope.recommend = function(index) {
      if ($scope.articleList[index].type !== 2) {
        return articleMark($scope.articleList[index].article_id, 'rec');
      }
    };
    $scope.$on('genPagination', function(event, p) {
      event.stopPropagation();
      if (id === 'articlepage') {
        if (searchMode || (typeof tagID !== "undefined" && tagID !== null)) {
          return $scope.search(p);
        } else {
          return $scope.getArticleList(p);
        }
      } else if (id === 'topiclistpage') {
        return getTopicList(p);
      }
    });
    tagID = $routeParams.tagid;
    if (tagID != null) {
      return $scope.search(0);
    } else {
      $scope.getArticleList(0);
      return getTopicList(0);
    }
  }
]);

var articleimportController;

articleimportController = app.controller('articleimportController', [
  'app', '$scope', '$rootScope', 'articleService', function(app, $scope, $rootScope, articleService) {
    var postaticleFail, refreshinput;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
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
    postaticleFail = function(reason) {
      return alert("导入不成功，" + reason);
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
        return articleService.articlepost("", $scope.title, userId, imglist, $scope.content, $scope.tag.id).then(refreshinput, postaticleFail);
      }
    };
  }
]);

var authenticationController;

authenticationController = app.controller('authenticationController', [
  'app', '$scope', '$routeParams', '$rootScope', 'authService', function(app, $scope, $routeParams, $rootScope, authService) {
    var authPost, authid, getAuthList, pageIndex, refreshmain;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    pageIndex = 0;
    authid = "";
    refreshmain = function() {
      if (window.location.href.indexOf('authdetail') > 0) {
        authid = $routeParams.authid;
        return authService.authinfo(authid).success(function(response) {
          if (checkRequest(response)) {
            return $scope.userinfo = response;
          }
        });
      } else {
        return getAuthList();
      }
    };
    getAuthList = function() {
      return authService.getauthlist('', pageIndex, 50).then(function(data) {
        $scope.pagination = data.pagination;
        $scope.authList = data.authlist;
        return $rootScope.checkAuthList(data);
      });
    };
    $scope.$on('genPagination', function(event, p) {
      event.stopPropagation();
      pageIndex = p;
      return refreshmain();
    });
    authPost = function(status, type) {
      var authStatus, authType, authreview;
      authType = "idcard";
      authStatus = status;
      authreview = $scope.userinfo.auth.idcard.auth_review;
      switch (type) {
        case 1:
          authType = "cert";
          authreview = $scope.userinfo.auth.cert.auth_review;
          break;
        case 2:
          authType = "record";
          authreview = $scope.userinfo.auth.record.auth_review;
      }
      return authService.authpost(authid, authType, authStatus, authreview).then(refreshmain);
    };
    $scope.approve = function(type) {
      return authPost("verified", type);
    };
    $scope.cancel = function(type) {
      return authPost("unverified", type);
    };
    $scope.reject = function(type) {
      return authPost("refused", type);
    };
    return refreshmain();
  }
]);

var configController;

configController = app.controller('configController', [
  'app', '$scope', 'configService', function(app, $scope, configService) {
    var getconfiginfo, setconfiginfo;
    $scope.addstate = false;
    $scope.petaddstate = false;
    $scope.adddata = {
      "title": "",
      "url": ""
    };
    $scope.petadddata = "";
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    getconfiginfo = function() {
      var listdata;
      listdata = configService.getconfig();
      $scope.configlist = listdata.videos;
      return $scope.petlist = listdata.pets;
    };
    setconfiginfo = function() {
      var configinfo, item, petlist, _i, _len, _ref;
      petlist = [];
      _ref = $scope.petlist;
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        item = _ref[_i];
        petlist.push(item.content);
      }
      configinfo = {
        "videos": $scope.configlist,
        "pets": petlist
      };
      return configService.setconfig(configinfo).then(getconfiginfo);
    };
    $scope.edit = function(type, index) {
      if (type === 0) {
        return $scope.configlist[index].isedit = true;
      } else if (type === 1) {
        return $scope.petlist[index].isedit = true;
      }
    };
    $scope.editsave = function(type, index) {
      if (type === 0) {
        $scope.configlist[index].isedit = false;
      } else if (type === 1) {
        $scope.petlist[index].isedit = false;
      }
      return setconfiginfo();
    };
    $scope["delete"] = function(type, index) {
      if (type === 0) {
        $scope.configlist.splice(index, 1);
      } else if (type === 1) {
        $scope.petlist.splice(index, 1);
      }
      return setconfiginfo();
    };
    $scope.add = function(type) {
      if (type === 0) {
        $scope.adddata.url = "";
        $scope.adddata.title = "";
        return $scope.addstate = true;
      } else if (type === 1) {
        $scope.petadddata = "";
        return $scope.petaddstate = true;
      }
    };
    $scope.addsave = function(type) {
      var item;
      if (type === 0) {
        $scope.addstate = false;
        $scope.configlist.push($scope.adddata);
      } else if (type === 1) {
        item = {
          "isedit": false,
          "content": $scope.petadddata
        };
        $scope.petlist.push(item);
        $scope.petaddstate = false;
      }
      return setconfiginfo();
    };
    $scope.cancelsave = function(type) {
      if (type === 0) {
        return $scope.addstate = false;
      } else if (type === 1) {
        return $scope.petaddstate = false;
      }
    };
    return getconfiginfo();
  }
]);

var coverController;

coverController = app.controller('coverController', [
  'app', '$scope', '$rootScope', 'articleService', function(app, $scope, $rootScope, articleService) {
    var getArticleList, getInterviewList, getList, getTopicList, pageCount, refreshCover, setFailed, topic_index, uploadComplete, uploadFailed;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    $scope.interviewid = "";
    $scope.previewinfo = {
      title: '',
      interview_id: '',
      cover_image: '',
      article_id: ''
    };
    pageCount = 20;
    topic_index = 0;
    getList = function(type, page_index) {
      var articleinfo;
      articleinfo = articleService.getarticlelist('', type, page_index, pageCount);
      if (type === 'topic') {
        $scope.interviewlist = articleinfo.articlelist;
        return $scope.pagination = articleinfo.pagination;
      } else if (type === 'rec') {
        $scope.articlelist = articleinfo.articlelist;
        return $scope.articlepagination = articleinfo.pagination;
      }
    };
    setFailed = function(reason) {
      return alert("操作不成功，" + reason);
    };
    getInterviewList = function(page_index) {
      return getList('topic', page_index);
    };
    getArticleList = function(page_index) {
      return getList('rec', page_index);
    };
    getTopicList = function(index) {
      var topiclist;
      topiclist = articleService.gettopiclist('', '', index, pageCount);
      topic_index = index;
      $scope.coversettinglist = topiclist.articlelist;
      return $scope.topiclistpagination = topiclist.pagination;
    };
    uploadFailed = function(evt) {
      return $scope.updateError = true;
    };
    refreshCover = function() {
      getInterviewList(0);
      getArticleList(0);
      getTopicList(topic_index);
      $scope.interviewid = '';
      $scope.articleid = '';
      $scope.cover_image = '';
      return $scope.fileUrl = '';
    };
    uploadComplete = function(evt) {
      var jsonData;
      jsonData = JSON.parse(evt.target.response);
      if (jsonData.error.error_id === 0) {
        $scope.fileUrl = jsonData.response_data.fileurl;
        $scope.$apply();
        return $scope.updateError = false;
      } else {
        return $scope.updateError = true;
      }
    };
    $scope.uploadFile = function() {
      return articleService.topicpost($scope.interviewid, $scope.articleid, $scope.fileUrl).then(refreshCover, setFailed);
    };
    $scope.uploadimg = function() {
      var fd, url, xhr;
      if ((document.getElementById('fileToUpload').files[0] != null) && (document.getElementById('fileToUpload').files[0].name != null)) {
        fd = new FormData();
        fd.append("filedata", document.getElementById('fileToUpload').files[0]);
        xhr = new XMLHttpRequest();
        xhr.addEventListener("load", uploadComplete, false);
        xhr.addEventListener("error", uploadFailed, false);
        xhr.addEventListener("abort", uploadFailed, false);
        url = Util.host + '/1/file/upload';
        xhr.open("POST", url);
        return xhr.send(fd);
      }
    };
    $scope.selectinterview = function(interviewid) {
      return $scope.interviewid = interviewid;
    };
    $scope.selectaritcle = function(articleid) {
      return $scope.articleid = articleid;
    };
    $scope["delete"] = function(index) {
      return articleService.deletearticle($scope.coversettinglist[index].article_id).then(refreshCover, setFailed);
    };
    $scope.$on('genPagination', function(event, p, id) {
      event.stopPropagation();
      if (id === 'interviewpage') {
        return getInterviewList(p);
      } else if (id === 'articlepage') {
        return getArticleList(p);
      } else if (id === 'topiclistpage') {
        return getTopicList(p);
      }
    });
    $scope.preview = function() {
      var item, _i, _len, _ref;
      $scope.previewinfo = {
        title: '',
        interview_id: '',
        cover_image: '',
        article_id: ''
      };
      $scope.previewinfo.cover_image = $scope.fileUrl;
      $scope.previewinfo.article_id = $scope.articleid;
      _ref = $scope.interviewlist;
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        item = _ref[_i];
        if (item.article_id === $scope.interviewid) {
          $scope.previewinfo.title = item.cover_text;
          break;
        }
      }
      return $scope.previewinfo.interview_id = $scope.interviewid;
    };
    $scope.previewhistory = function(index) {};
    return refreshCover();
  }
]);

var dashboardController;

dashboardController = app.controller('dashboardController', [
  'app', '$scope', '$routeParams', '$rootScope', 'dashboardService', function(app, $scope, $routeParams, $rootScope, dashboardService) {
    var days, getArticleList, getGameList, getOnlineList, getReportList, getsummaryinfo;
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    days = [1, 7, 15, 30];
    $scope.dropdowmItems = ["1 天", "7 天", "15 天", "30 天"];
    $scope.selectindex = 1;
    $scope.countPageChange = function(index) {
      $scope.selectindex = index;
      return getsummaryinfo(days[$scope.selectindex]);
    };
    getsummaryinfo = function(day) {
      return $scope.summary = dashboardService.getsummaryinfo(day);
    };
    getOnlineList = function(pageIndex) {
      var userlist;
      userlist = dashboardService.getlist('-onlinetime', pageIndex);
      $scope.onlinelist = userlist.userlist;
      return $scope.pagination = userlist.pagination;
    };
    getReportList = function(pageIndex) {
      var userlist;
      userlist = dashboardService.getlist('-record', pageIndex);
      $scope.reportlist = userlist.userlist;
      return $scope.reportpagination = userlist.pagination;
    };
    getArticleList = function(pageIndex) {
      var userlist;
      userlist = dashboardService.getlist('-post', pageIndex);
      $scope.articlelist = userlist.userlist;
      return $scope.articlepagination = userlist.pagination;
    };
    getGameList = function(pageIndex) {
      var userlist;
      userlist = dashboardService.getlist('-gametime', pageIndex);
      $scope.gamelist = userlist.userlist;
      return $scope.gamepagination = userlist.pagination;
    };
    $scope.$on('genPagination', function(event, p, id) {
      event.stopPropagation();
      if (id === 'onlinepage') {
        return getOnlineList(p);
      } else if (id === 'reportpage') {
        return getReportList(p);
      } else if (id === 'articlepage') {
        return getArticleList(p);
      } else if (id === 'gamepage') {
        return getGameList(p);
      }
    });
    getsummaryinfo(days[$scope.selectindex]);
    getOnlineList(0);
    getReportList(0);
    getArticleList(0);
    return getGameList(0);
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
          var data;
          if (checkRequest(retData)) {
            data = {
              isLogin: true,
              access_token: retData.access_token,
              userid: retData.userid
            };
            app.checkUser(data);
            userObj.userid = retData.userid;
            userObj.getInfo(retData.access_token, function(userInfo) {
              if (checkRequest(userInfo)) {
                data.profile = userInfo.profile;
                return app.checkUser(data);
              }
            });
            window.location.href = "#/0";
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
      if ($rootScope.isLogin) {
        return window.location.href = "#/0";
      }
    };
  }
]);

var tasklistController;

tasklistController = app.controller('tasklistController', [
  'app', '$scope', '$rootScope', 'taskService', 'utils', function(app, $scope, $rootScope, taskService, utils) {
    var pageIndex, refreshtable, taskReason, taskfinished, timeline;
    $scope.checked = true;
    $scope.itemsByPage = 50;
    taskfinished = false;
    pageIndex = 0;
    taskReason = {
      "accept": "不错呦，加油！",
      "reject": "您上传的资料有误，请重新检查！"
    };
    $scope.searchData = {
      "data": ""
    };
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    refreshtable = function() {
      if (!$scope.checked) {
        $scope.checked = true;
        taskfinished = fasle;
      }
      if (window.location.href.indexOf('tasklisthistory') > 0) {
        if ($scope.checked) {
          $scope.checked = false;
          taskfinished = true;
        }
      }
      return taskService.gettasklist(taskfinished, pageIndex).then(function(data) {
        $scope.rowCollection = data.tasklist;
        $scope.displayedCollection = [].concat($scope.rowCollection);
        $scope.pagination = data.pagination;
        if (!taskfinished) {
          return $rootScope.checkTaskList(data);
        }
      });
    };
    timeline = function() {
      var tasklistInfo;
      tasklistInfo = taskService.searchtask($scope.searchData.data, !$scope.checked, pageIndex);
      $scope.rowCollection = tasklistInfo.tasklist;
      $scope.displayedCollection = [].concat($scope.rowCollection);
      return $scope.pagination = tasklistInfo.pagination;
    };
    $scope.showImgs = function(index) {
      var str;
      utils.removeItem("task_imgs");
      utils.setItem("task_imgs", $scope.rowCollection[index].images);
      return str = utils.getItem("task_imgs");
    };
    $scope.Approve = function(row) {
      this.reason = row.reason.trim();
      if (this.reason === "") {
        return alert("please input the reason for the rejection");
      } else {
        pageIndex = 0;
        return taskService.taskapprove(row.userid, row.taskid, this.reason).then(refreshtable);
      }
    };
    $scope.Reject = function(row) {
      this.reason = row.reason.trim();
      if (this.reason === "") {
        return alert("please input the reason for the rejection");
      } else {
        pageIndex = 0;
        return taskService.taskreject(row.userid, row.taskid, this.reason).then(refreshtable);
      }
    };
    $scope.searchChange = function() {
      if (($scope.searchData.data != null) && $scope.searchData.data.length > 0) {
        pageIndex = 0;
        return timeline();
      } else {
        return refreshtable();
      }
    };
    $scope.dealAll = function() {
      var authlist, item, itemtmp, _i, _len, _ref;
      authlist = [];
      _ref = $scope.displayedCollection;
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        item = _ref[_i];
        if (item.pass >= 0 && item.reason.length > 0) {
          itemtmp = {
            userid: item.userid,
            task_id: item.taskid,
            reason: item.reason
          };
          if (item.pass === 1) {
            itemtmp.pass = true;
          } else {
            itemtmp.pass = false;
          }
          authlist.push(itemtmp);
        }
      }
      if (authlist.length > 0) {
        return taskService.dealalltask(authlist).then(refreshtable);
      }
    };
    $scope.approveSelect = function(index) {
      $scope.displayedCollection[index].pass = 1;
      return $scope.displayedCollection[index].reason = taskReason.accept;
    };
    $scope.rejectSelect = function(index) {
      $scope.displayedCollection[index].pass = 0;
      return $scope.displayedCollection[index].reason = taskReason.reject;
    };
    $scope.cancelSelect = function(index) {
      $scope.displayedCollection[index].pass = -1;
      if ($scope.displayedCollection[index].reason.length > 0) {
        return $scope.displayedCollection[index].reason = "";
      }
    };
    $scope.changeContent = function(index) {
      if ($scope.displayedCollection[index].pass === 1) {
        return taskReason.accept = $scope.displayedCollection[index].reason;
      } else if ($scope.displayedCollection[index].pass === 0) {
        return taskReason.reject = $scope.displayedCollection[index].reason;
      }
    };
    $scope.$on('genPagination', function(event, p) {
      event.stopPropagation();
      pageIndex = p;
      if (($scope.searchData.data != null) && $scope.searchData.data.length > 0) {
        return timeline();
      } else {
        return refreshtable();
      }
    });
    return refreshtable();
  }
]);

var userdetailController;

userdetailController = app.controller('userdetailController', [
  'app', '$scope', '$routeParams', '$rootScope', 'utils', 'dashboardService', function(app, $scope, $routeParams, $rootScope, utils, dashboardService) {
    var getArrayString, getArticleByUser, getFollowers, getFollows, getSportByUser, userinfoId;
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
    $scope.banUser = function(nState) {
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
    getArticleByUser = function(pageIndex) {
      var listdata;
      listdata = dashboardService.getaritlclebyuser(userinfoId, pageIndex);
      $scope.articlesList = listdata.articlelist;
      return $scope.pagination = listdata.pagination;
    };
    getSportByUser = function(pageIndex) {
      var listdata;
      listdata = dashboardService.getsportlistbyuser(userinfoId, pageIndex, 'run');
      $scope.sportList = listdata.sportlist;
      return $scope.sportpagination = listdata.pagination;
    };
    getFollows = function(pageIndex) {
      var listdata;
      listdata = dashboardService.getfriendshiplistbyuser(userinfoId, pageIndex, 'follows');
      $scope.followslist = listdata.friendshiplist;
      return $scope.followspagination = listdata.pagination;
    };
    getFollowers = function(pageIndex) {
      var listdata;
      listdata = dashboardService.getfriendshiplistbyuser(userinfoId, pageIndex, 'followers');
      $scope.followerslist = listdata.friendshiplist;
      return $scope.followerspagination = listdata.pagination;
    };
    $scope.showImgs = function(index) {
      var str;
      utils.removeItem("task_imgs");
      utils.setItem("task_imgs", $scope.rowCollection[index].images);
      return str = utils.getItem("task_imgs");
    };
    $scope.$on('genPagination', function(event, p, id) {
      event.stopPropagation();
      if (id === 'articlepage') {
        return getArticleByUser(p);
      } else if (id === 'sportpage') {
        return getSportByUser(p);
      } else if (id === 'followerspage') {
        return getFollowers(p);
      } else if (id === 'followspage') {
        return getFollows(p);
      }
    });
    getArticleByUser(0);
    getSportByUser(0);
    getFollows(0);
    return getFollowers(0);
  }
]);

var userlistController,
  __indexOf = [].indexOf || function(item) { for (var i = 0, l = this.length; i < l; i++) { if (i in this && this[i] === item) return i; } return -1; };

userlistController = app.controller('userlistController', [
  'app', '$scope', '$rootScope', '$routeParams', function(app, $scope, $rootScope, $routeParams) {
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
    $scope.dropdowmItems = ["50项", "100项", "200项"];
    $scope.selectType = ["选择", "性别", "年龄", "状态", "账号类型", "角色"];
    $scope.selectItemList = [["选择"], ["男", "女"], ["< 20岁", "20～40岁", "> 40岁"], ["正常", "禁言", "拉黑"], ["手机", "微博", "邮箱"], ["管理员", "教练", "普通用户"]];
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
      "gender": 0,
      "email": 0
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
        } else if ($scope.typeIndex === 4) {
          if ($scope.filtItemIndex === 0) {
            searchDetail.role = "phone";
          } else if ($scope.filtItemIndex === 1) {
            searchDetail.role = "weibo";
          } else {
            searchDetail.role = "email";
          }
        } else if ($scope.typeIndex === 5) {
          if ($scope.filtItemIndex === 0) {
            searchDetail.actor = "admin";
          } else if ($scope.filtItemIndex === 1) {
            searchDetail.actor = "coach";
          } else {
            searchDetail.actor = "user";
          }
        }
      } else {
        searchDetail.gender = "";
        searchDetail.age = "";
        searchDetail.ban_status = "";
      }
      return User.search(userToken, searchDetail.keyword, searchDetail.gender, searchDetail.age, searchDetail.ban_status, searchDetail.role, searchDetail.actor, sortStr, function(retData, page_index, page_total, total_count) {
        var useritem, _i, _j, _len, _results;
        if (checkRequest(retData)) {
          $scope.total_num = total_count;
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
          $scope.total_num = total_count;
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
          "gender": 0,
          "email": 0
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

app.factory('authq', [
  '$http', function($http) {
    return {
      getauthlist: function(sort, page_index, page_count) {
        return $http.get(Util.host + "/admin/user/auth/list", {
          params: {
            sort: sort,
            page_index: page_index,
            page_count: page_count,
            access_token: userToken
          }
        });
      },
      authactionpost: function(userid, auth_type, auth_status, authreview) {
        return $http.post(Util.host + "/admin/user/auth", {
          userid: userid,
          auth_type: auth_type,
          auth_status: auth_status,
          auth_review: authreview,
          access_token: userToken
        });
      },
      authinfo: function(userid) {
        return $http.get(Util.host + "/admin/user/info", {
          params: {
            userid: userid,
            access_token: userToken
          }
        });
      }
    };
  }
]);

app.factory('dashboardq', [
  '$http', function($http) {
    return {
      getsummary: function(day) {
        return $http.get(Util.host + "/admin/stat/summary", {
          params: {
            days: day,
            access_token: userToken
          }
        });
      },
      getuserlist: function(sort, pageIndex, pageCount) {
        if (pageIndex == null) {
          pageIndex = 0;
        }
        if (pageCount == null) {
          pageCount = 20;
        }
        return $http.get(Util.host + "/admin/user/list", {
          params: {
            sort: sort,
            access_token: userToken,
            page_index: pageIndex,
            page_count: pageCount
          }
        });
      },
      getsportlistbyuser: function(userId, type, pageIndex, pageCount) {
        if (pageIndex == null) {
          pageIndex = 0;
        }
        if (pageCount == null) {
          pageCount = 20;
        }
        return $http.get(Util.host + "/admin/record/timeline", {
          params: {
            userid: userId,
            type: type,
            access_token: userToken,
            page_index: pageIndex,
            page_count: pageCount
          }
        });
      },
      getfriendshipbyuser: function(userId, type, pageIndex, pageCount) {
        if (pageIndex == null) {
          pageIndex = 0;
        }
        if (pageCount == null) {
          pageCount = 20;
        }
        return $http.get(Util.host + "/admin/user/friendship", {
          params: {
            userid: userId,
            type: type,
            access_token: userToken,
            page_index: pageIndex,
            page_count: pageCount
          }
        });
      },
      getarticlebyuser: function(userid, pageIndex, pageCount) {
        if (pageIndex == null) {
          pageIndex = 0;
        }
        if (pageCount == null) {
          pageCount = 50;
        }
        return $http.get(Util.host + "/admin/article/timeline", {
          params: {
            userid: userid,
            page_index: pageIndex,
            page_count: pageCount,
            access_token: userToken
          }
        });
      }
    };
  }
]);

app.factory('configq', [
  '$http', function($http) {
    return {
      getconfig: function() {
        return $http.get(Util.host + "/admin/config/get", {
          params: {
            access_token: userToken
          }
        });
      },
      setconfig: function(configinfo) {
        return $http.post(Util.host + '/admin/config/set', {
          config: configinfo,
          access_token: userToken
        });
      }
    };
  }
]);

app.factory('taskService', [
  '$q', 'taskq', function($q, $taskq) {
    return {
      gettasklist: function(finished, page_index) {
        var deferred, tasklistInfo;
        deferred = $q.defer();
        tasklistInfo = {
          'tasklist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $taskq.gettasklist(page_index, finished).success(function(response) {
          var task, taskitem, taskjson, _i, _j, _k, _len, _len1, _ref, _ref1, _ref2, _results;
          if (checkRequest(response)) {
            _ref = response.users;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              taskitem = _ref[_i];
              _ref1 = taskitem.tasks;
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
                  profile: taskitem.profile,
                  begin_time: task.begin_time,
                  end_time: task.end_time,
                  distance: task.distance,
                  source: task.source,
                  duration: task.duration,
                  pass: -1
                };
                tasklistInfo.tasklist.push(taskjson);
              }
            }
            tasklistInfo.pagination.pageIndex = response.page_index;
            tasklistInfo.pagination.total = response.total_number;
            tasklistInfo.pagination.showPages = (function() {
              _results = [];
              for (var _k = 0, _ref2 = response.page_total; 0 <= _ref2 ? _k < _ref2 : _k > _ref2; 0 <= _ref2 ? _k++ : _k--){ _results.push(_k); }
              return _results;
            }).apply(this);
            tasklistInfo.pagination.pagetotal = response.page_total;
            return deferred.resolve(tasklistInfo);
          } else {
            return deferred.reject(response);
          }
        });
        return deferred.promise;
      },
      taskapprove: function(userid, taskid, reason) {
        return $taskq.taskaudit(userid, taskid, true, reason).success();
      },
      taskreject: function(userid, taskid, reason) {
        return $taskq.taskaudit(userid, taskid, false, reason).success();
      },
      dealalltask: function(authlist) {
        return $taskq.taskauditall(authlist).success();
      },
      searchtask: function(nickname, finish, page_index) {
        var tasklistInfo;
        tasklistInfo = {
          'tasklist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $taskq.search(nickname, finish, 50, page_index).success(function(response) {
          var task, taskitem, taskjson, _i, _j, _k, _len, _len1, _ref, _ref1, _ref2, _results;
          if (response.users != null) {
            _ref = response.users;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              taskitem = _ref[_i];
              _ref1 = taskitem.tasks;
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
                  profile: taskitem.profile,
                  begin_time: task.begin_time,
                  end_time: task.end_time,
                  distance: task.distance
                };
                if (finish) {
                  if (taskjson.status === "FINISH" || taskjson.status === "UNFINISH") {
                    tasklistInfo.tasklist.push(taskjson);
                  }
                } else {
                  if (taskjson.status === "AUTHENTICATION") {
                    tasklistInfo.tasklist.push(taskjson);
                  }
                }
              }
            }
            tasklistInfo.pagination.pageIndex = response.page_index;
            tasklistInfo.pagination.total = response.total_number;
            tasklistInfo.pagination.showPages = (function() {
              _results = [];
              for (var _k = 0, _ref2 = response.page_total; 0 <= _ref2 ? _k < _ref2 : _k > _ref2; 0 <= _ref2 ? _k++ : _k--){ _results.push(_k); }
              return _results;
            }).apply(this);
            return tasklistInfo.pagination.pagetotal = response.page_total;
          }
        });
        return tasklistInfo;
      }
    };
  }
]);

app.factory('authService', [
  '$q', 'authq', function($q, $authq) {
    return {
      getauthlist: function(sort, page_index, page_count) {
        var authlistinfo, deferred;
        deferred = $q.defer();
        authlistinfo = {
          'authlist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $authq.getauthlist(sort, page_index, page_count).success(function(response) {
          var item, _i, _j, _len, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.users;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              authlistinfo.authlist.push(item);
            }
            authlistinfo.pagination.pageIndex = response.page_index;
            authlistinfo.pagination.total = response.total_number;
            authlistinfo.pagination.showPages = (function() {
              _results = [];
              for (var _j = 0, _ref1 = response.page_total; 0 <= _ref1 ? _j < _ref1 : _j > _ref1; 0 <= _ref1 ? _j++ : _j--){ _results.push(_j); }
              return _results;
            }).apply(this);
            authlistinfo.pagination.pagetotal = response.page_total;
            return deferred.resolve(authlistinfo);
          } else {
            return deferred.reject(response);
          }
        });
        return deferred.promise;
      },
      authpost: function(userid, auth_type, auth_status, authreview) {
        return $authq.authactionpost(userid, auth_type, auth_status, authreview).success(function(response) {
          return console.log(response);
        });
      },
      authinfo: function(userid) {
        return $authq.authinfo(userid);
      }
    };
  }
]);

app.factory('dashboardService', [
  '$q', 'dashboardq', function($q, $dashboardq) {
    return {
      getsummaryinfo: function(day) {
        var summary;
        summary = {
          "summarylist": [],
          "users": 0,
          "onlines": 0,
          "online_coaches": 0
        };
        $dashboardq.getsummary(day).success(function(response) {
          var i, item, items, _i, _j, _len, _results;
          if (checkRequest(response)) {
            items = (function() {
              _results = [];
              for (var _i = 0; 0 <= day ? _i < day : _i > day; 0 <= day ? _i++ : _i--){ _results.push(_i); }
              return _results;
            }).apply(this);
            for (_j = 0, _len = items.length; _j < _len; _j++) {
              i = items[_j];
              item = {};
              item.reg_phone = response.reg_phone[i];
              item.reg_email = response.reg_email[i];
              item.reg_weibo = response.reg_weibo[i];
              item.logins = response.logins[i];
              item.actives = response.actives[i];
              item.post_users = response.post_users[i];
              item.posts = response.posts[i];
              item.gamers = response.gamers[i];
              item.game_time = response.game_time[i];
              item.record_users = response.record_users[i];
              item.auth_coaches = response.auth_coaches[i];
              item.coach_logins = response.coach_logins[i];
              item.coins = response.coins[i];
              item.dateTime = (new Date()).valueOf() - 24 * 60 * 60 * 1000 * i;
              summary.summarylist.push(item);
            }
            summary.users = response.users;
            summary.onlines = response.onlines;
            return summary.online_coaches = response.online_coaches;
          }
        });
        return summary;
      },
      getlist: function(sort, pageIndex) {
        var listdata;
        listdata = {
          'userlist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $dashboardq.getuserlist(sort, pageIndex, 10).success(function(response) {
          var item, userItem, _i, _j, _len, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.users;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              userItem = {};
              userItem.userid = item.userid;
              userItem.nickname = item.nickname;
              userItem.profile = item.profile;
              userItem.gender = item.gender;
              userItem.birthday = item.birthday;
              if (item.stat !== null && item.stat !== void 0) {
                userItem.onlinetime = item.stat.onlinetime;
                userItem.report = item.stat.records;
                userItem.article = item.stat.articles;
                userItem.game = item.stat.gametime;
              } else {
                userItem.onlinetime = 0;
                userItem.report = 0;
                userItem.article = 0;
                userItem.game = 0;
              }
              listdata.userlist.push(userItem);
            }
            listdata.pagination.pageIndex = response.page_index;
            listdata.pagination.total = response.total_number;
            listdata.pagination.showPages = (function() {
              _results = [];
              for (var _j = 0, _ref1 = response.page_total; 0 <= _ref1 ? _j < _ref1 : _j > _ref1; 0 <= _ref1 ? _j++ : _j--){ _results.push(_j); }
              return _results;
            }).apply(this);
            return listdata.pagination.pagetotal = response.page_total;
          }
        });
        return listdata;
      },
      getsportlistbyuser: function(userid, pageIndex, type) {
        var listdata;
        listdata = {
          'sportlist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $dashboardq.getsportlistbyuser(userid, type, pageIndex).success(function(response) {
          var item, _i, _j, _len, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.records;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              listdata.sportlist.push(item);
            }
            listdata.pagination.pageIndex = response.page_index;
            listdata.pagination.total = response.total_number;
            listdata.pagination.showPages = (function() {
              _results = [];
              for (var _j = 0, _ref1 = response.page_total; 0 <= _ref1 ? _j < _ref1 : _j > _ref1; 0 <= _ref1 ? _j++ : _j--){ _results.push(_j); }
              return _results;
            }).apply(this);
            return listdata.pagination.pagetotal = response.page_total;
          }
        });
        return listdata;
      },
      getfriendshiplistbyuser: function(userid, pageIndex, type) {
        var listdata;
        listdata = {
          'friendshiplist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $dashboardq.getfriendshipbyuser(userid, type, pageIndex).success(function(response) {
          var item, _i, _j, _len, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.users;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              listdata.friendshiplist.push(item);
            }
            listdata.pagination.pageIndex = response.page_index;
            listdata.pagination.total = response.total_number;
            listdata.pagination.showPages = (function() {
              _results = [];
              for (var _j = 0, _ref1 = response.page_total; 0 <= _ref1 ? _j < _ref1 : _j > _ref1; 0 <= _ref1 ? _j++ : _j--){ _results.push(_j); }
              return _results;
            }).apply(this);
            return listdata.pagination.pagetotal = response.page_total;
          }
        });
        return listdata;
      },
      getaritlclebyuser: function(userId, pageIndex) {
        var listdata;
        listdata = {
          'articlelist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $dashboardq.getarticlebyuser(userId, pageIndex).success(function(response) {
          var item, _i, _j, _len, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.articles;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              listdata.articlelist.push(item);
            }
            listdata.pagination.pageIndex = response.page_index;
            listdata.pagination.total = response.total_number;
            listdata.pagination.showPages = (function() {
              _results = [];
              for (var _j = 0, _ref1 = response.page_total; 0 <= _ref1 ? _j < _ref1 : _j > _ref1; 0 <= _ref1 ? _j++ : _j--){ _results.push(_j); }
              return _results;
            }).apply(this);
            return listdata.pagination.pagetotal = response.page_total;
          }
        });
        return listdata;
      }
    };
  }
]);

app.factory('configService', [
  '$q', 'configq', function($q, $configq) {
    return {
      getconfig: function() {
        var retList;
        retList = {
          "videos": [],
          "pets": []
        };
        $configq.getconfig().success(function(response) {
          var item, itemdata, _i, _j, _len, _len1, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.videos;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              item.isedit = false;
              retList.videos.push(item);
            }
            _ref1 = response.pets;
            _results = [];
            for (_j = 0, _len1 = _ref1.length; _j < _len1; _j++) {
              item = _ref1[_j];
              itemdata = {
                "content": item,
                "isedit": false
              };
              _results.push(retList.pets.push(itemdata));
            }
            return _results;
          }
        });
        return retList;
      },
      setconfig: function(info) {
        return $configq.setconfig(info).success();
      }
    };
  }
]);

app.factory('articleService', [
  '$q', 'articleq', function($q, $articleq) {
    return {
      getarticlelist: function(sort, tag, page_index, page_count) {
        var articleinfo;
        articleinfo = {
          'articlelist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $articleq.getarticlelist(sort, tag, page_index, page_count).success(function(response) {
          var item, _i, _j, _len, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.articles;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              item.type = 0;
              if (item.tags != null) {
                if (item.tags.indexOf('rec') !== -1) {
                  item.type = 2;
                } else if (item.tags.indexOf('topic') !== -1) {
                  item.type = 1;
                }
              }
              articleinfo.articlelist.push(item);
            }
            articleinfo.pagination.pageIndex = response.page_index;
            articleinfo.pagination.total = response.total_number;
            articleinfo.pagination.showPages = (function() {
              _results = [];
              for (var _j = 0, _ref1 = response.page_total; 0 <= _ref1 ? _j < _ref1 : _j > _ref1; 0 <= _ref1 ? _j++ : _j--){ _results.push(_j); }
              return _results;
            }).apply(this);
            return articleinfo.pagination.pagetotal = response.page_total;
          }
        });
        return articleinfo;
      },
      searcharticle: function(keyword, tag, sort, pageIndex, pageCount) {
        var articleinfo;
        articleinfo = {
          'articlelist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $articleq.searcharticle(keyword, tag, sort, pageIndex, pageCount).success(function(response) {
          var item, _i, _j, _len, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.articles;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              if (item.tags.indexOf('rec') !== -1) {
                item.type = 2;
              } else if (item.tags.indexOf('topic') !== -1) {
                item.type = 1;
              } else {
                item.type = 0;
              }
              articleinfo.articlelist.push(item);
            }
            articleinfo.pagination.pageIndex = response.page_index;
            articleinfo.pagination.total = response.total_number;
            articleinfo.pagination.showPages = (function() {
              _results = [];
              for (var _j = 0, _ref1 = response.page_total; 0 <= _ref1 ? _j < _ref1 : _j > _ref1; 0 <= _ref1 ? _j++ : _j--){ _results.push(_j); }
              return _results;
            }).apply(this);
            return articleinfo.pagination.pagetotal = response.page_total;
          }
        });
        return articleinfo;
      },
      articlepost: function(articleId, title, author, imglist, contents, tag) {
        var deferred;
        deferred = $q.defer();
        $articleq.articlepost(articleId, title, author, imglist, contents, tag).success(function(response) {
          if (checkRequest(response)) {
            return deferred.resolve(response);
          } else {
            return deferred.reject(response.error_desc);
          }
        });
        return deferred.promise;
      },
      getimagelist: function(contents) {
        var elem, imglist;
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
      },
      articlemark: function(article_id, type) {
        var deferred;
        deferred = $q.defer();
        $articleq.articlemark(article_id, type).success(function(response) {
          if (checkRequest(response)) {
            return deferred.resolve(response);
          } else {
            return deferred.reject(response.error_desc);
          }
        });
        return deferred.promise;
      },
      topicpost: function(topic, rec, image) {
        var deferred;
        deferred = $q.defer();
        $articleq.topicpost(topic, rec, image).success(function(response) {
          if (checkRequest(response)) {
            return deferred.resolve(response);
          } else {
            return deferred.reject(response.error_desc);
          }
        });
        return deferred.promise;
      },
      gettopiclist: function(sort, tag, page_index, page_count) {
        var topiclist;
        topiclist = {
          'articlelist': [],
          'pagination': {
            'total': 0,
            'pageIndex': 0,
            'pagetotal': 0,
            'showPages': []
          }
        };
        $articleq.gettopiclist(sort, tag, page_index, page_count).success(function(response) {
          var item, _i, _j, _len, _ref, _ref1, _results;
          if (checkRequest(response)) {
            _ref = response.articles;
            for (_i = 0, _len = _ref.length; _i < _len; _i++) {
              item = _ref[_i];
              topiclist.articlelist.push(item);
            }
            topiclist.pagination.pageIndex = response.page_index;
            topiclist.pagination.total = response.total_number;
            topiclist.pagination.showPages = (function() {
              _results = [];
              for (var _j = 0, _ref1 = response.page_total; 0 <= _ref1 ? _j < _ref1 : _j > _ref1; 0 <= _ref1 ? _j++ : _j--){ _results.push(_j); }
              return _results;
            }).apply(this);
            return topiclist.pagination.pagetotal = response.page_total;
          }
        });
        return topiclist;
      },
      deletearticle: function(article_id) {
        var deferred;
        deferred = $q.defer();
        $articleq.deletearticle(article_id).success(function(response) {
          if (checkRequest(response)) {
            return deferred.resolve(response);
          } else {
            return deferred.reject(response.error_desc);
          }
        });
        return deferred.promise;
      },
      getarticleinfo: function(article_id) {
        var deferred;
        deferred = $q.defer();
        $articleq.getarticleinfo(article_id).success(function(response) {
          if (checkRequest(response)) {
            return deferred.resolve(response);
          } else {
            return deferred.reject(response.error_desc);
          }
        });
        return deferred.promise;
      }
    };
  }
]);
