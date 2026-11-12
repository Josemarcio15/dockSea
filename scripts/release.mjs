#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import readline from 'readline';
import { execSync } from 'child_process';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, '..');

const versionFilePath = path.join(rootDir, 'VERSION');
const packageJsonPath = path.join(rootDir, 'package.json');
const frontendPkgJsonPath = path.join(rootDir, 'frontend', 'package.json');
const changelogPath = path.join(rootDir, 'CHANGELOG.md');

function readCurrentVersion() {
  if (fs.existsSync(versionFilePath)) {
    return fs.readFileSync(versionFilePath, 'utf-8').trim();
  }
  const pkg = JSON.parse(fs.readFileSync(packageJsonPath, 'utf-8'));
  return pkg.version;
}

function extractChangelogNotes(version) {
  if (!fs.existsSync(changelogPath)) return '';
  const content = fs.readFileSync(changelogPath, 'utf-8');
  
  const cleanVersion = version.replace(/^v/i, '').trim();
  const escaped = cleanVersion.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const sectionRegex = new RegExp(`##\\s*\\[?v?${escaped}\\]?([\\s\\S]*?)(?=(?:\\n##\\s|$))`, 'i');
  
  const match = content.match(sectionRegex);
  if (match && match[1]) {
    return match[1].trim();
  }
  return '';
}

function suggestNextVersion(current) {
  const match = current.match(/^(\d+\.\d+\.)(\d+)(.*)$/);
  if (match) {
    const nextPatch = parseInt(match[2], 10) + 1;
    return `${match[1]}${nextPatch}${match[3]}`;
  }
  return current;
}

function updateVersionsInFiles(newVersion) {
  fs.writeFileSync(versionFilePath, newVersion + '\n');

  if (fs.existsSync(packageJsonPath)) {
    const pkg = JSON.parse(fs.readFileSync(packageJsonPath, 'utf-8'));
    pkg.version = newVersion;
    fs.writeFileSync(packageJsonPath, JSON.stringify(pkg, null, 2) + '\n');
  }

  if (fs.existsSync(frontendPkgJsonPath)) {
    const fpkg = JSON.parse(fs.readFileSync(frontendPkgJsonPath, 'utf-8'));
    fpkg.version = newVersion;
    fs.writeFileSync(frontendPkgJsonPath, JSON.stringify(fpkg, null, 2) + '\n');
  }

  execSync('go run scripts/sync_version.go', { cwd: rootDir, stdio: 'inherit' });
}

function askPreFilled(query, defaultValue) {
  return new Promise((resolve) => {
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
    });

    rl.question(query, (answer) => {
      rl.close();
      resolve(answer.trim() || defaultValue);
    });

    if (defaultValue) {
      rl.write(defaultValue);
    }
  });
}

async function main() {
  console.log('\n\x1b[36m[DOCKSEA]\x1b[0m Release & Bump CLI\n');

  const currentVersion = readCurrentVersion();
  const nextVersionSuggestion = suggestNextVersion(currentVersion);

  try {
    const noCommit = process.argv.includes('--no-commit');

    console.log(`Versão atual : \x1b[33m${currentVersion}\x1b[0m`);
    const targetVersion = await askPreFilled('Nova versão  : ', nextVersionSuggestion);

    let targetTag = '';
    let commitMsg = '';
    let changelogNotes = extractChangelogNotes(targetVersion);

    if (changelogNotes) {
      console.log(`\n\x1b[32m[Changelog encontrado para ${targetVersion}]:\x1b[0m\n${changelogNotes}\n`);
    }

    if (!noCommit) {
      const suggestedTag = targetVersion.startsWith('v') ? targetVersion : `v${targetVersion}`;
      targetTag = await askPreFilled('Tag Git      : ', suggestedTag);

      const suggestedCommitMsg = `release: ${targetTag}`;
      commitMsg = await askPreFilled('Mensagem     : ', suggestedCommitMsg);
    }

    const tagAnnotation = changelogNotes
      ? `${commitMsg || `release: ${targetTag}`}\n\n${changelogNotes}`
      : (commitMsg || `release: ${targetTag}`);

    console.log('\n------------------------------------------------------------');
    console.log(`Versão : \x1b[33m${targetVersion}\x1b[0m`);
    if (!noCommit) {
      console.log(`Tag    : \x1b[32m${targetTag}\x1b[0m`);
      console.log(`Commit : \x1b[36m${commitMsg}\x1b[0m`);
    } else {
      console.log(`Git    : \x1b[33m--no-commit (sem commit/tag/push)\x1b[0m`);
    }
    console.log('------------------------------------------------------------');
    console.log('Build  : \x1b[36mOs pacotes (.deb, .AppImage, .exe) serão gerados pelo GitHub Actions\x1b[0m');
    console.log('------------------------------------------------------------\n');

    const confirmMsg = noCommit
      ? 'Confirma as alterações? (s/N): '
      : 'Confirma as alterações e enviar para o Git? (s/N): ';
    const confirm = await askPreFilled(confirmMsg, 's');
    if (confirm.toLowerCase() !== 's') {
      console.log('Operação cancelada.');
      process.exit(0);
    }

    console.log('\n[1/3] Atualizando VERSION, configs do Wails e package.json...');
    updateVersionsInFiles(targetVersion);
    console.log('   OK: Arquivos de versão sincronizados!');

    if (noCommit) {
      console.log('\n\x1b[33m[Flag --no-commit detectada: pulando git commit, tag e push]\x1b[0m\n');
      return;
    }

    console.log('\n[2/3] Criando commit e tag Git...');
    execSync('git add -A', { cwd: rootDir, stdio: 'inherit' });

    try {
      execSync(`git commit -m "${commitMsg.replace(/"/g, '\\"')}"`, {
        cwd: rootDir,
        stdio: 'inherit',
      });
    } catch {
      console.log('   (Nenhuma mudança nos arquivos para commitar)');
    }

    const tempTagFile = path.join(rootDir, '.git', 'TAG_EDITMSG_TEMP');
    fs.writeFileSync(tempTagFile, tagAnnotation, 'utf-8');
    try {
      execSync(`git tag -a ${targetTag} -F "${tempTagFile}"`, {
        cwd: rootDir,
        stdio: 'inherit',
      });
    } finally {
      if (fs.existsSync(tempTagFile)) fs.unlinkSync(tempTagFile);
    }
    console.log(`   OK: Tag ${targetTag} criada com sucesso!`);

    console.log('\n[3/3] Envio para o GitHub:');
    const currentBranch = execSync('git rev-parse --abbrev-ref HEAD', { cwd: rootDir }).toString().trim();
    const doPush = await askPreFilled(`Enviar branch ${currentBranch} e tag ${targetTag} para origin? (S/n): `, 's');

    if (doPush.toLowerCase() === 's') {
      console.log(`Enviando branch ${currentBranch}...`);
      execSync(`git push origin ${currentBranch}`, { cwd: rootDir, stdio: 'inherit' });
      console.log(`Enviando tag ${targetTag}...`);
      execSync(`git push origin ${targetTag}`, { cwd: rootDir, stdio: 'inherit' });
      console.log(`\n\x1b[32mSucesso total! Tag ${targetTag} enviada para o GitHub!\x1b[0m`);
      console.log('O GitHub Actions irá compilar os binários e publicar no Releases.\n');
    } else {
      console.log('\nTag criada apenas localmente. Para enviar manualmente:\n');
      console.log(`   git push origin ${currentBranch} && git push origin ${targetTag}\n`);
    }
  } catch (err) {
    console.error('\nErro durante o processo:', err.message);
    process.exit(1);
  }
}

main();

