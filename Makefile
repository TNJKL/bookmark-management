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
COVERAGE_THRESHOLD = 70

# Chạy toàn bộ unit test và integration test
test:
	go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
       if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
          echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
          exit 1; \
       else \
          echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
       fi


docker-build:
	docker build -t bookmark-service-test:latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down










