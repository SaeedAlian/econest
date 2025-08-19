import { Outlet } from "react-router";

import { SidebarProvider } from "@/components/ui/sidebar";
import { Sidebar } from "@/components/dashboard/sidebar";
import { Navbar } from "@/components/dashboard/navbar";
import { Footer } from "@/components/dashboard/footer";

function DashboardLayout() {
  return (
    <>
      <SidebarProvider>
        <Sidebar />
        <div className="flex flex-col items-center w-full min-h-screen">
          <Navbar />
          <div className="w-full mt-2 flex flex-col items-center">
            <Outlet />
          </div>
          <Footer />
        </div>
      </SidebarProvider>
    </>
  );
}

export default DashboardLayout;
