# FlowCargo Frontend - Quick Start Guide

**Last Updated**: September 30, 2024  
**Status**: Phase 1 Complete ✅

---

## 🚀 Quick Start

### Prerequisites
- Node.js v20+ installed
- npm v9+ installed
- Go backend running (optional for now)

### Get Started in 3 Steps

```bash
# 1. Navigate to frontend directory
cd FlowCargo/frontend

# 2. Install dependencies (if not already done)
npm install

# 3. Start development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

---

## 📁 Project Structure

```
frontend/
├── src/
│   ├── app/              # Next.js pages (App Router)
│   │   ├── layout.tsx    # Root layout
│   │   ├── page.tsx      # Home page
│   │   └── globals.css   # Global styles
│   ├── components/
│   │   └── ui/          # shadcn/ui components (to be added)
│   ├── lib/             # Utilities & API client (to be added)
│   ├── services/        # API services (to be added)
│   ├── types/           # TypeScript types (to be added)
│   └── hooks/           # Custom React hooks (to be added)
├── public/              # Static assets
└── Configuration files
```

---

## 🛠️ Available Commands

```bash
# Development
npm run dev          # Start dev server (http://localhost:3000)
npm run build        # Build for production
npm run start        # Start production server

# Code Quality
npm run lint         # Run ESLint
npm run type-check   # TypeScript type checking
```

---

## ⚙️ Configuration

### Environment Variables

Located in `.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=FlowCargo
NEXT_PUBLIC_APP_VERSION=0.1.0
NODE_ENV=development
```

### Key Technologies

- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS 3.4.1
- **Components**: shadcn/ui (to be installed in Phase 2)

---

## 🔄 Development Workflow

### Running Full Stack

**Terminal 1 - Backend:**
```bash
cd FlowCargo
make run
# Backend runs on http://localhost:8080
```

**Terminal 2 - Frontend:**
```bash
cd FlowCargo/frontend
npm run dev
# Frontend runs on http://localhost:3000
```

---

## 📦 Current Features

### ✅ What's Working
- Next.js 15 with App Router
- TypeScript with strict mode
- Tailwind CSS configured
- Basic home page
- Development and production builds
- ESLint and type checking

### ⏳ Coming Next (Phase 2)
- shadcn/ui component library
- API client for Go backend
- Admin dashboard pages
- User-facing pages
- CORS configuration

---

## 🐛 Troubleshooting

### Port 3000 Already in Use
```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9

# Or use different port
PORT=3001 npm run dev
```

### Clear Cache and Reinstall
```bash
rm -rf .next node_modules package-lock.json
npm install
```

### Check Build Errors
```bash
npm run build
# Look for specific error messages
```

---

## 📚 Important Files

| File | Purpose |
|------|---------|
| `src/app/layout.tsx` | Root layout wrapper |
| `src/app/page.tsx` | Home page |
| `src/app/globals.css` | Global styles + Tailwind |
| `tailwind.config.ts` | Tailwind configuration |
| `tsconfig.json` | TypeScript configuration |
| `next.config.js` | Next.js configuration |
| `.env.local` | Local environment variables |

---

## 🎯 Path Aliases

TypeScript is configured with path aliases:

```typescript
// Instead of: import Button from '../../components/ui/button'
import Button from '@/components/ui/button'

// @/* maps to src/*
```

---

## 📖 Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [shadcn/ui Components](https://ui.shadcn.com)

---

## ✅ Verification Checklist

Before proceeding to Phase 2:

- [ ] `npm install` runs without errors
- [ ] `npm run dev` starts successfully
- [ ] `npm run build` completes successfully
- [ ] `npm run lint` passes
- [ ] `npm run type-check` passes
- [ ] Home page loads at http://localhost:3000
- [ ] No console errors in browser

---

## 🚦 Status

**Current Phase**: Phase 1 Complete  
**Next Phase**: Phase 2 - shadcn/ui & API Integration  
**Ready to Proceed**: ✅ Yes

See `CHECKPOINT_01_FRONTEND_STRUCTURE.md` for detailed information.