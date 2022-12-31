import{S as Ct,i as St,s as Ot,C as I,N as Tt,e as r,w as y,b as m,c as Ae,f as T,g as a,h as i,m as Be,x as U,O as Pe,P as ut,k as Mt,Q as $t,n as Rt,t as pe,a as fe,o,d as Fe,R as qt,p as Dt,r as ce,u as Ht,y as G}from"./index.89a3f554.js";import{S as Lt}from"./SdkTabs.0a6ad1c9.js";function bt(p,t,l){const s=p.slice();return s[7]=t[l],s}function mt(p,t,l){const s=p.slice();return s[7]=t[l],s}function _t(p,t,l){const s=p.slice();return s[12]=t[l],s}function yt(p){let t;return{c(){t=r("p"),t.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",T(t,"class","txt-hint txt-sm txt-right")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function kt(p){let t,l,s,b,u,d,f,k,C,v,O,D,A,F,M,N,B;return{c(){t=r("tr"),t.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',l=m(),s=r("tr"),s.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>username</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The username of the auth record.</td>`,b=m(),u=r("tr"),u.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>email</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The auth record email address.
                    <br/>
                    This field can be updated only by admins or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,d=m(),f=r("tr"),f.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>emailVisibility</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Whether to show/hide the auth record email when fetching the record data.</td>`,k=m(),C=r("tr"),C.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>oldPassword</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Admins and auth records with
                    &quot;Manage&quot; access can skip this field.</td>`,v=m(),O=r("tr"),O.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>password</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>New auth record password.</td>`,D=m(),A=r("tr"),A.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>passwordConfirm</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>New auth record password confirmation.</td>`,F=m(),M=r("tr"),M.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>verified</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,N=m(),B=r("tr"),B.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>'},m(c,_){a(c,t,_),a(c,l,_),a(c,s,_),a(c,b,_),a(c,u,_),a(c,d,_),a(c,f,_),a(c,k,_),a(c,C,_),a(c,v,_),a(c,O,_),a(c,D,_),a(c,A,_),a(c,F,_),a(c,M,_),a(c,N,_),a(c,B,_)},d(c){c&&o(t),c&&o(l),c&&o(s),c&&o(b),c&&o(u),c&&o(d),c&&o(f),c&&o(k),c&&o(C),c&&o(v),c&&o(O),c&&o(D),c&&o(A),c&&o(F),c&&o(M),c&&o(N),c&&o(B)}}}function Pt(p){let t;return{c(){t=r("span"),t.textContent="Optional",T(t,"class","label label-warning")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function At(p){let t;return{c(){t=r("span"),t.textContent="Required",T(t,"class","label label-success")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function Bt(p){var u;let t,l=((u=p[12].options)==null?void 0:u.maxSelect)>1?"ids":"id",s,b;return{c(){t=y("User "),s=y(l),b=y(".")},m(d,f){a(d,t,f),a(d,s,f),a(d,b,f)},p(d,f){var k;f&1&&l!==(l=((k=d[12].options)==null?void 0:k.maxSelect)>1?"ids":"id")&&U(s,l)},d(d){d&&o(t),d&&o(s),d&&o(b)}}}function Ft(p){var u;let t,l=((u=p[12].options)==null?void 0:u.maxSelect)>1?"ids":"id",s,b;return{c(){t=y("Relation record "),s=y(l),b=y(".")},m(d,f){a(d,t,f),a(d,s,f),a(d,b,f)},p(d,f){var k;f&1&&l!==(l=((k=d[12].options)==null?void 0:k.maxSelect)>1?"ids":"id")&&U(s,l)},d(d){d&&o(t),d&&o(s),d&&o(b)}}}function Nt(p){let t,l,s,b,u;return{c(){t=y("File object."),l=r("br"),s=y(`
                        Set to `),b=r("code"),b.textContent="null",u=y(" to delete already uploaded file(s).")},m(d,f){a(d,t,f),a(d,l,f),a(d,s,f),a(d,b,f),a(d,u,f)},p:G,d(d){d&&o(t),d&&o(l),d&&o(s),d&&o(b),d&&o(u)}}}function jt(p){let t;return{c(){t=y("URL address.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function Et(p){let t;return{c(){t=y("Email address.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function It(p){let t;return{c(){t=y("JSON array or object.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function Ut(p){let t;return{c(){t=y("Number value.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function Jt(p){let t;return{c(){t=y("Plain text value.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function ht(p,t){let l,s,b,u,d,f=t[12].name+"",k,C,v,O,D=I.getFieldValueType(t[12])+"",A,F,M,N;function B(h,L){return h[12].required?At:Pt}let c=B(t),_=c(t);function K(h,L){if(h[12].type==="text")return Jt;if(h[12].type==="number")return Ut;if(h[12].type==="json")return It;if(h[12].type==="email")return Et;if(h[12].type==="url")return jt;if(h[12].type==="file")return Nt;if(h[12].type==="relation")return Ft;if(h[12].type==="user")return Bt}let H=K(t),S=H&&H(t);return{key:p,first:null,c(){l=r("tr"),s=r("td"),b=r("div"),_.c(),u=m(),d=r("span"),k=y(f),C=m(),v=r("td"),O=r("span"),A=y(D),F=m(),M=r("td"),S&&S.c(),N=m(),T(b,"class","inline-flex"),T(O,"class","label"),this.first=l},m(h,L){a(h,l,L),i(l,s),i(s,b),_.m(b,null),i(b,u),i(b,d),i(d,k),i(l,C),i(l,v),i(v,O),i(O,A),i(l,F),i(l,M),S&&S.m(M,null),i(l,N)},p(h,L){t=h,c!==(c=B(t))&&(_.d(1),_=c(t),_&&(_.c(),_.m(b,u))),L&1&&f!==(f=t[12].name+"")&&U(k,f),L&1&&D!==(D=I.getFieldValueType(t[12])+"")&&U(A,D),H===(H=K(t))&&S?S.p(t,L):(S&&S.d(1),S=H&&H(t),S&&(S.c(),S.m(M,null)))},d(h){h&&o(l),_.d(),S&&S.d()}}}function vt(p,t){let l,s=t[7].code+"",b,u,d,f;function k(){return t[6](t[7])}return{key:p,first:null,c(){l=r("button"),b=y(s),u=m(),T(l,"class","tab-item"),ce(l,"active",t[1]===t[7].code),this.first=l},m(C,v){a(C,l,v),i(l,b),i(l,u),d||(f=Ht(l,"click",k),d=!0)},p(C,v){t=C,v&4&&s!==(s=t[7].code+"")&&U(b,s),v&6&&ce(l,"active",t[1]===t[7].code)},d(C){C&&o(l),d=!1,f()}}}function wt(p,t){let l,s,b,u;return s=new Tt({props:{content:t[7].body}}),{key:p,first:null,c(){l=r("div"),Ae(s.$$.fragment),b=m(),T(l,"class","tab-item"),ce(l,"active",t[1]===t[7].code),this.first=l},m(d,f){a(d,l,f),Be(s,l,null),i(l,b),u=!0},p(d,f){t=d;const k={};f&4&&(k.content=t[7].body),s.$set(k),(!u||f&6)&&ce(l,"active",t[1]===t[7].code)},i(d){u||(pe(s.$$.fragment,d),u=!0)},o(d){fe(s.$$.fragment,d),u=!1},d(d){d&&o(l),Fe(s)}}}function gt(p){var it,at,ot,dt;let t,l,s=p[0].name+"",b,u,d,f,k,C,v,O=p[0].name+"",D,A,F,M,N,B,c,_,K,H,S,h,L,Ne,ae,W,je,ue,oe=p[0].name+"",be,Ee,me,Ie,_e,X,ye,Z,ke,ee,he,J,ve,Ue,g,we,j=[],Je=new Map,Te,te,Ce,V,Se,ge,Oe,x,Me,Ve,$e,xe,$,Qe,Y,ze,Ke,We,Re,Ye,qe,Ge,De,Xe,He,le,Le,Q,se,E=[],Ze=new Map,et,ne,P=[],tt=new Map,z;_=new Lt({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${p[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},p[3],I.dummyCollectionSchemaData(p[0])),null,4)};

const record = await pb.collection('${(it=p[0])==null?void 0:it.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${p[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},p[3],I.dummyCollectionSchemaData(p[0])),null,2)};

final record = await pb.collection('${(at=p[0])==null?void 0:at.name}').update('RECORD_ID', body: body);
    `}});let R=p[5]&&yt(),q=((ot=p[0])==null?void 0:ot.isAuth)&&kt(),de=(dt=p[0])==null?void 0:dt.schema;const lt=e=>e[12].name;for(let e=0;e<de.length;e+=1){let n=_t(p,de,e),w=lt(n);Je.set(w,j[e]=ht(w,n))}Y=new Tt({props:{content:"?expand=relField1,relField2.subRelField21"}});let re=p[2];const st=e=>e[7].code;for(let e=0;e<re.length;e+=1){let n=mt(p,re,e),w=st(n);Ze.set(w,E[e]=vt(w,n))}let ie=p[2];const nt=e=>e[7].code;for(let e=0;e<ie.length;e+=1){let n=bt(p,ie,e),w=nt(n);tt.set(w,P[e]=wt(w,n))}return{c(){t=r("h3"),l=y("Update ("),b=y(s),u=y(")"),d=m(),f=r("div"),k=r("p"),C=y("Update a single "),v=r("strong"),D=y(O),A=y(" record."),F=m(),M=r("p"),M.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,N=m(),B=r("p"),B.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,c=m(),Ae(_.$$.fragment),K=m(),H=r("h6"),H.textContent="API details",S=m(),h=r("div"),L=r("strong"),L.textContent="PATCH",Ne=m(),ae=r("div"),W=r("p"),je=y("/api/collections/"),ue=r("strong"),be=y(oe),Ee=y("/records/"),me=r("strong"),me.textContent=":id",Ie=m(),R&&R.c(),_e=m(),X=r("div"),X.textContent="Path parameters",ye=m(),Z=r("table"),Z.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to update.</td></tr></tbody>`,ke=m(),ee=r("div"),ee.textContent="Body Parameters",he=m(),J=r("table"),ve=r("thead"),ve.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,Ue=m(),g=r("tbody"),q&&q.c(),we=m();for(let e=0;e<j.length;e+=1)j[e].c();Te=m(),te=r("div"),te.textContent="Query parameters",Ce=m(),V=r("table"),Se=r("thead"),Se.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,ge=m(),Oe=r("tbody"),x=r("tr"),Me=r("td"),Me.textContent="expand",Ve=m(),$e=r("td"),$e.innerHTML='<span class="label">String</span>',xe=m(),$=r("td"),Qe=y(`Auto expand relations when returning the updated record. Ex.:
                `),Ae(Y.$$.fragment),ze=y(`
                Supports up to 6-levels depth nested relations expansion. `),Ke=r("br"),We=y(`
                The expanded relations will be appended to the record under the
                `),Re=r("code"),Re.textContent="expand",Ye=y(" property (eg. "),qe=r("code"),qe.textContent='"expand": {"relField1": {...}, ...}',Ge=y(`). Only
                the relations that the user has permissions to `),De=r("strong"),De.textContent="view",Xe=y(" will be expanded."),He=m(),le=r("div"),le.textContent="Responses",Le=m(),Q=r("div"),se=r("div");for(let e=0;e<E.length;e+=1)E[e].c();et=m(),ne=r("div");for(let e=0;e<P.length;e+=1)P[e].c();T(t,"class","m-b-sm"),T(f,"class","content txt-lg m-b-sm"),T(H,"class","m-b-xs"),T(L,"class","label label-primary"),T(ae,"class","content"),T(h,"class","alert alert-warning"),T(X,"class","section-title"),T(Z,"class","table-compact table-border m-b-base"),T(ee,"class","section-title"),T(J,"class","table-compact table-border m-b-base"),T(te,"class","section-title"),T(V,"class","table-compact table-border m-b-lg"),T(le,"class","section-title"),T(se,"class","tabs-header compact left"),T(ne,"class","tabs-content"),T(Q,"class","tabs")},m(e,n){a(e,t,n),i(t,l),i(t,b),i(t,u),a(e,d,n),a(e,f,n),i(f,k),i(k,C),i(k,v),i(v,D),i(k,A),i(f,F),i(f,M),i(f,N),i(f,B),a(e,c,n),Be(_,e,n),a(e,K,n),a(e,H,n),a(e,S,n),a(e,h,n),i(h,L),i(h,Ne),i(h,ae),i(ae,W),i(W,je),i(W,ue),i(ue,be),i(W,Ee),i(W,me),i(h,Ie),R&&R.m(h,null),a(e,_e,n),a(e,X,n),a(e,ye,n),a(e,Z,n),a(e,ke,n),a(e,ee,n),a(e,he,n),a(e,J,n),i(J,ve),i(J,Ue),i(J,g),q&&q.m(g,null),i(g,we);for(let w=0;w<j.length;w+=1)j[w].m(g,null);a(e,Te,n),a(e,te,n),a(e,Ce,n),a(e,V,n),i(V,Se),i(V,ge),i(V,Oe),i(Oe,x),i(x,Me),i(x,Ve),i(x,$e),i(x,xe),i(x,$),i($,Qe),Be(Y,$,null),i($,ze),i($,Ke),i($,We),i($,Re),i($,Ye),i($,qe),i($,Ge),i($,De),i($,Xe),a(e,He,n),a(e,le,n),a(e,Le,n),a(e,Q,n),i(Q,se);for(let w=0;w<E.length;w+=1)E[w].m(se,null);i(Q,et),i(Q,ne);for(let w=0;w<P.length;w+=1)P[w].m(ne,null);z=!0},p(e,[n]){var rt,pt,ft,ct;(!z||n&1)&&s!==(s=e[0].name+"")&&U(b,s),(!z||n&1)&&O!==(O=e[0].name+"")&&U(D,O);const w={};n&25&&(w.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${e[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},e[3],I.dummyCollectionSchemaData(e[0])),null,4)};

const record = await pb.collection('${(rt=e[0])==null?void 0:rt.name}').update('RECORD_ID', data);
    `),n&25&&(w.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${e[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},e[3],I.dummyCollectionSchemaData(e[0])),null,2)};

final record = await pb.collection('${(pt=e[0])==null?void 0:pt.name}').update('RECORD_ID', body: body);
    `),_.$set(w),(!z||n&1)&&oe!==(oe=e[0].name+"")&&U(be,oe),e[5]?R||(R=yt(),R.c(),R.m(h,null)):R&&(R.d(1),R=null),(ft=e[0])!=null&&ft.isAuth?q||(q=kt(),q.c(),q.m(g,we)):q&&(q.d(1),q=null),n&1&&(de=(ct=e[0])==null?void 0:ct.schema,j=Pe(j,n,lt,1,e,de,Je,g,ut,ht,null,_t)),n&6&&(re=e[2],E=Pe(E,n,st,1,e,re,Ze,se,ut,vt,null,mt)),n&6&&(ie=e[2],Mt(),P=Pe(P,n,nt,1,e,ie,tt,ne,$t,wt,null,bt),Rt())},i(e){if(!z){pe(_.$$.fragment,e),pe(Y.$$.fragment,e);for(let n=0;n<ie.length;n+=1)pe(P[n]);z=!0}},o(e){fe(_.$$.fragment,e),fe(Y.$$.fragment,e);for(let n=0;n<P.length;n+=1)fe(P[n]);z=!1},d(e){e&&o(t),e&&o(d),e&&o(f),e&&o(c),Fe(_,e),e&&o(K),e&&o(H),e&&o(S),e&&o(h),R&&R.d(),e&&o(_e),e&&o(X),e&&o(ye),e&&o(Z),e&&o(ke),e&&o(ee),e&&o(he),e&&o(J),q&&q.d();for(let n=0;n<j.length;n+=1)j[n].d();e&&o(Te),e&&o(te),e&&o(Ce),e&&o(V),Fe(Y),e&&o(He),e&&o(le),e&&o(Le),e&&o(Q);for(let n=0;n<E.length;n+=1)E[n].d();for(let n=0;n<P.length;n+=1)P[n].d()}}}function Vt(p,t,l){let s,b,{collection:u=new qt}=t,d=200,f=[],k={};const C=v=>l(1,d=v.code);return p.$$set=v=>{"collection"in v&&l(0,u=v.collection)},p.$$.update=()=>{var v,O;p.$$.dirty&1&&l(5,s=(u==null?void 0:u.updateRule)===null),p.$$.dirty&1&&l(2,f=[{code:200,body:JSON.stringify(I.dummyCollectionRecord(u),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${(O=(v=u==null?void 0:u.schema)==null?void 0:v[0])==null?void 0:O.name}": {
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
            `}]),p.$$.dirty&1&&(u.isAuth?l(3,k={username:"test_username_update",emailVisibility:!1,password:"87654321",passwordConfirm:"87654321",oldPassword:"12345678"}):l(3,k={}))},l(4,b=I.getApiExampleUrl(Dt.baseUrl)),[u,d,f,k,b,s,C]}class zt extends Ct{constructor(t){super(),St(this,t,Vt,gt,Ot,{collection:0})}}export{zt as default};
