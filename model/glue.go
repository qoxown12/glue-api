package model

import (
	"github.com/gofrs/uuid"
)

// GlueVersion
// @Description Glue의 버전
type GlueVersion struct {
	Mon interface {
	} `json:"mon"`
	Mgr interface {
	} `json:"mgr"`
	Osd interface {
	} `json:"osd"`
	Mds interface {
	} `json:"mds"`
	RbdMirror interface {
	} `json:"rbd-mirror"`
	Rgw interface {
	} `json:"rgw"`
	Overall interface {
	} `json:"overall"`
} //@name GlueVersion

// GlueStatus model info
// @Description Glue의 상태를 나타내는 구조체
type GlueStatus struct {
	Fsid   uuid.UUID `json:"fsid" example:"9980ffe8-4bc1-11ee-9b1f-002481004170" format:"uuid"` //Glue클러스터를 구분하는 ID
	Health struct {
		Status string `json:"status" example:"HEALTH_WARN" format:"string"`
		Checks interface {
		} `json:"checks"`
		Mutes interface{} `json:"mutes"`
	} `json:"health"`
	ElectionEpoch int      `json:"election_epoch" example:"148" format:"uint32"`
	Quorum        []int    `json:"quorum"`
	QuorumNames   []string `json:"quorum_names"`
	QuorumAge     int      `json:"quorum_age"`
	Monmap        struct {
		Epoch             int    `json:"epoch"`
		MinMonReleaseName string `json:"min_mon_release_name"`
		NumMons           int    `json:"num_mons"`
	} `json:"monmap"`
	Osdmap struct {
		Epoch          int `json:"epoch"`
		NumOsds        int `json:"num_osds"`
		NumUpOsds      int `json:"num_up_osds"`
		OsdUpSince     int `json:"osd_up_since"`
		NumInOsds      int `json:"num_in_osds"`
		OsdInSince     int `json:"osd_in_since"`
		NumRemappedPgs int `json:"num_remapped_pgs"`
	} `json:"osdmap"`
	Pgmap struct {
		PgsByState []struct {
			StateName string `json:"state_name"`
			Count     int    `json:"count"`
		} `json:"pgs_by_state"`
		NumPgs        int   `json:"num_pgs"`
		NumPools      int   `json:"num_pools"`
		NumObjects    int   `json:"num_objects"`
		DataBytes     int64 `json:"data_bytes"`
		BytesUsed     int64 `json:"bytes_used"`
		BytesAvail    int64 `json:"bytes_avail"`
		BytesTotal    int64 `json:"bytes_total"`
		ReadBytesSec  int   `json:"read_bytes_sec"`
		WriteBytesSec int   `json:"write_bytes_sec"`
		ReadOpPerSec  int   `json:"read_op_per_sec"`
		WriteOpPerSec int   `json:"write_op_per_sec"`
	} `json:"pgmap"`
	Fsmap struct {
		Epoch     int           `json:"epoch"`
		ByRank    []interface{} `json:"by_rank"`
		UpStandby int           `json:"up:standby"`
	} `json:"fsmap"`
	Mgrmap struct {
		Available   bool     `json:"available"`
		NumStandbys int      `json:"num_standbys"`
		Modules     []string `json:"modules"`
		Services    struct {
			Dashboard  string `json:"dashboard"`
			Prometheus string `json:"prometheus"`
		} `json:"services"`
	} `json:"mgrmap"`
	Servicemap struct {
		Epoch    int         `json:"epoch"`
		Modified string      `json:"modified"`
		Services interface{} `json:"services"`
	} `json:"servicemap"`
	ProgressEvents struct {
	} `json:"progress_events"`
} // @name GlueStatus

type GluePools interface{} // @name GluePools

type ServiceLs interface{} // @name ServiceLs

type InfoImage struct {
	Name            string   `json:"name"`
	ID              string   `json:"id"`
	Size            int64    `json:"size"`
	Objects         int      `json:"objects"`
	Order           int      `json:"order"`
	ObjectSize      int      `json:"object_size"`
	SnapshotCount   int      `json:"snapshot_count"`
	BlockNamePrefix string   `json:"block_name_prefix"`
	Format          int      `json:"format"`
	Features        []string `json:"features"`
	OpFeatures      []any    `json:"op_features"`
	Flags           []any    `json:"flags"`
	CreateTimestamp string   `json:"create_timestamp"`
	AccessTimestamp string   `json:"access_timestamp"`
	ModifyTimestamp string   `json:"modify_timestamp"`
} // @name InfoImage
type HostList []struct {
	Ip_Address string `json:"ip_address"`
	Addr       string `json:"addr"`
	Hostname   string `json:"hostname"`
	Status     string `json:"status"`
} // @name HostList

type Images []struct {
	Image  string `json:"image"`
	ID     string `json:"id"`
	Size   int64  `json:"size"`
	Format int    `json:"format"`
} // @name Images

type ImageCommon interface{} // @name ImageCommon

type PwEncryption []struct {
	Pass_Word string `json:"pass_word"`
} // @name PwEncryption