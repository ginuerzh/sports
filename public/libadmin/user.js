// Generated by CoffeeScript 1.8.0
(function() {
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

}).call(this);

//# sourceMappingURL=user.js.map