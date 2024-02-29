package userhandle

import (
	"log"
	"net/http"
	"strconv"
	user "webserver/models/usermodel"

	"github.com/gin-gonic/gin"
)

// Retorna todos os usuários ou o usuário com uma determinada id (ex.: http://rota?id=valor)
func GetUser(c *gin.Context) {
	var u = user.UserModel{}
	var idquery, err = strconv.ParseInt(c.Query("id"), 10, 32)

	if err == nil {
		u.Id = int(idquery)
		allusers, err := u.GetByID()

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Error on select user!")
			return
		}
		c.IndentedJSON(http.StatusOK, allusers)

	} else {

		allusers, err := u.GetAll()

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Error on select all users!")
			return
		}
		c.IndentedJSON(http.StatusOK, allusers)
	}
}

// Insere ou edita um usuário (se o id vier na requisição do post ele editará o usuário)
func RegisterUser(c *gin.Context) {
	var user_new *user.UserModel

	if err := c.BindJSON(&user_new); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Data in wrong format!")
		return
	}

	log.Println("User received => ", user_new)

	user_new, err := user_new.Save()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Error on insert/update user!")
		return
	}

	//retornando o JSON
	c.IndentedJSON(http.StatusCreated, user_new)
}

func DeleteUser(c *gin.Context) {
	var u = user.UserModel{}
	var idquery, err = strconv.ParseInt(c.Query("id"), 10, 32)

	if err == nil {
		u.Id = int(idquery)
		resp, err := u.Delete()

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Error on delete user!")
			return
		}
		c.IndentedJSON(http.StatusOK, resp)

	}
}
