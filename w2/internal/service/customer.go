package service

import (
	"eniqlostore/commons"
	"eniqlostore/internal/entity"
	"fmt"
	"net/http"
)

type ICustomerRepository interface {
	Insert(entity.Customer) (entity.Customer, error)
	CheckExistByPhoneNumber(phoneNumber string) (isExist bool, err error)
	FindBy(name string, phoneNumber string) ([]entity.Customer, error)
}

type CustomerService struct {
	customerRepository ICustomerRepository
}

type CreateCustomerRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

type GetCustomerRequst struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

func NewCustomerService(custRepo ICustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepository: custRepo,
	}
}

func (c *CustomerService) CreateCustomer(in CreateCustomerRequest) (entity.Customer, error) {
	var emptyCustomer entity.Customer

	if err := entity.PhoneNumber(in.PhoneNumber).Valid(); err != nil {
		return emptyCustomer, commons.CustomError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	isExist, err := c.customerRepository.CheckExistByPhoneNumber(in.PhoneNumber)

	if err != nil {
		return emptyCustomer, err
	}

	if isExist {
		return emptyCustomer, commons.CustomError{
			Message: "customer phone number already exist",
			Code:    http.StatusConflict,
		}
	}

	newCustomer, err := c.customerRepository.Insert(entity.Customer{
		Name:        in.Name,
		PhoneNumber: in.PhoneNumber,
	})

	if err != nil {
		return emptyCustomer, err
	}

	return newCustomer, nil
}

func (c *CustomerService) GetCustomers(getOpts GetCustomerRequst) (customers []entity.Customer, err error) {
	phoneNumberWithPrefixSign := fmt.Sprintf("+%s", getOpts.PhoneNumber)
	customers, err = c.customerRepository.FindBy(getOpts.Name, phoneNumberWithPrefixSign)

	return customers, err
}
