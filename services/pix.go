package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Neiva07/queue-management-service/persistence"
	"github.com/gerencianet/gn-api-sdk-go-pix/gerencianet"
)

type location struct {
	ID       uint64 `json:"id"`
	Location string `json:"location"`
}

type Calendar struct {
	CreatedAt time.Time `json:"criacao"`
	ExpiresIn uint64    `json:"expiracao"` //seconds
}

type Value struct {
	Original string `json:"original"`
}

type createImmediateChargeResponse struct {
	Loc        location `json:"loc"`
	Calendar   Calendar `json:"calendario"`
	Status     string   `json:"status"`
	TxID       string   `json:"txid"`
	PIXKey     string   `json:"chave"`
	TotalValue Value    `json:"value"`
}

type QRCode struct {
	Key   string `json:"qrcode"`
	Image string `json:"imagemQrcode"`
}

type GeneratePixResponse struct {
	QRCode     *QRCode
	Status     string
	PIXKey     string
	TotalValue string
	TXID       string
}

type ChargeDetailsResponse struct {
	Status string `json:"status"`
	PIXKey string `json:"chave"`
	TXID   string `json:"txid"`
}

var clientId, clientSecret, pixKey string
var credentials map[string]interface{}

func loadCredentials() {
	var isValid bool

	clientId, isValid = os.LookupEnv("CLIENT_ID")

	clientSecret, isValid = os.LookupEnv("CLIENT_SECRET")

	pixKey, isValid = os.LookupEnv("PIX_KEY")

	if !isValid {
		log.Fatalf("Invalid Environemnt varible to Generate PIX payment")

	}

	credentials = map[string]interface{}{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"sandbox":       false,
		"timeout":       20,
		"CA":            "prod.crt.pem",
		"Key":           "prod.key.pem",
	}
}

func GeneratePixKey(user *persistence.User, value uint64) (*GeneratePixResponse, error) {

	loadCredentials()

	gn := gerencianet.NewGerencianet(credentials)

	body := map[string]interface{}{

		"calendario": map[string]interface{}{
			"expiracao": 3600,
		},
		"devedor": map[string]interface{}{

			"cpf":  user.CPF,
			"nome": user.Name,
		},
		"valor": map[string]interface{}{

			"original": convertToDecimalString(value),
		},
		"chave":              pixKey,
		"solicitacaoPagador": "Teste.",
	}

	res, err := gn.CreateImmediateCharge(body)

	log.Println(res)

	if err != nil {
		log.Println("Error creating immediate charge ", err)
		return nil, err
	}

	cr := createImmediateChargeResponse{}

	err = json.Unmarshal([]byte(res), &cr)

	if err != nil {
		log.Println("Error unparsing immediate charge response ", err)
		return nil, err
	}

	qrCodeResponse, err := gn.GenerateQRCode(fmt.Sprint(cr.Loc.ID))

	if err != nil {
		log.Printf("Error generating qr code for location %v and user %v", cr.Loc.ID, user.UserID)
		log.Println(err)
		return nil, err
	}

	qr := QRCode{}

	err = json.Unmarshal([]byte(qrCodeResponse), &qr)

	if err != nil {
		log.Println("Error  unparsing generate qr code response", err)
		return nil, err
	}

	log.Println("Sucessfully generated pix charge")

	return &GeneratePixResponse{
		QRCode:     &qr,
		Status:     cr.Status,
		PIXKey:     cr.PIXKey,
		TotalValue: cr.TotalValue.Original,
		TXID:       cr.TxID,
	}, nil
}

func CheckCharge(txID string) *ChargeDetailsResponse {
	loadCredentials()

	gn := gerencianet.NewGerencianet(credentials)

	result, err := gn.DetailCharge(txID)

	if err != nil {
		log.Println(err)
	}

	chargeDetails := ChargeDetailsResponse{}

	err = json.Unmarshal([]byte(result), &chargeDetails)

	log.Println(result, err)

	return &chargeDetails
}

func AddPIXWebhook(url string) (string, error) {
	loadCredentials()

	gn := gerencianet.NewGerencianet(credentials)

	body := map[string]interface{}{

		"webhookUrl": url,
		// "webhookUrl": "https://events.hookdeck.com/e/src_rQnCeEscnQJ5",
	}

	s, err := gn.UpdateWebhook(pixKey, body)

	if err != nil {
		log.Println("Error updating web hook", err)
		return "", err
	}

	log.Println("Webhook updated successfully: ", s)

	return s, nil
}

func convertToDecimalString(value uint64) string {
	reais := fmt.Sprint(value / 100)

	v := fmt.Sprint(value / 10)
	i := fmt.Sprint(value)

	decimalCents := string(v[len(v)-1])

	return fmt.Sprintf("%v.%v%v", reais, decimalCents, string(i[len(i)-1]))

}
