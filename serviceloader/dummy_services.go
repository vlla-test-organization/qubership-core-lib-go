package serviceloader

import (
	"context"

	"github.com/golang-jwt/jwt"

	"github.com/netcracker/qubership-core-lib-go/v3/logging"
)

var logger logging.Logger

type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
	ValidateToken() func(ctx context.Context, token string) (*jwt.Token, error)
	GetClaimValue(token *jwt.Token, key string) (interface{}, error)
	IsSecured() bool
}

type Token struct {
}

type TlsConfig interface {
	IsTlsEnabled() bool
}

func init() {
	logger = logging.GetLogger("dummy-services")
}

func (s *Token) GetToken(ctx context.Context) (string, error) {
	logger.Info("Empty token value implementation ")
	return "", nil
}

func (s *Token) GetClaimValue(token *jwt.Token, key string) (interface{}, error) {
	logger.Info("Claim value 'nil' sent for key [%s] from dummy service", key)
	return nil, nil
}

func (s *Token) ValidateToken() func(ctx context.Context, token string) (*jwt.Token, error) {
	logger.Info("Token parsed unverified")
	return func(ctx context.Context, token string) (*jwt.Token, error) {
		parser := jwt.Parser{}
		parsedToken, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
		return parsedToken, err
	}
}

func (s *Token) IsSecured() bool {
	logger.Info("security options is not enabled by default")
	return false
}
