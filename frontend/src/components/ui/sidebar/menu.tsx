"use client"

import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"
import { Skeleton } from "@/components/ui/skeleton"
import {
    Tooltip,
    TooltipContent,
    TooltipTrigger,
} from "@/components/ui/tooltip"
import { useSidebar } from "./context"
import { sidebarMenuButtonVariants } from "./variants"

export function SidebarMenu({
    className,
    ref,
    ...props
}: React.ComponentProps<"ul"> & {
    ref?: React.Ref<HTMLUListElement>
}) {
    return (
        <ul
            ref={ref}
            data-sidebar="menu"
            className={cn("flex w-full min-w-0 flex-col gap-1", className)}
            {...props}
        />
    )
}

export function SidebarMenuItem({
    className,
    ref,
    ...props
}: React.ComponentProps<"li"> & {
    ref?: React.Ref<HTMLLIElement>
}) {
    return (
        <li
            ref={ref}
            data-sidebar="menu-item"
            className={cn("group/menu-item relative", className)}
            {...props}
        />
    )
}

export function SidebarMenuButton({
    asChild = false,
    isActive = false,
    variant = "default",
    size = "default",
    tooltip,
    className,
    onClick,
    ref,
    ...props
}: React.ComponentProps<"button"> & {
    asChild?: boolean
    isActive?: boolean
    tooltip?: string | React.ComponentProps<typeof TooltipContent>
    ref?: React.Ref<HTMLButtonElement>
} & VariantProps<typeof sidebarMenuButtonVariants>) {
    const Comp = asChild ? Slot : "button"
    const { isMobile, state, setOpen } = useSidebar()

    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        // If sidebar is collapsed, expand it on click
        if (state === "collapsed" && !isMobile) {
            setOpen(true)
        } else {
            setOpen(false)
        }
        // Call the original onClick if provided
        onClick?.(event)
    }

    const button = (
        <Comp
            ref={ref}
            data-sidebar="menu-button"
            data-size={size}
            data-active={isActive}
            className={cn(sidebarMenuButtonVariants({ variant, size }), className)}
            onClick={handleClick}
            {...props}
        />
    )

    if (!tooltip) {
        return button
    }

    if (typeof tooltip === "string") {
        tooltip = {
            children: tooltip,
        }
    }

    return (
        <Tooltip>
            <TooltipTrigger asChild>{button}</TooltipTrigger>
            <TooltipContent
                side="right"
                align="center"
                hidden={state !== "collapsed" || isMobile}
                {...tooltip}
            />
        </Tooltip>
    )
}

export function SidebarMenuAction({
    className,
    asChild = false,
    showOnHover = false,
    ref,
    ...props
}: React.ComponentProps<"button"> & {
    asChild?: boolean
    showOnHover?: boolean
    ref?: React.Ref<HTMLButtonElement>
}) {
    const Comp = asChild ? Slot : "button"

    return (
        <Comp
            ref={ref}
            data-sidebar="menu-action"
            className={cn(
                "absolute right-1 top-1.5 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground outline-none ring-sidebar-ring transition-transform hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 peer-hover/menu-button:text-sidebar-accent-foreground [&>svg]:size-4 [&>svg]:shrink-0",
                // Increases the hit area of the button on mobile.
                "after:absolute after:-inset-2 after:md:hidden",
                "peer-data-[size=sm]/menu-button:top-1",
                "peer-data-[size=default]/menu-button:top-1.5",
                "peer-data-[size=lg]/menu-button:top-2.5",
                "group-data-[collapsible=icon]:hidden",
                showOnHover &&
                "group-focus-within/menu-item:opacity-100 group-hover/menu-item:opacity-100 data-[state=open]:opacity-100 peer-data-[active=true]/menu-button:text-sidebar-accent-foreground md:opacity-0",
                className
            )}
            {...props}
        />
    )
}

export function SidebarMenuBadge({
    className,
    ref,
    ...props
}: React.ComponentProps<"div"> & {
    ref?: React.Ref<HTMLDivElement>
}) {
    return (
        <div
            ref={ref}
            data-sidebar="menu-badge"
            className={cn(
                "pointer-events-none absolute right-1 flex h-5 min-w-5 select-none items-center justify-center rounded-md px-1 text-xs font-medium tabular-nums text-sidebar-foreground",
                "peer-hover/menu-button:text-sidebar-accent-foreground peer-data-[active=true]/menu-button:text-sidebar-accent-foreground",
                "peer-data-[size=sm]/menu-button:top-1",
                "peer-data-[size=default]/menu-button:top-1.5",
                "peer-data-[size=lg]/menu-button:top-2.5",
                "group-data-[collapsible=icon]:hidden",
                className
            )}
            {...props}
        />
    )
}

export function SidebarMenuSkeleton({
    className,
    showIcon = false,
    ref,
    ...props
}: React.ComponentProps<"div"> & {
    showIcon?: boolean
    ref?: React.Ref<HTMLDivElement>
}) {
    // Random width between 50 to 90%.
    const width = React.useMemo(() => {
        return `${Math.floor(Math.random() * 40) + 50}%`
    }, [])

    return (
        <div
            ref={ref}
            data-sidebar="menu-skeleton"
            className={cn("flex h-8 items-center gap-2 rounded-md px-2", className)}
            {...props}
        >
            {showIcon && (
                <Skeleton
                    className="size-4 rounded-md"
                    data-sidebar="menu-skeleton-icon"
                />
            )}
            <Skeleton
                className="h-4 max-w-[--skeleton-width] flex-1"
                data-sidebar="menu-skeleton-text"
                style={
                    {
                        "--skeleton-width": width,
                    } as React.CSSProperties
                }
            />
        </div>
    )
}

export function SidebarMenuSub({
    className,
    ref,
    ...props
}: React.ComponentProps<"ul"> & {
    ref?: React.Ref<HTMLUListElement>
}) {
    return (
        <ul
            ref={ref}
            data-sidebar="menu-sub"
            className={cn(
                "mx-3.5 flex min-w-0 translate-x-px flex-col gap-1 border-l border-sidebar-border px-2.5 py-0.5",
                "group-data-[collapsible=icon]:hidden",
                className
            )}
            {...props}
        />
    )
}

export function SidebarMenuSubItem({
    ref,
    ...props
}: React.ComponentProps<"li"> & {
    ref?: React.Ref<HTMLLIElement>
}) {
    return <li ref={ref} {...props} />
}

export function SidebarMenuSubButton({
    asChild = false,
    size = "md",
    isActive,
    className,
    ref,
    ...props
}: React.ComponentProps<"a"> & {
    asChild?: boolean
    size?: "sm" | "md"
    isActive?: boolean
    ref?: React.Ref<HTMLAnchorElement>
}) {
    const Comp = asChild ? Slot : "a"

    return (
        <Comp
            ref={ref}
            data-sidebar="menu-sub-button"
            data-size={size}
            data-active={isActive}
            className={cn(
                "flex h-7 min-w-0 -translate-x-px items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground outline-none ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-sidebar-accent-foreground",
                "data-[active=true]:bg-sidebar-accent data-[active=true]:text-sidebar-accent-foreground",
                size === "sm" && "text-xs",
                size === "md" && "text-sm",
                "group-data-[collapsible=icon]:hidden",
                className
            )}
            {...props}
        />
    )
}
