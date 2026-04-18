export const filesToDeleteProp = "@@filesToDelete"; // symbols are not used because they are not reactive

// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
//     payload: {},
// }
export function onrecordsave(props) {
    const files = app.utils.toArray(props.payload[props.field.name]);

    const toDelete = app.utils.toArray(props.record[filesToDeleteProp]?.[props.field.name]);

    for (let filename of toDelete) {
        const index = files.indexOf(filename);
        if (index >= 0) {
            files.splice(index, 1);
        }
    }

    props.payload[props.field.name] = files;
}
