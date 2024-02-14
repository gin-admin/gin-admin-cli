{{- $name := .Name}}
import { PageContainer } from '@ant-design/pro-components';
import React, { useRef, useReducer } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { ProTable{{if .Extra.IndexProComponentsImport}}, {{.Extra.IndexProComponentsImport}}{{end}} } from '@ant-design/pro-components';
import { Space, message{{if .Extra.IndexAntdImport}}, {{.Extra.IndexAntdImport}}{{end}}} from 'antd';
import { fetch{{$name}}, del{{$name}} } from '@/services/{{.Extra.ImportService}}';
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
  const actionRef = useRef<ActionType>();
  const addTitle = "{{.Extra.AddTitle}}";
  const editTitle = "{{.Extra.EditTitle}}";
  const delTip = "{{.Extra.DelTip}}";

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
    {{- $fieldName := .Name}}
    {{- $fieldLabel :=.Extra.Label}}
    {{- if eq .Extra.InColumn "true"}}
    {{- if .Extra.ColumnComponent}}
    {{.Extra.ColumnComponent}},
    {{- else}}
    {
      title: "{{$fieldLabel}}",
      dataIndex: '{{lowerUnderline $fieldName}}',
      {{- if .Extra.Ellipsis}}
      ellipsis: true,
      {{- end}}
      {{- if .Extra.Width}}
      width: {{.Extra.Width}},
      {{- end}}
      {{- if .Extra.SearchKey}}
      key: '{{.Extra.SearchKey}}',
      {{- else}}
      search: false,
      {{- end}}
      {{- if .Extra.ValueType}}
      valueType: '{{.Extra.ValueType}}',
      {{- end}}
    },
    {{- end}}
    {{- end}}
    {{- end}}
    {
      title: "{{.Extra.ActionText}}",
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
                message.success("{{.Extra.DeleteSuccessMessage}}");
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
        pagination={{`{{`}} pageSize: 10, showSizeChanger: true {{`}}`}}
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
