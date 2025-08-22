package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/model"
	"github.com/google/uuid"
)

func GenerateOrder() model.Order {
	orderUID := uuid.New().String()
	trackNumber := "WBILM" + fmt.Sprintf("%d", rand.Intn(1000000)) + "TEST"

	numItems := rand.Intn(3) + 1
	items := make([]model.Item, numItems)

	var goodsTotal int
	for i := 0; i < numItems; i++ {
		price := rand.Intn(5000)
		sale := rand.Intn(50)
		totalPrice := price * (100 - sale)
		goodsTotal += totalPrice

		items[i] = model.Item{
			ChrtID:      rand.Intn(10000000),
			TrackNumber: trackNumber,
			Price:       price,
			Rid:         uuid.New().String(),
			Name:        fmt.Sprintf("Random Item %d", i+1),
			Sale:        sale,
			Size:        "0",
			TotalPrice:  totalPrice,
			NmID:        rand.Intn(3000000),
			Brand:       "Some Brand",
			Status:      202,
		}
	}

	deliveryCost := rand.Intn(2000)

	order := model.Order{
		OrderUID:        orderUID,
		TrackNumber:     trackNumber,
		Entry:           "WBIL",
		Locale:          "en",
		CustomerID:      "test_customer",
		DeliveryService: "meest",
		Shardkey:        fmt.Sprintf("%d", rand.Intn(10)),
		SmID:            rand.Intn(100),
		DateCreated:     time.Now().Format(time.RFC3339),
		OofShard:        "1",
		Delivery: model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  orderUID,
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       goodsTotal + deliveryCost,
			PaymentDt:    int(time.Now().Unix()),
			Bank:         "alpha",
			DeliveryCost: deliveryCost,
			GoodsTotal:   goodsTotal,
			CustomFee:    0,
		},
		Items: items,
	}

	return order
}
