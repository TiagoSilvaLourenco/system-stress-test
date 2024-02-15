# Use a imagem oficial do Golang como base
FROM golang:latest

# Copie o código fonte para o diretório de trabalho no contêiner
COPY . /go/src/app
WORKDIR /go/src/app

# Baixe as dependências do Go
RUN go mod download

# Construa o executável
RUN go build -o main .

# Exponha a porta (se necessário)
EXPOSE 8080

# Comando padrão para executar o aplicativo quando o contêiner for iniciado
ENTRYPOINT ["./main"]