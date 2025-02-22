mock:
	mockgen -package mockrepo -destination internal/mock/user/user.go  user/internal/repo UserImply
	mockgen -package mockapi -destination internal/mock/api/api.go  user/adapter/apihook Connector
	mockgen -package mockemailer -destination internal/mock/emailer/emailer.go  user/adapter/emailer Emailer
lint:
	golangci-lint run
compose-dev:
	docker compose -f docker-dev.yml watch
test:
	go test -cover ./...
server:
	go run main.go
swagger:
	swag init
migrate:
	migrate create -ext sql -dir migrations -seq -digits 6 $(NAME)
