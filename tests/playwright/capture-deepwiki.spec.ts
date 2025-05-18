import { test, expect } from '@playwright/test';
import fs from 'fs';
import path from 'path';
import fetch from 'node-fetch';

const BASE_URL = 'https://deepwiki.com/ArvinLovegood/go-stock/1-overview';
const OUTPUT_DIR = 'docs/capture';
const IMAGE_DIR = path.join(OUTPUT_DIR, 'images');

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

function htmlToText(html: string): string {
  // 粗略去除HTML标签，仅保留文本内容
  return html.replace(/<style[\s\S]*?<\/style>/gi, '')
    .replace(/<script[\s\S]*?<\/script>/gi, '')
    .replace(/<[^>]+>/g, '')
    .replace(/\n{2,}/g, '\n')
    .trim();
}

test('抓取deepwiki go-stock文档内容', async ({ page }) => {
  await ensureDir(OUTPUT_DIR);
  await ensureDir(IMAGE_DIR);
  await page.goto(BASE_URL);

  // 获取所有侧边栏条目
  const navLinks = await page.locator('nav a, aside a').all();
  for (const nav of navLinks) {
    const label = await nav.innerText();
    await nav.click();
    await page.waitForTimeout(800);
    // 获取主内容
    const main = await page.locator('main').first();
    const contentHtml = await main.innerHTML();
    const contentText = htmlToText(contentHtml);
    const safeLabel = label.replace(/\s+/g, '_').replace(/[^\w_\-]/g, '');
    const mdPath = path.join(OUTPUT_DIR, `${safeLabel}.md`);
    fs.writeFileSync(mdPath, contentText, 'utf-8');
    // 下载图片
    const imgs = await main.locator('img').all();
    for (let i = 0; i < imgs.length; i++) {
      const src = await imgs[i].getAttribute('src');
      if (src && !src.startsWith('data:')) {
        const imgUrl = src.startsWith('http') ? src : 'https://deepwiki.com' + src;
        const imgName = `${safeLabel}_${i}${path.extname(src) || '.png'}`;
        const imgPath = path.join(IMAGE_DIR, imgName);
        await downloadImage(imgUrl, imgPath);
      }
    }
  }
}); 