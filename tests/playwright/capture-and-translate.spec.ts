import { test, expect } from '@playwright/test';
import fs from 'fs';
import path from 'path';
import fetch from 'node-fetch';

const BASE_URL = 'http://localhost:3000'; // TODO: 替换为实际文档站点URL
const OUTPUT_DIR = 'docs/capture';
const IMAGE_DIR = path.join(OUTPUT_DIR, 'images');
const NAV_LABELS = [
  'Overview',
  'System Architecture',
  'Frontend Architecture',
  'Backend Architecture',
  'Data Flow',
  'Stock Monitoring',
  'Stock Data Sources',
  'Stock Groups',
  'Fund Monitoring',
  'AI-Powered Analysis',
  'Prompt Templates',
  'Configuration',
  'Notification System',
  'Scheduled Tasks',
  'Database System',
  'Developer Guide',
];

async function ensureDir(dir: string) {
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }
}

async function downloadImage(url: string, dest: string) {
  const res = await fetch(url);
  if (!res.ok) throw new Error(`Failed to fetch image: ${url}`);
  const buffer = await res.buffer();
  fs.writeFileSync(dest, buffer);
}

async function translateToChinese(text: string): Promise<string> {
  // TODO: 替换为实际翻译API或本地AI调用
  return text + '\n\n[此处为中文翻译，需接入翻译API]';
}

test('批量抓取文档内容并翻译', async ({ page }) => {
  await ensureDir(OUTPUT_DIR);
  await ensureDir(IMAGE_DIR);
  await page.goto(BASE_URL);

  for (const label of NAV_LABELS) {
    // 点击侧边栏条目
    const nav = page.getByRole('link', { name: label });
    await nav.click();
    await page.waitForTimeout(500); // 等待页面加载

    // 获取主内容
    const content = await page.locator('main').innerHTML();
    const mdPath = path.join(OUTPUT_DIR, `${label.replace(/\s+/g, '_')}.md`);
    fs.writeFileSync(mdPath, content, 'utf-8');

    // 下载图片
    const imgs = await page.locator('main img').all();
    for (let i = 0; i < imgs.length; i++) {
      const src = await imgs[i].getAttribute('src');
      if (src && !src.startsWith('data:')) {
        const imgUrl = src.startsWith('http') ? src : BASE_URL + src;
        const imgName = `${label.replace(/\s+/g, '_')}_${i}${path.extname(src) || '.png'}`;
        const imgPath = path.join(IMAGE_DIR, imgName);
        await downloadImage(imgUrl, imgPath);
      }
    }

    // 翻译内容
    const zhContent = await translateToChinese(content);
    const zhMdPath = path.join(OUTPUT_DIR, `${label.replace(/\s+/g, '_')}.zh.md`);
    fs.writeFileSync(zhMdPath, zhContent, 'utf-8');
  }
}); 