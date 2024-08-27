import{S as $t,i as Mt,s as qt,C as I,O as Z,N as Ot,e as r,v as b,b as f,c as he,f as v,g as i,h as s,m as ye,w as J,P as Ee,Q as _t,k as Ht,R as Rt,n as Dt,t as ce,a as pe,o as d,d as ke,A as Lt,q as ve,r as Pt,x as ee}from"./index-D0DO79Dq.js";import{S as Ft}from"./SdkTabs-DC6EUYpr.js";import{F as At}from"./FieldsQueryParam-BwleQAus.js";function ht(c,e,t){const n=c.slice();return n[8]=e[t],n}function yt(c,e,t){const n=c.slice();return n[8]=e[t],n}function kt(c,e,t){const n=c.slice();return n[13]=e[t],n}function vt(c){let e;return{c(){e=r("p"),e.innerHTML=`<em>Note that in case of a password change all previously issued tokens for the current record
                will be automatically invalidated and if you want your user to remain signed in you need to
                reauthenticate manually after the update call.</em>`},m(t,n){i(t,e,n)},d(t){t&&d(e)}}}function gt(c){let e;return{c(){e=r("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",v(e,"class","txt-hint txt-sm txt-right")},m(t,n){i(t,e,n)},d(t){t&&d(e)}}}function wt(c){let e,t,n,u,m,o,p,h,w,S,g,$,P,E,M,U,F;return{c(){e=r("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',t=f(),n=r("tr"),n.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>username</span></div></td> <td><span class="label">String</span></td> <td>The username of the auth record.</td>',u=f(),m=r("tr"),m.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address.
                    <br/>
                    This field can be updated only by admins or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,o=f(),p=r("tr"),p.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>emailVisibility</span></div></td> <td><span class="label">Boolean</span></td> <td>Whether to show/hide the auth record email when fetching the record data.</td>',h=f(),w=r("tr"),w.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>oldPassword</span></div></td> <td><span class="label">String</span></td> <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Admins and auth records with
                    &quot;Manage&quot; access can skip this field.</td>`,S=f(),g=r("tr"),g.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>New auth record password.</td>',$=f(),P=r("tr"),P.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>New auth record password confirmation.</td>',E=f(),M=r("tr"),M.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,U=f(),F=r("tr"),F.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>'},m(y,_){i(y,e,_),i(y,t,_),i(y,n,_),i(y,u,_),i(y,m,_),i(y,o,_),i(y,p,_),i(y,h,_),i(y,w,_),i(y,S,_),i(y,g,_),i(y,$,_),i(y,P,_),i(y,E,_),i(y,M,_),i(y,U,_),i(y,F,_)},d(y){y&&(d(e),d(t),d(n),d(u),d(m),d(o),d(p),d(h),d(w),d(S),d(g),d($),d(P),d(E),d(M),d(U),d(F))}}}function Nt(c){let e;return{c(){e=r("span"),e.textContent="Optional",v(e,"class","label label-warning")},m(t,n){i(t,e,n)},d(t){t&&d(e)}}}function Bt(c){let e;return{c(){e=r("span"),e.textContent="Required",v(e,"class","label label-success")},m(t,n){i(t,e,n)},d(t){t&&d(e)}}}function jt(c){var m;let e,t=((m=c[13].options)==null?void 0:m.maxSelect)>1?"ids":"id",n,u;return{c(){e=b("User "),n=b(t),u=b(".")},m(o,p){i(o,e,p),i(o,n,p),i(o,u,p)},p(o,p){var h;p&1&&t!==(t=((h=o[13].options)==null?void 0:h.maxSelect)>1?"ids":"id")&&J(n,t)},d(o){o&&(d(e),d(n),d(u))}}}function Et(c){var m;let e,t=((m=c[13].options)==null?void 0:m.maxSelect)>1?"ids":"id",n,u;return{c(){e=b("Relation record "),n=b(t),u=b(".")},m(o,p){i(o,e,p),i(o,n,p),i(o,u,p)},p(o,p){var h;p&1&&t!==(t=((h=o[13].options)==null?void 0:h.maxSelect)>1?"ids":"id")&&J(n,t)},d(o){o&&(d(e),d(n),d(u))}}}function Ut(c){let e,t,n,u,m;return{c(){e=b("File object."),t=r("br"),n=b(`
                        Set to `),u=r("code"),u.textContent="null",m=b(" to delete already uploaded file(s).")},m(o,p){i(o,e,p),i(o,t,p),i(o,n,p),i(o,u,p),i(o,m,p)},p:ee,d(o){o&&(d(e),d(t),d(n),d(u),d(m))}}}function It(c){let e;return{c(){e=b("URL address.")},m(t,n){i(t,e,n)},p:ee,d(t){t&&d(e)}}}function Jt(c){let e;return{c(){e=b("Email address.")},m(t,n){i(t,e,n)},p:ee,d(t){t&&d(e)}}}function Vt(c){let e;return{c(){e=b("JSON array or object.")},m(t,n){i(t,e,n)},p:ee,d(t){t&&d(e)}}}function Qt(c){let e;return{c(){e=b("Number value.")},m(t,n){i(t,e,n)},p:ee,d(t){t&&d(e)}}}function xt(c){let e;return{c(){e=b("Plain text value.")},m(t,n){i(t,e,n)},p:ee,d(t){t&&d(e)}}}function Tt(c,e){let t,n,u,m,o,p=e[13].name+"",h,w,S,g,$=I.getFieldValueType(e[13])+"",P,E,M,U;function F(C,O){return C[13].required?Bt:Nt}let y=F(e),_=y(e);function A(C,O){if(C[13].type==="text")return xt;if(C[13].type==="number")return Qt;if(C[13].type==="json")return Vt;if(C[13].type==="email")return Jt;if(C[13].type==="url")return It;if(C[13].type==="file")return Ut;if(C[13].type==="relation")return Et;if(C[13].type==="user")return jt}let N=A(e),T=N&&N(e);return{key:c,first:null,c(){t=r("tr"),n=r("td"),u=r("div"),_.c(),m=f(),o=r("span"),h=b(p),w=f(),S=r("td"),g=r("span"),P=b($),E=f(),M=r("td"),T&&T.c(),U=f(),v(u,"class","inline-flex"),v(g,"class","label"),this.first=t},m(C,O){i(C,t,O),s(t,n),s(n,u),_.m(u,null),s(u,m),s(u,o),s(o,h),s(t,w),s(t,S),s(S,g),s(g,P),s(t,E),s(t,M),T&&T.m(M,null),s(t,U)},p(C,O){e=C,y!==(y=F(e))&&(_.d(1),_=y(e),_&&(_.c(),_.m(u,m))),O&1&&p!==(p=e[13].name+"")&&J(h,p),O&1&&$!==($=I.getFieldValueType(e[13])+"")&&J(P,$),N===(N=A(e))&&T?T.p(e,O):(T&&T.d(1),T=N&&N(e),T&&(T.c(),T.m(M,null)))},d(C){C&&d(t),_.d(),T&&T.d()}}}function Ct(c,e){let t,n=e[8].code+"",u,m,o,p;function h(){return e[7](e[8])}return{key:c,first:null,c(){t=r("button"),u=b(n),m=f(),v(t,"class","tab-item"),ve(t,"active",e[1]===e[8].code),this.first=t},m(w,S){i(w,t,S),s(t,u),s(t,m),o||(p=Pt(t,"click",h),o=!0)},p(w,S){e=w,S&4&&n!==(n=e[8].code+"")&&J(u,n),S&6&&ve(t,"active",e[1]===e[8].code)},d(w){w&&d(t),o=!1,p()}}}function St(c,e){let t,n,u,m;return n=new Ot({props:{content:e[8].body}}),{key:c,first:null,c(){t=r("div"),he(n.$$.fragment),u=f(),v(t,"class","tab-item"),ve(t,"active",e[1]===e[8].code),this.first=t},m(o,p){i(o,t,p),ye(n,t,null),s(t,u),m=!0},p(o,p){e=o;const h={};p&4&&(h.content=e[8].body),n.$set(h),(!m||p&6)&&ve(t,"active",e[1]===e[8].code)},i(o){m||(ce(n.$$.fragment,o),m=!0)},o(o){pe(n.$$.fragment,o),m=!1},d(o){o&&d(t),ke(n)}}}function zt(c){var ct,pt,ut;let e,t,n=c[0].name+"",u,m,o,p,h,w,S,g=c[0].name+"",$,P,E,M,U,F,y,_,A,N,T,C,O,ue,Ue,fe,Y,Ie,ge,me=c[0].name+"",we,Je,Te,Ve,Ce,te,Se,le,Oe,ne,$e,V,Me,Qe,Q,qe,B=[],xe=new Map,He,ae,Re,x,De,ze,se,z,Le,Ke,Pe,We,q,Ye,G,Ge,Xe,Ze,Fe,et,Ae,tt,Ne,lt,nt,X,Be,ie,je,K,de,j=[],at=new Map,st,oe,H=[],it=new Map,W,R=c[6]&&vt();A=new Ft({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${c[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},c[3],I.dummyCollectionSchemaData(c[0])),null,4)};

const record = await pb.collection('${(ct=c[0])==null?void 0:ct.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${c[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},c[3],I.dummyCollectionSchemaData(c[0])),null,2)};

final record = await pb.collection('${(pt=c[0])==null?void 0:pt.name}').update('RECORD_ID', body: body);
    `}});let D=c[5]&&gt(),L=c[6]&&wt(),be=Z((ut=c[0])==null?void 0:ut.schema);const dt=l=>l[13].name;for(let l=0;l<be.length;l+=1){let a=kt(c,be,l),k=dt(a);xe.set(k,B[l]=Tt(k,a))}G=new Ot({props:{content:"?expand=relField1,relField2.subRelField21"}}),X=new At({});let _e=Z(c[2]);const ot=l=>l[8].code;for(let l=0;l<_e.length;l+=1){let a=yt(c,_e,l),k=ot(a);at.set(k,j[l]=Ct(k,a))}let re=Z(c[2]);const rt=l=>l[8].code;for(let l=0;l<re.length;l+=1){let a=ht(c,re,l),k=rt(a);it.set(k,H[l]=St(k,a))}return{c(){e=r("h3"),t=b("Update ("),u=b(n),m=b(")"),o=f(),p=r("div"),h=r("p"),w=b("Update a single "),S=r("strong"),$=b(g),P=b(" record."),E=f(),M=r("p"),M.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,U=f(),F=r("p"),F.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,y=f(),R&&R.c(),_=f(),he(A.$$.fragment),N=f(),T=r("h6"),T.textContent="API details",C=f(),O=r("div"),ue=r("strong"),ue.textContent="PATCH",Ue=f(),fe=r("div"),Y=r("p"),Ie=b("/api/collections/"),ge=r("strong"),we=b(me),Je=b("/records/"),Te=r("strong"),Te.textContent=":id",Ve=f(),D&&D.c(),Ce=f(),te=r("div"),te.textContent="Path parameters",Se=f(),le=r("table"),le.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to update.</td></tr></tbody>',Oe=f(),ne=r("div"),ne.textContent="Body Parameters",$e=f(),V=r("table"),Me=r("thead"),Me.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',Qe=f(),Q=r("tbody"),L&&L.c(),qe=f();for(let l=0;l<B.length;l+=1)B[l].c();He=f(),ae=r("div"),ae.textContent="Query parameters",Re=f(),x=r("table"),De=r("thead"),De.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',ze=f(),se=r("tbody"),z=r("tr"),Le=r("td"),Le.textContent="expand",Ke=f(),Pe=r("td"),Pe.innerHTML='<span class="label">String</span>',We=f(),q=r("td"),Ye=b(`Auto expand relations when returning the updated record. Ex.:
                `),he(G.$$.fragment),Ge=b(`
                Supports up to 6-levels depth nested relations expansion. `),Xe=r("br"),Ze=b(`
                The expanded relations will be appended to the record under the
                `),Fe=r("code"),Fe.textContent="expand",et=b(" property (eg. "),Ae=r("code"),Ae.textContent='"expand": {"relField1": {...}, ...}',tt=b(`). Only
                the relations that the user has permissions to `),Ne=r("strong"),Ne.textContent="view",lt=b(" will be expanded."),nt=f(),he(X.$$.fragment),Be=f(),ie=r("div"),ie.textContent="Responses",je=f(),K=r("div"),de=r("div");for(let l=0;l<j.length;l+=1)j[l].c();st=f(),oe=r("div");for(let l=0;l<H.length;l+=1)H[l].c();v(e,"class","m-b-sm"),v(p,"class","content txt-lg m-b-sm"),v(T,"class","m-b-xs"),v(ue,"class","label label-primary"),v(fe,"class","content"),v(O,"class","alert alert-warning"),v(te,"class","section-title"),v(le,"class","table-compact table-border m-b-base"),v(ne,"class","section-title"),v(V,"class","table-compact table-border m-b-base"),v(ae,"class","section-title"),v(x,"class","table-compact table-border m-b-lg"),v(ie,"class","section-title"),v(de,"class","tabs-header compact combined left"),v(oe,"class","tabs-content"),v(K,"class","tabs")},m(l,a){i(l,e,a),s(e,t),s(e,u),s(e,m),i(l,o,a),i(l,p,a),s(p,h),s(h,w),s(h,S),s(S,$),s(h,P),s(p,E),s(p,M),s(p,U),s(p,F),s(p,y),R&&R.m(p,null),i(l,_,a),ye(A,l,a),i(l,N,a),i(l,T,a),i(l,C,a),i(l,O,a),s(O,ue),s(O,Ue),s(O,fe),s(fe,Y),s(Y,Ie),s(Y,ge),s(ge,we),s(Y,Je),s(Y,Te),s(O,Ve),D&&D.m(O,null),i(l,Ce,a),i(l,te,a),i(l,Se,a),i(l,le,a),i(l,Oe,a),i(l,ne,a),i(l,$e,a),i(l,V,a),s(V,Me),s(V,Qe),s(V,Q),L&&L.m(Q,null),s(Q,qe);for(let k=0;k<B.length;k+=1)B[k]&&B[k].m(Q,null);i(l,He,a),i(l,ae,a),i(l,Re,a),i(l,x,a),s(x,De),s(x,ze),s(x,se),s(se,z),s(z,Le),s(z,Ke),s(z,Pe),s(z,We),s(z,q),s(q,Ye),ye(G,q,null),s(q,Ge),s(q,Xe),s(q,Ze),s(q,Fe),s(q,et),s(q,Ae),s(q,tt),s(q,Ne),s(q,lt),s(se,nt),ye(X,se,null),i(l,Be,a),i(l,ie,a),i(l,je,a),i(l,K,a),s(K,de);for(let k=0;k<j.length;k+=1)j[k]&&j[k].m(de,null);s(K,st),s(K,oe);for(let k=0;k<H.length;k+=1)H[k]&&H[k].m(oe,null);W=!0},p(l,[a]){var ft,mt,bt;(!W||a&1)&&n!==(n=l[0].name+"")&&J(u,n),(!W||a&1)&&g!==(g=l[0].name+"")&&J($,g),l[6]?R||(R=vt(),R.c(),R.m(p,null)):R&&(R.d(1),R=null);const k={};a&25&&(k.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},l[3],I.dummyCollectionSchemaData(l[0])),null,4)};

const record = await pb.collection('${(ft=l[0])==null?void 0:ft.name}').update('RECORD_ID', data);
    `),a&25&&(k.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},l[3],I.dummyCollectionSchemaData(l[0])),null,2)};

final record = await pb.collection('${(mt=l[0])==null?void 0:mt.name}').update('RECORD_ID', body: body);
    `),A.$set(k),(!W||a&1)&&me!==(me=l[0].name+"")&&J(we,me),l[5]?D||(D=gt(),D.c(),D.m(O,null)):D&&(D.d(1),D=null),l[6]?L||(L=wt(),L.c(),L.m(Q,qe)):L&&(L.d(1),L=null),a&1&&(be=Z((bt=l[0])==null?void 0:bt.schema),B=Ee(B,a,dt,1,l,be,xe,Q,_t,Tt,null,kt)),a&6&&(_e=Z(l[2]),j=Ee(j,a,ot,1,l,_e,at,de,_t,Ct,null,yt)),a&6&&(re=Z(l[2]),Ht(),H=Ee(H,a,rt,1,l,re,it,oe,Rt,St,null,ht),Dt())},i(l){if(!W){ce(A.$$.fragment,l),ce(G.$$.fragment,l),ce(X.$$.fragment,l);for(let a=0;a<re.length;a+=1)ce(H[a]);W=!0}},o(l){pe(A.$$.fragment,l),pe(G.$$.fragment,l),pe(X.$$.fragment,l);for(let a=0;a<H.length;a+=1)pe(H[a]);W=!1},d(l){l&&(d(e),d(o),d(p),d(_),d(N),d(T),d(C),d(O),d(Ce),d(te),d(Se),d(le),d(Oe),d(ne),d($e),d(V),d(He),d(ae),d(Re),d(x),d(Be),d(ie),d(je),d(K)),R&&R.d(),ke(A,l),D&&D.d(),L&&L.d();for(let a=0;a<B.length;a+=1)B[a].d();ke(G),ke(X);for(let a=0;a<j.length;a+=1)j[a].d();for(let a=0;a<H.length;a+=1)H[a].d()}}}function Kt(c,e,t){let n,u,m,{collection:o}=e,p=200,h=[],w={};const S=g=>t(1,p=g.code);return c.$$set=g=>{"collection"in g&&t(0,o=g.collection)},c.$$.update=()=>{var g,$;c.$$.dirty&1&&t(6,n=(o==null?void 0:o.type)==="auth"),c.$$.dirty&1&&t(5,u=(o==null?void 0:o.updateRule)===null),c.$$.dirty&1&&t(2,h=[{code:200,body:JSON.stringify(I.dummyCollectionRecord(o),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${($=(g=o==null?void 0:o.schema)==null?void 0:g[0])==null?void 0:$.name}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `},{code:403,body:`
                {
                  "code": 403,
                  "message": "You are not allowed to perform this request.",
                  "data": {}
                }
            `},{code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}]),c.$$.dirty&1&&(o.type==="auth"?t(3,w={username:"test_username_update",emailVisibility:!1,password:"87654321",passwordConfirm:"87654321",oldPassword:"12345678"}):t(3,w={}))},t(4,m=I.getApiExampleUrl(Lt.baseUrl)),[o,p,h,w,m,u,n,S]}class Xt extends $t{constructor(e){super(),Mt(this,e,Kt,zt,qt,{collection:0})}}export{Xt as default};
