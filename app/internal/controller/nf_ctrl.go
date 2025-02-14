package controller

import (
	"app/pkg/model"
	stability "app/pkg/stability"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var mu sync.Mutex

type NFController struct {
	nfs                map[string][]*model.NF
	customerController *CustomerController
	token              string
	mu                 sync.RWMutex
}

func NewNFController(token string) *NFController {
	return &NFController{
		nfs:                make(map[string][]*model.NF),
		customerController: NewCustomerController(token),
		token:              token,
	}
}

func (c *NFController) GetByName(ctx context.Context, name string) ([]*model.NF, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if nf, exists := c.nfs[name]; exists {
		return nf, nil
	}
	// Need to convert name to cpfcnpj for search
	cpfcnpj := c.customerController.GetCPFCNPJFromPartialName(ctx, name)
	if cpfcnpj == "" {
		return nil, fmt.Errorf("no customer found with name: %s", name)
	}

	nfs, err := c.GetByCPFCNPJ(ctx, cpfcnpj)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFs by CPFCNPJ: %w", err)
	}
	return nfs, nil
}

func (c *NFController) GetByNumero(ctx context.Context, numero string) ([]*model.NF, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if nf, exists := c.nfs[numero]; exists {
		return nf, nil
	}
	url := fmt.Sprintf("%s/notas?numero=%s", OlistERPURL, numero)
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
	}

	// Add the Bearer token to the request header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
	}
	// Parse the response
	var apiResponse model.NFResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no NF found with number: %s", numero)
	}

	nfs := parseItensFromNFResponse(c, apiResponse, numero)
	return nfs, nil
}

func (c *NFController) GetByCPFCNPJ(ctx context.Context, cpfcnpj string) ([]*model.NF, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if nf, exists := c.nfs[cpfcnpj]; exists {
		return nf, nil
	}
	url := fmt.Sprintf("%s/notas?cpfCnpj=%s", OlistERPURL, cpfcnpj)
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
	}

	// Add the Bearer token to the request header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
	}
	// Parse the response
	var apiResponse model.NFResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}
	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no NF found with CPF/CNPJ: %s", cpfcnpj)
	}

	nfs := parseItensFromNFResponse(c, apiResponse, cpfcnpj)
	return nfs, nil
}

var nfSituacaoMap = map[string]string{
	"1":  "Pendente",
	"2":  "Emitida",
	"3":  "Cancelada",
	"4":  "Enviada Aguardando Recibo",
	"5":  "Rejeitada",
	"6":  "Autorizada",
	"7":  "Emitida Danfe",
	"8":  "Registrada",
	"9":  "Enviada Aguardando Protocolo",
	"10": "Denegada",
}

func parseItensFromNFResponse(c *NFController, apiResponse model.NFResponse, key string) []*model.NF {
	fmt.Println("Reached parseItensFromNFResponse")
	var nfs []*model.NF
	cb := stability.NewBreaker()

	var wg sync.WaitGroup
	jobs := make(chan int, 4)
	res := make(chan results)
	numWorkers := 4

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				result, err := NFLink(c, job, cb)

				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				res <- result
			}
		}(i)
	}
	/*go func() {
		wg.Wait()
		close(res)
	}()*/

	go func() {
		for r := range res {
			link, err := r.Link, r.err
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				wg.Done()
				continue
			}
			for _, nf := range nfs {
				if nf.ID == r.id {
					nf.Link = link
					fmt.Printf("NFs: %v\n", nf)
				}
			}
			wg.Done()
		}
	}()

	for _, item := range apiResponse.Itens {
		situacaoDesc, ok := nfSituacaoMap[item.Situacao]
		if !ok {
			situacaoDesc = "Unknown"
		}

		nf := &model.NF{
			ID:          item.ID,
			Situacao:    situacaoDesc,
			Numero:      item.Numero,
			Serie:       item.Serie,
			ChaveAcesso: item.ChaveAcesso,
			DataEmissao: item.DataEmissao,
			Cliente: model.CustomerOrderEnpoint{
				CompanyOrderEndpoint: model.CompanyOrderEndpoint{
					TipoCadastro: "CLIENTE",
					ID:           item.Cliente.ID,
					Codigo:       model.Code(item.Cliente.Codigo),
					TipoPessoa:   item.Cliente.TipoPessoa,
					RazaoSocial:  item.Cliente.RazaoSocial,
					NomeFantasia: item.Cliente.NomeFantasia,
					CNPJCPF:      model.CNPJCPF(item.Cliente.CNPJCPF),
					Telefone:     model.Telefone(item.Cliente.Telefone),
					Email:        model.Email(item.Cliente.Email),
					Endereco: model.Endereco{
						Rua:         item.Cliente.Endereco.Rua,
						Numero:      item.Cliente.Endereco.Numero,
						Complemento: item.Cliente.Endereco.Complemento,
						Bairro:      item.Cliente.Endereco.Bairro,
						Municipio:   item.Cliente.Endereco.Municipio,
						Cep:         item.Cliente.Endereco.Cep,
						Uf:          item.Cliente.Endereco.Uf,
						Pais:        item.Cliente.Endereco.Pais,
					},
				},
			},
			Valor:              item.Valor,
			ValorProdutos:      item.ValorProdutos,
			CodigoRastreamento: item.CodigoRastreamento,
			UrlRastreamento:    item.UrlRastreamento,
			FretePorConta:      item.FretePorConta,
			Link:               "",
		}

		wg.Add(1)
		jobs <- item.ID

		nfs = append(nfs, nf)
	}

	close(jobs)

	wg.Wait()
	close(res)
	c.nfs[key] = nfs

	return nfs
}

func getNFLink(c *NFController, id int) (string, error) {
	url := fmt.Sprintf("%s/notas/%d/link", OlistERPURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return "", fmt.Errorf("failed to create request to Olist ERP API 1: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	failures := 1
	for sleep := 1 * time.Second; resp.StatusCode == http.StatusTooManyRequests; failures++ {
		waitingTime := calculateBackOff(failures) + sleep
		fmt.Printf("Rate limited, retrying in %s\n", waitingTime)
		time.Sleep(waitingTime)
		resp, err = client.Do(req)
		if failures > 5 {
			return "", fmt.Errorf("failed to create request to Olist ERP API 2: %w, id: %d, resp: %d", err, id, resp.StatusCode)
		}
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create request to Olist ERP API 2: %w, id: %d, resp: %d", err, id, resp.StatusCode)
	}
	defer resp.Body.Close()

	// Parse the response
	var secondaryResponse model.NFLink
	if err := json.NewDecoder(resp.Body).Decode(&secondaryResponse); err != nil {
		return "", fmt.Errorf("failed to create request to Olist ERP API 3: %w", err)
	}
	return secondaryResponse.Link, nil
}

type results struct {
	id   int
	Link string
	err  error
}

func NFLink(c *NFController, id int, cb stability.Circuitbreaker) (results, error) {
	time.Sleep(calculateBackOff(1))
	link, err := cb.Execute(func() (interface{}, error) {
		return getNFLink(c, id)
	})
	linkStr, ok := link.(string)
	if !ok {
		result := results{
			id:   id,
			Link: "",
			err:  fmt.Errorf("failed to convert link to string"),
		}
		return result, err

	} else {
		result := results{
			id:   id,
			Link: linkStr,
			err:  err,
		}
		return result, err
	}
}

/*
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Tiny ERP API response body: %w", err)
	}
	fmt.Printf("Response: %s\n", string(bodyBytes))
*/
