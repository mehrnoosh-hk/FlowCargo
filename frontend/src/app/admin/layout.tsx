import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Admin Dashboard",
    description: "Admin. dashboard page",
};

export default function AdminDashboard({ children }: { children: React.ReactNode }) {

    return (
        {children}
    );
}
