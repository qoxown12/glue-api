# Glue-API
Glue의 기능을 제어하기 위한 REST API 입니다.

본 문서는 swagger로 작성된 API 목록을 [swagger-markdown-ui](https://swagger-markdown-ui.netlify.app)를 사용해 작성된 README입니다.

## API 목록


| Method | API                                                                    |        진행도         | 비고                |
|--------|------------------------------------------------------------------------|:------------------:|-------------------|
| GET    | [api/v1/version](#version)                                             | :white_check_mark: | Version           |
| GET    | [api/v1/glue](#apiv1glue)                                              | :white_check_mark: | GlueStatus        |
| GET    | [api/v1/glue/version](#apiv1glueversion)                               | :white_check_mark: | GlueVersion       |
| GET    | [api/v1/glue/pool](#apiv1gluepool)                                     | :white_check_mark: | ListPools         |
| DELETE | [api/v1/glue/pool/:poolname](#apiv1gluepool)                           | :white_check_mark: | PoolDelete        |
| GET    | [api/v1/glue/rbd/:poolname](#apiv1gluerbdpool)                         | :white_check_mark: | ListImages        |
| GET    | [api/v1/gluefs](#apiv1gluefs)                                          | :white_check_mark: | FsStatus          |
| GET    | [api/v1/gluefs/info/:fs_name](#apiv1gluefsinfo)                        | :white_check_mark: | FsGetInfo         |
| GET    | [api/v1/gluefs/list](#apiv1gluefslist)                                 | :white_check_mark: | FsList            |
| POST   | [api/v1/gluefs/:fs_name](#apiv1gluefsname)                             | :white_check_mark: | FsCreate          |
| DELETE | [api/v1/gluefs/:fs_name](#apiv1gluefsname)                             | :white_check_mark: | FsDelete          |
| GET    | [api/v1/mirror](#apiv1mirror)                                          | :white_check_mark: | MirrorStatus      |
| POST   | [api/v1/mirror](#apiv1mirror)                                          | :white_check_mark: | MirrorSetup       |
| PATCH  | [api/v1/mirror]()                                                      |                    |                   |
| DELETE | [api/v1/mirror](#apiv1mirror)                                          |     :recycle:      | MirrorDelete      |
| GET    | [api/v1/mirror/image](#apiv1mirrorimage)                               | :white_check_mark: | MirrorImageList   |
| GET    | [api/v1/mirror/image/:pool/:imagename]()                               |                    |                   |
| POST   | [api/v1/mirror/image/:pool/:imagename]()                               |                    |                   |
| PATCH  | [api/v1/mirror/image/:pool/:imagename]()                               |                    |                   |
| DELETE | [api/v1/mirror/image/:pool/:imagename](#apiv1mirrorimagepoolimagename) | :white_check_mark: | MirrorImageDelete |
| POST   | [api/v1/mirror/image/prymary/:pool/:image]()                           |                    |                   |
| DELETE | [api/v1/mirror/image/prymary/:pool/:image]()                           |                    |                   |
| GET    | [api/v1/mirror/image/prymary/:pool/:image]()                           |                    |                   |
| GET    | [api/v1/nfs](#apiv1nfs)                                                | :white_check_mark: | NfsClusterLs      |
| GET    | [api/v1/nfs/export/:cluster_id](#apiv1nfsexportget)                    | :white_check_mark: | NfsExportDetailed |
| PUT    | [api/v1/nfs/export/:cluster_id](#apiv1nfsexportput)                    | :white_check_mark: | NfsExportUpdate   |
| POST   | [api/v1/nfs/export/:cluster_id](#apiv1nfsexportporst)                  | :white_check_mark: | NfsExportCreate   |
| DELETE | [api/v1/nfs/export/:cluster_id/:export_id](#apiv1nfsexportdel)         | :white_check_mark: | NfsExportDelete   |
| GET    | [api/v1/nfs/:cluster_id](#apiv1nfsclusterget)                          | :white_check_mark: | NfsClusterInfo    |
| DELETE | [api/v1/nfs/:cluster_id](#apiv1nfsclusterdel)                          | :white_check_mark: | NfsClusterDelete  |
| POST   | [api/v1/nfs/:cluster_id/:port]()                                       | :white_check_mark: | NfsClusterCreate  |
| ANY    | swagger/index.html                                                     | :white_check_mark: |                   |

### /api/v1/glue

#### GET
##### Summary:

Show Status of Glue

##### Description:

Glue 의 상태값을 보여줍니다.

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description     | Required | Schema |
|-----------|------------|-----------------|----------|--------|
| pool_name | path       | Pool Name       | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description     | Required | Schema |
|-----------|------------|-----------------|----------|--------|
| pool_name | path       | Pool Name       | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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
|------|-----------------------|-----------------------------------------------------------|
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

| Name    | Located in | Description     | Required | Schema |
|---------|------------|-----------------|----------|--------|
| fs_name | path       | Glue FS Name    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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
|------|-----------------------|-----------------------------------------------------------|
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

| Name    | Located in | Description     | Required | Schema |
|---------|------------|-----------------|----------|--------|
| fs_name | path       | Glue FS Name    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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

| Name    | Located in | Description     | Required | Schema |
|---------|------------|-----------------|----------|--------|
| fs_name | path       | Glue FS Name    | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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
|------|-----------------------|-----------------------------------------------------------|
| 200  | OK                    | [GlueVersion](#GlueVersion)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror

#### DELETE
##### Summary:

Delete Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터를 제거합니다.

##### Parameters

| Name           | Located in | Description                 | Required | Schema |
|----------------|------------|-----------------------------|----------|--------|
| host           | formData   | Remote Cluster Host Address | Yes      | string |
| privateKeyFile | formData   | Remote Cluster PrivateKey   | Yes      | file   |
| mirrorPool     | formData   | Pool Name for Mirroring     | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
| 200  | OK                    | [MirrorSetup](#MirrorSetup)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### GET
##### Summary:

Show Status of Mirror

##### Description:

Glue 의 미러링 상태를 보여줍니다.

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
| 200  | OK                    | [MirrorStatus](#MirrorStatus)                             |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

#### POST
##### Summary:

Setup Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터를 설정합니다..

##### Parameters

| Name              | Located in | Description                 | Required | Schema |
|-------------------|------------|-----------------------------|----------|--------|
| localClusterName  | formData   | Local Cluster Name          | Yes      | string |
| remoteClusterName | formData   | Remote Cluster Name         | Yes      | string |
| host              | formData   | Remote Cluster Host Address | Yes      | string |
| privateKeyFile    | formData   | Remote Cluster PrivateKey   | Yes      | file   |
| mirrorPool        | formData   | Pool Name for Mirroring     | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
| 200  | OK                    | [MirrorSetup](#MirrorSetup)                               |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image

#### GET
##### Summary:

Show List of Mirrored Image

##### Description:

미러링중인 이미지의 목록과 상태를 보여줍니다.

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
| 200  | OK                    | [MirrorList](#MirrorList)                                 |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### /api/v1/mirror/image/{pool}/{imagename}

#### DELETE
##### Summary:

Delete Mirrored Image

##### Description:

이미지의 미러링을 비활성화 합니다.

##### Parameters

| Name      | Located in | Description | Required | Schema |
|-----------|------------|-------------|----------|--------|
| imageName | path       | imageName   | Yes      | string |
| pool      | path       | pool        | Yes      | string |

##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description             | Required | Schema |
|-----------|------------|-------------------------|----------|--------|
| cluster_id| path       | NFS Cluster Identifier  | Yes      | string |


##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description             | Required | Schema |
|-----------|------------|-------------------------|----------|--------|
| cluster_id| path       | NFS Cluster Identifier  | Yes      | string |


##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description             | Required | Schema |
|-----------|------------|-------------------------|----------|--------|
| cluster_id| path       | NFS Cluster Identifier  | Yes      | string |
| port      | path       | NFS Cluster Port        | Yes      | string |


##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description             | Required | Schema |
|-----------|------------|-------------------------|----------|--------|
| cluster_id| path       | NFS Cluster Identifier  | Yes      | string |


##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description             | Required | Schema |
|-----------|------------|-------------------------|----------|--------|
| cluster_id| path       | NFS Cluster Identifier  | Yes      | string |
| json_file | body       | NFS Cluster JSON File   | Yes      | string |


##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description             | Required | Schema |
|-----------|------------|-------------------------|----------|--------|
| cluster_id| path       | NFS Cluster Identifier  | Yes      | string |
| json_file | body       | NFS Cluster JSON File   | Yes      | string |


##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
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

| Name      | Located in | Description             | Required | Schema |
|-----------|------------|-------------------------|----------|--------|
| cluster_id| path       | NFS Cluster Identifier  | Yes      | string |
| export_id | path       | NFS Export ID           | Yes      | string |


##### Responses

| Code | Description           | Schema                                                    |
|------|-----------------------|-----------------------------------------------------------|
| 200  | OK                    | [NfsExportDetailed](#NfsExportDetailed)                   |
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
|------|-----------------------|-----------------------------------------------------------|
| 200  | OK                    | [Version](#Version)                                       |
| 400  | Bad Request           | [HTTP400BadRequest](#HTTP400BadRequest)                   |
| 404  | Not Found             | [HTTP404NotFound](#HTTP404NotFound)                       |
| 500  | Internal Server Error | [HTTP500InternalServerError](#HTTP500InternalServerError) |

### Models


#### GlueStatus

| Name            | Type             | Description       | Required |
|-----------------|------------------|-------------------|----------|
| election_epoch  | integer (uint32) |                   | No       |
| fsid            | string (uuid)    | Glue클러스터를 구분하는 ID | No       |
| fsmap           | object           |                   | No       |
| health          | object           |                   | No       |
| mgrmap          | object           |                   | No       |
| monmap          | object           |                   | No       |
| osdmap          | object           |                   | No       |
| pgmap           | object           |                   | No       |
| progress_events | object           |                   | No       |
| quorum          | [ integer ]      |                   | No       |
| quorum_age      | integer          |                   | No       |
| quorum_names    | [ string ]       |                   | No       |
| servicemap      | object           |                   | No       |

#### GlueVersion

| Name       | Type           | Description | Required |
|------------|----------------|-------------|----------|
| mgr        | object         |             | No       |
| mon        | object         |             | No       |
| osd        | object         |             | No       |
| overall    | object         |             | No       |
| rbd-mirror | object         |             | No       |
| rgw        | object         |             | No       |

#### HTTP400BadRequest

| Name    | Type           | Description | Required |
|---------|----------------|-------------|----------|
| code    | integer        |             | No       |
| message | string         |             | No       |

#### HTTP404NotFound

| Name    | Type           | Description | Required |
|---------|----------------|-------------|----------|
| code    | integer        |             | No       |
| message | string         |             | No       |

#### HTTP500InternalServerError

| Name    | Type           | Description | Required |
|---------|----------------|-------------|----------|
| code    | integer        |             | No       |
| message | string         |             | No       |

#### Message

| Name    | Type   | Description | Required |
|---------|--------|-------------|----------|
| message | string |             | No       |

#### MirrorImage

| Name      | Type                                    | Description | Required |
|-----------|-----------------------------------------|-------------|----------|
| image     | string                                  |             | No       |
| items     | [ [MirrorImageItem](#MirrorImageItem) ] |             | No       |
| namespace | string                                  |             | No       |
| pool      | string                                  |             | No       |

#### MirrorImageItem

| Name       | Type   | Description | Required |
|------------|--------|-------------|----------|
| interval   | string |             | No       |
| start_time | string |             | No       |

#### MirrorList

| Name   | Type                            | Description | Required |
|--------|---------------------------------|-------------|----------|
| Local  | [ [MirrorImage](#MirrorImage) ] |             | No       |
| Remote | [ [MirrorImage](#MirrorImage) ] |             | No       |
| debug  | boolean (bool)                  | Debug info  | No       |

#### MirrorSetup

| Name              | Type           | Description | Required |
|-------------------|----------------|-------------|----------|
| host              | string         |             | No       |
| localClusterName  | string         | 미러링 상태      | No       |
| localToken        | string         |             | No       |
| mirrorPool        | string         |             | No       |
| privateKeyFile    | object         |             | No       |
| remoteClusterName | string         | 미러링 상태      | No       |
| remoteToken       | string         |             | No       |

#### MirrorStatus

| Name          | Type           | Description | Required |
|---------------|----------------|-------------|----------|
| daemon_health | string         | 미러링 데몬 상태   | No       |
| health        | string         | 미러링 상태      | No       |
| image_health  | string         | 이미지 상태      | No       |
| states        | object         | 이미지 상세      | No       |

#### FsStatus

| Name       | Type           | Description | Required |
|------------|----------------|-------------|----------|
| clients    | object         |             | No       |
| mdsversion | object         |             | No       |
| mdsmap     | object         |             | No       |
| pools      | object         |             | No       |

#### FsStatus

| Name       | Type           | Description | Required |
|------------|----------------|-------------|----------|
| mdsmap     | object         |             | No       |
| id         | integer        |             | No       |

#### FsList

| Name          | Type           | Description | Required |
|---------------|----------------|-------------|----------|
| name          | string         |             | No       |
| metadatapool  | string         |             | No       |
| metadatapoolid| string         |             | No       |
| datapoolids   | []integer      |             | No       |
| datapools     | []string       |             | No       |

#### NfsClusterLs

| Name          | Type           | Description | Required |
|---------------|----------------|-------------|----------|
| name          | string         |             | No       |

#### NfsClusterInfo

| Name          | Type           | Description | Required |
|---------------|----------------|-------------|----------|
|               | object         |             | No       |

#### NfsExportDetailed

| Name          | Type           | Description | Required |
|---------------|----------------|-------------|----------|
| accesstype    | string         |             | No       |
| clients       | []string       |             | No       |
| clusterid     | string         |             | No       |
| exportid      | integer        |             | No       |
| fsal          | object         |             | No       |
| path          | string         |             | No       |
| protocols     | []string       |             | No       |
| pseudo        | string         |             | No       |
| securitylabel | boolean        |             | No       |
| squash        | string         |             | No       |
| transports    | []string       |             | No       |

#### Version

| Name    | Type            | Description | Required |
|---------|-----------------|-------------|----------|
| version | string (string) |             | No       |
