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

   Also configure `SECRET_KEY` and `TOKEN_LIFESPAN` in `.env` or the environment. `SECRET_KEY` is required to sign and verify JWTs; use a long, random value and do not commit it. `TOKEN_LIFESPAN` is the token lifetime in hours and must be a positive integer.

3. Set a port and start the server:

   ```bash
   export PORT=8085
   export SECRET_KEY='your-long-random-secret'
   export TOKEN_LIFESPAN=24
   export ADMIN_USER_IDS=1,5
   go run ./cmd/server
   ```

   The API listens on `http://localhost:8085`. Login returns a `token` field. Send it to protected endpoints as `Authorization: Bearer <token>`. `ADMIN_USER_IDS` is an optional comma-separated allowlist of `user_account.id` values with admin override permissions.

   User profile endpoints operate on the authenticated account; the account ID is read from the JWT, not from the URL:

   - `GET /user/me`
   - `PUT /user/me`
   - `DELETE /user/me`
   - `GET /user/me/freelancer`
   - `GET /user/me/hirer`

### Reading the Job CRUD code

Follow one request through these folders:

1. [`internal/model/model.go`](./internal/model/model.go) defines the Job fields used in JSON and SQL.
2. [`internal/routes/routes.go`](./internal/routes/routes.go) maps HTTP methods and URLs to controller functions.
3. [`internal/controller/jobController.go`](./internal/controller/jobController.go) validates input and runs SQL queries against the `job` table.

Hire Managers can create jobs and edit/delete only jobs owned by their profile. Freelancers can submit applications. Users in `ADMIN_USER_IDS` can manage any job. `GET /job` and `GET /job/:jobid` remain public.

The job routes are `POST /job`, `PUT /job/:jobid`, and `DELETE /job/:jobid`. Create and update accept the job reference IDs and an optional `required_skill_ids` array. The main skill is automatically included in the job's required skills. `POST /job/:jobid/apply` accepts `payment_type_id`, `payment_amount`, and an optional `freelancer_comment`; it creates a proposal with `Pending` status. Ensure `status_catalog` contains a row whose `status_name` is `Pending`.

| Operation | Request | Repository action |
|---|---|---|
| Create | `POST /job` | Insert job and required skills |
| Update | `PUT /job/:jobid` | Replace job fields and required skills |
| Delete | `DELETE /job/:jobid` | Delete job |

### Endpoints

All responses use JSON. IDs in the request body must reference existing rows in the database.

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/api/v1/health` | Check API and database connectivity |
| GET | `/job` | Browse jobs with `page` and `limit` query parameters |
| GET | `/job/:jobid` | Get a job |
| POST | `/job` | Create a job |
| PUT | `/job/:jobid` | Replace a job |
| DELETE | `/job/:jobid` | Delete a job |
| POST | `/job/:jobid/apply` | Apply to a job as a freelancer |

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
  "title": "Frontend developer needed",
  "required_skill_ids": [1, 2, 3]
}
```

Application request example:

```json
{
  "payment_type_id": 1,
  "payment_amount": 500,
  "freelancer_comment": "I have experience delivering similar projects."
}
```

Import [`docs/postman_collection.json`](./docs/postman_collection.json) into Postman to try the endpoints. Before creating a job, update the sample IDs to records that exist in your database. The collection saves the created job ID for the subsequent get/update/delete requests. Successful creates return `201` with the new `job_id`; reads, updates, and deletes return `200`. Invalid input returns `400` and missing records return `404`.

### Swagger documentation

Swagger documentation is generated from Swaggo annotations in the API code. Run `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/server/main.go -o internal/swagger/docs --parseInternal --parseDependency` after changing annotations, then start the server and open `http://localhost:8085/swagger/index.html`. The page documents and lets you try the currently registered endpoints (`/health`, `/auth/register`, and `/auth/login`), and protected endpoints show an Authorize button for the `Authorization: Bearer <token>` header.
