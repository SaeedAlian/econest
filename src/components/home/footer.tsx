import { Link } from "react-router";
import { FaInstagram } from "react-icons/fa6";
import { FaThreads } from "react-icons/fa6";
import { BsTwitterX } from "react-icons/bs";

import { Button } from "@/components/ui/button";
import logo from "@/assets/logo.svg";

export type FooterLink = { text: string; href: string; target?: string };

export type FooterProps = {
  instagramLink?: string;
  threadsLink?: string;
  twitterLink?: string;
  links: FooterLink[];
};

function Footer({
  links,
  instagramLink,
  threadsLink,
  twitterLink,
}: FooterProps) {
  return (
    <footer className="mt-auto w-full bg-foreground flex flex-col items-center">
      <img src={logo} alt="EcoNest" className="w-24 h-24 m-[-10px]" />
      <div className="flex flex-row items-center w-full flex-wrap justify-center gap-1 mb-4">
        {links.map((l, i) => (
          <Button
            key={`${l}-${i}`}
            asChild
            variant="link"
            className="text-background"
          >
            <Link to={l.href} target={l.target ?? "_self"}>
              {l.text}
            </Link>
          </Button>
        ))}
      </div>
      <div className="flex flex-row items-center w-full flex-wrap justify-center gap-x-4 gap-y-3 mb-12">
        <Link to={instagramLink ?? ""} target="_blank">
          <FaInstagram className="text-primary" />
        </Link>
        <Link to={twitterLink ?? ""} target="_blank">
          <BsTwitterX className="text-primary" />
        </Link>
        <Link to={threadsLink ?? ""} target="_blank">
          <FaThreads className="text-primary" />
        </Link>
      </div>
      <p className="text-background/30 py-3 text-xs">
        Copyright © 2025 EcoNest
      </p>
    </footer>
  );
}

export { Footer };
