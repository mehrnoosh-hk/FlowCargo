import { SquarePen } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip"

export default function NewIssueButton() {
    return (
        <div className="w-full p-6 flex justify-center">
            <Tooltip>
                <TooltipTrigger asChild>
                    {/* 
                        BUTTON COLOR OPTIONS:
                        
                        Option 1: Use built-in variant (RECOMMENDED)
                        - variant="default"     → Primary brand color (purple)
                        - variant="secondary"   → Secondary color (cyan)
                        - variant="destructive" → Red for delete/error
                        - variant="outline"     → Subtle with border (current)
                        - variant="ghost"       → No background
                        
                        Option 2: Custom color with Tailwind classes
                        - Add className with bg-* and text-* utilities
                        
                        Option 3: Modify theme in globals.css
                        - Change --primary or --secondary colors globally
                    */}

                    {/* CURRENT: Primary brand color (purple/magenta) */}
                    <Button variant="default" size="icon">
                        <SquarePen />
                    </Button>
                </TooltipTrigger>
                <TooltipContent>
                    Create a new issue
                </TooltipContent>
            </Tooltip>
        </div>
    )
}