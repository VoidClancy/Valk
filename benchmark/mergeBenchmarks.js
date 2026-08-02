const fs = require('fs');
const path = require('path');

const target = path.join(__dirname, 'benchmarks.json');
let report = { updatedAt: '', domains: {} };

if (fs.existsSync(target)) {
    try {
        report = JSON.parse(fs.readFileSync(target, 'utf8'));
    } catch (e) {}
}

if (!report.domains) {
    report.domains = {};
}

function findFiles(dir, fileList = []) {
    if (!fs.existsSync(dir)) return fileList;
    const files = fs.readdirSync(dir);
    for (const file of files) {
        const filePath = path.join(dir, file);
        if (fs.statSync(filePath).isDirectory()) {
            findFiles(filePath, fileList);
        } else if (file === 'benchmarks.json') {
            fileList.push(filePath);
        }
    }
    return fileList;
}

const artifactsDir = path.join(process.cwd(), 'artifacts');
const jsonFiles = findFiles(artifactsDir);

for (const p of jsonFiles) {
    try {
        const data = JSON.parse(fs.readFileSync(p, 'utf8'));
        if (data.domains) {
            for (const [dKey, dVal] of Object.entries(data.domains)) {
                if (!report.domains[dKey]) {
                    report.domains[dKey] = dVal;
                } else if (dVal.databases) {
                    if (!report.domains[dKey].databases) {
                        report.domains[dKey].databases = {};
                    }
                    for (const [dbKey, dbVal] of Object.entries(dVal.databases)) {
                        report.domains[dKey].databases[dbKey] = dbVal;
                    }
                }
            }
        }
    } catch (e) {}
}

report.updatedAt = new Date().toISOString();
fs.writeFileSync(target, JSON.stringify(report, null, 2));
console.log('Merged benchmarks successfully written to:', target);
