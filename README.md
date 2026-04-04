# Airball · NBA Analytics

Plataforma de analytics da NBA com shot charts interativos, métricas avançadas e dados em tempo real da temporada 2025–26.

![Status](https://img.shields.io/badge/status-em%20desenvolvimento-orange)
![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

---

## Visão Geral

O Airball é uma aplicação fullstack de analytics da NBA que permite visualizar shot charts, analisar métricas avançadas e acompanhar qualquer jogador da liga atual. Os dados são obtidos diretamente da NBA Stats API via browser, e o backend gerencia autenticação, cache e persistência.

### Funcionalidades

- **Shot Chart 2D** — quadra interativa com dots por arremesso, heatmap por zona e filtros por tipo de jogada (3pts, mid-range, floater, pull-up, stepback, catch & shoot, layup, crossover)
- **Drives & Passes** — mapa de passes assistidos (xAST) e mapa de penetrações com caminhos de drive
- **Defesa** — heatmap defensivo, DAI (Defensive Activity Index) e métricas on-ball
- **Métricas Avançadas** — OCS (Offensive Creation Score), Gravity Index, PIR (Playmaking Impact Rating)
- **On/Off Court** — net rating com/sem o jogador, impacto por métrica e lineup synergy
- **Busca de jogadores** — lista completa da NBA em tempo real via NBA Stats API
- **Autenticação** — registro e login com JWT
- **UI responsiva** — desktop, tablet e mobile

---

## Stack

| Camada | Tecnologia |
|--------|-----------|
| Backend | Go 1.22 + Gin |
| Banco de dados | PostgreSQL 16 |
| Cache | Redis 7 |
| Auth | JWT (golang-jwt) |
| Frontend | HTML + CSS + JS vanilla |
| Dados NBA | NBA Stats API (fetch direto no browser) |
| Infra local | Docker |

---

## Estrutura do Projeto

```
airball/
├── cmd/
│   └── api/
│       └── main.go                  # Entrypoint — servidor Gin com graceful shutdown
├── internal/
│   ├── cache/
│   │   └── redis.go                 # Cliente Redis
│   ├── config/
│   │   └── config.go                # Leitura de variáveis de ambiente
│   ├── handlers/
│   │   ├── auth.go                  # POST /auth/register, POST /auth/login
│   │   ├── player.go                # GET /leaders/:category, GET /players/:id/shotchart
│   │   └── search.go                # GET /players/search?q=
│   ├── httpclient/
│   │   └── nba.go                   # Cliente HTTP para NBA Stats API
│   ├── middleware/
│   │   ├── auth.go                  # Validação JWT
│   │   ├── cors.go                  # CORS para o frontend
│   │   └── logger.go                # Request logging
│   ├── models/
│   │   ├── player.go                # Model de jogador
│   │   └── user.go                  # Model de usuário
│   ├── repository/
│   │   ├── db.go                    # Conexão PostgreSQL (pgx)
│   │   └── user_repo.go             # Queries de usuário
│   └── service/
│       ├── auth_service.go          # Lógica de autenticação + hash de senha
│       └── player_service.go        # Lógica de dados de jogadores
├── migrations/
│   ├── 001_create_users.sql         # Tabela de usuários
│   └── 002_create_favorites.sql     # Tabela de favoritos
├── pkg/
│   └── logger/
│       └── logger.go                # Logger estruturado
├── frontend/
│   └── airball.html                 # SPA completa (sem dependências externas)
├── .env                             # Variáveis de ambiente (não commitado)
├── .gitignore
├── go.mod
└── go.sum
```

---

## Setup Local

### Pré-requisitos

- [Go 1.22+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

### 1. Subir banco e cache

```powershell
docker run -d --name airball-postgres `
  -e POSTGRES_DB=airball `
  -e POSTGRES_USER=postgres `
  -e POSTGRES_PASSWORD=postgres `
  -p 5432:5432 postgres:16-alpine

docker run -d --name airball-redis `
  -p 6379:6379 redis:7-alpine
```

Se os containers já existem:

```powershell
docker start airball-postgres airball-redis
```

### 2. Aplicar migrations

```powershell
Get-Content migrations\001_create_users.sql | docker exec -i airball-postgres psql -U postgres -d airball
Get-Content migrations\002_create_favorites.sql | docker exec -i airball-postgres psql -U postgres -d airball
```

### 3. Instalar dependências

```powershell
go mod tidy
```

### 4. Subir o backend

```powershell
$env:DB_PASSWORD="postgres"
$env:DB_USER="postgres"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_NAME="airball"
$env:JWT_SECRET="sua-chave-secreta"
$env:PORT="8080"
$env:ENV="development"

go run ./cmd/api/main.go
```

Backend disponível em `http://localhost:8080`.

### 5. Abrir o frontend

Abra `frontend/airball.html` diretamente no navegador. Nenhum servidor necessário.

---

## Endpoints da API

```
GET  /health                           Status do servidor
POST /api/v1/auth/register             Criar conta
POST /api/v1/auth/login                Login, retorna JWT
GET  /api/v1/players/search?q={nome}   Buscar jogadores
GET  /api/v1/leaders/{categoria}       Líderes por estatística
GET  /api/v1/players/{id}/shotchart    Shot chart de um jogador
```

> Os endpoints de players atualmente retornam dados vazios porque a NBA bloqueia requisições server-side. A busca de jogadores está implementada diretamente no browser para contornar esse bloqueio.

---

## Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `DB_HOST` | Host do PostgreSQL | `localhost` |
| `DB_PORT` | Porta do PostgreSQL | `5432` |
| `DB_USER` | Usuário do banco | `postgres` |
| `DB_PASSWORD` | Senha do banco | — |
| `DB_NAME` | Nome do banco | `airball` |
| `REDIS_ADDR` | Endereço do Redis | `localhost:6379` |
| `JWT_SECRET` | Chave secreta para tokens | — |
| `PORT` | Porta do servidor | `8080` |
| `ENV` | Ambiente (`development` / `production`) | `development` |

---

## Roadmap

### Backend

- [ ] Proxy reverso para NBA Stats API — contornar bloqueio server-side
- [ ] Endpoints de favoritos — salvar e remover times e jogadores por usuário
- [ ] Cache por jogador + temporada no Redis com TTL configurável
- [ ] Rate limiting por usuário autenticado
- [ ] Refresh token e controle de expiração de sessão
- [ ] Testes unitários nos services e handlers
- [ ] `docker-compose.yml` para subir toda a stack com um comando

### Frontend

- [ ] Conectar shot chart com dados reais da NBA Stats API
- [ ] Página de perfil com favoritos persistidos no backend
- [ ] Comparativo lado a lado entre dois jogadores
- [ ] Filtro por temporada — regular season vs playoffs
- [ ] Animação de trajetória dos arremessos no shot chart
- [ ] PWA — instalável no mobile

### Analytics

- [ ] Dados reais de drives, passes e defesa via NBA tracking data
- [ ] Cálculo de OCS, DAI e PIR com dados reais ao invés de mockados
- [ ] Histórico de temporadas por jogador
- [ ] Rankings e líderes por categoria com atualização ao vivo

### Infra

- [ ] Deploy do backend (Railway / Render / Fly.io)
- [ ] CI/CD com GitHub Actions
- [ ] Secrets de ambiente via GitHub Actions
- [ ] Logs estruturados com monitoramento

---

## Licença

MIT — veja [LICENSE](LICENSE) para detalhes.
