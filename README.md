
# **🔥 Section 1: Java Backend Development & System Architecture (80 Points)**

## **1.1 - User Registration API (30 Points)**
**📌 Mô tả:**  
Bạn cần triển khai một **REST API** để xử lý đăng ký người dùng. Service cần:
- **Validate input:**
    - `username` phải là **alphanumeric**.
    - `email` phải hợp lệ.
    - `password` tối thiểu **8 ký tự**.
- **Lưu dữ liệu vào PostgreSQL/MySQL**.
- **Gửi email xác nhận (Mocked bằng Console Logging)**.
- **Viết Unit Test cho service**.

**📌 Yêu cầu kỹ thuật:**
- Sử dụng **Spring Boot + JPA + Hibernate / NodeJS**
- **OOP & Design Pattern:** Áp dụng **Design Pattern** cho Email Service.
- **Viết Unit Test (JUnit + Mockito / jest).**
- **Exception Handling (Custom Exception, Global Error Handler).**

👉 **Follow-up Questions:**
- **Tại sao chọn Factory Pattern thay vì Singleton?**
- **Làm thế nào để triển khai gửi email thực tế với Spring Mail?**

---

## **1.2 - Scalable CDN Architecture (20 Points)**
**📌 Mô tả:**  
Thiết kế **kiến trúc CDN** phục vụ hàng triệu người dùng trên toàn cầu.

**📌 Yêu cầu:**
- **Modularization:** Phân tách thành **Edge Servers, Cache Servers, Load Balancers**.
- **Caching Strategies:** Áp dụng **Cache Invalidation, Content Hashing, CDN Tiering**.
- **Redundancy & Failover:** Dùng **Geo Load Balancing, Active-Passive Replication**.
- **Data Consistency:** **Event-driven updates**, sử dụng **Kafka** để cập nhật real-time.

---

## **1.3 - File Compression Service (30 Points)**
**📌 Mô tả:**  
Viết một API để **nén file văn bản lớn**, sử dụng **Spring WebFlux + Streaming API**.

**📌 Yêu cầu:**
- **API `/compress`**: Nhận file → Nén → Lưu vào **S3 hoặc Disk**.
- **Sử dụng Stream API để tiết kiệm bộ nhớ**.
- **Dùng `GZIPOutputStream` hoặc `Snappy` để nén dữ liệu**.
- **Viết Unit Test**.

---

# **🔥 Section 2: Database Management & Optimization (25 Points)**

## **2.1 - SQL Query Optimization (15 Points)**
**📌 Mô tả:**
- Viết SQL truy vấn lấy **top 5 người dùng trẻ nhất**, đã đăng ký **trong 1 năm qua**.
- **Tối ưu indexing strategy** để query nhanh nhất.

---

## **2.2 - Time-Series Data Best Practices (10 Points)**
**📌 Mô tả:**  
Bạn đang lưu dữ liệu **IoT sensor** real-time.  
Hãy giải thích cách **tối ưu truy vấn dữ liệu thời gian thực** trong PostgreSQL.

**📌 Yêu cầu:**
- **Partitioning (Time-based Partitioning)**.
- **Indexing Strategies (`BRIN`, `GIN`, `HASH` indexes).**
- **Retention Policy (Data TTL + Archiving).**

---

# **🔥 Section 3: System Monitoring & Debugging (30 Points)**

## **3.1 - API Performance Monitoring (30 Points)**
**📌 Mô tả:**  
Bạn cần triển khai **middleware** để đo thời gian phản hồi API & cảnh báo nếu mất quá nhiều thời gian.

**📌 Yêu cầu:**
- **Spring Boot Filter** để log `responseTime`.
- **Cảnh báo nếu `responseTime > 200ms`**.
- **Lưu log vào PostgreSQL hoặc ELK Stack**.
- **Viết Unit Test**.

---

# **🔥 Tổng kết (135 Points)**
| **Section** | **Task** | **Time** | **Points** |
|------------|---------|---------|------------|
| Java Backend | User Registration API | 30 min | 30 |
| Java Backend | CDN Architecture | 20 min | 20 |
| Java Backend | File Compression API | 30 min | 30 |
| Database | SQL Query Optimization | 15 min | 15 |
| Database | Time-Series Best Practices | 10 min | 10 |
| Monitoring | API Performance Monitoring | 30 min | 30 |
| **Total** | 6 Tasks | **135 min** | **135 Points** |
