package model

// FsStatus model info
// @Description GlueFS의 상태 구조체
type FsStatus struct {
	Clients []struct {
		Clients int    `json:"clients"`
		Fs      string `json:"fs"`
	} `json:"clients"`
	MdsVersion []struct {
		Daemon  []string `json:"daemon"`
		Version string   `json:"version"`
	} `json:"mds_version"`
	Mdsmap []struct {
		Name  string `json:"name"`
		State string `json:"state"`
	} `json:"mdsmap"`
	Pools []struct {
		Avail int64  `json:"avail"`
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Used  int    `json:"used"`
	} `json:"pools"`
} //@name FsStatus
type FsUpdate struct {
	Service_type string      `yaml:"service_type"`
	Service_id   string      `yaml:"service_id"`
	Placement    FsPlacement `yaml:"placement"`
} //@name FsUpdate
type FsPlacement struct {
	Hosts []string `yaml:"hosts"`
}

// FsList model info
// @Description GlueFS 리스트 구조체
type FsList []struct {
	Name           string   `json:"name"`
	MetadataPool   string   `json:"metadata_pool"`
	MetadataPoolID int      `json:"metadata_pool_id"`
	DataPoolIds    []int    `json:"data_pool_ids"`
	DataPools      []string `json:"data_pools"`
} //@name FsList
type FsSum struct {
	FsStatus FsStatus `json:"status"`
	FsList   FsList   `json:"list"`
}

// FsGetInfo model info
// @Description GlueFS의 상세정보 구조체
type FsGetInfo struct {
	Mdsmap struct {
		Epoch      int `json:"epoch"`
		Flags      int `json:"flags"`
		FlagsState struct {
			Joinable            bool `json:"joinable"`
			AllowSnaps          bool `json:"allow_snaps"`
			AllowMultimdsSnaps  bool `json:"allow_multimds_snaps"`
			AllowStandbyReplay  bool `json:"allow_standby_replay"`
			RefuseClientSession bool `json:"refuse_client_session"`
		} `json:"flags_state"`
		EverAllowedFeatures       int    `json:"ever_allowed_features"`
		ExplicitlyAllowedFeatures int    `json:"explicitly_allowed_features"`
		Created                   string `json:"created"`
		Modified                  string `json:"modified"`
		Tableserver               int    `json:"tableserver"`
		Root                      int    `json:"root"`
		SessionTimeout            int    `json:"session_timeout"`
		SessionAutoclose          int    `json:"session_autoclose"`
		RequiredClientFeatures    struct {
		} `json:"required_client_features"`
		MaxFileSize         int64 `json:"max_file_size"`
		LastFailure         int   `json:"last_failure"`
		LastFailureOsdEpoch int   `json:"last_failure_osd_epoch"`
		MaxMds              int   `json:"max_mds"`
		In                  []int `json:"in"`
		Up                  struct {
			Mds0 int `json:"mds_0"`
		} `json:"up"`
		Failed  []any `json:"failed"`
		Damaged []any `json:"damaged"`
		Stopped []any `json:"stopped"`
		Info    interface {
		} `json:"info"`
		DataPools          []int  `json:"data_pools"`
		MetadataPool       int    `json:"metadata_pool"`
		Enabled            bool   `json:"enabled"`
		FsName             string `json:"fs_name"`
		Balancer           string `json:"balancer"`
		BalRankMask        string `json:"bal_rank_mask"`
		StandbyCountWanted int    `json:"standby_count_wanted"`
	} `json:"mdsmap"`
	ID int `json:"id"`
} //@name FsGetInfo
type CephHost []struct {
	Addr     string `json:"addr"`
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
}
