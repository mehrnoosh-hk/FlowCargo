import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "FlowCargo",
  description: "FlowCargo - Shipment and Tenant Management System",
};

// app/layout.tsx
export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="bg-neutral-800 text-purple-200">
        {children}
      </body>
    </html>
  )
}
