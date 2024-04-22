package model

type RgwDaemon []struct {
	ID             string `json:"id"`
	ServiceMapID   string `json:"service_map_id"`
	Version        string `json:"version"`
	ServerHostname string `json:"server_hostname"`
	RealmName      string `json:"realm_name"`
	ZonegroupName  string `json:"zonegroup_name"`
	ZoneName       string `json:"zone_name"`
	Default        bool   `json:"default"`
	Port           int    `json:"port"`
} // @name RgwDaemon
type RgwUserList []string

type RgwUserStat struct {
	Stats struct {
		Size         int `json:"size"`
		SizeActual   int `json:"size_actual"`
		SizeKb       int `json:"size_kb"`
		SizeKbActual int `json:"size_kb_actual"`
		NumObjects   int `json:"num_objects"`
	} `json:"stats"`
}
type RgwUserInfo struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Suspended   int    `json:"suspended"`
	MaxBuckets  int    `json:"max_buckets"`
	Subusers    []any  `json:"subusers"`
	Keys        []struct {
		User      string `json:"user"`
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
	} `json:"keys"`
	SwiftKeys           []any  `json:"swift_keys"`
	Caps                []any  `json:"caps"`
	OpMask              string `json:"op_mask"`
	System              bool   `json:"system"`
	DefaultPlacement    string `json:"default_placement"`
	DefaultStorageClass string `json:"default_storage_class"`
	PlacementTags       []any  `json:"placement_tags"`
	BucketQuota         struct {
		Enabled    bool `json:"enabled"`
		CheckOnRaw bool `json:"check_on_raw"`
		MaxSize    int  `json:"max_size"`
		MaxSizeKb  int  `json:"max_size_kb"`
		MaxObjects int  `json:"max_objects"`
	} `json:"bucket_quota"`
	UserQuota struct {
		Enabled    bool `json:"enabled"`
		CheckOnRaw bool `json:"check_on_raw"`
		MaxSize    int  `json:"max_size"`
		MaxSizeKb  int  `json:"max_size_kb"`
		MaxObjects int  `json:"max_objects"`
	} `json:"user_quota"`
	TempURLKeys []any  `json:"temp_url_keys"`
	Type        string `json:"type"`
	MfaIds      []any  `json:"mfa_ids"`
} // @name RgwUserInfo

type RgwUserInfoAndStat struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Suspended   int    `json:"suspended"`
	MaxBuckets  int    `json:"max_buckets"`
	Subusers    []any  `json:"subusers"`
	Keys        []struct {
		User      string `json:"user"`
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
	} `json:"keys"`
	SwiftKeys           []any  `json:"swift_keys"`
	Caps                []any  `json:"caps"`
	OpMask              string `json:"op_mask"`
	System              bool   `json:"system"`
	DefaultPlacement    string `json:"default_placement"`
	DefaultStorageClass string `json:"default_storage_class"`
	PlacementTags       []any  `json:"placement_tags"`
	BucketQuota         struct {
		Enabled    bool `json:"enabled"`
		CheckOnRaw bool `json:"check_on_raw"`
		MaxSize    int  `json:"max_size"`
		MaxSizeKb  int  `json:"max_size_kb"`
		MaxObjects int  `json:"max_objects"`
	} `json:"bucket_quota"`
	UserQuota struct {
		Enabled    bool `json:"enabled"`
		CheckOnRaw bool `json:"check_on_raw"`
		MaxSize    int  `json:"max_size"`
		MaxSizeKb  int  `json:"max_size_kb"`
		MaxObjects int  `json:"max_objects"`
	} `json:"user_quota"`
	TempURLKeys []any       `json:"temp_url_keys"`
	Type        string      `json:"type"`
	MfaIds      []any       `json:"mfa_ids"`
	Stats       RgwUserStat `json:"stats"`
} // @name RgwUserInfoAndStat

type RgwUpdate struct {
	Service_type string             `yaml:"service_type"`
	Service_id   string             `yaml:"service_id"`
	Placement    RgwUpdatePlacement `yaml:"placement"`
	Spec         RgwUpdateSpec      `yaml:"spec"`
}
type RgwUpdatePlacement struct {
	Hosts []string `yaml:"hosts"`
}
type RgwUpdateSpec struct {
	Rgw_realm         string `yaml:"rgw_realm"`
	Rgw_zonegroup     string `yaml:"rgw_zonegroup"`
	Rgw_zone          string `yaml:"rgw_zone"`
	Rgw_frontend_port string `yaml:"rgw_frontend_port"`
}

type RgwBucketCreate struct {
	Bucket                     string `json:"bucket"`
	Uid                        string `json:"uid"`
	Lock_enabled               string `json:"lock_enabled"`
	Lock_mode                  string `json:"lock_mode"`
	Lock_retention_period_days string `json:"lock_retention_period_days"`
}

type RGwCommon interface{} //@name RgwCommon

type RgwBucketUpdate struct {
	Bucket_id                  string `json:"bucket_id"`
	Uid                        string `json:"uid"`
	Versioning_state           string `json:"versioning_state"`
	Lock_mode                  string `json:"lock_mode"`
	Lock_retention_period_days string `json:"lock_retention_period_days"`
}
