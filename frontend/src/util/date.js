import moment from "moment";


// 将 YYYYMMDD 格式的字符串转换为 Date 对象
const parseDate = (dateStr) => {
    return moment(dateStr, "YYYYMMDD", true).valueOf();
};

// 将 Date 对象转换为 YYYYMMDD 格式的字符串
const formatDate = (date) => {
   return moment(date).format("YYYYMMDD");
};

/**
 *
 *
 * @param pre  format(YYYYMMDD)
 * @param after format(YYYYMMDD)
 * @returns {boolean}
 * @constructor
 */
 const DayBeforeExceedOneDay=(pre,after)=> {
    let  preDate =parseDate(pre);
    let afterDate = parseDate(after);
    return afterDate - preDate > 0 && ( afterDate-preDate) / (1000 * 60 * 60 * 24) >= 1;
}
/**
 * 将原来的时间的基础上增加days天，
 * @param now format(YYYYMMDD)
 * @param days
 * @constructor
 */
const  AddDays = (now,days)=>{
    let preData  = parseDate(now);
    // 添加天数（转换为毫秒数）
    preData += days * (1000 * 60 * 60 * 24);

    // 创建新的日期对象并格式化为 YYYYMMDD
    return moment(preData).format("YYYYMMDD");
}

 const DayAfterExceedOneDay=(pre,after)=> {
    let preDate =parseDate(pre);
    let afterDate = parseDate(after);

}


const visibleDate = (dateStr) => {
    // 解析输入的日期字符串
    const date = moment(dateStr, "YYYYMMDD", true);

    // 检查日期是否有效
    if (!date.isValid()) {
        console.error("Invalid date format provided");
        return null;
    }

    // 格式化为 'yyyy年mm月dd日'
    return date.format("YYYY年MM月DD日");
};


const SubDays=(now,days)=>{
     let preData  = parseDate(now);
     preData -= days*(1000 * 60 * 60 * 24);
    return moment(preData).format("YYYYMMDD");
}

function isToday(timestamp) {
    const date = new Date(timestamp);

    const today = new Date();
    const currentYear = today.getFullYear();
    const currentMonth = today.getMonth();
    const currentDate = today.getDate();

    const timestampYear = date.getFullYear();
    const timestampMonth = date.getMonth();
    const timestampDate = date.getDate();
    return currentYear === timestampYear && currentMonth === timestampMonth && currentDate === timestampDate;
}

export {DayBeforeExceedOneDay,DayAfterExceedOneDay,AddDays,SubDays,visibleDate,isToday};