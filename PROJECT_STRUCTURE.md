# 📁 Estructura del Proyecto Finanzas Domésticas

```
finanzas-app/
│
├── 📄 README.md                          # Descripción general del proyecto
├── 📄 SETUP.md                           # Guía de instalación y deployment ⭐
├── 📄 PROJECT_STRUCTURE.md               # Este archivo
├── 📄 .gitignore                         # Archivos a ignorar en git
│
├── 📄 DATABASE_SCHEMA.sql                # ⭐ Schema PostgreSQL completo
│                                          # - Users, Households, Categories
│                                          # - Transactions, Currency Exchanges
│                                          # - Recurring Expenses, Services
│                                          # - Budgets, Audit Logs
│
├── 📄 docker-compose.yml                 # ⭐ Orquestación local 3 servicios
│                                          # - PostgreSQL
│                                          # - Backend (Go)
│                                          # - Frontend (React + Nginx)
│
│
├── 📁 backend/                           # Go + Gin API
│   ├── 📄 go.mod                         # Módulos Go
│   ├── 📄 go.sum                         # Checksums de dependencias
│   ├── 📄 main.go                        # ⭐ Punto de entrada, setup routes
│   ├── 📄 Dockerfile                     # ⭐ Build multistage para producción
│   ├── 📄 .env.example                   # ⭐ Variables de entorno ejemplo
│   │
│   ├── 📁 config/
│   │   └── 📄 config.go                  # Cargar configuración de .env
│   │
│   ├── 📁 models/
│   │   └── 📄 models.go                  # ⭐ Todas las estructuras de datos
│   │                                      # - User, Household, Category
│   │                                      # - Transaction, CurrencyExchange
│   │                                      # - RecurringExpense, Reports
│   │                                      # - DTOs para API
│   │
│   ├── 📁 database/
│   │   └── 📄 database.go                # ⭐ Conexión PostgreSQL
│   │                                      # - Connect, InitializeSchema
│   │                                      # - CreateDefaultCategories
│   │
│   ├── 📁 middleware/
│   │   └── 📄 middleware.go              # ⭐ CORS, JWT Auth, Error Handling
│   │                                      # - AuthMiddleware
│   │                                      # - CORSMiddleware
│   │
│   ├── 📁 handlers/                      # ⭐ Endpoints REST
│   │   ├── 📄 auth.go                    # POST /register, POST /login
│   │   ├── 📄 transactions.go            # CRUD transacciones
│   │   └── 📄 handlers.go                # Otros handlers (Currency, Service, etc)
│   │                                      # - GetCategories, CreateCategory
│   │                                      # - GetExchanges, CreateExchange
│   │                                      # - GetRecurringExpenses
│   │                                      # - Reportes (Summary, Monthly, Currency)
│   │
│   └── 📁 logs/                          # Logs de aplicación (generados)
│
│
├── 📁 frontend/                          # React 18 + TypeScript + Vite
│   ├── 📄 package.json                   # ⭐ Dependencias Node
│   ├── 📄 tsconfig.json                  # Configuración TypeScript
│   ├── 📄 tsconfig.node.json             # TS para vite.config.ts
│   ├── 📄 vite.config.ts                 # ⭐ Configuración Vite (build, proxy)
│   ├── 📄 index.html                     # ⭐ Entrada HTML
│   ├── 📄 Dockerfile                     # ⭐ Build + Nginx production
│   ├── 📄 nginx.conf                     # ⭐ Configuración Nginx
│   ├── 📄 .env.example                   # Variables de entorno (dev)
│   │
│   ├── 📁 public/                        # Assets estáticos
│   │   └── vite.svg
│   │
│   └── 📁 src/
│       ├── 📄 App.tsx                    # ⭐ Componente raíz
│       ├── 📄 main.tsx                   # ⭐ Punto de entrada React
│       ├── 📄 index.css                  # Estilos globales
│       │
│       ├── 📁 types/
│       │   └── 📄 index.ts               # ⭐ Interfaces TypeScript completas
│       │                                  # - User, Household, Category
│       │                                  # - Transaction, CurrencyExchange
│       │                                  # - RecurringExpense, Reports
│       │                                  # - DTOs
│       │
│       ├── 📁 services/
│       │   └── 📄 api.ts                 # ⭐ Cliente Axios tipado
│       │                                  # - Login, Register
│       │                                  # - CRUD Transactions
│       │                                  # - CRUD Currency
│       │                                  # - Reportes
│       │
│       ├── 📁 stores/
│       │   └── 📄 auth.ts                # ⭐ Estado Zustand
│       │                                  # - User, Household, Token
│       │                                  # - Login, Register, Logout
│       │
│       ├── 📁 pages/                     # ⭐ Páginas principales
│       │   ├── 📄 LoginPage.tsx          # (a completar)
│       │   ├── 📄 DashboardPage.tsx      # (a completar)
│       │   ├── 📄 TransactionsPage.tsx   # (a completar)
│       │   ├── 📄 ReportsPage.tsx        # (a completar)
│       │   └── 📄 CurrencyPage.tsx       # (a completar)
│       │
│       ├── 📁 components/                # ⭐ Componentes reutilizables
│       │   ├── 📄 TransactionForm.tsx    # (a completar)
│       │   ├── 📄 TransactionList.tsx    # (a completar)
│       │   ├── 📄 IncomeExpenseChart.tsx # (a completar)
│       │   ├── 📄 CurrencyExchangeForm.tsx # (a completar)
│       │   └── 📄 Navigation.tsx         # (a completar)
│       │
│       ├── 📁 hooks/                     # Custom React hooks (si aplica)
│       │
│       └── 📁 utils/                     # Funciones auxiliares
│           └── 📄 formatters.ts          # Format moneda, fechas, etc.

```

---

## 📝 Notas Importantes

### Backend (Go)
- **Entrada:** `backend/main.go` - Aquí se configura todo
- **Modelos:** `backend/models/models.go` - Todas las estructuras
- **Database:** `backend/database/database.go` - Conexión y schema
- **Handlers:** `backend/handlers/*.go` - Lógica de endpoints
- **Middleware:** `backend/middleware/middleware.go` - Auth, CORS

### Frontend (React)
- **Entrada:** `frontend/src/main.tsx` - Bootstrap de React
- **Tipos:** `frontend/src/types/index.ts` - TypeScript interfaces
- **API:** `frontend/src/services/api.ts` - Cliente Axios
- **Estado:** `frontend/src/stores/auth.ts` - Zustand store

### Base de Datos
- **Schema:** `DATABASE_SCHEMA.sql` - Se ejecuta automáticamente al iniciar
- Tablas: Users, Households, Categories, Transactions, CurrencyExchanges, RecurringExpenses, AuditLogs

---

## 🚀 Archivos Clave (Prioritarios para Editar)

1. **Para agregar nuevos endpoints:** `backend/handlers/` → crear archivo, agregar en `main.go`
2. **Para modificar base de datos:** `DATABASE_SCHEMA.sql` → reiniciar con `docker-compose down -v`
3. **Para nuevas páginas React:** `frontend/src/pages/` → crear `.tsx`, agregar ruta en `App.tsx`
4. **Para llamadas API:** `frontend/src/services/api.ts` → agregar método, usar en página

---

## 📦 Stack Final Resumido

| Componente | Tecnología | Archivo |
|-----------|-----------|---------|
| **Servidor API** | Go 1.21 + Gin | `backend/main.go` |
| **Base de Datos** | PostgreSQL 16 | `DATABASE_SCHEMA.sql` |
| **Frontend** | React 18 + TypeScript | `frontend/src/main.tsx` |
| **Build Tool** | Vite 5 | `frontend/vite.config.ts` |
| **Estado** | Zustand | `frontend/src/stores/auth.ts` |
| **HTTP Client** | Axios | `frontend/src/services/api.ts` |
| **Container** | Docker + Compose | `docker-compose.yml` |
| **Web Server** | Nginx | `frontend/nginx.conf` |

---

## ✅ Checklist Deploy

- [ ] Editar `.env` con credenciales PostgreSQL
- [ ] Ejecutar `docker-compose up -d` o setup manual
- [ ] Crear primer usuario en `/register`
- [ ] Agregar primeras transacciones
- [ ] Verificar reportes funcionan
- [ ] Crear repo GitHub y push
- [ ] Conectar a Railway y agregar PostgreSQL
- [ ] Configurar variables de entorno en Railway
- [ ] Testar URLs de producción

---

¡Proyecto listo para desarrollar y desplegar! 🚀
