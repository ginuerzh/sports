// wallet
package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	errs "errors"
	btcaddr "github.com/ginuerzh/gimme-bitcoin-address"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"labix.org/v2/mgo/bson"
	//"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	CoinAddr = "http://localhost:8087"
)

func BindWalletApi(m *martini.ClassicMartini) {
	m.Get("/1/wallet/get",
		binding.Form(walletForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		getWalletHandler)
	m.Get("/1/wallet/balance",
		binding.Form(walletForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		balanceHandler)
	m.Get("/1/wallet/newaddr",
		binding.Form(walletForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		newAddrHandler)
	m.Post("/1/wallet/send",
		binding.Json(txForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		txHandler)
	m.Get("/1/wallet/txs",
		binding.Form(addrTxsForm{}),
		addrTxsHandler)
}

type walletForm struct {
	parameter
}

func getWalletHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	wal, err := getWallet(user.Wallet.Id, user.Wallet.Key)
	if err != nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, err.Error()))
		return
	}

	writeResponse(r.RequestURI, w, wal, nil)
}

func newAddrHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	wal, err := getWallet(user.Wallet.Id, user.Wallet.Key)
	if err != nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, err.Error()))
		return
	}
	k, err := wal.GenKey("")
	if err != nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, err.Error()))
		return
	}
	wal.AddKey(k)
	if _, err = saveWallet(user.Wallet.Id, wal); err != nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, err.Error()))
		return
	}

	user.Wallet.Addrs = append(user.Wallet.Addrs, k.PubKey)
	if err = user.AddWalletAddr(k.PubKey); err != nil {
		writeResponse(r.RequestURI, w, nil, err)
		return
	}

	writeResponse(r.RequestURI, w, k, nil)
}

func balanceHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	/*
		wal, err := getWallet(user.Wallet.Id, user.Wallet.Key)
		if err != nil {
			writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, err.Error()))
			return
		}

		var addrs []string
		for _, key := range wal.Keys {
			addrs = append(addrs, key.PubKey)
		}
	*/
	addrs := user.Wallet.Addrs
	balance, _ := getBalance(addrs)
	writeResponse(r.RequestURI, w, balance, nil)
}

type addrTxsForm struct {
	Addr string `form:"addr"`
}

func addrTxsHandler(r *http.Request, w http.ResponseWriter, redis *models.RedisLogger, form addrTxsForm) {
	txs, err := getAddrTxs(form.Addr)
	if err != nil {
		writeResponse(r.RequestURI, w, txs, err)
		return
	}

	//user := &models.Account{}
	//user.FindByWalletAddr(form.Addr)

	for i, tx := range txs {
		send := false
		if tx.Vin[0].PrevOut.Address == form.Addr {
			send = true
		}

		for _, out := range tx.Vout {
			if send && out.Address != tx.Vin[0].PrevOut.Address {
				txs[i].Amount -= out.Value
			}
			if !send && out.Address == form.Addr {
				txs[i].Amount += out.Value
			}
		}
	}

	writeResponse(r.RequestURI, w, map[string]interface{}{"txs": txs}, nil)
}

type txForm struct {
	Type     string `json:"trade_type"`
	Id       string `json:"article_id"`
	FromAddr string `json:"from"`
	ToAddr   string `json:"to" binding:"required"`
	Value    int64  `json:"value" binding:"required"`
	parameter
}

func txHandler(r *http.Request, w http.ResponseWriter,
	client *ApnClient, redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(txForm)

	if form.FromAddr == form.ToAddr {
		writeResponse(r.RequestURI, w, map[string]string{"txid": ""}, nil)
		return
	}

	receiver := &models.Account{}
	if find, err := receiver.FindByWalletAddr(form.ToAddr); !find {
		e := errors.NewError(errors.NotFoundError, "无效的收款地址")
		if err != nil {
			e = errors.NewError(errors.DbError, "无效的钱包地址")
		}
		writeResponse(r.RequestURI, w, nil, e)
		return
	}

	wal, err := getWallet(user.Wallet.Id, user.Wallet.Key)
	if err != nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, "获取钱包失败"))
		return
	}

	outputs, amount, err := getUnspent(form.FromAddr, wal.Keys, form.Value)
	if err != nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, "获取账户信息失败"))
		return
	}
	log.Println("amount:", amount, "value:", form.Value)

	if form.Value > amount {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.AccessError, "余额不足"))
		return
	}

	article := &models.Article{}
	if strings.ToLower(form.Type) == "reward" && len(form.Id) > 0 {
		if b, e := article.FindById(form.Id); !b {
			err := errors.NewError(errors.NotExistsError, "文章不存在!")
			if e != nil {
				err = errors.NewError(errors.DbError)
			}
			writeResponse(r.RequestURI, w, nil, err)
			return
		}
	}

	changeAddr := form.FromAddr
	if len(changeAddr) == 0 {
		changeAddr = wal.Keys[0].PubKey
	}
	rawtx, err := CreateRawTx2(outputs, amount, form.Value, form.ToAddr, changeAddr)
	if err != nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError))
		return
	}

	txid, err := sendRawTx(rawtx)
	if err != nil {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError))
		return
	}

	redis.Transaction(user.Id, receiver.Id, form.Value)
	// ws push
	event := &models.Event{
		Type: models.EventWallet,
		Time: time.Now().Unix(),
	}

	//log.Println("tx type:", strings.ToLower(form.Type))

	switch strings.ToLower(form.Type) {
	case "reward":

		article.Reward(user.Id, form.Value)

		event.Data = models.EventData{
			Type: models.EventReward,
			Id:   article.Id.Hex(),
			From: user.Id,
			To:   receiver.Id,
			Body: []models.MsgBody{
				{Type: "nikename", Content: user.Nickname},
				{Type: "image", Content: article.Image},
				{Type: "total_count", Content: strconv.FormatInt(article.TotalReward, 10)},
				{Type: "new_count", Content: strconv.Itoa(models.EventCount(models.EventReward, article.Id.Hex()) + 1)},
			},
		}
	default:
		event.Data = models.EventData{
			Type: models.EventTx,
			Id:   user.Id,
			From: user.Id,
			To:   receiver.Id,
			Body: []models.MsgBody{
				{Type: "nikename", Content: user.Nickname},
				{Type: "image", Content: user.Profile},
				{Type: "total_count", Content: strconv.FormatInt(form.Value, 10)},
			},
		}
	}

	if user.Id != receiver.Id {
		if err := event.Save(); err == nil {
			redis.IncrEventCount(receiver.Id, event.Data.Type, 1)
		}

		event.Data.Body = append(event.Data.Body,
			models.MsgBody{Type: "new_count", Content: strconv.Itoa(models.EventCount(models.EventTx, user.Id))})
		redis.PubMsg("wallet", receiver.Id, event.Bytes())
	}

	if receiver.Push {
		msg := user.Nickname + "通过转帐发送给你" +
			strconv.FormatFloat(float64(form.Value)/float64(models.Satoshi), 'f', 8, 64) + "个贝币"
		if event.Data.Type == models.EventReward {
			msg = user.Nickname + "给你的文章打赏了" +
				strconv.FormatFloat(float64(form.Value)/float64(models.Satoshi), 'f', 8, 64) + "个贝币"
		}
		go sendApn(client, msg, receiver.Devs...)
	}

	writeResponse(r.RequestURI, w, map[string]string{"txid": txid}, nil)
}

type Vin struct {
	Txid     string `json:"txid"`
	Sequence uint32 `json:"sequence"`
	Script   string `json:"script"`
	PrevOut  Vout   `json:"prev_out"`
}

type Vout struct {
	Value      int64  `json:"value"`
	N          uint32 `json:"n"`
	Script     string `json:"script"`
	ScriptType string `json:"type"`
	Address    string `json:"addr"`
}

type Tx struct {
	Hash    string  `json:"tx_hash"`
	Block   string  `json:"block"`
	Version int32   `json:"version"`
	Index   int64   `json:"tx_index"`
	Time    int64   `json:"time"`
	Vin     []*Vin  `json:"inputs"`
	Vout    []*Vout `json:"outputs"`
	Amount  int64   `json:"amount"`
}

func getAddrTxs(addr string) ([]Tx, error) {
	txs := []Tx{}

	if len(addr) == 0 {
		return txs, nil
	}
	resp, err := http.Get(CoinAddr + "/addr_txs?addr=" + addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = decodeJson(resp.Body, &txs)

	return txs, err
}

type txResp struct {
	Error string `json:"error"`
	Txid  string `json:"txid"`
}

func sendRawTx(rawtx string) (txid string, err error) {
	resp, err := http.PostForm(CoinAddr+"/pushtx", url.Values{"rawtx": {rawtx}})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	txr := new(txResp)
	if err = decodeJson(resp.Body, txr); err != nil {
		return
	}

	if len(txr.Error) > 0 {
		err = errs.New(txr.Error)
		return
	}
	return txr.Txid, nil
}

type output struct {
	TxHash  string `json:"tx_hash"`
	TxN     uint32 `json:"tx_output_n"`
	Script  string `json:"script"`
	Value   int64  `json:"value"`
	Address string `json:"address"`
	PrivKey string `json:"-"`
}

type unspent struct {
	Outputs []output `json:"unspent_outputs"`
}

func getUnspent(addr string, keys []*key, value int64) (outputs []output, amount int64, err error) {
	var addrs []string
	var wallet string
	var keyMap = make(map[string]string)
	for _, k := range keys {
		addrs = append(addrs, k.PubKey)
		keyMap[k.PubKey] = k.PrivKey
	}
	/*
		if len(addrs) > 0 {
			/wallet = strings.Join(addrs, "|")
		}
	*/

	if _, ok := keyMap[addr]; ok {
		addrs = []string{addr}
	}
	resp, err := http.Get(CoinAddr + "/unspent?addr=" + addr + "&keys=" + wallet)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	us := new(unspent)
	if err = decodeJson(resp.Body, us); err != nil {
		return
	}

	for _, op := range us.Outputs {
		amount += op.Value
		op.PrivKey = keyMap[op.Address]
		outputs = append(outputs, op)
		if amount >= value {
			break
		}
	}
	return
}

type balance struct {
	Address     string `json:"address"`
	Confirmed   int64  `json:"confirmed"`
	Unconfirmed int64  `json:"unconfirmed"`
}

type balanceAddrs struct {
	Addrs []balance `json:"addresses"`
}

func getBalance(addrs []string) (b *balanceAddrs, err error) {
	if len(addrs) == 0 {
		return
	}
	resp, err := http.Get(CoinAddr + "/multiaddr?addr=" + strings.Join(addrs, "|"))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	ba := new(balanceAddrs)
	if err = decodeJson(resp.Body, ba); err != nil {
		log.Println(err)
		return
	}

	b = ba

	return
}

type key struct {
	PubKey  string `json:"pubKey"`
	PrivKey string `json:"privKey"`
	Label   string `json:"label"`
}

type wallet struct {
	SharedKey string `json:"sharedKey"`
	Keys      []*key `json:"keys"`
}

func NewWallet() *wallet {
	w := &wallet{SharedKey: Uuid()}
	k, err := w.GenKey("")
	if err != nil {
		log.Println(err)
	}
	w.AddKey(k)
	return w
}

func (w *wallet) GenKey(label string) (*key, error) {
	privKey, pubKey, err := btcaddr.Bitcoin_GenerateKeypair()
	if err != nil {
		return nil, err
	}

	return &key{
		PubKey:  btcaddr.Bitcoin_Pubkey2Address(pubKey, 0),
		PrivKey: btcaddr.Bitcoin_Prikey2WIF(privKey),
		Label:   label,
	}, nil
}

func (w *wallet) AddKey(k *key) {
	if k == nil {
		return
	}
	w.Keys = append(w.Keys, k)
}

func getWallet(id, sharedKey string) (*wallet, error) {
	resp, err := http.Get(CoinAddr + "/wallet?wallet_id=" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		Id      string `json:"wallet_id"`
		Payload string `json:"payload"`
	}
	if err = decodeJson(resp.Body, &r); err != nil {
		return nil, err
	}

	return decryptWallet(r.Payload, strings.Join(strings.Split(sharedKey, "-"), ""))
}

func saveWallet(id string, w *wallet) (string, error) {
	s, err := encryptWallet(w, strings.Join(strings.Split(w.SharedKey, "-"), ""))
	if err != nil {
		return id, err
	}

	resp, err := http.PostForm(CoinAddr+"/wallet", url.Values{"wallet_id": {id}, "payload": {s}})
	if err != nil {
		return id, err
	}
	defer resp.Body.Close()

	var r struct {
		WalletId string `json:"wallet_id"`
		Status   string `json:"status"`
	}
	if err = decodeJson(resp.Body, &r); err != nil {
		return id, err
	}
	if r.Status != "ok" {
		return id, errs.New(r.Status)
	}
	return r.WalletId, nil
}

func decryptWallet(payload string, key string) (*wallet, error) {
	w := &wallet{}
	if len(payload) == 0 || len(key) == 0 {
		return w, nil
	}
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	data, err = Decrypt(data, key)
	if err != nil {
		return nil, err
	}
	err = decodeJson(bytes.NewBuffer(data), w)
	return w, err
}

func encryptWallet(w *wallet, key string) (string, error) {
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(w); err != nil {
		return "", err
	}
	enc, err := Encrypt(b.Bytes(), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}
