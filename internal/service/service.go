package service

import (
	"errors"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/errapp"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/loteria"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/model"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/repository"
)

type ServiceInterface interface {
	GetContest() (*model.Concurso, *errapp.ErrAppReq)
	GetNumberContest(str string) (*model.Concurso, *errapp.ErrAppReq)
	GenerateDailyContest() (*model.Concurso, *errapp.ErrAppReq)
	GenerateDailyGames() (*model.RespostaJogos, *errapp.ErrAppReq)
}

type Service struct {
	repo    repository.RepositoryInterface
	loteria loteria.LoteriaClient
}

func NewService(repo repository.RepositoryInterface, loteria loteria.LoteriaClient) *Service {
	return &Service{repo: repo, loteria: loteria}
}

func (s *Service) GetContest() (*model.Concurso, *errapp.ErrAppReq) {
	contest, err := s.loteria.GetLastContest()
	if err != nil {
		errappRes := errapp.NewErrAppReq(err.Err, err.Message, http.StatusInternalServerError)
		return nil, errappRes
	}

	return contest, nil
}

func (s *Service) GetNumberContest(str string) (*model.Concurso, *errapp.ErrAppReq) {
	c, err := strconv.Atoi(str)
	if err != nil {
		errappRes := errapp.NewErrAppReq(err, "Problema na requisição", http.StatusBadRequest)
		return nil, errappRes
	}

	contest, errRes := s.loteria.GetNumberContest(c)
	if errRes != nil {
		errappRes := errapp.NewErrAppReq(errRes.Err, errRes.Message, http.StatusInternalServerError)
		return nil, errappRes
	}

	return contest, nil
}

func (s *Service) GenerateDailyContest() (*model.Concurso, *errapp.ErrAppReq) {
	c2 := &model.Concurso{}
	c1, errRes := s.loteria.GetLastContest()
	if errRes != nil {
		errappRes := errapp.NewErrAppReq(errRes.Err, errRes.Message, http.StatusInternalServerError)
		return nil, errappRes
	}

	c2 = c1

	for {
		c2, errRes := s.loteria.GetNumberContest(c2.NumeroConcursoProximo)
		if errRes != nil {
			break
		}

		err := s.repo.CreateContest(c2)
		if err != nil {
			errappRes := errapp.NewErrAppReq(err, "erro ao salvar novo concurso", http.StatusInternalServerError)
			return nil, errappRes
		}
		if len(c2.ListaDezenas) == 15 {
			c2.StatusResultado = model.ResultadoConfirmado
		}
		c1 = c2
	}

	if len(c1.ListaDezenas) == 15 {
		c1.StatusResultado = model.ResultadoConfirmado
	}
	err := s.repo.CreateContest(c1)
	if err != nil {
		errappRes := errapp.NewErrAppReq(err, "erro ao salvar novo concurso", http.StatusInternalServerError)
		return nil, errappRes
	}

	return c1, nil
}

func (s *Service) GenerateDailyGames() (*model.RespostaJogos, *errapp.ErrAppReq) {
	c, errRes := s.loteria.GetLastContest()
	if errRes != nil {
		errappRes := errapp.NewErrAppReq(errRes.Err, errRes.Message, http.StatusInternalServerError)
		return nil, errappRes
	}

	f := &model.Jogo{}
	if err := s.repo.FindGame(f, string(model.EstrategiaFixa), c.Numero); err != nil {
		if !errors.Is(err, errapp.ErrNoDocuments) {
			errappRes := errapp.NewErrAppReq(err, "Erro de servidor", http.StatusInternalServerError)
			return nil, errappRes
		}
		f.Numero = c.NumeroConcursoProximo
		f.TipoEstrategia = model.EstrategiaFixa
		f.CriadoEm = time.Now()
		f.ListaDezenas = s.generateLotofacilNumbers15()
		err = s.repo.CreateGame(f)
		if err != nil {
			errappRes := errapp.NewErrAppReq(err, "Erro de servidor", http.StatusInternalServerError)
			return nil, errappRes
		}
	}

	fixo := f

	fr, errRes := s.loteria.GetFrequentNumbers()
	if errRes != nil {
		errappRes := errapp.NewErrAppReq(errRes.Err, errRes.Message, http.StatusInternalServerError)
		return nil, errappRes
	}

	nf := &model.Jogo{
		Numero:         c.NumeroConcursoProximo,
		ListaDezenas:   fr.Palpite,
		TipoEstrategia: model.EstrategiaMaiorFrequencia,
		CriadoEm:       time.Now(),
	}

	err := s.repo.CreateGame(nf)
	if err != nil {
		errappRes := errapp.NewErrAppReq(err, "Erro de servidor", http.StatusInternalServerError)
		return nil, errappRes
	}

	al := &model.Jogo{
		Numero:         c.NumeroConcursoProximo,
		ListaDezenas:   s.generateLotofacilNumbers15(),
		TipoEstrategia: model.EstrategiaAleatoria,
		CriadoEm:       time.Now(),
	}

	err = s.repo.CreateGame(al)
	if err != nil {
		errappRes := errapp.NewErrAppReq(err, "Erro de servidor", http.StatusInternalServerError)
		return nil, errappRes
	}

	jogos := &model.RespostaJogos{
		EstrategiaFixa:            fixo,
		EstrategiaMaiorFrequencia: nf,
		EstrategiaAleatoria:       al,
		Frequentes:                fr,
	}

	return jogos, nil
}

func (s *Service) generateLotofacilNumbers15() []string {
	l := rand.Perm(25)

	l = l[:15]

	for i := range l {
		l[i]++
	}

	slices.Sort(l)

	r := []string{}

	for _, v := range l {
		r = append(r, strconv.Itoa(v))
	}

	return r
}
