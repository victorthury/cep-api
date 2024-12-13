package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/victorthury/cep-api/internal/dto"
)

type CepHandler struct {
	BrasilApiUrl string
	ViaCepUrl    string
}

func NewCepHandler(brasilApiUrl, viaCepUrl string) *CepHandler {
	return &CepHandler{
		BrasilApiUrl: brasilApiUrl,
		ViaCepUrl:    viaCepUrl,
	}
}

func (c *CepHandler) MakeRequestToApi(w http.ResponseWriter, r *http.Request, requestURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	return resBody, nil
}

func (c *CepHandler) GetCepFromBrasilApi(w http.ResponseWriter, r *http.Request, cep string, ch chan dto.GetCepOutput) (*dto.GetBrasilApiOutput, error) {
	requestURL := fmt.Sprintf("%s/api/cep/v1/%s", c.BrasilApiUrl, cep)

	resBody, err := c.MakeRequestToApi(w, r, requestURL)
	if err != nil {
		return nil, err
	}

	var cepBrasilApi dto.GetBrasilApiOutput
	err = json.Unmarshal(resBody, &cepBrasilApi)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	cepOutput := dto.GetCepOutput{Cep: cepBrasilApi.Cep, State: cepBrasilApi.State, City: cepBrasilApi.City, Neighborhood: cepBrasilApi.Neighborhood, Street: cepBrasilApi.Street, Source: "BrasilApi"}
	ch <- cepOutput

	return &cepBrasilApi, nil
}

func (c *CepHandler) GetCepFromViaCep(w http.ResponseWriter, r *http.Request, cep string, ch chan dto.GetCepOutput) (*dto.GetViaCepOutput, error) {
	requestURL := fmt.Sprintf("%s/ws/%s/json/", c.ViaCepUrl, cep)

	resBody, err := c.MakeRequestToApi(w, r, requestURL)
	if err != nil {
		return nil, err
	}

	var cepViaCep dto.GetViaCepOutput
	err = json.Unmarshal(resBody, &cepViaCep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	cepOutput := dto.GetCepOutput{Cep: cepViaCep.Cep, State: cepViaCep.Uf, City: cepViaCep.Localidade, Neighborhood: cepViaCep.Bairro, Street: cepViaCep.Logradouro, Source: "ViaCep"}
	ch <- cepOutput

	return &cepViaCep, nil
}

func (c *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if len(cep) != 8 {
		w.WriteHeader(http.StatusBadRequest)
	}
	brasilApiChannel := make(chan dto.GetCepOutput)
	viaCepChannel := make(chan dto.GetCepOutput)

	go c.GetCepFromBrasilApi(w, r, cep, brasilApiChannel)
	go c.GetCepFromViaCep(w, r, cep, viaCepChannel)

	w.Header().Set("Content-Type", "application/json")

	select {
	case msg := <-brasilApiChannel:
		json.NewEncoder(w).Encode(msg)
		w.WriteHeader(http.StatusOK)
	case msg := <-viaCepChannel:
		json.NewEncoder(w).Encode(msg)
		w.WriteHeader(http.StatusOK)
	case <-time.After(time.Second):
		json.NewEncoder(w).Encode("timeout")
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
