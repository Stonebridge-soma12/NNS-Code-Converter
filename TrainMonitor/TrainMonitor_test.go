package TrainMonitor

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"
)

func getDBInfo() string {
	id := os.Getenv("id")
	pw := os.Getenv("pw")
	url := os.Getenv("url")

	return fmt.Sprintf("%s:%s@tcp(%s)/nns?parseTime=true&charset=utf8mb4", id, pw, url)
}

func TestGetEpochsFromDB(t *testing.T) {
	db, err := sqlx.Connect("mysql", getDBInfo())
	if err != nil {
		t.Error(err)
	}

	epochs, err := GetEpochsFromDB(db, WithTrainID(1))
	if err != nil {
		t.Error(err)
	}

	fmt.Println(epochs)
}

func TestEpoch_PushToDB(t *testing.T) {
	db, err := sqlx.Connect("mysql", getDBInfo())
	if err != nil {
		t.Error(err)
	}

	epoch := Epoch{}
	err = epoch.PushToDB(db)
	if err != nil {
		t.Error(err)
	}
}

func TestTrain_PushToDB(t *testing.T) {
	db, err := sqlx.Connect("mysql", getDBInfo())
	if err != nil {
		t.Error(err)
	}

	train := Train{}
	err = train.PushToDB(db)
	if err != nil {
		t.Error(err)
	}
}

func TestTrain_UpdateToDB(t *testing.T) {
	db, err := sqlx.Connect("mysql", getDBInfo())
	if err != nil {
		t.Error(err)
	}

	train := Train{Status: true}
	err = train.UpdateToDB(db, 1)
	if err != nil {
		t.Error(err)
	}
}