<script lang="ts">
    import { Info, TriangleAlert, OctagonX, BadgeCheck } from "lucide-svelte";

    type AlertType = "info" | "warning" | "error" | "success";

    let {
        type = "info",
        title,
        children,
    } = $props<{
        type?: AlertType;
        title: string;
        children: any;
    }>();

    const styles: Record<
        AlertType,
        {
            container: string;
            title: string;
            icon: typeof Info;
            iconColor: string;
        }
    > = {
        info: {
            container: "bg-blue-900/30 border-blue-500",
            title: "text-blue-400",
            icon: Info,
            iconColor: "text-blue-400",
        },
        warning: {
            container: "bg-yellow-900/30 border-yellow-500",
            title: "text-yellow-500",
            icon: TriangleAlert,
            iconColor: "text-yellow-500",
        },
        error: {
            container: "bg-red-900/30 border-red-500",
            title: "text-red-500",
            icon: OctagonX,
            iconColor: "text-red-500",
        },
        success: {
            container: "bg-green-900/30 border-green-500",
            title: "text-green-500",
            icon: BadgeCheck,
            iconColor: "text-green-500",
        },
    };

    let currentStyle = $derived(styles[type as AlertType]);
    let Icon = $derived(currentStyle.icon);
</script>

<div class="border-l-4 p-4 my-4 rounded-r-lg {currentStyle.container}">
    <div class="flex items-center gap-2 mb-2">
        <Icon class="h-5 w-5 {currentStyle.iconColor}" />
        <p class="font-bold {currentStyle.title}">{title}</p>
    </div>
    <div class="text-slate-300">
        {@render children()}
    </div>
</div>
