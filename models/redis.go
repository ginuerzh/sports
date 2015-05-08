// log
package models

import (
	//"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	//"strconv"
	//"encoding/json"
	//"labix.org/v2/mgo/bson"
	//"strings"
	"time"
)

const (
	redisPrefix                 = "sports"
	redisStatVisitorPrefix      = redisPrefix + ":stat:visitors:"       // set per day
	redisStatPvPrefix           = redisPrefix + ":stat:pv:"             // sorted set per day
	redisStatRegisterPrefix     = redisPrefix + ":stat:registers:"      // set per day, register users per day
	redisStatRegPhonePrefix     = redisPrefix + ":stat:register:phone:" // set per day, phone register users per day
	redisStatRegEmailPrefix     = redisPrefix + ":stat:register:email:" // set per day, email register users per day
	redisStatRegWeiboPrefix     = redisPrefix + ":stat:register:weibo:" // set per day, weibo register users per day
	redisStatLoginPrefix        = redisPrefix + ":stat:logins:"         // set per day, login users per day
	redisStatCoachLoginPrefix   = redisPrefix + ":stat:coach:logins:"   // set per day, login coaches per day
	redisStatOnlines            = redisPrefix + ":stat:onlines"         // set, current online users
	redisStatCoachOnlines       = redisPrefix + ":stat:coach:onlines"   // set, current online coaches
	redisStatOnlineTime         = redisPrefix + ":stat:onlinetime"      // sorted set, users total online time
	redisStatUserArticlesPrefix = redisPrefix + ":stat:articles:"       // sorted set per day, user articles per day, user:articless
	redisStatUserCommentsPrefix = redisPrefix + ":stat:comments:"       // sorted set per day, user comments per day, user:comments
	redisStatUserPostsPrefix    = redisPrefix + ":stat:posts:"          // sorted set per day, user posts per day, user:posts
	redisStatPosts              = redisPrefix + ":stat:posts"           // sorted set, total posts per day
	redisStatUserTotalArticles  = redisPrefix + ":stat:articles"        // sorted set, total user articles
	redisStatUserTotalComments  = redisPrefix + ":stat:comments"        // sorted set, total user comments
	redisStatUserTotalPosts     = redisPrefix + ":stat:posts"           // sorted set, total user posts (articles + comments)
	redisStatGamersPrefix       = redisPrefix + ":stat:gamers:"         // sorted set per day, game time per day (gamer:time)
	redisStatTotalGamers        = redisPrefix + ":stat:gamers"          // sorted set, total game time
	redisStatGameTime           = redisPrefix + ":stat:gametime"        // sorted set, total game time per day
	redisStatUserRecordsPrefix  = redisPrefix + ":stat:records:"        // sorted set per day, user records per day
	redisStatUserTotalRecords   = redisPrefix + ":stat:records"         // sorted set total user records
	redisStatAuthCoachesPrefix  = redisPrefix + ":stat:authcoaches:"    // set per day, auth coaches per day
	redisStatUserCoinsPrefix    = redisPrefix + ":stat:coins:"          // sorted set per day, coins send by system per day
	redisStatCoins              = redisPrefix + ":stat:coins"           // sorted set, total coins sended per day

	//redisUserOnlinesPrefix = redisPrefix + ":user:onlines:" // set per half an hour, current online users
	redisUserTokens     = redisPrefix + ":user:tokens" // hash, online user token <->userid
	RedisUserInfoPrefix = redisPrefix + ":user:info:"  // hashs per user, user's event box at now
	RedisUserCoins      = redisPrefix + ":user:coins"  // sorted set
	//redisUserGuest            = redisPrefix + ":user:guest"    // hashes for all guests
	//redisUserMessagePrefix    = redisPrefix + ":user:msgs:"         // list per user
	redisUserFollowPrefix    = redisPrefix + ":user:follow:"       // set per user
	redisUserFollowerPrefix  = redisPrefix + ":user:follower:"     // set per user
	redisUserBlacklistPrefix = redisPrefix + ":user:blacklist:"    // set per user
	redisUserWBImportPrefix  = redisPrefix + ":user:import:weibo:" // set per user
	redisUserGroupPrefix     = redisPrefix + ":user:group:"        // hash per user
	redisGroupPrefix         = redisPrefix + ":group:"             // set per group

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

	redisGameLB01 = redisPrefix + ":lb:game:01" // 77 jump
	redisGameLB02 = redisPrefix + ":lb:game:02" // escape
	redisGameLB03 = redisPrefix + ":lb:game:03" // bear
	redisGameLB04 = redisPrefix + ":lb:game:04" // line
	redisGameLB05 = redisPrefix + ":lb:game:05" // turn

	redisPubSubGroup = redisPrefix + ":pubsub:group:"
	redisPubSubUser  = redisPrefix + ":pubsub:user:"

	redisNoticeChannel = redisPrefix + ":pubsub:notice"
)

const (
	onlineUserExpire = 30 * 24 * 60 * 60 // 1mon online user timeout
	onlinesExpire    = 120 * 60          // 60m online set timeout
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

func (logger *RedisLogger) Notice(msg []byte) {
	conn := logger.pool.Get()
	defer conn.Close()
	conn.Do("PUBLISH", redisNoticeChannel, msg)
}

func (logger *RedisLogger) PubMsg(typ string, to string, msg []byte) {
	conn := logger.pool.Get()
	defer conn.Close()

	switch typ {
	case "groupchat":
		conn.Do("PUBLISH", redisPubSubGroup+to, msg)
	default:
		conn.Do("PUBLISH", redisPubSubUser+to, msg)
		//log.Println("publish to", to, string(msg), reply, err)
	}
}

/*
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
*/
func (logger *RedisLogger) OnlineUser(token string) (id string) {
	conn := logger.conn

	id, _ = redis.String(conn.Do("HGET", redisUserTokens, token))
	return
}

func (logger *RedisLogger) SetOnlineUser(token string, userid string) {
	if len(token) == 0 || len(userid) == 0 {
		return
	}

	logger.conn.Do("HSET", redisUserTokens, token, userid)
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

func (logger *RedisLogger) SetRelationship(userid string, peers []string, relation string, enable bool) {
	if len(userid) == 0 || len(peers) == 0 {
		return
	}
	//log.Println("set relationship", userid, peers, relation, enable)
	conn := logger.conn
	conn.Send("MULTI")

	for _, peer := range peers {
		if peer == userid || len(peer) == 0 {
			continue
		}
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
				conn.Send("SREM", redisUserFollowPrefix+peer, userid)
				conn.Send("SREM", redisUserFollowerPrefix+peer, userid)
				conn.Send("SREM", redisUserFollowerPrefix+userid, peer)
				conn.Send("SADD", redisUserBlacklistPrefix+userid, peer)
			} else {
				conn.Send("SREM", redisUserBlacklistPrefix+userid, peer)
			}
		default:
		}
	}

	conn.Do("EXEC")
}

/*
func (logger *RedisLogger) SetWBImport(userid, wb string) {
	logger.conn.Do("SADD", redisUserWBImportPrefix+userid, wb)
}

func (logger *RedisLogger) ImportFriend(userid, friend string) {
	conn := logger.conn
	logger.SetRelationship(userid, []string{friend}, RelFollowing, true)
	conn.Do("SREM", redisUserWBImportPrefix+friend, userid)
}
*/
func (logger *RedisLogger) Friends(types string, userid string) (users []string) {
	var key string
	switch types {
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

func (logger *RedisLogger) FriendCount(userid string) (follows, followers, friends, blacklist int) {
	conn := logger.conn

	conn.Send("MULTI")
	conn.Send("SCARD", redisUserFollowPrefix+userid)
	conn.Send("SCARD", redisUserFollowerPrefix+userid)
	conn.Send("SCARD", redisUserBlacklistPrefix+userid)
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
		return
	}
	counts := make([]int, 3)

	if err := redis.ScanSlice(values, &counts); err != nil {
		log.Println(err)
		return
	}

	follows = counts[0]
	followers = counts[1]
	blacklist = counts[2]
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

func (logger *RedisLogger) DelOnlineUser(token string) {
	conn := logger.conn

	userid, _ := redis.String(conn.Do("HGET", redisUserTokens, token))
	conn.Send("MULTI")
	conn.Send("HDEL", redisUserTokens, token)
	//conn.Send("SREM", redisUserOnlinesPrefix+onlineTimeString(), userid)
	conn.Send("SREM", redisStatOnlines, userid)
	conn.Do("EXEC")
}

func (logger *RedisLogger) IsOnline(userid string) bool {
	conn := logger.conn
	online, _ := redis.Bool(conn.Do("SISMEMBER", redisStatOnlines, userid))
	return online
}

func (logger *RedisLogger) Onlines() int {
	count, _ := redis.Int(logger.conn.Do("SCARD", redisStatOnlines))
	return count
}

func (logger *RedisLogger) scards(key string, days int) []int64 {
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

func (logger *RedisLogger) zcards(key string, days int) []int64 {
	if days <= 0 {
		days = 1
	}

	t := time.Now()
	d, _ := time.ParseDuration("-24h")

	conn := logger.conn

	conn.Send("MULTI")
	conn.Send("ZCARD", key+DateString(t))
	for i := 1; i < days; i++ {
		t = t.Add(d)
		conn.Send("ZCARD", key+DateString(t))
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

/*
func (logger *RedisLogger) EventCount(userid string) (counts []int) {
	counts = make([]int, 6)
	conn := logger.conn
	values, err := redis.Values(conn.Do("HMGET", RedisUserInfoPrefix+userid,
		"event_chat", "event_comment", "event_thumb", "event_reward", "event_subscribe", "event_tx"))
	if err != nil {
		log.Println(err)
		return
	}
	if err = redis.ScanSlice(values, &counts); err != nil {
		log.Println(err)
	}
	return
}


func (logger *RedisLogger) IncrEventCount(userid string, eventType string, count int) {
	if len(userid) == 0 || count == 0 {
		return
	}
	logger.conn.Do("HINCRBY", RedisUserInfoPrefix+userid, "event_"+eventType, count)
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
*/
// log unique visitors per day
func (logger *RedisLogger) LogVisitor(user string) {
	conn := logger.conn
	conn.Do("SADD", redisStatVisitorPrefix+DateString(time.Now()), user)
}

func (logger *RedisLogger) VisitorsCount(days int) []int64 {
	return logger.scards(redisStatVisitorPrefix, days)
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

// for sorting
type KVSlice []KV

func (p KVSlice) Len() int {
	return len(p)
}

func (p KVSlice) Less(i, j int) bool {
	return p[i].V < p[j].V
}

func (p KVSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
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

/*
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
*/
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

func (logger *RedisLogger) GetDisLB(start, stop int) []KV {
	values, _ := redis.Values(logger.conn.Do("ZREVRANGE", redisDisLeaderboard, start, stop, "WITHSCORES"))
	var s []KV

	if err := redis.ScanSlice(values, &s); err != nil {
		log.Println(err)
		return nil
	}
	return s
}

func (logger *RedisLogger) GetCoins(userid string) int64 {
	coins, _ := redis.Int64(logger.conn.Do("ZSCORE", RedisUserCoins, userid))
	return coins
}

func (logger *RedisLogger) SendCoins(userid string, coins int64) {
	if coins <= 0 {
		return
	}
	sdate := DateString(time.Now())

	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZINCRBY", RedisUserCoins, coins, userid)
	conn.Send("ZINCRBY", redisStatUserCoinsPrefix+sdate, coins, userid)
	conn.Send("ZINCRBY", redisStatCoins, coins, sdate)
	conn.Send("EXEC")
}

func (logger *RedisLogger) Transaction(from, to string, amount int64) {
	if len(from) == 0 || len(to) == 0 || amount <= 0 {
		return
	}
	conn := logger.conn
	conn.Send("MULTI")
	conn.Send("ZINCRBY", RedisUserCoins, -amount, from)
	conn.Send("ZINCRBY", RedisUserCoins, amount, to)
	conn.Do("EXEC")
}

func lbGameKey(typ int) string {
	var key string

	switch typ {
	case 0x01:
		key = redisGameLB01
	case 0x02:
		key = redisGameLB02
	case 0x03:
		key = redisGameLB03
	case 0x04:
		key = redisGameLB04
	case 0x05:
		key = redisGameLB05
	default:
		key = redisGameLB01
	}

	return key
}

func (logger *RedisLogger) SetGameScore(typ int, userid string, score int) {
	logger.conn.Do("ZADD", lbGameKey(typ), score, userid)
}

func (logger *RedisLogger) SetGameMaxScore(typ int, userid string, score int) {
	if len(userid) == 0 || score == 0 {
		return
	}

	conn := logger.conn

	key := lbGameKey(typ)
	max, _ := redis.Int(conn.Do("ZSCORE", key, userid))

	if score > max {
		conn.Do("ZINCRBY", key, score-max, userid)
	}
}

func (logger *RedisLogger) GameScores(typ int, skip, limit int) []KV {
	return logger.zrange(lbGameKey(typ), skip, skip+limit-1, true)
}

func (logger *RedisLogger) UserGameScores(typ int, userids ...string) []int64 {
	return logger.zscores(lbGameKey(typ), userids...)
}

func (logger *RedisLogger) UserGameRanks(typ int, userids ...string) []int {
	return logger.zrevrank(lbGameKey(typ), userids...)
}

func (logger *RedisLogger) GameUserCount(typ int) int {
	return logger.zcard(lbGameKey(typ))
}

func (logger *RedisLogger) zscores(key string, members ...string) (scores []int64) {
	conn := logger.conn

	if len(members) == 0 {
		return nil
	}

	conn.Send("MULTI")
	for _, member := range members {
		conn.Send("ZSCORE", key, member)
	}
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
		return nil
	}

	if err := redis.ScanSlice(values, &scores); err != nil {
		log.Println(err)
		return nil
	}

	return
}

func (logger *RedisLogger) zrange(key string, start, stop int, reverse bool) (kv []KV) {
	cmd := "ZRANGE"
	if reverse {
		cmd = "ZREVRANGE"
	}

	values, _ := redis.Values(logger.conn.Do(cmd, key, start, stop, "WITHSCORES"))

	if err := redis.ScanSlice(values, &kv); err != nil {
		log.Println(err)
		return nil
	}

	return
}

func (logger *RedisLogger) zrevrank(key string, members ...string) (ranks []int) {
	conn := logger.conn

	if len(members) == 0 {
		return nil
	}

	conn.Send("MULTI")
	for _, member := range members {
		conn.Send("ZREVRANK", key, member)
	}
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
		return nil
	}

	if err := redis.ScanSlice(values, &ranks); err != nil {
		log.Println(err)
		return nil
	}
	return
}

func (logger *RedisLogger) zcard(key string) int {
	n, _ := redis.Int(logger.conn.Do("ZCARD", key))
	return n
}

func (logger *RedisLogger) SetOnline(userid string, actor string, add bool, duration int64) {
	conn := logger.conn
	sdate := DateString(time.Now())
	//t := onlineTimeString()
	conn.Send("MULTI")
	//conn.Send("SADD", redisUserOnlinesPrefix+t, userid)
	//conn.Send("EXPIRE", redisUserOnlinesPrefix+t, onlinesExpire)
	if add {
		//conn.Send("SADD", redisStatOnlinesPrefix+DateString(time.Now()), userid)
		conn.Send("SADD", redisStatOnlines, userid)
		conn.Send("SADD", redisStatLoginPrefix+sdate, userid)
		if actor == ActorCoach {
			conn.Send("SADD", redisStatCoachOnlines, userid)
			conn.Send("SADD", redisStatCoachLoginPrefix+sdate, userid)
		}
	} else {
		conn.Send("SREM", redisStatOnlines, userid)
		conn.Send("ZINCRBY", redisStatOnlineTime, duration, userid)
		if actor == ActorCoach {
			conn.Send("SREM", redisStatCoachOnlines, userid)
		}
	}
	conn.Do("EXEC")
}

/*
func (logger *RedisLogger) LogLogin(userid string) {
	logger.conn.Do("SADD", redisStatLoginPrefix+DateString(time.Now()), userid)
}
*/

func (logger *RedisLogger) LoginCount(days int) []int64 {
	return logger.scards(redisStatLoginPrefix, days)
}
func (logger *RedisLogger) CoachLoginCount(days int) []int64 {
	return logger.scards(redisStatCoachLoginPrefix, days)
}

func (logger *RedisLogger) UserTotalOnlineTime(skip, limit int, desc bool) (kv []KV) {
	return logger.zrange(redisStatOnlineTime, skip, skip+limit-1, desc)
}

/*
func (logger *RedisLogger) OnlineUsersCount(days int) []int64 {
	t := time.Now()
	d, _ := time.ParseDuration("-24h")

	members := []string{DateString(t)}
	for i := 1; i < days; i++ {
		t = t.Add(d)
		members = append(members, DateString(t))
	}
	return logger.zscores(redisStatOnlinesPrefix, members...)
}
*/
// stats
// log register users per day
func (logger *RedisLogger) LogRegister(userid, types string) {
	sdate := DateString(time.Now())

	conn := logger.conn
	conn.Send("MULTI")
	switch types {
	case AccountEmail:
		conn.Send("SADD", redisStatRegEmailPrefix+sdate, userid)
	case AccountPhone:
		conn.Send("SADD", redisStatRegPhonePrefix+sdate, userid)
	case AccountWeibo:
		conn.Send("SADD", redisStatRegWeiboPrefix+sdate, userid)
	}
	conn.Send("SADD", redisStatRegisterPrefix+sdate, userid)
	conn.Do("EXEC")
}

func (logger *RedisLogger) RegisterCount(days int, types string) []int64 {
	key := redisStatRegisterPrefix

	switch types {
	case AccountEmail:
		key = redisStatRegEmailPrefix
	case AccountPhone:
		key = redisStatRegPhonePrefix
	case AccountWeibo:
		key = redisStatRegWeiboPrefix
	}
	return logger.scards(key, days)
}

func (logger *RedisLogger) AddPost(userid, types string, count int) {
	sdate := DateString(time.Now())

	conn := logger.conn
	conn.Send("MULTI")

	switch types {
	case "comment":
		conn.Send("ZINCRBY", redisStatUserCommentsPrefix+sdate, count, userid)
		conn.Send("ZINCRBY", redisStatUserTotalComments, count, userid)
	case "article":
		fallthrough
	default:
		conn.Send("ZINCRBY", redisStatUserArticlesPrefix+sdate, count, userid)
		conn.Send("ZINCRBY", redisStatUserTotalArticles, count, userid)
	}
	conn.Send("ZINCRBY", redisStatUserPostsPrefix+sdate, count, userid)
	conn.Send("ZINCRBY", redisStatUserTotalPosts, count, userid)
	conn.Send("ZINCRBY", redisStatPosts, count, sdate)

	conn.Send("EXEC")
}
func (logger *RedisLogger) PostUserCount(days int) []int64 {
	return logger.zcards(redisStatUserPostsPrefix, days)
}
func (logger *RedisLogger) PostsCount(days int) []int64 {
	t := time.Now()
	d, _ := time.ParseDuration("-24h")

	members := []string{DateString(t)}
	for i := 1; i < days; i++ {
		t = t.Add(d)
		members = append(members, DateString(t))
	}
	return logger.zscores(redisStatPosts, members...)
}
func (logger *RedisLogger) UserTotalPosts(skip, limit int, desc bool) (kv []KV) {
	return logger.zrange(redisStatUserTotalPosts, skip, skip+limit-1, desc)
}

func (logger *RedisLogger) AddGameTime(userid string, seconds int) {
	sdate := DateString(time.Now())
	conn := logger.conn
	conn.Send("MULTI")

	conn.Send("ZINCRBY", redisStatGamersPrefix+sdate, seconds, userid)
	conn.Send("ZINCRBY", redisStatGameTime, seconds, sdate)
	conn.Send("ZINCRBY", redisStatTotalGamers, seconds, userid)
	conn.Send("EXEC")
}
func (logger *RedisLogger) GamersCount(days int) []int64 {
	return logger.zcards(redisStatGamersPrefix, days)
}
func (logger *RedisLogger) GameTime(days int) []int64 {
	t := time.Now()
	d, _ := time.ParseDuration("-24h")

	members := []string{DateString(t)}
	for i := 1; i < days; i++ {
		t = t.Add(d)
		members = append(members, DateString(t))
	}
	return logger.zscores(redisStatGameTime, members...)
}
func (logger *RedisLogger) UserTotalGameTime(skip, limit int, desc bool) (kv []KV) {
	return logger.zrange(redisStatTotalGamers, skip, skip+limit-1, desc)
}

func (logger *RedisLogger) AddRecord(userid string, count int) {
	sdate := DateString(time.Now())

	conn := logger.conn
	conn.Send("MULTI")

	conn.Send("ZINCRBY", redisStatUserRecordsPrefix+sdate, count, userid)
	conn.Send("ZINCRBY", redisStatUserTotalRecords, count, userid)
	conn.Send("EXEC")
}
func (logger *RedisLogger) RecordUsersCount(days int) []int64 {
	return logger.zcards(redisStatUserRecordsPrefix, days)
}
func (logger *RedisLogger) UserTotalRecords(skip, limit int, desc bool) (kv []KV) {
	return logger.zrange(redisStatUserTotalRecords, skip, skip+limit-1, desc)
}

func (logger *RedisLogger) AddAuthCoach(userid string) {
	logger.conn.Do("SADD", redisStatAuthCoachesPrefix+DateString(time.Now()), userid)
}
func (logger *RedisLogger) AuthCoachesCount(days int) []int64 {
	return logger.scards(redisStatAuthCoachesPrefix, days)
}

func (logger *RedisLogger) CoinsCount(days int) []int64 {
	t := time.Now()
	d, _ := time.ParseDuration("-24h")

	members := []string{DateString(t)}
	for i := 1; i < days; i++ {
		t = t.Add(d)
		members = append(members, DateString(t))
	}
	return logger.zscores(redisStatCoins, members...)
}

func (logger *RedisLogger) Retention(date time.Time) []int {
	conn := logger.conn

	counts := make([]int, 8)

	counts[0], _ = redis.Int(conn.Do("SCARD", redisStatRegisterPrefix+DateString(date)))
	d := date.AddDate(0, 0, 1)
	for i := 0; i < 7; i++ {
		s, _ := redis.Strings(conn.Do("SINTER", redisStatRegisterPrefix+DateString(date),
			redisStatLoginPrefix+DateString(d)))
		counts[i+1] = len(s)
		d = d.AddDate(0, 0, 1)
	}

	return counts
}
