{{- $name := .Name}}
{{- $lowerCamelName := lowerCamel .Name}}
{{- $parentName := .Extra.ParentName}}
import { PageContainer } from '@ant-design/pro-components';
import React, { useRef, useReducer } from 'react';
import { useIntl } from '@umijs/max';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { ProTable{{if .Extra.IndexProComponentsImport}}, {{.Extra.IndexProComponentsImport}}{{end}} } from '@ant-design/pro-components';
import { Space, message{{if .Extra.IndexAntdImport}}, {{.Extra.IndexAntdImport}}{{end}} } from 'antd';
import { fetch{{$name}}, del{{$name}} } from '@/services/{{$parentName}}/{{$lowerCamelName}}';
import {{$name}}Modal from './components/SaveForm';
import { AddButton, EditIconButton, DelIconButton } from '@/components/Button';

enum ActionTypeEnum {
  ADD,
  EDIT,
  CANCEL,
}

interface Action {
  type: ActionTypeEnum;
  payload?: API.{{$name}};
}

interface State {
  visible: boolean;
  title: string;
  id?: string;
}

const {{$name}}: React.FC = () => {
  const intl = useIntl();
  const actionRef = useRef<ActionType>();
  const addTitle = intl.formatMessage({ id: 'pages.{{$parentName}}.{{$lowerCamelName}}.add', defaultMessage: 'Add {{$name}}' });
  const editTitle = intl.formatMessage({ id: 'pages.{{$parentName}}.{{$lowerCamelName}}.edit', defaultMessage: 'Edit {{$name}}' });
  const delTip = intl.formatMessage({ id: 'pages.{{$parentName}}.{{$lowerCamelName}}.delTip', defaultMessage: 'Are you sure you want to delete this record?' });

  const [state, dispatch] = useReducer(
    (pre: State, action: Action) => {
      switch (action.type) {
        case ActionTypeEnum.ADD:
          return {
            visible: true,
            title: addTitle,
          };
        case ActionTypeEnum.EDIT:
          return {
            visible: true,
            title: editTitle,
            id: action.payload?.id,
          };
        case ActionTypeEnum.CANCEL:
          return {
            visible: false,
            title: '',
            id: undefined,
          };
        default:
          return pre;
      }
    },
    { visible: false, title: '' },
  );

  const columns: ProColumns<API.{{$name}}>[] = [
    {{- range .Fields}}
    {{- if .Extra.ColumnComponent}}
    {{.Extra.ColumnComponent}},
    {{- end}}
    {{- end}}
    {{- if .Extra.IncludeCreatedAt}}
    {
      title: intl.formatMessage({ id: 'pages.table.column.created_at' }),
      dataIndex: 'created_at',
      valueType: 'dateTime',
      search: false,
      width: 160,
    },
    {{- end}}
    {{- if .Extra.IncludeUpdatedAt}}
    {
      title: intl.formatMessage({ id: 'pages.table.column.updated_at' }),
      dataIndex: 'updated_at',
      valueType: 'dateTime',
      search: false,
      width: 160,
    },
    {{- end}}
    {
      title: intl.formatMessage({ id: 'pages.table.column.operation' }),
      valueType: 'option',
      key: 'option',
      width: 130,
      render: (_, record) => (
        <Space size={2}>
          <EditIconButton
            key="edit"
            code="edit"
            onClick={() => {
              dispatch({ type: ActionTypeEnum.EDIT, payload: record });
            {{`}}`}}
          />
          <DelIconButton
            key="delete"
            code="delete"
            title={delTip}
            onConfirm={async () => {
              const res = await del{{$name}}(record.id!);
              if (res.success) {
                message.success(intl.formatMessage({ id: 'component.message.success.delete' }));
                actionRef.current?.reload();
              }
            {{`}}`}}
          />
        </Space>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.{{$name}}, API.PaginationParam>
        columns={columns}
        actionRef={actionRef}
        request={fetch{{$name}}}
        rowKey="id"
        cardBordered
        search={{`{{`}}
          labelWidth: 'auto',
        {{`}}`}}
        pagination={{`{{`}} defaultPageSize: 10, showSizeChanger: true {{`}}`}}
        options={{`{{`}}
          density: true,
          fullScreen: true,
          reload: true,
        {{`}}`}}
        dateFormatter="string"
        toolBarRender={() => [
          <AddButton
            key="add"
            code="add"
            onClick={() => {
              dispatch({ type: ActionTypeEnum.ADD });
            {{`}}`}}
          />,
        ]}
      />

      <{{$name}}Modal
        visible={state.visible}
        title={state.title}
        id={state.id}
        onCancel={() => {
          dispatch({ type: ActionTypeEnum.CANCEL });
        {{`}}`}}
        onSuccess={() => {
          dispatch({ type: ActionTypeEnum.CANCEL });
          actionRef.current?.reload();
        {{`}}`}}
      />
    </PageContainer>
  );
};

export default {{$name}};
