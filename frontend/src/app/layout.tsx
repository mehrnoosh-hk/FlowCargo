import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "FlowCargo",
  description: "FlowCargo - Shipment and Tenant Management System",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
