import { Outlet } from "react-router";
import { Navbar, NavbarLink } from "@/components/home/navbar";
import { Footer, FooterLink } from "@/components/home/footer";

const navbarLinks: NavbarLink[] = [
  {
    text: "Contact",
    href: "/contact",
  },
  {
    text: "About Us",
    href: "/about",
  },
  {
    text: "Login",
    href: "/auth/login",
  },
  {
    text: "Register",
    href: "/auth/register",
    bordered: true,
  },
];

const footerLinks: FooterLink[] = [
  {
    text: "Home",
    href: "/",
  },
  {
    text: "Contact",
    href: "/contact",
  },
  {
    text: "About Us",
    href: "/about",
  },
  {
    text: "Login",
    href: "/auth/login",
  },
];

function HomeLayout() {
  return (
    <>
      <Navbar links={navbarLinks} />
      <Outlet />
      <Footer links={footerLinks} />
    </>
  );
}

export default HomeLayout;
