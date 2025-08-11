<script setup lang="ts">
import { keysApi } from "@/api/keys";
import { settingsApi } from "@/api/settings";
import ProxyKeysInput from "@/components/common/ProxyKeysInput.vue";
import type { Group, GroupConfigOption, UpstreamInfo } from "@/types/models";
import { Add, Close, HelpCircleOutline, Remove } from "@vicons/ionicons5";
import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NInputNumber,
  NModal,
  NSelect,
  NTooltip,
  useMessage,
  type FormRules,
} from "naive-ui";
import { computed, reactive, ref, watch } from "vue";

interface Props {
  show: boolean;
  group?: Group | null;
}

interface Emits {
  (e: "update:show", value: boolean): void;
  (e: "success", value: Group): void;
  (e: "switchToGroup", groupId: number): void;
}

// 配置项类型
interface ConfigItem {
  key: string;
  value: number | string;
}

const props = withDefaults(defineProps<Props>(), {
  group: null,
});

const emit = defineEmits<Emits>();

const message = useMessage();
const loading = ref(false);
const formRef = ref();

// 表单数据接口
interface GroupFormData {
  name: string;
  display_name: string;
  description: string;
  upstreams: UpstreamInfo[];
  channel_type: "anthropic" | "gemini" | "openai";
  sort: number;
  test_model: string;
  validation_endpoint: string;
  param_overrides: string;
  config: Record<string, number | string>;
  configItems: ConfigItem[];
  proxy_keys: string;
}

// 表单数据
const formData = reactive<GroupFormData>({
  name: "",
  display_name: "",
  description: "",
  upstreams: [
    {
      url: "",
      weight: 1,
    },
  ] as UpstreamInfo[],
  channel_type: "openai",
  sort: 1,
  test_model: "",
  validation_endpoint: "",
  param_overrides: "",
  config: {},
  configItems: [] as ConfigItem[],
  proxy_keys: "",
});

const channelTypeOptions = ref<{ label: string; value: string }[]>([]);
const configOptions = ref<GroupConfigOption[]>([]);
const channelTypesFetched = ref(false);
const configOptionsFetched = ref(false);

// 跟踪用户是否已手动修改过字段（仅在新增模式下使用）
const userModifiedFields = ref({
  test_model: false,
  upstream: false,
});

// Hasilkan placeholder secara dinamis berdasarkan jenis saluran
const testModelPlaceholder = computed(() => {
  switch (formData.channel_type) {
    case "openai":
      return "gpt-4.1-nano";
    case "gemini":
      return "gemini-2.0-flash-lite";
    case "anthropic":
      return "claude-3-haiku-20240307";
    default:
      return "Silakan masukkan nama model";
  }
});

const upstreamPlaceholder = computed(() => {
  switch (formData.channel_type) {
    case "openai":
      return "https://api.openai.com";
    case "gemini":
      return "https://generativelanguage.googleapis.com";
    case "anthropic":
      return "https://api.anthropic.com";
    default:
      return "Silakan masukkan alamat upstream";
  }
});

const validationEndpointPlaceholder = computed(() => {
  switch (formData.channel_type) {
    case "openai":
      return "/v1/chat/completions";
    case "anthropic":
      return "/v1/messages";
    case "gemini":
      return ""; // Gemini tidak menampilkan bidang ini
    default:
      return "Silakan masukkan jalur endpoint validasi";
  }
});

// Aturan validasi formulir
const rules: FormRules = {
  name: [
    {
      required: true,
      message: "Silakan masukkan nama grup",
      trigger: ["blur", "input"],
    },
    {
      pattern: /^[a-z0-9_-]{3,30}$/,
      message: "Hanya boleh berisi huruf kecil, angka, tanda hubung, atau garis bawah, dengan panjang 3-30 karakter",
      trigger: ["blur", "input"],
    },
  ],
  channel_type: [
    {
      required: true,
      message: "Silakan pilih jenis saluran",
      trigger: ["blur", "change"],
    },
  ],
  test_model: [
    {
      required: true,
      message: "Silakan masukkan model tes",
      trigger: ["blur", "input"],
    },
  ],
  upstreams: [
    {
      type: "array",
      min: 1,
      message: "Diperlukan setidaknya satu alamat upstream",
      trigger: ["blur", "change"],
    },
  ],
};

// 监听弹窗显示状态
watch(
  () => props.show,
  show => {
    if (show) {
      if (!channelTypesFetched.value) {
        fetchChannelTypes();
      }
      if (!configOptionsFetched.value) {
        fetchGroupConfigOptions();
      }
      resetForm();
      if (props.group) {
        loadGroupData();
      }
    }
  }
);

// 监听渠道类型变化，在新增模式下智能更新默认值
watch(
  () => formData.channel_type,
  (_newChannelType, oldChannelType) => {
    if (!props.group && oldChannelType) {
      // 仅在新增模式且不是初始设置时处理
      // 检查测试模型是否应该更新（为空或是旧渠道类型的默认值）
      if (
        !userModifiedFields.value.test_model ||
        formData.test_model === getOldDefaultTestModel(oldChannelType)
      ) {
        formData.test_model = testModelPlaceholder.value;
        userModifiedFields.value.test_model = false;
      }

      // 检查第一个上游地址是否应该更新
      if (
        formData.upstreams.length > 0 &&
        (!userModifiedFields.value.upstream ||
          formData.upstreams[0].url === getOldDefaultUpstream(oldChannelType))
      ) {
        formData.upstreams[0].url = upstreamPlaceholder.value;
        userModifiedFields.value.upstream = false;
      }
    }
  }
);

// 获取旧渠道类型的默认值（用于比较）
function getOldDefaultTestModel(channelType: string): string {
  switch (channelType) {
    case "openai":
      return "gpt-4.1-nano";
    case "gemini":
      return "gemini-2.0-flash-lite";
    case "anthropic":
      return "claude-3-haiku-20240307";
    default:
      return "";
  }
}

function getOldDefaultUpstream(channelType: string): string {
  switch (channelType) {
    case "openai":
      return "https://api.openai.com";
    case "gemini":
      return "https://generativelanguage.googleapis.com";
    case "anthropic":
      return "https://api.anthropic.com";
    default:
      return "";
  }
}

// 重置表单
function resetForm() {
  const isCreateMode = !props.group;
  const defaultChannelType = "openai";

  // 先设置渠道类型，这样 computed 属性能正确计算默认值
  formData.channel_type = defaultChannelType;

  Object.assign(formData, {
    name: "",
    display_name: "",
    description: "",
    upstreams: [
      {
        url: isCreateMode ? upstreamPlaceholder.value : "",
        weight: 1,
      },
    ],
    channel_type: defaultChannelType,
    sort: 1,
    test_model: isCreateMode ? testModelPlaceholder.value : "",
    validation_endpoint: "",
    param_overrides: "",
    config: {},
    configItems: [],
    proxy_keys: "",
  });

  // 重置用户修改状态追踪
  if (isCreateMode) {
    userModifiedFields.value = {
      test_model: false,
      upstream: false,
    };
  }
}

// 加载分组数据（编辑模式）
function loadGroupData() {
  if (!props.group) {
    return;
  }

  const configItems = Object.entries(props.group.config || {}).map(([key, value]) => {
    return {
      key,
      value,
    };
  });
  Object.assign(formData, {
    name: props.group.name || "",
    display_name: props.group.display_name || "",
    description: props.group.description || "",
    upstreams: props.group.upstreams?.length
      ? [...props.group.upstreams]
      : [{ url: "", weight: 1 }],
    channel_type: props.group.channel_type || "openai",
    sort: props.group.sort || 1,
    test_model: props.group.test_model || "",
    validation_endpoint: props.group.validation_endpoint || "",
    param_overrides: JSON.stringify(props.group.param_overrides || {}, null, 2),
    config: {},
    configItems,
    proxy_keys: props.group.proxy_keys || "",
  });
}

async function fetchChannelTypes() {
  const options = (await settingsApi.getChannelTypes()) || [];
  channelTypeOptions.value =
    options?.map((type: string) => ({
      label: type,
      value: type,
    })) || [];
  channelTypesFetched.value = true;
}

// 添加上游地址
function addUpstream() {
  formData.upstreams.push({
    url: "",
    weight: 1,
  });
}

// Hapus alamat upstream
function removeUpstream(index: number) {
  if (formData.upstreams.length > 1) {
    formData.upstreams.splice(index, 1);
  } else {
    message.warning("Diperlukan setidaknya satu alamat upstream");
  }
}

async function fetchGroupConfigOptions() {
  const options = await keysApi.getGroupConfigOptions();
  configOptions.value = options || [];
  configOptionsFetched.value = true;
}

// 添加配置项
function addConfigItem() {
  formData.configItems.push({
    key: "",
    value: "",
  });
}

// 删除配置项
function removeConfigItem(index: number) {
  formData.configItems.splice(index, 1);
}

// 当配置项的key改变时，设置默认值
function handleConfigKeyChange(index: number, key: string) {
  const option = configOptions.value.find(opt => opt.key === key);
  if (option) {
    formData.configItems[index].value = option.default_value;
  }
}

const getConfigOption = (key: string) => {
  return configOptions.value.find(opt => opt.key === key);
};

// 关闭弹窗
function handleClose() {
  emit("update:show", false);
}

// 提交表单
async function handleSubmit() {
  if (loading.value) {
    return;
  }

  try {
    await formRef.value?.validate();

    loading.value = true;

    // Validasi format JSON
    let paramOverrides = {};
    if (formData.param_overrides) {
      try {
        paramOverrides = JSON.parse(formData.param_overrides);
      } catch {
        message.error("Timpa parameter harus dalam format JSON yang valid");
        return;
      }
    }

    // 将configItems转换为config对象
    const config: Record<string, number | string> = {};
    formData.configItems.forEach((item: ConfigItem) => {
      if (item.key && item.key.trim()) {
        const option = configOptions.value.find(opt => opt.key === item.key);
        if (option && typeof option.default_value === "number" && typeof item.value === "string") {
          const numValue = Number(item.value);
          config[item.key] = isNaN(numValue) ? 0 : numValue;
        } else {
          config[item.key] = item.value;
        }
      }
    });

    // 构建提交数据
    const submitData = {
      name: formData.name,
      display_name: formData.display_name,
      description: formData.description,
      upstreams: formData.upstreams.filter((upstream: UpstreamInfo) => upstream.url.trim()),
      channel_type: formData.channel_type,
      sort: formData.sort,
      test_model: formData.test_model,
      validation_endpoint: formData.validation_endpoint,
      param_overrides: paramOverrides,
      config,
      proxy_keys: formData.proxy_keys,
    };

    let res: Group;
    if (props.group?.id) {
      // 编辑模式
      res = await keysApi.updateGroup(props.group.id, submitData);
    } else {
      // 新建模式
      res = await keysApi.createGroup(submitData);
    }

    emit("success", res);
    // 如果是新建模式，发出切换到新分组的事件
    if (!props.group?.id && res.id) {
      emit("switchToGroup", res.id);
    }
    handleClose();
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <n-modal :show="show" @update:show="handleClose" class="group-form-modal">
    <n-card
      class="group-form-card"
      :title="group ? 'Edit Grup' : 'Buat Grup'"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
    >
      <template #header-extra>
        <n-button quaternary circle @click="handleClose">
          <template #icon>
            <n-icon :component="Close" />
          </template>
        </n-button>
      </template>

      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-placement="left"
        label-width="120px"
        require-mark-placement="right-hanging"
        class="group-form"
      >
        <!-- Informasi Dasar -->
        <div class="form-section">
          <h4 class="section-title">Informasi Dasar</h4>

          <!-- Nama grup dan nama tampilan dalam satu baris -->
          <div class="form-row">
            <n-form-item label="Nama Grup" path="name" class="form-item-half">
              <template #label>
                <div class="form-label-with-tooltip">
                  Nama Grup
                  <n-tooltip trigger="hover" placement="top">
                    <template #trigger>
                      <n-icon :component="HelpCircleOutline" class="help-icon" />
                    </template>
                    Sebagai bagian dari rute API, hanya boleh berisi huruf kecil, angka, tanda hubung, atau garis bawah, dengan panjang 3-30 karakter. Contoh: gemini, openai-2
                  </n-tooltip>
                </div>
              </template>
              <n-input v-model:value="formData.name" placeholder="gemini" />
            </n-form-item>

            <n-form-item label="Nama Tampilan" path="display_name" class="form-item-half">
              <template #label>
                <div class="form-label-with-tooltip">
                  Nama Tampilan
                  <n-tooltip trigger="hover" placement="top">
                    <template #trigger>
                      <n-icon :component="HelpCircleOutline" class="help-icon" />
                    </template>
                    Nama yang ramah untuk ditampilkan di antarmuka, dapat berisi karakter Cina dan khusus. Jika tidak diisi, nama grup akan digunakan sebagai nama tampilan
                  </n-tooltip>
                </div>
              </template>
              <n-input v-model:value="formData.display_name" placeholder="Google Gemini" />
            </n-form-item>
          </div>

          <!-- Jenis saluran dan urutan dalam satu baris -->
          <div class="form-row">
            <n-form-item label="Jenis Saluran" path="channel_type" class="form-item-half">
              <template #label>
                <div class="form-label-with-tooltip">
                  Jenis Saluran
                  <n-tooltip trigger="hover" placement="top">
                    <template #trigger>
                      <n-icon :component="HelpCircleOutline" class="help-icon" />
                    </template>
                    Pilih jenis penyedia API, yang menentukan format permintaan dan metode otentikasi. Mendukung layanan AI utama seperti OpenAI, Gemini, Anthropic, dll.
                  </n-tooltip>
                </div>
              </template>
              <n-select
                v-model:value="formData.channel_type"
                :options="channelTypeOptions"
                placeholder="Silakan pilih jenis saluran"
              />
            </n-form-item>

            <n-form-item label="Urutkan" path="sort" class="form-item-half">
              <template #label>
                <div class="form-label-with-tooltip">
                  Urutkan
                  <n-tooltip trigger="hover" placement="top">
                    <template #trigger>
                      <n-icon :component="HelpCircleOutline" class="help-icon" />
                    </template>
                    Menentukan urutan tampilan grup dalam daftar, angka yang lebih kecil akan ditampilkan lebih dulu. Disarankan untuk menggunakan angka interval seperti 10, 20, 30 untuk memudahkan penyesuaian di masa mendatang
                  </n-tooltip>
                </div>
              </template>
              <n-input-number
                v-model:value="formData.sort"
                :min="0"
                placeholder="Nilai urutan"
                style="width: 100%"
              />
            </n-form-item>
          </div>

          <!-- Model tes dan jalur tes dalam satu baris -->
          <div class="form-row">
            <n-form-item label="Model Tes" path="test_model" class="form-item-half">
              <template #label>
                <div class="form-label-with-tooltip">
                  Model Tes
                  <n-tooltip trigger="hover" placement="top">
                    <template #trigger>
                      <n-icon :component="HelpCircleOutline" class="help-icon" />
                    </template>
                    Nama model yang digunakan untuk memvalidasi kunci API. Sistem akan menggunakan model ini untuk mengirim permintaan tes untuk memeriksa apakah kunci tersedia, harap gunakan model yang ringan dan cepat
                  </n-tooltip>
                </div>
              </template>
              <n-input
                v-model:value="formData.test_model"
                :placeholder="testModelPlaceholder"
                @input="() => !props.group && (userModifiedFields.test_model = true)"
              />
            </n-form-item>

            <n-form-item
              label="Jalur Tes"
              path="validation_endpoint"
              class="form-item-half"
              v-if="formData.channel_type !== 'gemini'"
            >
              <template #label>
                <div class="form-label-with-tooltip">
                  Jalur Tes
                  <n-tooltip trigger="hover" placement="top">
                    <template #trigger>
                      <n-icon :component="HelpCircleOutline" class="help-icon" />
                    </template>
                    <div>
                      Jalur endpoint API kustom untuk memvalidasi kunci. Jika tidak diisi, jalur default akan digunakan:
                      <br />
                      • OpenAI: /v1/chat/completions
                      <br />
                      • Anthropic: /v1/messages
                      <br />
                      Jika Anda perlu menggunakan jalur non-standar, silakan isi jalur API lengkap di sini
                    </div>
                  </n-tooltip>
                </div>
              </template>
              <n-input
                v-model:value="formData.validation_endpoint"
                :placeholder="validationEndpointPlaceholder || 'Opsional, jalur API kustom untuk memvalidasi kunci'"
              />
            </n-form-item>

            <!-- Ketika saluran gemini, jalur tes tidak ditampilkan, diperlukan div placeholder untuk menjaga tata letak -->
            <div v-else class="form-item-half" />
          </div>

          <!-- Kunci Proxy -->
          <n-form-item label="Kunci Proxy" path="proxy_keys">
            <template #label>
              <div class="form-label-with-tooltip">
                Kunci Proxy
                <n-tooltip trigger="hover" placement="top">
                  <template #trigger>
                    <n-icon :component="HelpCircleOutline" class="help-icon" />
                  </template>
                  Kunci proxy khusus grup, digunakan untuk mengakses endpoint proxy grup ini. Pisahkan beberapa kunci dengan koma.
                </n-tooltip>
              </div>
            </template>
            <proxy-keys-input
              v-model="formData.proxy_keys"
              placeholder="Pisahkan beberapa kunci dengan koma bahasa Inggris , "
              size="medium"
            />
          </n-form-item>

          <!-- Deskripsi menempati satu baris -->
          <n-form-item label="Deskripsi" path="description">
            <template #label>
              <div class="form-label-with-tooltip">
                Deskripsi
                <n-tooltip trigger="hover" placement="top">
                  <template #trigger>
                    <n-icon :component="HelpCircleOutline" class="help-icon" />
                  </template>
                  Deskripsi terperinci grup, membantu anggota tim memahami tujuan dan fitur grup. Mendukung teks multibaris
                </n-tooltip>
              </div>
            </template>
            <n-input
              v-model:value="formData.description"
              type="textarea"
              placeholder=""
              :rows="1"
              :autosize="{ minRows: 1, maxRows: 5 }"
              style="resize: none"
            />
          </n-form-item>
        </div>

        <!-- Alamat Upstream -->
        <div class="form-section" style="margin-top: 10px">
          <h4 class="section-title">Alamat Upstream</h4>
          <n-form-item
            v-for="(upstream, index) in formData.upstreams"
            :key="index"
            :label="`Upstream ${index + 1}`"
            :path="`upstreams[${index}].url`"
            :rule="{
              required: true,
              message: '',
              trigger: ['blur', 'input'],
            }"
          >
            <template #label>
              <div class="form-label-with-tooltip">
                Upstream {{ index + 1 }}
                <n-tooltip trigger="hover" placement="top">
                  <template #trigger>
                    <n-icon :component="HelpCircleOutline" class="help-icon" />
                  </template>
                  Alamat URL lengkap server API. Beberapa upstream dapat mencapai penyeimbangan beban dan failover, meningkatkan ketersediaan layanan
                </n-tooltip>
              </div>
            </template>
            <div class="upstream-row">
              <div class="upstream-url">
                <n-input
                  v-model:value="upstream.url"
                  :placeholder="upstreamPlaceholder"
                  @input="() => !props.group && index === 0 && (userModifiedFields.upstream = true)"
                />
              </div>
              <div class="upstream-weight">
                <span class="weight-label">Bobot</span>
                <n-tooltip trigger="hover" placement="top" style="width: 100%">
                  <template #trigger>
                    <n-input-number
                      v-model:value="upstream.weight"
                      :min="1"
                      placeholder="Bobot"
                      style="width: 100%"
                    />
                  </template>
                  Bobot penyeimbangan beban, semakin besar nilainya, semakin tinggi probabilitas untuk dipilih. Misalnya: probabilitas upstream dengan bobot 2 untuk dipilih adalah dua kali lipat dari upstream dengan bobot 1
                </n-tooltip>
              </div>
              <div class="upstream-actions">
                <n-button
                  v-if="formData.upstreams.length > 1"
                  @click="removeUpstream(index)"
                  type="error"
                  quaternary
                  circle
                  size="small"
                >
                  <template #icon>
                    <n-icon :component="Remove" />
                  </template>
                </n-button>
              </div>
            </div>
          </n-form-item>

          <n-form-item>
            <n-button @click="addUpstream" dashed style="width: 100%">
              <template #icon>
                <n-icon :component="Add" />
              </template>
              Tambahkan Alamat Upstream
            </n-button>
          </n-form-item>
        </div>

        <!-- Konfigurasi Lanjutan -->
        <div class="form-section" style="margin-top: 10px">
          <n-collapse>
            <n-collapse-item name="advanced">
              <template #header>Konfigurasi Lanjutan</template>
              <div class="config-section">
                <h5 class="config-title-with-tooltip">
                  Konfigurasi Grup
                  <n-tooltip trigger="hover" placement="top">
                    <template #trigger>
                      <n-icon :component="HelpCircleOutline" class="help-icon config-help" />
                    </template>
                    Parameter konfigurasi khusus untuk grup ini, seperti waktu habis, jumlah percobaan ulang, dll. Konfigurasi ini akan menimpa pengaturan default global
                  </n-tooltip>
                </h5>

                <div class="config-items">
                  <n-form-item
                    v-for="(configItem, index) in formData.configItems"
                    :key="index"
                    class="config-item-row"
                    :label="`Konfigurasi ${index + 1}`"
                    :path="`configItems[${index}].key`"
                    :rule="{
                      required: true,
                      message: '',
                      trigger: ['blur', 'change'],
                    }"
                  >
                    <template #label>
                      <div class="form-label-with-tooltip">
                        Konfigurasi {{ index + 1 }}
                        <n-tooltip trigger="hover" placement="top">
                          <template #trigger>
                            <n-icon :component="HelpCircleOutline" class="help-icon" />
                          </template>
                          Pilih jenis parameter yang akan dikonfigurasi, lalu atur nilai yang sesuai. Parameter yang berbeda memiliki fungsi dan rentang nilai yang berbeda
                        </n-tooltip>
                      </div>
                    </template>
                    <div class="config-item-content">
                      <div class="config-select">
                        <n-select
                          v-model:value="configItem.key"
                          :options="
                            configOptions.map(opt => ({
                              label: opt.name,
                              value: opt.key,
                              disabled:
                                formData.configItems
                                  .map((item: ConfigItem) => item.key)
                                  ?.includes(opt.key) && opt.key !== configItem.key,
                            }))
                          "
                          placeholder="Silakan pilih parameter konfigurasi"
                          @update:value="value => handleConfigKeyChange(index, value)"
                          clearable
                        />
                      </div>
                      <div class="config-value">
                        <n-tooltip trigger="hover" placement="top">
                          <template #trigger>
                            <n-input-number
                              v-if="typeof configItem.value === 'number'"
                              v-model:value="configItem.value"
                              placeholder="Nilai parameter"
                              :precision="0"
                              style="width: 100%"
                            />
                            <n-input v-else v-model:value="configItem.value" placeholder="Nilai parameter" />
                          </template>
                          {{ getConfigOption(configItem.key)?.description || "Atur nilai item konfigurasi ini" }}
                        </n-tooltip>
                      </div>
                      <div class="config-actions">
                        <n-button
                          @click="removeConfigItem(index)"
                          type="error"
                          quaternary
                          circle
                          size="small"
                        >
                          <template #icon>
                            <n-icon :component="Remove" />
                          </template>
                        </n-button>
                      </div>
                    </div>
                  </n-form-item>
                </div>

                <div style="margin-top: 12px; padding-left: 120px">
                  <n-button
                    @click="addConfigItem"
                    dashed
                    style="width: 100%"
                    :disabled="formData.configItems.length >= configOptions.length"
                  >
                    <template #icon>
                      <n-icon :component="Add" />
                    </template>
                    Tambahkan Parameter Konfigurasi
                  </n-button>
                </div>
              </div>
              <div class="config-section">
                <n-form-item path="param_overrides">
                  <template #label>
                    <div class="form-label-with-tooltip">
                      Timpa Parameter
                      <n-tooltip trigger="hover" placement="top">
                        <template #trigger>
                          <n-icon :component="HelpCircleOutline" class="help-icon config-help" />
                        </template>
                        Gunakan format JSON untuk menentukan parameter permintaan API yang akan ditimpa. Contoh: {"temperature":
                        0.7}. Parameter ini akan digabungkan ke dalam parameter asli saat mengirim permintaan
                      </n-tooltip>
                    </div>
                  </template>
                  <n-input
                    v-model:value="formData.param_overrides"
                    type="textarea"
                    placeholder='{"temperature": 0.7}'
                    :rows="4"
                  />
                </n-form-item>
              </div>
            </n-collapse-item>
          </n-collapse>
        </div>
      </n-form>

      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px">
          <n-button @click="handleClose">Batal</n-button>
          <n-button type="primary" @click="handleSubmit" :loading="loading">
            {{ group ? "Perbarui" : "Buat" }}
          </n-button>
        </div>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped>
.group-form-modal {
  width: 800px;
  --n-color: rgba(255, 255, 255, 0.95);
}

.form-section {
  margin-top: 20px;
}

.section-title {
  font-size: 1rem;
  font-weight: 600;
  color: #374151;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid rgba(102, 126, 234, 0.1);
}

:deep(.n-form-item-label) {
  font-weight: 500;
}

:deep(.n-form-item-blank) {
  flex-grow: 1;
}

:deep(.n-input) {
  --n-border-radius: 6px;
}

:deep(.n-select) {
  --n-border-radius: 6px;
}

:deep(.n-input-number) {
  --n-border-radius: 6px;
}

:deep(.n-card-header) {
  border-bottom: 1px solid rgba(239, 239, 245, 0.8);
  padding: 10px 20px;
}

:deep(.n-card__content) {
  max-height: calc(100vh - 68px - 61px - 50px);
  overflow-y: auto;
}

:deep(.n-card__footer) {
  border-top: 1px solid rgba(239, 239, 245, 0.8);
  padding: 10px 15px;
}

:deep(.n-form-item-feedback-wrapper) {
  min-height: 10px;
}

.config-section {
  margin-top: 16px;
}

.config-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: #374151;
  margin: 0 0 12px 0;
}

.form-label {
  margin-left: 25px;
  margin-right: 10px;
  height: 34px;
  line-height: 34px;
  font-weight: 500;
}

/* Tooltip相关样式 */
.form-label-with-tooltip {
  display: flex;
  align-items: center;
  gap: 6px;
}

.help-icon {
  color: #9ca3af;
  font-size: 14px;
  cursor: help;
  transition: color 0.2s ease;
}

.help-icon:hover {
  color: #667eea;
}

.section-title-with-tooltip {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.section-help {
  font-size: 16px;
}

.collapse-header-with-tooltip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 500;
}

.collapse-help {
  font-size: 14px;
}

.config-title-with-tooltip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  color: #374151;
  margin: 0 0 12px 0;
}

.config-help {
  font-size: 13px;
}

/* 增强表单样式 */
:deep(.n-form-item-label) {
  font-weight: 500;
  color: #374151;
}

:deep(.n-input) {
  --n-border-radius: 8px;
  --n-border: 1px solid #e5e7eb;
  --n-border-hover: 1px solid #667eea;
  --n-border-focus: 1px solid #667eea;
  --n-box-shadow-focus: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

:deep(.n-select) {
  --n-border-radius: 8px;
}

:deep(.n-input-number) {
  --n-border-radius: 8px;
}

:deep(.n-button) {
  --n-border-radius: 8px;
}

/* 美化tooltip */
:deep(.n-tooltip__trigger) {
  display: inline-flex;
  align-items: center;
}

:deep(.n-tooltip) {
  --n-font-size: 13px;
  --n-border-radius: 8px;
}

:deep(.n-tooltip .n-tooltip__content) {
  max-width: 320px;
  line-height: 1.5;
}

:deep(.n-tooltip .n-tooltip__content div) {
  white-space: pre-line;
}

/* 折叠面板样式优化 */
:deep(.n-collapse-item__header) {
  font-weight: 500;
  color: #374151;
}

:deep(.n-collapse-item) {
  --n-title-padding: 16px 0;
}

:deep(.n-base-selection-label) {
  height: 40px;
}

/* 表单行布局 */
.form-row {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.form-item-half {
  flex: 1;
  width: 50%;
}

/* 上游地址行布局 */
.upstream-row {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.upstream-url {
  flex: 1;
}

.upstream-weight {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 0 0 140px;
}

.weight-label {
  font-weight: 500;
  color: #374151;
  white-space: nowrap;
}

.upstream-actions {
  flex: 0 0 32px;
  display: flex;
  justify-content: center;
}

/* 配置项行布局 */
.config-item-row {
  margin-bottom: 12px;
}

.config-item-content {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.config-select {
  flex: 0 0 200px;
}

.config-value {
  flex: 1;
}

.config-actions {
  flex: 0 0 32px;
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .group-form-card {
    width: 100vw !important;
  }

  .group-form {
    width: auto !important;
  }

  .form-row {
    flex-direction: column;
    gap: 0;
  }

  .form-item-half {
    width: 100%;
  }

  .section-title {
    font-size: 0.9rem;
  }

  .upstream-row,
  .config-item-content {
    flex-direction: column;
    gap: 8px;
    align-items: stretch;
  }

  .upstream-weight {
    flex: 1;
    flex-direction: column;
    align-items: flex-start;
  }

  .config-value {
    flex: 1;
  }

  .upstream-actions,
  .config-actions {
    justify-content: flex-end;
  }
}
</style>
