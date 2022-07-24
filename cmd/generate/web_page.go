package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

func getWebPageFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/src/pages/%s/index.js", dir, strings.ToLower(name))
	return fullname
}

func genWebPage(ctx context.Context, cmd *Command, item TplItem) error {
	data := map[string]interface{}{
		"Name":          item.StructName,
		"PluralName":    util.ToPlural(item.StructName),
		"Comment":       item.Comment,
		"UnderLineName": util.ToLowerUnderlinedNamer(item.StructName),
		"Fields":        item.Fields,
		"ReplaceLeft":   "{{",
		"ReplaceRight":  "}}",
	}

	buf, err := execParseTpl(webPageTpl, data)
	if err != nil {
		return err
	}

	fullname := getWebPageFileName(cmd.cfg.React, item.StructName)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const webPageTpl = `import React, { useRef, useState } from 'react';
import { Button, message, Tag, Tooltip } from 'antd';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import { BetaSchemaForm } from '@ant-design/pro-form';
import { {{.Name}} as API } from '@/services';
import Confirm from '@/pages/components/Confirm';
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { KeepAlive, useAccess } from 'umi';

const preProcessFields = fields => {
  if (fields.price > 0) {
    fields.price = parseInt(fields.price * 100);
  }
  return { ...fields };
};

const handleAdd = async (fields) => {
  const action = '添加';
  const hide = message.loading(` + "`正在$" + `{action}` + "`" + `);
  const msg = await API.create(preProcessFields(fields));
  hide();
  if (msg?.status === 'ok') {
    message.success(` + "`$" + `{action}成功` + "`" + `);
    return true;
  } else {
    message.error(` + "`$" + `{action}失败请重试！` + "`" + `);
    return false;
  }
};

const handleUpdate = async (fields) => {
  const action = '更新';
  const hide = message.loading(` + "`正在$" + `{action}` + "`" + `);
  const msg = await API.update(preProcessFields(fields));
  hide();
  if (msg?.status === 'ok') {
    message.success(` + "`$" + `{action}成功` + "`" + `);
    return true;
  } else {
    message.error(` + "`$" + `{action}失败请重试！` + "`" + `);
    return false;
  }
};

const handleRemove = async (selectedRows) => {
  const action = '删除';
  if (!selectedRows) return true;

  try {
    for (const { id, phone } of selectedRows) {
      let hide = message.loading(` + "`正在$" + `{action} ${phone}` + "`" + `);
      await API.remove({ id });
      hide();
    }
    message.success(` + "`$" + `{action}成功` + "`" + `);
    return true;
  } catch (error) {
    message.error(` + "`$" + `{action}失败，请重试` + "`" + `);
    return false;
  }
};

const TableList = props => {
  const { route: { name: headerTitle } } = props;
  const access = useAccess();
  const [layoutType] = useState('ModalForm');
  const [createModalVisit, setCreateModalVisit] = useState(false);
  const [updateModalVisit, setUpdateModalVisit] = useState(false);
  const [selectedRows, setSelectedRows] = useState([]);
  const [maskClosable, setMaskClosable] = useState(true);
  const [stepFormValues, setStepFormValues] = useState({});
  const [dataSource, setDataSource] = useState([]);

  const actionRef = useRef();

  const columns = [
    {
      title: '操作',
      hideInForm: true,
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => (<div style={{.ReplaceLeft}} whiteSpace: 'nowrap' {{.ReplaceRight}}>
        {/*["pink", "red", "yellow", "orange", "cyan", "green", "blue", "purple", "geekblue", "magenta", "volcano", "gold", "lime"]*/}
        <Tooltip title={'编辑'}>
          <Tag color={'blue'}>
            <a onClick={() => {
              setUpdateModalVisit(true);
              setStepFormValues(record);
            {{.ReplaceRight}}><EditOutlined style={{.ReplaceLeft}} color: '#91D5FF' {{.ReplaceRight}} /></a>
          </Tag>
        </Tooltip>
        <Tooltip title={'删除'}>
          <Tag color={'red'}>
            <a onClick={() => {
              const remove = async () => {
                const hide = message.loading('正在请求');
                const msg = await API.remove({ id: record.id });
                hide();
                if (msg?.status === 'ok') {
                  setDataSource(dataSource.filter((item) => {
                    return item.id !== record.id;
                  }));
                  message.success('请求成功');
                } else {
                  message.error('请求失败请重试！');
                }
              };
              return Confirm({
                title: '确认删除？',
                content: <>ID：{record.id}</>,
                confirmValue: record.id,
                placeholder: '输入“ID”确认删除',
                onOk: remove,
              });
            {{.ReplaceRight}}><DeleteOutlined style={{.ReplaceLeft}} color: 'red' {{.ReplaceRight}} /></a>
          </Tag>
        </Tooltip>
      </div>),
    },
    {
      title: '唯一标识',
      dataIndex: 'id',
    },` +
	"{{range .Fields}}" +
	"\n    {\n" +
	"      title: '{{.Comment}}',\n" +
	"      dataIndex: '{{fieldToLowerUnderlinedName .StructFieldName}}',\n" +
	"{{if eq .Condition false}}" +
	"      hideInSearch: true,\n" +
	"{{end}}" +
	"{{if .HideInTable}}" +
	"      hideInTable: true,\n" +
	"{{end}}" +
	"{{if .HideInForm}}" +
	"      hideInForm: true,\n" +
	"{{end}}" +
	"{{if .StructFieldRequired}}" +
	"      formItemProps: {\n" +
	"        rules: [{ required: true }],\n" +
	"      },\n" +
	"{{end}}" +
	"{{if .ValueType}}" +
	"      valueType: '{{.ValueType}}',\n" +
	"{{else if eq \"int\" .StructFieldType }}" +
	"      valueType: 'digit',\n" +
	"{{else if eq \"bool\" .StructFieldType }}" +
	"      valueType: 'switch',\n" +
	"{{end}}" +
	"    }," +
	"{{end}}" + `
  ];


  return (
    <PageContainer>
      <ProTable
        headerTitle={headerTitle}
        actionRef={actionRef}
        rowKey='id'
        pagination={{.ReplaceLeft}}
          pageSize: 10,
          position: ['bottomCenter'],
        {{.ReplaceRight}}
        form={{.ReplaceLeft}}
          syncToUrl: (values, type) => {
            if (type === 'get') {
              return { ...values };
            }
          },
        {{.ReplaceRight}}
        columnsState={{.ReplaceLeft}}
          persistenceKey: 'pro-table-client',
          persistenceType: 'localStorage',
        {{.ReplaceRight}}
        search={{.ReplaceLeft}} labelWidth: 'auto' {{.ReplaceRight}}
        toolBarRender={() => {
          return [
            (access.canClient ? <Button type='primary' key='primary' onClick={() => {
              setCreateModalVisit(true);
            {{.ReplaceRight}}><PlusOutlined /> 新建</Button> : null),
          ];
        {{.ReplaceRight}}
        dataSource={dataSource}
        onLoad={async dataSource => {
          setDataSource(dataSource.map(item => {
            item.price = (item.price / 100).toFixed(2);
            return item;
          }));
        {{.ReplaceRight}}
        request={API.query}
        columns={columns}
        tableAlertRender={false}
        rowSelection={access.canClient ? { onChange: (_, selectedRows) => setSelectedRows(selectedRows) } : false}
      />

      <BetaSchemaForm
        title='新建'
        visible={createModalVisit}
        onVisibleChange={setCreateModalVisit}
        modalProps={{.ReplaceLeft}}
          maskClosable: maskClosable,
          destroyOnClose: true,
        {{.ReplaceRight}}
        width={600}
        layoutType={layoutType}
        onFinish={async values => {
          const success = await handleAdd(values);
          if (success) {
            setCreateModalVisit(false);
            actionRef.current?.reload?.();
          }
        {{.ReplaceRight}}
        columns={columns.filter(item => {
          return item.dataIndex !== 'id';
        })}
      />

      {stepFormValues?.id && (
        <BetaSchemaForm
          title='编辑'
          visible={updateModalVisit}
          onVisibleChange={setUpdateModalVisit}
          modalProps={{.ReplaceLeft}}
            maskClosable: maskClosable,
            destroyOnClose: true,
            onClose: () => {
              setMaskClosable(true);
            },
          {{.ReplaceRight}}
          initialValues={{.ReplaceLeft}}
            ...stepFormValues,
            district: [stepFormValues.province_id, stepFormValues.city_id, stepFormValues.district_id, stepFormValues.street_id].filter(val => val > 0),
          {{.ReplaceRight}}
          onValuesChange={() => {
            setMaskClosable(false);
          {{.ReplaceRight}}
          width={600}
          layoutType={layoutType}
          onFinish={async values => {
            const success = await handleUpdate(values);
            if (success) {
              setMaskClosable(true);
              setUpdateModalVisit(false);
              setDataSource(dataSource.map((item) => {
                if (item.id === values.id) {
                  return { ...item, ...values };
                }
                return item;
              }));
            }
          {{.ReplaceRight}}
          columns={columns.map(item => {
            const readonlyFields = ['id', 'phone'];
            if (readonlyFields.indexOf(item.dataIndex) > -1) {
              return { ...item, readonly: true };
            }
            return item;
          })}
        />)}

    </PageContainer>
  );
};

export default props => {
  const { location, route: { name } } = props;
  return <KeepAlive name={name} location={location} saveScrollPosition='screen'>
    <TableList {...props} />
  </KeepAlive>;
};
`
