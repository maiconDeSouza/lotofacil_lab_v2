package loteria

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"slices"
	"time"

	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/env"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/errapp"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/model"
)

type LoteriaClient struct {
	env    *env.Env
	client *http.Client
}

func NewLoteriaClient(env *env.Env) *LoteriaClient {
	return &LoteriaClient{
		env: env,
		client: &http.Client{
			Timeout: time.Duration(5) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (l *LoteriaClient) GetLastContest() (*model.Concurso, *errapp.ErrApp) {
	req, err := http.NewRequest(http.MethodGet, l.env.URL, nil)
	if err != nil {
		errapp := errapp.NewErrApp(err, "Erro ao criar requisição")
		return nil, errapp
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	res, err := l.client.Do(req)
	if err != nil {
		errapp := errapp.NewErrApp(err, "Erro ao acessar a API")
		return nil, errapp
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errapp := errapp.NewErrApp(nil, fmt.Sprintf(
			"Erro ao acessar a API retornou status %d, mas era esperado %d", res.StatusCode, http.StatusOK,
		),
		)
		return nil, errapp
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		errapp := errapp.NewErrApp(err, "Erro ao ler resposta da API")
		return nil, errapp
	}

	contest := model.Concurso{CriadoEm: time.Now()}
	if err := json.Unmarshal(body, &contest); err != nil {
		errapp := errapp.NewErrApp(err, "API retornou json inesperado")
		return nil, errapp
	}

	return &contest, nil
}

func (l *LoteriaClient) GetNumberContest(c int) (*model.Concurso, *errapp.ErrApp) {
	url := fmt.Sprintf("%s/%d", l.env.URL, c)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		errapp := errapp.NewErrApp(err, "Erro ao criar requisição")
		return nil, errapp
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	res, err := l.client.Do(req)
	if err != nil {
		errapp := errapp.NewErrApp(err, "Erro ao acessar a API")
		return nil, errapp
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errapp := errapp.NewErrApp(nil, fmt.Sprintf(
			"Erro ao acessar a API retornou status %d, mas era esperado %d", res.StatusCode, http.StatusOK,
		),
		)
		return nil, errapp
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		errapp := errapp.NewErrApp(err, "Erro ao ler resposta da API")
		return nil, errapp
	}

	contest := model.Concurso{CriadoEm: time.Now()}
	if err := json.Unmarshal(body, &contest); err != nil {
		errapp := errapp.NewErrApp(err, "API retornou json inesperado")
		return nil, errapp
	}

	return &contest, nil
}

func (l *LoteriaClient) GetFrequentNumbers() (*model.ResultadoFrequencia, *errapp.ErrApp) {
	concurso, errRes := l.GetLastContest()
	if errRes != nil {
		return nil, errRes
	}

	m := make(map[string]int)
	for x := 1; x <= 25; x++ {
		m[fmt.Sprintf("%02d", x)] = 0
	}

	numero := concurso.Numero
	limite := numero - 29 // Busca os últimos 30 concursos (inclusivo)

	for x := numero; x >= limite; x-- {
		c, errRes := l.GetNumberContest(x)
		if errRes != nil {
			return nil, errRes
		}
		for _, n := range c.ListaDezenas {
			m[n]++
		}
	}

	ranking := l.ordenarFrequencia(m)

	palpite, err := l.generateFrequent(ranking)
	if err != nil {
		// Ajuste a criação do ErrApp conforme o construtor do seu pacote errapp
		return nil, errapp.NewErrApp(err, "Deu ruim")
	}

	return &model.ResultadoFrequencia{
		Palpite: palpite,
		Ranking: ranking,
	}, nil
}

func (l *LoteriaClient) ordenarFrequencia(m map[string]int) []model.DezenaFrequencia {
	ranking := make([]model.DezenaFrequencia, 0, len(m))

	for num, qtd := range m {
		ranking = append(ranking, model.DezenaFrequencia{
			Numero:     num,
			Quantidade: qtd,
		})
	}

	slices.SortFunc(ranking, func(a, b model.DezenaFrequencia) int {
		if a.Quantidade != b.Quantidade {
			return b.Quantidade - a.Quantidade // Maior quantidade primeiro
		}
		if a.Numero < b.Numero {
			return -1
		}
		if a.Numero > b.Numero {
			return 1
		}
		return 0
	})

	return ranking
}

func (l *LoteriaClient) generateFrequent(numbers []model.DezenaFrequencia) (model.ListaDezenas, error) {
	if len(numbers) != 25 {
		return nil, fmt.Errorf("quantidade inválida de dezenas: esperado 25, recebido %d", len(numbers))
	}

	// Define as 5 faixas e quantos números sortear de cada uma (5+4+3+2+1 = 15)
	tiers := []struct {
		items    []model.DezenaFrequencia
		quantity int
	}{
		{numbers[0:5], 5},
		{numbers[5:10], 4},
		{numbers[10:15], 3},
		{numbers[15:20], 2},
		{numbers[20:25], 1},
	}

	var result model.ListaDezenas

	for _, t := range tiers {
		var pool []string
		for _, item := range t.items {
			pool = append(pool, item.Numero)
		}

		drawn, err := l.drawNumbers(pool, t.quantity)
		if err != nil {
			return nil, fmt.Errorf("erro ao sortear faixa: %w", err)
		}
		result = append(result, drawn...)
	}

	if len(result) != 15 {
		return nil, fmt.Errorf("erro ao gerar números: total diferente de 15")
	}

	// Ordena o palpite final em ordem crescente ("01", "02", "05"...)
	slices.Sort(result)

	return result, nil
}

func (l *LoteriaClient) drawNumbers(source []string, quantity int) ([]string, error) {
	if quantity > len(source) {
		return nil, fmt.Errorf("quantidade solicitada (%d) é maior que o número de elementos disponíveis (%d)", quantity, len(source))
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("a quantidade a sortear deve ser maior que zero")
	}

	copiedNumbers := make([]string, len(source))
	copy(copiedNumbers, source)

	rand.Shuffle(len(copiedNumbers), func(i, j int) {
		copiedNumbers[i], copiedNumbers[j] = copiedNumbers[j], copiedNumbers[i]
	})

	return copiedNumbers[:quantity], nil
}
