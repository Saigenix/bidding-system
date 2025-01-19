package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrParsingClaim = errors.New("Failed to parse claims")
	ErrParsingToken = errors.New("Error parsing token")
)

// JWTGenerator is a struct that generates JWT tokens
//
// Usage:
//
// g := NewJWTGenerator(WithSecret([]byte("secret")))
// token, err := g.GenerateToken("userID")
//
// customization for example
// token, err := g.GenerateToken("userID", WithExpiration(10 * time.Minute))
type JWTGenerator struct {
	expiration time.Duration
	secret     []byte
	algorithm  jwt.SigningMethod
	claims     jwt.MapClaims
}

// NewJWTGenerator creates a new JWTGenerator with the given options
//
// Usage:
//
//	g := NewJWTGenerator(
//			[]byte("secret"),
//			WithExpiration(10 * time.Minute),
//	)
//
// Note: Theer are no default values for the secret
func NewJWTGenerator(secret []byte, opts ...JWTGenOpts) JWTGenerator {
	g := JWTGenerator{
		expiration: 10 * time.Minute, // Default expiration time
		algorithm:  jwt.SigningMethodHS256,
		secret:     secret,
	}

	for _, opt := range opts {
		opt(&g)
	}
	return g
}

type JWTGenOpts func(*JWTGenerator)

// WithExpiration sets the secret for the JWTGenerator
func WithExpiration(expiration time.Duration) JWTGenOpts {
	return func(g *JWTGenerator) {
		g.expiration = expiration
	}
}

// WithAlgorithm sets the algorithm for the JWTGenerator
func WithAlgorithm(algorithm jwt.SigningMethod) JWTGenOpts {
	return func(g *JWTGenerator) {
		g.algorithm = algorithm
	}
}

// WithClaims sets the claims for the JWTGenerator
func WithClaims(claims jwt.MapClaims) JWTGenOpts {
	return func(g *JWTGenerator) {
		for key, value := range claims {
			g.claims[key] = value
		}
	}
}

// WithSubject sets the subject for the JWTGenerator
func WithSubject(subject string) JWTGenOpts {
	return func(g *JWTGenerator) {
		g.claims["sub"] = subject
	}
}

// WithIssuer sets the issuer for the JWTGenerator
func WithIssuer(issuer string) JWTGenOpts {
	return func(g *JWTGenerator) {
		g.claims["iss"] = issuer
	}
}

// GenerateToken generates a JWT token for the given userID
func (g JWTGenerator) GenerateToken(userID string, opts ...JWTGenOpts) (string, error) {

	for _, opt := range opts {
		opt(&g)
	}

	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(g.expiration).Unix(),
	}

	token := jwt.NewWithClaims(g.algorithm, claims)

	return token.SignedString(g.secret)

}

// ParseToken parses the given token string and returns the token
//
// Note: This function does not support customizations
func (g JWTGenerator) ParseToken(tokenString string) (ParsedToken, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return g.secret, nil
	})

	if err != nil {
		return ParsedToken{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return ParsedToken{}, errors.Join(ErrParsingToken, ErrParsingClaim)
	}

	return ParsedToken{
		token:  token,
		Valid:  token.Valid,
		claims: claims,
	}, err
}

// ParsedToken is a struct that holds the parsed token
type ParsedToken struct {
	token *jwt.Token
	Valid bool

	claims jwt.MapClaims
}

// UserId returns the userID from the token
func (pt ParsedToken) UserId() string {
	id := pt.claims["userID"].(string)
	return id
}
