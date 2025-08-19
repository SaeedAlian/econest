import { Routes, Route } from "react-router";

import MainLayout from "@/layouts/main";
import HomeLayout from "@/layouts/home";

import Home from "@/pages/home";
import Contact from "@/pages/contact";
import AboutUs from "@/pages/aboutus";

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

      </Routes>
    </>
  );
}

export default App;
