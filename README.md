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