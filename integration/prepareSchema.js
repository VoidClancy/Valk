const fs = require('fs');
const path = require('path');

function prepare(mode) {
    const schemaPath = path.join(__dirname, 'schema.prisma');
    const content = fs.readFileSync(schemaPath, 'utf8');
    const lines = content.split(/\r?\n/);

    const out = [];

    for (const line of lines) {
        let currentLine = line;
        const trimmed = currentLine.trim();

        // Swap provider string
        if (trimmed.startsWith('provider =')) {
            if (mode === 'sqlite') {
                out.push('  provider = "sqlite"');
            } else {
                out.push('  provider = "postgres"');
            }
            continue;
        }

        if (mode === 'sqlite') {
            // Delete any postgres-specific Unsupported fields
            if (currentLine.includes('Unsupported(')) {
                continue;
            }

            // Strip any @db.something attributes
            currentLine = currentLine.replace(/@db\.[A-Za-z0-9_]+(?:\([^)]*\))?/g, '');
        }

        out.push(currentLine);
    }

    fs.writeFileSync(schemaPath, out.join('\n'));
}

const mode = process.argv[2];
if (!mode || (mode !== 'sqlite' && mode !== 'postgres')) {
    console.error('Usage: node prepareSchema.js [sqlite|postgres]');
    process.exit(1);
}

prepare(mode);
