Language: [EN](https://github.com/EvansTrein/RESTful_exchangerServer/blob/main/README.md) [RU](https://github.com/EvansTrein/RESTful_exchangerServer/blob/main/readmeRU.md)

<div align="center">
  <h1>Currency exchanger API</h1>
</div>

<div>
  <h2>How do I start it up?</h2>
</div>

This runs in docker, so you need: Postgres, Redis and a 3rd party gRPC service to run (the image will be downloaded from my dockerHub, it weighs 50 MB)

Clone or download the repository - type `make run-docker-compose`

If you don't use make - type `docker compose --env-file config.env up --build -d`

go to `http://localhost:8000/swagger/index.html`

<div>
  <h2>How's it work?</h2>
</div>

Upon registration, the user is automatically created accounts in currencies (USD, EUR, CNY, RUB). Each account can be interacted with - replenish, debit, exchange one currency for another. For these operations, it is necessary to log in (JWT token is issued). The exchange rate comes from the gRPC service and is cached so that you don't have to go to the gRPC service again when you request it again. 

<div>
  <h2>What's being used here and how?</h2>
</div>

**Framework** - <u>Gin</u>, originally, this was a test assignment from a company and Gin was like a requirement. But then, the assignment lost relevance, but I had almost done it, by this point. It could have been done without it. Each handler has one, its own method. That is, a handler can do one thing. 

**Logger** - <u>slog</u>, but its own wrapper is written. 

**Database** - <u>Postgres</u>, 3 tables. Users, currencies and accounts (one-to-one relationship, one user can have one account in each currency). The tables are created via migrations at server startup (we are talking about running in docker, there is a separate command to run migrations manually), using `github.com/golang-migrate/migrate/v4`. Currencies are added by a separate migration. When working with accounts, transactions and ACID are used so that the business logic is not broken.

**gRPC server** - written by myself, `https://github.com/EvansTrein/gRPC_exchangerServer`. From it we get currency rates for exchange. The server's response is cached so that we don't have to go to it every time.

**Сache** - <u>Redis</u>, total of 2 operations. Save by key, retrieve by key.

**Service Auth** - registration, issuing JWT token for access to protected resources and possibility to delete user. The service has a separate, specific for it, database interface, it contains only those methods that it needs. Middleware is used to check access during requests.

**Service Wallet** - responsible for account operations. Deposit, withdraw, exchange, get current balance, get exchange rate for all currencies. It interacts with gRPC server as a client, for this purpose a separate interface is written. The service has a separate, specific for it, database interface, it has only those methods that it needs. Also, it has a separate interface for accessing the cache.

Structures (models folder) are used for data transfer between all the above described. Tests and swagger documentation are written.