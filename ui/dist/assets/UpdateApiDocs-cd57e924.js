import{S as Ct,i as St,s as Ot,C as E,N as Tt,e as r,w as y,b,c as be,f as w,g as o,h as a,m as me,x as U,O as je,P as ut,k as $t,Q as Mt,n as qt,t as de,a as re,o as d,d as _e,p as Dt,r as ye,u as Ht,y as X}from"./index-e9c07574.js";import{S as Pt}from"./SdkTabs-31f0c5d0.js";import{F as Rt}from"./FieldsQueryParam-152f9b5c.js";function bt(f,t,l){const n=f.slice();return n[8]=t[l],n}function mt(f,t,l){const n=f.slice();return n[8]=t[l],n}function _t(f,t,l){const n=f.slice();return n[13]=t[l],n}function yt(f){let t;return{c(){t=r("p"),t.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",w(t,"class","txt-hint txt-sm txt-right")},m(l,n){o(l,t,n)},d(l){l&&d(t)}}}function kt(f){let t,l,n,u,_,s,p,k,C,S,T,$,F,A,M,g,B;return{c(){t=r("tr"),t.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',l=b(),n=r("tr"),n.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>username</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The username of the auth record.</td>`,u=b(),_=r("tr"),_.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>email</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The auth record email address.
                    <br/>
                    This field can be updated only by admins or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,s=b(),p=r("tr"),p.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>emailVisibility</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Whether to show/hide the auth record email when fetching the record data.</td>`,k=b(),C=r("tr"),C.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>oldPassword</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Admins and auth records with
                    &quot;Manage&quot; access can skip this field.</td>`,S=b(),T=r("tr"),T.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>password</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>New auth record password.</td>`,$=b(),F=r("tr"),F.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>passwordConfirm</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>New auth record password confirmation.</td>`,A=b(),M=r("tr"),M.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>verified</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,g=b(),B=r("tr"),B.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>'},m(c,m){o(c,t,m),o(c,l,m),o(c,n,m),o(c,u,m),o(c,_,m),o(c,s,m),o(c,p,m),o(c,k,m),o(c,C,m),o(c,S,m),o(c,T,m),o(c,$,m),o(c,F,m),o(c,A,m),o(c,M,m),o(c,g,m),o(c,B,m)},d(c){c&&d(t),c&&d(l),c&&d(n),c&&d(u),c&&d(_),c&&d(s),c&&d(p),c&&d(k),c&&d(C),c&&d(S),c&&d(T),c&&d($),c&&d(F),c&&d(A),c&&d(M),c&&d(g),c&&d(B)}}}function Lt(f){let t;return{c(){t=r("span"),t.textContent="Optional",w(t,"class","label label-warning")},m(l,n){o(l,t,n)},d(l){l&&d(t)}}}function Ft(f){let t;return{c(){t=r("span"),t.textContent="Required",w(t,"class","label label-success")},m(l,n){o(l,t,n)},d(l){l&&d(t)}}}function Bt(f){var _;let t,l=((_=f[13].options)==null?void 0:_.maxSelect)>1?"ids":"id",n,u;return{c(){t=y("User "),n=y(l),u=y(".")},m(s,p){o(s,t,p),o(s,n,p),o(s,u,p)},p(s,p){var k;p&1&&l!==(l=((k=s[13].options)==null?void 0:k.maxSelect)>1?"ids":"id")&&U(n,l)},d(s){s&&d(t),s&&d(n),s&&d(u)}}}function Nt(f){var _;let t,l=((_=f[13].options)==null?void 0:_.maxSelect)>1?"ids":"id",n,u;return{c(){t=y("Relation record "),n=y(l),u=y(".")},m(s,p){o(s,t,p),o(s,n,p),o(s,u,p)},p(s,p){var k;p&1&&l!==(l=((k=s[13].options)==null?void 0:k.maxSelect)>1?"ids":"id")&&U(n,l)},d(s){s&&d(t),s&&d(n),s&&d(u)}}}function jt(f){let t,l,n,u,_;return{c(){t=y("File object."),l=r("br"),n=y(`
                        Set to `),u=r("code"),u.textContent="null",_=y(" to delete already uploaded file(s).")},m(s,p){o(s,t,p),o(s,l,p),o(s,n,p),o(s,u,p),o(s,_,p)},p:X,d(s){s&&d(t),s&&d(l),s&&d(n),s&&d(u),s&&d(_)}}}function At(f){let t;return{c(){t=y("URL address.")},m(l,n){o(l,t,n)},p:X,d(l){l&&d(t)}}}function gt(f){let t;return{c(){t=y("Email address.")},m(l,n){o(l,t,n)},p:X,d(l){l&&d(t)}}}function Et(f){let t;return{c(){t=y("JSON array or object.")},m(l,n){o(l,t,n)},p:X,d(l){l&&d(t)}}}function Ut(f){let t;return{c(){t=y("Number value.")},m(l,n){o(l,t,n)},p:X,d(l){l&&d(t)}}}function It(f){let t;return{c(){t=y("Plain text value.")},m(l,n){o(l,t,n)},p:X,d(l){l&&d(t)}}}function ht(f,t){let l,n,u,_,s,p=t[13].name+"",k,C,S,T,$=E.getFieldValueType(t[13])+"",F,A,M,g;function B(h,L){return h[13].required?Ft:Lt}let c=B(t),m=c(t);function K(h,L){if(h[13].type==="text")return It;if(h[13].type==="number")return Ut;if(h[13].type==="json")return Et;if(h[13].type==="email")return gt;if(h[13].type==="url")return At;if(h[13].type==="file")return jt;if(h[13].type==="relation")return Nt;if(h[13].type==="user")return Bt}let R=K(t),O=R&&R(t);return{key:f,first:null,c(){l=r("tr"),n=r("td"),u=r("div"),m.c(),_=b(),s=r("span"),k=y(p),C=b(),S=r("td"),T=r("span"),F=y($),A=b(),M=r("td"),O&&O.c(),g=b(),w(u,"class","inline-flex"),w(T,"class","label"),this.first=l},m(h,L){o(h,l,L),a(l,n),a(n,u),m.m(u,null),a(u,_),a(u,s),a(s,k),a(l,C),a(l,S),a(S,T),a(T,F),a(l,A),a(l,M),O&&O.m(M,null),a(l,g)},p(h,L){t=h,c!==(c=B(t))&&(m.d(1),m=c(t),m&&(m.c(),m.m(u,_))),L&1&&p!==(p=t[13].name+"")&&U(k,p),L&1&&$!==($=E.getFieldValueType(t[13])+"")&&U(F,$),R===(R=K(t))&&O?O.p(t,L):(O&&O.d(1),O=R&&R(t),O&&(O.c(),O.m(M,null)))},d(h){h&&d(l),m.d(),O&&O.d()}}}function vt(f,t){let l,n=t[8].code+"",u,_,s,p;function k(){return t[7](t[8])}return{key:f,first:null,c(){l=r("button"),u=y(n),_=b(),w(l,"class","tab-item"),ye(l,"active",t[1]===t[8].code),this.first=l},m(C,S){o(C,l,S),a(l,u),a(l,_),s||(p=Ht(l,"click",k),s=!0)},p(C,S){t=C,S&4&&n!==(n=t[8].code+"")&&U(u,n),S&6&&ye(l,"active",t[1]===t[8].code)},d(C){C&&d(l),s=!1,p()}}}function wt(f,t){let l,n,u,_;return n=new Tt({props:{content:t[8].body}}),{key:f,first:null,c(){l=r("div"),be(n.$$.fragment),u=b(),w(l,"class","tab-item"),ye(l,"active",t[1]===t[8].code),this.first=l},m(s,p){o(s,l,p),me(n,l,null),a(l,u),_=!0},p(s,p){t=s;const k={};p&4&&(k.content=t[8].body),n.$set(k),(!_||p&6)&&ye(l,"active",t[1]===t[8].code)},i(s){_||(de(n.$$.fragment,s),_=!0)},o(s){re(n.$$.fragment,s),_=!1},d(s){s&&d(l),_e(n)}}}function Jt(f){var ot,dt,rt;let t,l,n=f[0].name+"",u,_,s,p,k,C,S,T=f[0].name+"",$,F,A,M,g,B,c,m,K,R,O,h,L,Ae,fe,W,ge,ke,pe=f[0].name+"",he,Ee,ve,Ue,we,Z,Te,ee,Ce,te,Se,I,Oe,Ie,J,$e,N=[],Je=new Map,Me,le,qe,V,De,Ve,se,x,He,xe,Pe,Qe,q,ze,Y,Ke,We,Ye,Re,Ge,Le,Xe,Fe,Ze,et,G,Be,ne,Ne,Q,ie,j=[],tt=new Map,lt,ae,D=[],st=new Map,z;m=new Pt({props:{js:`
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
    `}});let H=f[5]&&yt(),P=f[6]&&kt(),ce=(rt=f[0])==null?void 0:rt.schema;const nt=e=>e[13].name;for(let e=0;e<ce.length;e+=1){let i=_t(f,ce,e),v=nt(i);Je.set(v,N[e]=ht(v,i))}Y=new Tt({props:{content:"?expand=relField1,relField2.subRelField21"}}),G=new Rt({});let ue=f[2];const it=e=>e[8].code;for(let e=0;e<ue.length;e+=1){let i=mt(f,ue,e),v=it(i);tt.set(v,j[e]=vt(v,i))}let oe=f[2];const at=e=>e[8].code;for(let e=0;e<oe.length;e+=1){let i=bt(f,oe,e),v=at(i);st.set(v,D[e]=wt(v,i))}return{c(){t=r("h3"),l=y("Update ("),u=y(n),_=y(")"),s=b(),p=r("div"),k=r("p"),C=y("Update a single "),S=r("strong"),$=y(T),F=y(" record."),A=b(),M=r("p"),M.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,g=b(),B=r("p"),B.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,c=b(),be(m.$$.fragment),K=b(),R=r("h6"),R.textContent="API details",O=b(),h=r("div"),L=r("strong"),L.textContent="PATCH",Ae=b(),fe=r("div"),W=r("p"),ge=y("/api/collections/"),ke=r("strong"),he=y(pe),Ee=y("/records/"),ve=r("strong"),ve.textContent=":id",Ue=b(),H&&H.c(),we=b(),Z=r("div"),Z.textContent="Path parameters",Te=b(),ee=r("table"),ee.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to update.</td></tr></tbody>`,Ce=b(),te=r("div"),te.textContent="Body Parameters",Se=b(),I=r("table"),Oe=r("thead"),Oe.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,Ie=b(),J=r("tbody"),P&&P.c(),$e=b();for(let e=0;e<N.length;e+=1)N[e].c();Me=b(),le=r("div"),le.textContent="Query parameters",qe=b(),V=r("table"),De=r("thead"),De.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Ve=b(),se=r("tbody"),x=r("tr"),He=r("td"),He.textContent="expand",xe=b(),Pe=r("td"),Pe.innerHTML='<span class="label">String</span>',Qe=b(),q=r("td"),ze=y(`Auto expand relations when returning the updated record. Ex.:
                `),be(Y.$$.fragment),Ke=y(`
                Supports up to 6-levels depth nested relations expansion. `),We=r("br"),Ye=y(`
                The expanded relations will be appended to the record under the
                `),Re=r("code"),Re.textContent="expand",Ge=y(" property (eg. "),Le=r("code"),Le.textContent='"expand": {"relField1": {...}, ...}',Xe=y(`). Only
                the relations that the user has permissions to `),Fe=r("strong"),Fe.textContent="view",Ze=y(" will be expanded."),et=b(),be(G.$$.fragment),Be=b(),ne=r("div"),ne.textContent="Responses",Ne=b(),Q=r("div"),ie=r("div");for(let e=0;e<j.length;e+=1)j[e].c();lt=b(),ae=r("div");for(let e=0;e<D.length;e+=1)D[e].c();w(t,"class","m-b-sm"),w(p,"class","content txt-lg m-b-sm"),w(R,"class","m-b-xs"),w(L,"class","label label-primary"),w(fe,"class","content"),w(h,"class","alert alert-warning"),w(Z,"class","section-title"),w(ee,"class","table-compact table-border m-b-base"),w(te,"class","section-title"),w(I,"class","table-compact table-border m-b-base"),w(le,"class","section-title"),w(V,"class","table-compact table-border m-b-lg"),w(ne,"class","section-title"),w(ie,"class","tabs-header compact left"),w(ae,"class","tabs-content"),w(Q,"class","tabs")},m(e,i){o(e,t,i),a(t,l),a(t,u),a(t,_),o(e,s,i),o(e,p,i),a(p,k),a(k,C),a(k,S),a(S,$),a(k,F),a(p,A),a(p,M),a(p,g),a(p,B),o(e,c,i),me(m,e,i),o(e,K,i),o(e,R,i),o(e,O,i),o(e,h,i),a(h,L),a(h,Ae),a(h,fe),a(fe,W),a(W,ge),a(W,ke),a(ke,he),a(W,Ee),a(W,ve),a(h,Ue),H&&H.m(h,null),o(e,we,i),o(e,Z,i),o(e,Te,i),o(e,ee,i),o(e,Ce,i),o(e,te,i),o(e,Se,i),o(e,I,i),a(I,Oe),a(I,Ie),a(I,J),P&&P.m(J,null),a(J,$e);for(let v=0;v<N.length;v+=1)N[v]&&N[v].m(J,null);o(e,Me,i),o(e,le,i),o(e,qe,i),o(e,V,i),a(V,De),a(V,Ve),a(V,se),a(se,x),a(x,He),a(x,xe),a(x,Pe),a(x,Qe),a(x,q),a(q,ze),me(Y,q,null),a(q,Ke),a(q,We),a(q,Ye),a(q,Re),a(q,Ge),a(q,Le),a(q,Xe),a(q,Fe),a(q,Ze),a(se,et),me(G,se,null),o(e,Be,i),o(e,ne,i),o(e,Ne,i),o(e,Q,i),a(Q,ie);for(let v=0;v<j.length;v+=1)j[v]&&j[v].m(ie,null);a(Q,lt),a(Q,ae);for(let v=0;v<D.length;v+=1)D[v]&&D[v].m(ae,null);z=!0},p(e,[i]){var ft,pt,ct;(!z||i&1)&&n!==(n=e[0].name+"")&&U(u,n),(!z||i&1)&&T!==(T=e[0].name+"")&&U($,T);const v={};i&25&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${e[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},e[3],E.dummyCollectionSchemaData(e[0])),null,4)};

const record = await pb.collection('${(ft=e[0])==null?void 0:ft.name}').update('RECORD_ID', data);
    `),i&25&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${e[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},e[3],E.dummyCollectionSchemaData(e[0])),null,2)};

final record = await pb.collection('${(pt=e[0])==null?void 0:pt.name}').update('RECORD_ID', body: body);
    `),m.$set(v),(!z||i&1)&&pe!==(pe=e[0].name+"")&&U(he,pe),e[5]?H||(H=yt(),H.c(),H.m(h,null)):H&&(H.d(1),H=null),e[6]?P||(P=kt(),P.c(),P.m(J,$e)):P&&(P.d(1),P=null),i&1&&(ce=(ct=e[0])==null?void 0:ct.schema,N=je(N,i,nt,1,e,ce,Je,J,ut,ht,null,_t)),i&6&&(ue=e[2],j=je(j,i,it,1,e,ue,tt,ie,ut,vt,null,mt)),i&6&&(oe=e[2],$t(),D=je(D,i,at,1,e,oe,st,ae,Mt,wt,null,bt),qt())},i(e){if(!z){de(m.$$.fragment,e),de(Y.$$.fragment,e),de(G.$$.fragment,e);for(let i=0;i<oe.length;i+=1)de(D[i]);z=!0}},o(e){re(m.$$.fragment,e),re(Y.$$.fragment,e),re(G.$$.fragment,e);for(let i=0;i<D.length;i+=1)re(D[i]);z=!1},d(e){e&&d(t),e&&d(s),e&&d(p),e&&d(c),_e(m,e),e&&d(K),e&&d(R),e&&d(O),e&&d(h),H&&H.d(),e&&d(we),e&&d(Z),e&&d(Te),e&&d(ee),e&&d(Ce),e&&d(te),e&&d(Se),e&&d(I),P&&P.d();for(let i=0;i<N.length;i+=1)N[i].d();e&&d(Me),e&&d(le),e&&d(qe),e&&d(V),_e(Y),_e(G),e&&d(Be),e&&d(ne),e&&d(Ne),e&&d(Q);for(let i=0;i<j.length;i+=1)j[i].d();for(let i=0;i<D.length;i+=1)D[i].d()}}}function Vt(f,t,l){let n,u,_,{collection:s}=t,p=200,k=[],C={};const S=T=>l(1,p=T.code);return f.$$set=T=>{"collection"in T&&l(0,s=T.collection)},f.$$.update=()=>{var T,$;f.$$.dirty&1&&l(6,n=(s==null?void 0:s.type)==="auth"),f.$$.dirty&1&&l(5,u=(s==null?void 0:s.updateRule)===null),f.$$.dirty&1&&l(2,k=[{code:200,body:JSON.stringify(E.dummyCollectionRecord(s),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${($=(T=s==null?void 0:s.schema)==null?void 0:T[0])==null?void 0:$.name}": {
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
            `}]),f.$$.dirty&1&&(s.type==="auth"?l(3,C={username:"test_username_update",emailVisibility:!1,password:"87654321",passwordConfirm:"87654321",oldPassword:"12345678"}):l(3,C={}))},l(4,_=E.getApiExampleUrl(Dt.baseUrl)),[s,p,k,C,_,u,n,S]}class Kt extends Ct{constructor(t){super(),St(this,t,Vt,Jt,Ot,{collection:0})}}export{Kt as default};
