package auth

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const JWTKey = "2Ye5_w1z3zpD4dSGdRp3s98ZipCNQqmFDr9vioOx54"
const JWTExpire = 96
const UserAuthKey contextKey = "user"

type contextKey string

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.StandardClaims
}

var jwtKey = []byte(JWTKey)

func CreateToken(user_id int, role uint32) string {
	// Declare the expiration time of the token here
	expirationTime := time.Now().Add(time.Duration(JWTExpire) * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		UserID:   strconv.Itoa(int(user_id)),
		UserRole: strconv.Itoa(int(role)),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return ""
	}
	return tokenString
}
