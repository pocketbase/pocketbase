// {
//     collection: {},
//     field: {},
//     originalRecord: {},
//     clone: {},
// }
export function onrecordduplicate(props) {
    delete props.clone[props.field.name];
}
