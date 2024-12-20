package services

import (
	"projectGO/database"
	"projectGO/models"

	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.StandardClaims
}

func Register(c *fiber.Ctx) error {
var data map[string] string
	if err := c.BodyParser(&data); err != nil{
		return err
	}

	if (data["password"] != data["password_confirm"]){
		c.Status(400);
		return c.JSON(fiber.Map{
			"message": "Passwords dont match!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)


	user := models.User{
		FirstName: data["first_name"],
		LastName: data["last_name"],
		Email: data["email"],
		Password: password,
	}

	database.DB.Create(&user)

	return c.JSON(user)
}	

func Login(c *fiber.Ctx) error {
var data map[string] string
	if err := c.BodyParser(&data); err != nil{
		return err
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)

	
	if user.ID == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Incorrect password!",
		})
	}

	claims := jwt.RegisteredClaims{
		Issuer:        strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	
	return c.JSON(fiber.Map{
		"jwt": token,
	})

}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	claims := token.Claims.(*Claims)
	id := claims.Issuer

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		// Retorna erro se o usuário não for encontrado
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	filteredUser := fiber.Map{
		"id": user.ID,
        "first_name": user.FirstName,
        "last_name":  user.LastName,
    }
	return c.JSON(filteredUser)
}

func Logout(c *fiber.Ctx) error {
    cookie := fiber.Cookie{
        //Definir o valor do cookie como vazio e adicionar uma data de expiração no passado.
        Name: "jwt",
        Value: "",
        //No código da função de logout, remover o cookie definindo o mesmo cookie no passado ( '-' ).
        Expires: time.Now().Add(-time.Hour),
        HTTPOnly: true,
    }
    c.Cookie(&cookie)
    //Retornar uma resposta de sucesso em formato JSON.
    return c.JSON(fiber.Map{
        "message": "success",
    })
}