# 📦 CryptoLabs Akademi — API Documentation

**Base URL:** `https://backend.cryptolabsakademi.site/api`  
**Auth:** Bearer Token (via `Authorization: Bearer <token>` header)

---

## 🔐 Authentication

### POST `/api/register`
Register a new user.

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "secret123",
  "nama_lengkap": "John Doe",
  "nama_discord": "johndoe#1234",
  "nomor_hp": "08123456789",
  "from": "instagram"
}
```

**Response `200`:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "role": "user",
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

### POST `/api/login`
Login with username/email and password.

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "login": "johndoe",
  "password": "secret123"
}
```
> `login` can be either **username** or **email**.

**Response `200`:**
```json
{
  "message": "login success",
  "token": "abc123...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "user@example.com",
    "role": "user"
  }
}
```

---

### PUT `/api/update-profile`
Update authenticated user's profile.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "nama_lengkap": "John Doe Updated",
  "nama_discord": "johndoe#9999",
  "nomor_hp": "08129999999",
  "from": "twitter"
}
```

**Response `200`:** Updated user object.

---

### GET `/api/me`
Get the currently authenticated user's info.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "user@example.com",
  "role": "user",
  "member_expired_at": "2026-01-01T00:00:00Z",
  "membershipStatus": "active"
}
```

---

## 👥 Users (Admin Only)

### GET `/api/users`
Get all users.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "data": [
    { "id": 1, "username": "johndoe", "email": "...", "role": "user" }
  ]
}
```

---

### POST `/api/users`
Create a new user (admin).

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "username": "newuser",
  "email": "new@example.com",
  "password": "pass123",
  "role": "user",
  "phone": "08123456789",
  "name": "New User"
}
```

**Response `201`:**
```json
{
  "data": { "id": 2, "username": "newuser", ... }
}
```

---

### PUT `/api/users/:id`
Update a user by ID.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body (all fields optional):**
```json
{
  "username": "updateduser",
  "email": "updated@example.com",
  "password": "newpassword",
  "role": "admin",
  "phone": "08111111111",
  "name": "Updated Name"
}
```

**Response `200`:**
```json
{
  "data": { "id": 2, "username": "updateduser", ... }
}
```

---

### DELETE `/api/users/:id`
Delete a user by ID (admin cannot be deleted).

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "user deleted"
}
```

---

## 💳 Payments

### GET `/api/payments`
Get payments. Admins see all; users see their own.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "method": "bank",
      "status": "pending",
      "amount": 200000,
      "month_count": 1,
      "bukti": "bukti/1_...",
      "discord_code": null
    }
  ]
}
```

---

### POST `/api/checkout`
Submit a membership checkout with payment proof.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `method` | string | ✅ | Payment method (e.g. `bank`, `ewallet`, `crypto`) |
| `month_count` | integer | ✅ | Number of months to subscribe |
| `bukti` | file | ✅ | Payment proof image/PDF (jpg, jpeg, png, webp, heic, pdf) |
| `coupon_code` | string | ❌ | Optional coupon code |

**Response `200`:**
```json
{
  "message": "Checkout berhasil, menunggu konfirmasi admin",
  "payment": {
    "id": 1,
    "status": "pending",
    "amount": 200000,
    "original_amount": 200000,
    "discount": 0
  }
}
```

---

### POST `/api/payments/:id/approve`
Approve a payment (admin only). Generates a Discord code.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Payment berhasil di-approve",
  "discord_code": "ABCDEFGH"
}
```

---

### POST `/api/payments/:id/reject`
Reject a payment (admin only).

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "reason": "Bukti pembayaran tidak valid"
}
```

**Response `200`:**
```json
{
  "message": "Payment berhasil ditolak"
}
```

---

### DELETE `/api/payments/:id`
Delete a payment (admin only). Cannot delete paid payments.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Payment berhasil dihapus"
}
```

---

## 🏷️ Coupons

### GET `/api/coupons/check`
Validate a coupon code publicly.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Params:**
```
?code=DISC10&month=3
```

**Response `200`:**
```json
{
  "message": "Coupon valid",
  "coupon": {
    "id": 1,
    "code": "DISC10",
    "type": "percent",
    "value": 10,
    "max_discount": 50000,
    "quota": 100,
    "used_count": 5,
    "expired_at": null
  }
}
```

---

### GET `/api/coupons`
Get all coupons (admin only).

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "data": [ { "id": 1, "code": "DISC10", ... } ]
}
```

---

### GET `/api/coupons/:id`
Get a single coupon by ID (admin only).

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "coupon": { "id": 1, "code": "DISC10", ... }
}
```

---

### POST `/api/coupons`
Create a coupon (admin only).

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "code": "SAVE20",
  "type": "percent",
  "value": 20,
  "max_discount": 100000,
  "quota": 50,
  "trigger": "checkout",
  "expired_at": "2025-12-31",
  "is_active": true,
  "min_month": 3
}
```
> `type` must be `"percent"` or `"fixed"`.

**Response `200`:**
```json
{
  "message": "Coupon berhasil dibuat",
  "coupon": { ... }
}
```

---

### PUT `/api/coupons/:id`
Update a coupon (admin only). All fields optional.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "code": "NEWSAVE",
  "type": "fixed",
  "value": 50000,
  "max_discount": null,
  "quota": 100,
  "trigger": null,
  "expired_at": "2026-06-30",
  "is_active": true,
  "min_month": 1
}
```

**Response `200`:**
```json
{
  "message": "Coupon berhasil diupdate"
}
```

---

### DELETE `/api/coupons/:id`
Delete a coupon (admin only).

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Coupon berhasil dihapus"
}
```

---

## 💰 Membership Pricing

### GET `/api/membership/pricing`
Get pricing for given month durations.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Params:**
```
?months=1,3,6,12,1000
```

**Response `200`:**
```json
{
  "aturan_db": true,
  "data": {
    "1":    { "price": 200000, "price_rupiah": "200.000" },
    "3":    { "price": 550000, "price_rupiah": "550.000" },
    "6":    { "price": 1000000, "price_rupiah": "1.000.000" },
    "12":   { "price": 1800000, "price_rupiah": "1.800.000" },
    "1000": { "price": 2500000, "price_rupiah": "2.500.000" }
  }
}
```

---

## 📐 Rule Pricing (Admin Only)

### GET `/api/rule-pricing`
Get all pricing rules.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "success": true,
  "message": "Berhasil Mendapatkan Data Aturan Pembayaran",
  "data": [
    { "id": 1, "min_month": 1, "max_month": 2, "total_price": 200000, "is_active": true }
  ]
}
```

---

### GET `/api/rule-pricing/:id`
Get a single pricing rule by ID.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{ "id": 1, "min_month": 1, "max_month": 2, "total_price": 200000 }
```

---

### POST `/api/rule-pricing`
Create a pricing rule.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "min_month": 1,
  "max_month": 2,
  "total_price": 200000,
  "is_active": true
}
```

**Response `200`:**
```json
{
  "data": { "id": 5, "min_month": 1, "max_month": 2, "total_price": 200000 }
}
```

---

### PUT `/api/rule-pricing/:id`
Update a pricing rule.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "min_month": 1,
  "max_month": 2,
  "total_price": 250000,
  "is_active": true
}
```

**Response `200`:**
```json
{
  "data": { ... }
}
```

---

### DELETE `/api/rule-pricing/:id`
Delete a pricing rule.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Rule pricing berhasil dihapus"
}
```

---

## 📚 Module Groups

### GET `/api/module-groups`
Get all module groups. Non-members only see free content.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Crypto Basics",
      "description": "...",
      "is_active": true,
      "for_member": false
    }
  ]
}
```

---

### POST `/api/module-groups`
Create a module group (admin only).

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "Title": "Advanced Trading",
  "Description": "Materi trading lanjutan",
  "IsActive": true,
  "ForMember": true
}
```

**Response `200`:**
```json
{
  "message": "Module group berhasil dibuat",
  "data": { ... }
}
```

---

### PUT `/api/module-groups/:id`
Update a module group (admin only).

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "Title": "Advanced Trading v2",
  "Description": "Updated description",
  "IsActive": true,
  "ForMember": true
}
```

**Response `200`:**
```json
{
  "message": "Module group berhasil diupdate"
}
```

---

### DELETE `/api/module-groups/:id`
Delete a module group (admin only).

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Module group berhasil dihapus"
}
```

---

## 🎥 Modules

### GET `/api/modules/group/:group_id`
Get all modules in a group. Non-members see only free modules. Admins see `youtube_id`; regular users do not.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "data": {
    "id": "uuid",
    "title": "Crypto Basics",
    "modules": [
      {
        "id": "uuid",
        "title": "Intro to Bitcoin",
        "description": "...",
        "youtube_id": "abc123",
        "for_member": false
      }
    ]
  }
}
```

---

### POST `/api/modules`
Create a module (admin only).

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "module_group_id": "uuid-group",
  "title": "DeFi Introduction",
  "description": "Pengenalan DeFi",
  "youtube_id": "dQw4w9WgXcQ",
  "is_active": true,
  "for_member": true
}
```

**Response `200`:**
```json
{
  "message": "Module berhasil dibuat",
  "data": { ... }
}
```

---

### PUT `/api/modules/:id`
Update a module (admin only).

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Updated DeFi",
  "description": "Updated description",
  "youtube_id": "newVideoID",
  "is_active": true,
  "for_member": true
}
```

**Response `200`:**
```json
{
  "message": "Module berhasil diupdate"
}
```

---

### DELETE `/api/modules/:id`
Delete a module (admin only).

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Module berhasil dihapus"
}
```

---

## ▶️ Module Streaming

### GET `/api/stream/:module_id`
Get the YouTube embed URL for a module. Requires active membership.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "title": "DeFi Introduction",
  "embed_url": "https://www.youtube.com/embed/abc123?controls=1&..."
}
```

---

## 📊 Module Progress

### POST `/api/module-progress`
Track user progress on a module.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "module_id": "uuid-module",
  "status": "completed"
}
```
> `status` must be `"watching"` or `"completed"`.

**Response `200`:**
```json
{
  "message": "Progress tersimpan"
}
```

---

## 💸 Expenses (Admin Only)

### GET `/api/expenses`
Get all expenses with optional filters.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Params:**
```
?month=2026-01&category=operasional
```

**Response `200`:**
```json
{
  "data": [
    {
      "id": 1,
      "description": "Server hosting",
      "amount": 500000,
      "category": "operasional",
      "spent_at": "2026-01-15T00:00:00Z"
    }
  ]
}
```

---

### GET `/api/expenses/:id`
Get a single expense by ID.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "expense": { "id": 1, "description": "Server hosting", ... }
}
```

---

### POST `/api/expenses`
Create an expense record.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "description": "Domain renewal",
  "amount": 250000,
  "category": "operasional",
  "spent_at": "2026-01-10"
}
```

**Response `200`:**
```json
{
  "message": "Pengeluaran berhasil dibuat",
  "expense": { ... }
}
```

---

### PUT `/api/expenses/:id`
Update an expense.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body (all optional):**
```json
{
  "description": "Updated description",
  "amount": 300000,
  "category": "marketing",
  "spent_at": "2026-01-20"
}
```

**Response `200`:**
```json
{
  "message": "Pengeluaran berhasil diupdate",
  "expense": { ... }
}
```

---

### DELETE `/api/expenses/:id`
Delete an expense.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Pengeluaran berhasil dihapus"
}
```

---

## 📣 Announcements (Admin Only)

### GET `/api/announcements`
Get all announcements.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Maintenance Notice",
      "content": "...",
      "type": "info",
      "target": { "audience": "all" }
    }
  ]
}
```

---

### POST `/api/announcements`
Create and send an announcement.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Server Maintenance",
  "content": "Halo [NAMA], server akan maintenance besok.",
  "type": "warning",
  "channels": ["email", "discord"],
  "target": { "audience": "all" }
}
```
> `channels` options: `"email"`, `"discord"`  
> `target.audience` options: `"all"`, `"active"`, `"expired"`  
> Use `[NAMA]` in content as a placeholder for the user's username.

**Response `200`:**
```json
{
  "message": "Announcement berhasil dibuat"
}
```

---

## 🤖 Auto Messager (Admin Only)

### GET `/api/automessager`
Get all auto messager configs.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:** Array of AutoMessager objects.

---

### GET `/api/automessager/:id`
Get a single auto messager by ID.

**Headers:**
```
Authorization: Bearer <token>
```

---

### POST `/api/automessager`
Create an auto messager.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Daily Reminder",
  "message": "Halo semua! Jangan lupa belajar hari ini 🚀",
  "bot_id": "bot-uuid",
  "server_id": "discord-server-id",
  "channel_id": "discord-channel-id",
  "run_time": "08:30",
  "days_of_week": ["Mon", "Tue", "Wed", "Thu", "Fri"],
  "timezone": "Asia/Jakarta",
  "image": null
}
```

**Response `200`:** Created AutoMessager object.

---

### PUT `/api/automessager/:id`
Update an auto messager.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:** Same structure as POST.

**Response `200`:**
```json
{
  "message": "Updated successfully"
}
```

---

### PATCH `/api/automessager/:id/toggle`
Toggle active/inactive state.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:** Updated AutoMessager object.

---

### DELETE `/api/automessager/:id`
Delete an auto messager.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Deleted successfully"
}
```

---

## 🤖 Bots (Admin Only)

### GET `/api/bots`
Get all bots (tokens masked).

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
[
  { "id": "uuid", "name": "Main Bot", "token": "********", "is_active": true }
]
```

---

### GET `/api/bots/:id`
Get a single bot by ID.

**Headers:**
```
Authorization: Bearer <token>
```

---

### POST `/api/bots`
Create a bot.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "name": "News Bot",
  "token": "your-discord-bot-token"
}
```

**Response `200`:** Created Bot object (token masked).

---

### PUT `/api/bots/:id`
Update a bot.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Updated Bot Name",
  "token": "new-token-if-changing"
}
```

**Response `200`:**
```json
{
  "message": "Bot updated"
}
```

---

### PATCH `/api/bots/:id/toggle`
Toggle bot active/inactive.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:** Updated Bot object.

---

### DELETE `/api/bots/:id`
Delete a bot.

**Headers:**
```
Authorization: Bearer <token>
```

**Response `200`:**
```json
{
  "message": "Bot deleted"
}
```

---

## 🔔 Notifications (SSE)

### GET `/api/notif`
Subscribe to real-time notifications via Server-Sent Events.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Params:**
```
?user_id=1
```

**Response:** `text/event-stream` stream.

```
data: {"id":"uuid","title":"Pembayaran Berhasil","message":"...","type":"success","user_id":1}
```

---

## 🌐 Webhook

### POST `/webhook/test`
Test webhook endpoint.

**Body:** Any JSON or raw body.

**Response `200`:**
```json
{
  "message": "Webhook test successful",
  "data": "..."
}
```

---

## 🏥 Health Check

### GET `/health`
Check server status.

**Response `200`:**
```json
{
  "status": "oke"
}
```

### HEAD `/health`
Lightweight health check. Returns `200 OK` with no body.

---

## ⚠️ Common Error Responses

| Status | Meaning |
|--------|---------|
| `400` | Bad Request — invalid body or missing fields |
| `401` | Unauthorized — missing or invalid token |
| `403` | Forbidden — insufficient permissions |
| `404` | Not Found — resource does not exist |
| `429` | Too Many Requests — rate limit exceeded (10 req / 10s) |
| `500` | Internal Server Error |

**Error body format:**
```json
{
  "error": "Human readable error message"
}
```
