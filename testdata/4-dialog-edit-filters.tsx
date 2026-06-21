import React from "react";

function classNames(...classes: string[]) {
    return classes.filter(Boolean).join(' ');
}

function Header({className, ...rest}: React.ComponentProps<"div">) {
    const isNameInvalid = false;
    const isPatternInvalid = false;
    const isRegex = false;
    return (<>
        <div className="select-none mt-4 px-5 grid grid-cols-2 gap-1">
            Name
        </div>

        <div className={`h-8 ${isNameInvalid ? 'border-red-500 focus-visible:ring-red-500' : ''}`}>
            Pattern
        </div>

        <div className={`h-8 pr-8 ${isPatternInvalid ? 'border-red-500 focus-visible:ring-red-500' : ''}`}>
            Pattern
        </div>

        <div className={`border-0 border-l border-l-border rounded-r-[3px] rounded-l-none ${isRegex ? 'text-primary bg-primary/10' : 'hover:bg-primary/5'} ${isPatternInvalid ? 'border-l-red-500' : ''}`}>
            Pattern
        </div>

        <div className={classNames("mx-5 mt-1 mb-1 text-xs text-muted-foreground text-balance", className)} {...rest}>
            Pattern
        </div>
    </>);
}
