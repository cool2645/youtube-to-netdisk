function prefixInteger(num, n) {
    return (Array(n).join(0) + num).slice(-n);
}

function formatDateTime(time) {
    let datetime = new Date(time);
    let year = datetime.getFullYear();
    let month = prefixInteger(datetime.getMonth() + 1, 2);
    let date = prefixInteger(datetime.getDate(), 2);
    let hour = prefixInteger(datetime.getHours(), 2);
    let minute = prefixInteger(datetime.getMinutes(), 2);
    let second = prefixInteger(datetime.getSeconds(), 2);
    return year + "-" + month + "-" + date + " " + hour + ":" + minute + ":" + second + "";
}

function formatDateTimeFromDatetimeString(datetime) {
    formatDateTime(new Date(datetime).getTime())
}

export {formatDateTime, formatDateTimeFromDatetimeString}
export default formatDateTime