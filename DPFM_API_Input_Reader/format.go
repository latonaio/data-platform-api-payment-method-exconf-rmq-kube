package dpfm_api_input_reader

import (
	"data-platform-api-payment-method-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToPaymentMethod() *requests.PaymentMethod {
	data := sdc.PaymentMethod
	return &requests.PaymentMethod{
		PaymentMethod: data.PaymentMethod,
	}
}
