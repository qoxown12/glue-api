# glue-api



### /api/v1/glue

#### GET
##### Summary:

Show Status of Glue

##### Description:

Glue 의 상태값을 보여줍니다.

##### Responses

| Code | Description           | Schema                                                                      |
|------|-----------------------|-----------------------------------------------------------------------------|
| 200  | OK                    | [GlueStatus](#GlueStatus)                                                   |
| 400  | Bad Request           | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest)                   |
| 404  | Not Found             | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound)                       |
| 500  | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

### /api/v1/glue/pool

#### GET
##### Summary:

List Pools of Glue

##### Description:

Glue 의 스토리지 풀 목록을 보여줍니다.

##### Responses

| Code | Description           | Schema                                                                      |
|------|-----------------------|-----------------------------------------------------------------------------|
| 200  | OK                    | [GlueVersion](#GlueVersion)                                                 |
| 400  | Bad Request           | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest)                   |
| 404  | Not Found             | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound)                       |
| 500  | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

### /api/v1/glue/pool/{pool}

#### GET
##### Summary:

List Images of Pool Glue

##### Description:

Glue 스토리지 풀의 이미지 목록을 보여줍니다.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| pool | path | pool | Yes | string |

##### Responses

| Code | Description           | Schema                                                                      |
|------|-----------------------|-----------------------------------------------------------------------------|
| 200  | OK                    | [GlueVersion](#GlueVersion)                                                 |
| 400  | Bad Request           | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest)                   |
| 404  | Not Found             | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound)                       |
| 500  | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

### /api/v1/glue/version

#### GET
##### Summary:

Show Versions of Glue

##### Description:

Glue 의 버전을 보여줍니다.

##### Responses

| Code | Description           | Schema                                                                      |
|------|-----------------------|-----------------------------------------------------------------------------|
| 200  | OK                    | [GlueVersion](#GlueVersion)                                                 |
| 400  | Bad Request           | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest)                   |
| 404  | Not Found             | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound)                       |
| 500  | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

### /api/v1/mirror

#### DELETE
##### Summary:

Delete Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터를 제거합니다.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| host | formData | Remote Cluster Host Address | Yes | string |
| privateKeyFile | formData | Remote Cluster PrivateKey | Yes | file |
| mirrorPool | formData | Pool Name for Mirroring | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [MirrorSetup](#MirrorSetup) |
| 400 | Bad Request | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest) |
| 404 | Not Found | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound) |
| 500 | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

#### GET
##### Summary:

Show Status of Mirror

##### Description:

Glue 의 미러링 상태를 보여줍니다.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [MirrorStatus](#MirrorStatus) |
| 400 | Bad Request | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest) |
| 404 | Not Found | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound) |
| 500 | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

#### POST
##### Summary:

Setup Mirroring Cluster

##### Description:

Glue 의 미러링 클러스터를 설정합니다.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| localClusterName | formData | Local Cluster Name | Yes | string |
| remoteClusterName | formData | Remote Cluster Name | Yes | string |
| host | formData | Remote Cluster Host Address | Yes | string |
| privateKeyFile | formData | Remote Cluster PrivateKey | Yes | file |
| mirrorPool | formData | Pool Name for Mirroring | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [MirrorSetup](#MirrorSetup) |
| 400 | Bad Request | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest) |
| 404 | Not Found | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound) |
| 500 | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

### /api/v1/mirror/image

#### GET
##### Summary:

Show List of Mirrored Image

##### Description:

미러링중인 이미지의 목록과 상태를 보여줍니다.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [MirrorList](#MirrorList) |
| 400 | Bad Request | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest) |
| 404 | Not Found | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound) |
| 500 | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

### /api/v1/mirror/image/{pool}/{imagename}

#### DELETE
##### Summary:

Delete Mirrored Image

##### Description:

이미지의 미러링을 비활성화 합니다.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| imageName | path | imageName | Yes | string |
| pool | path | pool | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controller.Message](#controller.Message) |
| 400 | Bad Request | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest) |
| 404 | Not Found | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound) |
| 500 | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

### /version

#### GET
##### Summary:

Show Versions of API

##### Description:

API 의 버전을 보여줍니다.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [Version](#Version) |
| 400 | Bad Request | [httputil.HTTP400BadRequest](#httputil.HTTP400BadRequest) |
| 404 | Not Found | [httputil.HTTP404NotFound](#httputil.HTTP404NotFound) |
| 500 | Internal Server Error | [httputil.HTTP500InternalServerError](#httputil.HTTP500InternalServerError) |

### Models


#### GlueStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| debug | boolean (bool) | Debug info | No |
| election_epoch | integer (uint32) |  | No |
| fsid | string (uuid) | Glue클러스터를 구분하는 ID | No |
| fsmap | object |  | No |
| health | object |  | No |
| mgrmap | object |  | No |
| monmap | object |  | No |
| osdmap | object |  | No |
| pgmap | object |  | No |
| progress_events | object |  | No |
| quorum | [ integer ] |  | No |
| quorum_age | integer |  | No |
| quorum_names | [ string ] |  | No |
| servicemap | object |  | No |

#### GlueVersion

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| debug | boolean (bool) | Debug info | No |
| mgr | object |  | No |
| mon | object |  | No |
| osd | object |  | No |
| overall | object |  | No |
| rbd-mirror | object |  | No |
| rgw | object |  | No |

#### MirrorImage

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| image | string |  | No |
| items | [ [MirrorImageItem](#MirrorImageItem) ] |  | No |
| namespace | string |  | No |
| pool | string |  | No |

#### MirrorImageItem

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| interval | string |  | No |
| start_time | string |  | No |

#### MirrorList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Local | [ [MirrorImage](#MirrorImage) ] |  | No |
| Remote | [ [MirrorImage](#MirrorImage) ] |  | No |
| debug | boolean (bool) | Debug info | No |

#### MirrorSetup

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| debug | boolean (bool) | Debug info | No |
| host | string |  | No |
| localClusterName | string | 미러링 상태 | No |
| localToken | string |  | No |
| mirrorPool | string |  | No |
| privateKeyFile | object |  | No |
| remoteClusterName | string | 미러링 상태 | No |
| remoteToken | string |  | No |

#### MirrorStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| daemon_health | string | 미러링 데몬 상태 | No |
| debug | boolean (bool) | Debug info | No |
| health | string | 미러링 상태 | No |
| image_health | string | 이미지 상태 | No |
| states | object | 이미지 상세 | No |

#### Version

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| debug | boolean (bool) | Debug info | No |
| version | string (string) |  | No |

#### controller.Message

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### httputil.HTTP400BadRequest

| Name    | Type           | Description | Required |
|---------|----------------|-------------|----------|
| code    | integer        |             | No       |
| debug   | boolean (bool) | Debug info  | No       |
| message | string         |             | No       |

#### httputil.HTTP404NotFound

| Name    | Type           | Description | Required |
|---------|----------------|-------------|----------|
| code    | integer        |             | No       |
| debug   | boolean (bool) | Debug info  | No       |
| message | string         |             | No       |

#### httputil.HTTP500InternalServerError

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code  | integer |  | No |
| debug | boolean (bool) | Debug info | No |
| message | string |  | No |