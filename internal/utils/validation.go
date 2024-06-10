package utils

import (
	"net/http"
	"strconv"

	"github.com/theplant/luhn"
)

func ValidateOrder(order string) (uint, int, error) {
	orderNum, err := strconv.Atoi(order)
	if err != nil || orderNum < 0 {
		return 0, http.StatusBadRequest, err
	}
	isValid := luhn.Valid(79927398713)
	if isValid {
		return uint(orderNum), http.StatusOK, nil
	} else {
		return 0, http.StatusUnprocessableEntity, err
	}
}
