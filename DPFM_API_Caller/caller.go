package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-payment-method-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-payment-method-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-payment-method-exconf-rmq-kube/database"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.PaymentMethod {
	paymentMethod := *input.PaymentMethod.PaymentMethod
	notKeyExistence := make([]string, 0, 1)
	KeyExistence := make([]string, 0, 1)

	existData := &dpfm_api_output_formatter.PaymentMethod{
		PaymentMethod: paymentMethod,
		ExistenceConf: false,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !e.confPaymentMethod(paymentMethod) {
			notKeyExistence = append(notKeyExistence, paymentMethod)
			return
		}
		KeyExistence = append(KeyExistence, paymentMethod)
	}()

	wg.Wait()

	if len(KeyExistence) == 0 {
		return existData
	}
	if len(notKeyExistence) > 0 {
		return existData
	}

	existData.ExistenceConf = true
	return existData
}

func (e *ExistenceConf) confPaymentMethod(val string) bool {
	rows, err := e.db.Query(
		`SELECT PaymentMethod 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_payment_method_payment_method_data 
		WHERE PaymentMethod = ?;`, val,
	)
	if err != nil {
		e.l.Error(err)
		return false
	}

	for rows.Next() {
		var paymentMethod string
		err := rows.Scan(&paymentMethod)
		if err != nil {
			e.l.Error(err)
			continue
		}
		if paymentMethod == val {
			return true
		}
	}
	return false
}
