# Bookmark Management Service

Dự án Bookmark Management cung cấp dịch vụ rút gọn liên kết (Shorten URL) sử dụng **Golang (Gin)** và **Redis** làm lưu trữ tạm thời.

---

## API Endpoints (Rút gọn URL)

### 1. Tạo Link Rút Gọn (Shorten URL)
* **Endpoint:** `POST /v1/links/shorten`
* **Content-Type:** `application/json`
* **Request Body:**
  ```json
  {
    "url": "https://google.com",
    "exp": 60000
  }
  ```
  *(Lưu ý: `exp` là thời gian hết hạn tính bằng giây, bắt buộc phải `>= 50000`)*
* **Response (200 OK):**
  ```json
  {
    "code": "Songoku",
    "message": "Shorten URL generated successfully!"
  }
  ```

### 2. Chuyển Hướng Link (Redirect URL)
* **Endpoint:** `GET /v1/links/redirect/{code}`
* **Params:** `code` (mã rút gọn nhận được từ API trên)
* **Response (302 Found):** Tự động chuyển hướng (Redirect) về URL gốc ban đầu.

---

## Chạy Ứng Dụng Với Docker

Dự án đã cấu hình sẵn Dockerfile (Multi-stage build) và Docker Compose giúp khởi chạy nhanh chóng cả server API và Redis.

* **Cách chạy toàn bộ hệ thống:**
  ```bash
  make docker-up
  ```
  Lệnh này sẽ khởi động container Redis (port `6379`) và Bookmark Service (port `8080`).

* **Cách dừng hệ thống:**
  ```bash
  make docker-down
  ```

* **Xây dựng lại Docker Image:**
  ```bash
  make docker-build
  ```

---

## Các Lệnh Trong Makefile

Dưới đây là các lệnh tiện ích được cấu hình sẵn trong `Makefile` để phục vụ phát triển cục bộ (local development):

| Lệnh | Ý nghĩa |
|---|---|
| `make run` | Khởi chạy nhanh dự án bằng lệnh Go thông thường (`go run cmd/api/main.go`). |
| `make swag` | Tạo/Cập nhật lại tài liệu Swagger UI tự động trong thư mục `docs/`. |
| `make dev-run` | Vừa cập nhật Swagger vừa chạy ứng dụng (`swag run`). |
| `make mock` | Tự động quét và tạo các Mock struct phục vụ viết Unit Test. |
| `make test` | Chạy toàn bộ test, tính toán độ phủ code (coverage, yêu cầu `>= 80%`) và xuất ra file báo cáo HTML `coverage.html`. |