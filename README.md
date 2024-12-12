# Glue-API

Glue의 기능을 제어하기 위한 REST API 입니다.

본 문서는 swagger로 작성된 API 목록을 [swagger-markdown-ui](https://swagger-markdown-ui.netlify.app)를 사용해 작성된 README입니다.

## API 목록

| Method | API                                                                    |       진행도       | 비고                        |
| ------ | ---------------------------------------------------------------------- | :----------------: | --------------------------- |
| GET    | [api/v1/version](#version)                                             | :white_check_mark: | Version                     |
| GET    | [api/v1/glue](#apiv1glue)                                              | :white_check_mark: | GlueStatus                  |
| GET    | [api/v1/glue/version](#apiv1glueversion)                               | :white_check_mark: | GlueVersion                 |
| GET    | [api/v1/glue/pool](#apiv1gluepool)                                     | :white_check_mark: | ListPools                   |
| DELETE | [api/v1/glue/pool/:poolname](#apiv1gluepool-)                          | :white_check_mark: | PoolDelete                  |
| GET    | [api/v1/glue/rbd/:poolname](#apiv1gluerbd-)                            | :white_check_mark: | ListImages                  |
| GET    | [api/v1/gluefs](#apiv1gluefs)                                          | :white_check_mark: | FsStatus                    |
| GET    | [api/v1/gluefs/info/:fs_name](#apiv1gluefsinfo-)                       | :white_check_mark: | FsGetInfo                   |
| GET    | [api/v1/gluefs/list](#apiv1gluefslist)                                 | :white_check_mark: | FsList                      |
| POST   | [api/v1/gluefs/:fs_name](#apiv1gluefs-)                                | :white_check_mark: | FsCreate                    |
| DELETE | [api/v1/gluefs/:fs_name](#apiv1gluefs-)                                | :white_check_mark: | FsDelete                    |
| GET    | [api/v1/mirror](#apiv1mirror)                                          | :white_check_mark: | MirrorStatus                |
| POST   | [api/v1/mirror](#apiv1mirror)                                          | :white_check_mark: | MirrorSetup                 |
| PUT    | [api/v1/mirror](#apiv1mirror)                                          | :white_check_mark: | MirrorUpdate                |
| DELETE | [api/v1/mirror](#apiv1mirror)                                          | :white_check_mark: | MirrorDelete                |
| POST   | [api/v1/mirror/:pool](#apiv1mirror)                                    | :white_check_mark: | MirrorPoolEnable            |
| DELETE | [api/v1/mirror/:pool](#apiv1mirror)                                    | :white_check_mark: | MirrorPoolDisable           |
| DELETE | [api/v1/mirror/garbage](#apiv1mirror)                                  | :white_check_mark: | MirrorDeleteGarbage         |
| GET    | [api/v1/mirror/image/:pool](#apiv1mirrorimage)                         | :white_check_mark: | MirrorImageList             |
| GET    | [api/v1/mirror/image/info/:pool/:imageName]()                          | :white_check_mark: | MirrorImageParentInfo       |
| GET    | [api/v1/mirror/image/status/:pool/:imageName]()                        | :white_check_mark: | MirrorImageStatus           |
| POST   | [api/v1/mirror/image/promote/:pool/:imageName]()                       | :white_check_mark: | MirrorImagePromote          |
| POST   | [api/v1/mirror/image/promote/peer/:pool/:imageName]()                  | :white_check_mark: | MirrorImagePromotePeer      |
| DELETE | [api/v1/mirror/image/demote/:pool/:imageName]()                        | :white_check_mark: | MirrorImageDemote           |
| DELETE | [api/v1/mirror/image/demote/peer/:pool/:imageName]()                   | :white_check_mark: | MirrorImageDemotePeer       |
| PUT    | [api/v1/mirror/image/resync/:pool/:imageName]()                        | :white_check_mark: | MirrorImageResync           |
| PUT    | [api/v1/mirror/image/resync/peer/:pool/:imageName]()                   | :white_check_mark: | MirrorImageResyncPeer       |
| GET    | [api/v1/mirror/image/:pool/:imageName]()                               | :white_check_mark: | MirrorImageInfo             |
| DELETE | [api/v1/mirror/image/:pool/:imageName]()                               | :white_check_mark: | MirrorImageScheduleDelete   |
| POST   | [api/v1/mirror/image/:pool/:imageName/:hostName:/vmName]()             | :white_check_mark: | MirrorImageScheduleSetup    |
| POST   | [api/v1/mirror/image/snapshot/:pool/:vmName]()                         | :white_check_mark: | MirrorImageSnap             |
| GET    | [api/v1/nfs](#apiv1nfs)                                                | :white_check_mark: | NfsClusterLs,NfsClusterInfo |
| DELETE | [api/v1/nfs/:cluster_id](#apiv1nfscluster)                             | :white_check_mark: | NfsClusterDelete            |
| POST   | [api/v1/nfs/:cluster_id/:port](#apiv1nfscluster-)                      | :white_check_mark: | NfsClusterCreate            |
| GET    | [api/v1/nfs/export](#apiv1nfsexport)                                   | :white_check_mark: | NfsExportDetailed           |
| PUT    | [api/v1/nfs/export/:cluster_id](#apiv1nfsexport-)                      | :white_check_mark: | NfsExportUpdate             |
| POST   | [api/v1/nfs/export/:cluster_id](#apiv1nfsexport-)                      | :white_check_mark: | NfsExportCreate             |
| DELETE | [api/v1/nfs/export/:cluster_id/:export_id](#apiv1nfsexport--)          | :white_check_mark: | NfsExportDelete             |
| POST   | [api/v1/iscsi](#apiv1iscsi)                                            | :white_check_mark: | IscsiServiceCreate          |
| PUT    | [api/v1/iscsi/discovery](#apiv1iscsidiscovery)                         | :white_check_mark: | IscsiUpdateDiscoveryAuth    |
| GET    | [api/v1/iscsi/discovery](#apiv1iscsidiscovery)                         | :white_check_mark: | IscsiGetDiscoveryAuth       |
| GET    | [api/v1/iscsi/target](#apiv1iscsitarget)                               | :white_check_mark: | IscsiTargetList             |
| POST   | [api/v1/iscsi/target](#apiv1iscsitarget)                               | :white_check_mark: | IscsiTargetCreate           |
| DELETE | [api/v1/iscsi/target](#apiv1iscsitarget)                               | :white_check_mark: | IscsiTargetDelete           |
| PUT    | [api/v1/iscsi/target](#apiv1iscsitarget)                               | :white_check_mark: | IscsiTargetUpdate           |
| GET    | [api/v1/smb](#apiv1smb)                                                | :white_check_mark: | SmbStatus                   |
| POST   | [api/v1/smb](#apiv1smb)                                                | :white_check_mark: | SmbCreate                   |
| DELETE | [api/v1/smb](#apiv1smb)                                                | :white_check_mark: | SmbDelete                   |
| POST   | [api/v1/smb/user](#apiv1smbuser)                                       | :white_check_mark: | SmbUserCreate               |
| PUT    | [api/v1/smb/user](#apiv1smbuser)                                       | :white_check_mark: | SmbUserUpdate               |
| DELETE | [api/v1/smb/user](#apiv1smbuser)                                       | :white_check_mark: | SmbUserDelete               |
| ANY    | swagger/index.html                                                     | :white_check_mark: |                             |

### /api/v1/glue

#### GET

##### Summary:

Show Status of Glue

##### Description:

Glue 의 상태값을 보여줍니다.

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [GlueStatus](#GlueStatus)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/glue/pool

#### GET

##### Summary:

List Pools of Glue

##### Description:

Glue 의 스토리지 풀 목록을 보여줍니다..

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [GlueVersion](#GlueVersion)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/glue/rbd/{pool_name}

#### GET

##### Summary:

List Images of Pool Glue

##### Description:

Glue 스토리지 풀의 이미지 목록을 보여줍니다..

##### Parameters

| Name      | Located in | Description | Required | Schema |
| --------- | ---------- | ----------- | -------- | ------ |
| pool_name | path       | Pool Name   | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [ListImages](#ListImages)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/glue/pool/{pool_name}

#### DELETE

##### Summary:

Delete of Pool Glue

##### Description:

Glue 스토리지 풀을 삭제합니다..

##### Parameters

| Name      | Located in | Description | Required | Schema |
| --------- | ---------- | ----------- | -------- | ------ |
| pool_name | path       | Pool Name   | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [PoolDelete](#PoolDelete)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/gluefs

#### GET

##### Summary:

Show Status of GlueFS

##### Description:

GlueFS의 상태값을 보여줍니다..

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [FsStatus](#FsStatus)                                     |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/gluefs/info/{fs_name}

#### GET

##### Summary:

Show Info of GlueFS

##### Description:

GlueFS의 상세정보를 보여줍니다..

##### Parameters

| Name    | Located in | Description  | Required | Schema |
| ------- | ---------- | ------------ | -------- | ------ |
| fs_name | path       | Glue FS Name | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [FsGetInfo](#FsGetInfo)                                   |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/gluefs/list

#### GET

##### Summary:

Show List of GlueFS

##### Description:

GlueFS의 리스트를 보여줍니다..

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [FsList](#FsList)                                         |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/gluefs/{fs_name}

#### POST

##### Summary:

Create of GlueFS

##### Description:

GlueFS를 생성합니다..

##### Parameters

| Name    | Located in | Description  | Required | Schema |
| ------- | ---------- | ------------ | -------- | ------ |
| fs_name | path       | Glue FS Name | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [FsCreate](#FsCreate)                                     |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### DELETE

##### Summary:

Delete of GlueFS

##### Description:

GlueFS를 삭제합니다..

##### Parameters

| Name    | Located in | Description  | Required | Schema |
| ------- | ---------- | ------------ | -------- | ------ |
| fs_name | path       | Glue FS Name | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [FsDelete](#FsDelete)                                     |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/glue/version

#### GET

##### Summary:

Show Versions of Glue

##### Description:

Glue 의 버전을 보여줍니다.

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [GlueVersion](#GlueVersion)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror

#### GET

##### Summary:

Show Status of Mirror

##### Description:

Glue 의 미러링 상태를 보여줍니다.

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorStatus](#MirrorStatus)                             |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### POST

##### Summary:

Setup Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터를 설정합니다.

##### Parameters

| Name              | Located in | Description                 | Required | Schema |
| ----------------- | ---------- | --------------------------- | -------- | ------ |
| localClusterName  | formData   | Local Cluster Name          | Yes      | string |
| remoteClusterName | formData   | Remote Cluster Name         | Yes      | string |
| host              | formData   | Remote Cluster Host Address | Yes      | string |
| privateKeyFile    | formData   | Remote Cluster PrivateKey   | Yes      | file   |
| mirrorPool        | formData   | Pool Name for Mirroring     | Yes      | string |
| moldUrl           | formData   | Mold Url                    | Yes      | string |
| moldApiKey        | formData   | Mold Api Key                | Yes      | string |
| moldSecretKey     | formData   | Mold Secret Key             | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorSetup](#MirrorSetup)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### PUT

##### Summary:

Put Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터 설정을 변경합니다.

##### Parameters

| Name              | Located in | Description                 | Required | Schema |
| ----------------- | ---------- | --------------------------- | -------- | ------ |
| interval          | formData   | Mirroring Schedule Interval | Yes      | string |
| moldUrl           | formData   | Mold API request URL        | Yes      | string |
| moldApiKey        | formData   | Mold Admin Api Key          | Yes      | string |
| moldSecretKey     | formData   | Mold Admin Secret Key       | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorSetup](#MirrorSetup)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### DELETE

##### Summary:

Delete Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터를 제거합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| host           | formData   | Remote Cluster Host Address | Yes      | string |
| privateKeyFile | formData   | Remote Cluster PrivateKey   | Yes      | file   |
| mirrorPool     | formData   | Pool Name for Mirroring     | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorSetup](#MirrorSetup)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/{mirrorPool}

#### POST

##### Summary:

Enable Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터를 활성화합니다.

##### Parameters

| Name              | Located in | Description                 | Required | Schema |
| ----------------- | ---------- | --------------------------- | -------- | ------ |
| localClusterName  | formData   | Local Cluster Name          | Yes      | string |
| remoteClusterName | formData   | Remote Cluster Name         | Yes      | string |
| host              | formData   | Remote Cluster Host Address | Yes      | string |
| privateKeyFile    | formData   | Remote Cluster PrivateKey   | Yes      | file   |
| mirrorPool        | formData   | Pool Name for Mirroring     | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorSetup](#MirrorSetup)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### DELETE

##### Summary:

Disable Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터를 비활성화합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| host           | formData   | Remote Cluster Host Address | Yes      | string |
| privateKeyFile | formData   | Remote Cluster PrivateKey   | Yes      | file   |
| mirrorPool     | formData   | Pool Name for Mirroring     | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorSetup](#MirrorSetup)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/garbage

#### DELETE

##### Summary:

Delete Garbage Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터 가비지를 삭제합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorSetup](#MirrorSetup)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/{pool}

#### GET

##### Summary:

Show List of Mirrored Snapshot

##### Description:

미러링중인 이미지의 목록과 상태를 보여줍니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/info/{pool}/{imageName}

#### GET

##### Summary:

Show Mirroring Image Parent Info

##### Description:

Glue 의 이미지에 미러링 정보를 확인합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |
| imageName      | path       | Image Name for Mirroring    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/status/{pool}/{imageName}

#### GET

##### Summary:

Show Mirroring Image Status

##### Description:

Glue 의 이미지에 미러링상태를 확인합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |
| imageName      | path       | Image Name for Mirroring    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/promote/{pool}/{imageName}

#### POST

##### Summary:

Promote Image Mirroring

##### Description:

Glue 의 이미지를 Promote 합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |
| imageName      | path       | Image Name for Mirroring    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/promote/peer/{pool}/{imageName}

#### POST

##### Summary:

Promote Peer Image Mirroring

##### Description:

Peer Glue 의 이미지를 Promote 합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |
| imageName      | path       | Image Name for Mirroring    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/demote/{pool}/{imageName}

#### DELETE

##### Summary:

Demote Image Mirroring

##### Description:

Glue 의 이미지를 Demote 합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |
| imageName      | path       | Image Name for Mirroring    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/demote/peer/{pool}/{imageName}

#### DELETE

##### Summary:

Demote Peer Image Mirroring

##### Description:

Peer Glue 의 이미지를 Demote 합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |
| imageName      | path       | Image Name for Mirroring    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/resync/{pool}/{imageName}

#### PUT

##### Summary:

Resync Image Mirroring

##### Description:

Glue 의 이미지를 Resync 합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |
| imageName      | path       | Image Name for Mirroring    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/resync/peer/{pool}/{imageName}

#### PUT

##### Summary:

Resync Peer Image Mirroring

##### Description:

Peer Glue 의 이미지를 Resync 합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
| -------------- | ---------- | --------------------------- | -------- | ------ |
| mirrorPool     | path       | Pool Name for Mirroring     | Yes      | string |
| imageName      | path       | Image Name for Mirroring    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/{pool}/{imagename}

#### GET

##### Summary:

Show Infomation of Mirrored Snapshot

##### Description:

미러링중인 이미지의 정보를 보여줍니다.

##### Parameters

| Name      | Located in | Description | Required | Schema |
| --------- | ---------- | ----------- | -------- | ------ |
| imageName | path       | imageName   | Yes      | string |
| mirrorPool| path       | mirrorPool  | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [Message](#Message)                                       |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/{pool}/{imagename}

#### DELETE

##### Summary:

Delete Mirrored Snapshot Schedule

##### Description:

이미지의 미러링 스케줄링을 비활성화 합니다.

##### Parameters

| Name      | Located in | Description | Required | Schema |
| --------- | ---------- | ----------- | -------- | ------ |
| imageName | path       | imageName   | Yes      | string |
| mirrorPool| path       | mirrorPool  | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [Message](#Message)                                       |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/{pool}/{imagename}/{hostname}/{vmname}

#### POST

##### Summary:

Setup Image Mirroring Schedule

##### Description:

Glue의 이미지에 미러링 스케줄을 설정합니다.

##### Parameters

| Name      | Located in | Description                | Required | Schema |
| --------- | ---------- | -------------------------- | -------- | ------ |
| mirrorPool| path       | Pool Name of Mirroring     | Yes      | string |
| imageName | path       | Image Name for Mirroring   | Yes      | string |
| hostName  | path       | Host Name for Mirroring VM | Yes      | string |
| vmName    | path       | VM Name for Mirroring Image| Yes      | string |
| volType   | formData   | Volume Type                | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [Message](#Message)                                       |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/snapshot/{pool}/{imagename}

#### POST

##### Summary:

Take Image Mirroring Snapshot or Setup Image Mirroring Snapshot Schedule

##### Description:

Glue의 이미지에 미러링 스냅샷을 생성하거나 스케줄을 설정합니다.

##### Parameters

| Name      | Located in | Description                              | Required | Schema |
| --------- | ---------- | ---------------------------------------- | -------- | ------ |
| mirrorPool| path       | Pool Name of Mirroring                   | Yes      | string |
| vmName    | path       | VM Name for Mirroring Image              | Yes      | string |
| hostName  | formData   | Host Name for Mirroring VM               | No       | string |
| imageName | formData   | Image Name for Mirroring Image (Schedule)| No       | string |
| imageList | formData   | Image List for Mirroring (Manual)        | No       | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [Message](#Message)                                       |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/nfs

#### GET

##### Summary:

Show List of Glue NFS Cluster

##### Description:

Glue NFS Cluster의 리스트를 보여줍니다..

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [NfsClusterLs](#NfsClusterLs)                             |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/nfs/{cluster_id}

#### GET

##### Summary:

Show Info of Glue NFS Cluster

##### Description:

Glue NFS Cluster의 상세정보를 보여줍니다..

##### Parameters

| Name       | Located in | Description            | Required | Schema |
| ---------- | ---------- | ---------------------- | -------- | ------ |
| cluster_id | path       | NFS Cluster Identifier | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [NfsClusterInfo](#NfsClusterInfo)                         |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### DELETE

##### Summary:

Delete of Glue NFS Cluster

##### Description:

Glue NFS Cluster를 삭제합니다..

##### Parameters

| Name       | Located in | Description            | Required | Schema |
| ---------- | ---------- | ---------------------- | -------- | ------ |
| cluster_id | path       | NFS Cluster Identifier | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [NfsClusterDelete](#NfsClusterDelete)                     |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/nfs/{cluster_id}/{port}

#### POST

##### Summary:

Create of Glue NFS Cluster

##### Description:

Glue NFS Cluster를 생성합니다..

##### Parameters

| Name       | Located in | Description            | Required | Schema |
| ---------- | ---------- | ---------------------- | -------- | ------ |
| cluster_id | path       | NFS Cluster Identifier | Yes      | string |
| port       | path       | NFS Cluster Port       | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [NfsClusterCreate](#NfsClusterCreate)                     |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/nfs/export/{cluster_id}

#### GET

##### Summary:

Show Detailed of Glue NFS Export

##### Description:

Glue NFS Export의 상세정보를 보여줍니다..

##### Parameters

| Name       | Located in | Description            | Required | Schema |
| ---------- | ---------- | ---------------------- | -------- | ------ |
| cluster_id | path       | NFS Cluster Identifier | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [NfsExportDetailed](#NfsExportDetailed)                   |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### PUT

##### Summary:

Update of Glue NFS Export

##### Description:

Glue NFS Export를 수정합니다..

##### Parameters

| Name       | Located in | Description            | Required | Schema |
| ---------- | ---------- | ---------------------- | -------- | ------ |
| cluster_id | path       | NFS Cluster Identifier | Yes      | string |
| json_file  | body       | NFS Cluster JSON File  | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [NfsExportUpdate](#NfsExportUpdate)                       |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### POST

##### Summary:

Create of Glue NFS Export

##### Description:

Glue NFS Export를 생성합니다..

##### Parameters

| Name       | Located in | Description            | Required | Schema |
| ---------- | ---------- | ---------------------- | -------- | ------ |
| cluster_id | path       | NFS Cluster Identifier | Yes      | string |
| json_file  | body       | NFS Cluster JSON File  | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [NfsExportCreate](#NfsExportCreate)                       |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/nfs/export/{cluster_id}/{export_id}

#### DELETE

##### Summary:

Delete of Glue NFS Export

##### Description:

Glue NFS Export를 삭제합니다..

##### Parameters

| Name       | Located in | Description            | Required | Schema |
| ---------- | ---------- | ---------------------- | -------- | ------ |
| cluster_id | path       | NFS Cluster Identifier | Yes      | string |
| export_id  | path       | NFS Export ID          | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [NfsExportDetailed](#NfsExportDetailed)                   |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/iscsi

#### POST

##### Summary:

Create of Iscsi Servcie Daemon

##### Description:

Iscsi 서비스 데몬을 생성합니다.

##### Parameters

| Name         | Located in | Description                | Required | Schema   |
| ------------ | ---------- | -------------------------- | -------- | -------- |
| hosts        | formData   | Host Name                  | Yes      | []string |
| service_id   | formData   | ISCSI Service Name         | Yes      | string   |
| service_id   | formData   | ISCSI Service Name         | Yes      | string   |
| pool         | formData   | Pool Name                  | Yes      | string   |
| api_port     | formData   | ISCSI API Port             | Yes      | int      |
| api_user     | formData   | ISCSI API User             | Yes      | string   |
| api_password | formData   | ISCSI API Password         | Yes      | string   |
| count        | formData   | Iscsi Service Daemon Count | Yes      | int      |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [IscsiServiceCreate](#IscsiServiceCreate)                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/iscsi/discovery

#### GET

##### Summary:

Show of Iscsi Discovery Auth Details

##### Description:

Iscsi 계정 정보를 가져옵니다.

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [IscsiGetDiscoveryAuth](#IscsiGetDiscoveryAuth)           |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### PUT

##### Summary:

Update of Iscsi Discovery Auth Details

##### Description:

Iscsi 계정 정보를 수정합니다.

##### Parameters

| Name            | Located in | Description                                   | Required | Schema |
| --------------- | ---------- | --------------------------------------------- | -------- | ------ |
| user            | formData   | Iscsi Discovery Authorization Username        | No       | string |
| password        | formData   | Iscsi Discovery Authorization Password        | No       | string |
| mutual_user     | formData   | Iscsi Discovery Authorization Mutual Username | No       | string |
| mutual_password | formData   | Iscsi Discovery Authorization Mutual Password | No       | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [IscsiUpdateDiscoveryAuth](#IscsiUpdateDiscoveryAuth)     |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/iscsi/target

#### GET

##### Summary:

Show List of Iscsi Target

##### Description:

Iscsi 타겟 리스트를 가져옵니다.

##### Parameters

| Name   | Located in | Description           | Required | Schema |
| ------ | ---------- | --------------------- | -------- | ------ |
| iqn_id | query      | Iscsi Target IQN Name | No       | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [IscsiCommon](#IscsiCommon)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### DELETE

##### Summary:

Delete of Iscsi Target

##### Description:

Iscsi 타겟을 삭제합니다.

##### Parameters

| Name   | Located in | Description           | Required | Schema |
| ------ | ---------- | --------------------- | -------- | ------ |
| iqn_id | query      | Iscsi Target IQN Name | No       | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [IscsiCommon](#IscsiCommon)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### POST

##### Summary:

Create of Iscsi Target

##### Description:

Iscsi 타겟을 생성합니다.

##### Parameters

| Name            | Located in | Description                | Required | Schema   |
| --------------- | ---------- | -------------------------- | -------- | -------- |
| iqn_id          | formData   | Iscsi Target IQN Name      | Yes      | string   |
| hostname        | formData   | Gateway Host Name          | Yes      | []string |
| ip_address      | formData   | Gateway Host IP Address    | Yes      | []string |
| pool_name       | formData   | Glue Pool Name             | No       | []string |
| image_name      | formData   | Glue Image Name            | No       | []string |
| acl_enabled     | formData   | scsi Authentication        | Yes      | boolean  |
| username        | formData   | Iscsi Auth User            | No       | string   |
| password        | formData   | Iscsi Auth Password        | No       | string   |
| mutual_username | formData   | Iscsi Auth Mutual User     | No       | string   |
| mutual_password | formData   | Iscsi Auth Mutaul Password | No       | string   |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [IscsiCommon](#IscsiCommon)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/smb

#### GET

##### Summary:

Show Status of Smb Servcie Daemon

##### Description:

SMB 서비스 데몬 상태를 조회합니다.

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [SmbStatus](#SmbStatus)                                   |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### DELETE

##### Summary:

Delete of Smb Service

##### Description:

SMB 서비스 전체를 삭제합니다.

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | ["Success"]                                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### POST

##### Summary:

Create of Smb Service

##### Description:

SMB 서비스 전체를 생성합니다.

##### Parameters

| Name        | Located in | Description                   | Required | Schema |
| ----------- | ---------- | ----------------------------- | -------- | ------ |
| username    | formData   | SMB Username                  | Yes      | string |
| password    | formData   | SMB Password                  | Yes      | string |
| folder_name | formData   | SMB Share Folder Name         | Yes      | string |
| path        | formData   | SMB Server Actual Shared Path | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | ["Success"]                                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/smb/user

#### PUT

##### Summary:

Update User of Smb Service

##### Description:

SMB 서비스 사용자의 패스워드를 변경합니다.

##### Parameters

| Name     | Located in | Description  | Required | Schema |
| -------- | ---------- | ------------ | -------- | ------ |
| username | formData   | SMB Username | Yes      | string |
| password | formData   | SMB Password | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | ["Success"]                                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### DELETE

##### Summary:

Delete User of Smb Service

##### Description:

SMB 서비스 사용자를 삭제합니다.

##### Parameters

| Name     | Located in | Description  | Required | Schema |
| -------- | ---------- | ------------ | -------- | ------ |
| username | formData   | SMB Username | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | ["Success"]                                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### POST

##### Summary:

Create User of Smb Service

##### Description:

SMB 서비스 사용자를 생성합니다.

##### Parameters

| Name     | Located in | Description  | Required | Schema |
| -------- | ---------- | ------------ | -------- | ------ |
| username | formData   | SMB Username | Yes      | string |
| password | formData   | SMB Password | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | ["Success"]                                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /version

#### GET

##### Summary:

Show Versions of API

##### Description:

API 의 버전을 보여줍니다.

##### Responses

| Code | Description           | Schema                                                    |
| ---- | --------------------- | --------------------------------------------------------- |
| 200  | OK                    | [Version](#Version)                                       |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### Models

#### GlueStatus

| Name            | Type             | Description                | Required |
| --------------- | ---------------- | -------------------------- | -------- |
| election_epoch  | integer (uint32) |                            | No       |
| fsid            | string (uuid)    | Glue클러스터를 구분하는 ID | No       |
| fsmap           | object           |                            | No       |
| health          | object           |                            | No       |
| mgrmap          | object           |                            | No       |
| monmap          | object           |                            | No       |
| osdmap          | object           |                            | No       |
| pgmap           | object           |                            | No       |
| progress_events | object           |                            | No       |
| quorum          | [ integer ]      |                            | No       |
| quorum_age      | integer          |                            | No       |
| quorum_names    | [ string ]       |                            | No       |
| servicemap      | object           |                            | No       |

#### GlueVersion

| Name       | Type   | Description | Required |
| ---------- | ------ | ----------- | -------- |
| mgr        | object |             | No       |
| mon        | object |             | No       |
| osd        | object |             | No       |
| overall    | object |             | No       |
| rbd-mirror | object |             | No       |
| rgw        | object |             | No       |

#### HTTP400BadRequest

| Name    | Type    | Description | Required |
| ------- | ------- | ----------- | -------- |
| code    | integer |             | No       |
| message | string  |             | No       |

#### HTTP404NotFound

| Name    | Type    | Description | Required |
| ------- | ------- | ----------- | -------- |
| code    | integer |             | No       |
| message | string  |             | No       |

#### HTTP500InternalServerError

| Name    | Type    | Description | Required |
| ------- | ------- | ----------- | -------- |
| code    | integer |             | No       |
| message | string  |             | No       |

#### Message

| Name    | Type   | Description | Required |
| ------- | ------ | ----------- | -------- |
| message | string |             | No       |

#### MirrorImage

| Name      | Type                                    | Description | Required |
| --------- | --------------------------------------- | ----------- | -------- |
| image     | string                                  |             | No       |
| items     | [ [MirrorImageItem](#MirrorImageItem) ] |             | No       |
| namespace | string                                  |             | No       |
| pool      | string                                  |             | No       |

#### MirrorImageItem

| Name       | Type   | Description | Required |
| ---------- | ------ | ----------- | -------- |
| interval   | string |             | No       |
| start_time | string |             | No       |

#### MirrorList

| Name   | Type                            | Description | Required |
| ------ | ------------------------------- | ----------- | -------- |
| Local  | [ [MirrorImage](#MirrorImage) ] |             | No       |
| Remote | [ [MirrorImage](#MirrorImage) ] |             | No       |
| debug  | boolean (bool)                  | Debug info  | No       |

#### MirrorSetup

| Name              | Type   | Description | Required |
| ----------------- | ------ | ----------- | -------- |
| host              | string |             | No       |
| localClusterName  | string | 미러링 상태 | No       |
| localToken        | string |             | No       |
| mirrorPool        | string |             | No       |
| privateKeyFile    | object |             | No       |
| remoteClusterName | string | 미러링 상태 | No       |
| remoteToken       | string |             | No       |

#### MirrorStatus

| Name          | Type   | Description      | Required |
| ------------- | ------ | ---------------- | -------- |
| daemon_health | string | 미러링 데몬 상태 | No       |
| health        | string | 미러링 상태      | No       |
| image_health  | string | 이미지 상태      | No       |
| states        | object | 이미지 상세      | No       |

#### FsStatus

| Name       | Type   | Description | Required |
| ---------- | ------ | ----------- | -------- |
| clients    | object |             | No       |
| mdsversion | object |             | No       |
| mdsmap     | object |             | No       |
| pools      | object |             | No       |

#### FsStatus

| Name   | Type    | Description | Required |
| ------ | ------- | ----------- | -------- |
| mdsmap | object  |             | No       |
| id     | integer |             | No       |

#### FsList

| Name           | Type      | Description | Required |
| -------------- | --------- | ----------- | -------- |
| name           | string    |             | No       |
| metadatapool   | string    |             | No       |
| metadatapoolid | string    |             | No       |
| datapoolids    | []integer |             | No       |
| datapools      | []string  |             | No       |

#### NfsClusterLs

| Name | Type   | Description | Required |
| ---- | ------ | ----------- | -------- |
| name | string |             | No       |

#### NfsClusterInfo

| Name | Type   | Description | Required |
| ---- | ------ | ----------- | -------- |
|      | object |             | No       |

#### NfsExportDetailed

| Name          | Type     | Description | Required |
| ------------- | -------- | ----------- | -------- |
| accesstype    | string   |             | No       |
| clients       | []string |             | No       |
| clusterid     | string   |             | No       |
| exportid      | integer  |             | No       |
| fsal          | object   |             | No       |
| path          | string   |             | No       |
| protocols     | []string |             | No       |
| pseudo        | string   |             | No       |
| securitylabel | boolean  |             | No       |
| squash        | string   |             | No       |
| transports    | []string |             | No       |

#### IscsiServiceCreate

| Name         | Type   | Description | Required |
| ------------ | ------ | ----------- | -------- |
| service_type | string |             | No       |
| service_id   | string |             | No       |
| placement    | object |             | No       |
| spec         | object |             | No       |

#### IscsiTargetList

| Name    | Type   | Description | Required |
| ------- | ------ | ----------- | -------- |
| targets | object |             | No       |

#### IscsiDiscoveryAuth

| Name            | Type   | Description | Required |
| --------------- | ------ | ----------- | -------- |
| username        | string |             | No       |
| password        | string |             | No       |
| mutual_username | string |             | No       |
| mutual_password | string |             | No       |

#### SmbStatus

| Name        | Type   | Description | Required |
| ----------- | ------ | ----------- | -------- |
| names       | string |             | No       |
| description | string |             | No       |
| status      | string |             | No       |
| state       | string |             | No       |
| users       | object |             | No       |

#### IscsiCommon

| Name | Type   | Description | Required |
| ---- | ------ | ----------- | -------- |
|      | object |             | No       |

#### Version

| Name    | Type            | Description | Required |
| ------- | --------------- | ----------- | -------- |
| version | string (string) |             | No       |
