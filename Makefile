.PHONY: dev test mock swag clean

#Các phần liên quan tới Swagger  và MakeFile em nhờ AI gen thử còn  code bài tập là em dựa vào bài học rồi tự viết lại ạ


# Chạy ứng dụng ở chế độ local
dev:
	go run cmd/api/main.go

# Chạy toàn bộ unit test và integration test
test:
	go test ./internal/... -v

# Tạo lại các file Mock bằng Mockery
mock:
	go generate ./...

# Tạo lại tài liệu Swagger UI
swag:
	swag init -g cmd/api/main.go

