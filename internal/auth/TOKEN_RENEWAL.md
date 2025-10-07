# Token Renewal Strategy

## Overview

This document outlines the **silent token renewal strategy** used to manage authentication for both **server** and **API** access. The goal is to provide seamless session continuity while maintaining strong security and enforcing inactivity limits.

---

## Authentication Flow

### 1. Login Phase

Upon a successful login:
- A **session token** and a **refresh token** are issued.
- The **session token** is sent to the client as an `HttpOnly` cookie.
- The **refresh token** is securely stored in the server-side **database**.

This design ensures that:
- The client cannot access or manipulate token data directly.
- The refresh token remains protected against theft via client-side attacks.

---

## Silent Renewal

### 2. Request Handling

When the client makes authenticated requests:
- The server checks the **session token** in the request’s `HttpCookie`.
- If the session token is valid, the request proceeds as normal.
- If the session token is **nearing expiration**, a **silent renewal** process is triggered.

### 3. Renewal Process

The silent renewal process works as follows:

1. The server uses the **refresh token** (retrieved from the database) to issue:
   - A **new session token** (sent to the client).
   - A **new refresh token** (stored in the database, replacing the old one).

2. The **inactivity period** is reset at this point:
   - Each renewal extends the session lifetime.
   - The refresh token maintains a **90-day rolling expiration** window — if unused for 90 days, it becomes invalid.

---

## Security Considerations

- **HttpOnly Cookies:** Prevent JavaScript access to session tokens, reducing XSS risk.
- **Database-Stored Refresh Tokens:** Ensures tokens can be revoked or rotated on logout or breach.
- **Token Rotation:** Each renewal invalidates the previous refresh token, mitigating replay attacks.
- **Inactivity Enforcement:** The 90-day limit ensures inactive accounts naturally expire, reducing risk exposure.

---

## Summary

| Component        | Storage Location | Expiration | Renewal Trigger | Notes |
|------------------|------------------|-------------|------------------|-------|
| Session Token    | Client (HttpOnly cookie) | Short-lived (e.g., 15–30 mins) | Near-expiry during requests | Used for active session validation |
| Refresh Token    | Server (Database) | 90 days (rolling) | On session renewal | Rotated on each renewal to reset inactivity |

This approach ensures **secure**, **scalable**, and **user-friendly** authentication with minimal friction during token renewal.
