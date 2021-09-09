package TrainMonitor

import "github.com/jmoiron/sqlx"

type Train struct {
	Id      int     `db:"id" json:"id"`
	Status  bool    `db:"status" json:"status"`
	Acc     float64 `db:"acc" json:"acc"`
	Loss    float64 `db:"loss" json:"loss"`
	ValAcc  float64 `db:"val_acc" json:"val_acc"`
	ValLoss float64 `db:"val_loss" json:"val_loss"`
	Epochs  int     `db:"epochs" json:"epochs"`
	Name    string  `db:"name" json:"name"`
}

func pushTrain() Option {
	return optionFunc(func(o *options) {
		o.queryString = "insert into " +
			"Train (status, acc, loss, val_acc, val_loss, epochs, name) " +
			"values(:status, :acc, :loss, :val_acc, :val_loss, :epochs, :name)"
	})
}

func updateTrain() Option {
	return optionFunc(func(o *options) {
		o.queryString = "update Train " +
			"set status=:status, acc=:acc, loss=:loss, val_acc=:val_acc, val_loss=:val_loss, epochs=:epochs, name=:name " +
			"where id = :id"
	})
}

func WithId(id int) Option {
	return optionFunc(func (o *options) {
		o.queryString = "select * from Train where id = ?"
		o.args = append(o.args, id)
	})
}

func (t *Train) PushToDB(db *sqlx.DB) error {
	options := options{}
	pushTrain().apply(&options)

	_, err := db.NamedExec(options.queryString, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Train) UpdateToDB(db *sqlx.DB, id int) error {
	options := options{}
	updateTrain().apply(&options)
	t.Id = id

	_, err := db.NamedExec(options.queryString, t)
	if err != nil {
		return err
	}

	return nil
}

func GetTrainFromDB(db *sqlx.DB, opts ... Option) (Train, error) {
	options := options{}

	for _, o := range opts {
		o.apply(&options)
	}

	var res Train
	err := db.Get(&res, options.queryString, options.args)
	if err != nil {
		return Train{}, err
	}

	return res, err
}