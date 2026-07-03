.PHONY: run swag dev-run test mock

#Các phần liên quan tới Swagger  và MakeFile em nhờ AI gen thử còn  code bài tập là em dựa vào bài học rồi tự viết lại ạ


run:
	go run cmd/api/main.go

# Tạo lại tài liệu Swagger UI
swag:
	swag init -g cmd/api/main.go --output docs

dev-run: swag run



# Tạo lại các file Mock bằng Mockery
mock:
	go generate ./...



COVERAGE_EXCLUDE=mocks|main.go

# Chạy toàn bộ unit test và integration test
test:
	go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html










