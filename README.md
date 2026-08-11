# 💰 Finanzas Domésticas - App Web Robusta

Aplicación moderna y segura para gestionar finanzas domésticas en familia. Registra ingresos, egresos, compra/venta de dólares, impuestos, servicios y genera reportes.

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Status](https://img.shields.io/badge/status-active-green)

---

## ✨ Características

✅ **Multi-usuario** - Vos y tu mujer acceden con credenciales propias  
✅ **Transacciones** - Ingresos, egresos, transferencias  
✅ **Divisas** - Registra compra/venta de dólares  
✅ **Servicios e Impuestos** - Seguimiento de pagos recurrentes  
✅ **Reportes** - Dashboards y gráficos financieros  
✅ **Seguridad** - Autenticación JWT, datos encriptados  
✅ **Sin Mantenimiento** - Deploy a Railway, actualizaciones automáticas  

---

## 🛠️ Stack Tecnológico

| Capa | Tecnología | Por qué |
|------|-----------|--------|
| Backend | **Go + Gin** | Robustez extrema, sin runtime pesado, concurrencia nativa |
| Frontend | **React + TypeScript** | Type-safe, componentes reutilizables, moderno |
| Database | **PostgreSQL** | Estándar de oro para datos financieros |
| Hosting | **Railway** | Serverless, auto-scaling, sin DevOps |

---

## 🚀 Inicio Rápido

### Opción 1: Local (Desarrollo)

```bash
# Clonar repo
git clone https://github.com/tu-usuario/finanzas-app.git
cd finanzas-app

# Con Docker Compose (recomendado)
docker-compose up -d

# O setup manual (ver SETUP.md)
```

**Acceso:**
- Frontend: http://localhost:3000 o http://localhost
- Backend: http://localhost:8080
- Base de datos: localhost:5432

### Opción 2: Railway (Producción)

```bash
# Crear repo en GitHub
git push origin main

# En Railway:
# 1. Conectar repo GitHub
# 2. Agregar PostgreSQL database
# 3. Deploy automático en cada push
```

**URLs de producción** (obtenidas de Railway):
- Frontend: `https://tu-app.railway.app`
- Backend: `https://tu-api.railway.app`

Ver [SETUP.md](./SETUP.md) para instrucciones detalladas.

---

## 📋 Primeros Pasos Después de Deploy

1. **Crear Cuenta**
   - Ve a `/register`
   - Email, nombre, contraseña (mín 8 caracteres)
   - Nombre del "Household" (grupo familiar)

2. **Añadir Transacciones**
   - Dashboard → "Nueva Transacción"
   - Selecciona categoría, monto, fecha
   - Tipo: Ingreso/Egreso/Transferencia

3. **Registrar Dólares**
   - Menú → "Divisas"
   - Compra/Venta con tipo de cambio

4. **Ver Reportes**
   - Menú → "Reportes"
   - Resumen, desglose mensual, balance en divisas

---

## 🔐 Seguridad

- **Autenticación:** JWT (7 días de expiración)
- **Contraseñas:** Hasheadas con bcrypt
- **Conexión:** HTTPS en producción (Railway)
- **Base de datos:** Encriptada en tránsito
- **CORS:** Restringido a orígenes permitidos

---

## 📂 Estructura

```
finanzas-app/
├── DATABASE_SCHEMA.sql          # Definición completa de tablas
├── SETUP.md                     # Guía de instalación y deployment
├── docker-compose.yml           # Orquestación local
│
├── backend/                     # Go + Gin API
│   ├── main.go                  # Punto de entrada
│   ├── config/                  # Configuración
│   ├── models/                  # Estructuras de datos
│   ├── handlers/                # Endpoints REST
│   ├── middleware/              # Auth, CORS, logging
│   └── database/                # Conexión PostgreSQL
│
└── frontend/                    # React + TypeScript
    ├── src/
    │   ├── pages/               # Login, Dashboard, Reportes
    │   ├── components/          # Forms, Charts, Cards
    │   ├── stores/              # Zustand state management
    │   ├── services/            # Cliente API
    │   └── types/               # TypeScript interfaces
    └── Dockerfile               # Build para producción
```

---

## 📊 API Endpoints Principales

### Auth
```http
POST   /api/v1/auth/register
POST   /api/v1/auth/login
```

### Transacciones
```http
GET    /api/v1/households/:id/transactions
POST   /api/v1/households/:id/transactions
PUT    /api/v1/households/:id/transactions/:id
DELETE /api/v1/households/:id/transactions/:id
```

### Reportes
```http
GET    /api/v1/households/:id/reports/summary
GET    /api/v1/households/:id/reports/monthly-breakdown
GET    /api/v1/households/:id/reports/currency-balance
```

Ver [SETUP.md](./SETUP.md) para documentación completa.

---

## 🐛 Troubleshooting

**"Cannot connect to backend"**
```bash
# Verifica que backend esté corriendo
curl http://localhost:8080/health
```

**"Database error"**
```bash
# Resetea PostgreSQL
docker-compose down -v
docker-compose up -d
```

**"JWT Invalid"**
- Verifica que `JWT_SECRET` en `.env` sea consistente
- Limpia localStorage del navegador

Ver [SETUP.md](./SETUP.md#troubleshooting) para más soluciones.

---

## 🚀 Próximas Características

- [ ] Invitación de usuarios por email
- [ ] Exportar reportes a PDF
- [ ] Categorías personalizadas por usuario
- [ ] Sincronización de tipos de cambio reales
- [ ] Historial de cambios (auditoría)
- [ ] Notificaciones de pagos vencidos

---

## 📝 Licencia

Proyecto personal. Úsalo libremente.

---

## 💬 Preguntas o Issues

Contacta al desarrollador o abre un issue en GitHub.

---

**¡Listo para usar!** 🎉

Para comenzar:
1. Lee [SETUP.md](./SETUP.md)
2. Elige: Local con Docker o Railway
3. Crea tu primer usuario
4. ¡Comienza a registrar finanzas!
