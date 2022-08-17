package data

import (
	"fmt"
	"log"
	"middleman-capstone/domain"

	"gorm.io/gorm"
)

type inventoryData struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.InventoryData {
	return &inventoryData{
		db: db,
	}
}

func (ind *inventoryData) CekStok(newRecap []domain.InventoryProduct, id int, gen string) bool {
	var product = FromIP2(newRecap, id, gen)
	something := domain.ProductUser{}
	for _, val := range product {
		res2 := ind.db.Model(&domain.ProductUser{}).Select("stock").Where("id = ? AND id_user = ?", val.ProductID, id).First(&something)
		if res2.Error != nil {
			log.Println("Cannot retrieve object", res2.Error.Error())
			return false
		}
		cekstok := something.Stock - val.Qty
		if cekstok < 0 {
			return false
		}
	}
	return true
}

func (ind *inventoryData) CreateUserInventoryData(newRecap domain.Inventory, id int, gen string) domain.Inventory {
	var inven = FromModel(newRecap, id, gen)
	fmt.Println("inventories", inven)
	res := ind.db.Create(inven)

	if res.Error != nil {
		log.Println("cannot create data")
		return domain.Inventory{}
	}

	if res.RowsAffected == 0 {
		log.Println("failed to insert data")
		return domain.Inventory{}
	}
	return inven.ToI()
}

// func (ind *inventoryData) CreateUserInventoryData(newRecap domain.InventoryProduct) domain.InventoryProduct {
func (ind *inventoryData) CreateUserDetailInventoryData(newRecap []domain.InventoryProduct, id int, gen string) []domain.InventoryProduct {
	var product = FromIP2(newRecap, id, gen)
	err := ind.db.Create(product)

	if err.Error != nil {
		log.Println("cannot create data", err.Error.Error())
		return []domain.InventoryProduct{}
	}

	if err.RowsAffected == 0 {
		log.Println("failed to insert data")
		return []domain.InventoryProduct{}
	}
	return ParsePUToArr(product)
}

func (ind *inventoryData) RekapStock(newRecap []domain.InventoryProduct, id int, gen string) bool {
	var product = FromIP2(newRecap, id, gen)
	something := InventoryProduct{}
	for _, val := range product {
		res2 := ind.db.Model(&InventoryProduct{}).Select("inventory_products.user_id, inventory_products.product_id, product_users.name, inventory_products.qty, inventory_products.unit, product_users.stock").Joins("left join product_users on product_users.id = inventory_products.product_id").Where("product_id = ? AND idip = ?", val.ProductID, gen).First(&something)
		if res2.Error != nil {
			log.Println("Cannot retrieve object", res2.Error.Error())
			return false
		}
		res3 := ind.db.Model(&InventoryProduct{}).Where("product_id = ? AND user_id = ? AND idip = ?", val.ProductID, id, gen).Updates(&something)
		if res3.Error != nil {
			log.Println("Cannot retrieve object", res2.Error.Error())
			return false
		}
		stock := something.Stock - something.Qty
		res4 := ind.db.Model(&domain.ProductUser{}).Where("id = ? AND id_user = ?", val.ProductID, id).Update("stock", stock)
		if res4.Error != nil {
			log.Println("Cannot retrieve object", res2.Error.Error())
			return false
		}
	}
	return true
}

func (ind *inventoryData) DeleteInOutBound(id int) (err string) {
	res := ind.db.Where("id_user = ?", id).Delete(&domain.InOutBounds{})

	if res.Error != nil {
		log.Println("cannot delete data")
		return "cannot delete data"
	}
	if res.RowsAffected < 1 {
		log.Println("no data deleted")
		return "no data deleted"
	}
	return ""
}

func (ind *inventoryData) ReadUserOutBoundDetailData(id int, outboundIDGenerate string) []domain.InventoryProduct {
	var product []InventoryProduct
	err := ind.db.Where("user_id = ? AND idip = ?", id, outboundIDGenerate).Find(&product)
	if err.Error != nil {
		log.Println("cannot read data", err.Error.Error())
		return []domain.InventoryProduct{}
	}
	if err.RowsAffected == 0 {
		log.Println("data not found")
		return []domain.InventoryProduct{}
	}
	return ParsePUToArr(product)
}
