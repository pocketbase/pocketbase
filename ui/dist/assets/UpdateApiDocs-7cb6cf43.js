import{S as $t,i as Ot,s as Mt,C as E,M as St,e as r,w as y,b as m,c as be,f as T,g as a,h as i,m as me,x as U,N as Ne,P as mt,k as qt,Q as Dt,n as Ht,t as de,a as re,o,d as _e,T as Pt,p as Rt,r as ye,u as Lt,y as X}from"./index-a65ca895.js";import{S as Ft}from"./SdkTabs-ad912c8f.js";import{F as At}from"./FieldsQueryParam-ba250473.js";function _t(f,t,l){const s=f.slice();return s[7]=t[l],s}function yt(f,t,l){const s=f.slice();return s[7]=t[l],s}function kt(f,t,l){const s=f.slice();return s[12]=t[l],s}function ht(f){let t;return{c(){t=r("p"),t.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",T(t,"class","txt-hint txt-sm txt-right")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function vt(f){let t,l,s,b,u,d,p,k,C,w,$,P,F,j,O,g,A;return{c(){t=r("tr"),t.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',l=m(),s=r("tr"),s.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>username</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The username of the auth record.</td>`,b=m(),u=r("tr"),u.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>email</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The auth record email address.
                    <br/>
                    This field can be updated only by admins or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,d=m(),p=r("tr"),p.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>emailVisibility</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Whether to show/hide the auth record email when fetching the record data.</td>`,k=m(),C=r("tr"),C.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>oldPassword</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Admins and auth records with
                    &quot;Manage&quot; access can skip this field.</td>`,w=m(),$=r("tr"),$.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>password</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>New auth record password.</td>`,P=m(),F=r("tr"),F.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>passwordConfirm</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>New auth record password confirmation.</td>`,j=m(),O=r("tr"),O.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>verified</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,g=m(),A=r("tr"),A.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>'},m(c,_){a(c,t,_),a(c,l,_),a(c,s,_),a(c,b,_),a(c,u,_),a(c,d,_),a(c,p,_),a(c,k,_),a(c,C,_),a(c,w,_),a(c,$,_),a(c,P,_),a(c,F,_),a(c,j,_),a(c,O,_),a(c,g,_),a(c,A,_)},d(c){c&&o(t),c&&o(l),c&&o(s),c&&o(b),c&&o(u),c&&o(d),c&&o(p),c&&o(k),c&&o(C),c&&o(w),c&&o($),c&&o(P),c&&o(F),c&&o(j),c&&o(O),c&&o(g),c&&o(A)}}}function Bt(f){let t;return{c(){t=r("span"),t.textContent="Optional",T(t,"class","label label-warning")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function Nt(f){let t;return{c(){t=r("span"),t.textContent="Required",T(t,"class","label label-success")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function jt(f){var u;let t,l=((u=f[12].options)==null?void 0:u.maxSelect)>1?"ids":"id",s,b;return{c(){t=y("User "),s=y(l),b=y(".")},m(d,p){a(d,t,p),a(d,s,p),a(d,b,p)},p(d,p){var k;p&1&&l!==(l=((k=d[12].options)==null?void 0:k.maxSelect)>1?"ids":"id")&&U(s,l)},d(d){d&&o(t),d&&o(s),d&&o(b)}}}function gt(f){var u;let t,l=((u=f[12].options)==null?void 0:u.maxSelect)>1?"ids":"id",s,b;return{c(){t=y("Relation record "),s=y(l),b=y(".")},m(d,p){a(d,t,p),a(d,s,p),a(d,b,p)},p(d,p){var k;p&1&&l!==(l=((k=d[12].options)==null?void 0:k.maxSelect)>1?"ids":"id")&&U(s,l)},d(d){d&&o(t),d&&o(s),d&&o(b)}}}function Et(f){let t,l,s,b,u;return{c(){t=y("File object."),l=r("br"),s=y(`
                        Set to `),b=r("code"),b.textContent="null",u=y(" to delete already uploaded file(s).")},m(d,p){a(d,t,p),a(d,l,p),a(d,s,p),a(d,b,p),a(d,u,p)},p:X,d(d){d&&o(t),d&&o(l),d&&o(s),d&&o(b),d&&o(u)}}}function Ut(f){let t;return{c(){t=y("URL address.")},m(l,s){a(l,t,s)},p:X,d(l){l&&o(t)}}}function It(f){let t;return{c(){t=y("Email address.")},m(l,s){a(l,t,s)},p:X,d(l){l&&o(t)}}}function Jt(f){let t;return{c(){t=y("JSON array or object.")},m(l,s){a(l,t,s)},p:X,d(l){l&&o(t)}}}function Vt(f){let t;return{c(){t=y("Number value.")},m(l,s){a(l,t,s)},p:X,d(l){l&&o(t)}}}function xt(f){let t;return{c(){t=y("Plain text value.")},m(l,s){a(l,t,s)},p:X,d(l){l&&o(t)}}}function wt(f,t){let l,s,b,u,d,p=t[12].name+"",k,C,w,$,P=E.getFieldValueType(t[12])+"",F,j,O,g;function A(h,L){return h[12].required?Nt:Bt}let c=A(t),_=c(t);function K(h,L){if(h[12].type==="text")return xt;if(h[12].type==="number")return Vt;if(h[12].type==="json")return Jt;if(h[12].type==="email")return It;if(h[12].type==="url")return Ut;if(h[12].type==="file")return Et;if(h[12].type==="relation")return gt;if(h[12].type==="user")return jt}let R=K(t),S=R&&R(t);return{key:f,first:null,c(){l=r("tr"),s=r("td"),b=r("div"),_.c(),u=m(),d=r("span"),k=y(p),C=m(),w=r("td"),$=r("span"),F=y(P),j=m(),O=r("td"),S&&S.c(),g=m(),T(b,"class","inline-flex"),T($,"class","label"),this.first=l},m(h,L){a(h,l,L),i(l,s),i(s,b),_.m(b,null),i(b,u),i(b,d),i(d,k),i(l,C),i(l,w),i(w,$),i($,F),i(l,j),i(l,O),S&&S.m(O,null),i(l,g)},p(h,L){t=h,c!==(c=A(t))&&(_.d(1),_=c(t),_&&(_.c(),_.m(b,u))),L&1&&p!==(p=t[12].name+"")&&U(k,p),L&1&&P!==(P=E.getFieldValueType(t[12])+"")&&U(F,P),R===(R=K(t))&&S?S.p(t,L):(S&&S.d(1),S=R&&R(t),S&&(S.c(),S.m(O,null)))},d(h){h&&o(l),_.d(),S&&S.d()}}}function Tt(f,t){let l,s=t[7].code+"",b,u,d,p;function k(){return t[6](t[7])}return{key:f,first:null,c(){l=r("button"),b=y(s),u=m(),T(l,"class","tab-item"),ye(l,"active",t[1]===t[7].code),this.first=l},m(C,w){a(C,l,w),i(l,b),i(l,u),d||(p=Lt(l,"click",k),d=!0)},p(C,w){t=C,w&4&&s!==(s=t[7].code+"")&&U(b,s),w&6&&ye(l,"active",t[1]===t[7].code)},d(C){C&&o(l),d=!1,p()}}}function Ct(f,t){let l,s,b,u;return s=new St({props:{content:t[7].body}}),{key:f,first:null,c(){l=r("div"),be(s.$$.fragment),b=m(),T(l,"class","tab-item"),ye(l,"active",t[1]===t[7].code),this.first=l},m(d,p){a(d,l,p),me(s,l,null),i(l,b),u=!0},p(d,p){t=d;const k={};p&4&&(k.content=t[7].body),s.$set(k),(!u||p&6)&&ye(l,"active",t[1]===t[7].code)},i(d){u||(de(s.$$.fragment,d),u=!0)},o(d){re(s.$$.fragment,d),u=!1},d(d){d&&o(l),_e(s)}}}function Qt(f){var ot,dt,rt,ft;let t,l,s=f[0].name+"",b,u,d,p,k,C,w,$=f[0].name+"",P,F,j,O,g,A,c,_,K,R,S,h,L,je,fe,W,ge,ke,pe=f[0].name+"",he,Ee,ve,Ue,we,Z,Te,ee,Ce,te,Se,I,$e,Ie,J,Oe,B=[],Je=new Map,Me,le,qe,V,De,Ve,se,x,He,xe,Pe,Qe,M,ze,Y,Ke,We,Ye,Re,Ge,Le,Xe,Fe,Ze,et,G,Ae,ne,Be,Q,ie,N=[],tt=new Map,lt,ae,q=[],st=new Map,z;_=new Ft({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${f[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},f[3],E.dummyCollectionSchemaData(f[0])),null,4)};

const record = await pb.collection('${(ot=f[0])==null?void 0:ot.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${f[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},f[3],E.dummyCollectionSchemaData(f[0])),null,2)};

final record = await pb.collection('${(dt=f[0])==null?void 0:dt.name}').update('RECORD_ID', body: body);
    `}});let D=f[5]&&ht(),H=((rt=f[0])==null?void 0:rt.isAuth)&&vt(),ce=(ft=f[0])==null?void 0:ft.schema;const nt=e=>e[12].name;for(let e=0;e<ce.length;e+=1){let n=kt(f,ce,e),v=nt(n);Je.set(v,B[e]=wt(v,n))}Y=new St({props:{content:"?expand=relField1,relField2.subRelField21"}}),G=new At({});let ue=f[2];const it=e=>e[7].code;for(let e=0;e<ue.length;e+=1){let n=yt(f,ue,e),v=it(n);tt.set(v,N[e]=Tt(v,n))}let oe=f[2];const at=e=>e[7].code;for(let e=0;e<oe.length;e+=1){let n=_t(f,oe,e),v=at(n);st.set(v,q[e]=Ct(v,n))}return{c(){t=r("h3"),l=y("Update ("),b=y(s),u=y(")"),d=m(),p=r("div"),k=r("p"),C=y("Update a single "),w=r("strong"),P=y($),F=y(" record."),j=m(),O=r("p"),O.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,g=m(),A=r("p"),A.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,c=m(),be(_.$$.fragment),K=m(),R=r("h6"),R.textContent="API details",S=m(),h=r("div"),L=r("strong"),L.textContent="PATCH",je=m(),fe=r("div"),W=r("p"),ge=y("/api/collections/"),ke=r("strong"),he=y(pe),Ee=y("/records/"),ve=r("strong"),ve.textContent=":id",Ue=m(),D&&D.c(),we=m(),Z=r("div"),Z.textContent="Path parameters",Te=m(),ee=r("table"),ee.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to update.</td></tr></tbody>`,Ce=m(),te=r("div"),te.textContent="Body Parameters",Se=m(),I=r("table"),$e=r("thead"),$e.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,Ie=m(),J=r("tbody"),H&&H.c(),Oe=m();for(let e=0;e<B.length;e+=1)B[e].c();Me=m(),le=r("div"),le.textContent="Query parameters",qe=m(),V=r("table"),De=r("thead"),De.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Ve=m(),se=r("tbody"),x=r("tr"),He=r("td"),He.textContent="expand",xe=m(),Pe=r("td"),Pe.innerHTML='<span class="label">String</span>',Qe=m(),M=r("td"),ze=y(`Auto expand relations when returning the updated record. Ex.:
                `),be(Y.$$.fragment),Ke=y(`
                Supports up to 6-levels depth nested relations expansion. `),We=r("br"),Ye=y(`
                The expanded relations will be appended to the record under the
                `),Re=r("code"),Re.textContent="expand",Ge=y(" property (eg. "),Le=r("code"),Le.textContent='"expand": {"relField1": {...}, ...}',Xe=y(`). Only
                the relations that the user has permissions to `),Fe=r("strong"),Fe.textContent="view",Ze=y(" will be expanded."),et=m(),be(G.$$.fragment),Ae=m(),ne=r("div"),ne.textContent="Responses",Be=m(),Q=r("div"),ie=r("div");for(let e=0;e<N.length;e+=1)N[e].c();lt=m(),ae=r("div");for(let e=0;e<q.length;e+=1)q[e].c();T(t,"class","m-b-sm"),T(p,"class","content txt-lg m-b-sm"),T(R,"class","m-b-xs"),T(L,"class","label label-primary"),T(fe,"class","content"),T(h,"class","alert alert-warning"),T(Z,"class","section-title"),T(ee,"class","table-compact table-border m-b-base"),T(te,"class","section-title"),T(I,"class","table-compact table-border m-b-base"),T(le,"class","section-title"),T(V,"class","table-compact table-border m-b-lg"),T(ne,"class","section-title"),T(ie,"class","tabs-header compact left"),T(ae,"class","tabs-content"),T(Q,"class","tabs")},m(e,n){a(e,t,n),i(t,l),i(t,b),i(t,u),a(e,d,n),a(e,p,n),i(p,k),i(k,C),i(k,w),i(w,P),i(k,F),i(p,j),i(p,O),i(p,g),i(p,A),a(e,c,n),me(_,e,n),a(e,K,n),a(e,R,n),a(e,S,n),a(e,h,n),i(h,L),i(h,je),i(h,fe),i(fe,W),i(W,ge),i(W,ke),i(ke,he),i(W,Ee),i(W,ve),i(h,Ue),D&&D.m(h,null),a(e,we,n),a(e,Z,n),a(e,Te,n),a(e,ee,n),a(e,Ce,n),a(e,te,n),a(e,Se,n),a(e,I,n),i(I,$e),i(I,Ie),i(I,J),H&&H.m(J,null),i(J,Oe);for(let v=0;v<B.length;v+=1)B[v]&&B[v].m(J,null);a(e,Me,n),a(e,le,n),a(e,qe,n),a(e,V,n),i(V,De),i(V,Ve),i(V,se),i(se,x),i(x,He),i(x,xe),i(x,Pe),i(x,Qe),i(x,M),i(M,ze),me(Y,M,null),i(M,Ke),i(M,We),i(M,Ye),i(M,Re),i(M,Ge),i(M,Le),i(M,Xe),i(M,Fe),i(M,Ze),i(se,et),me(G,se,null),a(e,Ae,n),a(e,ne,n),a(e,Be,n),a(e,Q,n),i(Q,ie);for(let v=0;v<N.length;v+=1)N[v]&&N[v].m(ie,null);i(Q,lt),i(Q,ae);for(let v=0;v<q.length;v+=1)q[v]&&q[v].m(ae,null);z=!0},p(e,[n]){var pt,ct,ut,bt;(!z||n&1)&&s!==(s=e[0].name+"")&&U(b,s),(!z||n&1)&&$!==($=e[0].name+"")&&U(P,$);const v={};n&25&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${e[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},e[3],E.dummyCollectionSchemaData(e[0])),null,4)};

const record = await pb.collection('${(pt=e[0])==null?void 0:pt.name}').update('RECORD_ID', data);
    `),n&25&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${e[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},e[3],E.dummyCollectionSchemaData(e[0])),null,2)};

final record = await pb.collection('${(ct=e[0])==null?void 0:ct.name}').update('RECORD_ID', body: body);
    `),_.$set(v),(!z||n&1)&&pe!==(pe=e[0].name+"")&&U(he,pe),e[5]?D||(D=ht(),D.c(),D.m(h,null)):D&&(D.d(1),D=null),(ut=e[0])!=null&&ut.isAuth?H||(H=vt(),H.c(),H.m(J,Oe)):H&&(H.d(1),H=null),n&1&&(ce=(bt=e[0])==null?void 0:bt.schema,B=Ne(B,n,nt,1,e,ce,Je,J,mt,wt,null,kt)),n&6&&(ue=e[2],N=Ne(N,n,it,1,e,ue,tt,ie,mt,Tt,null,yt)),n&6&&(oe=e[2],qt(),q=Ne(q,n,at,1,e,oe,st,ae,Dt,Ct,null,_t),Ht())},i(e){if(!z){de(_.$$.fragment,e),de(Y.$$.fragment,e),de(G.$$.fragment,e);for(let n=0;n<oe.length;n+=1)de(q[n]);z=!0}},o(e){re(_.$$.fragment,e),re(Y.$$.fragment,e),re(G.$$.fragment,e);for(let n=0;n<q.length;n+=1)re(q[n]);z=!1},d(e){e&&o(t),e&&o(d),e&&o(p),e&&o(c),_e(_,e),e&&o(K),e&&o(R),e&&o(S),e&&o(h),D&&D.d(),e&&o(we),e&&o(Z),e&&o(Te),e&&o(ee),e&&o(Ce),e&&o(te),e&&o(Se),e&&o(I),H&&H.d();for(let n=0;n<B.length;n+=1)B[n].d();e&&o(Me),e&&o(le),e&&o(qe),e&&o(V),_e(Y),_e(G),e&&o(Ae),e&&o(ne),e&&o(Be),e&&o(Q);for(let n=0;n<N.length;n+=1)N[n].d();for(let n=0;n<q.length;n+=1)q[n].d()}}}function zt(f,t,l){let s,b,{collection:u=new Pt}=t,d=200,p=[],k={};const C=w=>l(1,d=w.code);return f.$$set=w=>{"collection"in w&&l(0,u=w.collection)},f.$$.update=()=>{var w,$;f.$$.dirty&1&&l(5,s=(u==null?void 0:u.updateRule)===null),f.$$.dirty&1&&l(2,p=[{code:200,body:JSON.stringify(E.dummyCollectionRecord(u),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${($=(w=u==null?void 0:u.schema)==null?void 0:w[0])==null?void 0:$.name}": {
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
            `}]),f.$$.dirty&1&&(u.$isAuth?l(3,k={username:"test_username_update",emailVisibility:!1,password:"87654321",passwordConfirm:"87654321",oldPassword:"12345678"}):l(3,k={}))},l(4,b=E.getApiExampleUrl(Rt.baseUrl)),[u,d,p,k,b,s,C]}class Gt extends $t{constructor(t){super(),Ot(this,t,zt,Qt,Mt,{collection:0})}}export{Gt as default};
