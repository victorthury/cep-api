# cep-api

Uma api que pega o cep ou da brasil api ou do via cep. Caso demore mais de 1
segundo, retorna `"timeout"`.

O foco deste repositório é exercitar multithreading com go

## Como executar

Entre na pasta `cmd/server` e execute o seguinte comando:

```
go run main.go
```

##

Em um cliente http, execute a request em `test/cep.http`
