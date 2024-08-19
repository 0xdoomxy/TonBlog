import moment from "moment";


function FormatDateDistance(date) {
    let d = moment(date);
    return d.fromNow();
}

export {FormatDateDistance};