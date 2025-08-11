<script setup lang="ts">
import { copy } from "@/utils/clipboard";
import { Copy, Key } from "@vicons/ionicons5";
import { NButton, NIcon, NInput, NInputNumber, NModal, NSpace, useMessage } from "naive-ui";
import { ref } from "vue";

interface Props {
  modelValue: string;
  placeholder?: string;
  size?: "small" | "medium" | "large";
}

interface Emits {
  (e: "update:modelValue", value: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: "Pisahkan beberapa kunci dengan koma bahasa Inggris",
  size: "small",
});

const emit = defineEmits<Emits>();

const message = useMessage();

// 密钥生成弹窗相关
const showKeyGeneratorModal = ref(false);
const keyCount = ref(1);
const isGenerating = ref(false);

// 生成随机字符串
function generateRandomString(length: number): string {
  const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

// 生成密钥
function generateKeys(): string[] {
  const keys: string[] = [];
  for (let i = 0; i < keyCount.value; i++) {
    keys.push(`sk-${generateRandomString(48)}`);
  }
  return keys;
}

// 打开密钥生成器弹窗
function openKeyGenerator() {
  showKeyGeneratorModal.value = true;
  keyCount.value = 1;
}

// 确认生成密钥
function confirmGenerateKeys() {
  if (isGenerating.value) {
    return;
  }

  try {
    isGenerating.value = true;
    const newKeys = generateKeys();
    const currentValue = props.modelValue || "";

    let updatedValue = currentValue.trim();

    // 处理逗号兼容情况
    if (updatedValue && !updatedValue.endsWith(",")) {
      updatedValue += ",";
    }

    // 添加新生成的密钥
    if (updatedValue) {
      updatedValue += newKeys.join(",");
    } else {
      updatedValue = newKeys.join(",");
    }

    emit("update:modelValue", updatedValue);
    showKeyGeneratorModal.value = false;

    message.success(`Berhasil membuat ${keyCount.value} kunci`);
  } finally {
    isGenerating.value = false;
  }
}

// 复制代理密钥
async function copyProxyKeys() {
  const proxyKeys = props.modelValue || "";
  if (!proxyKeys.trim()) {
    message.warning("Tidak ada kunci untuk disalin");
    return;
  }

  // 将逗号分隔的密钥转换为换行分隔
  const formattedKeys = proxyKeys
    .split(",")
    .map(key => key.trim())
    .filter(key => key.length > 0)
    .join("\n");

  const success = await copy(formattedKeys);
  if (success) {
    message.success("Kunci disalin ke clipboard");
  } else {
    message.error("Gagal menyalin, silakan salin secara manual");
  }
}

// 处理输入框值变化
function handleInput(value: string) {
  emit("update:modelValue", value);
}
</script>

<template>
  <div class="proxy-keys-input">
    <n-input
      :value="modelValue"
      :placeholder="placeholder"
      clearable
      :size="size"
      @update:value="handleInput"
    >
      <template #suffix>
        <n-space :size="4" :wrap-item="false">
          <n-button text type="primary" :size="size" @click="openKeyGenerator">
            <template #icon>
              <n-icon :component="Key" />
            </template>
            Hasilkan
          </n-button>
          <n-button text type="tertiary" :size="size" @click="copyProxyKeys" style="opacity: 0.7">
            <template #icon>
              <n-icon :component="Copy" />
            </template>
            Salin
          </n-button>
        </n-space>
      </template>
    </n-input>

    <!-- 密钥生成器弹窗 -->
    <n-modal
      v-model:show="showKeyGeneratorModal"
      preset="dialog"
      title="Hasilkan Kunci Proksi"
      positive-text="Konfirmasi Pembuatan"
      negative-text="Batal"
      :positive-button-props="{ loading: isGenerating }"
      @positive-click="confirmGenerateKeys"
    >
      <n-space vertical :size="16">
        <div>
          <p style="margin: 0 0 8px 0; color: #666; font-size: 14px">
            Silakan masukkan jumlah kunci yang akan dibuat (maksimal 100):
          </p>
          <n-input-number
            v-model:value="keyCount"
            :min="1"
            :max="100"
            placeholder="Silakan masukkan jumlah"
            style="width: 100%"
            :disabled="isGenerating"
          />
        </div>
        <div style="color: #999; font-size: 12px; line-height: 1.4">
          <p>Kunci yang dihasilkan akan dimasukkan di akhir konten kotak input saat ini, dipisahkan oleh koma</p>
        </div>
      </n-space>
    </n-modal>
  </div>
</template>

<style scoped>
.proxy-keys-input {
  width: 100%;
}
</style>
