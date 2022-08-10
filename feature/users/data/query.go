package data

import (
	"errors"
	"fmt"
	"log"
	"middleman-capstone/domain"
	"middleman-capstone/feature/common"
	user "middleman-capstone/feature/users"

	_bcrypt "golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type userData struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.UserData {
	return &userData{
		db: db,
	}
}

func (ud *userData) Insert(newUser domain.User) (row int, err error) {
	user := FromModel(newUser)
	passwordHashed, errorHash := _bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if errorHash != nil {
		fmt.Println("Error hash", errorHash.Error())
	}
	user.Password = string(passwordHashed)
	resultCreate := ud.db.Create(&user)
	if resultCreate.Error != nil {
		return 0, resultCreate.Error
	}
	if resultCreate.RowsAffected != 1 {
		return 0, errors.New("failed to insert data, your email is already registered")
	}
	return int(resultCreate.RowsAffected), nil
}

func (ud *userData) LoginData(authData user.LoginModel) (data map[string]interface{}, err error) {
	userData := User{}
	result := ud.db.Where("email = ?", authData.Email).First(&userData)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected != 1 {
		return nil, errors.New("failed to login")
	}

	errCrypt := _bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(authData.Password))

	if errCrypt != nil {
		return nil, errors.New("password incorrect")
	}
	token := common.GenerateToken(int(userData.ID), userData.Role)

	var dataToken = map[string]interface{}{}
	dataToken["id"] = int(userData.ID)
	dataToken["name"] = userData.Name
	dataToken["email"] = userData.Email
	dataToken["role"] = userData.Role
	dataToken["token"] = token
	return dataToken, nil
}
func (ud *userData) GetSpecific(userID int) (domain.User, error) {
	var tmp User
	err := ud.db.Where("ID = ?", userID).First(&tmp).Error
	if err != nil {
		log.Println("There is a problem with data", err.Error())
		return domain.User{}, err
	}

	return tmp.ToModel(), nil
}
func (ud *userData) DeleteData(userID int) (row int, err error) {
	res := ud.db.Delete(&User{}, userID)
	if res.Error != nil {
		log.Println("cannot delete data", res.Error.Error())
		return 0, res.Error
	}
	if res.RowsAffected < 1 {
		log.Println("no data deleted", res.Error.Error())
		return 0, errors.New("failed to delete data ")
	}
	return int(res.RowsAffected), nil
}
func (ud *userData) UpdateData(data map[string]interface{}, idFromToken int) (row int, err error) {
	res := ud.db.Model(&User{}).Where("id = ?", idFromToken).Updates(data)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected != 1 {
		return 0, errors.New("failed update data")
	}
	return int(res.RowsAffected), nil
}

func (ud *userData) CreateProductData(newProduct domain.ProductUser) domain.ProductUser {
	var product = FromPU(newProduct)
	err := ud.db.Create(&product)

	if err.Error != nil {
		log.Println("cannot create data", err.Error.Error())
		return domain.ProductUser{}
	}

	if err.RowsAffected == 0 {
		log.Println("failed to insert data")
		return domain.ProductUser{}
	}
	return product.ToPU()
}

func (ud *userData) ReadAllProductData(id int) []domain.ProductUser {
	var product []ProductUser
	err := ud.db.Where("id_user = ?", id).Find(&product)
	if err.Error != nil {
		log.Println("cannot read data", err.Error.Error())
		return []domain.ProductUser{}
	}
	if err.RowsAffected == 0 {
		log.Println("data not found")
		return []domain.ProductUser{}
	}
	return ParsePUToArr(product)
}

func (ud *userData) UpdateProductData(updatedData domain.ProductUser) domain.ProductUser {
	var products = FromPU(updatedData)
	err := ud.db.Model(&ProductUser{}).Where("id = ? AND id_user = ?", products.ID, products.IdUser).Updates(products)

	if err.Error != nil {
		log.Println("cannot update data", err.Error.Error())
		return domain.ProductUser{}
	}

	if err.RowsAffected == 0 {
		log.Println("data not found")
		return domain.ProductUser{}
	}

	return products.ToPU()
}

func (ud *userData) DeleteProductData(productid, id int) (row int, err error) {
	res := ud.db.Where("id = ? AND id_user = ?", productid, id).Delete(&ProductUser{})

	if res.Error != nil {
		log.Println("cannot delete data", res.Error.Error())
		return 0, res.Error
	}
	if res.RowsAffected < 1 {
		log.Println("no data deleted", res.Error.Error())
		return 0, errors.New("failed to delete data ")
	}
	return int(res.RowsAffected), nil
}
