package services

import (
	"reflect"
	"testing"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/EvansTrein/RESTful_exchangerServer/pkg/logs"
)

func TestWallet_CurrencyExchangeLogic(t *testing.T) {
	discardLogger := logs.NewDiscardLogger()
	testWallet := &Wallet{log: discardLogger}

	type args struct {
		data *models.CurrencyExchangeData
	}
	tests := []struct {
		name    string
		w       *Wallet
		args    args
		want    *models.CurrencyExchangeResult
		wantErr bool
	}{
		{
            name: "Valid exchange",
            w:    testWallet,
            args: args{
                data: &models.CurrencyExchangeData{
                    BaseBalance:  1000.0,
                    ToBalance:    500.0,
                    ExchangeRate: 1.2,
                    Amount:       100.0,
                },
            },
            want: &models.CurrencyExchangeResult{
                NewBaseBalance: 900.0,
                NewToBalance:  620.0,
                Received:       120.0,
            },
            wantErr: false,
        },
		{
            name: "Valid large values exchange",
            w:    testWallet,
            args: args{
                data: &models.CurrencyExchangeData{
                    BaseBalance:  1000000.0,
                    ToBalance:    500000.0,
                    ExchangeRate: 1.5,
                    Amount:       50000.0,
                },
            },
            want: &models.CurrencyExchangeResult{
                NewBaseBalance: 950000.0,
                NewToBalance:  575000.0,
                Received:       75000.0,
            },
            wantErr: false,
        },
		{
            name: "Valid fractional values exchange",
            w:    testWallet,
            args: args{
                data: &models.CurrencyExchangeData{
                    BaseBalance:  100.5,
                    ToBalance:    50.25,
                    ExchangeRate: 1.1,
                    Amount:       10.5,
                },
            },
            want: &models.CurrencyExchangeResult{
                NewBaseBalance: 90,
                NewToBalance:  61.8,
                Received:       11.55,
            },
            wantErr: false,
        },
		{
            name: "InValid negative amount",
            w:    testWallet,
            args: args{
                data: &models.CurrencyExchangeData{
                    BaseBalance:  1000.0,
                    ToBalance:    500.0,
                    ExchangeRate: 1.2,
                    Amount:       -100.0,
                },
            },
            want:    nil,
            wantErr: true,
        },
		{
            name: "InValid negative exchange rate",
            w:    testWallet,
            args: args{
                data: &models.CurrencyExchangeData{
                    BaseBalance:  1000.0,
                    ToBalance:    500.0,
                    ExchangeRate: -1.2,
                    Amount:       100.0,
                },
            },
            want:    nil,
            wantErr: true,
        },
		{
            name: "InValid insufficient base balance",
            w:    testWallet,
            args: args{
                data: &models.CurrencyExchangeData{
                    BaseBalance:  50.0,
                    ToBalance:    500.0,
                    ExchangeRate: 1.2,
                    Amount:       100.0,
                },
            },
            want:    nil,
            wantErr: true,
        },
		{
            name: "InValid zero exchange rate",
            w:    testWallet,
            args: args{
                data: &models.CurrencyExchangeData{
                    BaseBalance:  1000.0,
                    ToBalance:    500.0,
                    ExchangeRate: 0.0,
                    Amount:       100.0,
                },
            },
            want:    nil,
            wantErr: true,
        },
		{
            name: "InValid zero amount",
            w:    testWallet,
            args: args{
                data: &models.CurrencyExchangeData{
                    BaseBalance:  1000.0,
                    ToBalance:    500.0,
                    ExchangeRate: 1.2,
                    Amount:       0.0,
                },
            },
            want:    nil,
            wantErr: true,
        },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.CurrencyExchangeLogic(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wallet.CurrencyExchangeLogic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wallet.CurrencyExchangeLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}
