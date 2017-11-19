package auth

import (
    "errors"
    "github.com/dgrijalva/jwt-go"
    "time"
)

type JwtStore struct {
    tokenKey    []byte
    expireAfter time.Duration
}

type JwtToken struct {
    tokenKey []byte
    jwt.Token
}

func (t *JwtToken) Claims(key string) interface{} {
    claims := t.Token.Claims.(jwt.MapClaims)
    return claims[key]
}

func (t *JwtToken) SetClaim(key string, value interface{}) ClaimSetter {
    claims := t.Token.Claims.(jwt.MapClaims)
    claims[key] = value
    return t
}

func (t *JwtToken) Expiry() time.Time {
    expt := t.Claims("exp")
    var exp time.Time
    switch t := expt.(type) {
    case float64:
        exp = time.Unix(int64(t), 0)
    case int64:
        exp = time.Unix(t, 0)
    default:
        exp = time.Now()
    }
    return exp
}

func (t *JwtToken) IsExpired() bool {
    exp := t.Expiry()
    return time.Now().After(exp)
}

func (t *JwtToken) String() string {
    tokenStr, _ := t.Token.SignedString(t.tokenKey)
    return tokenStr
}

func (s *JwtStore) NewToken(id interface{}) *JwtToken {
    token := jwt.New(jwt.GetSigningMethod("HS256"))
  claims := token.Claims.(jwt.MapClaims)
    claims["exp"] = time.Now().Add(s.expireAfter).Unix()
    t := &JwtToken{
        tokenKey: s.tokenKey,
        Token:    *token,
    }
    return t
}

func (s *JwtStore) CheckToken(token string) (Token, error) {
    t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        return s.tokenKey, nil
    })
    if err != nil {
        return nil, err
    }
    jtoken := &JwtToken{s.tokenKey, *t}
    if jtoken.IsExpired() {
        return nil, errors.New("Token expired")
    }
    return jtoken, nil
}

func NewJwtStore(tokenKey string, expireAfter time.Duration) *JwtStore {
    return &JwtStore{
        []byte(tokenKey),
        expireAfter,
    }
}
