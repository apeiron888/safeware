# SafeWare - Secure Inventory Management System

<div align="center">

**A high-security, role-based inventory management system with strict access control and comprehensive audit logging.**

[Features](#-features) • [Screenshots](#-screenshots) • [Tech Stack](#-tech-stack) • [API Documentation](#-api-documentation) • [Quick Start](#-quick-start) • [Security](#-security)

</div>

---

## 📸 Screenshots

### Login Page
![Login Page](./screenshots/login.png)

### Manager Dashboard
![Manager Dashboard](./screenshots/manager-dashboard.png)

### Manager - Warehouses
![Manager Warehouses](./screenshots/manager-warehouses.png)

### Manager - Items
![Manager Items](./screenshots/manager-items.png)

### Supervisor Dashboard
![Supervisor Dashboard](./screenshots/supervisor-dashboard.png)

### Staff Dashboard
![Staff Dashboard](./screenshots/staff-dashboard.png)

### Auditor Dashboard
![Auditor Dashboard](./screenshots/auditor-dashboard.png)

### Auditor - Audit Logs
![Auditor Audit Logs](./screenshots/auditor-audit-logs.png)
---

## 🚀 Features

### 🔐 Multi-Level Access Control

SafeWare implements three layers of access control for maximum security:

#### **RBAC (Role-Based Access Control)**
- **Manager**: Complete system access including company setup, warehouse management, employee management, and global oversight
- **Supervisor**: Manages assigned warehouse, supervises staff, and controls warehouse inventory
- **Staff**: Handles inventory items within assigned warehouse
- **Auditor**: Read-only access to all company data for compliance and monitoring

#### **RuBAC (Rule-Based Access Control)**
- **Time Enforcement**: Non-manager employees can only access the system between 8:00 AM - 6:00 PM
- Prevents unauthorized access outside business hours

#### **MAC (Mandatory Access Control)**
- **Warehouse Boundaries**: Staff and Supervisors are strictly restricted to their assigned warehouse
- Cannot view or modify items in other warehouses
- Data isolation ensures security and prevents cross-contamination

### 🛡️ Security & Compliance

- **AES-256 Encryption**: All sensitive audit log details are encrypted at rest using AES-GCM
- **JWT Authentication**: Secure token-based authentication with role and warehouse context
- **Comprehensive Auditing**: Every action (Create, Update, Delete, Login) is logged asynchronously
- **Audit Trail**: Complete history of all system activities for compliance and security monitoring
- **Encrypted Storage**: Sensitive data encrypted before storage in MongoDB

### 📊 Management Capabilities

- **Warehouse Management**: Create, update, and organize multiple warehouse locations
- **Employee Management**: Hire, promote/demote, and manage employees across different roles
- **Inventory Tracking**: Real-time inventory management with status tracking
- **Kanban Boards**: Visual inventory organization with drag-and-drop functionality
- **Audit Logs**: Searchable, filterable audit logs with decryption on-demand

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go (Golang) 1.21+
- **Framework**: Gin Web Framework
- **Database**: MongoDB (with MongoDB Atlas support)
- **Authentication**: JWT (JSON Web Tokens)
- **Encryption**: AES-256-GCM
- **Architecture**: Clean Architecture with separation of concerns

### Frontend
- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v6
- **State Management**: React Context API
- **UI Components**: 
  - Headless UI for accessible components
  - Hero Icons for iconography
  - React DnD for drag-and-drop
  - Formik & Yup for forms and validation
- **Styling**: Tailwind CSS with dark mode support
- **HTTP Client**: Axios with interceptors

### DevOps & Testing
- **Containerization**: Docker & Docker Compose
- **API Testing**: Postman collections included
- **Verification**: Automated shell scripts for testing workflows

---

## 📦 API Documentation

A complete Postman collection is available for testing all API endpoints.

### Files
- **Collection**: `tests/backend/safeware_backend.postman_collection.json`
- **Environment**: `tests/backend/safeware_environment.postman_environment.json`

### API Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register new company (Manager)
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

#### Manager Endpoints
- `GET /api/v1/manager/warehouses` - List all warehouses
- `POST /api/v1/manager/warehouses` - Create warehouse
- `PUT /api/v1/manager/warehouses/:id` - Update warehouse
- `DELETE /api/v1/manager/warehouses/:id` - Delete warehouse
- `GET /api/v1/manager/warehouses/:id` - Get warehouse details
- `GET /api/v1/manager/employees` - List all employees
- `POST /api/v1/manager/employees` - Create employee
- `PUT /api/v1/manager/employees/:id` - Update employee
- `DELETE /api/v1/manager/employees/:id` - Delete employee
- `GET /api/v1/manager/items` - List all items
- `POST /api/v1/manager/items` - Create item
- `PUT /api/v1/manager/items/:id` - Update item
- `DELETE /api/v1/manager/items/:id` - Delete item
- `GET /api/v1/manager/audit-logs` - View audit logs (decrypted)

#### Supervisor Endpoints
- `GET /api/v1/supervisor/employees` - List warehouse staff
- `GET /api/v1/supervisor/items` - List warehouse items
- `POST /api/v1/supervisor/items` - Create item in warehouse
- `PUT /api/v1/supervisor/items/:id` - Update warehouse item
- `DELETE /api/v1/supervisor/items/:id` - Delete warehouse item

#### Staff Endpoints
- `GET /api/v1/staff/items` - List warehouse items
- `PUT /api/v1/staff/items/:id` - Update item status/details

#### Auditor Endpoints
- `GET /api/v1/auditor/warehouses` - View all warehouses (read-only)
- `GET /api/v1/auditor/items` - View all items (read-only)
- `GET /api/v1/auditor/audit-logs` - View audit logs (read-only)

---

## 🚀 Quick Start

### Prerequisites
- **Go**: 1.21 or higher
- **Node.js**: 18.0 or higher
- **MongoDB**: Local instance or MongoDB Atlas account
- **npm** or **yarn** package manager

### 1. Clone the Repository
```bash
git clone <repository-url>
cd safeware
```

### 2. Backend Setup

#### Configure Environment
```bash
cd backend
cp .env.example .env
```

Edit `.env` and configure:
```env
MONGODB_URI=mongodb://localhost:27017/safeware
# OR for MongoDB Atlas:
# MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/safeware

JWT_SECRET=your-super-secret-jwt-key-change-this
ENCRYPTION_KEY=your-32-byte-encryption-key-here
```

#### Install Dependencies & Run
```bash
go mod download
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`

### 3. Frontend Setup

#### Configure Environment
```bash
cd frontend
cp .env.example .env
```

Edit `.env` and configure:
```env
VITE_API_URL=http://localhost:8080/api/v1
```

#### Install Dependencies & Run
```bash
npm install
npm run dev
```

The frontend will start on `http://localhost:5173`

### 4. Access the Application

Open your browser and navigate to `http://localhost:5173`

**First Time Setup:**
1. Click "Register here" on the login page
2. Create a Manager account for your company
3. Login with your Manager credentials
4. Create warehouses, hire employees, and start managing inventory!

### 5. Verify Installation (Optional)

Run the automated verification script to test the complete workflow:
```bash
./verify_vaultory.sh
```

This script will:
- Register a Manager account
- Create a warehouse
- Hire Supervisor and Staff
- Test access control (time & location restrictions)
- Verify audit logging
- Test all CRUD operations

---

## 📂 Project Structure

```
safeware/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # Application entry point
│   ├── internal/
│   │   ├── audit/                   # Audit service (encryption & logging)
│   │   │   ├── audit.go
│   │   │   └── encryption.go
│   │   ├── auth/                    # Authentication & JWT
│   │   │   └── jwt.go
│   │   ├── database/                # Database connection
│   │   │   └── mongodb.go
│   │   ├── handlers/                # API route handlers
│   │   │   ├── auth.go
│   │   │   ├── manager.go
│   │   │   ├── supervisor.go
│   │   │   ├── staff.go
│   │   │   ├── auditor.go
│   │   │   └── item.go
│   │   ├── middleware/              # HTTP middleware
│   │   │   ├── auth.go              # JWT validation
│   │   │   ├── time_restriction.go  # Time-based access control
│   │   │   └── warehouse_access.go  # Warehouse boundary enforcement
│   │   └── models/                  # Data models
│   │       ├── user.go
│   │       ├── warehouse.go
│   │       ├── item.go
│   │       └── audit_log.go
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/              # Reusable components
│   │   │   ├── layout/
│   │   │   │   └── Layout.tsx       # Main layout with sidebar
│   │   │   └── ProtectedRoute.tsx   # Route protection
│   │   ├── context/                 # React Context
│   │   │   └── AuthContext.tsx      # Authentication state
│   │   ├── pages/                   # Page components
│   │   │   ├── auth/
│   │   │   │   ├── Login.tsx
│   │   │   │   └── Register.tsx
│   │   │   ├── manager/
│   │   │   │   ├── Dashboard.tsx
│   │   │   │   ├── Warehouses.tsx
│   │   │   │   ├── WarehouseDetails.tsx
│   │   │   │   ├── Employees.tsx
│   │   │   │   └── AuditLogs.tsx
│   │   │   ├── supervisor/
│   │   │   │   ├── Dashboard.tsx
│   │   │   │   └── Employees.tsx
│   │   │   ├── staff/
│   │   │   │   └── Dashboard.tsx
│   │   │   └── auditor/
│   │   │       ├── Dashboard.tsx
│   │   │       ├── Warehouses.tsx
│   │   │       └── AuditLogs.tsx
│   │   ├── services/                # API services
│   │   │   └── api.ts               # Axios configuration
│   │   ├── types/                   # TypeScript types
│   │   │   ├── auth.ts
│   │   │   ├── warehouse.ts
│   │   │   └── item.ts
│   │   ├── App.tsx                  # Root component
│   │   └── main.tsx                 # Application entry
│   ├── package.json
│   └── vite.config.ts
├── tests/
│   └── backend/
│       ├── safeware_backend.postman_collection.json
│       └── safeware_environment.postman_environment.json
├── docker-compose.yml               # Docker configuration
├── verify_vaultory.sh              # Automated verification script
└── README.md                        # This file
```

---

## 🔒 Security Notes

### Time Restrictions
- **Manager** accounts can access the system 24/7
- **Supervisor, Staff, and Auditor** accounts are restricted to 8:00 AM - 6:00 PM
- Attempts to access outside allowed hours will return `403 Forbidden`

### Warehouse Isolation
- Staff and Supervisors are bound to their assigned warehouse
- They cannot view or modify data from other warehouses
- Warehouse ID is embedded in JWT tokens for validation

### Audit Logs
- All logs are encrypted using AES-256-GCM encryption
- Decryption happens on-the-fly when viewing through Manager/Auditor endpoints
- Logs capture: timestamp, user, action, resource, and detailed changes

### Authentication
- Passwords are hashed using bcrypt before storage
- JWT tokens include role and warehouse context
- Access tokens expire after 15 minutes
- Refresh tokens valid for 7 days
- Token refresh mechanism prevents session interruption

---

## 🤝 Contributing

This project follows enterprise security standards. When contributing:

1. Never commit sensitive data (`.env` files, credentials)
2. Follow the existing code structure and patterns
3. Add tests for new features
4. Update API documentation for endpoint changes
5. Ensure all security middleware is properly applied

---

## 📝 License

This project is proprietary software developed for secure inventory management.

---

## 👥 Support

For questions or issues:
- Review the API documentation in `tests/backend/`
- Check the verification script `verify_vaultory.sh` for example workflows
- Consult the Postman collection for API usage examples

---

<div align="center">

**Purpose-built for secure inventory management**

Go • React • TypeScript • MongoDB

</div>
