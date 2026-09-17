package model

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Passageiro struct {
	Nome               string `json:"nome"`
	Senha              string `json:"senha"`
	Id                 string `json:"id"`
	IdCaronaPassageiro string
}

type Motorista struct {
	Nome     string `json:"nome"`
	Senha    string `json:"senha"`
	Carro    string `json:"carro"`
	Id       string `json:"id"`
	IdCarona string
}

// struct de corrida para o passageiro
// lista de motoristas, lista de trechos
// isso será construído ao fazermos a bfs no grafo para montar uma rota para o passageiro
type CaronaPassageiro struct {
	Rota []Trecho
}

type PairTrecho struct {
	Origem  string
	Destino string
}

//peso do grafo
type Trecho struct {
	IdMotorista   string
	Data          string //vai ficar aqui por enquanto fds
	IdPassageiros []string
	Trecho        PairTrecho 
	Assentos      int
}

// Este caba aqui será construído por nós
// percorrendo o grafo
type Carona struct {
	IdMotorista string
	IdCarona    string
	Data        string
	Horario     string
	ValorTrecho float32
	Assentos    int
	Rota        []string // ou Trecho
}

type InfoCarona struct {
	Data        string
	Horario     string
	ValorTrecho float32
	Assentos    int
	Rota        string
}

func GerarID() string {
	// 1. Pega o tempo atual em milissegundos
	timestamp := time.Now().UnixMilli()

	// 2. Gera 6 bytes de entropia (resultando em 12 caracteres hexadecimais)
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback de segurança caso a leitura de bytes falhe
		return fmt.Sprintf("%d", timestamp)
	}

	// 3. Concatena tudo em uma string formatada
	return fmt.Sprintf("%d-%s", timestamp, hex.EncodeToString(b))
}

//a porra do graf



type Grafo struct {
	// A estrutura: Origem -> Destino -> Slice de Trechos (N conexões)
	conexoes map[string]map[string][]Trecho
}

// Inicializa o grafo na memória
func NovoGrafo() *Grafo {
	return &Grafo{
		conexoes: make(map[string]map[string][]Trecho),
	}
}

// Adiciona uma aresta direcionada (Origem -> Destino)
func (g *Grafo) AdicionarAresta(origem, destino string, peso Trecho) {
	// Se o nó de origem ainda não existe no mapa, precisamos instanciá-lo
	if g.conexoes[origem] == nil {
		g.conexoes[origem] = make(map[string][]Trecho)
	}

	// Adiciona o novo trecho no slice daquela rota específica
	g.conexoes[origem][destino] = append(g.conexoes[origem][destino], peso)
}

// Retorna todas as conexões (e seus pesos) entre dois vértices em O(1)
func (g *Grafo) ConsultarArestas(origem, destino string) []Trecho {
	// Se a origem existir, ele busca o destino.
	// Se não existir rota, o Go naturalmente retorna um slice vazio (nil), sem quebrar o código.
	return g.conexoes[origem][destino]
}

// ConstruirPairTrechos agora recebe um slice de strings e gera os pares de origem e destino
func ConstruirPairTrechos(rotas []string) ([]PairTrecho, error) {
	// Validação básica: precisa de pelo menos 2 pontos para ter origem e destino
	if len(rotas) < 2 {
		return nil, errors.New("são necessários pelo menos dois locais para formar um trecho")
	}

	// Limpa os espaços e valida se há alguma string vazia no meio do slice
	for i := range rotas {
		rotas[i] = strings.TrimSpace(rotas[i])
		if rotas[i] == "" {
			return nil, errors.New("foi encontrado um local vazio na lista de rotas")
		}
	}

	var pares []PairTrecho

	// Monta os pares iterando até o penúltimo elemento
	for i := 0; i < len(rotas)-1; i++ {
		par := PairTrecho{
			Origem:  rotas[i],
			Destino: rotas[i+1],
		}
		pares = append(pares, par)
	}

	return pares, nil
}

// 2. Função que recebe os pares e o número de assentos para montar as arestas (Trechos)
func ConstruirTrechos(pares []PairTrecho, nAssentos int) ([]Trecho, error) {
	// Validação de segurança
	if len(pares) == 0 {
		return nil, errors.New("o slice de pares está vazio, impossível gerar trechos")
	}

	if nAssentos < 1 {
		return nil, errors.New("o número de assentos deve ser maior que zero")
	}

	var trechos []Trecho

	for _, par := range pares {
		novoTrecho := Trecho{
			IdPassageiros: make([]string, 0), // Inicia pronto para receber appends
			Trecho:        par,
			Assentos:      nAssentos,
		}
		trechos = append(trechos, novoTrecho)
	}

	return trechos, nil
}

// ProcessarLista recebe uma string, valida o separador, e retorna um slice em maiúsculo
func ProcessarLista(entrada string) ([]string, error) {
	// 1. Remove os espaços das pontas da string principal
	entrada = strings.TrimSpace(entrada)

	// 2. Valida se a string está vazia
	if entrada == "" {
		return nil, errors.New("a string de entrada está vazia")
	}

	// 3. Valida se a string usa vírgula como separador
	// (Se não tiver vírgula, consideramos que usou outro separador ou é inválida)
	if !strings.Contains(entrada, ",") {
		return nil, errors.New("formato inválido: a string deve ser separada por vírgulas")
	}

	// 4. Faz a separação da string
	itens := strings.Split(entrada, ",")
	var resultado []string

	// 5. Percorre os itens para limpar espaços vazios e converter para Upper Case
	for _, item := range itens {
		itemLimpo := strings.TrimSpace(item)
		
		// Ignora pedaços vazios caso a string venha com vírgulas duplas (ex: "A,,B")
		if itemLimpo != "" {
			resultado = append(resultado, strings.ToUpper(itemLimpo))
		}
	}

	// Validação extra caso a string fosse apenas vírgulas (ex: ",,,")
	if len(resultado) == 0 {
		return nil, errors.New("nenhum valor válido encontrado na string")
	}

	return resultado, nil
}

// AdicionarTrechosEmLote percorre um slice de Trechos e os insere diretamente 
// em um grafo mapeado (map[string]map[string][]Trecho).
func AdicionarTrechosEmLote(trechos []Trecho, grafo map[string]map[string][]Trecho) {
	for _, t := range trechos {
		origem := t.Trecho.Origem
		destino := t.Trecho.Destino

		// Verifica se o nó de origem existe; se não, inicializa
		if grafo[origem] == nil {
			grafo[origem] = make(map[string][]Trecho)
		}

		// Adiciona o novo trecho no slice daquela rota específica
		grafo[origem][destino] = append(grafo[origem][destino], t)
	}
}