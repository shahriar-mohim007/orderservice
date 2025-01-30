package utilis

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/rand"
	"time"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateOrderConsignmentID(prefix string) string {

	date := time.Now().Format("060102")

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomPart := make([]byte, 6)
	for i := range randomPart {
		randomPart[i] = charset[rand.Intn(len(charset))]
	}

	return fmt.Sprintf("%s%s%s", prefix, date, string(randomPart))
}

func CalculateDeliveryFee(cityID int, weight float64) float64 {
	baseFee := 60.0
	if cityID != 1 {
		baseFee = 100.0
	}
	if weight > 0.5 && weight <= 1.0 {
		return baseFee + 10.0
	} else if weight > 1.0 {
		return baseFee + 10.0 + 15.0*math.Ceil(weight-1.0)
	}
	return baseFee
}
