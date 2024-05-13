cd ./server && go build && cd .. && go run ./server/server.go --mainport 50003 & 
go run ./server/server.go --mainport 50002 & 
go run ./server/server.go --ports "50002,50003"