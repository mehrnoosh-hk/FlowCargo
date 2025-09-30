export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-8">
      <div className="max-w-4xl w-full space-y-8 text-center">
        <h1 className="text-6xl font-bold text-foreground">
          Welcome to FlowCargo
        </h1>
        <p className="text-xl text-foreground/70">
          Shipment and Tenant Management System
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-12">
          <a
            href="/admin"
            className="p-8 border border-foreground/20 rounded-lg hover:border-foreground/40 transition-colors"
          >
            <h2 className="text-2xl font-semibold mb-2">Admin Dashboard →</h2>
            <p className="text-foreground/60">
              Manage tenants, shipments, and system configuration
            </p>
          </a>

          <a
            href="/tracking"
            className="p-8 border border-foreground/20 rounded-lg hover:border-foreground/40 transition-colors"
          >
            <h2 className="text-2xl font-semibold mb-2">Track Shipment →</h2>
            <p className="text-foreground/60">
              Track your shipments in real-time
            </p>
          </a>
        </div>

        <div className="mt-12 p-6 bg-foreground/5 rounded-lg">
          <p className="text-sm text-foreground/60">
            Backend API: <code className="px-2 py-1 bg-foreground/10 rounded">{process.env.NEXT_PUBLIC_API_URL}</code>
          </p>
        </div>
      </div>
    </div>
  );
}
