<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
</script>

<tr>
    <td id="query-page">fields</td>
    <td>
        <span class="label">String</span>
    </td>
    <td>
        <p>
            Comma separated string of the fields to return in the JSON response
            <em>(by default returns all fields)</em>. Ex.:
            <CodeBlock
                content={`
                    ?fields=*,expand.relField.name
                `}
            />
        </p>
        <p>
            <code>*</code> targets all keys from the specific depth level.
        </p>
        <p>In addition, the following field modifiers are also supported:</p>
        <ul>
            <li>
                <code>:excerpt(maxLength, withEllipsis?)</code>
                <br />
                Returns a short plain text version of the field string value.
                <br />
                Ex.:
                <code>?fields=*,description:excerpt(200,true)</code>
            </li>
        </ul>
    </td>
</tr>
