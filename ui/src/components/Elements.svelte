<script>
    import tooltip from "@/actions/tooltip";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import { addInfoToast } from "@/stores/toasts";

    let popupActive = true;

    setTimeout(function () {
        addInfoToast("Hello world");
    }, 500);
</script>

<div class="form-field">
    <label for="">EXAMPLE</label>
    <ObjectSelect multiple searchable items={["test1", "test2"]} />
</div>
<hr />
<div class="form-field">
    <label for="">EXAMPLE</label>
    <ObjectSelect searchable items={["test1", "test2"]} />
</div>
<hr />
<div class="form-field disabled">
    <label for="">EXAMPLE</label>
    <ObjectSelect disabled searchable items={["test1", "test2"]} />
</div>
<hr />

<div class="alert">
    <div class="icon">
        <i class="ri-information-line" />
    </div>
    <div class="content">Hello world!</div>
    <div class="close">
        <i class="ri-close-line" />
    </div>
</div>

<div class="alert alert-info">
    <div class="icon">
        <i class="ri-information-line" />
    </div>
    <div class="content">Hello world!</div>
    <div class="close">
        <i class="ri-close-line" />
    </div>
</div>

<div class="alert alert-danger">
    <div class="icon">
        <i class="ri-information-line" />
    </div>
    <div class="content">Hello world!</div>
    <div class="close">
        <i class="ri-close-line" />
    </div>
</div>

<div class="alert alert-warning">
    <div class="icon">
        <i class="ri-error-warning-line" />
    </div>
    <div class="content">Hello world!</div>
    <div class="close">
        <i class="ri-close-line" />
    </div>
</div>

<div class="alert alert-success">
    <div class="icon">
        <i class="ri-checkbox-circle-line" />
    </div>
    <div class="content">Hello world!</div>
    <div class="close">
        <i class="ri-close-line" />
    </div>
</div>

<hr />

<h1>H1 title</h1>
<p>Lorem Ipsum dolor sit amet...</p>
<h2>H2 title</h2>
<p>Lorem Ipsum dolor sit amet...</p>
<h3>H3 title</h3>
<p>Lorem Ipsum dolor sit amet...</p>
<h4>H4 title</h4>
<p>Lorem Ipsum dolor sit amet...</p>
<h5>H5 title</h5>
<p>Lorem Ipsum dolor sit amet...</p>
<h6>H6 title</h6>
<p>Lorem Ipsum dolor sit amet...</p>

<hr />

<div class="grid">
    <div class="col-6">COL1</div>
    <div class="col-6">COL2</div>
</div>

<p>
    Lorem Ipsum is <a href="/">simply dummy</a> text of the printing and typesetting industry. Lorem Ipsum has
    been the industry's
    <strong>standard</strong> dummy text ever since the 1500s, when an unknown printer took a galley of type
    and <em>scrambled</em> it to make a type specimen book. It has survived not only five centuries, but also
    the leap into electronic typesetting, remaining<sup>1</sup> essentially<sub>2</sub> unchanged.
</p>
<p>
    It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and
    more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum
</p>

<ul>
    <li><small>Option 1</small></li>
    <li>Option 2</li>
    <li>Option 3</li>
</ul>

<ol>
    <li>Option 1</li>
    <li>Option 2</li>
    <li>Option 3</li>
</ol>

<hr />

<span use:tooltip={"My tooltip"}>Lorem Ipsum</span>

<hr />

<button class="btn">Button default</button>
<button class="btn btn-danger">Button danger</button>
<button class="btn btn-warning">Button warning</button>
<button class="btn btn-success">Button success</button>
<button class="btn btn-info">Button info</button>
<button class="btn btn-hint">Button hint</button>

<hr />

<button class="btn btn-secondary">Button default</button>
<button class="btn btn-secondary btn-danger">Button danger</button>
<button class="btn btn-secondary btn-warning">Button danger</button>
<button class="btn btn-secondary btn-success">Button success</button>
<button class="btn btn-secondary btn-info">Button info</button>
<button class="btn btn-secondary btn-hint">Button hint</button>

<hr />

<button class="btn btn-outline">Button default</button>
<button class="btn btn-outline btn-danger">Button danger</button>
<button class="btn btn-outline btn-warning">Button danger</button>
<button class="btn btn-outline btn-success">Button success</button>
<button class="btn btn-outline btn-info">Button info</button>
<button class="btn btn-outline btn-hint">Button hint</button>

<hr />

<button disabled class="btn">Button default</button>
<button disabled class="btn btn-danger">Button danger</button>
<button disabled class="btn btn-warning">Button warning</button>
<button disabled class="btn btn-success">Button success</button>
<button disabled class="btn btn-info">Button info</button>
<button disabled class="btn btn-hint">Button hint</button>

<button disabled class="btn btn-secondary">Button default</button>
<button disabled class="btn btn-secondary btn-danger">Button danger</button>
<button disabled class="btn btn-secondary btn-warning">Button danger</button>
<button disabled class="btn btn-secondary btn-success">Button success</button>
<button disabled class="btn btn-secondary btn-info">Button info</button>
<button disabled class="btn btn-secondary btn-hint">Button hint</button>
<hr />

<button class="btn">
    <i class="ri-mail-line" />
    <span class="txt">Button default</span>
</button>
<button class="btn btn-danger">
    <i class="ri-mail-line" />
    <span class="txt">Button danger</span>
</button>
<button class="btn btn-warning">
    <i class="ri-mail-line" />
    <span class="txt">Button warning</span>
</button>
<button class="btn btn-success">
    <i class="ri-mail-line" />
    <span class="txt">Button success</span>
</button>
<button class="btn btn-hint">
    <i class="ri-mail-line" />
    <span class="txt">Button hint</span>
</button>

<hr />

<button class="btn btn-sm">
    <i class="ri-mail-line" />
    <span class="txt">Button default</span>
</button>
<button class="btn btn-danger btn-sm">
    <i class="ri-mail-line" />
    <span class="txt">Button danger</span>
</button>
<button class="btn btn-warning btn-sm">
    <i class="ri-mail-line" />
    <span class="txt">Button warning</span>
</button>
<button class="btn btn-success btn-sm">
    <i class="ri-mail-line" />
    <span class="txt">Button success</span>
</button>
<button class="btn btn-hint btn-sm">
    <i class="ri-mail-line" />
    <span class="txt">Button hint</span>
</button>

<hr />

<button class="btn btn-lg">
    <i class="ri-mail-line" />
    <span class="txt">Button default</span>
</button>
<button class="btn btn-danger btn-lg">
    <i class="ri-mail-line" />
    <span class="txt">Button danger</span>
</button>
<button class="btn btn-warning btn-lg">
    <i class="ri-mail-line" />
    <span class="txt">Button warning</span>
</button>
<button class="btn btn-success btn-lg">
    <i class="ri-mail-line" />
    <span class="txt">Button success</span>
</button>
<button class="btn btn-hint btn-lg">
    <i class="ri-mail-line" />
    <span class="txt">Button hint</span>
</button>

<hr />

<button class="btn btn-circle">
    <i class="ri-mail-line" />
</button>
<button class="btn btn-sm btn-circle">
    <i class="ri-mail-line" />
</button>
<button class="btn btn-lg btn-circle">
    <i class="ri-mail-line" />
</button>

<button class="btn btn-secondary btn-circle">
    <i class="ri-mail-line" />
</button>
<button class="btn btn-secondary btn-sm btn-circle">
    <i class="ri-mail-line" />
</button>
<button class="btn btn-secondary btn-lg btn-circle">
    <i class="ri-mail-line" />
</button>

<hr />

<button class="btn btn-loading">
    <i class="ri-mail-line" />
    <span class="txt">Button Loading</span>
</button>
<button class="btn btn-loading btn-primary btn-sm">
    <i class="ri-mail-line" />
    <span class="txt">Button Loading</span>
</button>
<button class="btn btn-loading btn-danger btn-lg">
    <i class="ri-mail-line" />
    <span class="txt">Button Loading</span>
</button>

<button class="btn btn-loading btn-circle">
    <i class="ri-mail-line" />
</button>
<button class="btn btn-loading btn-primary btn-sm btn-circle">
    <i class="ri-mail-line" />
</button>
<button class="btn btn-loading btn-danger btn-lg btn-circle">
    <i class="ri-mail-line" />
</button>

<hr />

<button disabled class="btn btn-loading">
    <i class="ri-mail-line" />
    <span class="txt">Button Loading</span>
</button>
<button disabled class="btn btn-loading btn-primary btn-sm">
    <i class="ri-mail-line" />
    <span class="txt">Button Loading</span>
</button>
<button disabled class="btn btn-loading btn-danger btn-lg">
    <i class="ri-mail-line" />
    <span class="txt">Button Loading</span>
</button>

<button disabled class="btn btn-loading btn-circle">
    <i class="ri-mail-line" />
</button>
<button disabled class="btn btn-loading btn-primary btn-sm btn-circle">
    <i class="ri-mail-line" />
</button>
<button disabled class="btn btn-loading btn-danger btn-lg btn-circle">
    <i class="ri-mail-line" />
</button>

<hr />

<input type="text" />

<hr />
<select>
    <option value="1" selected>Option 1</option>
    <option value="">Option 2</option>
    <option value="">Option 3</option>
</select>
<hr />

<textarea cols="30" rows="10" />

<hr />

<div class="form-field required">
    <label for="field_1">Name</label>
    <input type="text" id="field_1" placeholder="Name 123" />
</div>

<div class="form-field required">
    <label for="field_2">Description</label>
    <textarea id="field_2" />
</div>

<div class="form-field">
    <label for="field_3">Choose value</label>
    <select id="field_3">
        <option value="1" selected>Option 1</option>
        <option value="">Option 2</option>
        <option value="">Option 3</option>
    </select>
</div>

<hr />

<div class="form-field">
    <input type="text" placeholder="Lorem ipsum dolor sit amet..." />
</div>

<div class="form-field">
    <textarea />
</div>

<div class="form-field">
    <select>
        <option value="1" selected>Option 1</option>
        <option value="">Option 2</option>
        <option value="">Option 3</option>
    </select>
</div>

<hr />

<div class="form-field">
    <input type="checkbox" id="field_check" />
    <label for="field_check">
        I agree with the <a href="/">terms and conditions</a>
    </label>
</div>

<div class="form-field">
    <input type="radio" name="radio_check" id="field_radio1" value="1" />
    <label for="field_radio1">Radio 1</label>
</div>
<div class="form-field">
    <input type="radio" name="radio_check" id="field_radio2" value="2" />
    <label for="field_radio2">Radio 2</label>
</div>

<div class="form-field form-field-toggle">
    <input type="checkbox" id="field_toggle" />
    <label for="field_toggle">Toggle check</label>
</div>

<hr />

<div class="form-field error">
    <label for="field_error">Name + error</label>
    <input type="text" id="field_error" />
    <div class="help-block help-block-error">
        <p>Something went wrong</p>
    </div>
</div>

<div class="form-field">
    <label for="field_hint">Name + hint</label>
    <input type="text" id="field_hint" />
    <div class="help-block">
        <p>Lorem ipsum dolor <a href="/">sit</a> amet</p>
    </div>
</div>

<div class="form-field">
    <input type="checkbox" id="field_check_hint" />
    <label for="field_check_hint">Checkbox hint</label>
    <div class="help-block">Lorem ipsum</div>
</div>

<div class="form-field has-error">
    <input type="checkbox" id="field_check_error" />
    <label for="field_check_error">Checkbox error</label>
    <div class="help-block help-error">Lorem ipsum</div>
</div>

<hr />

<div class="form-field disabled">
    <input type="checkbox" id="field_check" disabled />
    <label for="field_check">
        I agree with the <a href="/">terms and conditions</a>
    </label>
</div>

<div class="form-field">
    <input type="radio" name="radio_check" id="field_radio1" value="1" disabled />
    <label for="field_radio1">Radio 1</label>
</div>
<div class="form-field">
    <input type="radio" name="radio_check" id="field_radio2" value="2" disabled />
    <label for="field_radio2">Radio 2</label>
</div>

<div class="form-field form-field-toggle disabled">
    <input type="checkbox" id="field_toggle" disabled />
    <label for="field_toggle" />
</div>

<hr />

<div class="form-field disabled">
    <label for="field_addon1">Name</label>
    <div class="form-field-addon">
        <i class="ri-mail-line" />
    </div>
    <input disabled type="text" id="field_addon1" />
</div>

<div class="form-field">
    <div class="form-field-addon">
        <i class="ri-mail-line" />
    </div>
    <input type="text" />
</div>

<hr />

<div class="form-group">
    <div class="form-field">
        <label for="field_group1">Name</label>
        <div class="form-field-addon">
            <div class="btn btn-circle btn-secondary">
                <i class="ri-mail-line" />
            </div>
        </div>
        <input type="text" id="field_group1" />
    </div>
    <div class="form-field">
        <label for="field_group2">Password</label>
        <div class="form-field-addon">
            <i class="ri-mail-line" />
        </div>
        <input type="password" id="field_group2" />
    </div>
</div>

<OverlayPanel bind:active={popupActive} popup={false}>
    <h4 slot="header">My title</h4>
    <p>Lorem ipsum dolor sit amet...</p>
    <svelte:fragment slot="footer">
        <button class="btn btn-secondary">Cancel</button>
        <button class="btn btn-expanded">Save</button>
    </svelte:fragment>
</OverlayPanel>
