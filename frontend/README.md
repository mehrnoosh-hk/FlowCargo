# FlowCargo Frontend

Next.js-based frontend application for FlowCargo, featuring both user-facing pages and an admin dashboard.

## Tech Stack

- **Framework**: Next.js 15+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui
- **HTTP Client**: Axios (to be added)
- **Forms**: React Hook Form + Zod (to be added)

## Project Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── (user)/            # User-facing pages route group
│   │   ├── admin/             # Admin dashboard routes
│   │   ├── layout.tsx         # Root layout
│   │   └── globals.css        # Global styles
│   ├── components/
│   │   └── ui/                # shadcn/ui components
│   ├── lib/                   # Utilities and API client
│   ├── services/              # API service layer
│   ├── types/                 # TypeScript type definitions
│   └── hooks/                 # Custom React hooks
├── public/                    # Static assets
└── Configuration files
```

## Prerequisites

- Node.js 18.17+ or 20+
- npm 9+
- Go backend running on http://localhost:8080

## Getting Started

### 1. Install Dependencies

```bash
npm install
```

### 2. Environment Variables

Copy `.env.example` to `.env.local` and configure:

```bash
cp .env.example .env.local
```

Edit `.env.local` as needed:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=FlowCargo
```

### 3. Run Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Available Scripts

- `npm run dev` - Start development server (port 3000)
- `npm run build` - Build production bundle
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run type-check` - Run TypeScript type checking

## Development Workflow

### Running Full Stack

**Terminal 1 - Backend:**
```bash
cd FlowCargo
make run
```

**Terminal 2 - Frontend:**
```bash
cd FlowCargo/frontend
npm run dev
```

### Using Makefile (from FlowCargo root)

```bash
# Install frontend dependencies
make frontend-install

# Run frontend dev server
make frontend-dev

# Build frontend
make frontend-build
```

## Routes

### User-Facing Routes
- `/` - Home page
- `/about` - About page
- `/tracking` - Shipment tracking
- (More to be added)

### Admin Routes
- `/admin` - Admin dashboard
- `/admin/tenants` - Tenant management
- `/admin/shipments` - Shipment management
- (More to be added)

## API Integration

The frontend communicates with the Go backend via REST API:
- **Backend URL**: http://localhost:8080
- **API Client**: Located in `src/lib/api.ts`
- **Services**: API calls organized in `src/services/`

### Example API Call

```typescript
import api from '@/lib/api';

const response = await api.get('/api/tenants');
```

## Adding shadcn/ui Components

```bash
npx shadcn@latest add button
npx shadcn@latest add card
npx shadcn@latest add input
```

Components will be added to `src/components/ui/`.

## Styling

This project uses Tailwind CSS for styling:

- Global styles: `src/app/globals.css`
- Tailwind config: `tailwind.config.ts`
- Component styles: Inline with Tailwind utility classes

## TypeScript

- Strict mode enabled
- Path aliases configured: `@/*` maps to `src/*`
- Type definitions: `src/types/`

## Future Enhancements

- [ ] Authentication (NextAuth.js or JWT)
- [ ] Form validation (React Hook Form + Zod)
- [ ] State management (React Context/Zustand)
- [ ] Data fetching (React Query/SWR)
- [ ] Testing (Jest + React Testing Library)
- [ ] E2E testing (Playwright/Cypress)

## Deployment

### Build for Production

```bash
npm run build
npm run start
```

### Environment Variables for Production

Set the following in your production environment:
- `NEXT_PUBLIC_API_URL` - Production API URL
- `NODE_ENV=production`

## Troubleshooting

### Port 3000 already in use
```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9

# Or use a different port
PORT=3001 npm run dev
```

### CORS Issues
Ensure the Go backend has CORS middleware configured to allow `http://localhost:3000`.

### Module not found errors
```bash
# Clear Next.js cache and reinstall
rm -rf .next node_modules
npm install
```

## Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [React Documentation](https://react.dev)
- [Tailwind CSS](https://tailwindcss.com)
- [shadcn/ui](https://ui.shadcn.com)
- [TypeScript](https://www.typescriptlang.org)

## License

[Your License Here]