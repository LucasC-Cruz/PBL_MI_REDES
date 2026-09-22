package main

import (
	"Golang/Lgcc"
	"Golang/Model"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type InfoReserva struct {
	Origem  string `json:"origem"`
	Destino string `json:"destin"`
	Data    string `json:"data"`
}

var (
	telas                map[string]string
	passageiros          map[string]Model.Passageiro
	motoristas           map[string]Model.Motorista
	motoristasOnline     map[string]string
	passageirosOnline    map[string]string
	caronas              map[string]Model.Carona
	grafo                map[string]map[string][]Model.Trecho
	caronasPassageiros   map[string]Model.CaronaPassageiro
	caronasPassageirosMu sync.Mutex

	caronasMu    sync.Mutex
	passageiroMu sync.Mutex
	motoristaMu  sync.Mutex
	grafoMu      sync.Mutex
	once         sync.Once
	once2        sync.Once
	once3        sync.Once
	once4        sync.Once
	once5        sync.Once
	once6        sync.Once
	once7        sync.Once
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Comando incorreto!")
		fmt.Println("Talvez queira dizer: go run <arquivo> <porta>?")
		os.Exit(1)
	}

	porta := fmt.Sprintf(":%s", os.Args[1])

	listener, err := net.Listen("tcp", porta)
	if err != nil {
		fmt.Println("ERRO ao criar socket ouvidor, err: ", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Estou escutando no endereço %s\n", listener.Addr())

	iniciarArqTela()
	iniciarArqPassageiros()
	iniciarArqMotoristas()
	iniciarArqGrafo()
	iniciarArqCaronas()
	iniciarArqCaronasPassageiros()
	motoristasOnline = make(map[string]string)
	passageirosOnline = make(map[string]string)

	mux := InicializaMux(&telas)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection, err", err)
			continue
		}

		go handleConnection(conn, &mux)
	}
}

func handleConnection(conn net.Conn, mux *map[string]func(net.Conn, []byte, string)) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		headers := make(map[string]string)
		var err error

		headers, err = Lgcc.OneForAllParser(reader, headers)
		if err != nil {
			fmt.Println("Deu ruim o parser")
		}

		if headers["KeepAlive"] == "true" {
			switch headers["Protocolo"] {
			case "get":
				(*mux)[headers["Caminho"]](conn, nil, headers["Token"])
			case "post":
				n, _ := strconv.Atoi(headers["Tamanho-Conteudo"])
				body := make([]byte, n)
				if _, err := io.ReadFull(reader, body); err != nil {
					fmt.Println("Deu ruim a leitura no server post")
				}
				(*mux)[headers["Caminho"]](conn, body, headers["Token"])
			default:
				fmt.Println("é o default ")
			}
		} else {
			delete(motoristasOnline, headers["Token"])
			delete(passageirosOnline, headers["Token"])
			fmt.Println("Encerrando conexão graciosamente ;)")
			break
		}
	}
}

func iniciarArqTela() {
	once.Do(func() {
		dados, err := os.ReadFile("../data/telas.json")
		if err != nil {
			fmt.Println("Erro ao ler arquivo", err)
			os.Exit(1)
		}
		err = json.Unmarshal(dados, &telas)
		if err != nil {
			fmt.Println("Erro na conversão de json para map")
			os.Exit(1)
		}
	})
}

func iniciarArqPassageiros() {
	once3.Do(func() {
		dados, err := os.ReadFile("../data/Passageiro.json")
		if err != nil {
			fmt.Println("Erro ao ler arquivo Passageiro.json", err)
			os.Exit(1)
		}
		err = json.Unmarshal(dados, &passageiros)
		if err != nil {
			fmt.Println("Erro na conversão de json para slice de passageiros")
			os.Exit(1)
		}
	})
}

func iniciarArqMotoristas() {
	once4.Do(func() {
		dados, err := os.ReadFile("../data/Motoristas.json")
		if err != nil {
			fmt.Println("Erro ao ler arquivo Motoristas.json", err)
			os.Exit(1)
		}
		err = json.Unmarshal(dados, &motoristas)
		if err != nil {
			fmt.Println("Erro na conversão de json para slice de usuarios")
			os.Exit(1)
		}
	})
}

func iniciarArqGrafo() {
	once5.Do(func() {
		if grafo == nil {
			grafo = make(map[string]map[string][]Model.Trecho)
		}
		dados, err := os.ReadFile("../data/grafo.json")
		if err != nil {
			if os.IsNotExist(err) {
				os.WriteFile("../data/grafo.json", []byte("{}"), 0644)
				return
			}
			os.Exit(1)
		}
		if len(dados) == 0 {
			return
		}
		json.Unmarshal(dados, &grafo)
		if grafo == nil {
			grafo = make(map[string]map[string][]Model.Trecho)
		}
	})
}

func iniciarArqCaronas() {
	once6.Do(func() {
		if caronas == nil {
			caronas = make(map[string]Model.Carona)
		}
		dados, err := os.ReadFile("../data/carona.json")
		if err != nil {
			if os.IsNotExist(err) {
				os.WriteFile("../data/carona.json", []byte("{}"), 0644)
				return
			}
			os.Exit(1)
		}
		if len(dados) == 0 {
			return
		}
		json.Unmarshal(dados, &caronas)
		if caronas == nil {
			caronas = make(map[string]Model.Carona)
		}
	})
}

func iniciarArqCaronasPassageiros() {
	once7.Do(func() {
		if caronasPassageiros == nil {
			caronasPassageiros = make(map[string]Model.CaronaPassageiro)
		}
		dados, err := os.ReadFile("../data/caronas_passageiros.json")
		if err != nil {
			if os.IsNotExist(err) {
				os.WriteFile("../data/caronas_passageiros.json", []byte("{}"), 0644)
				return
			}
			return
		}
		if len(dados) == 0 {
			return
		}
		json.Unmarshal(dados, &caronasPassageiros)
		if caronasPassageiros == nil {
			caronasPassageiros = make(map[string]Model.CaronaPassageiro)
		}
	})
}

func InicializaMux(tela *map[string]string) map[string]func(net.Conn, []byte, string) {
	var mux map[string]func(net.Conn, []byte, string)

	once2.Do(func() {
		telas := *tela

		mux = map[string]func(net.Conn, []byte, string){
			"/": func(conn net.Conn, dados []byte, tk string) {
				valor, existe := telas["/"]
				if existe {
					body, tamanho, tipo, _ := Lgcc.GeraBody(valor)
					header := Lgcc.GeraHeaderFull("OK", 200, tipo, tamanho, "")
					Lgcc.EscritaResp(conn, header, body)
				} else {
					conn.Close()
				}
			},
			"/tela/login": func(conn net.Conn, dados []byte, tk string) {
				valor, existe := telas["/tela/login"]
				if existe {
					body, tamanho, tipo, _ := Lgcc.GeraBody(valor)
					header := Lgcc.GeraHeaderFull("OK", 200, tipo, tamanho, "")
					Lgcc.EscritaResp(conn, header, body)
				} else {
					conn.Close()
				}
			},
			"/tela/cadastroPassageiro": func(conn net.Conn, dados []byte, tk string) {
				valor, existe := telas["/tela/cadastroPassageiro"]
				if existe {
					body, tamanho, tipo, _ := Lgcc.GeraBody(valor)
					header := Lgcc.GeraHeaderFull("OK", 200, tipo, tamanho, "")
					Lgcc.EscritaResp(conn, header, body)
				} else {
					conn.Close()
				}
			},
			"/tela/cadastroMotorista": func(conn net.Conn, dados []byte, tk string) {
				valor, existe := telas["/tela/cadastroMotorista"]
				if existe {
					body, tamanho, tipo, _ := Lgcc.GeraBody(valor)
					header := Lgcc.GeraHeaderFull("OK", 200, tipo, tamanho, "")
					Lgcc.EscritaResp(conn, header, body)
				} else {
					conn.Close()
				}
			},
			"/tela/menuMotorista": func(conn net.Conn, dados []byte, tk string) {
				valor, existe := telas["/tela/menuMotorista"]
				if existe {
					body, tamanho, tipo, _ := Lgcc.GeraBody(valor)
					header := Lgcc.GeraHeaderFull("OK", 200, tipo, tamanho, "")
					Lgcc.EscritaResp(conn, header, body)
				} else {
					conn.Close()
				}
			},
			"/tela/menuMotorista/cidades": func(conn net.Conn, dados []byte, tk string) {
				valor, existe := telas["/tela/menuMotorista/cidades"]
				if existe {
					body, tamanho, tipo, _ := Lgcc.GeraBody(valor)
					header := Lgcc.GeraHeaderFull("OK", 200, tipo, tamanho, "")
					Lgcc.EscritaResp(conn, header, body)
				} else {
					conn.Close()
				}
			},
			"/cadastro/passageiro": func(conn net.Conn, dados2 []byte, tk string) {
				var p Model.Passageiro
				err := json.Unmarshal(dados2, &p)
				if err != nil {
					Lgcc.GeraRespostaErro(999, conn, "")
					return
				}
				passageiroMu.Lock()
				_, existe := passageiros[p.Nome]
				if existe {
					passageiroMu.Unlock()
					Lgcc.GeraRespostaErro(999, conn, "")
					return
				} else {
					p.Id = Model.GerarID()
					passageiros[p.Nome] = p
					cu, _ := json.Marshal(passageiros)
					os.WriteFile("../data/Passageiro.json", cu, 0644)
					passageiroMu.Unlock()

					header := Lgcc.GeraHeaderRespPost("OK", 200, "")
					Lgcc.EscritaResp(conn, header, nil)
				}
			},
			"/cadastro/motorista": func(conn net.Conn, dados2 []byte, tk string) {
				var m Model.Motorista
				err := json.Unmarshal(dados2, &m)
				if err != nil {
					Lgcc.GeraRespostaErro(999, conn, "")
					return
				}
				motoristaMu.Lock()
				_, existe := motoristas[m.Nome]
				if existe {
					motoristaMu.Unlock()
					Lgcc.GeraRespostaErro(999, conn, "")
					return
				} else {
					m.Id = Model.GerarID()
					motoristas[m.Nome] = m
					cu, _ := json.Marshal(motoristas)
					os.WriteFile("../data/Motoristas.json", cu, 0644)
					motoristaMu.Unlock()

					header := Lgcc.GeraHeaderRespPost("OK", 200, "")
					Lgcc.EscritaResp(conn, header, nil)
				}
			},
			"/login/passageiro": func(conn net.Conn, dados2 []byte, tka string) {
				var p Model.Passageiro
				err := json.Unmarshal(dados2, &p)
				if err != nil {
					Lgcc.GeraRespostaErro(999, conn, "")
					return
				}
				_, existe := passageiros[p.Nome]
				if existe {
					tk, _ := Lgcc.GerarToken(16)
					passageirosOnline[tk] = p.Nome
					header := Lgcc.GeraHeaderRespPost("OK", 200, tk)
					Lgcc.EscritaResp(conn, header, nil)
				} else {
					Lgcc.GeraRespostaErro(999, conn, "")
					return
				}
			},
			"/login/motorista": func(conn net.Conn, dados2 []byte, tka string) {
				var m Model.Motorista
				err := json.Unmarshal(dados2, &m)
				if err != nil {
					Lgcc.GeraRespostaErro(999, conn, "")
					return
				}
				_, existe := motoristas[m.Nome]
				if existe {
					tk, _ := Lgcc.GerarToken(16)
					motoristasOnline[tk] = m.Nome
					header := Lgcc.GeraHeaderRespPost("OK", 200, tk)
					Lgcc.EscritaResp(conn, header, nil)
				} else {
					Lgcc.GeraRespostaErro(999, conn, "")
					return
				}
			},
			"/passageiro/listarCorridasCadastradas": func(conn net.Conn, dados2 []byte, tk string) {
				nomePassageiro, online := passageirosOnline[tk]
				if !online {
					Lgcc.GeraRespostaErro(999, conn, tk)
					return
				}

				passageiroMu.Lock()
				pass, existe := passageiros[nomePassageiro]
				passageiroMu.Unlock()

				if !existe || pass.IdCaronaPassageiro == "" {
					emptyBytes, _ := json.Marshal([]any{})
					body, tamanho, tipo, _ := Lgcc.GeraBody(emptyBytes)
					header := Lgcc.GeraHeaderFull("OK", 200, tipo, tamanho, tk)
					Lgcc.EscritaResp(conn, header, body)
					return
				}

				caronasPassageirosMu.Lock()
				caronaPass, existeCp := caronasPassageiros[pass.IdCaronaPassageiro]
				caronasPassageirosMu.Unlock()

				var dadosEnvio []byte
				if !existeCp {
					dadosEnvio, _ = json.Marshal([]any{})
				} else {
					dadosEnvio, _ = json.Marshal(caronaPass.Rota)
				}

				body, tamanho, tipo, err := Lgcc.GeraBody(dadosEnvio)
				if err != nil {
					Lgcc.GeraRespostaErro(999, conn, tk)
					return
				}

				header := Lgcc.GeraHeaderFull("OK", 200, tipo, tamanho, tk)
				Lgcc.EscritaResp(conn, header, body)
			},
			"/passageiro/reservarPassagem": func(conn net.Conn, dados2 []byte, tk string) {
				var info InfoReserva
				err := json.Unmarshal(dados2, &info)
				if err != nil {
					Lgcc.GeraRespostaErro(999, conn, tk)
					return
				}

				nomePassageiro, online := passageirosOnline[tk]
				if !online {
					Lgcc.GeraRespostaErro(999, conn, tk)
					return
				}

				origem := strings.ToUpper(strings.TrimSpace(info.Origem))
				destino := strings.ToUpper(strings.TrimSpace(info.Destino))
				dataDesejada := strings.TrimSpace(info.Data)

				grafoMu.Lock()
				trechosEncontrados := buscarMenorCaminhoGrafo(grafo, origem, destino, dataDesejada)
				grafoMu.Unlock()

				if len(trechosEncontrados) == 0 {
					Lgcc.GeraRespostaErro(404, conn, tk)
					return
				}

				grafoMu.Lock()
				for _, tEncontrado := range trechosEncontrados {
					orig := tEncontrado.Trecho.Origem
					dest := tEncontrado.Trecho.Destino
					sliceTrechos := grafo[orig][dest]
					for i := range sliceTrechos {
						if sliceTrechos[i].IdMotorista == tEncontrado.IdMotorista &&
							strings.TrimSpace(sliceTrechos[i].Data) == dataDesejada &&
							sliceTrechos[i].Trecho.Origem == orig &&
							sliceTrechos[i].Trecho.Destino == dest {
							if sliceTrechos[i].Assentos > 0 {
								sliceTrechos[i].Assentos--
								sliceTrechos[i].IdPassageiros = append(sliceTrechos[i].IdPassageiros, nomePassageiro)
							}
						}
					}
				}
				if gBytes, errGrafo := json.Marshal(grafo); errGrafo == nil {
					os.WriteFile("../data/grafo.json", gBytes, 0644)
				}
				grafoMu.Unlock()

				idCaronaPassageiro := Model.GerarID()

				caronaPass := Model.CaronaPassageiro{
					Rota: trechosEncontrados,
				}

				caronasPassageirosMu.Lock()
				if caronasPassageiros == nil {
					caronasPassageiros = make(map[string]Model.CaronaPassageiro)
				}
				caronasPassageiros[idCaronaPassageiro] = caronaPass

				cpBytes, errCp := json.Marshal(caronasPassageiros)
				if errCp == nil {
					os.WriteFile("../data/caronas_passageiros.json", cpBytes, 0644)
				}
				caronasPassageirosMu.Unlock()

				passageiroMu.Lock()
				pass := passageiros[nomePassageiro]
				pass.IdCaronaPassageiro = idCaronaPassageiro
				passageiros[nomePassageiro] = pass
				if pBytes, errP := json.Marshal(passageiros); errP == nil {
					os.WriteFile("../data/Passageiro.json", pBytes, 0644)
				}
				passageiroMu.Unlock()

				header := Lgcc.GeraHeaderRespPost("OK", 200, idCaronaPassageiro)
				Lgcc.EscritaResp(conn, header, nil)
			},
			"/passageiro/cancelarPassagem": func(conn net.Conn, dados2 []byte, tk string) {
				nomePassageiro, online := passageirosOnline[tk]
				if !online {
					Lgcc.GeraRespostaErro(999, conn, tk)
					return
				}

				passageiroMu.Lock()
				pass, existe := passageiros[nomePassageiro]
				if !existe || pass.IdCaronaPassageiro == "" {
					passageiroMu.Unlock()
					Lgcc.GeraRespostaErro(404, conn, tk)
					return
				}
				idReserva := pass.IdCaronaPassageiro
				pass.IdCaronaPassageiro = ""
				passageiros[nomePassageiro] = pass
				if pBytes, errP := json.Marshal(passageiros); errP == nil {
					os.WriteFile("../data/Passageiro.json", pBytes, 0644)
				}
				passageiroMu.Unlock()

				caronasPassageirosMu.Lock()
				caronaPass, caronaExiste := caronasPassageiros[idReserva]
				if caronaExiste {
					delete(caronasPassageiros, idReserva)
					if cpBytes, errCp := json.Marshal(caronasPassageiros); errCp == nil {
						os.WriteFile("../data/caronas_passageiros.json", cpBytes, 0644)
					}
				}
				caronasPassageirosMu.Unlock()

				if caronaExiste {
					grafoMu.Lock()
					for _, tEncontrado := range caronaPass.Rota {
						orig := tEncontrado.Trecho.Origem
						dest := tEncontrado.Trecho.Destino
						sliceTrechos := grafo[orig][dest]
						for i := range sliceTrechos {
							if sliceTrechos[i].IdMotorista == tEncontrado.IdMotorista &&
								strings.TrimSpace(sliceTrechos[i].Data) == strings.TrimSpace(tEncontrado.Data) &&
								sliceTrechos[i].Trecho.Origem == orig &&
								sliceTrechos[i].Trecho.Destino == dest {
								sliceTrechos[i].Assentos++
								newPassageiros := []string{}
								for _, pName := range sliceTrechos[i].IdPassageiros {
									if pName != nomePassageiro {
										newPassageiros = append(newPassageiros, pName)
									}
								}
								sliceTrechos[i].IdPassageiros = newPassageiros
							}
						}
					}
					if gBytes, errGrafo := json.Marshal(grafo); errGrafo == nil {
						os.WriteFile("../data/grafo.json", gBytes, 0644)
					}
					grafoMu.Unlock()
				}

				header := Lgcc.GeraHeaderRespPost("OK", 200, tk)
				Lgcc.EscritaResp(conn, header, nil)
			},
			"/motorista/cadastrarCorrida": func(conn net.Conn, dados2 []byte, tk string) {
				var ca Model.InfoCarona
				var carona Model.Carona
				err := json.Unmarshal(dados2, &ca)
				if err != nil {
					Lgcc.GeraRespostaErro(999, conn, tk)
					return
				}

				nome := motoristasOnline[tk]

				motoristaMu.Lock()
				motorista := motoristas[nome]
				motoristaMu.Unlock()
				idMotorista := motorista.Nome

				idCarona := Model.GerarID()
				motorista.IdCarona = idCarona

				rotas, err := Model.ProcessarLista(ca.Rota)
				if err != nil {
					fmt.Println("Erro ao converter string rotas em slice")
				}

				carona.Data = ca.Data
				carona.Horario = ca.Horario
				carona.Rota = rotas
				carona.Assentos = ca.Assentos
				carona.ValorTrecho = ca.ValorTrecho
				carona.IdMotorista = idMotorista

				caronasMu.Lock()
				caronas[idCarona] = carona
				cu, _ := json.Marshal(caronas)
				os.WriteFile("../data/carona.json", cu, 0644)
				caronasMu.Unlock()

				motoristaMu.Lock()
				motoristas[idMotorista] = motorista
				cu, _ = json.Marshal(motoristas)
				os.WriteFile("../data/Motoristas.json", cu, 0644)
				motoristaMu.Unlock()

				pares, _ := Model.ConstruirPairTrechos(rotas)
				trechos, _ := Model.ConstruirTrechos(pares, ca.Assentos)

				for i := range trechos {
					trechos[i].Data = ca.Data
					trechos[i].IdMotorista = idMotorista
				}

				grafoMu.Lock()
				Model.AdicionarTrechosEmLote(trechos, grafo)
				gBytes, _ := json.Marshal(grafo)
				os.WriteFile("../data/grafo.json", gBytes, 0644)
				grafoMu.Unlock()

				header := Lgcc.GeraHeaderRespPost("OK", 200, tk)
				Lgcc.EscritaResp(conn, header, nil)
			},
			"/motorista/cancelarCorrida": func(conn net.Conn, dados2 []byte, tk string) {
				err := CancelarCorridaMotorista(tk)
				if err != nil {
					Lgcc.GeraRespostaErro(999, conn, tk)
					return
				}
				header := Lgcc.GeraHeaderRespPost("OK", 200, tk)
				Lgcc.EscritaResp(conn, header, nil)
			},
		}
	})
	return mux
}

func buscarMenorCaminhoGrafo(grafo map[string]map[string][]Model.Trecho, origem, destino, dataDesejada string) []Model.Trecho {
	type queueItem struct {
		city    string
		path    []Model.Trecho
		visited map[string]bool
	}

	queue := []queueItem{
		{
			city:    origem,
			path:    []Model.Trecho{},
			visited: map[string]bool{origem: true},
		},
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.city == destino && len(curr.path) > 0 {
			return curr.path
		}

		destinosMap, existe := grafo[curr.city]
		if !existe {
			continue
		}

		for nextCity, trechos := range destinosMap {
			if curr.visited[nextCity] {
				continue
			}

			for _, t := range trechos {
				if strings.TrimSpace(t.Data) == strings.TrimSpace(dataDesejada) && t.Assentos > 0 {
					newVisited := make(map[string]bool)
					for k, v := range curr.visited {
						newVisited[k] = v
					}
					newVisited[nextCity] = true

					newPath := make([]Model.Trecho, len(curr.path))
					copy(newPath, curr.path)
					newPath = append(newPath, t)

					queue = append(queue, queueItem{
						city:    nextCity,
						path:    newPath,
						visited: newVisited,
					})
				}
			}
		}
	}
	return nil
}

func CancelarCorridaMotorista(tokenMotorista string) error {
	nomeMotorista, online := motoristasOnline[tokenMotorista]
	if !online {
		return fmt.Errorf("motorista não está online ou token inválido")
	}

	motoristaMu.Lock()
	motorista, existe := motoristas[nomeMotorista]
	motoristaMu.Unlock()

	if !existe {
		return fmt.Errorf("motorista não encontrado no banco de dados")
	}

	idCarona := motorista.IdCarona
	if idCarona == "" {
		return fmt.Errorf("o motorista não possui uma corrida ativa para cancelar")
	}

	caronasMu.Lock()
	carona, caronaExiste := caronas[idCarona]
	if !caronaExiste {
		caronasMu.Unlock()
		return fmt.Errorf("carona não encontrada")
	}

	motorista.IdCarona = ""
	motoristaMu.Lock()
	motoristas[nomeMotorista] = motorista
	motoristaMu.Unlock()

	rotas := carona.Rota
	delete(caronas, idCarona)
	caronasMu.Unlock()

	grafoMu.Lock()
	defer grafoMu.Unlock()

	passageiroMu.Lock()
	defer passageiroMu.Unlock()

	removido := false
	for i := 0; i < len(rotas)-1; i++ {
		origem := rotas[i]
		destino := rotas[i+1]

		trechos := grafo[origem][destino]

		for idxTrecho, trecho := range trechos {
			// LOG DE TESTE: Verifique no terminal do servidor se isso imprime o match correto
			fmt.Printf("Comparando Trecho Motorista [%s] com Motorista Atual [%s]\n", trecho.IdMotorista, motorista.Nome)

			if trecho.IdMotorista == motorista.Nome {
				for _, idPassageiro := range trecho.IdPassageiros {
					pass, passExiste := passageiros[idPassageiro]
					if passExiste {
						pass.IdCaronaPassageiro = ""
						passageiros[idPassageiro] = pass
					}
				}

				grafo[origem][destino] = append(trechos[:idxTrecho], trechos[idxTrecho+1:]...)
				removido = true
				break
			}
		}
	}

	if !removido {
		fmt.Println("AVISO: Nenhum trecho correspondente a este motorista foi encontrado no grafo!")
	}

	if cuPass, err := json.Marshal(passageiros); err == nil {
		os.WriteFile("../data/Passageiro.json", cuPass, 0644)
	} else {
		fmt.Println("Erro ao serializar Passageiro.json:", err)
	}

	if cuGrafo, err := json.Marshal(grafo); err == nil {
		errWrite := os.WriteFile("../data/grafo.json", cuGrafo, 0644)
		if errWrite != nil {
			fmt.Println("ERRO CRÍTICO ao salvar grafo", errWrite)
		} else {
			fmt.Println("Sucesso: grafo.json atualizado e salvo em disco!")
		}
	} else {
		fmt.Println("Erro ao serializar grafo:", err)
	}

	return nil
}
