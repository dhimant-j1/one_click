# Insurance Quote Application System Design Document

## Overview

This document outlines the design for a web-based insurance quote generation application based on the provided database schema. The system supports multiple user roles: MasterAdmin, AgencyAdmin, LocationAdmin, Agent, and Customer. However, as per the requirements, the application is divided into **three main portals** sharing the same database:

1. **Master Admin Portal**: For system-wide administration, managing agencies, users, and providers.
2. **Agency Admin Portal**: For agency-level management, including locations, agents, and provider access (encompassing AgencyAdmin and LocationAdmin roles).
3. **Agent Portal**: For agents to manage customers, vehicles, quotes, and policies.

Customers do not have a dedicated portal but interact via agents (e.g., through shared quotes or email notifications). All portals use the same backend APIs, with role-based access control (RBAC) enforced via the `Users.Role` field. Authentication is handled via sessions (`UserSessions`) and password resets (`PasswordResetTokens`).

The database schema (as provided in the ER diagram) will be used for all data storage and retrieval. Key entities include Users, Agencies, Vehicles, Quotes, Policies, etc. Database interactions are via SQL queries (assuming a relational DB like PostgreSQL), with JSON fields (`AdditionalData`) for extensibility.

The system assumes a RESTful API backend (e.g., built with Node.js/Express or similar) that all portals consume. APIs are secured with JWT tokens derived from `UserSessions.SessionToken`.

## 1. Master Admin Portal

### Purpose

The Master Admin oversees the entire system: creating/managing agencies, insurance providers, high-level users (e.g., AgencyAdmins), and system configurations.

### Pages

- **Dashboard**: Overview of all agencies, total quotes/policies, user stats.
- **Agencies Management**: List, create, edit, delete agencies.
- **Insurance Providers Management**: List, create, edit, delete providers; manage agency access.
- **Users Management**: List all users (filter by role/agency), create/edit/delete users (e.g., assign AgencyAdmins).
- **Reports**: Generate reports on quotes, policies, and user activity.
- **Settings**: System-wide settings (e.g., coverage options, roles).
- **Login/Logout**: Authentication page.
- **Profile**: View/edit own profile.

### Workflow

1. **Login**: MasterAdmin logs in using Email/PasswordHash; session created in UserSessions.
2. **Manage Agencies**: View list from Agencies table; create new agency (insert into Agencies); assign users (update Users.AgencyID).
3. **Manage Providers**: Add providers to InsuranceProviders; grant access via AgencyProviderAccess.
4. **User Oversight**: Search users from Users table; promote/demote roles; monitor sessions via UserSessions.
5. **Reporting**: Query aggregates from Quotes, Policies, Vehicles for dashboards.
6. **Logout**: Invalidate session in UserSessions.

Error handling: If an agency is deleted, cascade updates to Users (set AgencyID to NULL) and AgencyLocations.

### Database Interactions

- **Reads**: SELECT from Users (JOIN Agencies, AgencyLocations), InsuranceProviders (JOIN AgencyProviderAccess), Quotes/Policies for reports.
- **Writes**: INSERT/UPDATE/DELETE on Agencies, InsuranceProviders, AgencyProviderAccess, Users (e.g., Role, AgencyID).
- **Constraints**: Enforce UK on Users.Email; FKs for AgencyID/LocationID.
- **Extensibility**: Use AdditionalData JSON for custom agency/provider fields.

### API Listings

All APIs prefixed with `/api/master-admin/`. Role check: Must be 'MasterAdmin'.

- **GET /agencies**: List all agencies (query Agencies table).
- **POST /agencies**: Create agency (body: {AgencyName, AgentCode, AdditionalData}; insert into Agencies).
- **PUT /agencies/:id**: Update agency.
- **DELETE /agencies/:id**: Delete agency (cascade to related tables).
- **GET /providers**: List providers (query InsuranceProviders).
- **POST /providers**: Create provider (body: {ProviderName, ContactInfo, AdditionalData}).
- **PUT /providers/:id**: Update provider.
- **DELETE /providers/:id**: Delete provider.
- **POST /provider-access**: Grant access (body: {AgencyID, ProviderID}; insert into AgencyProviderAccess).
- **DELETE /provider-access**: Revoke access.
- **GET /users**: List users (query Users with filters: role, agencyID).
- **POST /users**: Create user (body: {FirstName, LastName, Email, PasswordHash, Role, AgencyID, ...}; insert into Users).
- **PUT /users/:id**: Update user (e.g., change Role).
- **DELETE /users/:id**: Delete user.
- **GET /reports/quotes**: Aggregate quotes stats (query Quotes JOIN Users).
- **GET /reports/policies**: Aggregate policies stats.
- **POST /auth/login**: Authenticate and create session (insert UserSessions).
- **POST /auth/logout**: Invalidate session.
- **POST /auth/reset-password**: Generate reset token (insert PasswordResetTokens).
- **PUT /auth/reset-password/:token**: Use token to reset password.

## 2. Agency Admin Portal

### Purpose

AgencyAdmins manage their agency, locations, agents, and LocationAdmins. LocationAdmins (sub-role) focus on location-specific tasks but share the portal with restricted views.

### Pages

- **Dashboard**: Agency stats (e.g., agents, locations, quotes generated).
- **Locations Management**: List, create, edit, delete agency locations.
- **Agents Management**: List, assign, edit agents (Users with Role='Agent').
- **Users Management**: Manage LocationAdmins and Agents within the agency.
- **Provider Access**: View/manage accessible providers (read-only for most, edit if allowed).
- **Reports**: Agency-specific reports on quotes/policies.
- **Login/Logout**: Authentication.
- **Profile**: View/edit own profile and agency details.

### Workflow

1. **Login**: AgencyAdmin/LocationAdmin logs in; session scoped to AgencyID/LocationID.
2. **Manage Locations**: (AgencyAdmin only) Create locations (insert AgencyLocations); assign users (update Users.LocationID).
3. **Manage Agents/Users**: Add agents (insert Users with Role='Agent', AgencyID); view driving/insurance history.
4. **Provider Review**: View granted providers via AgencyProviderAccess; request access (notify MasterAdmin via email or internal system).
5. **Reporting**: Filter queries to own AgencyID.
6. **Logout**: Invalidate session.

For LocationAdmin: Views filtered to their LocationID; cannot manage agencies.

Error handling: Prevent cross-agency access via WHERE AgencyID = current_user.AgencyID.

### Database Interactions

- **Reads**: SELECT from AgencyLocations (JOIN Agencies WHERE AgencyID = current), Users (WHERE AgencyID = current OR LocationID = current), AgencyProviderAccess.
- **Writes**: INSERT/UPDATE on AgencyLocations, Users (limited to own agency; e.g., assign Agent role).
- **Constraints**: FK enforcement on AgencyID/LocationID; Role-based restrictions in queries.
- **Extensibility**: AdditionalData for location-specific custom fields.

### API Listings

All APIs prefixed with `/api/agency-admin/`. Role check: 'AgencyAdmin' or 'LocationAdmin'; scope to AgencyID/LocationID.

- **GET /locations**: List locations for agency (query AgencyLocations WHERE AgencyID = current).
- **POST /locations**: Create location (body: {Address, AdditionalData}; insert with AgencyID).
- **PUT /locations/:id**: Update location.
- **DELETE /locations/:id**: Delete location.
- **GET /agents**: List agents (query Users WHERE Role='Agent' AND AgencyID = current).
- **POST /agents**: Create agent (body: similar to users POST; set AgencyID).
- **PUT /agents/:id**: Update agent.
- **DELETE /agents/:id**: Delete agent.
- **GET /users**: List users in agency/location.
- **POST /users**: Create LocationAdmin/Agent.
- **PUT /users/:id**: Update user in scope.
- **GET /providers**: List accessible providers (query AgencyProviderAccess JOIN InsuranceProviders).
- **GET /reports/quotes**: Agency quotes stats (query Quotes WHERE AgentID IN agency users).
- **GET /reports/policies**: Agency policies stats.
- **POST /auth/login**: Same as above.
- **POST /auth/logout**: Same.
- **POST /auth/reset-password**: Same.

## 3. Agent Portal

### Purpose

Agents interact with customers to create vehicle quotes, manage driving/insurance history, bind policies, and track status.

### Pages

- **Dashboard**: Agent's quotes/policies overview, customer list.
- **Customers Management**: List, create, edit customers (Users with Role='Customer').
- **Vehicles Management**: Add/edit vehicles for customers.
- **Drivers Management**: Assign drivers to vehicles, manage history.
- **Quotes Management**: Create, edit, present quotes; select providers/coverages.
- **Policies Management**: Bind quotes to policies, view active/expired.
- **Reports**: Personal reports on quotes bound.
- **Login/Logout**: Authentication.
- **Profile**: View/edit own profile.

### Workflow

1. **Login**: Agent logs in; session scoped to AgencyID/LocationID.
2. **Manage Customers**: Add customer (insert Users with Role='Customer'); update details like Address, DOB.
3. **Add Vehicles/Drivers**: Create vehicle (insert Vehicles with CustomerID); assign drivers (insert VehicleDrivers); add history (DrivingHistory, InsuranceHistory).
4. **Generate Quote**: Create quote (insert Quotes with AgentID, CustomerID, VehicleID); add line items (QuoteLineItems with ProviderID, CoverageID).
5. **Bind Policy**: Update quote status to 'Bound'; create policy (insert Policies linking to QuoteID).
6. **View History**: Pull driving/insurance records for underwriting.
7. **Logout**: Invalidate session.

Error handling: Only access own agency's providers via AgencyProviderAccess; validate VIN uniqueness.

### Database Interactions

- **Reads**: SELECT from Users (Customers), Vehicles (JOIN Users), VehicleDrivers, DrivingHistory, InsuranceHistory, Quotes (JOIN QuoteLineItems, Coverages, InsuranceProviders WHERE ProviderID IN agency access).
- **Writes**: INSERT/UPDATE on Vehicles, VehicleDrivers, DrivingHistory, InsuranceHistory, Quotes, QuoteLineItems, Policies (set Status).
- **Constraints**: FKs on CustomerID, VehicleID, ProviderID; enforce min_faves etc. via business logic.
- **Extensibility**: AdditionalData for vehicle/coverage customizations.

### API Listings

All APIs prefixed with `/api/agent/`. Role check: 'Agent'; scope to AgencyID.

- **GET /customers**: List customers (query Users WHERE Role='Customer' AND AgencyID = current).
- **POST /customers**: Create customer (body: {FirstName, LastName, Email, ...}; insert with AgencyID).
- **PUT /customers/:id**: Update customer.
- **DELETE /customers/:id**: Delete (soft delete or archive).
- **GET /vehicles**: List vehicles for customer (query Vehicles WHERE CustomerID = :id).
- **POST /vehicles**: Create vehicle (body: {VIN, Make, Model, ...}; insert with CustomerID).
- **PUT /vehicles/:id**: Update vehicle.
- **POST /drivers**: Assign driver to vehicle (insert VehicleDrivers).
- **POST /driving-history**: Add incident (insert DrivingHistory with UserID).
- **POST /insurance-history**: Add history (insert InsuranceHistory).
- **GET /quotes**: List quotes (query Quotes WHERE AgentID = current).
- **POST /quotes**: Create quote (body: {CustomerID, VehicleID}; insert with AgentID).
- **PUT /quotes/:id**: Update quote status.
- **POST /quote-line-items**: Add coverage (insert QuoteLineItems with ProviderID from accessible).
- **GET /coverages**: List available coverages (query Coverages).
- **GET /providers**: List accessible providers.
- **POST /policies**: Bind quote to policy (insert Policies with QuoteID).
- **GET /policies**: List policies.
- **GET /reports**: Agent stats.
- **POST /auth/login**: Same.
- **POST /auth/logout**: Same.
- **POST /auth/reset-password**: Same.

## General System Notes

- **Authentication Across Portals**: Shared APIs for login/logout/reset using UserSessions and PasswordResetTokens. JWT tokens include Role, AgencyID, LocationID for RBAC.
- **Database Workflow**: All portals query/write via the same DB. Use transactions for quote binding (e.g., update Quotes.Status and insert Policies atomically). Indexing on FKs (e.g., UserID, AgencyID) for performance.
- **Security**: API middleware checks Role and scopes (e.g., WHERE AgencyID = req.user.AgencyID). Hash passwords with bcrypt.
- **Extensibility**: JSON fields allow adding fields without schema changes.
- **Frontend**: Assume React/Vue for portals, consuming APIs via fetch/Axios.
- **Deployment**: Single backend, multiple frontend apps or routed via /master, /agency, /agent.
