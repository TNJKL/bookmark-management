.PHONY: run swag dev-run test mock

GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMG_TAG := latest
IMG_NAME=ththi/bookmark_service
ifneq ($(GIT_TAG),)
   IMG_TAG := $(GIT_TAG)
endif


export IMG_TAG


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



COVERAGE_EXCLUDE=mocks|main.go|docs|test|config.go
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
	docker build -t $(IMG_NAME):$(IMG_TAG) .

docker-up:
	docker compose up -d

docker-down:
	docker compose down


COVERAGE_FOLDER=./test-output
docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="${COVERAGE_EXCLUDE}" --build-arg COVERAGE_THRESHOLD="${COVERAGE_THRESHOLD}" --progress=plain --target test -t test:test --output ./test-output .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	  echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	  exit 1; \
    else \
	  echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
    fi


DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$${DOCKER_PASSWORD}" | docker login -u "$${DOCKER_USERNAME}" --password-stdin

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)








