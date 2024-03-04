package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-auth/models"
	"go-auth/utils"
	"time"
)

var jwtKey = []byte("my_secret_key")

func SignUp(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}
	query := `
    	SELECT *
		FROM auth_user
    	WHERE email = $1
    `
	rows, err := models.DB.Query(query, user.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var id, name, email, password, role string
	for rows.Next() {
		err := rows.Scan(&id, &name, &email, &password, &role)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}
	if rows.Err() != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if id != "" {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}
	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		c.JSON(500, gin.H{"error": "could not generate password hash"})
		return
	}
	query = `
		INSERT INTO auth_user (name, email, password, role)
		VALUES ($1, $2, $3, $4)
    `
	_, err = models.DB.Exec(query, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": "user created"})
}

func Home(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	claims, err := utils.ParserToken(cookie)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(200, gin.H{"success": "home page", "role": claims.Role})
}

func Premium(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	claims, err := utils.ParserToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	if claims.Role != "admin" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(200, gin.H{"success": "premium page", "role": claims.Role})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged out"})
}

func Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var existingUser models.User
	query := `
		SELECT *
		FROM auth_user
		WHERE email = $1
    `
	rows, err := models.DB.Query(query, user.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var id, name, email, password, role string
	for rows.Next() {
		err := rows.Scan(&id, &name, &email, &password, &role)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return

		}
	}
	existingUser.ID = id
	existingUser.Email = email
	existingUser.Name = name
	existingUser.Password = password
	existingUser.Role = role
	if rows.Err() != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}
	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		c.JSON(400, gin.H{"error": "invalid password"})
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return

	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged in"})
}
