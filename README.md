# Weather Forecast API 🌤️

Uma API em Go que recebe um CEP, identifica a cidade correspondente e retorna as temperaturas atuais em Celsius, Fahrenheit e Kelvin.

## 📋 Funcionalidades

- Validação de CEP (8 dígitos)
- Busca de cidade através do CEP usando a API ViaCEP
- Consulta de temperatura atual usando WeatherAPI
- Conversão automática entre escalas de temperatura
- Containerização com Docker
- Pronto para deploy no Google Cloud Run

## 🚀 Como usar

### Endpoint Principal

```
GET /weather?cep={CEP}
```

### Domain Google Cloud
```
https://weather-forecast-604797684824.southamerica-east1.run.app
```
#### Exemplo de Uso
```
curl "https://weather-forecast-604797684824.southamerica-east1.run.app/weather?cep=01001000"
```
### Exemplos de Uso

#### Sucesso (200)
```bash
curl "http://localhost:8080/weather?cep=01310100"
```

**Resposta:**
```json
{
  "temp_C": 25.0,
  "temp_F": 77.0,
  "temp_K": 298.0
}
```

#### CEP Inválido (422)
```bash
curl "http://localhost:8080/weather?cep=123"
```

**Resposta:**
```json
{
  "message": "invalid zipcode"
}
```

#### CEP Não Encontrado (404)
```bash
curl "http://localhost:8080/weather?cep=00000000"
```

**Resposta:**
```json
{
  "message": "can not find zipcode"
}
```

**Resposta:** `OK`

## 🛠️ Configuração e Instalação

### Pré-requisitos

- Go 1.22+
- Docker e Docker Compose
- Chave da API WeatherAPI (gratuita em https://www.weatherapi.com/)

### 1. Clone o repositório

```bash
git clone <repository-url>
cd weather_forecast
```

### 2. Configure a variável de ambiente

Obtenha sua chave gratuita em [WeatherAPI](https://www.weatherapi.com/) e configure:

```bash
export WEATHER_API_KEY="sua_chave_aqui"
```

### 3. Executar localmente

#### Opção A: Com Go

```bash
go mod tidy
go run main.go
```

#### Opção B: Com Docker

```bash
docker build -t weather-forecast .
docker run -p 8080:8080 -e WEATHER_API_KEY="sua_chave_aqui" weather-forecast
```

#### Opção C: Com Docker Compose

```bash
# Edite o docker-compose.yml com sua chave da API
docker-compose up --build
```

## 🔧 Variáveis de Ambiente

| Variável | Descrição | Obrigatória | Padrão |
|----------|-----------|-------------|---------|
| `WEATHER_API_KEY` | Chave da API WeatherAPI | ✅ Sim | - |

## 📊 APIs Utilizadas

### ViaCEP
- **URL:** https://viacep.com.br/
- **Uso:** Busca de informações de endereço por CEP
- **Gratuita:** Sim
- **Limite:** Sem limite oficial

### WeatherAPI
- **URL:** https://www.weatherapi.com/
- **Uso:** Consulta de dados meteorológicos
- **Gratuita:** Sim (até 1M requests/mês)
- **Documentação:** https://www.weatherapi.com/docs/

## 🔄 Conversões de Temperatura

As conversões são realizadas automaticamente:

- **Celsius para Fahrenheit:** `F = C × 1.8 + 32`
- **Celsius para Kelvin:** `K = C + 273`
