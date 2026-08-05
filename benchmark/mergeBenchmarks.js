const fs = require('fs');
const path = require('path');

const target = path.join(__dirname, 'benchmarks.json');
let report = { updatedAt: new Date().toISOString(), domains: { orm: { name: "ORM Performance", description: "Performance benchmarks for orm", databases: {} } } };

if (fs.existsSync(target)) {
    try {
        report = JSON.parse(fs.readFileSync(target, 'utf8'));
    } catch (e) {}
}

if (!report.domains) {
    report.domains = {};
}
if (!report.domains.orm) {
    report.domains.orm = { name: "ORM Performance", description: "Performance benchmarks for orm", databases: {} };
}
if (!report.domains.orm.databases) {
    report.domains.orm.databases = {};
}

const artifactsDir = path.join(process.cwd(), 'artifacts');

if (fs.existsSync(artifactsDir)) {
    const subdirs = fs.readdirSync(artifactsDir);
    for (const subdir of subdirs) {
        // subdir is "sqlite" or "postgres"
        const filePath = path.join(artifactsDir, subdir, 'benchmarks.json');
        if (fs.existsSync(filePath)) {
            try {
                const artifactData = JSON.parse(fs.readFileSync(filePath, 'utf8'));
                if (artifactData.domains && artifactData.domains.orm && artifactData.domains.orm.databases) {
                    const dbs = artifactData.domains.orm.databases;
                    // Explicitly pull only the database matching the artifact directory name
                    if (dbs[subdir]) {
                        report.domains.orm.databases[subdir] = dbs[subdir];
                        console.log(`Successfully merged fresh '${subdir}' benchmark data.`);
                    }
                }
            } catch (e) {
                console.error(`Error reading artifact for ${subdir}:`, e);
            }
        }
    }
}

report.updatedAt = new Date().toISOString();
fs.writeFileSync(target, JSON.stringify(report, null, 2));
console.log('Merged benchmarks successfully written to:', target);
