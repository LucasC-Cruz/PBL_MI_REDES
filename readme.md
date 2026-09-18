# PBL - MI: Redes de Computadores (Sistema de Caronas)

Este repositório contém o código-fonte de um sistema distribuído de gerenciamento de caronas desenvolvido em Go (Golang) como parte de um Problema Baseado em Laboratório (PBL) da disciplina de Redes de Computadores. O sistema adota uma arquitetura cliente-servidor para gerenciar motoristas, passageiros, viagens e rotas utilizando dados persistidos em arquivos JSON.

## Estrutura Geral do Projeto

A organização dos diretórios do projeto está estruturada da seguinte forma:

```text
.
├── .vscode/
│   └── launch.json
├── apiclient/
│   └── client.go
├── Client/
│   └── main.go
├── Docker/
│   ├── Client/
│   │   └── Dockerfile
│   └── Server/
│       └── Dockerfile
├── Lgcc/
│   └── lgcc.go
├── Model/
│   └── model.go
├── Server/
│   ├── carona.json
│   ├── caronas_passageiros.json
│   ├── grafo.json
│   ├── Motoristas.json
│   ├── Passageiro.json
│   ├── servidor.go
│   ├── telas.json
│   └── Viagens.json
├── test/
│   └── servidor_concorrencia_test.go
├── Makefile
├── go.mod
└── README.md
```

## Manual de Uso (Execução Local)

Se você deseja rodar a aplicação diretamente na sua máquina sem utilizar containers, siga os passos abaixo:

### Pré-requisitos

* Go (Golang) instalado (versão 1.18 ou superior recomendada).

### 1. Executando o Servidor

Abra o terminal, navegue até a raiz do projeto e execute o servidor passando a porta desejada (ex: `6742`):

```bash
cd Server
go run servidor.go 6742
```

### 2. Executando o Cliente

Abra outro terminal, navegue até a pasta do cliente e execute o programa:

```bash
Client
go run .
```

## Manual de Uso com Docker e Makefile (Multi-computador)

Para garantir que o ambiente seja idêntico em qualquer máquina e facilitar o deploy/execução em computadores diferentes, utilizamos o Docker em conjunto com o Makefile.

### Pré-requisitos

* Docker instalado e em execução.

### Comandos do Makefile

Na raiz do projeto, você pode utilizar os seguintes comandos no terminal:

1. **Construir as imagens Docker:**
   ```bash
   make build
   ```
   *(Cria as imagens `server-app` e `client-app` utilizando os Dockerfiles localizados na pasta `Docker/`).*

2. **Executar o Servidor:**
   ```bash
   make run-server
   ```
   *(Cria a rede interna do Docker `rede-projeto` e inicia o container do servidor escutando na porta `6742`).*

3. **Executar o Cliente:**
   ```bash
   make run-client
   ```
   *(Inicia o container do cliente de forma interativa, permitindo a comunicação com o servidor).*

4. **Limpar o Ambiente:**
   ```bash
   make clean
   ```
   *(Para e remove os containers ativos, além de apagar a rede Docker criada).*

### Executando em Computadores Distintos (Rede Local / Diferentes Máquinas)

Caso o Servidor esteja rodando em um computador físico da rede (por exemplo, com IP `192.168.1.50`) e você deseje rodar o Cliente a partir de outro computador:

1. Certifique-se de que a porta `6742` está liberada no firewall do computador onde o servidor está hospedado.
2. No código de conexão do cliente (dentro da pasta `apiclient/`), certifique-se de que o endereço de conexão aponta para o IP real da máquina servidora em vez de `localhost`.