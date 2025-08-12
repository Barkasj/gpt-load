package types

// ConfigManager defines the interface for configuration management
type ConfigManager interface {
	IsMaster() bool
	GetAuthConfig() AuthConfig
	GetCORSConfig() CORSConfig
	GetPerformanceConfig() PerformanceConfig
	GetLogConfig() LogConfig
	GetDatabaseConfig() DatabaseConfig
	GetEffectiveServerConfig() ServerConfig
	GetRedisDSN() string
	Validate() error
	DisplayServerConfig()
	ReloadConfig() error
}

// SystemSettings defines all system configuration items
type SystemSettings struct {
	// Parameter Dasar
	AppUrl                         string `json:"app_url" default:"http://localhost:3001" name:"Alamat Proyek" category:"Parameter Dasar" desc:"URL dasar proyek, digunakan untuk menyusun alamat endpoint grup. Konfigurasi sistem diprioritaskan dibanding variabel lingkungan APP_URL." validate:"required"`
	RequestLogRetentionDays        int    `json:"request_log_retention_days" default:"7" name:"Durasi Simpan Log (hari)" category:"Parameter Dasar" desc:"Jumlah hari log permintaan disimpan di database; 0 berarti tidak dilakukan pembersihan." validate:"required,min=0"`
	RequestLogWriteIntervalMinutes int    `json:"request_log_write_interval_minutes" default:"1" name:"Periode Penulisan Log Tertunda (menit)" category:"Parameter Dasar" desc:"Periode (menit) untuk menulis log permintaan dari cache ke database; 0 berarti tulis real-time." validate:"required,min=0"`
	ProxyKeys                      string `json:"proxy_keys" name:"Kunci Proksi Global" category:"Parameter Dasar" desc:"Kunci proksi global untuk mengakses semua endpoint grup. Pisahkan beberapa kunci dengan koma." validate:"required"`

	// Pengaturan Permintaan
	RequestTimeout        int    `json:"request_timeout" default:"600" name:"Batas Waktu Permintaan (detik)" category:"Pengaturan Permintaan" desc:"Batas waktu (detik) untuk keseluruhan siklus hidup permintaan yang diteruskan." validate:"required,min=1"`
	ConnectTimeout        int    `json:"connect_timeout" default:"15" name:"Batas Waktu Koneksi (detik)" category:"Pengaturan Permintaan" desc:"Batas waktu (detik) untuk membuat koneksi baru ke layanan upstream." validate:"required,min=1"`
	IdleConnTimeout       int    `json:"idle_conn_timeout" default:"120" name:"Batas Waktu Koneksi Menganggur (detik)" category:"Pengaturan Permintaan" desc:"Batas waktu (detik) koneksi menganggur pada klien HTTP." validate:"required,min=1"`
	ResponseHeaderTimeout int    `json:"response_header_timeout" default:"600" name:"Batas Waktu Header Respons (detik)" category:"Pengaturan Permintaan" desc:"Waktu maksimum (detik) menunggu header respons dari layanan upstream." validate:"required,min=1"`
	MaxIdleConns          int    `json:"max_idle_conns" default:"100" name:"Maksimum Koneksi Menganggur" category:"Pengaturan Permintaan" desc:"Jumlah maksimum total koneksi menganggur yang diizinkan di pool koneksi klien HTTP." validate:"required,min=1"`
	MaxIdleConnsPerHost   int    `json:"max_idle_conns_per_host" default:"50" name:"Maksimum Koneksi Menganggur per Host" category:"Pengaturan Permintaan" desc:"Jumlah maksimum koneksi menganggur untuk tiap host upstream pada pool klien HTTP." validate:"required,min=1"`
	ProxyURL              string `json:"proxy_url" name:"Alamat Server Proksi" category:"Pengaturan Permintaan" desc:"Alamat server proksi HTTP/HTTPS global, contoh: http://user:pass@host:port. Jika kosong, gunakan konfigurasi dari variabel lingkungan."`

	// Konfigurasi Kunci
	MaxRetries                   int `json:"max_retries" default:"3" name:"Jumlah Ulangi Maksimum" category:"Konfigurasi Kunci" desc:"Jumlah percobaan ulang maksimum untuk satu permintaan dengan kunci berbeda; 0 berarti tidak mengulang." validate:"required,min=0"`
	BlacklistThreshold           int `json:"blacklist_threshold" default:"3" name:"Ambang Daftar Hitam" category:"Konfigurasi Kunci" desc:"Jumlah kegagalan beruntun sebelum suatu kunci dimasukkan ke daftar hitam; 0 berarti tidak diblokir." validate:"required,min=0"`
	KeyValidationIntervalMinutes int `json:"key_validation_interval_minutes" default:"60" name:"Interval Validasi Kunci (menit)" category:"Konfigurasi Kunci" desc:"Interval bawaan (menit) untuk memvalidasi kunci di latar belakang." validate:"required,min=1"`
	KeyValidationConcurrency     int `json:"key_validation_concurrency" default:"10" name:"Jumlah Paralel Validasi Kunci" category:"Konfigurasi Kunci" desc:"Jumlah paralel saat validasi terjadwal terhadap kunci tidak valid di latar belakang. Jika memakai SQLite atau performa lingkungan kurang, jaga di bawah 20 untuk menghindari ketidakkonsistenan data." validate:"required,min=1"`
	KeyValidationTimeoutSeconds  int `json:"key_validation_timeout_seconds" default:"20" name:"Batas Waktu Validasi Kunci (detik)" category:"Konfigurasi Kunci" desc:"Batas waktu permintaan API (detik) saat memvalidasi satu kunci pada proses latar belakang." validate:"required,min=1"`

	// For cache
	ProxyKeysMap map[string]struct{} `json:"-"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port                    int    `json:"port"`
	Host                    string `json:"host"`
	IsMaster                bool   `json:"is_master"`
	ReadTimeout             int    `json:"read_timeout"`
	WriteTimeout            int    `json:"write_timeout"`
	IdleTimeout             int    `json:"idle_timeout"`
	GracefulShutdownTimeout int    `json:"graceful_shutdown_timeout"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Key string `json:"key"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	Enabled          bool     `json:"enabled"`
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
}

// PerformanceConfig represents performance configuration
type PerformanceConfig struct {
	MaxConcurrentRequests int `json:"max_concurrent_requests"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	EnableFile bool   `json:"enable_file"`
	FilePath   string `json:"file_path"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	DSN string `json:"dsn"`
}

type RetryError struct {
	StatusCode         int    `json:"status_code"`
	ErrorMessage       string `json:"error_message"`
	ParsedErrorMessage string `json:"-"`
	KeyValue           string `json:"key_value"`
	Attempt            int    `json:"attempt"`
	UpstreamAddr       string `json:"-"`
}
