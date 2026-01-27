export function generateCronFromRecurring(time: string, days: number[]): string {
    const [hours, minutes] = time.split(':').map(Number);
    const validDays = days.sort((a, b) => a - b).join(',');

    // Format: mm hh * * days
    // Days in Cron (robfig): 0-6 (Sun-Sat)
    return `${minutes} ${hours} * * ${validDays === '' ? '*' : validDays}`;
}

export function generateCronFromInterval(value: number, unit: 'm' | 'h'): string {
    return `@every ${value}${unit}`;
}

export function explainCron(expression: string): string {
    if (expression.startsWith('@every')) {
        const match = expression.match(/@every (\d+)([mh])/);
        if (match) {
            const val = match[1];
            const unit = match[2] === 'm' ? 'minute(s)' : 'hour(s)';
            return `Every ${val} ${unit}`;
        }
        return expression;
    }

    const parts = expression.split(' ');
    if (parts.length === 5) {
        // Simple explanation for standard cron
        const [min, hour, , , days] = parts;
        return `At ${hour.padStart(2, '0')}:${min.padStart(2, '0')} on days: ${days === '*' ? 'All' : days}`;
    }

    return expression;
}
