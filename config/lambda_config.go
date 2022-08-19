package config

type lambdaConfig struct {
	SchemaLoadMode        string `json:"schema_load_mode"`
	ModuleName            string `json:"module_name"`
	ProjectKey            string `json:"project_key"`
	LambdaMainServicePath string `json:"lambda_main_service_path"`
	MicroserviceDev       bool   `json:"microservice_dev"`
	Theme                 string `json:"theme"`
	Domain                string `json:"domain"`
	Title                 string `json:"title"`
	SubTitle              string `json:"subTitle"`
	Copyright             string `json:"copyright"`
	Favicon               string `json:"favicon"`
	Bg                    string `json:"bg"`
	Logo                  string `json:"logo"`
	LogoDark              string `json:"logo_dark"`
	LogoText              string `json:"logoText"`
	SuperURL              string `json:"super_url"`
	AppURL                string `json:"app_url"`
	HasLanguage           bool   `json:"has_language"`
	WithCrudLog           bool   `json:"withCrudLog"`
	KrudPublic            bool   `json:"krud_public"`
	ControlPanel          struct {
		LogoLight    string   `json:"logoLight"`
		LogoDark     string   `json:"logoDark"`
		BrandBtnURL  string   `json:"brandBtnUrl"`
		ThemeMode    string   `json:"themeMode"`
		PrimaryColor string   `json:"primaryColor"`
		ThemeColors  []string `json:"themeColors"`
		ExtraStyles  []string `json:"extraStyles"`
		ExtraScripts []string `json:"extraScripts"`
	} `json:"controlPanel"`
	Languages []struct {
		Label string `json:"label"`
		Code  string `json:"code"`
	} `json:"languages"`
	DefaultLanguage string `json:"default_language"`
	AdminRoles      []int64
	RoleRedirects   []struct {
		RoleID int64  `json:"role_id"`
		URL    string `json:"url"`
	} `json:"role-redirects"`
	UserDataFields         []string `json:"user_data_fields"`
	DataFormCustomElements []struct {
		Element string `json:"element"`
	} `json:"data_form_custom_elements"`
	DataGridCustomElements []struct {
		Element string `json:"element"`
	} `json:"data_grid_custom_elements"`
	PasswordResetTimeOut int                    `json:"password_reset_time_out"`
	StaticWords          map[string]interface{} `json:"static_words"`
	Notify               struct {
		FirebaseConfig struct {
			APIKey            string `json:"apiKey"`
			PublicKey         string `json:"publicKey"`
			AuthDomain        string `json:"authDomain"`
			DatabaseURL       string `json:"databaseURL"`
			ProjectID         string `json:"projectId"`
			StorageBucket     string `json:"storageBucket"`
			MessagingSenderID string `json:"messagingSenderId"`
			AppID             string `json:"appId"`
			MeasurementID     string `json:"measurementId"`
		} `json:"firebaseConfig"`
		ServerKey string `json:"serverKey"`
		Sound     string `json:"sound"`
		Icon      string `json:"icon"`
	} `json:"notify"`
}
type LambdaConfigFile struct {
	SchemaLoadMode string `json:"schema_load_mode"`
	ProjectKey     string `json:"project_key"`
	Theme          string `json:"theme"`
	Domain         string `json:"domain"`
	Title          string `json:"title"`
	SubTitle       string `json:"subTitle"`
	Copyright      string `json:"copyright"`
	Favicon        string `json:"favicon"`
	Bg             string `json:"bg"`
	Logo           string `json:"logo"`
	LogoText       string `json:"logoText"`
	SuperURL       string `json:"super_url"`
	AppURL         string `json:"app_url"`
	HasLanguage    bool   `json:"has_language"`
	WithCrudLog    bool   `json:"withCrudLog"`
	KrudPublic     bool   `json:"krud_public"`
	ControlPanel   struct {
		LogoLight    string   `json:"logoLight"`
		LogoDark     string   `json:"logoDark"`
		BrandBtnURL  string   `json:"brandBtnUrl"`
		ThemeMode    string   `json:"themeMode"`
		PrimaryColor string   `json:"primaryColor"`
		ThemeColors  []string `json:"themeColors"`
		ExtraStyles  []string `json:"extraStyles"`
		ExtraScripts []string `json:"extraScripts"`
	} `json:"controlPanel"`
	Languages []struct {
		Label string `json:"label"`
		Code  string `json:"code"`
	} `json:"languages"`
	DefaultLanguage string `json:"default_language"`
	RoleRedirects   []struct {
		RoleID int64  `json:"role_id"`
		URL    string `json:"url"`
	} `json:"role-redirects"`
	UserDataFields         []string `json:"user_data_fields"`
	DataFormCustomElements []struct {
		Element string `json:"element"`
	} `json:"data_form_custom_elements"`
	DataGridCustomElements []struct {
		Element string `json:"element"`
	} `json:"data_grid_custom_elements"`
	PasswordResetTimeOut int                    `json:"password_reset_time_out"`
	StaticWords          map[string]interface{} `json:"static_words"`
	Notify               struct {
		FirebaseConfig struct {
			APIKey            string `json:"apiKey"`
			PublicKey         string `json:"publicKey"`
			AuthDomain        string `json:"authDomain"`
			DatabaseURL       string `json:"databaseURL"`
			ProjectID         string `json:"projectId"`
			StorageBucket     string `json:"storageBucket"`
			MessagingSenderID string `json:"messagingSenderId"`
			AppID             string `json:"appId"`
			MeasurementID     string `json:"measurementId"`
		} `json:"firebaseConfig"`
		ServerKey string `json:"serverKey"`
		Sound     string `json:"sound"`
		Icon      string `json:"icon"`
	} `json:"notify"`
}
