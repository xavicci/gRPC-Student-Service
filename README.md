"# gRPC-Student-Service"
go get google.golang.org/protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go get github.com/lib/pq

-- go get google.golang.org/grpc

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative studentpb/student.proto

cd database
docker build . -t grpc-students
docker run -p 54321:5432 --name mypostgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=mibase grpc-students
go run .\server-student\main.go
