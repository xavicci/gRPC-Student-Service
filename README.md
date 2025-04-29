"# gRPC-Student-Service"
go get google.golang.org/protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go get github.com/lib/pq

-- go get google.golang.org/grpc

export PATH=$PATH:$(go env GOPATH)/bin 

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative studentpb/student.proto

cd database
docker build . -t grpc-students
docker run -p 5432:5432 --name mypostgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=mibase test-grpc-db
go run .\server-student\main.go

docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=mibase test-grpc-db

# Primero se crea el patron repositorio y luego PostGresRepository: 
    REPO: contiene el codigo de las funciones para modificar la base de datos.
     - Utiliza models
# Luego se llama al servicio que carga el repositorio (usa interface):
    SERVICIO: contiene la comunicacion con el servicio GRPC (req) y llama las funciones del repo cargando la info que recibe del *req
     -  Utiliza models
     -  Utiliza repo
     -  Utiliza proto

# Se crea el NewServerGrpc, RegisterServiceServer y el Reflection.Register
# Luego se inicia el servico de  grpc