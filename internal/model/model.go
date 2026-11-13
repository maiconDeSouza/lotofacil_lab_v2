package model

import (
	"time"
)

type ListaDezenas []string

type TipoEstrategia string

const (
	EstrategiaFixa            TipoEstrategia = "fixa"
	EstrategiaMaiorFrequencia TipoEstrategia = "maiorFrequencia"
	EstrategiaAleatoria       TipoEstrategia = "aleatoria"
)

type StatusResultado string

const (
	ResultadoConfirmado StatusResultado = "confirmado"
	ResultadoAguardando StatusResultado = "aguardando"
)

type DezenaFrequencia struct {
	Numero     string `json:"numero"`
	Quantidade int    `json:"quantidade"`
}

// ResultadoFrequencia agrupa o palpite final e o ranking estatístico
type ResultadoFrequencia struct {
	Palpite ListaDezenas       `json:"palpite"`
	Ranking []DezenaFrequencia `json:"ranking"`
}

type Premio struct {
	DescricaoFixa      string  `bson:"descricaoFaixa" json:"descricaoFaixa"`
	Faixa              int     `bson:"faixa" json:"faixa"`
	NumeroDeGanhadores int     `bson:"numeroDeGanhadores" json:"numeroDeGanhadores"`
	ValorPremio        float64 `bson:"valorPremio" json:"valorPremio"`
}

type Concurso struct {
	Numero                int             `bson:"numero" json:"numero"`
	NumeroConcursoProximo int             `bson:"numeroConcursoProximo" json:"numeroConcursoProximo"`
	TipoJogo              string          `bson:"tipoJogo" json:"tipoJogo"`
	DataApuracao          string          `bson:"dataApuracao" json:"dataApuracao"`
	ListaDezenas          ListaDezenas    `bson:"listaDezenas" json:"listaDezenas"`
	ListaRateioPremio     []Premio        `bson:"listaRateioPremio" json:"listaRateioPremio"`
	Acumulado             bool            `bson:"acumulado" json:"acumulado"`
	StatusResultado       StatusResultado `bson:"statusResultado"`
	CriadoEm              time.Time       `bson:"criadoEm"`
}

type Jogo struct {
	Numero           int            `bson:"numero" json:"numero"`
	ListaDezenas     ListaDezenas   `bson:"listaDezenas" json:"listaDezenas"`
	TipoEstrategia   TipoEstrategia `bson:"tipoEstrategia" json:"tipoEstrategia"`
	Acertos          int            `bson:"acertos" json:"acertos"`
	NumerosAcertados ListaDezenas   `bson:"NumerosAcertados" json:"NumerosAcertados"`
	CriadoEm         time.Time      `bson:"criadoEm" json:"criadoEm"`
}

type Pontuacao struct {
	EstrategiaFixa            int
	EstrategiaMaiorFrequencia int
	EstrategiaAleatoria       int
}

type RespostaJogos struct {
	EstrategiaFixa            *Jogo                `json:"estrategiaFixa "`
	EstrategiaMaiorFrequencia *Jogo                `json:"estrategiaMaiorFrequencia"`
	EstrategiaAleatoria       *Jogo                `json:"estrategiaAleatoria"`
	Frequentes                *ResultadoFrequencia `json:"frequentes"`
}
