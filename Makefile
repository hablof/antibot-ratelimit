# тесты долгие, 15+ секунд
test:
	go test ./... 

run:
	docker-compose up -d
stop:
	docker-compose down