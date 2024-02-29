package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserModel struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd,omitempty"`
	IsAdmin  bool   `json:"isadmin"`
	Hashcode string `json:"hashcode"`
}

var users_registred = make([]UserModel, 0)

// Função pra substituir o println da struct UserModel
func (u *UserModel) String() string {
	return fmt.Sprintf("Usuário: %s, Nome: %s, Email: %s ", u.User, u.Name, u.Email)
}

func GetUser(c *gin.Context) {
	//retornando o Array JSON
	c.IndentedJSON(http.StatusOK, users_registred)
}

func RegisterUser(c *gin.Context) {
	var user_new UserModel

	if err := c.BindJSON(&user_new); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Data in wrong format!")
		return
	}

	fmt.Println("User received => ", user_new)

	//[...] Operações do banco de dados [...]
	users_registred = append(users_registred, user_new)

	//retornando o JSON
	c.IndentedJSON(http.StatusCreated, user_new)
}

func main() {
	router := gin.Default()
	router.POST("/user", RegisterUser)
	router.GET("/user", GetUser)

	router.Run("0.0.0.0:8085")
}
