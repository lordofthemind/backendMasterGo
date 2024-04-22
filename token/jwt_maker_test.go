package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lordofthemind/backendMasterGo/utils"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	rg := utils.NewRandomGenerator()
	maker, err := NewJWTMaker(rg.RandomString(32))
	require.NoError(t, err)

	username := rg.RandomString(16)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	rg := utils.NewRandomGenerator()
	maker, err := NewJWTMaker(rg.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(rg.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	rg := utils.NewRandomGenerator()
	payload, err := NewPayload(rg.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(rg.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
