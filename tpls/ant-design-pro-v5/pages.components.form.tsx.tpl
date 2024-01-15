{{- $name := .Name}}
{{- $lowerCamelName := lowerCamel .Name}}
{{- $parentName := .Extra.ParentName}}
import React, { useEffect, useRef } from 'react';
import { useIntl } from '@umijs/max';
import { message{{if .Extra.FormAntdImport}}, {{.Extra.FormAntdImport}}{{end}} } from 'antd';
import { ModalForm{{if .Extra.FormProComponentsImport}}, {{.Extra.FormProComponentsImport}}{{end}} } from '@ant-design/pro-components';
import type { ProFormInstance } from '@ant-design/pro-components';
import { add{{$name}}, get{{$name}}, update{{$name}} } from '@/services/{{$parentName}}/{{$lowerCamelName}}';

type {{$name}}ModalProps = {
  onSuccess: () => void;
  onCancel: () => void;
  visible: boolean;
  title: string;
  id?: string;
};

const {{$name}}Modal: React.FC<{{$name}}ModalProps> = (props: {{$name}}ModalProps) => {
  const intl = useIntl();
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
          formRef.current?.setFieldsValue(data);
        }
      });
    }
  }, [props]);

  return (
    <ModalForm<API.{{$name}}>
      open={props.visible}
      title={props.title}
      width={800}
      formRef={formRef}
      layout="horizontal"
      grid={true}
      rowProps={{`{{`}} gutter: 20 {{`}}`}}
      submitTimeout={3000}
      submitter={{`{{`}}
        searchConfig: {
          submitText: intl.formatMessage({ id: 'button.confirm' }),
          resetText: intl.formatMessage({ id: 'button.cancel' }),
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
        if (props.id) {
          await update{{$name}}(props.id, values);
        } else {
          await add{{$name}}(values);
        }

        message.success(intl.formatMessage({ id: 'component.message.success.save' }));
        props.onSuccess();
        return true;
      {{`}}`}}
      initialValues={{`{{ }}`}}
    >
      {{- range .Fields}}
      {{- if .Form}}
      {{- if .Extra.FormComponent}}
      {{.Extra.FormComponent}}
      {{- end}}
      {{- end}}
      {{- end}}
    </ModalForm>
  );
};

export default {{$name}}Modal;
