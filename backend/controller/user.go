package controller

import (
	"blog/service"
	"blog/utils"
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
)

var networks = map[string]*liteapi.Client{}
var knownHashes = make(map[string]wallet.Version)

func init() {
	var err error
	networks["-239"], err = liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		logrus.Fatal(err)
	}
	networks["-3"], err = liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		logrus.Fatal(err)
	}
	for i := wallet.Version(0); i <= wallet.V4R2; i++ {
		ver := wallet.GetCodeHashByVer(i)
		knownHashes[hex.EncodeToString(ver[:])] = i
	}
	userController = newUser()
}

type user struct {
	secret     string
	payloadTtl time.Duration
}

func newUser() *user {
	return &user{
		secret:     viper.GetString("secret"),
		payloadTtl: time.Duration(viper.GetInt("payloadttlsec")) * time.Second,
	}
}

var userController *user

func GetUser() *user {
	return userController
}
func (u *user) PayloadHandler(c *gin.Context) {
	payload, err := generatePayload(u.secret, u.payloadTtl)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("系统出错"))
		return
	}

	c.JSON(http.StatusOK, struct {
		Payload string `json:"payload"`
	}{
		Payload: payload,
	})
}

func (u *user) ProofHandler(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	var tp utils.TonProof
	err = json.Unmarshal(b, &tp)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	parsed, err := convertTonProofMessage(c, &tp)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("验证出错"))
		return
	}

	net := networks[tp.Network]
	if net == nil {
		c.JSON(200, utils.NewFailedResponse("验证出错"))
		return
	}
	var hexpk []byte = make([]byte, 32)
	_, err = hex.Decode(hexpk, []byte(tp.PublicKey))
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	check, err := checkProof(c, ed25519.PublicKey(hexpk), net, parsed)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("验证出错"))
		return
	}
	if !check {
		c.JSON(200, utils.NewFailedResponse("验证出错"))
		return
	}

	claims := &utils.JwtCustomClaims{
		Address:   tp.Address,
		PublicKey: tp.PublicKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(u.secret))
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("验证出错"))
		return
	}
	err = service.GetUser().AutoCreateIfNotExist(c, tp.PublicKey, tp.PublicKey)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("登陆失败"))
		return

	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(t))
}

const (
	tonProofPrefix   = "ton-proof-item-v2/"
	tonConnectPrefix = "ton-connect"
)

func signatureVerify(pubkey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(pubkey, message, signature)
}

func convertTonProofMessage(ctx context.Context, tp *utils.TonProof) (*utils.ParsedMessage, error) {
	addr, err := tongo.ParseAddress(tp.Address)
	if err != nil {
		return nil, err
	}

	var parsedMessage utils.ParsedMessage

	sig, err := base64.StdEncoding.DecodeString(tp.Proof.Signature)
	if err != nil {
		logrus.Errorf("convert to ton proof message err:%s", err.Error())
		return nil, err
	}
	parsedMessage.Workchain = addr.ID.Workchain
	parsedMessage.Address = addr.ID.Address[:]
	parsedMessage.Domain = tp.Proof.Domain
	parsedMessage.Timstamp = tp.Proof.Timestamp
	parsedMessage.Signature = sig
	parsedMessage.Payload = tp.Proof.Payload
	parsedMessage.StateInit = tp.Proof.StateInit
	return &parsedMessage, nil
}

func createMessage(ctx context.Context, message *utils.ParsedMessage) ([]byte, error) {
	wc := make([]byte, 4)
	binary.BigEndian.PutUint32(wc, uint32(message.Workchain))

	ts := make([]byte, 8)
	binary.LittleEndian.PutUint64(ts, uint64(message.Timstamp))

	dl := make([]byte, 4)
	binary.LittleEndian.PutUint32(dl, message.Domain.LengthBytes)
	m := []byte(tonProofPrefix)
	m = append(m, wc...)
	m = append(m, message.Address...)
	m = append(m, dl...)
	m = append(m, []byte(message.Domain.Value)...)
	m = append(m, ts...)
	m = append(m, []byte(message.Payload)...)

	messageHash := sha256.Sum256(m)
	fullMes := []byte{0xff, 0xff}
	fullMes = append(fullMes, []byte(tonConnectPrefix)...)
	fullMes = append(fullMes, messageHash[:]...)
	res := sha256.Sum256(fullMes)
	return res[:], nil
}

func checkProof(ctx context.Context, pubKey ed25519.PublicKey, net *liteapi.Client, tonProofReq *utils.ParsedMessage) (bool, error) {

	if time.Now().After(time.Unix(tonProofReq.Timstamp, 0).Add(time.Duration(viper.GetInt("ton.prooflifetimesec")) * time.Second)) {
		msgErr := "proof has been expired"
		logrus.Error(msgErr)
		return false, fmt.Errorf(msgErr)
	}

	if tonProofReq.Domain.Value != viper.GetString("ton.domain") {
		msgErr := fmt.Sprintf("wrong domain: %v", tonProofReq.Domain)
		logrus.Error(msgErr)
		return false, fmt.Errorf(msgErr)
	}

	mes, err := createMessage(ctx, tonProofReq)
	if err != nil {
		logrus.Errorf("create message error: %v", err)
		return false, err
	}

	return signatureVerify(pubKey, mes, tonProofReq.Signature), nil
}
func parseStateInit(stateInit string) ([]byte, error) {
	cells, err := boc.DeserializeBocBase64(stateInit)
	if err != nil || len(cells) != 1 {
		return nil, err
	}
	var state tlb.StateInit
	err = tlb.Unmarshal(cells[0], &state)
	if err != nil {
		return nil, err
	}
	if !state.Data.Exists || !state.Code.Exists {
		return nil, fmt.Errorf("empty init state")
	}
	codeHash, err := state.Code.Value.Value.HashString()
	if err != nil {
		return nil, err
	}
	version, prs := knownHashes[codeHash]
	if !prs {
		return nil, fmt.Errorf("unknown code hash")
	}
	var pubKey tlb.Bits256
	switch version {
	case wallet.V1R1, wallet.V1R2, wallet.V1R3, wallet.V2R1, wallet.V2R2:
		var data wallet.DataV1V2
		err = tlb.Unmarshal(&state.Data.Value.Value, &data)
		if err != nil {
			return nil, err
		}
		pubKey = data.PublicKey
	case wallet.V3R1, wallet.V3R2, wallet.V4R1, wallet.V4R2:
		var data wallet.DataV3
		err = tlb.Unmarshal(&state.Data.Value.Value, &data)
		if err != nil {
			return nil, err
		}
		pubKey = data.PublicKey
	}

	return pubKey[:], nil
}

func compareStateInitWithAddress(a ton.Address, stateInit string) (bool, error) {
	cells, err := boc.DeserializeBocBase64(stateInit)
	if err != nil || len(cells) != 1 {
		return false, err
	}
	h, err := cells[0].Hash()
	if err != nil {
		return false, err
	}
	return bytes.Equal(h, a.ID.Address[:]), nil
}
func generatePayload(secret string, ttl time.Duration) (string, error) {
	payload := make([]byte, 16, 48)
	_, err := rand.Read(payload[:8])
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}
	binary.BigEndian.PutUint64(payload[8:16], uint64(time.Now().Add(ttl).Unix()))
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	payload = h.Sum(payload)
	return hex.EncodeToString(payload[:32]), nil
}

func getWalletPubKey(ctx context.Context, address ton.Address, net *liteapi.Client) (ed25519.PublicKey, error) {
	_, result, err := abi.GetPublicKey(ctx, net, address.ID)
	if err != nil {
		return nil, err
	}
	if r, ok := result.(abi.GetPublicKeyResult); ok {
		i := big.Int(r.PublicKey)
		b := i.Bytes()
		if len(b) < 24 || len(b) > 32 { //govno kakoe-to
			return nil, fmt.Errorf("invalid publock key")
		}
		return append(make([]byte, 32-len(b)), b...), nil //make padding if first bytes are empty
	}
	return nil, fmt.Errorf("can't get publick key")
}

func checkPayload(payload, secret string) error {
	b, err := hex.DecodeString(payload)
	if err != nil {
		return err
	}
	if len(b) != 32 {
		return fmt.Errorf("invalid payload length")
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(b[:16])
	sign := h.Sum(nil)
	if subtle.ConstantTimeCompare(b[16:], sign[:16]) != 1 {
		return fmt.Errorf("invalid payload signature")
	}
	if time.Since(time.Unix(int64(binary.BigEndian.Uint64(b[8:16])), 0)) > 0 {
		return fmt.Errorf("payload expired")
	}
	return nil
}
