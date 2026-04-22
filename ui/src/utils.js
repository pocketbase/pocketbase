const DEFAULT_RANDOM_ALPHABET = "abcdefghijklmnopqrstuvwxyz0123456789";

const REMEMBER_PATH_KEY = "pb_redirect";

const navigationStore = store({
    hash: window.location.hash,
});
window.addEventListener("hashchange", () => {
    navigationStore.hash = window.location.hash;
});

// https://prismjs.com/extending.html
// https://prismjs.com/tokens
Prism.languages.pbrule = {
    "string": Prism.languages.js.string,
    "number": Prism.languages.js.number,
    "function": Prism.languages.js.function,
    "boolean": /\b(?:true|false)\b/i,
    "constant": /\b(?:null)\b/i,
    "comment": {
        pattern: /\/\/.*/,
        greedy: true,
    },
    "italic": /_via_|\:\w+/,
    "keyword": /&&|\|\||\??(?:!~|!=|>=|<=|=|~|>|>|<)(?=[@\w\s]|$)/,
};

const utils = {
    /**
     * Checks whether value is plain object.
     *
     * @param  {Mixed} value
     * @return {boolean}
     */
    isObject(value) {
        return value !== null && typeof value === "object" && value.constructor === Object;
    },

    /**
     * Checks whether a value is empty. The following values are considered as empty:
     * - null
     * - undefined
     * - empty string
     * - empty array
     * - empty object
     *
     * @param  {Mixed} value
     * @return {boolean}
     */
    isEmpty(value) {
        return (
            value == null
            || value === ""
            || (Array.isArray(value) && value.length === 0)
            || (typeof value === "object" && app.utils.isEmptyObject(value))
        );
    },

    /**
     * Checks if an object doesn't have any properties.
     *
     * @param  {object obj
     * @return {boolean}
     */
    isEmptyObject(obj) {
        for (let i in obj) {
            return false;
        }
        return true;
    },

    /**
     * Normalizes and returns arr as a new array instance.
     *
     * @param  {Array}   arr
     * @param  {boolean} [allowEmpty]
     * @return {Array}
     */
    toArray(arr, allowEmpty = false) {
        if (Array.isArray(arr)) {
            return arr.slice();
        }

        return (allowEmpty || !utils.isEmpty(arr)) && typeof arr !== "undefined" ? [arr] : [];
    },

    /**
     * Removes single element from array by loosely comparying values.
     *
     * @param {Array} arr
     * @param {Mixed} value
     */
    removeByValue(arr, value) {
        if (!Array.isArray(arr)) {
            console.warn("[removeByValue] not an array:", arr);
            return;
        }

        for (let i = arr.length - 1; i >= 0; i--) {
            if (arr[i] == value) {
                arr.splice(i, 1);
                break;
            }
        }
    },

    /**
     * Removes single element from an objects array by matching a property value.
     *
     * @param {Array}  objectsArr
     * @param {string} key
     * @param {Mixed}  value
     */
    removeByKey(objectsArr, key, value) {
        if (!Array.isArray(objectsArr)) {
            console.warn("[removeByKey] not an array:", objectsArr);
            return;
        }

        for (let i in objectsArr) {
            if (objectsArr[i][key] == value) {
                objectsArr.splice(i, 1);
                break;
            }
        }
    },

    /**
     * Adds `value` in `arr` only if it's not added already.
     *
     * @param {Array} arr
     * @param {Mixed} value
     */
    pushUnique(arr, value) {
        if (!Array.isArray(arr)) {
            console.warn("[pushUnique] not an array:", arr);
            return;
        }

        if (!arr.includes(value)) {
            arr.push(value);
        }
    },

    /**
     * Merges all `valuesArr` items that don't exist in `targetArr`.
     *
     * @param {Array} targetArr
     * @param {Array} valuesArr
     */
    mergeUnique(targetArr, valuesArr) {
        for (let v of valuesArr) {
            app.utils.pushUnique(targetArr, v);
        }

        return targetArr;
    },

    /**
     * Adds or replace an object array element by comparing its key value.
     *
     * @param {Array}  objectsArr
     * @param {Object} item
     * @param {string} [key]
     */
    pushOrReplaceObject(objectsArr, item, key = "id") {
        for (let i = objectsArr.length - 1; i >= 0; i--) {
            if (objectsArr[i][key] == item[key]) {
                objectsArr[i] = item;
                return;
            }
        }

        objectsArr.push(item);
    },

    /**
     * Filters and returns a new objects array with duplicated elements removed.
     *
     * @param  {Array} objectsArr
     * @param  {string} key
     * @return {Array}
     */
    filterDuplicatesByKey(objectsArr, key = "id") {
        objectsArr = Array.isArray(objectsArr) ? objectsArr : [];

        const uniqueMap = {};

        for (const item of objectsArr) {
            uniqueMap[item[key]] = item;
        }

        return Object.values(uniqueMap);
    },

    /**
     * Filters and returns a new object with removed redacted props.
     *
     * @param  {Object} obj
     * @param  {string} [mask] Default to "******"
     * @return {Object}
     */
    filterRedactedProps(obj, mask = "******") {
        const result = JSON.parse(JSON.stringify(obj || {}));

        for (let prop in result) {
            if (typeof result[prop] === "object" && result[prop] !== null) {
                result[prop] = utils.filterRedactedProps(result[prop], mask);
            } else if (result[prop] === mask) {
                delete result[prop];
            }
        }

        return result;
    },

    /**
     * Safely access nested object/array key with dot-notation.
     *
     * @example
     * ```javascript
     * let myObj = {a: {b: {c: 3}}}
     * this.getByPath(myObj, "a.b.c");       // returns 3
     * this.getByPath(myObj, "a.b.c.d");     // returns null
     * this.getByPath(myObj, "a.b.c.d", -1); // returns -1
     * ```
     *
     * @param  {Object|Array} data
     * @param  {string}       path
     * @param  {Mixed}        [defaultVal]
     * @param  {string}       [delimiter]
     * @return {Mixed}
     */
    getByPath(data, path, defaultVal = null, delimiter = ".") {
        let result = data || {};
        let parts = (path || "").split(delimiter);

        for (const part of parts) {
            if ((!utils.isObject(result) && !Array.isArray(result)) || typeof result[part] === "undefined") {
                return defaultVal;
            }

            result = result[part];
        }

        return result;
    },

    /**
     * Sets a new value to an object (or array) by its key path.
     *
     * @example
     * ```javascript
     * this.setByPath({}, "a.b.c", 1);             // results in {a: b: {c: 1}}
     * this.setByPath({a: {b: {c: 3}}}, "a.b", 4); // results in {a: {b: 4}}
     * ```
     *
     * @param  {Array|Object} data
     * @param  {string}       path
     * @param  {string}       delimiter
     */
    setByPath(data, path, newValue, delimiter = ".") {
        if (data === null || typeof data !== "object") {
            console.warn("setByPath: data not an object or array.");
            return;
        }

        let result = data;
        let parts = path.split(delimiter);
        let lastPart = parts.pop();

        for (const part of parts) {
            if (
                (!utils.isObject(result) && !Array.isArray(result))
                || (!utils.isObject(result[part]) && !Array.isArray(result[part]))
            ) {
                result[part] = {};
            }

            result = result[part];
        }

        result[lastPart] = newValue;
    },

    /**
     * Recursively delete element from an object (or array) by its key path.
     * Empty array or object elements from the parents chain will be also removed.
     *
     * @example
     * ```javascript
     * this.deleteByPath({a: {b: {c: 3}}}, "a.b.c");       // results in {}
     * this.deleteByPath({a: {b: {c: 3, d: 4}}}, "a.b.c"); // results in {a: {b: {d: 4}}}
     * ```
     *
     * @param  {Array|Object} data
     * @param  {string}       path
     * @param  {string}       delimiter
     */
    deleteByPath(data, path, delimiter = ".") {
        let result = data || {};
        let parts = (path || "").split(delimiter);
        let lastPart = parts.pop();

        for (const part of parts) {
            if (
                (!utils.isObject(result) && !Array.isArray(result))
                || (!utils.isObject(result[part]) && !Array.isArray(result[part]))
            ) {
                result[part] = {};
            }

            result = result[part];
        }

        if (Array.isArray(result)) {
            result.splice(lastPart, 1);
        } else if (utils.isObject(result)) {
            delete result[lastPart];
        }

        // cleanup the parents chain
        if (
            parts.length > 0
            && ((Array.isArray(result) && !result.length) || (utils.isObject(result) && !Object.keys(result).length))
            && ((Array.isArray(data) && data.length > 0) || (utils.isObject(data) && Object.keys(data).length > 0))
        ) {
            utils.deleteByPath(data, parts.join(delimiter), delimiter);
        }
    },

    /**
     * Returns a new object with the zero value of the existing defined props.
     *
     * @param  {Object} obj
     * @param  {Array}  [preservedProps]
     * @return {Object}
     */
    emptyClone(obj, preservedProps = []) {
        const clone = JSON.parse(JSON.stringify(obj));

        for (let prop in clone) {
            if (preservedProps.includes(prop)) {
                continue;
            }

            if (typeof clone[prop] == "string") {
                clone[prop] = "";
            } else if (typeof clone[prop] == "number") {
                clone[prop] = 0;
            } else if (typeof clone[prop] == "boolean") {
                clone[prop] = false;
            } else if (Array.isArray(clone[prop])) {
                clone[prop] = [];
            } else if (app.utils.isObject(clone[prop])) {
                clone[prop] = {};
            }
        }

        return clone;
    },

    /**
     * Generates pseudo-random string (suitable for ids and keys).
     *
     * @param  {number} [length] The string of the resulting random string (default 8)
     * @return {string}
     */
    randomString(length = 8, alphabet = DEFAULT_RANDOM_ALPHABET) {
        let result = "";

        for (let i = 0; i < length; i++) {
            result += alphabet.charAt(Math.floor(Math.random() * alphabet.length));
        }

        return result;
    },

    /**
     * Attempts to generates cryptographically random secret when `crypto`
     * is supported, otherwise fallback to `app.utils.randomString`.
     *
     * @param  {number} [length] The string of the resulting random string (default 30)
     * @return {string}
     */
    randomSecret(length = 30) {
        if (typeof crypto === "undefined") {
            return app.utils.randomString(length);
        }

        const arr = new Uint8Array(length);
        crypto.getRandomValues(arr);

        const alphabet = "-_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"; // 64 to devide "cleanly" 256

        let result = "";

        for (let i = 0; i < length; i++) {
            result += alphabet.charAt(arr[i] % alphabet.length);
        }

        return result;
    },

    /**
     * Converts and normalizes string into a sentence.
     *
     * @param  {string}  str
     * @param  {boolean} [stopCheck]
     * @return {string}
     */
    sentenize(str, stopCheck = true) {
        if (typeof str !== "string") {
            return "";
        }

        str = str.trim().split("_").join(" ");
        if (str === "") {
            return str;
        }

        str = str[0].toUpperCase() + str.substring(1);

        if (stopCheck) {
            let lastChar = str[str.length - 1];
            if (lastChar !== "." && lastChar !== "?" && lastChar !== "!") {
                str += ".";
            }
        }

        return str;
    },

    /**
     * Trims the matching quotes from the provided value.
     *
     * The value will be returned unchanged if `val` is not
     * wrapped with quotes or it is not string.
     *
     * @param  {Mixed} val
     * @return {Mixed}
     */
    trimQuotedValue(val) {
        if (
            typeof val == "string"
            && (val[0] == `"` || val[0] == `'` || val[0] == "`")
            && val[0] == val[val.length - 1]
        ) {
            return val.slice(1, -1);
        }

        return val;
    },

    /**
     * Normalizes and converts the provided string to a slug.
     *
     * @param  {string} str
     * @param  {string} [delimiter]
     * @param  {Array}  [preserved] List of special characters to keep unmodified.
     * @return {string}
     */
    slugify(str, delimiter = "_", preserved = []) {
        if (str === "") {
            return "";
        }

        // special characters
        const specialCharsMap = {
            a: /а|à|á|å|â/gi,
            b: /б/gi,
            c: /ц|ç/gi,
            d: /д/gi,
            e: /е|è|é|ê|ẽ|ë/gi,
            f: /ф/gi,
            g: /г/gi,
            h: /х/gi,
            i: /й|и|ì|í|î/gi,
            j: /ж/gi,
            k: /к/gi,
            l: /л/gi,
            m: /м/gi,
            n: /н|ñ/gi,
            o: /о|ò|ó|ô|ø/gi,
            p: /п/gi,
            q: /я/gi,
            r: /р/gi,
            s: /с/gi,
            t: /т/gi,
            u: /ю|ù|ú|ů|û/gi,
            v: /в/gi,
            w: /в/gi,
            x: /ь/gi,
            y: /ъ/gi,
            z: /з/gi,
            ae: /ä|æ/gi,
            oe: /ö/gi,
            ue: /ü/gi,
            Ae: /Ä/gi,
            Ue: /Ü/gi,
            Oe: /Ö/gi,
            ss: /ß/gi,
            and: /&/gi,
        };

        // replace special characters
        for (let k in specialCharsMap) {
            str = str.replace(specialCharsMap[k], k);
        }

        return str
            .replace(new RegExp("[" + preserved.join("") + "]", "g"), " ") // replace preserved characters with spaces
            .replace(/[^\w\ ]/gi, "") // replaces all non-alphanumeric with empty string
            .replace(/\s+/g, delimiter); // collapse whitespaces and replace with `delimiter`
    },

    /**
     * Encodes the HTML entities of the specified string.
     *
     * @param  {string} str
     * @return {string}
     */
    encodeEntities(str) {
        if (!str) {
            return "";
        }

        return str
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll("\"", "&quot;")
            .replaceAll("'", "&#039;");
    },

    /**
     * Returns the plain text version (aka. strip tags) of the provided string.
     *
     * NB! HTML entities are preserved. If you want to remove them call
     * `app.utils.encodeEntities(result)` on the plainText result.
     *
     * @param  {string}  str
     * @return {string}
     */
    plainText(str) {
        if (!str) {
            return "";
        }

        const doc = new DOMParser().parseFromString(str, "text/html");

        return (doc.body.textContent || "").trim();
    },

    /**
     * Truncates the provided text to the specified max characters length.
     *
     * @param  {string}  str
     * @param  {Number}  [length]
     * @param  {boolean} [dots]
     * @return {string}
     */
    truncate(str, length = 150, dots = true) {
        str = "" + str;

        if (str.length <= length) {
            return str;
        }

        str = str.slice(0, length);

        if (dots) {
            while (str.endsWith(".")) {
                str = str.slice(0, -1);
            }
            str += "...";
        }

        return str;
    },

    /**
     * Returns a new object/array copy with truncated large text fields.
     *
     * @param  {Object|Array} objOrArr
     * @param  {Number}       [length]
     * @param  {boolean}      [dots]
     * @return {Object|Array}
     */
    truncateObject(objOrArr, length = 150, dots = true) {
        const truncated = Array.isArray(objOrArr) ? [] : {};

        for (let key in objOrArr) {
            let value = objOrArr[key];

            if (typeof value === "string") {
                value = app.utils.truncate(value, length, dots);
            } else if (Array.isArray(value)) {
                value = app.utils.truncateObject(value, length, dots);
            } else if (app.utils.isObject(value)) {
                value = app.utils.truncateObject(value, length, dots);
            }

            truncated[key] = value;
        }

        return truncated;
    },

    /**
     * Returns a stringified truncated version of the provided value
     * or fallback to `missingValue` in case it is empty.
     *
     * @param  {Mixed}  val
     * @param  {number} [truncateLength]
     * @param  {string} [missingValue]
     * @return {string}
     */
    displayValue(val, truncateLength = 150, missingValue = "N/A") {
        // check the raw value for "emptiness"
        if (utils.isEmpty(val)) {
            return missingValue;
        }

        if (typeof val == "string") {
            // already a string
        } else if (typeof val == "boolean") {
            val = val ? "True" : "False";
        } else if (Array.isArray(val) && typeof val[0] != "object") {
            // assuming primitive array values
            val = val.map((child) => utils.displayValue(child, truncateLength, missingValue)).join(", ");
        } else {
            try {
                val = JSON.stringify(val) || "";
            } catch (_) {
                val = "" + val;
            }
        }

        return val ? utils.truncate(val, truncateLength) : missingValue;
    },

    /**
     * Splits `str` and returns its non empty parts as an array.
     *
     * @param  {string} str
     * @param  {string} [separator]
     * @return {Array}
     */
    splitNonEmpty(str, separator = ",") {
        const items = (str || "").split(separator);
        const result = [];

        for (let item of items) {
            item = item.trim();
            if (!utils.isEmpty(item)) {
                result.push(item);
            }
        }

        return result;
    },

    /**
     * Returns a concatenated `items` string of only the none empty values.
     *
     * @param  {Array} items
     * @param  {string} [separator]
     * @return {Array}
     */
    joinNonEmpty(items, separator = ", ") {
        items = items || [];

        const result = [];

        for (let item of items) {
            item = typeof item === "string" ? item.trim() : item;
            if (!utils.isEmpty(item)) {
                result.push("" + item);
            }
        }

        return result.join(separator);
    },

    /**
     * Returns a human readable file size string from size in bytes.
     *
     * @param  {Number} size s
     * @return {string}
     */
    formattedFileSize(size) {
        const i = size ? Math.floor(Math.log(size) / Math.log(1024)) : 0;

        return (size / Math.pow(1024, i)).toFixed(2) * 1 + " " + ["B", "KB", "MB", "GB", "TB"][i];
    },

    /**
     * Returns a RFC3339 datetime formatted string (YYYY-MM-DD HH:mm:ss.nnnZ)
     * from the specified `datetime-local` input value (e.g. YYYY-MM-DDTHH:mm:ss).
     *
     * @param  {string|number|Date} strOrDate
     * @return {string}
     */
    toRFC3339Datetime(strOrDate) {
        if (!strOrDate) {
            return "";
        }

        let date;
        if (strOrDate instanceof Date) {
            date = strOrDate;
        } else if (typeof strOrDate == "string") {
            date = new Date(strOrDate.replace(" ", "T"));
        } else {
            date = new Date(strOrDate);
        }

        return date.toISOString().replace("T", " ");
    },

    toLocalDatetime(strOrDate) {
        if (!strOrDate) {
            return "";
        }

        let date;
        if (strOrDate instanceof Date) {
            date = strOrDate;
        } else if (typeof strOrDate == "string") {
            date = new Date(strOrDate.replace(" ", "T"));
        } else {
            date = new Date(strOrDate);
        }

        const year = date.getFullYear();
        if (isNaN(year)) {
            return ""; // invalid date
        }

        const month = (date.getMonth() + 1).toString().padStart(2, "0");
        const day = date.getDate().toString().padStart(2, "0");
        const hours = date.getHours().toString().padStart(2, "0");
        const minutes = date.getMinutes().toString().padStart(2, "0");
        const seconds = date.getSeconds().toString().padStart(2, "0");
        const milliseconds = date.getMilliseconds().toString().padStart(3, "0");

        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}.${milliseconds}`;
    },

    /**
     * Returns `datetime-local` input value string (YYYY-MM-DDTHH:mm:ss)
     * from the specified date (e.g. YYYY-MM-DD HH:mm:ss.nnnZ).
     *
     * @param  {string|number|Date} strOrDate
     * @return {string}
     */
    toDatetimeLocalInputValue(strOrDate) {
        if (!strOrDate) {
            return "";
        }

        let date;
        if (strOrDate instanceof Date) {
            date = strOrDate;
        } else if (typeof strOrDate == "string") {
            date = new Date(strOrDate.replaceAll(" ", "T"));
        } else {
            date = new Date(strOrDate);
        }

        const year = date.getFullYear();
        if (isNaN(year)) {
            return ""; // invalid date
        }

        const month = (date.getMonth() + 1).toString().padStart(2, "0");
        const day = date.getDate().toString().padStart(2, "0");
        const hours = date.getHours().toString().padStart(2, "0");
        const minutes = date.getMinutes().toString().padStart(2, "0");
        const seconds = date.getSeconds().toString().padStart(2, "0");

        return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}`;
    },

    /**
     * Copies the string representation of `val` to the user clipboard.
     *
     * @param  {Mixed} val
     * @return {Promise}
     */
    async copyToClipboard(val) {
        // normalizes val
        if (val === null || typeof val == "undefined") {
            val = "";
        } else if (val instanceof Date) {
            val = val.toISOString();
        } else if (typeof val == "object") {
            val = JSON.stringify(val);
        } else {
            val = "" + val;
        }

        if (!val.length || !window.navigator?.clipboard) {
            return;
        }

        return window.navigator.clipboard.writeText(val).catch((err) => {
            console.warn("Failed to copy.", err);
        });
    },

    /**
     * Forces the browser to start downloading the specified url.
     *
     * @param {string} url  The url of the file to download.
     * @param {string} name The result file name.
     */
    download(url, name) {
        let tempLink = document.createElement("a");
        tempLink.setAttribute("href", url);
        tempLink.setAttribute("download", name);
        tempLink.setAttribute("target", "_blank");
        tempLink.setAttribute("rel", "noopener noreferrer");
        tempLink.click();
        tempLink = null;
    },

    /**
     * Downloads a json file created from the provide object.
     *
     * @param {mixed} obj   The JS object to download.
     * @param {string} name The result file name.
     */
    downloadJSON(obj, name) {
        name = name.endsWith(".json") ? name : name + ".json";

        const blob = new Blob([JSON.stringify(obj, null, 2)], {
            type: "application/json",
        });

        const url = window.URL.createObjectURL(blob);

        utils.download(url, name);
    },

    /**
     * Returns a normalized API URL address that is used in the API example docs.
     *
     * @return {string}
     */
    getApiExampleURL() {
        let url;

        // use the JS SDK url if it is already an absolute url
        if (
            app.pb.baseURL.startsWith("http://")
            || app.pb.baseURL.startsWith("https://")
        ) {
            url = app.pb.baseURL;
        } else {
            url = window.location.href;

            // check for the ui path in case the app is under a subpath
            const start = url.indexOf("/_/");
            if (start >= 0) {
                url = url.substring(0, start);
            } else {
                url = window.location.origin;
            }
        }

        // for broader compatibility replace localhost with 127.0.0.1
        // (see https://github.com/pocketbase/js-sdk/issues/21)
        return url.replace("//localhost", "//127.0.0.1");
    },

    /**
     * @todo consider making part of Shablon or at least suggest as a snippet?
     *
     * Returns if the specified href address match the current hash route path.
     *
     * @param  {string}  href
     * @param  {boolean} [subPathPattern] Whether to use exact or sub path matching pattern (default to true).
     * @param  {string}  [customHash] Hash string to parse (if not set, fallbacks to `navigationStore.hash`).
     * @return {boolean}
     */
    isActivePath(href, subPathPattern = true, customHash = "") {
        customHash = customHash || navigationStore.hash;

        let pattern;
        if (subPathPattern) {
            pattern = new RegExp("^" + RegExp.escape(href) + "\\/?.*$");
        } else {
            pattern = new RegExp("^" + RegExp.escape(href) + "\\/?(?:\\?.+)?$");
        }

        return pattern.test(customHash);
    },

    /**
     * @todo move to Shablon?
     * @todo consider making the value arrays for consistency with the router
     *
     * Extracts the hash query parameters from the current url and
     * returns them as plain object.
     *
     * @param  {string} [customHash] Hash string to parse (if not set, fallbacks to `window.location.hash`).
     * @return {Object}
     */
    getHashQueryParams(customHash = "") {
        customHash = customHash || navigationStore.hash;

        let query = "";

        const queryStart = customHash.indexOf("?");
        if (queryStart > -1) {
            query = window.location.hash.substring(queryStart + 1);
        }

        return Object.fromEntries(new URLSearchParams(query));
    },

    /**
     * @todo move to Shablon?
     *
     * Replaces the current hash query parameters with the provided `params`.
     *
     * Empty `params` values are removed from the final qyery string.
     *
     * @param {Object} params The new parameters to replace.
     * @param {null|boolean} [updateHistory] Specifies what to do with window.history:
     *                                       - true: push a new history state
     *                                       - false: does nothing to window.history
     *                                       - null/undefined: replaces the current history state (default)
     * @return {string} A new absolute url with the replaced hash query params.
     */
    replaceHashQueryParams(params, updateHistory = null) {
        params = params || {};

        let query = "";

        let hash = window.location.hash;

        const queryStart = hash.indexOf("?");
        if (queryStart > -1) {
            query = hash.substring(queryStart + 1);
            hash = hash.substring(0, queryStart);
        }

        const parsed = new URLSearchParams(query);

        for (const key in params) {
            const val = params[key];
            if (utils.isEmpty(val)) {
                parsed.delete(key);
            } else {
                parsed.set(key, val);
            }
        }

        query = parsed.toString();
        if (query != "") {
            hash += "?" + query;
        }

        // replace the hash/fragment part with the updated one
        const original = window.location.href;
        let base = original;
        const hashIndex = base.indexOf("#");
        if (hashIndex > -1) {
            base = base.substring(0, hashIndex);
        }

        const newHref = base + hash;

        if (updateHistory === false) {
            // no-op...
        } else if (updateHistory === true) {
            window.history.pushState(null, "", newHref);
        } else {
            window.history.replaceState(null, "", newHref);
        }

        return newHref;
    },

    /**
     * Locally stores the current path for later redirect.
     */
    rememberPath() {
        window.localStorage.setItem(REMEMBER_PATH_KEY, window.location.hash);
    },

    /**
     * Redirect to a remembered local path.
     *
     * @param {string} [fallback] Fallback path if there is nothing stored.
     */
    toRememberedPath(fallback = "#/collections") {
        let path = window.localStorage.getItem(REMEMBER_PATH_KEY);
        if (path) {
            window.localStorage.removeItem(REMEMBER_PATH_KEY);
        }

        window.location.hash = path || fallback;
    },

    /**
     * Returns and deserializes a localStorage stored JSON value.
     *
     * @param  {string} key
     * @param  {Mixed}  [defaultVal]
     * @return {Mixed}  The deserialized found value (or `defaultVal` if missing).
     */
    getLocalHistory(key, defaultVal = null) {
        try {
            const raw = window.localStorage.getItem(key);
            if (raw) {
                return JSON.parse(raw) || defaultVal;
            }
        } catch (err) {
            console.log("failed to load local history:", key, err);
        }

        return defaultVal;
    },

    /**
     * Serializes and saves in localStorage the provided data.
     *
     * If `data` is "empty" value the localStorage entry will be removed.
     * If data is string it is saved as it is, otherwise with `JSON.stringify()`.
     *
     * @param {string} key
     * @param {Mixed}  data
     */
    saveLocalHistory(key, data) {
        try {
            if (app.utils.isEmpty(data)) {
                window.localStorage.removeItem(key);
            } else if (typeof data == "string") {
                window.localStorage.setItem(key, data);
            } else {
                window.localStorage.setItem(key, JSON.stringify(data));
            }
        } catch (err) {
            console.log("failed to save local history:", key, err);
        }
    },

    /**
     * Creates a thumbnail from `File` with the specified `width` and `height` params.
     * Returns a `Promise` with the generated base64 url.
     *
     * @param  {File}   file
     * @param  {Number} [width]
     * @param  {Number} [height]
     * @return {Promise}
     */
    generateThumb(file, width = 100, height = 100) {
        return new Promise((resolve) => {
            let reader = new FileReader();

            reader.onload = function(e) {
                let img = new Image();

                img.onload = function() {
                    let canvas = document.createElement("canvas");
                    let ctx = canvas.getContext("2d");
                    let imgWidth = img.width;
                    let imgHeight = img.height;

                    canvas.width = width;
                    canvas.height = height;

                    ctx.drawImage(
                        img,
                        imgWidth > imgHeight ? (imgWidth - imgHeight) / 2 : 0,
                        0, // top aligned
                        imgWidth > imgHeight ? imgHeight : imgWidth,
                        imgWidth > imgHeight ? imgHeight : imgWidth,
                        0,
                        0,
                        width,
                        height,
                    );

                    return resolve(canvas.toDataURL(file.type));
                };

                img.src = e.target.result;
            };

            reader.readAsDataURL(file);
        });
    },

    /**
     * Normalizes the search filter by converting a simple search term into
     * a wildcard filter expression using the provided fallback search fields.
     *
     * If searchTerm is already an expression it is returned without changes.
     *
     * @param  {string} searchTerm
     * @param  {Array}  fallbackFields
     * @return {string}
     */
    normalizeSearchFilter(searchTerm, fallbackFields = []) {
        searchTerm = (searchTerm || "").trim();
        if (!searchTerm || !fallbackFields.length) {
            return searchTerm;
        }

        const opChars = ["=", "!=", "~", "!~", ">", ">=", "<", "<="];

        // loosely check if it is already a filter expression
        for (const op of opChars) {
            if (searchTerm.includes(op)) {
                return searchTerm;
            }
        }

        searchTerm = isNaN(searchTerm) && searchTerm != "true" && searchTerm != "false"
            ? `"${searchTerm.replace(/^[\"\'\`]|[\"\'\`]$/gm, "")}"`
            : searchTerm;

        return fallbackFields.map((f) => `${f}~${searchTerm}`).join("||");
    },

    logLevels: {
        [-4]: {
            label: "DEBUG",
            class: "",
        },
        0: {
            label: "INFO",
            class: "success",
        },
        4: {
            label: "WARN",
            class: "warning",
        },
        8: {
            label: "ERROR",
            class: "danger",
        },
    },
    logDataFormatters: {
        execTime: function(log) {
            if (typeof log?.data?.execTime == "undefined") {
                return "N/A";
            }
            return log.data.execTime + "ms";
        },
    },

    /**
     * @todo consider defining as helper in shablon?
     *
     * Extends a reactive baseStore with the provided attrs by applying
     * similar rules as the attrs for DOM elements.
     *
     * Returns an array with the resulting watchers that you can use to unsubscribe
     * once you are done with the store.
     *
     * @param  {Proxy} baseStore
     * @param  {Object} attrs
     * @param  {Array} [exclude]
     * @return {Array}
     */
    extendStore(baseStore, attrs = {}, ...exclude) {
        const watchers = [];

        for (let key in attrs) {
            let val = attrs[key];

            if (
                typeof baseStore.__raw?.[key] == "function"
                || typeof val != "function"
                || (key.length > 2 && key.startsWith("on"))
                // @todo consider using exclude to skip prop loading?
                || exclude.includes(key)
            ) {
                baseStore[key] = val;
            } else {
                watchers.push(
                    watch(val, (result) => {
                        baseStore[key] = result;
                    }),
                );
            }
        }

        return watchers;
    },

    /**
     * Converts a CSS time string into a millisecond number.
     * Returns 0 if empty or invalid.
     *
     * @param  {string} cssTimeStr
     * @return {number}
     */
    cssTimeToMs(cssTimeStr) {
        if (!cssTimeStr) {
            return 0;
        }

        cssTimeStr = cssTimeStr.toLowerCase();

        if (cssTimeStr.endsWith("ms")) {
            return Number(cssTimeStr.substring(0, cssTimeStr.length - 2));
        }

        if (cssTimeStr.endsWith("s")) {
            return Number(cssTimeStr.substring(0, cssTimeStr.length - 1));
        }

        return Number(cssTimeStr) || 0;
    },

    /**
     * Very rudimentary check to evaluate if the provided HEX color is "dark",
     * aka. whether it is suitable as background for a white text.
     *
     * @see https://en.wikipedia.org/wiki/YIQ
     * @see https://24ways.org/2010/calculating-color-contrast/
     *
     * @param  {string} hexcolor The HEX color in its 6 chars expanded format (with or without the "#" prefix).
     * @return {boolean}
     */
    isDarkEnoughForWhiteText(hexcolor) {
        hexcolor = hexcolor?.startsWith("#") ? hexcolor.substring(1) : hexcolor;

        if (hexcolor?.length != 6) {
            return false;
        }

        const r = parseInt(hexcolor.substring(0, 2), 16);
        const g = parseInt(hexcolor.substring(2, 4), 16);
        const b = parseInt(hexcolor.substring(4, 6), 16);
        const yiq = ((r * 299) + (g * 587) + (b * 114)) / 1000;

        return yiq < 128;
    },

    // ---------------------------------------------------------------
    imageExtensions: [".jpg", ".jpeg", ".png", ".svg", ".gif", ".jfif", ".webp", ".avif"],
    videoExtensions: [".mp4", ".avi", ".mov", ".3gp", ".wmv"],
    audioExtensions: [".aa", ".aac", ".m4v", ".mp3", ".ogg", ".oga", ".mogg", ".amr"],
    documentExtensions: [".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odp", ".odt", ".ods", ".txt"],

    /**
     * Loosely check if a file has image extension.
     *
     * @param  {string} filename
     * @return {boolean}
     */
    hasImageExtension(filename) {
        filename = (filename || "").toLowerCase();
        return !!app.utils.imageExtensions.find((ext) => filename.endsWith(ext));
    },

    /**
     * Loosely check if a file has video extension.
     *
     * @param  {string} filename
     * @return {boolean}
     */
    hasVideoExtension(filename) {
        filename = (filename || "").toLowerCase();
        return !!app.utils.videoExtensions.find((ext) => filename.endsWith(ext));
    },

    /**
     * Loosely check if a file has audio extension.
     *
     * @param  {string} filename
     * @return {boolean}
     */
    hasAudioExtension(filename) {
        filename = (filename || "").toLowerCase();
        return !!app.utils.audioExtensions.find((ext) => filename.endsWith(ext));
    },

    /**
     * Loosely check if a file has document extension.
     *
     * @param  {string} filename
     * @return {boolean}
     */
    hasDocumentExtension(filename) {
        filename = (filename || "").toLowerCase();
        return !!app.utils.documentExtensions.find((ext) => filename.endsWith(ext));
    },

    /**
     * Returns the file type based on its filename.
     *
     * @param  {string} filename
     * @return {string}
     */
    getFileType(filename) {
        if (app.utils.hasImageExtension(filename)) return "image";
        if (app.utils.hasVideoExtension(filename)) return "video";
        if (app.utils.hasAudioExtension(filename)) return "audio";
        if (app.utils.hasDocumentExtension(filename)) return "document";
        return "file";
    },
    fileTypeIcons: {
        image: "ri-image-line",
        video: "ri-movie-line",
        audio: "ri-music-2-line",
        document: "ri-file-line",
        file: "ri-file-line",
    },

    // ---------------------------------------------------------------

    fallbackFieldIcon: "ri-puzzle-line",
    fallbackCollectionIcon: "ri-puzzle-line",
    fallbackProviderIcon: "ri-puzzle-line",
    fallbackPresentableProps: [
        "title",
        "name",
        "slug",
        "email",
        "username",
        "nickname",
        "displayName",
        "label",
        "subject",
        "topic",
        "message",
        "heading",
        "headline",
        "header",
        "caption",
        "key",
        "identifier",
        "id",
    ],

    /**
     * Returns a shallow copy of the specified collections ordered by their name.
     *
     * @param  {Array} collections
     * @return {Array}
     */
    sortedCollections(collections = []) {
        let underscoreA, underscoreB;
        function sortNames(a, b) {
            // order system collections last
            underscoreA = a.name.startsWith("_");
            underscoreB = b.name.startsWith("_");
            if (underscoreA && !underscoreB) {
                return 1;
            }
            if (!underscoreA && underscoreB) {
                return -1;
            }

            if (a.name > b.name) {
                return 1;
            }
            if (a.name < b.name) {
                return -1;
            }

            return 0;
        }

        return collections.slice().sort(sortNames);
    },

    /**
     * Groups and sorts collections array by type (auth, base, view) and name.
     *
     * @param  {Array} collections
     * @return {Array}
     */
    sortedCollectionsByType(collections = []) {
        const auth = [];
        const base = [];
        const view = [];

        for (const collection of collections) {
            if (collection.type === "auth") {
                auth.push(collection);
            } else if (collection.type === "base") {
                base.push(collection);
            } else {
                view.push(collection);
            }
        }

        return [].concat(
            app.utils.sortedCollections(auth),
            app.utils.sortedCollections(base),
            app.utils.sortedCollections(view),
        );
    },

    /**
     * Checks if the provided 2 collections has any change (ignoring root fields order).
     *
     * @param  {Object} oldCollection
     * @param  {Object} newCollection
     * @param  {boolean}    [withDeleteMissing] Skip missing fields from the newCollection.
     * @return {boolean}
     */
    hasCollectionChanges(oldCollection, newCollection, withDeleteMissing = false) {
        oldCollection = oldCollection || {};
        newCollection = newCollection || {};

        if (oldCollection.id != newCollection.id) {
            return true;
        }

        for (let prop in oldCollection) {
            if (prop !== "fields" && JSON.stringify(oldCollection[prop]) !== JSON.stringify(newCollection[prop])) {
                return true;
            }
        }

        const oldFields = Array.isArray(oldCollection.fields) ? oldCollection.fields : [];
        const newFields = Array.isArray(newCollection.fields) ? newCollection.fields : [];

        const removedFields = oldFields.filter((oldField) => {
            return oldField?.id && !newFields.find((f) => f.id == oldField.id);
        });

        const addedFields = newFields.filter((newField) => {
            return newField?.id && !oldFields.find((f) => f.id == newField.id);
        });

        const changedFields = newFields.filter((newField) => {
            const oldField = app.utils.isObject(newField) && oldFields.find((f) => f.id == newField.id);
            if (!oldField) {
                return false;
            }

            for (let prop in oldField) {
                if (JSON.stringify(newField[prop]) != JSON.stringify(oldField[prop])) {
                    return true;
                }
            }

            return false;
        });

        return !!(addedFields.length || changedFields.length || (withDeleteMissing && removedFields.length));
    },

    /**
     * Rudimentary SELECT query columns extractor.
     * Returns an array with the identifier aliases
     * (expressions wrapped in parenthesis are skipped).
     *
     * @param  {string} selectQuery
     * @return {Array}
     */
    extractColumnsFromQuery(selectQuery) {
        const groupReplacement = "__PBGROUP__";

        selectQuery = (selectQuery || "")
            // replace parenthesis/group expessions
            .replace(/\([\s\S]+?\)/gm, groupReplacement)
            // replace multi-whitespace characters with single space
            .replace(/[\t\r\n]|(?:\s\s)+/g, " ");

        const match = selectQuery.match(/select\s+([\s\S]+)\s+from/);

        const expressions = match?.[1]?.split(",") || [];

        const result = [];

        for (let expr of expressions) {
            const column = expr.trim().split(" ").pop(); // get only the alias
            if (column != "" && column != groupReplacement) {
                result.push(column.replace(/[\'\"\`\[\]\s]/g, ""));
            }
        }

        return result;
    },

    /**
     * Returns an array with all public collection identifiers (collection fields + type specific fields).
     *
     * @param  {Object} collection The collection to extract identifiers from.
     * @param  {string} [prefix]   Optional prefix for each found identified.
     * @return {Array}
     */
    getAllCollectionIdentifiers(collection, prefix = "") {
        if (!collection) {
            return [];
        }

        let result = [prefix + "id"];

        const isAuth = collection.type == "auth";
        const isView = collection.type === "view";

        if (isView) {
            for (let col of app.utils.extractColumnsFromQuery(collection.viewQuery)) {
                app.utils.pushUnique(result, prefix + col);
            }
        }

        const fields = collection.fields || [];
        for (const field of fields) {
            if (field.type == "password" || (isAuth && field.name == "tokenKey")) {
                continue;
            }

            if (app.fieldTypes[field.type]?.identifierExtractor) {
                const vals = app.utils.toArray(app.fieldTypes[field.type]?.identifierExtractor(field, prefix));
                for (let val of vals) {
                    app.utils.pushUnique(result, val);
                }
            } else {
                app.utils.pushUnique(result, prefix + field.name);
            }
        }

        return result;
    },

    /**
     * Returns a plain record object populated with dummy data (used usually in the API previews).
     *
     * @param  {Object}  collection
     * @param  {Boolean} [forSubmit]
     * @return {Object}
     */
    getDummyFieldsData(collection, forSubmit = false) {
        const fields = collection?.fields || [];

        const result = {};

        for (const field of fields) {
            if (field.hidden) {
                continue;
            }

            if (app.fieldTypes[field.type]?.dummyData) {
                const val = app.fieldTypes[field.type].dummyData(field, forSubmit);
                if (typeof val !== "undefined") {
                    result[field.name] = val;
                }
            } else {
                result[field.name] = "[[DATA]]";
            }
        }

        return result;
    },

    // SQL indexes
    // ---------------------------------------------------------------

    /**
     * Parses the specified SQL index and returns an object with its components.
     *
     * For example:
     *
     * ```js
     * parseIndex("CREATE UNIQUE INDEX IF NOT EXISTS schemaname.idxname on tablename (col1, col2) where expr")
     * // output:
     * {
     *   "unique":     true,
     *   "optional":   true,
     *   "schemaName": "schemaname"
     *   "indexName":  "idxname"
     *   "tableName":  "tablename"
     *   "columns":    [{name: "col1", "collate": "", "sort": ""}, {name: "col1", "collate": "", "sort": ""}]
     *   "where":      "expr"
     * }
     * ```
     *
     * @param  {string} idx
     * @return {Object}
     */
    parseIndex(idx) {
        const result = {
            unique: false,
            optional: false,
            schemaName: "",
            indexName: "",
            tableName: "",
            columns: [],
            where: "",
        };

        const indexRegex =
            /create\s+(unique\s+)?\s*index\s*(if\s+not\s+exists\s+)?(\S*)\s+on\s+(\S*)\s*\(([\s\S]*)\)(?:\s*where\s+([\s\S]*))?/gim;
        const matches = indexRegex.exec((idx || "").trim());

        if (matches?.length != 7) {
            return result;
        }

        const sqlQuoteRegex = /^[\"\'\`\[\{}]|[\"\'\`\]\}]$/gm;

        // unique
        result.unique = matches[1]?.trim().toLowerCase() === "unique";

        // optional
        result.optional = !app.utils.isEmpty(matches[2]?.trim());

        // schemaName and indexName
        const namePair = (matches[3] || "").split(".");
        if (namePair.length == 2) {
            result.schemaName = namePair[0].replace(sqlQuoteRegex, "");
            result.indexName = namePair[1].replace(sqlQuoteRegex, "");
        } else {
            result.schemaName = "";
            result.indexName = namePair[0].replace(sqlQuoteRegex, "");
        }

        // tableName
        result.tableName = (matches[4] || "").replace(sqlQuoteRegex, "");

        // columns
        const rawColumns = (matches[5] || "")
            .replace(/,(?=[^\(]*\))/gim, "{PB_TEMP}") // temporary replace comma within expressions for easier splitting
            .split(","); // split columns

        for (let col of rawColumns) {
            col = col.trim().replaceAll("{PB_TEMP}", ","); // revert temp replacement

            const colRegex = /^([\s\S]+?)(?:\s+collate\s+([\w]+))?(?:\s+(asc|desc))?$/gim;
            const colMatches = colRegex.exec(col);
            if (colMatches?.length != 4) {
                continue;
            }

            const colOrExpr = colMatches[1]?.trim()?.replace(sqlQuoteRegex, "");
            if (!colOrExpr) {
                continue;
            }
            result.columns.push({
                name: colOrExpr,
                collate: colMatches[2] || "",
                sort: colMatches[3]?.toUpperCase() || "",
            });
        }

        // WHERE expression
        result.where = matches[6] || "";

        return result;
    },

    /**
     * Builds an index expression from parsed index parts (see parseIndex()).
     *
     * @param  {Array} indexParts
     * @return {string}
     */
    buildIndex(indexParts) {
        let result = "CREATE ";

        if (indexParts.unique) {
            result += "UNIQUE ";
        }

        result += "INDEX ";

        if (indexParts.optional) {
            result += "IF NOT EXISTS ";
        }

        if (indexParts.schemaName) {
            result += `\`${indexParts.schemaName}\`.`;
        }

        result += `\`${indexParts.indexName || "idx_" + app.utils.randomString(10)}\` `;

        result += `ON \`${indexParts.tableName}\` (`;

        const nonEmptyCols = indexParts.columns.filter((col) => !!col?.name);

        if (nonEmptyCols.length > 1) {
            result += "\n  ";
        }

        result += nonEmptyCols
            .map((col) => {
                let item = "";

                if (col.name.includes("(") || col.name.includes(" ")) {
                    // most likely an expression
                    item += col.name;
                } else {
                    // regular identifier
                    item += "`" + col.name + "`";
                }

                if (col.collate) {
                    item += " COLLATE " + col.collate;
                }

                if (col.sort) {
                    item += " " + col.sort.toUpperCase();
                }

                return item;
            })
            .join(",\n  ");

        if (nonEmptyCols.length > 1) {
            result += "\n";
        }

        result += `)`;

        if (indexParts.where) {
            result += ` WHERE ${indexParts.where}`;
        }

        return result;
    },

    /**
     * Parses and merges the current index with the specified `newFields`.
     *
     * The `newFields` argument could be:
     * - plain object with the same props as `parseIndex`
     * - function that accepted the current parsed index and returns a new object with the fields to overwrite
     *
     * @param  {string} idx
     * @param  {Object|Function} newFields
     * @return {string}
     */
    replaceIndexFields(idx, newFields) {
        let parsed = app.utils.parseIndex(idx);

        if (typeof newFields == "function") {
            Object.assign(parsed, newFields(parsed) || {});
        } else {
            Object.assign(parsed, newFields || {});
        }

        return app.utils.buildIndex(parsed);
    },

    /**
     * Replaces an idx column name with a new one (if exists).
     *
     * @param  {string} idx
     * @param  {string} oldColumn
     * @param  {string} newColumn
     * @return {string}
     */
    replaceIndexColumn(idx, oldColumn, newColumn) {
        if (oldColumn === newColumn) {
            return idx; // no change
        }

        const parsed = app.utils.parseIndex(idx);

        let hasChange = false;
        for (let col of parsed.columns) {
            if (col.name === oldColumn) {
                col.name = newColumn;
                hasChange = true;
            }
        }

        return hasChange ? app.utils.buildIndex(parsed) : idx;
    },
};

window.app = window.app || {};
window.app.utils = utils;
