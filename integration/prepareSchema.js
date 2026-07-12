const fs = require('fs');
const path = require('path');

function prepare(mode) {
    const schemaPath = path.join(__dirname, 'schema.prisma');
    const backupPath = path.join(__dirname, 'schema.prisma.backup');

    if (mode === 'postgres') {
        if (fs.existsSync(backupPath)) {
            fs.copyFileSync(backupPath, schemaPath);
            fs.unlinkSync(backupPath);
            console.log('Restored schema.prisma from postgres backup.');
        } else {
            console.log('No backup found, schema.prisma is already in postgres mode.');
        }
        return;
    }

    // mode === 'sqlite'
    if (!fs.existsSync(backupPath)) {
        fs.copyFileSync(schemaPath, backupPath);
        console.log('Created backup of postgres schema.prisma.');
    }

    const content = fs.readFileSync(schemaPath, 'utf8');
    const lines = content.split(/\r?\n/);

    const out = [];

    for (const line of lines) {
        let currentLine = line;
        const trimmed = currentLine.trim();

        // Swap provider string
        if (trimmed.startsWith('provider =')) {
            out.push('  provider = "sqlite"');
            continue;
        }

        // Delete any postgres-specific Unsupported fields
        if (currentLine.includes('Unsupported(')) {
            continue;
        }

        // Strip any @db.something attributes
        currentLine = currentLine.replace(/@db\.[A-Za-z0-9_]+(?:\([^)]*\))?/g, '');

        out.push(currentLine);
    }

    fs.writeFileSync(schemaPath, out.join('\n'));
    console.log('Prepared schema.prisma for sqlite.');
}

const mode = process.argv[2];
if (!mode || (mode !== 'sqlite' && mode !== 'postgres')) {
    console.error('Usage: node prepareSchema.js [sqlite|postgres]');
    process.exit(1);
}

prepare(mode);
