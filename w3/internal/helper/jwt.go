package helper

import (
	"errors"
	"halosuster/configuration"
	"halosuster/domain/entity"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type IJWTManager interface {
	CreateToken(user entity.User) (string, error)
	GetClaim(tokenString string) (*JsonWebTokenClaims, error)
}

type jsonWebToken struct {
	expirationTimeInMinute int
	signingKey             string
	issuer                 string
}

type JsonWebTokenClaims struct {
	UserID string `json:"userId"`
	Role   string `json:"roleId"`
	NIP    int    `json:"nip"`
	jwt.RegisteredClaims
}

func NewJwt(jwtConfiguration configuration.IJWTConfiguration) *jsonWebToken {
	return &jsonWebToken{
		expirationTimeInMinute: jwtConfiguration.GetExpireInMinute(),
		signingKey:             jwtConfiguration.GetSigningKey(),
		issuer:                 jwtConfiguration.GetIssuer(),
	}
}

func (t *jsonWebToken) CreateToken(user entity.User) (string, error) {
	expiresAt := time.Now().Add(time.Duration(t.expirationTimeInMinute) * time.Minute)

	integer, _ := strconv.Atoi(user.NIP)
	claims := JsonWebTokenClaims{
		UserID: user.ID.String(),
		Role:   user.Role,
		NIP:    integer,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    t.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(t.signingKey))

	return ss, err
}

func (j *jsonWebToken) GetClaim(tokenString string) (*JsonWebTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JsonWebTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		} else if method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}

		return []byte(j.signingKey), nil

	}, jwt.WithExpirationRequired(), jwt.WithIssuer(j.issuer))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claim, ok := token.Claims.(*JsonWebTokenClaims)

	if !ok {
		return nil, ErrInvalidToken
	}

	return claim, nil
}
