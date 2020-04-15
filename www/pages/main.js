Vue.prototype.mainBus = new Vue();

Vue.config.productionTip = false;

Vue.directive('tooltip', function (el, binding) {
    $(el).tooltip({
        title: binding.value,
        placement: binding.arg,
        trigger: 'hover'
    })
})

Vue.filter('prettyBytes', (num) => {
    // jacked from: https://github.com/sindresorhus/pretty-bytes
    if (typeof num !== "number" || isNaN(num)) {
        throw new TypeError("Expected a number");
    }

    var exponent;
    var unit;
    var neg = num < 0;
    var units = ["B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

    if (neg) {
        num = -num;
    }

    if (num < 1) {
        return (neg ? "-" : "") + num + " B";
    }

    exponent = Math.min(
        Math.floor(Math.log(num) / Math.log(1024)),
        units.length - 1
    );
    num = (num / Math.pow(1024, exponent)).toFixed(2) * 1;
    unit = units[exponent];

    return (neg ? "-" : "") + num + " " + unit;
});


new Clipboard('.btn-copy');
