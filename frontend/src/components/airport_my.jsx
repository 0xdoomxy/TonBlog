import React, {useContext, useEffect, useRef, useState} from 'react';
import {Form, Input, Progress, Skeleton, Statistic, Table, Tag, Tooltip} from 'antd';
import {motion} from 'framer-motion';
import {CheckCircleOutlined, ExclamationCircleOutlined, InfoCircleOutlined, SyncOutlined} from "@ant-design/icons";
import {Star} from "./index.js";
import {isToday} from "../util/date.js";

const EditableContext = React.createContext(null);
const EditableRow = ({index, ...props}) => {
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
    };
    const save = async () => {
        try {
            const values = await form.validateFields();
            toggleEdit();
            handleSave({
                ...record,
                ...values,
            });
        } catch (errInfo) {
            console.log('Save failed:', errInfo);
        }
    };
    let childNode = children;
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
                <Input ref={inputRef} onPressEnter={save} onBlur={save}/>
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
const MyAirport = (props) => {
    const {isAdmin} = props;
    const [dataSource, setDataSource] = useState([
        {
            key: '0',
            name: 'Edward King 0',
            start_time: Date.now() - 1000 * 1000 * 60 * 24,
            end_time: Date.now() - 500 * 1000 * 60 * 24,
            final_time: Date.now() + 2000 * 1000 * 60 * 24,
            address: "www.baidu.com",
            tag: "区块链,AI",
            financing_balance: "3000$",
            financing_from: "a16z,binance",
            task_type: "拉人头,交互",
            balance: 1000,
            status: "SUCCESS"
        },
        {
            key: '1',
            name: 'Edward King 0',
            start_time: Date.now() - 1000 * 60 * 24,
            end_time: Date.now() - 250 * 1000 * 60 * 24,
            final_time: Date.now() + 2000 * 1000 * 60 * 24,
            address: "www.baidu.com",
            tag: "区块链,AI",
            financing_balance: "3000$",
            financing_from: "a16z,binance",
            task_type: "拉人头,交互",
            status: "PROCESSING",
            update_time: Date.now() - 250 * 1000 * 60 * 24
            // balance: 1000
        },
        {
            key: '1',
            name: 'Edward King 0',
            start_time: Date.now() - 1000 * 60 * 24,
            end_time: Date.now() - 250 * 1000 * 60 * 24,
            final_time: Date.now() + 2000 * 1000 * 60 * 24,
            address: "www.baidu.com",
            tag: "区块链,AI",
            financing_balance: "3000$",
            financing_from: "a16z,binance",
            task_type: "拉人头,交互",
            status: "PROCESSING",
            update_time: Date.now() -  1000 * 60 * 24
            // balance: 1000
        },
        {
            key: '2',
            name: 'Edward King 0',
            start_time: Date.now() - 1000 * 60 * 24,
            end_time: Date.now() - 250 * 1000 * 60 * 24,
            final_time: Date.now() + 2000 * 1000 * 60 * 24,
            address: "www.baidu.com",
            tag: "区块链,AI",
            financing_balance: "3000$",
            financing_from: "a16z,binance",
            task_type: "拉人头,交互",
            status: "UNCHECK"
            // balance: 1000
        },
    ]);
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
            title: "状态",
            align: "center",
            dataIndex: "status",
            render: (_, record) => {
                switch (record.status) {
                    case 'SUCCESS':
                        return (
                            <Tag icon={<CheckCircleOutlined/>} color="success">
                                空投领取
                            </Tag>
                        )
                    case 'PROCESSING':
                        return (
                            <Tag icon={<SyncOutlined spin/>} color="processing">
                                进行中
                            </Tag>
                        )
                    case 'UNCHECK':
                        return (
                            <Tag icon={<ExclamationCircleOutlined/>} color="warning">
                                待领取
                            </Tag>
                        )
                }
            }
        },
        {
            title: '进度',
            align: "center",
            render: (_, record) => {
                let end = record.end_time;
                let final = record.final_time;
                let now = Date.now();
                let p = Math.floor((now - end) / (final - end) * 100);
                return (
                    <Progress
                        format={(percent) => `${record.status === "PROCESSING" ? "空投进度:" : "领取进度:"} ${percent}%`}
                        percent={p} percentPosition={{align: 'center', type: 'outer'}} size={[100, 30]}/>
                );
            }
        },
        {
            title: '项目名',
            dataIndex: 'name',
            editable: true,
            align: "center",
        },
        {
            title: '官网地址',
            dataIndex: 'address',
            editable: true,
            align: "center",
            render: (_, record) => {
                return <a href={record.address}>官网地址</a>
            }
        },
        {
            title: '赛道',
            dataIndex: 'tag',
            editable: true,
            align: "center",
            render: (_, record) => {
                return <div className={"flex  justify-center items-center"}>
                    {record.tag.split(',').map((tag) => {
                        let color = tag.length > 5 ? 'geekblue' : 'green';
                        if (tag === 'loser') {
                            color = 'volcano';
                        }
                        return (
                            <Tag color={color} key={tag}>
                                {tag.toUpperCase()}
                            </Tag>
                        );
                    })}
                </div>
            }
        },
        {
            title: '融资金额',
            dataIndex: 'financing_balance',
            editable: true,
            align: "center",
        },
        {
            title: '融资来源方',
            dataIndex: 'financing_from',
            editable: true,
            align: "center",
            render: (_, record) => {
                return <div className={"flex  justify-center items-center"}>
                    {record.financing_from.split(',').map((tag) => {
                        let color = tag.length > 5 ? 'geekblue' : 'volcano';
                        return (
                            <Tag color={color} key={tag}>
                                {tag.toUpperCase()}
                            </Tag>
                        );
                    })}
                </div>
            }
        },
        {
            title: '教程',
            dataIndex: 'teaching',
            editable: true,
            align: "center",
            render: (_, record) => {
                return <a href={record.teaching}>教程链接</a>
            }
        },
        {
            title: <Tooltip placement={"rightTop"} color={"rgba(116,112,112,0.88)"}
                            title={"该空投在平台收集的空投中的评分"}>空投质量<InfoCircleOutlined
                className={"relative  bottom-3 left-2"}/></Tooltip>,
            dataIndex: 'weight',
            align: "center",
            render: (_, record) => {
                return (
                    <Star number={record.weight}/>
                )
            }
        },
        {
            title: '任务类型',
            dataIndex: 'task_type',
            editable: true,
            align: "center",
            render: (_, record) => {
                return <div className={"flex  justify-center items-center"}>
                    {record.task_type.split(',').map((tag) => {
                        let color = tag.length > 5 ? 'magenta' : 'purple';
                        return (
                            <Tag color={color} key={tag}>
                                {tag.toUpperCase()}
                            </Tag>
                        );
                    })}
                </div>
            }
        },
        {
            title: <Tooltip placement={"rightTop"} color={"rgba(116,112,112,0.88)"}
                            title={"平台用户获取该空投的总数量"}>空投数量<InfoCircleOutlined
                className={"relative  bottom-3 left-2"}/></Tooltip>,
            dataIndex: "balance",
            editable: true,
            align: "center",
            render: (_, record) => (record.status === "SUCCESS"&&record.balance)? (
                <Statistic  value={record.balance}/>
            ):<Skeleton paragraph={{
            rows: 1,
        }} active />,
            onCell: (item) => {
                if (item.status === "SUCCESS") {
                    return {colSpan: 2};
                } else {
                    return {colSpan: 0};
                }

            }
        },
        {
            title: '进展',
            dataIndex: 'operation',
            align: "center",
            onCell: (record) => {
                if (record.status === "SUCCESS") {
                    return {colSpan: 0};
                }
                if (record.status === "PROCESSING" && isToday(record.update_time)){
                    return {
                        colSpan: 0,
                    }
                }
                return {
                    colSpan:2,
                }
            },
            render: (_, record) =>
                dataSource.length >= 1 && (<>{
                        record.status === "PROCESSING" && !isToday(record.update_time) &&
                        <motion.button whileHover={{scale: 1.1}}
                                       whileTap={{scale: 0.9}}
                                       transition={{type: "spring", stiffness: 400, damping: 10}}
                                       className={"motion-button  px-1"} title="完成任务"
                                       style={{width: "80px", height: "40px" ,background:"#b7eb8f"}}
                                       key={record.key}
                                       onClick={() => handleComplete(record.key)}>
                            <a>今日完成</a>
                        </motion.button>
                    }
                        {record.status === "UNCHECK" && <motion.button whileHover={{scale: 1.1}}
                                                                       whileTap={{scale: 0.9}}
                                                                       transition={{
                                                                           type: "spring",
                                                                           stiffness: 400,
                                                                           damping: 10
                                                                       }}
                                                                       className={"motion-button  px-1"} title="领取空投"
                                                                       style={{width: "80px", height: "40px",background:"#faad14"}}
                                                                       key={record.key}
                                                                       onClick={() => handleComplete(record.key)}>
                            <a>领取</a>
                        </motion.button>}
                    </>
                )
        },
    ];
    const handleAdd = () => {

        setDataSource([...dataSource, newData]);
    };
    const handleSave = (row) => {
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
    return (
        <div className={"w-full h-full flex justify-center items-center flex-col"}>
            <div>
                <Table
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
    )
}
export default MyAirport;