import { DateTime } from "luxon";

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
            (value === "0001-01-01T00:00:00Z") || // zero time
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
     * Normalizes and returns arr as a valid array instance (if not already).
     *
     * @param  {Array}   arr
     * @param  {Boolean} [allowEmpty]
     * @return {Array}
     */
    static toArray(arr, allowEmpty = false) {
        if (Array.isArray(arr)) {
            return arr;
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
     * @param  {Array}  objectsArr
     * @param  {Object} item
     * @param  {Mixed}  [key]
     * @return {Array}
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
            if (typeof result[prop] === 'object' && result[prop] !== null) {
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
        let parts  = (path || '').split(delimiter);

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
        if (!CommonHelper.isObject(data) && !Array.isArray(data)) {
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
        let parts    = (path || '').split(delimiter);
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
     * Normalizes and converts the provided string to a slug.
     *
     * @param  {String} str
     * @param  {String} [delimiter]
     * @param  {Array}  [preserved]
     * @return {String}
     */
    static slugify(str, delimiter = '_', preserved = ['.', '=', '-']) {
        if (str === '') {
            return '';
        }

        // special characters
        const specialCharsMap = {
            'a': /а|à|á|å|â/gi,
            'b': /б/gi,
            'c': /ц|ç/gi,
            'd': /д/gi,
            'e': /е|è|é|ê|ẽ|ë/gi,
            'f': /ф/gi,
            'g': /г/gi,
            'h': /х/gi,
            'i': /й|и|ì|í|î/gi,
            'j': /ж/gi,
            'k': /к/gi,
            'l': /л/gi,
            'm': /м/gi,
            'n': /н|ñ/gi,
            'o': /о|ò|ó|ô|ø/gi,
            'p': /п/gi,
            'q': /я/gi,
            'r': /р/gi,
            's': /с/gi,
            't': /т/gi,
            'u': /ю|ù|ú|ů|û/gi,
            'v': /в/gi,
            'w': /в/gi,
            'x': /ь/gi,
            'y': /ъ/gi,
            'z': /з/gi,
            'ae': /ä|æ/gi,
            'oe': /ö/gi,
            'ue': /ü/gi,
            'Ae': /Ä/gi,
            'Ue': /Ü/gi,
            'Oe': /Ö/gi,
            'ss': /ß/gi,
            'and': /&/gi
        };

        // replace special characters
        for (let k in specialCharsMap) {
            str = str.replace(specialCharsMap[k], k);
        }

        const slug = str
            .replace(new RegExp('[' + preserved.join('') + ']', 'g'), ' ') // replace preserved characters with spaces
            .replace(/[^\w\ ]/gi, '')                                      // replaces all non-alphanumeric with empty string
            .replace(/\s+/g, delimiter);                                   // collapse whitespaces and replace with `delimiter`

        return slug.charAt(0).toLowerCase() + slug.slice(1);
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
     * Returns a DateTime instance from a date object/string.
     *
     * @param  {String|Date} date
     * @return {DateTime}
     */
    static getDateTime(date) {
        if (typeof date === 'string') {
            const sFormat = "yyyy-MM-dd HH:mm:ss";
            const msFormat = "yyyy-MM-dd HH:mm:ss.SSS";
            const format = date.length === msFormat.length ? msFormat : sFormat;
            return DateTime.fromFormat(date, format, { zone: 'UTC' });
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
    static formatToUTCDate(date, format = 'yyyy-MM-dd HH:mm:ss') {
        return CommonHelper.getDateTime(date).toUTC().toFormat(format);
    }

    /**
     * Returns formatted datetime string in the local timezone.
     *
     * @param  {String|Date} date
     * @param  {String}      [format] The result format (see https://moment.github.io/luxon/#/parsing?id=table-of-tokens)
     * @return {String}
     */
    static formatToLocalDate(date, format = 'yyyy-MM-dd HH:mm:ss') {
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
     * Downloads a json file created from the provide object.
     *
     * @param {mixed} obj   The JS object to download.
     * @param {String} name  The result file name.
     */
    static downloadJson(obj, name) {
        const encodedObj = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(obj, null, 2));

        const tempLink = document.createElement('a');
        tempLink.setAttribute("href", encodedObj);
        tempLink.setAttribute("download", name + ".json");
        tempLink.click();
        tempLink.remove();
    }

    /**
     * Parses and returns the decoded jwt payload data.
     *
     * @param  {String} jwt
     * @return {Object}
     */
    static getJWTPayload(jwt) {
        const raw = (jwt || '').split(".")[1] || '';
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
     * Loosely check if a file is an image based on its filename extension.
     *
     * @param  {String} filename
     * @return {Boolean}
     */
    static hasImageExtension(filename) {
        return /\.jpg|\.jpeg|\.png|\.svg|\.gif|\.webp|\.avif$/.test(filename)
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
     * Returns a dummy collection record object.
     *
     * @param  {Object} collection
     * @return {Object}
     */
    static dummyCollectionRecord(collection) {
        const fields = collection?.schema || [];

        const dummy = {
            "@collectionId": collection?.id,
            "@collectionName": collection?.name,
            "id": "RECORD_ID",
            "created": "2022-01-01 01:00:00",
            "updated": "2022-01-01 23:59:59",
        };

        for (const field of fields) {
            let val = null;
            if (field.type === 'number') {
                val = 123;
            } else if (field.type === "date") {
                val = "2022-01-01 10:00:00";
            } else if (field.type === "bool") {
                val = true;
            } else if (field.type === "email") {
                val = "test@example.com";
            } else if (field.type === "url") {
                val = "https://example.com";
            } else if (field.type === "json") {
                val = 'JSON (array/object)';
            } else if (field.type === "file") {
                val = 'filename.jpg';
                if (field.options?.maxSelect > 1) {
                    val = [val];
                }
            } else if (field.type === "select") {
                val = field.options?.values?.[0];
                if (field.options?.maxSelect > 1) {
                    val = [val];
                }
            } else if (field.type === "relation" || field.type === "user") {
                val = 'RELATION_RECORD_ID';
                if (field.options?.maxSelect > 1) {
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
        field = field || {};

        switch (field.type) {
            case 'bool':
                return 'Boolean';
            case 'number':
                return 'Number';
            case 'file':
                return 'File';
            case 'select':
            case 'relation':
            case 'user':
                if (field.options?.maxSelect > 1) {
                    return 'Array<String>';
                }
                return 'String';
            default:
                return 'String';
        }
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
}
