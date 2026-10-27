# Freelance-Marketplace-Ecosystem-

The Scenario: Connecting clients with freelancers, posting jobs, managing project milestones, contracts, and payment releases. DB Challenge: Handling multi-party transactions, milestone tracking, and user-rating aggregation tables.


Database Project Report Template
Due Date: [Insert Date - Week ….]
Project ID & Title: Freelance Marketplace Ecosystem (Mini-Upwork)
A. Project Identity
Team Name: chậm mà chắc
Team Members: Bùi Hữu Hải Đăng (n24dece058@student.ptithcm.edu.vn)
Vũ Anh Khoa (n24dece078@student.ptithcm.edu.vn)
Lê Trần Hoàng Minh (n24dece082@student.ptithcm.edu.vn)
Project Title: [A descriptive title for your specific implementation]

B. Report structure

1. Introduction & Project Scope (Adapted from ISO/IEC/IEEE 29148)
   1.1 System Objective
   1.2 Business Rules & Constraints
2. Database Design (ISO/IEC 19505 / IE Standards)
   2.1 Conceptual Model (ER/EER Diagram)
   2.2 Logical Schema Mapping
   2.3 Normalization Verification (Formal proofs for 1NF -> 3NF/BCNF)
3. Data Dictionary (Adapted from ISO/IEC 11179)
   [Tables and attribute metadata]
4. Database Implementation (SQL Style Guide Compliant)
   4.1 DDL Script (Tables, Views, Indexes, Triggers)
   4.2 Advanced Queries & Performance Test Cases
5. Verification & Security
   5.1 Test cases showing constraints prevent bad data insertion
   5.2 Role-Based Access Control (RBAC) definition (GRANT/REVOKE statements)

## Go + Gin REST API (MySQL)

The API uses the `freelance_platform` database created by [`scripts/database.sql`](./scripts/database.sql). It provides CRUD for job posts.

### Run

1. Install/import `scripts/database.sql` into MySQL.
2. Set `MYSQL_DSN` to a MySQL DSN in `.env` (the file is gitignored) or in your environment. The app loads `.env` automatically; `parseTime=true` is required for timestamps:

   ```bash
   MYSQL_DSN='your_user:your_password@tcp(127.0.0.1:3306)/freelance_platform?parseTime=true&charset=utf8mb4'
   ```

3. Start the server:

   ```bash
   go run ./cmd/server
   ```

   The API listens on `http://localhost:8085`.

### Endpoints

All responses use JSON. IDs in the request body must reference existing rows in the database.

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/api/v1/health` | Check API and database connectivity |
| GET | `/api/v1/jobs` | Browse jobs; optional `q`, `hire_manager_id`, `skill_id` filters |
| GET | `/api/v1/jobs/:id` | Get a job |
| POST | `/api/v1/jobs` | Create a job |
| PUT | `/api/v1/jobs/:id` | Replace a job |
| DELETE | `/api/v1/jobs/:id` | Delete a job |

Job request example:

```json
{
  "hire_manager_id": 1,
  "duration_id": 1,
  "complexity_id": 1,
  "description": "Build a responsive landing page",
  "main_skill_id": 1,
  "payment_type_id": 1,
  "payment_amount": 500,
  "title": "Frontend developer needed"
}
```

Import [`docs/postman_collection.json`](./docs/postman_collection.json) into Postman to try the endpoints. Before creating a job, update the sample IDs to records that exist in your database. The collection saves the created job ID for the subsequent get/update/delete requests. Successful creates return `201`, reads/updates return `200`, deletes return `204`; invalid input returns `400`, missing records `404`, and constraint conflicts `409`.
