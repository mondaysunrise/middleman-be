package usecase

import (
	"log"
	"middleman-capstone/domain"
	_data "middleman-capstone/feature/orders/data"
	_helper "middleman-capstone/helper"
	"strconv"

	"github.com/go-playground/validator"
)

type orderUseCase struct {
	orderData domain.OrderData
	validate  *validator.Validate
}

func New(od domain.OrderData, v *validator.Validate) domain.OrderUseCase {
	return &orderUseCase{
		orderData: od,
		validate:  v,
	}
}

func (oc *orderUseCase) GetAllAdmin(limit, offset int) (data []domain.Order, err error) {
	data, err = oc.orderData.SelectDataAdminAll(limit, offset)
	return data, err
}

func (oc *orderUseCase) GetAllUser(limit, offset, idUser int) (data []domain.Order, err error) {
	data, err = oc.orderData.SelectDataUserAll(limit, offset, idUser)
	return data, err
}
func (oc *orderUseCase) GetDetail(idUser, idOrder int) (grandTotal int, err error) {
	grandTotal, err = oc.orderData.GetDetailData(idUser, idOrder)
	if grandTotal == 0 {
		log.Println("error get data")
		return -1, nil
	}
	if err != nil {
		log.Println("failed to get data")
		return 400, nil
	}
	return grandTotal, nil
}
func (oc *orderUseCase) GetItems(idOrder int) (data []domain.Items, err error) {
	data, err = oc.orderData.GetDetailItems(idOrder)

	if err != nil {
		log.Println("failed to get data")
		return []domain.Items{}, nil
	}
	return data, nil
}

func (oc *orderUseCase) CreateOrder(dataOrder domain.Order, idUser int) int {

	validError := oc.validate.Var(dataOrder, "gt=0")
	if validError != nil {
		log.Println("Validation error : ", validError)
		return 400
	}

	idOrder, err := oc.orderData.Insert(dataOrder)
	if err != nil {
		log.Println("error adding data")
		return 400
	}

	create := oc.orderData.InsertData(dataOrder.Items, idOrder)
	if len(create) == 0 {
		log.Println("error after creating data")
		return 500
	}

	return 201
}

func (oc *orderUseCase) Payment(grandTotal, idUser int) (orderName, url, token string, dataUser domain.User) {
	user, _ := oc.orderData.GetUser(idUser)
	data := _data.OrderPayment{
		Email:      user.Email,
		Name:       user.Name,
		GrandTotal: grandTotal,
	}
	data.Phone, _ = strconv.Atoi(user.Phone)

	orderIDGen, trans := _helper.Payment(data)

	return orderIDGen, trans.RedirectURL, trans.Token, user
}

// func (oc *orderUseCase) CreateItems(data []domain.Items) (row int, err error) {
// 	res, err := oc.orderData.CreateItems(data, 1)
// 	return res, err
// }

func (oc *orderUseCase) AcceptPayment(data domain.PaymentWeb) (row int, err error) {
	row, err = oc.orderData.AcceptPaymentData(data)

	if row < 1 {
		log.Println("error payment")
		return -1, err
	}
	return row, err
}

func (oc *orderUseCase) ConfirmOrder(ordername string, userid int) (domain.Order, int) {
	order := oc.orderData.ConfirmOrderData(ordername)
	user, _ := oc.orderData.GetUser(userid)
	totalPayment := strconv.Itoa(order.GrandTotal)
	data := _helper.Recipient{
		OrderID:      ordername,
		Name:         user.Name,
		Email:        user.Email,
		Handphone:    user.Phone,
		TotalPayment: totalPayment,
	}
	if order.OrderName == "" {
		log.Println("Empty Data")
		return domain.Order{}, 404
	} else {
		_helper.SendEmail(data)
	}

	return order, 200
}

func (oc *orderUseCase) DoneOrder(ordername string) (domain.Order, int) {
	order := oc.orderData.DoneOrderData(ordername)
	if order.OrderName == "" {
		log.Println("Empty Data")
		return domain.Order{}, 404
	}

	updateadminstok := oc.orderData.UpdateStokAdmin(ordername)
	if !updateadminstok {
		log.Println("failed update data")
		return domain.Order{}, 500
	}

	cekown, id := oc.orderData.CekUser(ordername)
	if len(cekown) < 1 {
		log.Println("failed retrieve data")
		return domain.Order{}, 500
	}

	for _, val := range cekown {
		owned := oc.orderData.CekOwned(val, id)
		if !owned {
			product := oc.orderData.CreateNewProduct(val, id)
			if !product {
				log.Println("failed create data")
				return domain.Order{}, 500
			}
		} else {
			product := oc.orderData.UpdateNewProduct(val, id)
			if !product {
				log.Println("failed update data")
				return domain.Order{}, 500
			}
		}
	}

	delete := oc.orderData.DeleteCart(id)
	if !delete {
		log.Println("failed delete data")
		return domain.Order{}, 500
	}

	return order, 200
}
