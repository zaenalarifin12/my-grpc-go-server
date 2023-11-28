package main

import (
	"database/sql"
	uuid2 "github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/zaenalarifin12/my-grpc-go-server/db"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/adapter/database"
	mygrpc "github.com/zaenalarifin12/my-grpc-go-server/internal/adapter/grpc"
	app "github.com/zaenalarifin12/my-grpc-go-server/internal/application"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/domain/bank"
	"log"
	"math/rand"
	"time"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	sqlDB, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/bank?sslmode=disable")
	if err != nil {
		log.Fatal("Can't connect to the database: ", err)
	}
	defer sqlDB.Close()

	db.Migrate(sqlDB)

	databaseAdapter, err := database.NewDatabaseAdapter(sqlDB)

	if err != nil {
		log.Fatalln("Can't create database adapter : ", err)
	}

	//runDummyOrm(databaseAdapter)

	hs := &app.HelloService{}
	bs := app.NewBankService(databaseAdapter)
	rs := &app.ResiliencyService{}

	go generateExchangeRates(bs, "USD", "IDR", 5*time.Second)

	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, rs, 9000)
	grpcAdapter.Run()
}

func runDummyOrm(da *database.DatabaseAdapter) {
	now := time.Now()

	uuid, _ := da.Save(&database.DummyOrm{
		UserId:    uuid2.New(),
		UserName:  "Tim " + time.Now().Format("15:04:00"),
		CreatedAt: now,
		UpdatedAt: now,
	})

	res, _ := da.GetByUUID(&uuid)

	log.Println("res : ", res)
}

func generateExchangeRates(bs *app.BankService, fromCurrency, toCurrency string, duration time.Duration) {
	ticker := time.NewTicker(duration)

	for range ticker.C {
		now := time.Now()
		validFrom := now.Truncate(time.Second).Add(3 * time.Second)
		validTo := validFrom.Add(duration).Add(-1 * time.Millisecond)

		dummyRate := bank.ExchangeRate{
			FromCurrency:       fromCurrency,
			ToCurrency:         toCurrency,
			Rate:               2000 + float64(rand.Intn(300)),
			ValidFromTimestamp: validFrom,
			ValidToTimestamp:   validTo,
		}

		bs.CreateExchangeRate(dummyRate)
	}

}
