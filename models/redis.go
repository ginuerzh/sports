// log
package models

import (
	//"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	//"strconv"
	//"encoding/json"
	//"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

const (
	redisPrefix             = "sports"
	redisStatVisitorPrefix  = redisPrefix + ":stat:visitors:"  // set per day
	redisStatPvPrefix       = redisPrefix + ":stat:pv:"        // sorted set per day
	redisStatRegisterPrefix = redisPrefix + ":stat:registers:" // set per day

	redisUserOnlinesPrefix    = redisPrefix + ":user:onlines:"      // set per half an hour
	redisUserOnlineUserPrefix = redisPrefix + ":user:online:"       // hashs per user
	redisUserGuest            = redisPrefix + ":user:guest"         // hashes for all guests
	redisUserMessagePrefix    = redisPrefix + ":user:msgs:"         // list per user
	redisUserFollowPrefix     = redisPrefix + ":user:follow:"       // set per user
	redisUserFollowerPrefix   = redisPrefix + ":user:follower:"     // set per user
	redisUserBlacklistPrefix  = redisPrefix + ":user:blacklist:"    // set per user
	redisUserWBImportPrefix   = redisPrefix + ":user:import:weibo:" // set per user
	redisUserGroupPrefix      = redisPrefix + ":user:group:"        // hash per user
	redisGroupPrefix          = redisPrefix + ":group:"             // set per group

	redisStatArticleViewPrefix = redisPrefix + ":stat:articles:view:"  // sorted set per day
	redisStatArticleView       = redisPrefix + ":stat:articles:view"   // sorted set
	redisStatArticleReview     = redisPrefix + ":stat:articles:review" // sorted set
	redisStatArticleThumb      = redisPrefix + ":stat:articles:thumb"  // sorted set

	redisArticleCachePrefix   = redisPrefix + ":article:cache:"   // string per article
	redisArticleViewPrefix    = redisPrefix + ":article:view:"    // set per article
	redisArticleThumbPrefix   = redisPrefix + ":article:thumb:"   // set per article
	redisArticleReviewPrefix  = redisPrefix + ":article:review:"  // set per article
	redisArticleRelatedPrefix = redisPrefix + ":article:related:" // sorted set per article
	//redisUserArticlePrefix    = redisPrefix + ":user:articles:" // sorted set per user

	redisDisLeaderboard    = redisPrefix + ":lb:distance:total" // sorted set
	redisMaxDisLeaderboard = redisPrefix + ":lb:distance:max"   // sorted set
	redisDurLeaderboard    = redisPrefix + ":lb:duration:total" // sorted set
	redisScorePhysicalLB   = redisPrefix + ":lb:score:physical" // sorted set
	redisScoreLiteralLB    = redisPrefix + ":lb:score:literal"  // sorted set
	redisScoreMentalLB     = redisPrefix + ":lb:score:mental"   // sorted set
	redisScoreWealthLB     = redisPrefix + ":lb:score:wealth"   // sorted set

	redisPubSubGroup = redisPrefix + ":pubsub:group:"
	redisPubSubUser  = redisPrefix + ":pubsub:user:"
)

const (
	onlineUserExpire = 60 * 60  // 15m online user timeout
	onlinesExpire    = 120 * 60 // 60m online set timeout
)

type RedisLogger struct {
	pool *redis.Pool
	conn redis.Conn
}

func NewRedisLogger(pool *redis.Pool, conn redis.Conn) *RedisLogger {
	return &RedisLogger{pool, conn}
}

func (logger *RedisLogger) Close() error {
	return logger.conn.Close()
}

func (logger *RedisLogger) PubSub(userid string, groups ...string) *redis.PubSubConn {
	channels := []interface{}{redisPubSubUser + userid}
	for _, group := range groups {
		channels = append(channels, redisPubSubGroup+group)
	}
	conn := redis.PubSubConn{logger.pool.Get()}
	conn.Subscribe(channels...)
	return &conn
}

func (logger *RedisLogger) Subscribe(psc *redis.PubSubConn, groups ...string) error {
	var channels []interface{}
	for _, group := range groups {
		channels = append(channels, redisPubSubGroup+group)
	}
	return psc.Subscribe(channels...)
}

func (logger *RedisLogger) Unsubscribe(psc *redis.PubSubConn, groups ...string) error {
	var channels []interface{}
	for _, group := range groups {
		channels = append(channels, redisPubSubGroup+group)
	}
	return psc.Unsubscribe(channels...)
}

func (logger *RedisLogger) PubMsg(typ string, to string, msg []byte) {
	conn := logger.pool.Get()
	defer conn.Close()

	switch typ {
	case "groupchat":
		conn.Do("PUBLISH", redisPubSubGroup+to, msg)
	default:
		conn.Do("PUBLISH", redisPubSubUser+to, msg)
	}
}

func (logger *RedisLogger) Users() int {
	count, _ := redis.Int(logger.conn.Do("HLEN", redisUserGuest))
	return count
}

// log register users per day
func (logger *RedisLogger) LogRegister(userid string) {
	logger.conn.Do("SADD", redisStatRegisterPrefix+DateString(time.Now()), userid)
}

func (logger *RedisLogger) RegisterCount(days int) []int64 {
	return logger.setsCount(redisStatRegisterPrefix, days)
}

func onlineTimeString() string {
	now := time.Now()
	min := now.Minute()
	if min < 30 {
		now = now.Add(time.Duration(0-min) * time.Minute)
	} else {
		now = now.Add(time.Duration(30-min) * time.Minute)
	}
	return now.Format("200601021504")
}

type redisUser struct {
	Userid    string  `redis:"userid"`
	Nickname  string  `redis:"nickname"`
	Profile   string  `redis:"profile"`
	RegTime   int64   `redis:"reg_time"`
	Role      string  `redis:"role"`
	Lng       float64 `redis:"lng"`
	Lat       float64 `redis:"lat"`
	SetInfo   bool    `redis:"setinfo"`
	WalletId  string  `redis:"wallet_id"`
	Sharedkey string  `redis:"shared_key"`
	Addr      string  `redis:"recv_addr"`
	Addrs     string  `redis:"addrs"`
	Score     int     `redis:"score"`
	Level     int     `redis:"level"`
}

func (logger *RedisLogger) OnlineUser(accessToken string) *Account {
	if len(accessToken) == 0 {
		return nil
	}
	user := &redisUser{}
	conn := logger.conn

	if strings.HasPrefix(accessToken, GuestUserPrefix) {
		user.Userid, _ = redis.String(conn.Do("HGET", redisUserGuest, accessToken))
	} else {
		v, err := redis.Values(conn.Do("HGETALL", redisUserOnlineUserPrefix+accessToken))
		if err != nil {
			log.Println(err)
			return nil
		}
		if err := redis.ScanStruct(v, user); err != nil {
			log.Println(err)
			return nil
		}
	}

	if len(user.Userid) == 0 {
		return nil
	}

	addrs := strings.Split(user.Addrs, ",")

	return &Account{
		Id:       user.Userid,
		Nickname: user.Nickname,
		Profile:  user.Profile,
		RegTime:  time.Unix(user.RegTime, 0),
		Role:     user.Role,
		Loc:      &Location{Lng: user.Lng, Lat: user.Lat},
		Setinfo:  user.SetInfo,
		Wallet:   DbWallet{Id: user.WalletId, Addrs: addrs, Addr: addrs[0], Key: user.Sharedkey},
		Score:    user.Score,
		Level:    user.Level,
	}
}

func (logger *RedisLogger) LogOnlineUser(accessToken string, user *Account) {
	if user == nil {
		return
	}

	conn := logger.conn

	u := &redisUser{
		Userid:    user.Id,
		Nickname:  user.Nickname,
		Profile:   user.Profile,
		RegTime:   user.RegTime.Unix(),
		Role:      user.Role,
		SetInfo:   user.Setinfo,
		WalletId:  user.Wallet.Id,
		Sharedkey: user.Wallet.Key,
		Addr:      user.Wallet.Addr,
		Addrs:     strings.Join(user.Wallet.Addrs, ","),
		Score:     user.Score,
		Level:     user.Level,
	}
	if user.Loc != nil {
		u.Lat = user.Loc.Lat
		u.Lng = user.Loc.Lng
	}

	conn.Send("MULTI")
	if !strings.HasPrefix(accessToken, GuestUserPrefix) {
		conn.Send("HMSET", redis.Args{}.Add(redisUserOnlineUserPrefix+accessToken).AddFlat(u)...)
		conn.Send("EXPIRE", redisUserOnlineUserPrefix+accessToken, onlineUserExpire)
	} else {
		conn.Send("HSETNX", redisUserGuest, accessToken, user.Id)
	}

	timeStr := onlineTimeString()
	conn.Send("SADD", redisUserOnlinesPrefix+timeStr, user.Id)
	conn.Send("EXPIRE", redisUserOnlinesPrefix+timeStr, onlinesExpire)
	conn.Do("EXEC")
}

func (logger *RedisLogger) Relationship(userid, peer string) string {
	if userid == peer || len(userid) == 0 || len(peer) == 0 {
		return ""
	}

	following, _ := redis.Bool(logger.conn.Do("SISMEMBER", redisUserFollowPrefix+userid, peer))
	follower, _ := redis.Bool(logger.conn.Do("SISMEMBER", redisUserFollowerPrefix+userid, peer))
	if following && follower {
		return RelFriend
	} else if following {
		return RelFollowing
	} else if follower {
		return RelFollower
	}
	if black, _ := redis.Bool(logger.conn.Do("SISMEMBER", redisUserBlacklistPrefix+userid, peer)); black {
		return RelBlacklist
	}

	return ""
}

func (logger *RedisLogger) SetRelationship(userid, peer string, relation string, enable bool) {
	if userid == peer || len(userid) == 0 || len(peer) == 0 {
		return
	}
	conn := logger.conn
	conn.Send("MULTI")

	switch relation {
	case RelFollowing:
		if enable {
			conn.Send("SADD", redisUserFollowPrefix+userid, peer)
			conn.Send("SADD", redisUserFollowerPrefix+peer, userid)
		} else {
			conn.Send("SREM", redisUserFollowPrefix+userid, peer)
			conn.Send("SREM", redisUserFollowerPrefix+peer, userid)
		}
	case RelBlacklist:
		if enable {
			conn.Send("SREM", redisUserFollowPrefix+userid, peer)
			conn.Send("SREM", redisUserFollowerPrefix+peer, userid)
			conn.Send("SADD", redisUserBlacklistPrefix+userid, peer)
		} else {
			conn.Send("SREM", redisUserBlacklistPrefix+userid, peer)
		}
	default:
	}

	conn.Do("EXEC")
}

/*
func (logger *RedisLogger) Followed(userid, following string) bool {
	if len(userid) == 0 || len(following) == 0 {
		return false
	}
	followed, _ := redis.Bool(logger.conn.Do("SISMEMBER", redisUserFollowPrefix+userid, following))
	return followed
}


func (logger *RedisLogger) SetFollow(userid, following string, follow bool) {
	if len(userid) == 0 || len(following) == 0 {
		return
	}
	conn := logger.conn
	conn.Send("MULTI")
	if follow {
		conn.Send("SADD", redisUserFollowPrefix+userid, following)
		conn.Send("SADD", redisUserFollowerPrefix+following, userid)
	} else {
		conn.Send("SREM", redisUserFollowPrefix+userid, following)
		conn.Send("SREM", redisUserFollowerPrefix+following, userid)
	}
	conn.Do("EXEC")
}
*/
func (logger *RedisLogger) SetWBImport(userid, wb string) {
	logger.conn.Do("SADD", redisUserWBImportPrefix+userid, wb)
}

func (logger *RedisLogger) ImportFriend(userid, friend string) {
	conn := logger.conn
	logger.SetRelationship(userid, friend, RelFollowing, true)
	conn.Do("SREM", redisUserWBImportPrefix+friend, userid)
}

func (logger *RedisLogger) Friends(typ string, userid string) (users []string) {
	var key string
	switch typ {
	case RelFollowing:
		key = redisUserFollowPrefix + userid
	case RelFollower:
		key = redisUserFollowerPrefix + userid
	case RelFriend:
		users, _ = redis.Strings(logger.conn.Do("SINTER",
			redisUserFollowerPrefix+userid, redisUserFollowPrefix+userid))
		return
	case RelBlacklist:
		key = redisUserBlacklistPrefix + userid
	case "weibo":
		key = redisUserWBImportPrefix + userid
	default:
		return
	}

	users, _ = redis.Strings(logger.conn.Do("SMEMBERS", key))
	return
}

func (logger *RedisLogger) FriendCount(userid string) (follows int, followers int, friends int) {
	conn := logger.conn

	conn.Send("MULTI")
	conn.Send("SCARD", redisUserFollowPrefix+userid)
	conn.Send("SCARD", redisUserFollowerPrefix+userid)
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
		return
	}
	counts := make([]int, 2)

	if err := redis.ScanSlice(values, &counts); err != nil {
		log.Println(err)
		return
	}

	follows = counts[0]
	followers = counts[1]
	friends = len(logger.Friends("friend", userid))
	return
}

func (logger *RedisLogger) JoinGroup(userid, gid string, join bool) {
	conn := logger.conn
	conn.Send("MULTI")
	if join {
		conn.Send("HSET", redisUserGroupPrefix+userid, gid, time.Now().Unix())
		conn.Send("SADD", redisGroupPrefix+gid, userid)
	} else {
		conn.Send("HDEL", redisUserGroupPrefix+userid, gid)
		conn.Send("SREM", redisGroupPrefix+gid, userid)
	}
	conn.Do("EXEC")
}

func (logger *RedisLogger) Groups(userid string) []string {
	v, _ := redis.Strings(logger.conn.Do("HKEYS", redisUserGroupPrefix+userid))
	return v
}

func (logger *RedisLogger) DelOnlineUser(accessToken string) *Account {
	conn := logger.conn

	user := &Account{}
	v, err := redis.Values(conn.Do("HGETALL", redisUserOnlineUserPrefix+accessToken))
	if err != nil {
		log.Println(err)
		return nil
	}
	if err := redis.ScanStruct(v, user); err != nil {
		log.Println(err)
		return nil
	}
	conn.Send("MULTI")
	conn.Send("DEL", redisUserOnlineUserPrefix+accessToken)
	conn.Send("SREM", redisUserOnlinesPrefix+onlineTimeString(), user.Id)
	conn.Do("EXEC")

	return user
}

func (logger *RedisLogger) IsOnline(userid string) bool {
	conn := logger.conn
	online, _ := redis.Bool(conn.Do("SISMEMBER", redisUserOnlinesPrefix+onlineTimeString(), userid))
	return online
}

func (logger *RedisLogger) Onlines() int {
	count, _ := redis.Int(logger.conn.Do("SCARD", redisUserOnlinesPrefix+onlineTimeString()))
	return count
}

func (logger *RedisLogger) setsCount(key string, days int) []int64 {
	if days <= 0 {
		days = 1
	}

	t := time.Now()
	d, _ := time.ParseDuration("-24h")

	conn := logger.conn

	conn.Send("MULTI")
	conn.Send("SCARD", key+DateString(t))
	for i := 1; i < days; i++ {
		t = t.Add(d)
		conn.Send("SCARD", key+DateString(t))
	}
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
		return nil
	}

	counts := make([]int64, len(values))
	for i, v := range values {
		counts[i], _ = v.(int64)
	}

	return counts
}

func (logger *RedisLogger) LogUserMessages(userid string, msgs ...string) {
	args := redis.Args{}.Add(redisUserMessagePrefix + userid).AddFlat(msgs)
	conn := logger.conn
	conn.Do("LPUSH", args...)
}

func (logger *RedisLogger) MessageCount(userid string) int {
	conn := logger.conn

	count, err := redis.Int(conn.Do("LLEN", redisUserMessagePrefix+userid))
	if err != nil {
		log.Println(err)
	}
	return count
}

func (logger *RedisLogger) ClearMessages(userid string) {
	conn := logger.conn
	conn.Do("DEL", redisUserMessagePrefix+userid)
}

// log unique visitors per day
func (logger *RedisLogger) LogVisitor(user string) {
	conn := logger.conn
	conn.Do("SADD", redisStatVisitorPrefix+DateString(time.Now()), user)
}

func (logger *RedisLogger) VisitorsCount(days int) []int64 {
	return logger.setsCount(redisStatVisitorPrefix, days)
}

// log pv per day
func (logger *RedisLogger) LogPV(path string) {
	conn := logger.conn
	conn.Do("ZINCRBY", redisStatPvPrefix+DateString(time.Now()), 1, path)
}

type KV struct {
	K string `json:"path"`
	V int64  `json:"count"`
}

func (logger *RedisLogger) PVs(dates ...string) map[string][]KV {
	if len(dates) == 0 {
		dates = []string{DateString(time.Now())}
	}

	pvs := make(map[string][]KV, len(dates))

	for _, date := range dates {
		pvs[date] = logger.PV(date)
	}

	return pvs
}

func (logger *RedisLogger) PV(date string) []KV {
	if len(date) == 0 {
		return nil
	}

	conn := logger.conn
	count, _ := redis.Int(conn.Do("ZCARD", redisStatPvPrefix+date))
	values, err := redis.Values(conn.Do("ZREVRANGE", redisStatPvPrefix+date, 0, count, "WITHSCORES"))

	if err != nil {
		log.Println(err)
		return nil
	}

	var pvs []KV

	if err := redis.ScanSlice(values, &pvs); err != nil {
		log.Println(err)
		return nil
	}
	return pvs
}

func (logger *RedisLogger) ArticleCount(articleId string) (view, thumb, review int64) {
	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZSCORE", redisStatArticleView, articleId)
	//conn.Send(conn.Do("SCARD", redisArticleViewPrefix+articleId))
	conn.Send("ZSCORE", redisStatArticleThumb, articleId)
	conn.Send("ZSCORE", redisStatArticleReview, articleId)
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
		return
	}

	var counts []struct {
		Count int64
	}

	if err := redis.ScanSlice(values, &counts); err != nil {
		log.Println(err)
		return
	}

	view = counts[0].Count
	thumb = counts[1].Count
	review = counts[2].Count

	//	log.Println(view, thumb, review)

	return
}

func (logger *RedisLogger) LogArticleView(articleId string, userid string) {
	conn := logger.conn
	//log.Println("log article view", articleId, userid)
	conn.Send("MULTI")
	conn.Send("ZINCRBY", redisStatArticleViewPrefix+DateString(time.Now()), 1, articleId)
	conn.Send("ZINCRBY", redisStatArticleView, 1, articleId)
	conn.Send("SADD", redisArticleViewPrefix+articleId, userid)
	//conn.Send("ZADD", redisUserArticlePrefix+userid, AccessRate, articleId)
	conn.Do("EXEC")
}

func (logger *RedisLogger) ArticleViewers(articleId string) []string {
	if len(articleId) == 0 {
		return nil
	}

	conn := logger.conn
	viewers, _ := redis.Strings(conn.Do("SMEMBERS", redisArticleViewPrefix+articleId))

	return viewers
}

func (logger *RedisLogger) ArticleView(userid string, articles ...string) []bool {
	if len(userid) == 0 {
		return nil
	}

	conn := logger.conn
	conn.Send("MULTI")
	for _, article := range articles {
		conn.Send("SISMEMBER", redisArticleViewPrefix+article, userid)
	}
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil || len(values) != len(articles) {
		log.Println(err)
		return nil
	}

	views := make([]bool, len(articles))
	for i, v := range values {
		if b, ok := v.(int64); ok && b != 0 {
			views[i] = true
		}
	}
	return views
}

func (logger *RedisLogger) ArticleTopView(days, max int) []string {
	if days <= 0 {
		days = 1
	}
	if max <= 0 {
		max = 3
	}

	t := time.Now()
	d, _ := time.ParseDuration("-24h")

	keys := make([]string, days)
	keys[0] = redisStatArticleViewPrefix + DateString(t)
	for i := 1; i < days; i++ {
		t = t.Add(d)
		keys[i] = redisStatArticleViewPrefix + DateString(t)
	}

	args := redis.Args{}.Add(redisStatArticleViewPrefix + "out").Add(days).AddFlat(keys)
	//log.Println(args)
	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZUNIONSTORE", args...)
	conn.Send("ZREVRANGE", redisStatArticleViewPrefix+"out", 0, max, "WITHSCORES")
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
		return nil
	}

	var tops []KV
	s, _ := values[1].([]interface{})

	if err := redis.ScanSlice(s, &tops); err != nil {
		log.Println(err)
		return nil
	}

	articles := make([]string, len(tops))
	for i, _ := range tops {
		articles[i] = tops[i].K
	}

	return articles
}

func (logger *RedisLogger) LogArticleReview(userid, articleId string) {
	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZINCRBY", redisStatArticleReview, 1, articleId)
	conn.Send("SADD", redisArticleReviewPrefix+articleId, userid)
	//conn.Send("ZADD", redisUserArticlePrefix+userid, ReviewRate|AccessRate, articleId)
	conn.Do("EXEC")
}

func (logger *RedisLogger) ArticleReviewCount(articleId string) (count int) {
	conn := logger.conn
	count, _ = redis.Int(conn.Do("ZSCORE", redisStatArticleReview, articleId))
	return
}

func (logger *RedisLogger) ArticleTopReview(max int) []string {
	if max <= 0 {
		max = 1
	}
	conn := logger.conn
	articles, err := redis.Strings(conn.Do("ZREVRANGE", redisStatArticleReview, 0, max))
	if err != nil {
		log.Println(err)
		return nil
	}

	return articles
}

func (logger *RedisLogger) LogArticleThumb(userid, articleId string, thumb bool) {
	inc := 1
	if !thumb {
		inc = -1
	}
	conn := logger.conn
	//log.Println("log article thumb", userid, articleId, thumb)
	conn.Send("MULTI")
	conn.Send("ZINCRBY", redisStatArticleThumb, inc, articleId)
	if thumb {
		conn.Send("SADD", redisArticleThumbPrefix+articleId, userid)
		//conn.Send("ZADD", redisUserArticlePrefix+userid, ThumbRate|AccessRate, articleId)
	} else {
		conn.Send("SREM", redisArticleThumbPrefix+articleId, userid)
	}
	conn.Do("EXEC")
}

func (logger *RedisLogger) ArticleThumbers(articleId string) []string {
	if len(articleId) == 0 {
		return nil
	}

	conn := logger.conn
	thumbers, _ := redis.Strings(conn.Do("SMEMBERS", redisArticleThumbPrefix+articleId))

	return thumbers
}

func (logger *RedisLogger) ArticleThumbed(userid, articleId string) (b bool) {
	conn := logger.conn
	b, _ = redis.Bool(conn.Do("SISMEMBER", redisArticleThumbPrefix+articleId, userid))
	return
}

func (logger *RedisLogger) ArticleThumbCount(articleId string) (count int) {
	conn := logger.conn
	count, _ = redis.Int(conn.Do("SCARD", redisArticleThumbPrefix+articleId))
	return
}

func (logger *RedisLogger) ArticleTopThumb(max int) []string {
	if max <= 0 {
		max = 1
	}
	conn := logger.conn
	articles, err := redis.Strings(conn.Do("ZREVRANGE", redisStatArticleThumb, 0, max))
	if err != nil {
		log.Println(err)
		return nil
	}

	return articles
}

func (logger *RedisLogger) UpdateRecLB(userid string, distance, duration int) {
	if len(userid) == 0 {
		return
	}
	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZINCRBY", redisDisLeaderboard, distance, userid)
	conn.Send("ZINCRBY", redisDurLeaderboard, duration, userid)
	conn.Do("EXEC")
	if rec := logger.MaxDisRecord(userid); rec < distance {
		conn.Send("MULTI")
		conn.Send("ZREM", redisMaxDisLeaderboard, userid)
		conn.Send("ZINCRBY", redisMaxDisLeaderboard, distance, userid)
		conn.Do("EXEC")
	}
}

func (logger *RedisLogger) MaxDisRecord(userid string) int {
	max, _ := redis.Int(logger.conn.Do("ZSCORE", redisMaxDisLeaderboard, userid))
	return max
}

func (logger *RedisLogger) RecStats(userid string) (int, int) {
	dis, _ := redis.Int(logger.conn.Do("ZSCORE", redisDisLeaderboard, userid))
	dur, _ := redis.Int(logger.conn.Do("ZSCORE", redisDurLeaderboard, userid))
	return dis, dur
}

func (logger *RedisLogger) LBDisRank(userid string) int {
	if len(userid) == 0 {
		return -1
	}
	rank, err := redis.Int(logger.conn.Do("ZREVRANK", redisDisLeaderboard, userid))
	if err != nil {
		return -1
	}
	return rank
}

func (logger *RedisLogger) LBDisCard() int {
	count, _ := redis.Int(logger.conn.Do("ZCARD", redisDisLeaderboard))
	return count
}

func (logger *RedisLogger) LBDurRank(userid string) int {
	if len(userid) == 0 {
		return -1
	}
	rank, err := redis.Int(logger.conn.Do("ZREVRANK", redisDurLeaderboard, userid))
	if err != nil {
		return -1
	}
	return rank
}

func (logger *RedisLogger) UserProps(userid string) *Props {
	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZSCORE", redisScorePhysicalLB, userid)
	conn.Send("ZSCORE", redisScoreLiteralLB, userid)
	conn.Send("ZSCORE", redisScoreMentalLB, userid)
	conn.Send("ZSCORE", redisScoreWealthLB, userid)
	values, _ := redis.Values(conn.Do("EXEC"))

	var scores []int64
	if err := redis.ScanSlice(values, &scores); err != nil {
		log.Println(err)
		return nil
	}

	props := &Props{
		Physical: scores[0],
		Literal:  scores[1],
		Mental:   scores[2],
		Wealth:   scores[3],
	}

	props.Score = int64(UserScore(props))
	props.Level = int64(UserLevel(int(props.Score)))

	return props
}

func (logger *RedisLogger) Props(typ string, ids ...string) []int {
	conn := logger.conn
	var key string

	switch typ {
	case ScorePhysical:
		key = redisScorePhysicalLB
	case ScoreLiteral:
		key = redisScoreLiteralLB
	case ScoreMental:
		key = redisScoreMentalLB
	case ScoreWealth:
		key = redisScoreWealthLB
	default:
		key = redisScorePhysicalLB
	}
	conn.Send("MULTI")
	for _, id := range ids {
		conn.Send("ZSCORE", key, id)
	}
	values, _ := redis.Values(conn.Do("EXEC"))
	var scores []int

	if err := redis.ScanSlice(values, &scores); err != nil {
		log.Println(err)
		return nil
	}
	return scores
}

func (logger *RedisLogger) AddProps(userid string, props *Props) (*Props, error) {
	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZINCRBY", redisScorePhysicalLB, props.Physical, userid)
	conn.Send("ZINCRBY", redisScoreLiteralLB, props.Literal, userid)
	conn.Send("ZINCRBY", redisScoreMentalLB, props.Mental, userid)
	conn.Send("ZINCRBY", redisScoreWealthLB, props.Wealth, userid)
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var newScores []int64
	if err := redis.ScanSlice(values, &newScores); err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println("new scores:", newScores)
	props.Physical = newScores[0]
	props.Literal = newScores[1]
	props.Mental = newScores[2]
	props.Wealth = newScores[3]

	return props, nil
}

func (logger *RedisLogger) GetDisLB(start, stop int) []KV {
	values, _ := redis.Values(logger.conn.Do("ZREVRANGE", redisDisLeaderboard, start, stop, "WITHSCORES"))
	var s []KV

	if err := redis.ScanSlice(values, &s); err != nil {
		log.Println(err)
		return nil
	}
	return s
}

func (logger *RedisLogger) Transaction(from, to string, amount int64) {
	if len(from) == 0 || len(to) == 0 || amount <= 0 {
		return
	}
	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZINCRBY", redisScoreWealthLB, -amount, from)
	conn.Send("ZINCRBY", redisScoreWealthLB, amount, to)
	conn.Do("EXEC")
}
