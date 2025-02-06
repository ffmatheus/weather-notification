# API DE NOTIFICAÇÃO SOBRE O CLIMA

API desenvolvida com o intuito de oferecer um serviço de agendamento de notificações sobre o clima e o tempo.

## Arquitetura

O projeto segue uma arquitetura hexagonal (também conhecida pelo conceito de portas e adaptadores) para manter o domínio da aplicação isolado e independente de tecnologias externas. E levando também em consideração a pretenção de inclusão de diversos canais de notificação, esse tipo de arquitetura facilitaria a implementação.

![Arquitetura do Projeto](weather-notification-architeture.png)

### Estrutura da Arquitetura

- **Adaptadores Primários**: Serão os "startes" da aplicação
- **Core**: Possui as lógicas/regras de negócio
  - Portas: Interfaces que definem os contratos (implementações a serem realizadas)
  - Domínio: Serviços e entidades que implementam as regras de negócio
- **Adaptadores Secundários**: Implementações para serviços externos como a api do CPTEC e o RabbitMQ
  - PostgreSQL: Persistência de dados
  - RabbitMQ: Gerenciamento de filas e mensageria
  - CPTEC API: Serviço externo de previsão do tempo

## Principais Tecnologias

- Go 1.23
- PostgreSQL
- RabbitMQ
- Docker
- Docker Compose

## Configuração

### Pré-requisitos
- [Go 1.23+](https://go.dev/dl/)
- [Docker e Docker Compose](https://www.docker.com/)

### Instalação
1. Clone o repositório
2. Copie `.env.example` para `.env` e configure as variáveis
3. Execute: `docker-compose up -d`

## API

### Autenticação
Todas as rotas requerem o header `Authorization` com um token válido. Que pode ser obtido ou configurado através da env API_TOKEN

### Endpoints

#### Usuários
- `POST /api/users` - Criar usuário
- `PUT /api/users/{id}` - Atualizar usuário
- `PATCH /api/users/{id}/optout` - Atualizar opt-out
- `GET /api/users` - Listar usuários

#### Notificações
- `POST /api/notifications` - Agendar notificação
- `GET /api/notifications` - Listar notificações do usuário

#### Clima
- `GET /api/weather/search` - Buscar cidade
- `GET /api/weather/forecast` - Buscar previsão

#### Notificações Globais
- `POST /api/notifications/global` - Criar notificação global
- `GET /api/notifications/global` - Listar notificações globais

## Utilização

- Criação de usuário fornecendo o nome da cidade, nome do usuário e e-mail
- Ao buscar uma cidade, caso ela ainda não tenha sido armazenada na base de dados, é feita a persistência do dado
- Para buscar o uuid de uma cidade, basta usar o endpoint de busca/listagem
- As notificações globais notificam TODOS os usuários com opt-out FALSE, com as informações de suas respectivas cidades vinculadas no cadastro
- Nas notificações customizáveis, o usuário consegue criar horários específicos e adicionar notificações de outras cidades

### Documentação
Acesse a documentação completa da API em `/swagger/index.html`

Para autenticar as requisições, no Authorize adicionar "Bearer" + API_TOKEN (variável de ambiente)

## Acessando Interfaces

#### RabbitMQ
- URL: http://localhost:15672
- Usuario: teste (conforme configuração)
- Senha: teste (conforme configuração)

A interface permite monitorar as filas e mensagens

## Testes

A aplicação possui alguns exemplos de testes automatizados que podem ser referência para a criação de novos em futuras implementações.

### Estrutura
Os testes estão organizados junto aos arquivos que testam, seguindo o sufixo `_test.go`

### Executando os Testes
```bash

# Executa todos os testes
go test ./...

# Executa testes de um pacote específico
go test ./internal/domain/entity

# Executa testes com cobertura
go test -cover ./...
```