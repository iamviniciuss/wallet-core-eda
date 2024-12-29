# Wallet Core EDA - API

Este projeto é uma API para gerenciamento de uma carteira digital (wallet) baseada no padrão de arquitetura Event-Driven Architecture (EDA). Ele fornece uma estrutura escalável e flexível para lidar com transações financeiras, eventos e estados da carteira de forma assíncrona e desacoplada.

## Índice

- [Visão Geral](#vis%C3%A3o-geral)
- [Funcionalidades](#funcionalidades)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Como Executar](#como-executar)
- [Contribuição](#contribui%C3%A7%C3%A3o)
- [Licença](#licen%C3%A7a)

---

## Visão Geral

O **Wallet Core EDA - API** utiliza eventos para gerenciar o fluxo de informações e transações em tempo real. Essa abordagem permite que os componentes do sistema sejam independentes, o que melhora a escalabilidade e a manutenibilidade.

## Funcionalidades

- Criação e gerenciamento de contas de usuário.
- Processamento de depósitos, saques e transferências.
- Armazenamento e consulta de eventos históricos.
- Integração com sistemas externos por meio de eventos.
- Controle de saldo e validação de transações.

## Tecnologias Utilizadas

- **Golang**: Linguagem principal para o desenvolvimento da API.
- **Kafka/RabbitMQ**: Para gestão de eventos.
- **PostgreSQL**: Banco de dados relacional para persistência de dados.
- **Docker**: Para containerização da aplicação.
- **Clean Architecture**: Para organização e separação de responsabilidades do código.

## Estrutura do Projeto

```
api/
├── cmd/                # Entrypoints para iniciar o servidor
├── config/             # Configurações da aplicação
├── internal/           # Código de negócios principal
│   ├── domain/         # Entidades e contratos
│   ├── usecase/        # Casos de uso da aplicação
│   └── repository/     # Interação com o banco de dados
├── pkg/                # Pacotes compartilhados
├── test/               # Testes automatizados
└── main.go             # Arquivo principal para execução
```

## Como Executar

### Pré-requisitos

- **Docker** e **Docker Compose** instalados.
- Golang 1.20+ instalado.

### Passos

1. Clone o repositório:
   ```bash
   git clone https://github.com/iamviniciuss/wallet-core-eda.git
   cd wallet-core-eda/api
   ```

2. Configure as variáveis de ambiente:
   - Crie um arquivo `.env` na raiz do projeto com base no arquivo `.env.example` e preencha as informações necessárias.

3. Execute a aplicação com Docker Compose:
   ```bash
   docker-compose up --build
   ```

4. Acesse a API:
   - Por padrão, o servidor estará disponível em: `http://localhost:8080`.

### Testes

Execute os testes automatizados com o seguinte comando:
```bash
go test ./...
```

## Contribuição

Contribuições são bem-vindas! Siga os passos abaixo:

1. Faça um fork do repositório.
2. Crie uma branch para sua feature ou correção: `git checkout -b minha-feature`.
3. Commit suas alterações: `git commit -m 'Adicionando nova feature'`.
4. Envie para o repositório original: `git push origin minha-feature`.
5. Abra um Pull Request.

## Licença

Este projeto está licenciado sob a Licença MIT. Veja o arquivo [LICENSE](../LICENSE) para mais informações.