var app, httphost, taskmenulist, taskreason;

app = angular.module('app', ['ngRoute', 'ngMaterial']);

httphost = "http://172.24.222.54:8080";

taskreason = {
  accept: "不错呦，加油！",
  reject: "您上传的资料有误，请重新检查！"
};

taskmenulist = ["待审批", "已审批"];

app.factory('taskq', [
  '$http', function($http) {
    return {
      gettasklist: function(page_index) {
        return $http.get(httphost + '/admin/task/list', {
          params: {
            page_index: page_index,
            page_count: 50,
            access_token: userToken
          }
        });
      },
      taskaudit: function(userid, taskid, pass, reason) {
        return $http.post(httphost + '/admin/task/auth', {
          userid: userid,
          task_id: taskid,
          pass: pass,
          reason: reason,
          access_token: userToken
        });
      },
      search: function(nickname, finish, page_count, page_index) {
        return $http.get(httphost + '/admin/task/timeline', {
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

app.factory('userreq', [
  '$http', function($http) {
    return {
      userlogin: function(userid, password) {
        return $http.post(httphost + '/admin/login', {
          username: userid,
          password: password
        });
      },
      userlogout: function() {
        return $http.post(httphost + '/admin/logout', {
          access_token: userToken
        });
      },
      getuserinfo: function(userId) {
        return $http.get(httphost + '/admin/user/info', {
          params: {
            userid: userId,
            access_token: userToken
          }
        });
      }
    };
  }
]);

var checkDate, checkRequest, urlPath, userId, userToken;

app.constant('app', {
  version: Date.now()
});

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
  var login, tasklist;
  login = {
    templateUrl: 'html/user-login.html',
    controller: 'loginController'
  };
  tasklist = {
    templateUrl: 'html/task-list.html',
    controller: 'tasklistController'
  };
  return $routeProvider.when('/', login).when('/task', tasklist).when('/task/:id', tasklist);
});

app.run([
  'app', '$rootScope', 'utils', '$filter', 'userService', '$mdDialog', function(app, $rootScope, utils, $filter, userService, $mdDialog) {
    $rootScope.isLogin = utils.getItem('isLogin');
    $rootScope.profile = utils.getItem('profile');
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
    $rootScope.logout = function(ev) {
      var confirm;
      confirm = $mdDialog.confirm().parent(angular.element(document.body)).title('悦动力').content('你确定要退出吗？').ariaLabel('Lucky day').ok('确认').cancel('取消').targetEvent(ev);
      return $mdDialog.show(confirm).then(function() {
        var data;
        userService.userlogout();
        data = {
          isLogin: false,
          userid: '',
          access_token: '',
          profile: ''
        };
        app.checkUser(data);
        return window.location.href = "#/";
      }, function() {});
    };
    return app.rootScope = $rootScope;
  }
]);

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
    imagePath = "../images/lanhan.png";
    if (angular.isString(data)) {
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
}).filter("taskSource", function() {
  return function(data) {
    var taskSource;
    taskSource = "手动输入";
    if (angular.isString(data) && data.length > 0) {
      taskSource = "自动导入";
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
}).filter("distance", function() {
  return function(data) {
    return data / 1000;
  };
}).filter("descfilter", function() {
  return function(data) {
    var desc;
    desc = "（您还未填写运动心情）";
    if (angular.isString(data) && data.length > 0) {
      if (data.length > 24) {
        desc = data.substr(0, 24) + "..";
      } else {
        desc = data;
      }
    }
    return desc;
  };
}).filter("descfilterdetail", function() {
  return function(data) {
    var desc;
    desc = "（您还未填写运动心情）";
    if (angular.isString(data) && data.length > 0) {
      desc = data;
    }
    return desc;
  };
});

app.factory('utils', function() {
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

app.factory('taskService', [
  '$q', 'taskq', function($q, $taskq) {
    return {
      gettasklist: function(tasktype, page_index) {
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
        $taskq.gettasklist(page_index).success(function(response, status, headers, config) {
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
                  distance: task.distance
                };
                switch (tasktype) {
                  case "Auditting":
                    if (taskjson.status === "AUTHENTICATION") {
                      tasklistInfo.tasklist.push(taskjson);
                    }
                    break;
                  case "Audited":
                    if (taskjson.status === "FINISH" || taskjson.status === "UNFINISH") {
                      tasklistInfo.tasklist.push(taskjson);
                    }
                    break;
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
            tasklistInfo.pagination.pagetotal = response.page_total;
            return deferred.resolve(tasklistInfo);
          } else {
            return deferred.reject(response);
          }
        });
        return deferred.promise;
      },
      taskapprove: function(userid, taskid, reason) {
        var deferred;
        deferred = $q.defer();
        $taskq.taskaudit(userid, taskid, true, reason).success(function(response, status, headers, config) {
          if (checkRequest(response)) {
            return deferred.resolve(response);
          } else {
            return deferred.reject(response);
          }
        });
        return deferred.promise;
      },
      taskreject: function(userid, taskid, reason) {
        var deferred;
        deferred = $q.defer();
        $taskq.taskaudit(userid, taskid, false, reason).success(function(response, status, headers, config) {
          if (checkRequest(response)) {
            return deferred.resolve(response);
          } else {
            return deferred.reject(response);
          }
        });
        return deferred.promise;
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
                  distance: task.distance,
                  showdetail: false,
                  chatcontent: ""
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

app.factory('userService', [
  '$q', 'userreq', function($q, $userreq) {
    return {
      userlogin: function(userid, password) {
        var deferred, userData;
        deferred = $q.defer();
        userData = {
          'isLogin': false,
          'id': "",
          'access_token': "",
          'profile': ""
        };
        $userreq.userlogin(userid, password).success(function(data, status, headers, config) {
          if (checkRequest(data)) {
            userData.isLogin = true;
            userData.id = data.userid;
            userData.access_token = data.access_token;
            return deferred.resolve(userData);
          } else {
            return deferred.reject(data);
          }
        }).error(function(data, status, headers, config) {
          return deferred.reject(data);
        });
        return deferred.promise;
      },
      userlogout: function(userid) {
        var deferred;
        deferred = $q.defer();
        $userreq.userlogout(userid).success(function(data, status, headers, config) {
          if (checkRequest(data)) {
            return deferred.resolve(data);
          } else {
            return deferred.reject(data);
          }
        }).error(function(data, status, headers, config) {
          return deferred.reject(data);
        });
        return deferred.promise;
      },
      getuserinfo: function(userId) {
        var deferred;
        deferred = $q.defer();
        $userreq.getuserinfo(userid).success(function(data, status, headers, config) {
          if (checkRequest(data)) {
            return deferred.resolve(data);
          } else {
            return deferred.reject(data);
          }
        }).error(function(data, status, headers, config) {
          return deferred.reject(data);
        });
        return deferred.promise;
      }
    };
  }
]);

var loginController;

loginController = app.controller('loginController', [
  'app', '$scope', '$routeParams', '$rootScope', '$mdDialog', 'userService', function(app, $scope, $routeParams, $rootScope, $mdDialog, userService) {
    var checkLogin;
    $scope.alert = '';
    $rootScope.apptile = "悦动力（登录）";
    $scope.onLogin = function() {
      var promise;
      if (($scope.username != null) && ($scope.pwd != null)) {
        promise = userService.userlogin($scope.username, $scope.pwd);
        return promise.then(function(data) {
          var dataTmp, promiseuserinfo;
          dataTmp = {
            isLogin: true,
            access_token: retData.access_token,
            userid: retData.userid
          };
          app.checkUser(dataTmp);
          promiseuserinfo = userService.getuserinfo(retData.userid).then(function(userInfo) {
            dataTmp.profile = userInfo.profile;
            return app.checkUser(dataTmp);
          });
          $scope.alert = '';
          window.location.href = "#/task";
          $scope.username = "";
          return $scope.pwd = "";
        }, function(data) {
          return $scope.alert = '抱歉! 您的输入有误，请重新检查。';
        });
      }
    };
    $scope.enterLogin = function() {
      var event;
      event = window.event || arguments.callee.caller["arguments"][0];
      if (event.keyCode === 13) {
        return $scope.onLogin();
      }
    };
    checkLogin = function() {
      if ($rootScope.isLogin) {
        return window.location.href = "#/task";
      }
    };
    return checkLogin();
  }
]);

var tasklistController;

tasklistController = app.controller('tasklistController', [
  'app', '$scope', '$rootScope', 'taskService', 'utils', '$mdSidenav', '$routeParams', '$mdDialog', function(app, $scope, $rootScope, taskService, utils, $mdSidenav, $routeParams, $mdDialog) {
    var Approve, DialogController, Load, Reject, nDealIndex, pageIndex, refreshtable, taskStr, totalPage;
    $(function() {
      return $(window).scroll(function() {
        return Load();
      });
    });
    if (!app.getCookie("isLogin")) {
      window.location.href = "#/";
      return;
    }
    Load = function() {
      var documentHeight, loadHeight, scrollHight, windowHeight;
      if (!$scope.showwaiting) {
        loadHeight = 0;
        documentHeight = parseInt($(document).height(), 10);
        windowHeight = parseInt($(window).height(), 10);
        scrollHight = parseInt($(window).scrollTop(), 10);
        if (documentHeight - scrollHight - windowHeight <= loadHeight && pageIndex < totalPage - 1) {
          $scope.showwaiting = true;
          pageIndex++;
          refreshtable();
        }
        if (scrollHight <= 0 && pageIndex > 0) {
          return pageIndex--;
        }
      }
    };
    nDealIndex = 0;
    taskStr = 'Auditting';
    pageIndex = 0;
    totalPage = 0;
    $scope.taskid = $routeParams.id;
    $rootScope.taskmenu = taskmenulist;
    $rootScope.apptile = "运动纪录管理（待审批）";
    $scope.chatlist = [
      {
        "self": false,
        "content": "坚持一下"
      }, {
        "self": true,
        "content": "比较累啊"
      }, {
        "self": false,
        "content": "加油，稍作休整"
      }
    ];
    $scope.toggleRight = function(index) {
      nDealIndex = index;
      return $mdSidenav('right').toggle();
    };
    $scope.close = function() {
      $mdSidenav('right').close();
      return $scope.submitData = "";
    };
    $scope.reasondata = "";
    $scope.submitData = "";
    Approve = function() {
      return taskService.taskapprove($scope.taskList[nDealIndex].userid, $scope.taskList[nDealIndex].taskid, $scope.reasondata).then(refreshtable);
    };
    Reject = function() {
      return taskService.taskreject($scope.taskList[nDealIndex].userid, $scope.taskList[nDealIndex].taskid, $scope.reasondata).then(refreshtable);
    };
    refreshtable = function() {
      if ($scope.taskid == null) {
        return taskService.gettasklist(taskStr, pageIndex).then(function(data) {
          $scope.showwaiting = false;
          $scope.taskList = data.tasklist;
          pageIndex = data.pagination.pageIndex;
          totalPage = data.pagination.pagetotal;
          return document.getElementsByTagName('body')[0].scrollTop = 10;
        });
      }
    };
    $scope.submit = function() {
      if ($scope.submitData === "accept") {
        Approve();
      } else {
        Reject();
      }
      return $scope.close();
    };
    $scope.$watch("submitData", function(newData, oldData) {
      if ($scope.submitData === "accept") {
        return $scope.reasondata = taskreason.accept;
      } else if ($scope.submitData === "reject") {
        return $scope.reasondata = taskreason.reject;
      } else {
        return $scope.reasondata = "";
      }
    });
    $scope.showdetail = function(index) {
      return $scope.taskList[index].showdetail = !$scope.taskList[index].showdetail;
    };
    $scope.sendchat = function(index) {
      var item;
      if ($scope.taskList[index].chatcontent.length > 0) {
        item = {
          "self": true,
          "content": $scope.taskList[index].chatcontent
        };
        $scope.chatlist.push(item);
        return {
          "content": $scope.taskList[index].chatcontent = ""
        };
      }
    };
    $scope.sendchatbykey = function(index) {
      var event;
      event = window.event || arguments.callee.caller["arguments"][0];
      if (event.keyCode === 13) {
        return $scope.sendchat(index);
      }
    };
    $scope.showImgDetail = function(ev, path) {
      $scope.imgPath = path;
      return $mdDialog.show({
        controller: DialogController,
        templateUrl: 'html/task-image-detail.html',
        parent: angular.element(document.body),
        targetEvent: ev,
        clickOutsideToClose: true,
        locals: {
          imgPath: $scope.imgPath
        }
      }).then(function() {
        return $mdDialog.hide();
      });
    };
    $rootScope.taskmenuClick = function(index) {
      if (index === 0 && taskStr !== 'Auditting') {
        taskStr = 'Auditting';
        refreshtable();
        return $rootScope.apptile = "运动纪录管理（待审批）";
      } else if (index === 1 && taskStr !== 'Audited') {
        taskStr = 'Audited';
        refreshtable();
        return $rootScope.apptile = "运动纪录管理（已审批）";
      }
    };
    DialogController = function($scope, $mdDialog, imgPath) {
      $scope.imgPath = imgPath;
      return $scope.closeDialog = function() {
        return $mdDialog.hide();
      };
    };
    return refreshtable();
  }
]);
