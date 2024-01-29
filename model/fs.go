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
		Caps  int    `json:"caps,omitempty"`
		Dirs  int    `json:"dirs,omitempty"`
		DNS   int    `json:"dns,omitempty"`
		Inos  int    `json:"inos,omitempty"`
		Name  string `json:"name"`
		Rank  int    `json:"rank,omitempty"`
		Rate  int    `json:"rate,omitempty"`
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
		Compat              struct {
			Compat struct {
			} `json:"compat"`
			RoCompat struct {
			} `json:"ro_compat"`
			Incompat struct {
				Feature1  string `json:"feature_1"`
				Feature2  string `json:"feature_2"`
				Feature3  string `json:"feature_3"`
				Feature4  string `json:"feature_4"`
				Feature5  string `json:"feature_5"`
				Feature6  string `json:"feature_6"`
				Feature7  string `json:"feature_7"`
				Feature8  string `json:"feature_8"`
				Feature9  string `json:"feature_9"`
				Feature10 string `json:"feature_10"`
			} `json:"incompat"`
		} `json:"compat"`
		MaxMds int   `json:"max_mds"`
		In     []int `json:"in"`
		Up     struct {
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

// FsList model info
// @Description GlueFS 리스트 구조체
type FsList []struct {
	Name           string   `json:"name"`
	MetadataPool   string   `json:"metadata_pool"`
	MetadataPoolID int      `json:"metadata_pool_id"`
	DataPoolIds    []int    `json:"data_pool_ids"`
	DataPools      []string `json:"data_pools"`
} //@name FsList
