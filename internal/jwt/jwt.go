package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
)

const expiration = time.Hour * 24

// NewToken creates an authentication token for the specified user resource.
func NewToken(secret []byte, user *cargonaut.User) (token *cargonaut.Token, err error) {
	// Create unique JWT ID.
	id := uuid.NewV4()

	// Setup the standard claims.
	now := time.Now().UTC()
	exp := now.Add(expiration)
	claims := jws.Claims{}
	claims.SetIssuer("my-cargonaut.com")
	claims.SetSubject("authentication")
	claims.SetAudience("my-cargonaut.com")
	claims.SetExpiration(exp)
	claims.SetNotBefore(now)
	claims.SetIssuedAt(now)
	claims.SetJWTID(id.String())

	// Setup the custom user claims.
	claims.Set("user", map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.DisplayName,
	})

	// Create and serialize JWT.
	tokenStr, err := jws.NewJWT(claims, crypto.SigningMethodHS256).Serialize(secret)
	if err != nil {
		return nil, fmt.Errorf("serialize token: %w", err)
	}

	token = &cargonaut.Token{
		ID:        id,
		UserID:    user.ID,
		Token:     string(tokenStr),
		ExpiresAt: exp,
		CreatedAt: now,
	}
	return token, nil
}

// UserFromToken checks an authentication token for validity and returns the
// user which is associated with it from the token claims. The provided token is
// updated.
func UserFromToken(secret []byte, token *cargonaut.Token) (user *cargonaut.User, err error) {
	// Parse the token string.
	jwt, err := jws.ParseJWT([]byte(token.Token))
	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}

	// Setup the standard claims which are expected to match.
	claims := jws.Claims{}
	claims.SetIssuer("my-cargonaut.com")
	claims.SetSubject("authentication")
	claims.SetAudience("my-cargonaut.com")

	// Setup custom validator function to check the user claims.
	val := func(c jws.Claims) error {
		var ok bool

		// Get standard claims.
		tokenID, ok := c.JWTID()
		if !ok {
			return errors.New("missing token id")
		} else if token.ID, err = uuid.FromString(tokenID); err != nil {
			return fmt.Errorf("parse token id: %w", err)
		}

		expiresAt, ok := c.Expiration()
		if !ok {
			return errors.New("missing expiration time")
		}
		token.ExpiresAt = expiresAt.UTC()

		createdAt, ok := c.IssuedAt()
		if !ok {
			return errors.New("missing creation time")
		}
		token.CreatedAt = createdAt.UTC()

		// Get user claims.
		uClaims, ok := c.Get("user").(map[string]interface{})
		if !ok {
			return errors.New("could not find user claims")
		}
		user = new(cargonaut.User)

		userID, ok := uClaims["id"].(string)
		if !ok {
			return errors.New("missing user id")
		} else if user.ID, err = uuid.FromString(userID); err != nil {
			return fmt.Errorf("parse user id: %w", err)
		}

		if user.Email, ok = uClaims["email"].(string); !ok {
			return errors.New("missing user email")
		}

		if user.DisplayName, ok = uClaims["name"].(string); !ok {
			return errors.New("missing user name")
		}

		return nil
	}

	// Setup validator and validate the token.
	v := jws.NewValidator(claims, expiration, 0, val)
	if err = jwt.Validate(secret, crypto.SigningMethodHS256, v); err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}
	return user, nil
}
