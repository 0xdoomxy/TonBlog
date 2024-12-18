
export function debounce(fn, wait) {
    let timer;
    return function(...args) {
        let context = this;
        if (timer) {
            clearTimeout(timer);
        }
        timer = setTimeout(() => fn.call(context, ...args), wait);
    };
}
