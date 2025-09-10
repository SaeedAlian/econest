import { Routes, Route } from "react-router";

import MainLayout from "@/layouts/main";
import HomeLayout from "@/layouts/home";
import DashboardLayout from "@/layouts/dashboard";

import Home from "@/pages/home";
import Contact from "@/pages/contact";
import AboutUs from "@/pages/aboutus";
import Register from "@/pages/register";
import Login from "@/pages/login";
import ProductList from "@/pages/product-list";
import ProductView from "@/pages/product-view";

import "./App.css";

function App() {
  return (
    <>
      <Routes>
        <Route element={<MainLayout />}>
          <Route element={<HomeLayout />}>
            <Route index element={<Home />} />
            <Route path="about" element={<AboutUs />} />
            <Route path="contact" element={<Contact />} />
          </Route>
        </Route>
        <Route path="product" element={<DashboardLayout />}>
          <Route index element={<ProductList />} />
          <Route path=":productId" element={<ProductView />} />
        </Route>

        <Route path="auth">
          <Route path="login" element={<Login />} />
          <Route path="register" element={<Register />} />
        </Route>
      </Routes>
    </>
  );
}

export default App;
