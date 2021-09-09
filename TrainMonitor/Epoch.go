package TrainMonitor

import (
	"encoding/gob"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"net"
	"net/http"
)

type Epoch struct {
	TrainId      int     `db:"train_id" json:"train_id"`
	Acc          float64 `db:"acc" json:"acc"`
	Epoch        int     `db:"epoch" json:"epoch"`
	Loss         float64 `db:"loss" json:"loss"`
	ValAcc       float64 `db:"val_acc" json:"val_acc"`
	ValLoss      float64 `db:"val_loss" json:"val_loss"`
	LearningRate float64 `db:"learning_rate" json:"lr"`
}

func pushEpoch() Option {
	return optionFunc(func(o *options) {
		o.queryString = "insert into " +
			"Epoch(train_id, epoch, acc, loss, val_acc, val_loss, learning_rate) " +
			"values (:train_id, :epoch, :acc, :loss, :val_acc, :val_loss, :learning_rate)"
	})
}

func WithTrainID(trainId int) Option {
	return optionFunc(func(o *options) {
		o.queryString = "select * from Epoch where train_id = ?"
		o.args = append(o.args, trainId)
	})
}

func (e *Epoch) BindEpoch(r *http.Request) error {
	var epoch []byte
	_, err := r.Body.Read(epoch)
	if err != nil {
		return err
	}

	var res Epoch
	err = json.Unmarshal(epoch, &res)
	if err != nil {
		return err
	}

	return nil
}

func (e *Epoch) PushToDB(db *sqlx.DB) error {
	options := options{}
	pushEpoch().apply(&options)

	_, err := db.NamedExec(options.queryString, e)
	if err != nil {
		return err
	}

	return nil
}

func GetEpochsFromDB(db *sqlx.DB, opts ...Option) ([]Epoch, error) {
	options := options{}

	for _, o := range opts {
		o.apply(&options)
	}

	var epochs []Epoch
	rows, err := db.Queryx(options.queryString, options.args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var epoch Epoch
		err := rows.StructScan(&epoch)
		if err != nil {
			return nil, err
		}
		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func (e *Epoch) PushToSocket(conn net.Conn) error {
	encoder := gob.NewEncoder(conn)

	err := encoder.Encode(e)
	if err != nil {
		return err
	}

	return nil
}
