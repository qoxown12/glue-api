package model

type SubVolumeGroupInfo struct {
	Atime      string   `json:"atime"`
	BytesPcent string   `json:"bytes_pcent"`
	BytesQuota int64    `json:"bytes_quota"`
	BytesUsed  int      `json:"bytes_used"`
	CreatedAt  string   `json:"created_at"`
	Ctime      string   `json:"ctime"`
	DataPool   string   `json:"data_pool"`
	Gid        int      `json:"gid"`
	Mode       int      `json:"mode"`
	MonAddrs   []string `json:"mon_addrs"`
	Mtime      string   `json:"mtime"`
	UID        int      `json:"uid"`
} //@name SubVolumeGroupInfo

type SubVolumeAllLs []struct {
	Name string `json:"name"`
} //@name SubVolumeAllLs
type SubVolumeAllSnapLs []struct {
	Name string `json:"name"`
} //@name SubVolumeAllSnapLs

type SubVolumeAllSnap interface{} //@name SubVolumeAllSnap

type SubVolumeGroupList struct {
	Name     string `json:"name"`
	Info     SubVolumeGroupInfo
	Path     string   `json:"path"`
	Snapshot []string `json:"snapshot"`
} //@name SubVolumeGroupList

type SubVolumeInfo struct {
	Atime         string   `json:"atime"`
	BytesPcent    string   `json:"bytes_pcent"`
	BytesQuota    int64    `json:"bytes_quota"`
	BytesUsed     int      `json:"bytes_used"`
	CreatedAt     string   `json:"created_at"`
	Ctime         string   `json:"ctime"`
	DataPool      string   `json:"data_pool"`
	Features      []string `json:"features"`
	Gid           int      `json:"gid"`
	Mode          int      `json:"mode"`
	MonAddrs      []string `json:"mon_addrs"`
	Mtime         string   `json:"mtime"`
	Path          string   `json:"path"`
	PoolNamespace string   `json:"pool_namespace"`
	State         string   `json:"state"`
	Type          string   `json:"type"`
	UID           int      `json:"uid"`
} // @name SubVolumeInfo
type SubVolumeList struct {
	Name     string `json:"name"`
	Info     SubVolumeInfo
	SnapShot []string `json:"snapshot"`
} // @name SubVolumeList
