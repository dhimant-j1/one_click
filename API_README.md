# Insurance Management API - OpenAPI Specification

This directory contains the OpenAPI 3.0.3 specification for the Insurance Management API.

## Files

- `openapi.yaml` - Complete OpenAPI specification for all API endpoints

## API Overview

The Insurance Management API provides endpoints for:

### Authentication

- **POST** `/api/auth/login` - User authentication
- **POST** `/api/auth/reset-password` - Password reset request
- **PUT** `/api/auth/reset-password/{token}` - Password update with token

### Master Admin Operations

All master admin endpoints require authentication and MasterAdmin role.

#### Agencies Management

- **GET** `/api/master-admin/agencies` - List all agencies
- **POST** `/api/master-admin/agencies` - Create new agency
- **PUT** `/api/master-admin/agencies/{id}` - Update agency
- **DELETE** `/api/master-admin/agencies/{id}` - Delete agency

#### Insurance Providers Management

- **GET** `/api/master-admin/providers` - List all providers
- **POST** `/api/master-admin/providers` - Create new provider
- **PUT** `/api/master-admin/providers/{id}` - Update provider
- **DELETE** `/api/master-admin/providers/{id}` - Delete provider

#### Provider Access Management

- **POST** `/api/master-admin/provider-access` - Grant agency access to provider
- **DELETE** `/api/master-admin/provider-access` - Revoke agency access to provider

#### User Management

- **GET** `/api/master-admin/users` - List all users
- **POST** `/api/master-admin/users` - Create new user
- **PUT** `/api/master-admin/users/{id}` - Update user
- **DELETE** `/api/master-admin/users/{id}` - Delete user

#### Reports

- **GET** `/api/master-admin/reports/quotes` - Get quotes report
- **GET** `/api/master-admin/reports/policies` - Get policies report

#### Authentication

- **POST** `/api/master-admin/auth/logout` - Logout user

## Data Models

### User Roles

- `MasterAdmin` - Full system access
- `AgencyAdmin` - Agency-level administration
- `LocationAdmin` - Location-level administration
- `Agent` - Insurance agent
- `Customer` - Insurance customer

### Quote Status

- `Draft` - Quote in draft state
- `Presented` - Quote presented to customer
- `Bound` - Quote bound to policy

### Policy Status

- `Active` - Active policy
- `Expired` - Expired policy
- `Cancelled` - Cancelled policy

## Authentication

The API uses JWT Bearer tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Usage

1. **Login** to get a JWT token:

   ```bash
   curl -X POST http://localhost:8081/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@example.com","password":"password"}'
   ```

2. **Use the token** for authenticated requests:
   ```bash
   curl -X GET http://localhost:8081/api/master-admin/agencies \
     -H "Authorization: Bearer <your-token>"
   ```

## Viewing the Specification

You can view and interact with this API specification using:

1. **Swagger UI** - Import the `openapi.yaml` file into [Swagger Editor](https://editor.swagger.io/)
2. **Postman** - Import the OpenAPI specification to generate a collection
3. **VS Code** - Use OpenAPI extensions for syntax highlighting and validation

## Server Information

- **Development Server**: `http://localhost:8081/api`
- **Base Path**: `/api`

## Notes

- All master admin endpoints require both authentication and MasterAdmin role
- Request/response content type is `application/json`
- Timestamps use ISO 8601 format
- Database uses PostgreSQL with proper ENUM types for status fields
