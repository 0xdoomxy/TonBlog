import React, { useContext, useEffect, useRef, useState } from 'react';
import { Form, Input, Modal, Progress, Table } from 'antd';
import { motion } from 'framer-motion';
import { AirportClient } from '../agent/agent';
const EditableContext = React.createContext(null);
const EditableRow = ({ index, ...props }) => {
    const [form] = Form.useForm();
    return (
        <Form form={form} component={false}>
            <EditableContext.Provider value={form}>
                <tr {...props} />
            </EditableContext.Provider>
        </Form>
    );
};
const EditableCell = ({
    title,
    editable,
    children,
    dataIndex,
    record,
    handleSave,
    ...restProps
}) => {
    const [editing, setEditing] = useState(false);
    const inputRef = useRef(null);
    const form = useContext(EditableContext);
    useEffect(() => {
        if (editing) {
            inputRef.current?.focus();
        }
    }, [editing]);
    const toggleEdit = () => {
        setEditing(!editing);
        form.setFieldsValue({
            [dataIndex]: record[dataIndex],
        });
        console.log(record);
    };
    const save = async () => {
        try {
            const values = await form.validateFields();
            toggleEdit();
            handleSave({
                ...record,
                ...values,
            });
            console.log(values);
        } catch (errInfo) {
            console.log('Save failed:', errInfo);
        }
    };
    let childNode = children;
    console.log(childNode);
    if (editable) {
        childNode = editing ? (
            <Form.Item
                style={{
                    margin: 0,
                }}
                name={dataIndex}
                rules={[
                    {
                        required: true,
                        message: `${title} is required.`,
                    },
                ]}
            >
                <Input ref={inputRef} onPressEnter={save} onBlur={save} />
            </Form.Item>
        ) : (
            <div
                className="editable-cell-value-wrap"
                style={{
                    paddingInlineEnd: 24,
                }}
                onClick={toggleEdit}
            >
                {children}
            </div>
        );
    }
    return <td {...restProps}>{childNode}</td>;
};

const RunningAirport = (props) => {
    const { isAdmin } = props;
    const [dataSource, setDataSource] = useState(null);
    useEffect(()=>{
       
    },[])
    const findRunningAirport=(page,pageSize)=>{
        AirportClient.FindRunningAirport(page,pageSize).then((data)=>{
            if (data === undefined || data === null) {
                return;
            }
            if (!data.status) {
                let msg = data.message;
                if (msg === undefined || msg === null) {
                    msg = "系统出错啦";
                }
                toast.error(msg);
                return;
            }
            if (data.data === undefined || data.data === null ) {
                return;
            }
        })
    }
    const handleAdd = ()=>{

    }
    //TODO
    const handleDelete = (key) => {
        const newData = dataSource.filter((item) => item.key !== key);
        setDataSource(newData);
    };
    //TODO
    const handleComplete = (key) => {
        const newData = dataSource.filter((item) => item.key !== key);
        setDataSource(newData);
    }
    const defaultColumns = [
        {
            title: '进度',
            render: (_, record) => {
                let start = record.start_time;
                let end = record.end_time;
                let now = Date.now();
                let p = Math.floor((now - start) / (end - start) * 100);
                return (
                    <Progress percent={p} percentPosition={{ align: 'center', type: 'inner' }} size={[100, 20]} />
                );
            }
        },
        {
            title: '项目名',
            dataIndex: 'name',
            editable: isAdmin,
        },
        {
            title: '官网地址',
            dataIndex: 'address',
            editable: isAdmin,
            render: (_, record) => {
                return <a href={record.address}>官网地址</a>
            }
        },
        {
            title: '赛道',
            dataIndex: 'tag',
            editable: isAdmin,
        },
        {
            title: '融资金额',
            dataIndex: 'financing_balance',
            editable: isAdmin,
        },
        {
            title: '融资来源方',
            dataIndex: 'financing_from',
            editable: isAdmin,
        },
        {
            title: '教程',
            dataIndex: 'teaching',
            editable: isAdmin,
            render: (_, record) => {
                return <a href={record.teaching}>教程链接</a>
            }
        },
        {
            title: '任务类型',
            dataIndex: 'task_type',
            editable: isAdmin
        },
        {
            title: '进展',
            dataIndex: 'operation',
            render: (_, record) =>
                dataSource.length >= 1 ? (
                    <div className={"w-full justify-center items-center flex-col"}>
                        <motion.button whileHover={{ scale: 1.1 }}
                            whileTap={{ scale: 0.9 }}
                            style={{ width: "80px", height: "40px" }}
                            transition={{ type: "spring", stiffness: 400, damping: 10 }}
                            className={"motion-button  px-1"} title="今日完成"
                            key={record.key}
                            onClick={() => handleDelete(record.key)}
                        >
                            <a>今日完成</a>
                        </motion.button>
                        {isAdmin &&
                            <motion.button whileHover={{ scale: 1.1 }}
                                whileTap={{ scale: 0.9 }}
                                transition={{ type: "spring", stiffness: 400, damping: 10 }}
                                className={"motion-button  px-1"} title="结束空投"
                                style={{ width: "80px", height: "40px" }}
                                key={record.key}
                                onClick={() => handleComplete(record.key)}>
                                <a>结束空投</a>
                            </motion.button>}
                        {isAdmin &&
                            <motion.button whileHover={{ scale: 1.1 }}
                                whileTap={{ scale: 0.9 }}
                                transition={{ type: "spring", stiffness: 400, damping: 10 }}
                                className={"motion-button  px-1"} title="删除空投"
                                style={{ width: "80px", height: "40px" }}
                                key={record.key}
                                onClick={() => handleDelete(record.key)}>
                                <a>删除空投</a>
                            </motion.button>}
                    </div>

                ) : null,
        },
    ];
    const handleSave = (row) => {
        console.log(row)
        const newData = [...dataSource];
        const index = newData.findIndex((item) => row.key === item.key);
        const item = newData[index];
        newData.splice(index, 1, {
            ...item,
            ...row,
        });
        setDataSource(newData);
    };
    const components = {
        body: {
            row: EditableRow,
            cell: EditableCell,
        },
    };
    const columns = defaultColumns.map((col) => {
        if (!col.editable) {
            return col;
        }
        return {
            ...col,
            onCell: (record) => ({
                record,
                editable: col.editable,
                dataIndex: col.dataIndex,
                title: col.title,
                handleSave,
            }),
        };
    });
    return (<>
        <Modal></Modal>
        <div className={"w-full h-full flex justify-center items-center flex-col"}>
            {isAdmin && <div className={"w-full h-full"}>
                <div className={"w-full items-center justify-start flex pb-4 pl-4"}>
                    <motion.button
                        className={"motion-button  flex justify-center items-center  md:text-md lg:text-xl text-white"}
                        whileHover={{ scale: 1.1 }}
                        whileTap={{ scale: 0.9 }}
                        transition={{ type: "spring", stiffness: 400, damping: 10 }}
                        onClick={handleAdd}
                    >
                        新增空投
                    </motion.button>
                </div>
            </div>
            }
            <div>
                <Table
                   key={"running"}
                    tableLayout={"auto"}
                    components={components}
                    rowClassName={() => 'editable-row'}
                    bordered
                    className={"w-full flex justify-center items-center h-full"}
                    dataSource={dataSource}
                    columns={columns}
                />
            </div>
        </div>
        </>
    )
}

export default RunningAirport;