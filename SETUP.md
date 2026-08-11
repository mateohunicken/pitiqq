# Finanzas Domésticas - Guía de Setup y Deployment

## 📋 Descripción

Aplicación web robusta para gestionar finanzas domésticas en familia. Permite registrar ingresos, egresos, compra/venta de dólares, impuestos y servicios.

**Stack Tecnológico:**
- **Backend:** Go + Gin (robusto, performante, sin dependencias de runtime)
- **Frontend:** React + TypeScript (seguro, tipado)
- **Base de datos:** PostgreSQL (estándar de oro para finanzas)
- **Hosting:** Railway (serverless-like, sin mantenimiento)

---

## 🚀 Setup Local (Desarrollo)

### Requisitos Previos
- **Go 1.21+** - https://golang.org/dl
- **Node.js 18+** - https://nodejs.org
- **PostgreSQL 14+** - https://www.postgresql.org/download
- **Git**

### 1. Clonar o Descargar el Proyecto

```bash
cd finanzas-app
```

### 2. Setup de Base de Datos

```bash
# Crear base de datos PostgreSQL
createdb finanzas_app

# O en psql:
psql
CREATE DATABASE finanzas_app;
\q
```

### 3. Setup del Backend

```bash
cd backend

# Copiar archivo de configuración
cp .env.example .env

# Editar .env con tus credenciales de PostgreSQL
nano .env
# O en Windows:
# notepad .env
```

**Configuración en `.env`:**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_contraseña
DB_NAME=finanzas_app
JWT_SECRET=una-clave-super-secreta-muy-larga-minimo-32-caracteres
PORT=8080
```

```bash
# Descargar dependencias
go mod download

# Ejecutar servidor (creará tablas automáticamente)
go run main.go
```

**Debería ver:**
```
✓ Connected to PostgreSQL
✓ Database schema initialized
🚀 Finanzas App Backend running on http://localhost:8080
```

### 4. Setup del Frontend

```bash
cd ../frontend

# Instalar dependencias
npm install

# Crear archivo de configuración
cp .env.example .env.local
# Variables por defecto ya están configuradas para localhost

# Ejecutar servidor dev
npm run dev
```

**Acceso:** http://localhost:5173

---

## 🔑 Crear Primer Usuario

La app funciona con multi-usuario desde el inicio. Al registrar:

1. Ve a http://localhost:5173/register
2. Completa:
   - Email: `tú@email.com`
   - Nombre: `Tu Nombre`
   - Contraseña: Mínimo 8 caracteres
   - Nombre del Household: `Finanzas Familiares`
3. Click en "Registrarse"

Tu mujer puede luego invitarse a través de funcionalidad de invitaciones (endpoint en desarrollo).

---

## 📦 Deploy a Railway (Sin Mantenimiento)

Railway es una plataforma de hosting que requiere setup inicial pero después funciona sin intervención.

### Paso 1: Crear Cuenta en Railway
1. Ve a https://railway.app
2. Crea cuenta con GitHub (recomendado)

### Paso 2: Conectar Repositorio

```bash
# Crear repositorio git (si aún no existe)
git init
git add .
git commit -m "Initial commit"

# Crear repositorio en GitHub
# 1. Ve a https://github.com/new
# 2. Crea repo "finanzas-app"
# 3. Sigue instrucciones para push

git remote add origin https://github.com/tu-usuario/finanzas-app.git
git push -u origin main
```

### Paso 3: Configurar en Railway

1. Entra a https://railway.app/dashboard
2. Click en "New Project"
3. Selecciona "Deploy from GitHub"
4. Autoriza y selecciona tu repo `finanzas-app`
5. Click "Deploy"

### Paso 4: Agregar PostgreSQL

1. En Railway Dashboard, click "Add"
2. Selecciona "Database → PostgreSQL"
3. Se conecta automáticamente a tu app

### Paso 5: Configurar Variables de Entorno

En Railway Dashboard → Tu App → Variables:

```
DB_HOST=postgres_host_de_railway
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=auto_generada_por_railway
DB_NAME=finanzas_app
DB_SSLMODE=require
JWT_SECRET=genera-una-clave-aleatoria-muy-larga
ENV=production
PORT=8080
```

Railway obtiene automáticamente `DB_HOST`, `DB_USER`, `DB_PASSWORD` si añadiste PostgreSQL.

### Paso 6: Desplegar

Railway despliega automáticamente cuando haces push a `main`:

```bash
git push origin main
```

Chequea en Railway Dashboard → Deployments para seguimiento.

---

## 🌐 Acceso Después del Deploy

Después del deploy en Railway obtendrás URLs como:

- **Frontend:** `https://finanzas-app-production-xxx.railway.app`
- **Backend:** `https://finanzas-api-production-xxx.railway.app`

Actualiza la configuración de CORS en `.env` del backend:

```env
CORS_ORIGINS=https://finanzas-app-production-xxx.railway.app,https://tu-dominio-personalizado.com
```

---

## 🛡️ Seguridad

### Antes de ir a Producción:

1. **JWT Secret**: Genera una clave aleatoria fuerte (mínimo 32 caracteres)
   ```bash
   openssl rand -base64 32
   ```

2. **HTTPS**: Railway lo proporciona automáticamente

3. **Contraseñas**: Se hashean con bcrypt (ver `handlers/auth.go`)

4. **CORS**: Configura solo los orígenes permitidos

5. **Base de Datos**: PostgreSQL en Railway tiene SSL incluido

---

## 📊 Estructura de Carpetas

```
finanzas-app/
├── DATABASE_SCHEMA.sql     # Esquema PostgreSQL
├── backend/
│   ├── main.go             # Entrada principal
│   ├── go.mod              # Dependencias Go
│   ├── .env.example        # Configuración ejemplo
│   ├── config/config.go    # Cargar configuración
│   ├── models/models.go    # Definición de modelos
│   ├── database/database.go # Conexión BD
│   ├── handlers/           # Endpoints API
│   └── middleware/         # Autenticación, CORS
└── frontend/
    ├── package.json        # Dependencias Node
    ├── src/
    │   ├── types/          # Tipos TypeScript
    │   ├── services/       # Cliente API
    │   ├── stores/         # Estado Zustand
    │   ├── pages/          # Páginas React
    │   ├── components/     # Componentes reutilizables
    │   └── App.tsx         # Componente raíz
    └── index.html          # HTML base
```

---

## 🔧 Endpoints Principales

### Autenticación
- `POST /api/v1/auth/register` - Registrar usuario
- `POST /api/v1/auth/login` - Iniciar sesión

### Transacciones
- `GET /api/v1/households/:householdId/transactions` - Listar
- `POST /api/v1/households/:householdId/transactions` - Crear
- `PUT /api/v1/households/:householdId/transactions/:id` - Actualizar
- `DELETE /api/v1/households/:householdId/transactions/:id` - Eliminar

### Divisas
- `GET /api/v1/households/:householdId/currency-exchanges` - Listar cambios
- `POST /api/v1/households/:householdId/currency-exchanges` - Registrar cambio

### Reportes
- `GET /api/v1/households/:householdId/reports/summary` - Resumen financiero
- `GET /api/v1/households/:householdId/reports/monthly-breakdown` - Desglose mensual
- `GET /api/v1/households/:householdId/reports/currency-balance` - Balance en divisas

---

## 🐛 Troubleshooting

### "Cannot connect to database"
```bash
# Verificar que PostgreSQL está corriendo
psql -U postgres -d finanzas_app
# Si no puedes conectar, reinicia PostgreSQL
```

### "JWT token is invalid"
- Asegúrate de que `JWT_SECRET` es el mismo en backend y que esté guardado en `.env`

### "CORS Error"
- Backend y Frontend deben estar en mismas URLs o configurar CORS_ORIGINS correctamente

### Frontend no se conecta al backend
- Verifica que backend está en `http://localhost:8080`
- Checkea proxy en `vite.config.ts`

---

## 🚦 Próximos Pasos (Post-MVP)

1. **Invitación de usuarios**: Permitir que otros se unan al household
2. **Exportar a PDF**: Generar reportes para imprimir
3. **Notificaciones**: Alertas para pagos de servicios
4. **Sincronización de datos**: Cálculo automático de saldos
5. **Historial detallado**: Auditoría de cambios
6. **Multimoneda avanzada**: Conversión automática con tipos de cambio reales

---

## 📞 Soporte

Cualquier problema o duda, contacta al desarrollador.

**¡Aplicación lista para usar!** 🎉
