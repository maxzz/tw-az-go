#!/usr/bin/env node

import { spawnSync } from "node:child_process";
import { readFileSync, chmodSync, existsSync, mkdirSync, rmSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const binariesDir = join(root, "binaries");
const pkg = JSON.parse(readFileSync(join(root, "package.json"), "utf8"));
const buildAll = process.argv.includes("--all");

const allTargets = [
    { os: "darwin", arch: "arm64" },
    { os: "darwin", arch: "amd64" },
    { os: "linux", arch: "amd64" },
    { os: "linux", arch: "arm64" },
    { os: "windows", arch: "amd64" },
];

const windowsTargets = [{ os: "windows", arch: "amd64" }];
const targets = buildAll ? allTargets : windowsTargets;

if (buildAll) {
    rmSync(binariesDir, { recursive: true, force: true });
}
mkdirSync(binariesDir, { recursive: true });

for (const { os, arch } of targets) {
    const platformKey = `${os}-${arch}`;
    const outDir = join(binariesDir, platformKey);
    mkdirSync(outDir, { recursive: true });

    const ext = os === "windows" ? ".exe" : "";
    const outFile = join(outDir, `twaz${ext}`);

    const result = spawnSync(
        "go",
        [
            "build",
            "-trimpath",
            "-ldflags",
            `-s -w -X twaz/twaz.Version=${pkg.version}`,
            "-o",
            outFile,
            "./cmd/twaz",
        ],
        {
            cwd: root,
            env: {
                ...process.env,
                GOOS: os,
                GOARCH: arch,
                CGO_ENABLED: "0",
            },
            stdio: "inherit",
        },
    );

    if (result.status !== 0) {
        process.exit(result.status ?? 1);
    }

    if (os !== "windows") {
        chmodSync(outFile, 0o755);
    }
}

const checkTarget = buildAll
    ? join(binariesDir, "darwin-arm64", "twaz")
    : join(binariesDir, "windows-amd64", "twaz.exe");

if (!existsSync(checkTarget)) {
    console.error("build-binaries: no binaries produced");
    process.exit(1);
}

const scope = buildAll ? "all platforms" : "windows-amd64";
console.log(`Built ${targets.length} platform binary/binaries (${scope}) into binaries/`);
