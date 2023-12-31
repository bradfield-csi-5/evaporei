package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

// 4 bytes
type UserId int

type UserData struct {
	ages []uint8
	amounts []uint32 // paymentCents
}

// 12 bytes at least,  but probably 16 bytes w/ padding
type Address struct {
	// 8 bytes (ptr, word)
	fullAddress string
	// 4 bytes
	zip         int
}

// 16 bytes
type DollarAmount struct {
	dollars, cents uint64
}

// 40 bytes
type Payment struct {
	amount DollarAmount
	// 24 bytes
	time   time.Time
}

// 4 + 8 + 4 + 12 + 12
// 40 bytes
type User struct {
	id       UserId
	name     string
	age      int
	address  Address
	// 12 bytes
	payments []Payment
}

func AverageAge(users UserData) float64 {
	sum0, sum1, sum2, sum3 := uint64(0), uint64(0), uint64(0), uint64(0)
	for i := 0 ; i < len(users.ages) / 4 * 4 ; i += 4 {
		sum0 += uint64(users.ages[i])
		sum1 += uint64(users.ages[i + 1])
		sum2 += uint64(users.ages[i + 2])
		sum3 += uint64(users.ages[i + 3])
	}
	return float64(sum0 + sum1 + sum2 + sum3) / float64(len(users.ages))
}

func AveragePaymentAmount(users UserData) float64 {
	sum0, sum1, sum2, sum3 := uint64(0), uint64(0), uint64(0), uint64(0)
	for i := 0 ; i < len(users.amounts) / 4 * 4 ; i += 4 {
		sum0 += uint64(users.amounts[i])
		sum1 += uint64(users.amounts[i + 1])
		sum2 += uint64(users.amounts[i + 2])
		sum3 += uint64(users.amounts[i + 3])
	}
	return 0.01 * float64(sum0 + sum1 + sum2 + sum3) / float64(len(users.amounts))
}

// Compute the standard deviation of payment amounts
// Variance[X] = E[X^2] - E[X]^2
func StdDevPaymentAmount(users UserData) float64 {
	sumSquare0, sum0 := 0.0, 0.0
	sumSquare1, sum1 := 0.0, 0.0
	sumSquare2, sum2 := 0.0, 0.0
	sumSquare3, sum3 := 0.0, 0.0
	for i := 0 ; i < len(users.amounts) / 4 * 4 ; i += 4 {
		x := float64(users.amounts[i]) * 0.01
		sumSquare0 += x * x
		sum0 += x

		x = float64(users.amounts[i + 1]) * 0.01
		sumSquare1 += x * x
		sum1 += x

		x = float64(users.amounts[i + 2]) * 0.01
		sumSquare2 += x * x
		sum2 += x

		x = float64(users.amounts[i + 3]) * 0.01
		sumSquare3 += x * x
		sum3 += x
	}
	sumSquare := sumSquare0 + sumSquare1 + sumSquare2 + sumSquare3
	sum := sum0 + sum1 + sum2 + sum3

	count := float64(len(users.amounts))
	avgSquare := sumSquare / count
	avg := sum / count
	return math.Sqrt(avgSquare - avg * avg)
}

func LoadData() UserData {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	ages := make([]uint8, len(userLines))
	for i, line := range userLines {
		age, _ := strconv.Atoi(line[2])
		ages[i] = uint8(age)
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	amounts := make([]uint32, len(paymentLines))
	for i, line := range paymentLines {
		paymentCents, _ := strconv.ParseUint(line[0], 10, 32)
		amounts[i] = uint32(paymentCents)
	}

	return UserData {
		ages,
		amounts,
	}
}
