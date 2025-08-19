import { DOMAttributes } from "react";
import { Link } from "react-router";
import { IconType } from "react-icons/lib";
import { IoIosArrowForward } from "react-icons/io";

import { useIsMobile } from "@/hooks/use-mobile";
import {
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuGroup,
} from "@/components/ui/dropdown-menu";

export type SidebarSubItem = {
  title: string;
  isLink?: boolean;
  icon?: IconType;
  onClick?: DOMAttributes<HTMLDivElement>["onClick"];
  href: string;
  target?: string;
};

export type SidebarItem = SidebarSubItem & {
  subItems?: SidebarSubItem[];
  onClick?: DOMAttributes<HTMLButtonElement>["onClick"];
};

export type SidebarItemSectionProps = {
  label?: string;
  items?: SidebarItem[];
};

function SidebarItemSection({ label, items }: SidebarItemSectionProps) {
  const isMobile = useIsMobile();

  return (
    <>
      {label != null ? <SidebarGroupLabel>{label}</SidebarGroupLabel> : null}
      {items != null && items.length > 0 ? (
        <SidebarMenu>
          {items.map((i) => (
            <SidebarMenuItem key={i.title}>
              {i.subItems != null && i.subItems.length > 0 ? (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <SidebarMenuButton
                      tooltip={i.title}
                      onClick={i.onClick ? i.onClick : () => null}
                      asChild={i.isLink}
                    >
                      {i.isLink ? (
                        <Link to={i.href}>
                          {i.icon && <i.icon />}
                          <span>{i.title}</span>
                          <IoIosArrowForward className="ml-auto" />
                        </Link>
                      ) : (
                        <>
                          {i.icon && <i.icon />}
                          <span>{i.title}</span>
                          <IoIosArrowForward className="ml-auto" />
                        </>
                      )}
                    </SidebarMenuButton>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                    side={isMobile ? "bottom" : "right"}
                    align="end"
                    sideOffset={12}
                  >
                    <DropdownMenuGroup>
                      {i.subItems.map((si) => (
                        <DropdownMenuItem
                          key={`${i.title}-${si.title}`}
                          onClick={si.onClick ? si.onClick : () => null}
                          asChild={si.isLink}
                        >
                          {si.isLink ? (
                            <Link to={si.href}>
                              {si.icon && <si.icon />}
                              <span>{si.title}</span>
                            </Link>
                          ) : (
                            <>
                              {si.icon && <si.icon />}
                              <span>{si.title}</span>
                            </>
                          )}
                        </DropdownMenuItem>
                      ))}
                    </DropdownMenuGroup>
                  </DropdownMenuContent>
                </DropdownMenu>
              ) : (
                <SidebarMenuButton
                  tooltip={i.title}
                  onClick={i.onClick ? i.onClick : () => null}
                  asChild={i.isLink}
                >
                  {i.isLink ? (
                    <Link to={i.href}>
                      {i.icon && <i.icon />}
                      <span>{i.title}</span>
                    </Link>
                  ) : (
                    <>
                      {i.icon && <i.icon />}
                      <span>{i.title}</span>
                    </>
                  )}
                </SidebarMenuButton>
              )}
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      ) : null}
    </>
  );
}

export { SidebarItemSection };
