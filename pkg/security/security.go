package security

import (

)

// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(bytes), err
// }

// func CheckPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

// func GenerateJWT(userID int) (string, error) {
// 	claims := &jwt.RegisteredClaims{
// 		// Subject:   string(userID),
// 		Subject:   strconv.Itoa(userID),
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte("SECRET"))
// }

// func GenerateJWT(userID int) (string, error) {
// 	// ペイロードの作成
// 	claims := &jwt.MapClaims{
// 		"sub": userID,
// 		"exp": time.Now().Add(time.Hour * 24).Unix(),
// 	}

// 	// トークン生成
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// トークンに署名を付与
// 	tokenString, err := token.SignedString([]byte("SECRET"))
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokenString, nil
// }


// この関数使われてないかも
// func ParseJWT(tokenStr string) (*jwt.RegisteredClaims, error) {
// 	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
// 		return claims, nil
// 	}
// 	return nil, errors.New("invalid token")
// }
