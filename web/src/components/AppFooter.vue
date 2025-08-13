<script setup lang="ts">
import { ChatbubbleOutline } from "@vicons/ionicons5";
import { NIcon, NTooltip } from "naive-ui";
</script>

<template>
  <footer class="app-footer">
    <div class="footer-container">
      <!-- Area informasi utama -->
      <div class="footer-main">

        <!-- Area tautan -->
        <div class="links-container">
          <n-tooltip trigger="hover" placement="top">
            <template #trigger>
              <a
                href="https://wa.me/6281392492669?text=Halo%20Admin"
                target="_blank"
                rel="noopener noreferrer"
                class="footer-link"
              >
                <n-icon :component="ChatbubbleOutline" :size="14" class="link-icon" />
                <span>WhatsApp</span>
              </a>
            </template>
            Hubungi Kami
          </n-tooltip>
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.app-footer {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  padding: 12px 24px;
  font-size: 14px;
  min-height: 52px;
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-main {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  line-height: 1.4;
}

.project-info {
  color: #666;
  font-weight: 500;
}

.project-info a {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
}

.project-info a:hover {
  text-decoration: underline;
}

/* 版本信息区域 */
.version-container {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.version-icon {
  display: flex;
  align-items: center;
}

.version-text {
  font-weight: 500;
  font-size: 13px;
  color: #666;
  white-space: nowrap;
}

.version-clickable {
  cursor: pointer;
}

.version-clickable:hover {
  background: rgba(240, 160, 32, 0.1);
  transform: translateY(-1px);
}

.version-checking {
  opacity: 0.7;
}

/* 链接区域 */
.links-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.footer-link {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #666;
  text-decoration: none;
  padding: 4px 6px;
  border-radius: 4px;
  transition: all 0.2s ease;
  font-size: 13px;
  white-space: nowrap;
}

.footer-link:hover {
  color: var(--primary-color, #18a058);
  background: rgba(24, 160, 88, 0.1);
  transform: translateY(-1px);
}

.link-icon {
  display: flex;
  align-items: center;
}

/* 版权信息区域 */
.copyright-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.copyright-text {
  color: #888;
  font-size: 12px;
}

.license-text {
  color: #888;
  font-size: 12px;
}

.author-link {
  font-weight: 600;
  color: #667eea;
  text-decoration: none;
}

.author-link:hover {
  text-decoration: underline !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-footer {
    padding: 10px 16px;
    height: auto;
  }

  .footer-main {
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }

  .footer-main :deep(.n-divider) {
    display: none;
  }

  .links-container {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .footer-main {
    gap: 6px;
  }

  .links-container {
    flex-wrap: wrap;
    justify-content: center;
    gap: 12px;
  }

  .project-info {
    font-size: 12px;
  }

  .footer-link {
    font-size: 12px;
  }
}
</style>
