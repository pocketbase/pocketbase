import{S as Ct,i as St,s as Ot,C as U,N as Tt,e as r,w as y,b as m,c as Ae,f as T,g as a,h as i,m as Be,x as I,O as Pe,P as ut,k as $t,Q as Mt,n as qt,t as fe,a as pe,o,d as Fe,T as Dt,p as Ht,r as ce,u as Rt,y as G}from"./index-4eea3e34.js";import{S as Lt}from"./SdkTabs-5d6cc1d4.js";function bt(f,t,l){const s=f.slice();return s[7]=t[l],s}function mt(f,t,l){const s=f.slice();return s[7]=t[l],s}function _t(f,t,l){const s=f.slice();return s[12]=t[l],s}function yt(f){let t;return{c(){t=r("p"),t.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",T(t,"class","txt-hint txt-sm txt-right")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function kt(f){let t,l,s,b,u,d,p,k,C,w,O,R,A,j,$,E,B;return{c(){t=r("tr"),t.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',l=m(),s=r("tr"),s.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
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
                    &quot;Manage&quot; access can skip this field.</td>`,w=m(),O=r("tr"),O.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>password</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>New auth record password.</td>`,R=m(),A=r("tr"),A.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>passwordConfirm</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>New auth record password confirmation.</td>`,j=m(),$=r("tr"),$.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>verified</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,E=m(),B=r("tr"),B.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>'},m(c,_){a(c,t,_),a(c,l,_),a(c,s,_),a(c,b,_),a(c,u,_),a(c,d,_),a(c,p,_),a(c,k,_),a(c,C,_),a(c,w,_),a(c,O,_),a(c,R,_),a(c,A,_),a(c,j,_),a(c,$,_),a(c,E,_),a(c,B,_)},d(c){c&&o(t),c&&o(l),c&&o(s),c&&o(b),c&&o(u),c&&o(d),c&&o(p),c&&o(k),c&&o(C),c&&o(w),c&&o(O),c&&o(R),c&&o(A),c&&o(j),c&&o($),c&&o(E),c&&o(B)}}}function Pt(f){let t;return{c(){t=r("span"),t.textContent="Optional",T(t,"class","label label-warning")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function At(f){let t;return{c(){t=r("span"),t.textContent="Required",T(t,"class","label label-success")},m(l,s){a(l,t,s)},d(l){l&&o(t)}}}function Bt(f){var u;let t,l=((u=f[12].options)==null?void 0:u.maxSelect)>1?"ids":"id",s,b;return{c(){t=y("User "),s=y(l),b=y(".")},m(d,p){a(d,t,p),a(d,s,p),a(d,b,p)},p(d,p){var k;p&1&&l!==(l=((k=d[12].options)==null?void 0:k.maxSelect)>1?"ids":"id")&&I(s,l)},d(d){d&&o(t),d&&o(s),d&&o(b)}}}function Ft(f){var u;let t,l=((u=f[12].options)==null?void 0:u.maxSelect)>1?"ids":"id",s,b;return{c(){t=y("Relation record "),s=y(l),b=y(".")},m(d,p){a(d,t,p),a(d,s,p),a(d,b,p)},p(d,p){var k;p&1&&l!==(l=((k=d[12].options)==null?void 0:k.maxSelect)>1?"ids":"id")&&I(s,l)},d(d){d&&o(t),d&&o(s),d&&o(b)}}}function Nt(f){let t,l,s,b,u;return{c(){t=y("File object."),l=r("br"),s=y(`
                        Set to `),b=r("code"),b.textContent="null",u=y(" to delete already uploaded file(s).")},m(d,p){a(d,t,p),a(d,l,p),a(d,s,p),a(d,b,p),a(d,u,p)},p:G,d(d){d&&o(t),d&&o(l),d&&o(s),d&&o(b),d&&o(u)}}}function jt(f){let t;return{c(){t=y("URL address.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function Et(f){let t;return{c(){t=y("Email address.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function Ut(f){let t;return{c(){t=y("JSON array or object.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function It(f){let t;return{c(){t=y("Number value.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function Jt(f){let t;return{c(){t=y("Plain text value.")},m(l,s){a(l,t,s)},p:G,d(l){l&&o(t)}}}function ht(f,t){let l,s,b,u,d,p=t[12].name+"",k,C,w,O,R=U.getFieldValueType(t[12])+"",A,j,$,E;function B(h,P){return h[12].required?At:Pt}let c=B(t),_=c(t);function K(h,P){if(h[12].type==="text")return Jt;if(h[12].type==="number")return It;if(h[12].type==="json")return Ut;if(h[12].type==="email")return Et;if(h[12].type==="url")return jt;if(h[12].type==="file")return Nt;if(h[12].type==="relation")return Ft;if(h[12].type==="user")return Bt}let L=K(t),S=L&&L(t);return{key:f,first:null,c(){l=r("tr"),s=r("td"),b=r("div"),_.c(),u=m(),d=r("span"),k=y(p),C=m(),w=r("td"),O=r("span"),A=y(R),j=m(),$=r("td"),S&&S.c(),E=m(),T(b,"class","inline-flex"),T(O,"class","label"),this.first=l},m(h,P){a(h,l,P),i(l,s),i(s,b),_.m(b,null),i(b,u),i(b,d),i(d,k),i(l,C),i(l,w),i(w,O),i(O,A),i(l,j),i(l,$),S&&S.m($,null),i(l,E)},p(h,P){t=h,c!==(c=B(t))&&(_.d(1),_=c(t),_&&(_.c(),_.m(b,u))),P&1&&p!==(p=t[12].name+"")&&I(k,p),P&1&&R!==(R=U.getFieldValueType(t[12])+"")&&I(A,R),L===(L=K(t))&&S?S.p(t,P):(S&&S.d(1),S=L&&L(t),S&&(S.c(),S.m($,null)))},d(h){h&&o(l),_.d(),S&&S.d()}}}function vt(f,t){let l,s=t[7].code+"",b,u,d,p;function k(){return t[6](t[7])}return{key:f,first:null,c(){l=r("button"),b=y(s),u=m(),T(l,"class","tab-item"),ce(l,"active",t[1]===t[7].code),this.first=l},m(C,w){a(C,l,w),i(l,b),i(l,u),d||(p=Rt(l,"click",k),d=!0)},p(C,w){t=C,w&4&&s!==(s=t[7].code+"")&&I(b,s),w&6&&ce(l,"active",t[1]===t[7].code)},d(C){C&&o(l),d=!1,p()}}}function wt(f,t){let l,s,b,u;return s=new Tt({props:{content:t[7].body}}),{key:f,first:null,c(){l=r("div"),Ae(s.$$.fragment),b=m(),T(l,"class","tab-item"),ce(l,"active",t[1]===t[7].code),this.first=l},m(d,p){a(d,l,p),Be(s,l,null),i(l,b),u=!0},p(d,p){t=d;const k={};p&4&&(k.content=t[7].body),s.$set(k),(!u||p&6)&&ce(l,"active",t[1]===t[7].code)},i(d){u||(fe(s.$$.fragment,d),u=!0)},o(d){pe(s.$$.fragment,d),u=!1},d(d){d&&o(l),Fe(s)}}}function gt(f){var it,at,ot,dt;let t,l,s=f[0].name+"",b,u,d,p,k,C,w,O=f[0].name+"",R,A,j,$,E,B,c,_,K,L,S,h,P,Ne,ae,W,je,ue,oe=f[0].name+"",be,Ee,me,Ue,_e,X,ye,Z,ke,ee,he,J,ve,Ie,g,we,F=[],Je=new Map,Te,te,Ce,V,Se,ge,Oe,x,$e,Ve,Me,xe,M,Qe,Y,ze,Ke,We,qe,Ye,De,Ge,He,Xe,Re,le,Le,Q,se,N=[],Ze=new Map,et,ne,q=[],tt=new Map,z;_=new Lt({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${f[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},f[3],U.dummyCollectionSchemaData(f[0])),null,4)};

const record = await pb.collection('${(it=f[0])==null?void 0:it.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${f[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},f[3],U.dummyCollectionSchemaData(f[0])),null,2)};

final record = await pb.collection('${(at=f[0])==null?void 0:at.name}').update('RECORD_ID', body: body);
    `}});let D=f[5]&&yt(),H=((ot=f[0])==null?void 0:ot.isAuth)&&kt(),de=(dt=f[0])==null?void 0:dt.schema;const lt=e=>e[12].name;for(let e=0;e<de.length;e+=1){let n=_t(f,de,e),v=lt(n);Je.set(v,F[e]=ht(v,n))}Y=new Tt({props:{content:"?expand=relField1,relField2.subRelField21"}});let re=f[2];const st=e=>e[7].code;for(let e=0;e<re.length;e+=1){let n=mt(f,re,e),v=st(n);Ze.set(v,N[e]=vt(v,n))}let ie=f[2];const nt=e=>e[7].code;for(let e=0;e<ie.length;e+=1){let n=bt(f,ie,e),v=nt(n);tt.set(v,q[e]=wt(v,n))}return{c(){t=r("h3"),l=y("Update ("),b=y(s),u=y(")"),d=m(),p=r("div"),k=r("p"),C=y("Update a single "),w=r("strong"),R=y(O),A=y(" record."),j=m(),$=r("p"),$.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,E=m(),B=r("p"),B.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,c=m(),Ae(_.$$.fragment),K=m(),L=r("h6"),L.textContent="API details",S=m(),h=r("div"),P=r("strong"),P.textContent="PATCH",Ne=m(),ae=r("div"),W=r("p"),je=y("/api/collections/"),ue=r("strong"),be=y(oe),Ee=y("/records/"),me=r("strong"),me.textContent=":id",Ue=m(),D&&D.c(),_e=m(),X=r("div"),X.textContent="Path parameters",ye=m(),Z=r("table"),Z.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to update.</td></tr></tbody>`,ke=m(),ee=r("div"),ee.textContent="Body Parameters",he=m(),J=r("table"),ve=r("thead"),ve.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,Ie=m(),g=r("tbody"),H&&H.c(),we=m();for(let e=0;e<F.length;e+=1)F[e].c();Te=m(),te=r("div"),te.textContent="Query parameters",Ce=m(),V=r("table"),Se=r("thead"),Se.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,ge=m(),Oe=r("tbody"),x=r("tr"),$e=r("td"),$e.textContent="expand",Ve=m(),Me=r("td"),Me.innerHTML='<span class="label">String</span>',xe=m(),M=r("td"),Qe=y(`Auto expand relations when returning the updated record. Ex.:
                `),Ae(Y.$$.fragment),ze=y(`
                Supports up to 6-levels depth nested relations expansion. `),Ke=r("br"),We=y(`
                The expanded relations will be appended to the record under the
                `),qe=r("code"),qe.textContent="expand",Ye=y(" property (eg. "),De=r("code"),De.textContent='"expand": {"relField1": {...}, ...}',Ge=y(`). Only
                the relations that the user has permissions to `),He=r("strong"),He.textContent="view",Xe=y(" will be expanded."),Re=m(),le=r("div"),le.textContent="Responses",Le=m(),Q=r("div"),se=r("div");for(let e=0;e<N.length;e+=1)N[e].c();et=m(),ne=r("div");for(let e=0;e<q.length;e+=1)q[e].c();T(t,"class","m-b-sm"),T(p,"class","content txt-lg m-b-sm"),T(L,"class","m-b-xs"),T(P,"class","label label-primary"),T(ae,"class","content"),T(h,"class","alert alert-warning"),T(X,"class","section-title"),T(Z,"class","table-compact table-border m-b-base"),T(ee,"class","section-title"),T(J,"class","table-compact table-border m-b-base"),T(te,"class","section-title"),T(V,"class","table-compact table-border m-b-lg"),T(le,"class","section-title"),T(se,"class","tabs-header compact left"),T(ne,"class","tabs-content"),T(Q,"class","tabs")},m(e,n){a(e,t,n),i(t,l),i(t,b),i(t,u),a(e,d,n),a(e,p,n),i(p,k),i(k,C),i(k,w),i(w,R),i(k,A),i(p,j),i(p,$),i(p,E),i(p,B),a(e,c,n),Be(_,e,n),a(e,K,n),a(e,L,n),a(e,S,n),a(e,h,n),i(h,P),i(h,Ne),i(h,ae),i(ae,W),i(W,je),i(W,ue),i(ue,be),i(W,Ee),i(W,me),i(h,Ue),D&&D.m(h,null),a(e,_e,n),a(e,X,n),a(e,ye,n),a(e,Z,n),a(e,ke,n),a(e,ee,n),a(e,he,n),a(e,J,n),i(J,ve),i(J,Ie),i(J,g),H&&H.m(g,null),i(g,we);for(let v=0;v<F.length;v+=1)F[v]&&F[v].m(g,null);a(e,Te,n),a(e,te,n),a(e,Ce,n),a(e,V,n),i(V,Se),i(V,ge),i(V,Oe),i(Oe,x),i(x,$e),i(x,Ve),i(x,Me),i(x,xe),i(x,M),i(M,Qe),Be(Y,M,null),i(M,ze),i(M,Ke),i(M,We),i(M,qe),i(M,Ye),i(M,De),i(M,Ge),i(M,He),i(M,Xe),a(e,Re,n),a(e,le,n),a(e,Le,n),a(e,Q,n),i(Q,se);for(let v=0;v<N.length;v+=1)N[v]&&N[v].m(se,null);i(Q,et),i(Q,ne);for(let v=0;v<q.length;v+=1)q[v]&&q[v].m(ne,null);z=!0},p(e,[n]){var rt,ft,pt,ct;(!z||n&1)&&s!==(s=e[0].name+"")&&I(b,s),(!z||n&1)&&O!==(O=e[0].name+"")&&I(R,O);const v={};n&25&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${e[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},e[3],U.dummyCollectionSchemaData(e[0])),null,4)};

const record = await pb.collection('${(rt=e[0])==null?void 0:rt.name}').update('RECORD_ID', data);
    `),n&25&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${e[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},e[3],U.dummyCollectionSchemaData(e[0])),null,2)};

final record = await pb.collection('${(ft=e[0])==null?void 0:ft.name}').update('RECORD_ID', body: body);
    `),_.$set(v),(!z||n&1)&&oe!==(oe=e[0].name+"")&&I(be,oe),e[5]?D||(D=yt(),D.c(),D.m(h,null)):D&&(D.d(1),D=null),(pt=e[0])!=null&&pt.isAuth?H||(H=kt(),H.c(),H.m(g,we)):H&&(H.d(1),H=null),n&1&&(de=(ct=e[0])==null?void 0:ct.schema,F=Pe(F,n,lt,1,e,de,Je,g,ut,ht,null,_t)),n&6&&(re=e[2],N=Pe(N,n,st,1,e,re,Ze,se,ut,vt,null,mt)),n&6&&(ie=e[2],$t(),q=Pe(q,n,nt,1,e,ie,tt,ne,Mt,wt,null,bt),qt())},i(e){if(!z){fe(_.$$.fragment,e),fe(Y.$$.fragment,e);for(let n=0;n<ie.length;n+=1)fe(q[n]);z=!0}},o(e){pe(_.$$.fragment,e),pe(Y.$$.fragment,e);for(let n=0;n<q.length;n+=1)pe(q[n]);z=!1},d(e){e&&o(t),e&&o(d),e&&o(p),e&&o(c),Fe(_,e),e&&o(K),e&&o(L),e&&o(S),e&&o(h),D&&D.d(),e&&o(_e),e&&o(X),e&&o(ye),e&&o(Z),e&&o(ke),e&&o(ee),e&&o(he),e&&o(J),H&&H.d();for(let n=0;n<F.length;n+=1)F[n].d();e&&o(Te),e&&o(te),e&&o(Ce),e&&o(V),Fe(Y),e&&o(Re),e&&o(le),e&&o(Le),e&&o(Q);for(let n=0;n<N.length;n+=1)N[n].d();for(let n=0;n<q.length;n+=1)q[n].d()}}}function Vt(f,t,l){let s,b,{collection:u=new Dt}=t,d=200,p=[],k={};const C=w=>l(1,d=w.code);return f.$$set=w=>{"collection"in w&&l(0,u=w.collection)},f.$$.update=()=>{var w,O;f.$$.dirty&1&&l(5,s=(u==null?void 0:u.updateRule)===null),f.$$.dirty&1&&l(2,p=[{code:200,body:JSON.stringify(U.dummyCollectionRecord(u),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${(O=(w=u==null?void 0:u.schema)==null?void 0:w[0])==null?void 0:O.name}": {
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
            `}]),f.$$.dirty&1&&(u.$isAuth?l(3,k={username:"test_username_update",emailVisibility:!1,password:"87654321",passwordConfirm:"87654321",oldPassword:"12345678"}):l(3,k={}))},l(4,b=U.getApiExampleUrl(Ht.baseUrl)),[u,d,p,k,b,s,C]}class zt extends Ct{constructor(t){super(),St(this,t,Vt,gt,Ot,{collection:0})}}export{zt as default};
