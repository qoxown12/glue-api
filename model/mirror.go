package model

/*
type MirrorList struct {
	AbleModel
}
*/

type MirrorImageItem struct {
	Interval  string `json:"interval"`
	StartTime string `json:"start_time"`
} //@name MirrorImageItem

type MirrorImage struct {
	Pool      string            `json:"pool"`
	Namespace string            `json:"namespace"`
	Image     string            `json:"image"`
	Items     []MirrorImageItem `json:"items"`
} //@name MirrorImage

type MirrorList struct {
	AbleModel
	Local  []MirrorImage `json:"Local"`
	Remote []MirrorImage `json:"Remote"`
} //@name MirrorList

type MirrorSetup struct {
	AbleModel
	LocalClusterName  string      `json:"localClusterName"`  //미러링 상태
	RemoteClusterName string      `json:"remoteClusterName"` //미러링 상태
	Host              string      `json:"host"`
	PrivateKeyFile    interface{} `json:"privateKeyFile" type:"file"`
	MirrorPool        string      `json:"mirrorPool"`
	LocalToken        string      `json:"localToken"`
	RemoteToken       string      `json:"remoteToken"`
} //@name MirrorSetup

type MirrorConf struct {
	Name            string `json:"name"`            //미러링 상태
	ClusterFileName string `json:"clusterFileName"` //미러링 데몬 상태
	KeyFileName     string `json:"keyFileName"`     //이미지 상태
	ClusterName     string `json:"clusterName"`     //이미지 상세
	Mode            string `json:"mode"`
	SiteName        string `json:"site_name"`
	Peers           []struct {
		Uuid       string `json:"uuid"`
		Direction  string `json:"direction"`
		SiteName   string `json:"site_name"`
		MirrorUuid string `json:"mirror_uuid"`
		ClientName string `json:"client_name"`
		Key        string `json:"key"`
		MonHost    string `json:"mon_host"`
	} `json:"peers"`
} //@name MirrorConf

type MirrorStatus struct {
	AbleModel
	Health       string      `json:"health"`        //미러링 상태
	DaemonHealth string      `json:"daemon_health"` //미러링 데몬 상태
	ImageHealth  string      `json:"image_health"`  //이미지 상태
	States       interface{} `json:"states"`        //이미지 상세
} //@name MirrorStatus

type MirrorInfo struct {
	Mode     string `json:"mode"`
	SiteName string `json:"site_name"`
	Peers    []struct {
		Uuid       string `json:"uuid"`
		Direction  string `json:"direction"`
		SiteName   string `json:"site_name"`
		MirrorUuid string `json:"mirror_uuid"`
		ClientName string `json:"client_name"`
		Key        string `json:"key"`
		MonHost    string `json:"mon_host"`
	} `json:"peers"`
} // @name MirrorInfo

type MirrorToken struct {
	Fsid     string `json:"fsid"`
	ClientId string `json:"client_id"`
	Key      string `json:"key"`
	MonHost  string `json:"mon_host"`
} // @name MirrorToken

type AuthKey struct {
	Key string `json:"key"`
} //@name AuthKey

type ImageList struct {
	AbleModel
	Images []Image `json:"images"`
} // @name Images

type Image struct {
	Image  string `json:"image"`
	Id     string `json:"id"`
	Size   int64  `json:"size"`
	Format int    `json:"format"`
} //@name Image
