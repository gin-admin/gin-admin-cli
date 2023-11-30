{{- $name := .Name}}
{{- $lowerPluralName := lowerHyphensPlural .Name}}
// @ts-ignore
/* eslint-disable */
import { request } from 'umi';

/** 查询列表 GET /api/v1/{{$lowerPluralName}} */
export async function fetch{{$name}}(params: API.PaginationParam, options?: { [key: string]: any }) {
  return request<API.ResponseResult<API.{{$name}}[]>>('/api/v1/{{$lowerPluralName}}', {
    method: 'GET',
    params: {
      current: '1',
      pageSize: '10',
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建数据 POST /api/v1/{{$lowerPluralName}} */
export async function add{{$name}}(body: API.{{$name}}, options?: { [key: string]: any }) {
  return request<API.ResponseResult<API.{{$name}}>>('/api/v1/{{$lowerPluralName}}', {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/** 获取单条记录 GET /api/v1/{{$lowerPluralName}}/${id} */
export async function get{{$name}}(id: string, options?: { [key: string]: any }) {
  return request<API.ResponseResult<API.{{$name}}>>(`/api/v1/{{$lowerPluralName}}/${id}`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 更新记录 PUT /api/v1/{{$lowerPluralName}}/${id} */
export async function update{{$name}}(id: string, body: API.{{$name}}, options?: { [key: string]: any }) {
  return request<API.ResponseResult<any>>(`/api/v1/{{$lowerPluralName}}/${id}`, {
    method: 'PUT',
    data: body,
    ...(options || {}),
  });
}

/** 删除记录 DELETE /api/v1/{{$lowerPluralName}}/${id} */
export async function del{{$name}}(id: string, options?: { [key: string]: any }) {
  return request<API.ResponseResult<any>>(`/api/v1/{{$lowerPluralName}}/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}
