{{- $name := .Name}}
{{- $lowerPluralName := lowerHyphensPlural .Name}}
// @ts-ignore
/* eslint-disable */
import { request } from 'umi';

/** Query list GET /api/v1/{{$lowerPluralName}} */
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

/** Create record POST /api/v1/{{$lowerPluralName}} */
export async function add{{$name}}(body: API.{{$name}}, options?: { [key: string]: any }) {
  return request<API.ResponseResult<API.{{$name}}>>('/api/v1/{{$lowerPluralName}}', {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/** Get record by ID GET /api/v1/{{$lowerPluralName}}/${id} */
export async function get{{$name}}(id: string, options?: { [key: string]: any }) {
  return request<API.ResponseResult<API.{{$name}}>>(`/api/v1/{{$lowerPluralName}}/${id}`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** Update record by ID PUT /api/v1/{{$lowerPluralName}}/${id} */
export async function update{{$name}}(id: string, body: API.{{$name}}, options?: { [key: string]: any }) {
  return request<API.ResponseResult<any>>(`/api/v1/{{$lowerPluralName}}/${id}`, {
    method: 'PUT',
    data: body,
    ...(options || {}),
  });
}

/** Delete record by ID DELETE /api/v1/{{$lowerPluralName}}/${id} */
export async function del{{$name}}(id: string, options?: { [key: string]: any }) {
  return request<API.ResponseResult<any>>(`/api/v1/{{$lowerPluralName}}/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}
