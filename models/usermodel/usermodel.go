package usermodel

import (
	"encoding/json"
	"fmt"
	"time"
	db "webserver/library/database"
)

type UserModel struct {
	Id       int    `json:"id,omitempty,string"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd,omitempty"`
	Hashcode string `json:"hashcode"`
	Dtupdate string `json:"dtupdate"`
	Dtcreate string `json:"dtcreate"`
}

func (u UserModel) String() string {
	b, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return fmt.Sprintf(string(b))
}

func (u UserModel) Save() (*UserModel, error) {
	if u.Passwd != "" && u.Email != "" {
		u.Hashcode = db.GetHash(u.Email, u.Passwd)
		u.Passwd = ""
	}
	t := time.Now().UTC()
	u.Dtupdate = t.Format("2006-01-02 15:04:05")

	if u.Id == 0 {
		response, err := db.Set[UserModel]("user", u)
		if err != nil {
			fmt.Println("Error on insert!", err.Error())
			return nil, err
		}
		u.Id = int(response)
		resp, err := u.GetByID()

		if err != nil {
			fmt.Println("Error on Update!", err.Error())
			return nil, err
		}
		return resp, nil
	}

	_, err := db.Update[UserModel]("user", u, fmt.Sprintf("id = '%d'", u.Id))
	if err != nil {
		fmt.Println("Error on Update!", err.Error())
		return nil, err
	}

	fmt.Println("Id => ", u.Id)

	resp, err := u.GetByID()

	if err != nil {
		fmt.Println("Error on Update!", err.Error())
		return nil, err
	}

	return resp, nil

}

func (u *UserModel) GetAll() (*[]UserModel, error) {
	response, err := db.GetAll[UserModel]("user")

	if err != nil {
		fmt.Println("Error on GetAll!")
		return nil, err
	}

	return &response, nil
}

func (u *UserModel) GetByID() (*UserModel, error) {
	response, err := db.GetFirst[UserModel]("user", &db.GetProps{Where: fmt.Sprintf("id='%d'", u.Id)})

	if err != nil {
		fmt.Println("Error on GetAll!")
		return nil, err
	}
	return response, nil
}

func (u *UserModel) Delete() (bool, error) {
	response, err := db.Delete[UserModel]("user", fmt.Sprintf("id='%d'", u.Id))
	if err != nil {
		fmt.Println("Error on GetAll!")
		return false, err
	}

	return response > 0, nil
}
