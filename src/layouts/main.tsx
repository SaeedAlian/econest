import { Outlet } from "react-router";

function MainLayout() {
  return (
    <main className="flex flex-col items-center min-h-screen w-full">
      <Outlet />
    </main>
  );
}

export default MainLayout;
