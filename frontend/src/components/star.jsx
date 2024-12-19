import {Rate} from "antd";
import React from "react";


const visibleWeight = {
    1:'D',
    2:'C',
    3:'B',
    4:'A',
    5:'S'
}


const Star = ({number})=>{
    return (
        <Rate count={1} style={{ color:'#000',fontFamily: "sofia-pro, sans-serif", fontWeight: 400,fontStyle: "normal"}} value={1}  character={visibleWeight[Math.floor(number)]?visibleWeight[Math.floor(number)]:'D'} allowHalf disabled style={{ fontSize: 36 }} />
    )
}


export default Star;