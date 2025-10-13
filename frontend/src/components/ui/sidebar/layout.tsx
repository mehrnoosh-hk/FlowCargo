"use client"

import * as React from "react"
import { cn } from "@/lib/utils"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"

export function SidebarInput({
    className,
    ref,
    ...props
}: React.ComponentProps<typeof Input> & {
    ref?: React.Ref<HTMLInputElement>
}) {
    return (
        <Input
            ref={ref}
            data-sidebar="input"
            className={cn(
                "h-8 w-full bg-background shadow-none focus-visible:ring-2 focus-visible:ring-sidebar-ring",
                className
            )}
            {...props}
        />
    )
}

export function SidebarHeader({
    className,
    ref,
    ...props
}: React.ComponentProps<"div"> & {
    ref?: React.Ref<HTMLDivElement>
}) {
    return (
        <div
            ref={ref}
            data-sidebar="header"
            className={cn("flex flex-col gap-2 p-2", className)}
            {...props}
        />
    )
}

export function SidebarFooter({
    className,
    ref,
    ...props
}: React.ComponentProps<"div"> & {
    ref?: React.Ref<HTMLDivElement>
}) {
    return (
        <div
            ref={ref}
            data-sidebar="footer"
            className={cn("flex flex-col gap-2 p-2", className)}
            {...props}
        />
    )
}

export function SidebarSeparator({
    className,
    ref,
    ...props
}: React.ComponentProps<typeof Separator> & {
    ref?: React.Ref<HTMLHRElement>
}) {
    return (
        <Separator
            ref={ref}
            data-sidebar="separator"
            className={cn("mx-2 w-auto bg-sidebar-border", className)}
            {...props}
        />
    )
}

export function SidebarContent({
    className,
    ref,
    ...props
}: React.ComponentProps<"div"> & {
    ref?: React.Ref<HTMLDivElement>
}) {
    return (
        <div
            ref={ref}
            data-sidebar="content"
            className={cn(
                "flex min-h-0 flex-1 flex-col gap-2 overflow-auto group-data-[collapsible=icon]:overflow-hidden",
                className
            )}
            {...props}
        />
    )
}
