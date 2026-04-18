import { ClientResponseError } from "pocketbase";

// {
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
//     payload: {},
// }
export function onrecordsave(props) {
    try {
        const val = props.record[props.field.name];
        if (typeof val == "string") {
            JSON.parse(val);
        }
    } catch (err) {
        // simulate API error
        throw new ClientResponseError({
            status: 400,
            response: {
                message: "Invalid JSON data",
                data: {
                    [props.field.name]: {
                        code: "invalid_json",
                        message: err.toString(),
                    },
                },
            },
        });
    }
}
