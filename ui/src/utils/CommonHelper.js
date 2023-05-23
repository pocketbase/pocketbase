import { DateTime } from "luxon";

const imageExtensions = [
   ".jpg", ".jpeg", ".png", ".svg",
   ".gif", ".jfif", ".webp", ".avif",
];

const videoExtensions = [
    ".mp4", ".avi", ".mov", ".3gp", ".wmv",
];

const audioExtensions = [
    ".aa", ".aac", ".m4v", ".mp3",
    ".ogg", ".oga", ".mogg", ".amr",
];

const documentExtensions = [
    ".pdf", ".doc", ".docx", ".xls",
    ".xlsx", ".ppt", ".pptx", ".odp",
    ".odt", ".ods", ".txt",
];

export default class CommonHelper {
    /**
     * Checks whether value is plain object.
     *
     * @param  {Mixed} value
     * @return {Boolean}
     */
    static isObject(value) {
        return value !== null && typeof value === "object" && value.constructor === Object;
    }

    /**
     * Deep clones the provided value.
     *
     * @param  {Mixed} val
     * @return {Mixed}
     */
    static clone(value) {
        return typeof structuredClone !== "undefined" ? structuredClone(value) : JSON.parse(JSON.stringify(value));
    }

    /**
     * Checks whether a value is empty. The following values are considered as empty:
     * - null
     * - undefined
     * - empty string
     * - empty array
     * - empty object
     * - zero uuid, time and dates
     *
     * @param  {Mixed} value
     * @return {Boolean}
     */
    static isEmpty(value) {
        return (
            (value === "") ||
            (value === null) ||
            (value === "00000000-0000-0000-0000-000000000000") || // zero uuid
            (value === "0001-01-01 00:00:00.000Z") || // zero datetime
            (value === "0001-01-01") || // zero date
            (typeof value === "undefined") ||
            (Array.isArray(value) && value.length === 0) ||
            (CommonHelper.isObject(value) && Object.keys(value).length === 0)
        );
    }

    /**
     * Checks whether the provided dom element is a form field (input, textarea, select).
     *
     * @param  {Node} element
     * @return {Boolean}
     */
    static isInput(element) {
        let tagName = element && element.tagName ? element.tagName.toLowerCase() : "";

        return (
            tagName === "input" ||
            tagName === "select" ||
            tagName === "textarea" ||
            element.isContentEditable
        )
    }

    /**
     * Checks if an element is a common focusable one.
     *
     * @param  {Node} element
     * @return {Boolean}
     */
    static isFocusable(element) {
        let tagName = element && element.tagName ? element.tagName.toLowerCase() : "";

        return (
            CommonHelper.isInput(element) ||
            tagName === "button" ||
            tagName === "a" ||
            tagName === "details" ||
            element.tabIndex >= 0
        );
    }

    /**
     * Check if `obj` has at least one none empty property.
     *
     * @param  {Object} obj
     * @return {Boolean}
     */
    static hasNonEmptyProps(obj) {
        for (let i in obj) {
            if (!CommonHelper.isEmpty(obj[i])) {
                return true;
            }
        }

        return false;
    }

    /**
     * Normalizes and returns arr as a new array instance.
     *
     * @param  {Array}   arr
     * @param  {Boolean} [allowEmpty]
     * @return {Array}
     */
    static toArray(arr, allowEmpty = false) {
        if (Array.isArray(arr)) {
            return arr.slice();
        }

        return (allowEmpty || !CommonHelper.isEmpty(arr)) && typeof arr !== "undefined" ? [arr] : [];
    }

    /**
     * Loosely checks if value exists in an array.
     *
     * @param  {Array}  arr
     * @param  {String} value
     * @return {Boolean}
     */
    static inArray(arr, value) {
        arr = Array.isArray(arr) ? arr : [];

        for (let i = arr.length - 1; i >= 0; i--) {
            if (arr[i] == value) {
                return true;
            }
        }

        return false;
    }

    /**
     * Removes single element from array by loosely comparying values.
     *
     * @param {Array} arr
     * @param {Mixed} value
     */
    static removeByValue(arr, value) {
        arr = Array.isArray(arr) ? arr : [];

        for (let i = arr.length - 1; i >= 0; i--) {
            if (arr[i] == value) {
                arr.splice(i, 1);
                break;
            }
        }
    }

    /**
     * Adds `value` in `arr` only if it's not added already.
     *
     * @param {Array} arr
     * @param {Mixed} value
     */
    static pushUnique(arr, value) {
        if (!CommonHelper.inArray(arr, value)) {
            arr.push(value);
        }
    }

    /**
     * Returns single element from objects array by matching its key value.
     *
     * @param  {Array} objectsArr
     * @param  {Mixed} key
     * @param  {Mixed} value
     * @return {Object}
     */
    static findByKey(objectsArr, key, value) {
        objectsArr = Array.isArray(objectsArr) ? objectsArr : [];

        for (let i in objectsArr) {
            if (objectsArr[i][key] == value) {
                return objectsArr[i];
            }
        }

        return null;
    }

    /**
     * Group objects array by a specific key.
     *
     * @param  {Array}  objectsArr
     * @param  {String} key
     * @return {Object}
     */
    static groupByKey(objectsArr, key) {
        objectsArr = Array.isArray(objectsArr) ? objectsArr : [];

        const result = {};

        for (let i in objectsArr) {
            result[objectsArr[i][key]] = result[objectsArr[i][key]] || [];

            result[objectsArr[i][key]].push(objectsArr[i]);
        }

        return result;
    }

    /**
     * Removes single element from objects array by matching an item"s property value.
     *
     * @param {Array}  objectsArr
     * @param {String} key
     * @param {Mixed}  value
     */
    static removeByKey(objectsArr, key, value) {
        for (let i in objectsArr) {
            if (objectsArr[i][key] == value) {
                objectsArr.splice(i, 1);
                break;
            }
        }
    }

    /**
     * Adds or replace an object array element by comparing its key value.
     *
     * @param {Array}  objectsArr
     * @param {Object} item
     * @param {Mixed}  [key]
     */
    static pushOrReplaceByKey(objectsArr, item, key = "id") {
        for (let i = objectsArr.length - 1; i >= 0; i--) {
            if (objectsArr[i][key] == item[key]) {
                objectsArr[i] = item; // replace
                return;
            }
        }

        objectsArr.push(item);
    }

    /**
     * Filters and returns a new objects array with duplicated elements removed.
     *
     * @param  {Array} objectsArr
     * @param  {String} key
     * @return {Array}
     */
    static filterDuplicatesByKey(objectsArr, key = "id") {
        objectsArr = Array.isArray(objectsArr) ? objectsArr : [];

        const uniqueMap = {};

        for (const item of objectsArr) {
            uniqueMap[item[key]] = item;
        }

        return Object.values(uniqueMap)
    }

    /**
     * Filters and returns a new object with removed redacted props.
     *
     * @param  {Object} obj
     * @param  {String} [mask] Default to '******'
     * @return {Object}
     */
    static filterRedactedProps(obj, mask = "******") {
        const result = JSON.parse(JSON.stringify(obj || {}));

        for (let prop in result) {
            if (typeof result[prop] === "object" && result[prop] !== null) {
                result[prop] = CommonHelper.filterRedactedProps(result[prop], mask)
            } else if (result[prop] === mask) {
                delete result[prop];
            }
        }
        return result;
    }

    /**
     * Safely access nested object/array key with dot-notation.
     *
     * @example
     * ```javascript
     * var myObj = {a: {b: {c: 3}}}
     * this.getNestedVal(myObj, "a.b.c");       // returns 3
     * this.getNestedVal(myObj, "a.b.c.d");     // returns null
     * this.getNestedVal(myObj, "a.b.c.d", -1); // returns -1
     * ```
     *
     * @param  {Object|Array} data
     * @param  {string}       path
     * @param  {Mixed}        [defaultVal]
     * @param  {String}       [delimiter]
     * @return {Mixed}
     */
    static getNestedVal(data, path, defaultVal = null, delimiter = ".") {
        let result = data || {};
        let parts  = (path || "").split(delimiter);

        for (const part of parts) {
            if (
                (!CommonHelper.isObject(result) && !Array.isArray(result)) ||
                typeof result[part] === "undefined"
            ) {
                return defaultVal;
            }

            result = result[part];
        }

        return result;
    }

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
     * @param  {String}       delimiter
     */
    static setByPath(data, path, newValue, delimiter = ".") {
        if (data === null || typeof data !== "object") {
            console.warn("setByPath: data not an object or array.");
            return
        }

        let result   = data;
        let parts    = path.split(delimiter);
        let lastPart = parts.pop();

        for (const part of parts) {
            if (
                (!CommonHelper.isObject(result) && !Array.isArray(result)) ||
                (!CommonHelper.isObject(result[part]) && !Array.isArray(result[part]))
            ) {
                result[part] = {};
            }

            result = result[part];
        }

        result[lastPart] = newValue;
    }

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
     * @param  {String}       delimiter
     */
    static deleteByPath(data, path, delimiter = ".") {
        let result   = data || {};
        let parts    = (path || "").split(delimiter);
        let lastPart = parts.pop();

        for (const part of parts) {
            if (
                (!CommonHelper.isObject(result) && !Array.isArray(result)) ||
                (!CommonHelper.isObject(result[part]) && !Array.isArray(result[part]))
            ) {
                result[part] = {};
            }

            result = result[part];
        }

        if (Array.isArray(result)) {
            result.splice(lastPart, 1);
        } else if (CommonHelper.isObject(result)) {
            delete (result[lastPart]);
        }

        // cleanup the parents chain
        if (
            parts.length > 0 &&
            (
                (Array.isArray(result) && !result.length) ||
                (CommonHelper.isObject(result) && !Object.keys(result).length)
            ) &&
            (
                (Array.isArray(data) && data.length > 0) ||
                (CommonHelper.isObject(data) && Object.keys(data).length > 0)
            )
        ) {
            CommonHelper.deleteByPath(data, parts.join(delimiter), delimiter);
        }
    }

    /**
     * Generates random string (suitable for elements id and keys).
     *
     * @param  {Number} [length] Results string length (default 10)
     * @return {String}
     */
    static randomString(length) {
        length = length || 10;

        let result = "";
        let alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

        for (let i = 0; i < length; i++) {
            result += alphabet.charAt(Math.floor(Math.random() * alphabet.length));
        }

        return result;
    }

    /**
     * Converts and normalizes string into a sentence.
     *
     * @param  {String}  str
     * @param  {Boolean} [stopCheck]
     * @return {String}
     */
    static sentenize(str, stopCheck = true) {
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

        return str
    }

    /**
     * Trims the matching quotes from the provided value.
     *
     * The value will be returned unchanged if `val` is not
     * wrapped with quotes or it is not string.
     *
     * @param  {Mixed} val
     * @return {Mixed}
     */
    static trimQuotedValue(val) {
        if (
            typeof val == "string" &&
            (val[0] == `"`  || val[0] == `'` || val[0] == "`") &&
            val[0] == val[val.length-1]
        ) {
            return val.slice(1, -1);
        }

        return val
    }

    /**
     * Returns the plain text version (aka. strip tags) of the provided string.
     *
     * @param  {String} str
     * @return {String}
     */
    static plainText(str) {
        if (!str) {
            return "";
        }

        const doc = new DOMParser().parseFromString(str, "text/html");

        return (doc.body.innerText || "").trim();
    }

    /**
     * Truncates the provided text to the specified max characters length.
     *
     * @param  {String}  str
     * @param  {Number}  [length]
     * @param  {Boolean} [dots]
     * @return {String}
     */
    static truncate(str, length = 150, dots = true) {
        str = str || "";

        if (str.length <= length) {
            return str;
        }

        return str.substring(0, length) + (dots ? "..." : "");
    }

    /**
     * Returns a new object copy with truncated the large text fields.
     *
     * @param  {Object} obj
     * @return {Object}
     */
    static truncateObject(obj) {
        const truncated = {};

        for (let key in obj) {
            let value = obj[key];

            if (typeof value === "string") {
                value = CommonHelper.truncate(value, 150, true);
            }

            truncated[key] = value;
        }

        return truncated;
    }

    /**
     * Normalizes and converts the provided string to a slug.
     *
     * @param  {String} str
     * @param  {String} [delimiter]
     * @param  {Array}  [preserved]
     * @return {String}
     */
    static slugify(str, delimiter = "_", preserved = [".", "=", "-"]) {
        if (str === "") {
            return "";
        }

        // special characters
        const specialCharsMap = {
            "a": /а|à|á|å|â/gi,
            "b": /б/gi,
            "c": /ц|ç/gi,
            "d": /д/gi,
            "e": /е|è|é|ê|ẽ|ë/gi,
            "f": /ф/gi,
            "g": /г/gi,
            "h": /х/gi,
            "i": /й|и|ì|í|î/gi,
            "j": /ж/gi,
            "k": /к/gi,
            "l": /л/gi,
            "m": /м/gi,
            "n": /н|ñ/gi,
            "o": /о|ò|ó|ô|ø/gi,
            "p": /п/gi,
            "q": /я/gi,
            "r": /р/gi,
            "s": /с/gi,
            "t": /т/gi,
            "u": /ю|ù|ú|ů|û/gi,
            "v": /в/gi,
            "w": /в/gi,
            "x": /ь/gi,
            "y": /ъ/gi,
            "z": /з/gi,
            "ae": /ä|æ/gi,
            "oe": /ö/gi,
            "ue": /ü/gi,
            "Ae": /Ä/gi,
            "Ue": /Ü/gi,
            "Oe": /Ö/gi,
            "ss": /ß/gi,
            "and": /&/gi
        };

        // replace special characters
        for (let k in specialCharsMap) {
            str = str.replace(specialCharsMap[k], k);
        }

        return str
            .replace(new RegExp('[' + preserved.join("") + ']', 'g'), ' ') // replace preserved characters with spaces
            .replace(/[^\w\ ]/gi, "")                                      // replaces all non-alphanumeric with empty string
            .replace(/\s+/g, delimiter);                                   // collapse whitespaces and replace with `delimiter`
    }

    /**
     * Returns `str` with escaped regexp characters, making it safe to
     * embed in regexp as a whole literal expression.
     *
     * @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions#escaping
     * @param  {String} str
     * @return {String}
     */
    static escapeRegExp(str) {
      return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'); // $& means the whole matched string
    }

    /**
     * Splits `str` and returns its non empty parts as an array.
     *
     * @param  {String} str
     * @param  {String} [separator]
     * @return {Array}
     */
    static splitNonEmpty(str, separator = ",") {
        const items = (str || "").split(separator);
        const result = [];

        for (let item of items) {
            item = item.trim();
            if (!CommonHelper.isEmpty(item)) {
                result.push(item);
            }
        }

        return result;
    }

    /**
     * Returns a concatenated `items` string.
     *
     * @param  {String} items
     * @param  {String} [separator]
     * @return {Array}
     */
    static joinNonEmpty(items, separator = ", ") {
        const result = [];

        for (let item of items) {
            item = typeof item === "string" ? item.trim() : "";
            if (!CommonHelper.isEmpty(item)) {
                result.push(item);
            }
        }

        return result.join(separator);
    }

    /**
     * Extract the user initials from the provided username or email address
     * (eg. converts "john.doe@example.com" to "JD").
     *
     * @param  {String} str
     * @return {String}
     */
    static getInitials(str) {
        str = (str || "").split("@")[0].trim();

        if (str.length <= 2) {
            return str.toUpperCase();
        }

        const parts = str.split(/[\.\_\-\ ]/);

        if (parts.length >= 2) {
            return (parts[0][0] + parts[1][0]).toUpperCase();
        }

        return str[0].toUpperCase();
    }

    /**
     * Returns a human readable file size string from size in bytes.
     *
     * @param  {Number} size s
     * @return {String}
     */
    static formattedFileSize(size) {
        const i = size ? Math.floor(Math.log(size) / Math.log(1024)) : 0;

        return (size / Math.pow(1024, i)).toFixed(2) * 1 + " " + ["B", "KB", "MB", "GB", "TB"][i];
    }

    /**
     * Returns a DateTime instance from a date object/string.
     *
     * @param  {String|Date} date
     * @return {DateTime}
     */
    static getDateTime(date) {
        if (typeof date === "string") {
            const formats = {
                19: "yyyy-MM-dd HH:mm:ss",
                23: "yyyy-MM-dd HH:mm:ss.SSS",
                20: "yyyy-MM-dd HH:mm:ss'Z'",
                24: "yyyy-MM-dd HH:mm:ss.SSS'Z'",
            }
            const format = formats[date.length] || formats[19];
            return DateTime.fromFormat(date, format, { zone: "UTC" });
        }

        return DateTime.fromJSDate(date);
    }

    /**
     * Returns formatted datetime string in the UTC timezone.
     *
     * @param  {String|Date} date
     * @param  {String}      [format] The result format (see https://moment.github.io/luxon/#/parsing?id=table-of-tokens)
     * @return {String}
     */
    static formatToUTCDate(date, format = "yyyy-MM-dd HH:mm:ss") {
        return CommonHelper.getDateTime(date).toUTC().toFormat(format);
    }

    /**
     * Returns formatted datetime string in the local timezone.
     *
     * @param  {String|Date} date
     * @param  {String}      [format] The result format (see https://moment.github.io/luxon/#/parsing?id=table-of-tokens)
     * @return {String}
     */
    static formatToLocalDate(date, format = "yyyy-MM-dd HH:mm:ss") {
        return CommonHelper.getDateTime(date).toLocal().toFormat(format);
    }

    /**
     * Copies text to the user clipboard.
     *
     * @param  {String} text
     * @return {Promise}
     */
    static async copyToClipboard(text) {
        text = "" + text // ensure that text is string

        if (!text.length || !window?.navigator?.clipboard) {
            return;
        }

        return window.navigator.clipboard.writeText(text).catch((err) => {
            console.warn("Failed to copy.", err);
        })
    }

    /**
     * Forces the browser to start downloading the specified url.
     *
     * @param {String} url  The url of the file to download.
     * @param {String} name The result file name.
     */
    static download(url, name) {
        const tempLink = document.createElement("a");
        tempLink.setAttribute("href", url);
        tempLink.setAttribute("download", name);
        tempLink.click();
        tempLink.remove();
    }

    /**
     * Downloads a json file created from the provide object.
     *
     * @param {mixed} obj   The JS object to download.
     * @param {String} name The result file name.
     */
    static downloadJson(obj, name) {
        const encodedObj = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(obj, null, 2));

        name = name.endsWith(".json") ? name : (name + ".json");

        CommonHelper.download(encodedObj, name)
    }

    /**
     * Parses and returns the decoded jwt payload data.
     *
     * @param  {String} jwt
     * @return {Object}
     */
    static getJWTPayload(jwt) {
        const raw = (jwt || "").split(".")[1] || "";
        if (raw === "") {
            return {};
        }

        try {
            const encodedPayload = decodeURIComponent(atob(raw));
            return JSON.parse(encodedPayload) || {};
        } catch (err) {
            console.warn("Failed to parse JWT payload data.", err);
        }

        return  {};
    }

    /**
     * Loosely check if a file has image extension.
     *
     * @param  {String} filename
     * @return {Boolean}
     */
    static hasImageExtension(filename) {
        return !!imageExtensions.find((ext) => filename.toLowerCase().endsWith(ext));
    }

    /**
     * Loosely check if a file has video extension.
     *
     * @param  {String} filename
     * @return {Boolean}
     */
    static hasVideoExtension(filename) {
        return !!videoExtensions.find((ext) => filename.toLowerCase().endsWith(ext));
    }

    /**
     * Loosely check if a file has audio extension.
     *
     * @param  {String} filename
     * @return {Boolean}
     */
    static hasAudioExtension(filename) {
        return !!audioExtensions.find((ext) => filename.toLowerCase().endsWith(ext));
    }

    /**
     * Loosely check if a file has document extension.
     *
     * @param  {String} filename
     * @return {Boolean}
     */
    static hasDocumentExtension(filename) {
        return !!documentExtensions.find((ext) => filename.toLowerCase().endsWith(ext));
    }

    /**
     * Returns the file type based on its filename.
     *
     * @param  {String} filename
     * @return {String}
     */
    static getFileType(filename) {
        if (CommonHelper.hasImageExtension(filename)) return "image";
        if (CommonHelper.hasDocumentExtension(filename)) return "document";
        if (CommonHelper.hasVideoExtension(filename)) return "video";
        if (CommonHelper.hasAudioExtension(filename)) return "audio";
        return "file";
    }

    /**
     * Creates a thumbnail from `File` with the specified `width` and `height` params.
     * Returns a `Promise` with the generated base64 url.
     *
     * @param  {File}   file
     * @param  {Number} [width]
     * @param  {Number} [height]
     * @return {Promise}
     */
    static generateThumb(file, width = 100, height = 100) {
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
                        // imgHeight > imgWidth ? (imgHeight - imgWidth) / 2 : 0,
                        imgWidth > imgHeight ? imgHeight : imgWidth,
                        imgWidth > imgHeight ? imgHeight : imgWidth,
                        0,
                        0,
                        width,
                        height
                    );

                    return resolve(canvas.toDataURL(file.type));
                };

                img.src = e.target.result;
            }

            reader.readAsDataURL(file);
        });
    }

    /**
     * Normalizes and append a value to the provided form data.
     *
     * @param {FormData} formData
     * @param {string}   key
     * @param {mixed}    value
     */
    static addValueToFormData(formData, key, value) {
        if (typeof value === "undefined") {
            return;
        }

        if (CommonHelper.isEmpty(value)) {
            formData.append(key, "");
        } else if (Array.isArray(value)) {
            for (const item of value) {
                CommonHelper.addValueToFormData(formData, key, item);
            }
        } else if (value instanceof File) {
            formData.append(key, value);
        } else if (value instanceof Date) {
            formData.append(key, value.toISOString());
        } else if (CommonHelper.isObject(value)) {
            formData.append(key, JSON.stringify(value));
        } else {
            formData.append(key, "" + value);
        }
    }

    /**
     * Returns a dummy collection record object.
     *
     * @param  {Object} collection
     * @return {Object}
     */
    static dummyCollectionRecord(collection) {
        const fields = collection?.schema || [];

        const dummy = {
            "id": "RECORD_ID",
            "collectionId": collection?.id,
            "collectionName": collection?.name,
        };

        if (collection?.isAuth) {
            dummy["username"] = "username123";
            dummy["verified"] = false;
            dummy["emailVisibility"] = true;
            dummy["email"] = "test@example.com";
        }

        const hasCreated = !collection?.$isView || CommonHelper.extractColumnsFromQuery(collection?.options?.query).includes("created");
        if (hasCreated) {
            dummy["created"] = "2022-01-01 01:00:00.123Z";
        }

        const hasUpdated = !collection?.$isView || CommonHelper.extractColumnsFromQuery(collection?.options?.query).includes("updated");
        if (hasUpdated) {
            dummy["updated"] = "2022-01-01 23:59:59.456Z";
        }

        for (const field of fields) {
            let val = null;
            if (field.type === "number") {
                val = 123;
            } else if (field.type === "date") {
                val = "2022-01-01 10:00:00.123Z";
            } else if (field.type === "bool") {
                val = true;
            } else if (field.type === "email") {
                val = "test@example.com";
            } else if (field.type === "url") {
                val = "https://example.com";
            } else if (field.type === "json") {
                val = 'JSON';
            } else if (field.type === "file") {
                val = 'filename.jpg';
                if (field.options?.maxSelect !== 1) {
                    val = [val];
                }
            } else if (field.type === "select") {
                val = field.options?.values?.[0];
                if (field.options?.maxSelect !== 1) {
                    val = [val];
                }
            } else if (field.type === "relation") {
                val = 'RELATION_RECORD_ID';
                if (field.options?.maxSelect !== 1) {
                    val = [val];
                }
            } else {
                val = "test";
            }

            dummy[field.name] = val;
        }

        return dummy;
    }

    /**
     * Returns a dummy collection schema data object.
     *
     * @param  {Object} collection
     * @return {Object}
     */
    static dummyCollectionSchemaData(collection) {
        const fields = collection?.schema || [];

        const dummy = {};

        for (const field of fields) {
            let val = null;

            if (field.type === "number") {
                val = 123;
            } else if (field.type === "date") {
                val = "2022-01-01 10:00:00.123Z";
            } else if (field.type === "bool") {
                val = true;
            } else if (field.type === "email") {
                val = "test@example.com";
            } else if (field.type === "url") {
                val = "https://example.com";
            } else if (field.type === "json") {
                val = 'JSON';
            } else if (field.type === "file") {
                continue; // currently file upload is supported only via FormData
            } else if (field.type === "select") {
                val = field.options?.values?.[0];
                if (field.options?.maxSelect !== 1) {
                    val = [val];
                }
            } else if (field.type === "relation") {
                val = 'RELATION_RECORD_ID';
                if (field.options?.maxSelect !== 1) {
                    val = [val];
                }
            } else {
                val = "test";
            }

            dummy[field.name] = val;
        }

        return dummy;
    }

    /**
     * Returns a collection type icon.
     *
     * @param  {String} type
     * @return {String}
     */
    static getCollectionTypeIcon(type) {
        switch (type?.toLowerCase()) {
            case "auth":
                return "ri-group-line";
            case "view":
                return "ri-table-line";
            default:
                return "ri-folder-2-line";
        }
    }

    /**
     * Returns a field type icon.
     *
     * @param  {String} type
     * @return {String}
     */
    static getFieldTypeIcon(type) {
        switch (type?.toLowerCase()) {
            case "primary":
                return "ri-key-line";
            case "text":
                return "ri-text";
            case "number":
                return "ri-hashtag";
            case "date":
                return "ri-calendar-line";
            case "bool":
                return "ri-toggle-line";
            case "email":
                return "ri-mail-line";
            case "url":
                return "ri-link";
            case "editor":
                return "ri-edit-2-line";
            case "select":
                return "ri-list-check";
            case "json":
                return "ri-braces-line";
            case "file":
                return "ri-image-line";
            case "relation":
                return "ri-mind-map";
            case "user":
                return "ri-user-line";
            default:
                return "ri-star-s-line";
        }
    }

    /**
     * Returns the field value base type as text.
     *
     * @param  {Object} field
     * @return {String}
     */
    static getFieldValueType(field) {
        switch (field?.type) {
            case 'bool':
                return 'Boolean';
            case 'number':
                return 'Number';
            case 'file':
                return 'File';
            case 'select':
            case 'relation':
                if (field?.options?.maxSelect === 1) {
                    return 'String';
                }
                return 'Array<String>';
            default:
                return 'String';
        }
    }

    /**
     * Returns the zero-default string value of the provided field.
     *
     * @param  {Object} field
     * @return {String}
     */
    static zeroDefaultStr(field) {
        if (field?.type === "number") {
            return "0";
        }

        if (field?.type === "bool") {
            return "false";
        }

        if (field?.type === "json") {
            return 'null, "", [], {}';
        }

        // arrayable fields
        if (["select", "relation", "file"].includes(field?.type) && field?.options?.maxSelect != 1) {
            return "[]";
        }

        return '""';
    }

    /**
     * Returns an API url address extract from the current running instance.
     *
     * @param  {String} fallback Fallback url that will be used if the extractions fail.
     * @return {String}
     */
    static getApiExampleUrl(fallback) {
        let url = window.location.href.substring(0, window.location.href.indexOf("/_")) || fallback || '/';

        // for broader compatibility replace localhost with 127.0.0.1
        // (see https://github.com/pocketbase/js-sdk/issues/21)
        return url.replace('//localhost', '//127.0.0.1');
    }

    /**
     * Checks if the provided 2 collections has any change (ignoring root schema fields order).
     *
     * @param  {Collection} oldCollection
     * @param  {Collection} newCollection
     * @param  {Boolean}    withDeleteMissing Skip missing schema fields from the newCollection.
     * @return {Boolean}
     */
    static hasCollectionChanges(oldCollection, newCollection, withDeleteMissing = false) {
        oldCollection = oldCollection || {};
        newCollection = newCollection || {};

        if (oldCollection.id != newCollection.id) {
            return true;
        }

        for (let prop in oldCollection) {
            if (prop !== 'schema' && JSON.stringify(oldCollection[prop]) !== JSON.stringify(newCollection[prop])) {
                return true;
            }
        }

        const oldSchema = Array.isArray(oldCollection.schema) ? oldCollection.schema : [];
        const newSchema = Array.isArray(newCollection.schema) ? newCollection.schema : [];
        const removedFields = oldSchema.filter((oldField) => {
            return oldField?.id && !CommonHelper.findByKey(newSchema, "id", oldField.id);
        });
        const addedFields = newSchema.filter((newField) => {
            return newField?.id && !CommonHelper.findByKey(oldSchema, "id", newField.id);
        });
        const changedFields = newSchema.filter((newField) => {
            const oldField = CommonHelper.isObject(newField) && CommonHelper.findByKey(oldSchema, "id", newField.id);
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

        return !!(
            addedFields.length ||
            changedFields.length ||
            (withDeleteMissing && removedFields.length)
        );
    }

    /**
     * Groups and sorts collections array by type (auth, base, view).
     *
     * @param  {Array} collections
     * @return {Array}
     */
    static sortCollections(collections = []) {
        const auth = [];
        const base = [];
        const view = [];

        for (const collection of collections) {
            if (collection.type === 'auth') {
                auth.push(collection);
            } else if (collection.type === 'base') {
                base.push(collection);
            } else {
                view.push(collection);
            }
        }

        return [].concat(auth, base, view);
    }


    /**
     * "Yield" to the main thread to break long runing task into smaller ones.
     *
     * (see https://web.dev/optimize-long-tasks/)
     */
    static yieldToMain() {
        return new Promise((resolve) => {
            setTimeout(resolve, 0);
        });
    }

    /**
     * Returns the default Flatpickr initialization options.
     *
     * @return {Object}
     */
    static defaultFlatpickrOptions() {
        return {
            dateFormat: "Y-m-d H:i:S",
            disableMobile: true,
            allowInput: true,
            enableTime: true,
            time_24hr: true,
            locale: {
                firstDayOfWeek: 1,
            },
        }
    }

    /**
     * Returns the default rich editor options.
     *
     * @return {Object}
     */
    static defaultEditorOptions() {
        return {
            branding: false,
            promotion: false,
            menubar: false,
            min_height: 270,
            height: 270,
            max_height: 700,
            autoresize_bottom_margin: 30,
            skin: "pocketbase",
            content_style: "body { font-size: 14px }",
            plugins: [
                "autoresize",
                "autolink",
                "lists",
                "link",
                "image",
                "searchreplace",
                "fullscreen",
                "media",
                "table",
                "code",
                "codesample",
                "directionality",
            ],
            toolbar: "styles | alignleft aligncenter alignright | bold italic forecolor backcolor | bullist numlist | link image table codesample direction | code fullscreen",
            file_picker_types: "image",
            // @see https://www.tiny.cloud/docs/tinymce/6/file-image-upload/#interactive-example
            file_picker_callback: (cb, value, meta) => {
                const input = document.createElement("input");
                input.setAttribute("type", "file");
                input.setAttribute("accept", "image/*");

                input.addEventListener("change", (e) => {
                    const file = e.target.files[0];
                    const reader = new FileReader();

                    reader.addEventListener("load", () => {
                        if (!tinymce) {
                            return;
                        }

                        // We need to register the blob in TinyMCEs image blob registry.
                        // In future TinyMCE version this part will be handled internally.
                        const id = "blobid" + new Date().getTime();
                        const blobCache = tinymce.activeEditor.editorUpload.blobCache;
                        const base64 = reader.result.split(",")[1];
                        const blobInfo = blobCache.create(id, file, base64);
                        blobCache.add(blobInfo);

                        // call the callback and populate the Title field with the file name
                        cb(blobInfo.blobUri(), { title: file.name });
                    });

                    reader.readAsDataURL(file);
                });

                input.click();
            },
            setup: (editor) => {
                editor.on('keydown', (e) => {
                    // propagate save shortcut to the parent
                    if ((e.ctrlKey || e.metaKey) && e.code == "KeyS" && editor.formElement) {
                        e.preventDefault();
                        e.stopPropagation();
                        editor.formElement.dispatchEvent(new KeyboardEvent("keydown", e));
                    }
                });

                const lastDirectionKey = "tinymce_last_direction";

                // load last used text direction for blank editors
                editor.on('init', () => {
                    const lastDirection = window?.localStorage?.getItem(lastDirectionKey);
                    if (!editor.isDirty() && editor.getContent() == "" && lastDirection == "rtl") {
                        editor.execCommand("mceDirectionRTL");
                    }
                });

                // text direction dropdown
                editor.ui.registry.addMenuButton("direction", {
                    icon: "visualchars",
                    fetch: (callback) => {
                        const items = [
                            {
                                type: "menuitem",
                                text: "LTR content",
                                icon: "ltr",
                                onAction: () => {
                                    window?.localStorage?.setItem(lastDirectionKey, "ltr");
                                    tinymce.activeEditor.execCommand("mceDirectionLTR");
                                }
                            },
                            {
                                type: "menuitem",
                                text: "RTL content",
                                icon: "rtl",
                                onAction: () => {
                                    window?.localStorage?.setItem(lastDirectionKey, "rtl");
                                    tinymce.activeEditor.execCommand("mceDirectionRTL");
                                }
                            }
                        ];

                        callback(items);
                    }
                });
            },
        };
    }

    /**
     * Tries to output the first displayable field of the provided model.
     *
     * @param  {Object} model
     * @return {Any}
     */
    static displayValue(model, displayFields, missingValue = "N/A") {
        model = model || {};
        displayFields = displayFields || [];

        let result = [];

        for (const field of displayFields) {
            let val = model[field];

            if (typeof val === "undefined") {
                continue
            }

            if (CommonHelper.isEmpty(val)) {
                result.push(missingValue);
            } else if (typeof val === "boolean")  {
                result.push(val ? "True" : "False");
            } else if (typeof val === "string") {
                val = val.indexOf("<") >= 0 ? CommonHelper.plainText(val) : val;
                result.push(CommonHelper.truncate(val));
            } else {
                result.push(val);
            }
        }

        if (result.length > 0) {
            return result.join(", ");
        }

        const fallbackProps = [
            "title",
            "name",
            "slug",
            "email",
            "username",
            "label",
            "heading",
            "message",
            "key",
            "id",
        ];

        for (const prop of fallbackProps) {
            if (!CommonHelper.isEmpty(model[prop])) {
                return model[prop];
            }
        }

        return missingValue;
    }

    /**
     * Rudimentary SELECT query columns extractor.
     * Returns an array with the identifier aliases
     * (expressions wrapped in parenthesis are skipped).
     *
     * @param  {String} selectQuery
     * @return {Array}
     */
    static extractColumnsFromQuery(selectQuery) {
        const groupReplacement = "__GROUP__";

        selectQuery = (selectQuery || "").
            // replace parenthesis/group expessions
            replace(/\([\s\S]+?\)/gm, groupReplacement).
            // replace multi-whitespace characters with single space
            replace(/[\t\r\n]|(?:\s\s)+/g, " ");

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
    }

    /**
     * Returns an array with all public collection identifiers (schema + type specific fields).
     *
     * @param  {[type]} collection The collection to extract identifiers from.
     * @param  {String} prefix     Optional prefix for each found identified.
     * @return {Array}
     */
    static getAllCollectionIdentifiers(collection, prefix = "") {
        if (!collection) {
            return [];
        }

        let result = [prefix + "id"];

        if (collection.$isView) {
            for (let col of CommonHelper.extractColumnsFromQuery(collection.options.query)) {
                CommonHelper.pushUnique(result, prefix + col);
            }
        } else if (collection.$isAuth) {
            result.push(prefix + "username");
            result.push(prefix + "email");
            result.push(prefix + "emailVisibility");
            result.push(prefix + "verified");
            result.push(prefix + "created");
            result.push(prefix + "updated");
        } else {
            result.push(prefix + "created");
            result.push(prefix + "updated");
        }

        const schema = collection.schema || [];

        for (const field of schema) {
            CommonHelper.pushUnique(result, prefix + field.name);
        }

        return result;
    }

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
     * @param  {String} idx
     * @return {Object}
     */
    static parseIndex(idx) {
        const result = {
            unique:     false,
            optional:   false,
            schemaName: "",
            indexName:  "",
            tableName:  "",
            columns:    [],
            where:      "",
        };

        const indexRegex = /create\s+(unique\s+)?\s*index\s*(if\s+not\s+exists\s+)?(\S*)\s+on\s+(\S*)\s+\(([\s\S]*)\)(?:\s*where\s+([\s\S]*))?/gmi;
        const matches    = indexRegex.exec((idx || "").trim())

        if (matches?.length != 7) {
            return result;
        }

        const sqlQuoteRegex = /^[\"\'\`\[\{}]|[\"\'\`\]\}]$/gm

        // unique
        result.unique = matches[1]?.trim().toLowerCase() === "unique";

        // optional
        result.optional = !CommonHelper.isEmpty(matches[2]?.trim());

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
            .replace(/,(?=[^\(]*\))/gmi, "{PB_TEMP}") // temporary replace comma within expressions for easier splitting
            .split(",");                              // split columns

        for (let col of rawColumns) {
            col = col.trim().replaceAll("{PB_TEMP}", ",") // revert temp replacement

            const colRegex = /^([\s\S]+?)(?:\s+collate\s+([\w]+))?(?:\s+(asc|desc))?$/gmi
            const colMatches = colRegex.exec(col);
            if (colMatches?.length != 4) {
                continue
            }

            const colOrExpr = colMatches[1]?.trim()?.replace(sqlQuoteRegex, "");
            if (!colOrExpr) {
                continue;
            }
            result.columns.push({
                name:    colOrExpr,
                collate: colMatches[2] || "",
                sort:    colMatches[3]?.toUpperCase() || "",
            });
        }

        // WHERE expression
        result.where = matches[6] || "";

        return result;
    }

    /**
     * Builds an index expression from parsed index parts (see parseIndex()).
     *
     * @param  {Array} indexParts
     * @return {String}
     */
    static buildIndex(indexParts) {
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

        result += `\`${indexParts.indexName || "idx_" + CommonHelper.randomString(7)}\` `;

        result += `ON \`${indexParts.tableName}\` (`;

        const nonEmptyCols = indexParts.columns.filter((col) => !!col?.name);

        if (nonEmptyCols.length > 1) {
            result += "\n  ";
        }

        result += nonEmptyCols.map((col) => {
                let item = "";

                if (col.name.includes("(") || col.name.includes(" ")) {
                    // most likely an expression
                    item += col.name;
                } else {
                    // regular identifier
                    item += ("`" + col.name + "`");
                }

                if (col.collate) {
                    item += (" COLLATE " + col.collate);
                }

                if (col.sort) {
                    item += (" " + c.sort.toUpperCase());
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
    }

    /**
     * Replaces the idx table name with newTableName.
     *
     * @param  {String} idx
     * @param  {String} newTableName
     * @return {String}
     */
    static replaceIndexTableName(idx, newTableName) {
        const parsed = CommonHelper.parseIndex(idx);

        parsed.tableName = newTableName;

        return CommonHelper.buildIndex(parsed);
    }

    /**
     * Replaces an idx column name with a new one (if exists).
     *
     * @param  {String} idx
     * @param  {String} oldColumn
     * @param  {String} newColumn
     * @return {String}
     */
    static replaceIndexColumn(idx, oldColumn, newColumn) {
        if (oldColumn === newColumn) {
            return idx; // no change
        }

        const parsed = CommonHelper.parseIndex(idx);

        let hasChange = false;
        for (let col of parsed.columns) {
            if (col.name === oldColumn) {
                col.name = newColumn;
                hasChange = true;
            }
        }

        return hasChange ? CommonHelper.buildIndex(parsed) : idx;
    }

    /**
     * Normalizes the search filter by converting a simple search term into
     * a wildcard filter expression using the provided fallback search fields.
     *
     * If searchTerm is already an expression it is returned without changes.
     *
     * @param  {String} searchTerm
     * @param  {Array}  fallbackFields
     * @return {String}
     */
    static normalizeSearchFilter(searchTerm, fallbackFields) {
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
    }
}
