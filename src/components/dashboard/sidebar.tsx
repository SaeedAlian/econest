import { DOMAttributes } from "react";
import { Link } from "react-router";
import { CiLogin } from "react-icons/ci";
import { MdOutlineCancel } from "react-icons/md";

import { cn } from "@/lib/utils";
import { useIsMobile } from "@/hooks/use-mobile";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  SidebarItemSection,
  SidebarItem,
} from "@/components/dashboard/sidebar-item-section";
import {
  Sidebar as Sbar,
  SidebarGroup,
  SidebarGroupContent,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuGroup,
  DropdownMenuLabel,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import logo from "@/assets/logo-dark-notext.svg";

export type SidebarProps = {
  mainButton?: {
    text: string;
    onClick?: DOMAttributes<HTMLButtonElement>["onClick"];
  };
  topItems?: SidebarItem[];
  middleItems?: SidebarItem[];
  bottomItems?: SidebarItem[];
  topLabel?: string;
  middleLabel?: string;
  bottomLabel?: string;
  user?: {
    name: string;
    email: string;
    avatar?: string;
  };
  userItems?: SidebarItem[];
  onLogout?: () => void;
};

function Sidebar({
  mainButton,
  topItems,
  bottomItems,
  bottomLabel,
  middleItems,
  middleLabel,
  topLabel,
  user,
  userItems,
  onLogout,
}: SidebarProps) {
  const isMobile = useIsMobile();

  return (
    <>
      <Sbar collapsible="offcanvas">
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem className="flex items-center gap-2">
              <SidebarMenuButton
                asChild
                className="data-[slot=sidebar-menu-button]:!p-1.5"
              >
                <Link to="/">
                  <img
                    src={logo}
                    alt="EcoNest"
                    className={cn(
                      "w-24 h-24 m-[-30px]",
                      isMobile ? "w-20 h-20 m-[-25px]" : "",
                    )}
                  />
                  EcoNest Inc.
                </Link>
              </SidebarMenuButton>
              {isMobile ? <SidebarTrigger Icon={MdOutlineCancel} /> : null}
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupContent className="flex flex-col gap-2">
              {mainButton ? (
                <SidebarMenu>
                  <SidebarMenuItem className="flex items-center gap-2">
                    <SidebarMenuButton
                      tooltip="Quick Create"
                      className="bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground min-w-8 duration-200 ease-linear"
                      onClick={
                        mainButton.onClick ? mainButton.onClick : () => null
                      }
                    >
                      <span>{mainButton.text}</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </SidebarMenu>
              ) : null}
              <SidebarItemSection items={topItems} label={topLabel} />
            </SidebarGroupContent>
          </SidebarGroup>
          <SidebarGroup>
            <SidebarGroupContent>
              <SidebarItemSection items={middleItems} label={middleLabel} />
            </SidebarGroupContent>
          </SidebarGroup>
          <SidebarGroup className="mt-auto">
            <SidebarGroupContent>
              <SidebarItemSection items={bottomItems} label={bottomLabel} />
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <SidebarMenu>
            <SidebarMenuItem>
              {user != null ? (
                <DropdownMenu>
                  <SidebarMenuButton
                    size="lg"
                    className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                  >
                    <DropdownMenuTrigger asChild>
                      <div className="flex items-center gap-2">
                        <Avatar className="h-8 w-8 rounded-lg grayscale">
                          {user?.avatar != null ? (
                            <AvatarImage src={user.avatar} alt={user.name} />
                          ) : null}
                          <AvatarFallback className="rounded-lg">
                            {user.name.slice(0, 2).toUpperCase()}
                          </AvatarFallback>
                        </Avatar>
                        <div className="grid flex-1 text-left text-sm leading-tight">
                          <span className="truncate font-medium">
                            {user.name}
                          </span>
                          <span className="text-muted-foreground truncate text-xs">
                            {user.email}
                          </span>
                        </div>
                      </div>
                    </DropdownMenuTrigger>
                  </SidebarMenuButton>
                  <DropdownMenuContent
                    className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                    side={isMobile ? "bottom" : "right"}
                    align="end"
                    sideOffset={4}
                  >
                    {user != null ? (
                      <>
                        <DropdownMenuLabel className="p-0 font-normal">
                          <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                            <Avatar className="h-8 w-8 rounded-lg">
                              {user?.avatar != null ? (
                                <AvatarImage
                                  src={user.avatar}
                                  alt={user.name}
                                />
                              ) : null}
                              <AvatarFallback className="rounded-lg">
                                {user.name.slice(0, 2).toUpperCase()}
                              </AvatarFallback>
                            </Avatar>
                            <div className="grid flex-1 text-left text-sm leading-tight">
                              <span className="truncate font-medium">
                                {user.name}
                              </span>
                              <span className="text-muted-foreground truncate text-xs">
                                {user.email}
                              </span>
                            </div>
                          </div>
                        </DropdownMenuLabel>
                        <DropdownMenuSeparator />
                      </>
                    ) : null}
                    {userItems != null && userItems.length > 0 ? (
                      <>
                        <DropdownMenuGroup>
                          {userItems.map((i) => (
                            <DropdownMenuItem
                              key={i.title}
                              onClick={i.onClick ? i.onClick : () => null}
                              asChild={i.isLink}
                            >
                              {i.isLink && i.href != null ? (
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
                            </DropdownMenuItem>
                          ))}
                        </DropdownMenuGroup>
                        <DropdownMenuSeparator />
                      </>
                    ) : null}
                    {onLogout != null ? (
                      <DropdownMenuItem
                        onClick={onLogout}
                        variant="destructive"
                      >
                        Log out
                      </DropdownMenuItem>
                    ) : null}
                  </DropdownMenuContent>
                </DropdownMenu>
              ) : (
                <SidebarMenuButton size="lg">
                  <CiLogin /> Login
                </SidebarMenuButton>
              )}
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </Sbar>
    </>
  );
}

export { Sidebar };
