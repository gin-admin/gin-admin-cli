{{- $name := .Name}}
{{- $includeStatus := .Include.Status}}
{{- $statusEnabledText := .Extra.StatusEnabledText}}
{{- $statusDisabledText :=.Extra.StatusDisabledText}}
import React, { useEffect, useRef } from 'react';
import { message } from 'antd';
import {
  ModalForm,
  ProFormText,
  {{- if $includeStatus}}
  ProFormSwitch,
  {{- end}}
} from '@ant-design/pro-components';
import type { ProFormInstance } from '@ant-design/pro-components';
import { add{{$name}}, get{{$name}}, update{{$name}} } from '@/services/{{.Extra.ImportService}}';

type {{$name}}ModalProps = {
  onSuccess: () => void;
  onCancel: () => void;
  visible: boolean;
  title: string;
  id?: string;
};

const {{$name}}Modal: React.FC<{{$name}}ModalProps> = (props: {{$name}}ModalProps) => {
  const formRef = useRef<ProFormInstance<API.{{$name}}>>();

  useEffect(() => {
    if (!props.visible) {
      return;
    }

    formRef.current?.resetFields();
    if (props.id) {
      get{{$name}}(props.id).then(async (res) => {
        if (res.data) {
          const data = res.data;
          {{- if $includeStatus}}
          data.statusChecked = data.status === 'enabled';
          {{- end}}
          formRef.current?.setFieldsValue(data);
        }
      });
    }
  }, [props]);

  return (
    <ModalForm<API.{{$name}}>
      visible={props.visible}
      title={props.title}
      width={800}
      formRef={formRef}
      layout="horizontal"
      grid={true}
      submitTimeout={3000}
      submitter={{`{{`}}
        searchConfig: {
          submitText: '{{.Extra.SubmitText}}',
          resetText: '{{.Extra.ResetText}}',
        },
      {{`}}`}}
      modalProps={{`{{`}}
        destroyOnClose: true,
        maskClosable: false,
        onCancel: () => {
          props.onCancel();
        },
      {{`}}`}}
      onFinish={async (values: API.{{$name}}) => {
        {{- if $includeStatus}}
        values.status = values.statusChecked ? 'enabled' : 'disabled';
        delete values.statusChecked;
        {{- end}}

        if (props.id) {
          await update{{$name}}(props.id, values);
        } else {
          await add{{$name}}(values);
        }

        message.success('{{.Extra.SaveSuccessMessage}}');
        props.onSuccess();
        return true;
      }}
      initialValues={{`{{}}`}}
    >
      {{- range .Fields}}
      {{- $fieldName := .Name}}
      {{- $required := .Extra.Required}}
      {{- $fieldLabel := .Extra.Label}}
      {{- $fieldPlaceholder := .Extra.Placeholder}}
      {{- $fieldRulesMessage := .Extra.RulesMessage}}
      {{- if .Form}}
      {{- if eq $fieldName "Status"}}
      <ProFormSwitch
        name="statusChecked"
        label="{{$fieldLabel}}"
        fieldProps={{`{{`}}
          checkedChildren: '{{$statusEnabledText}}',
          unCheckedChildren: '{{$statusDisabledText}}',
        {{`}}`}}
        colProps={{`{{`}} span: 12 {{`}}`}}
      />
      {{- else}}
      <ProFormText
        name="{{lowerUnderline $fieldName}}"
        label="{{$fieldLabel}}"
        {{- if $fieldPlaceholder}}
        placeholder="{{$fieldPlaceholder}}"
        {{- end}}
        {{- if $required}}
        rules={[
          {
            required: {{$required}},
            message: '{{$fieldRulesMessage}}',
          },
        ]}
        {{- end}}
        colProps={{`{{`}} span: 12 {{`}}`}}
      />
      {{- end}}
      {{- end}}
      {{- end}}
    </ModalForm>
  );
};

export default {{$name}}Modal;
