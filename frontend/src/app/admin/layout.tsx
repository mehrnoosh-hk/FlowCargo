import { AdminSidebar } from "@/components/admin-sidebar";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Admin Dashboard",
    description: "Admin. dashboard page",
};

const defaultOpen: boolean = false

export default function AdminDashboard({ children }: { children: React.ReactNode }) {

    return (
        <SidebarProvider defaultOpen={defaultOpen}>
            <AdminSidebar />
            <main>
                {children}
            </main>
        </SidebarProvider>
    );
}
