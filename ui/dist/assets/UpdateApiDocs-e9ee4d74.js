import{S as $t,i as Ot,s as Mt,C as U,N as X,M as St,e as o,w as _,b,c as me,f as w,g as i,h as s,m as _e,x as E,P as Ne,Q as _t,k as qt,R as Rt,n as Dt,t as re,a as ce,o as d,d as he,U as Ht,p as Pt,r as ye,u as Lt,y as Z}from"./index-a084d9d7.js";import{S as Ft}from"./SdkTabs-ba0ec979.js";import{F as At}from"./FieldsQueryParam-71e01e64.js";function ht(r,e,t){const n=r.slice();return n[7]=e[t],n}function yt(r,e,t){const n=r.slice();return n[7]=e[t],n}function kt(r,e,t){const n=r.slice();return n[12]=e[t],n}function vt(r){let e;return{c(){e=o("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",w(e,"class","txt-hint txt-sm txt-right")},m(t,n){i(t,e,n)},d(t){t&&d(e)}}}function gt(r){let e,t,n,f,u,c,p,h,T,g,S,D,L,N,$,j,F;return{c(){e=o("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',t=b(),n=o("tr"),n.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>username</span></div></td> <td><span class="label">String</span></td> <td>The username of the auth record.</td>',f=b(),u=o("tr"),u.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address.
                    <br/>
                    This field can be updated only by admins or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,c=b(),p=o("tr"),p.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>emailVisibility</span></div></td> <td><span class="label">Boolean</span></td> <td>Whether to show/hide the auth record email when fetching the record data.</td>',h=b(),T=o("tr"),T.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>oldPassword</span></div></td> <td><span class="label">String</span></td> <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Admins and auth records with
                    &quot;Manage&quot; access can skip this field.</td>`,g=b(),S=o("tr"),S.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>New auth record password.</td>',D=b(),L=o("tr"),L.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>New auth record password confirmation.</td>',N=b(),$=o("tr"),$.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,j=b(),F=o("tr"),F.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>'},m(y,m){i(y,e,m),i(y,t,m),i(y,n,m),i(y,f,m),i(y,u,m),i(y,c,m),i(y,p,m),i(y,h,m),i(y,T,m),i(y,g,m),i(y,S,m),i(y,D,m),i(y,L,m),i(y,N,m),i(y,$,m),i(y,j,m),i(y,F,m)},d(y){y&&(d(e),d(t),d(n),d(f),d(u),d(c),d(p),d(h),d(T),d(g),d(S),d(D),d(L),d(N),d($),d(j),d(F))}}}function Bt(r){let e;return{c(){e=o("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,n){i(t,e,n)},d(t){t&&d(e)}}}function Nt(r){let e;return{c(){e=o("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,n){i(t,e,n)},d(t){t&&d(e)}}}function jt(r){var u;let e,t=((u=r[12].options)==null?void 0:u.maxSelect)>1?"ids":"id",n,f;return{c(){e=_("User "),n=_(t),f=_(".")},m(c,p){i(c,e,p),i(c,n,p),i(c,f,p)},p(c,p){var h;p&1&&t!==(t=((h=c[12].options)==null?void 0:h.maxSelect)>1?"ids":"id")&&E(n,t)},d(c){c&&(d(e),d(n),d(f))}}}function Ut(r){var u;let e,t=((u=r[12].options)==null?void 0:u.maxSelect)>1?"ids":"id",n,f;return{c(){e=_("Relation record "),n=_(t),f=_(".")},m(c,p){i(c,e,p),i(c,n,p),i(c,f,p)},p(c,p){var h;p&1&&t!==(t=((h=c[12].options)==null?void 0:h.maxSelect)>1?"ids":"id")&&E(n,t)},d(c){c&&(d(e),d(n),d(f))}}}function Et(r){let e,t,n,f,u;return{c(){e=_("File object."),t=o("br"),n=_(`
                        Set to `),f=o("code"),f.textContent="null",u=_(" to delete already uploaded file(s).")},m(c,p){i(c,e,p),i(c,t,p),i(c,n,p),i(c,f,p),i(c,u,p)},p:Z,d(c){c&&(d(e),d(t),d(n),d(f),d(u))}}}function It(r){let e;return{c(){e=_("URL address.")},m(t,n){i(t,e,n)},p:Z,d(t){t&&d(e)}}}function Jt(r){let e;return{c(){e=_("Email address.")},m(t,n){i(t,e,n)},p:Z,d(t){t&&d(e)}}}function Vt(r){let e;return{c(){e=_("JSON array or object.")},m(t,n){i(t,e,n)},p:Z,d(t){t&&d(e)}}}function xt(r){let e;return{c(){e=_("Number value.")},m(t,n){i(t,e,n)},p:Z,d(t){t&&d(e)}}}function Qt(r){let e;return{c(){e=_("Plain text value.")},m(t,n){i(t,e,n)},p:Z,d(t){t&&d(e)}}}function wt(r,e){let t,n,f,u,c,p=e[12].name+"",h,T,g,S,D=U.getFieldValueType(e[12])+"",L,N,$,j;function F(k,P){return k[12].required?Nt:Bt}let y=F(e),m=y(e);function K(k,P){if(k[12].type==="text")return Qt;if(k[12].type==="number")return xt;if(k[12].type==="json")return Vt;if(k[12].type==="email")return Jt;if(k[12].type==="url")return It;if(k[12].type==="file")return Et;if(k[12].type==="relation")return Ut;if(k[12].type==="user")return jt}let H=K(e),C=H&&H(e);return{key:r,first:null,c(){t=o("tr"),n=o("td"),f=o("div"),m.c(),u=b(),c=o("span"),h=_(p),T=b(),g=o("td"),S=o("span"),L=_(D),N=b(),$=o("td"),C&&C.c(),j=b(),w(f,"class","inline-flex"),w(S,"class","label"),this.first=t},m(k,P){i(k,t,P),s(t,n),s(n,f),m.m(f,null),s(f,u),s(f,c),s(c,h),s(t,T),s(t,g),s(g,S),s(S,L),s(t,N),s(t,$),C&&C.m($,null),s(t,j)},p(k,P){e=k,y!==(y=F(e))&&(m.d(1),m=y(e),m&&(m.c(),m.m(f,u))),P&1&&p!==(p=e[12].name+"")&&E(h,p),P&1&&D!==(D=U.getFieldValueType(e[12])+"")&&E(L,D),H===(H=K(e))&&C?C.p(e,P):(C&&C.d(1),C=H&&H(e),C&&(C.c(),C.m($,null)))},d(k){k&&d(t),m.d(),C&&C.d()}}}function Tt(r,e){let t,n=e[7].code+"",f,u,c,p;function h(){return e[6](e[7])}return{key:r,first:null,c(){t=o("button"),f=_(n),u=b(),w(t,"class","tab-item"),ye(t,"active",e[1]===e[7].code),this.first=t},m(T,g){i(T,t,g),s(t,f),s(t,u),c||(p=Lt(t,"click",h),c=!0)},p(T,g){e=T,g&4&&n!==(n=e[7].code+"")&&E(f,n),g&6&&ye(t,"active",e[1]===e[7].code)},d(T){T&&d(t),c=!1,p()}}}function Ct(r,e){let t,n,f,u;return n=new St({props:{content:e[7].body}}),{key:r,first:null,c(){t=o("div"),me(n.$$.fragment),f=b(),w(t,"class","tab-item"),ye(t,"active",e[1]===e[7].code),this.first=t},m(c,p){i(c,t,p),_e(n,t,null),s(t,f),u=!0},p(c,p){e=c;const h={};p&4&&(h.content=e[7].body),n.$set(h),(!u||p&6)&&ye(t,"active",e[1]===e[7].code)},i(c){u||(re(n.$$.fragment,c),u=!0)},o(c){ce(n.$$.fragment,c),u=!1},d(c){c&&d(t),he(n)}}}function zt(r){var ot,rt,ct,pt;let e,t,n=r[0].name+"",f,u,c,p,h,T,g,S=r[0].name+"",D,L,N,$,j,F,y,m,K,H,C,k,P,je,pe,W,Ue,ke,ue=r[0].name+"",ve,Ee,ge,Ie,we,ee,Te,te,Ce,le,Se,I,$e,Je,J,Oe,A=[],Ve=new Map,Me,ne,qe,V,Re,xe,ae,x,De,Qe,He,ze,O,Ke,Y,We,Ye,Ge,Pe,Xe,Le,Ze,Fe,et,tt,G,Ae,se,Be,Q,ie,B=[],lt=new Map,nt,de,M=[],at=new Map,z;m=new Ft({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${r[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},r[3],U.dummyCollectionSchemaData(r[0])),null,4)};

const record = await pb.collection('${(ot=r[0])==null?void 0:ot.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${r[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},r[3],U.dummyCollectionSchemaData(r[0])),null,2)};

final record = await pb.collection('${(rt=r[0])==null?void 0:rt.name}').update('RECORD_ID', body: body);
    `}});let q=r[5]&&vt(),R=((ct=r[0])==null?void 0:ct.isAuth)&&gt(),fe=X((pt=r[0])==null?void 0:pt.schema);const st=l=>l[12].name;for(let l=0;l<fe.length;l+=1){let a=kt(r,fe,l),v=st(a);Ve.set(v,A[l]=wt(v,a))}Y=new St({props:{content:"?expand=relField1,relField2.subRelField21"}}),G=new At({});let be=X(r[2]);const it=l=>l[7].code;for(let l=0;l<be.length;l+=1){let a=yt(r,be,l),v=it(a);lt.set(v,B[l]=Tt(v,a))}let oe=X(r[2]);const dt=l=>l[7].code;for(let l=0;l<oe.length;l+=1){let a=ht(r,oe,l),v=dt(a);at.set(v,M[l]=Ct(v,a))}return{c(){e=o("h3"),t=_("Update ("),f=_(n),u=_(")"),c=b(),p=o("div"),h=o("p"),T=_("Update a single "),g=o("strong"),D=_(S),L=_(" record."),N=b(),$=o("p"),$.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,j=b(),F=o("p"),F.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,y=b(),me(m.$$.fragment),K=b(),H=o("h6"),H.textContent="API details",C=b(),k=o("div"),P=o("strong"),P.textContent="PATCH",je=b(),pe=o("div"),W=o("p"),Ue=_("/api/collections/"),ke=o("strong"),ve=_(ue),Ee=_("/records/"),ge=o("strong"),ge.textContent=":id",Ie=b(),q&&q.c(),we=b(),ee=o("div"),ee.textContent="Path parameters",Te=b(),te=o("table"),te.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to update.</td></tr></tbody>',Ce=b(),le=o("div"),le.textContent="Body Parameters",Se=b(),I=o("table"),$e=o("thead"),$e.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',Je=b(),J=o("tbody"),R&&R.c(),Oe=b();for(let l=0;l<A.length;l+=1)A[l].c();Me=b(),ne=o("div"),ne.textContent="Query parameters",qe=b(),V=o("table"),Re=o("thead"),Re.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',xe=b(),ae=o("tbody"),x=o("tr"),De=o("td"),De.textContent="expand",Qe=b(),He=o("td"),He.innerHTML='<span class="label">String</span>',ze=b(),O=o("td"),Ke=_(`Auto expand relations when returning the updated record. Ex.:
                `),me(Y.$$.fragment),We=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ye=o("br"),Ge=_(`
                The expanded relations will be appended to the record under the
                `),Pe=o("code"),Pe.textContent="expand",Xe=_(" property (eg. "),Le=o("code"),Le.textContent='"expand": {"relField1": {...}, ...}',Ze=_(`). Only
                the relations that the user has permissions to `),Fe=o("strong"),Fe.textContent="view",et=_(" will be expanded."),tt=b(),me(G.$$.fragment),Ae=b(),se=o("div"),se.textContent="Responses",Be=b(),Q=o("div"),ie=o("div");for(let l=0;l<B.length;l+=1)B[l].c();nt=b(),de=o("div");for(let l=0;l<M.length;l+=1)M[l].c();w(e,"class","m-b-sm"),w(p,"class","content txt-lg m-b-sm"),w(H,"class","m-b-xs"),w(P,"class","label label-primary"),w(pe,"class","content"),w(k,"class","alert alert-warning"),w(ee,"class","section-title"),w(te,"class","table-compact table-border m-b-base"),w(le,"class","section-title"),w(I,"class","table-compact table-border m-b-base"),w(ne,"class","section-title"),w(V,"class","table-compact table-border m-b-lg"),w(se,"class","section-title"),w(ie,"class","tabs-header compact left"),w(de,"class","tabs-content"),w(Q,"class","tabs")},m(l,a){i(l,e,a),s(e,t),s(e,f),s(e,u),i(l,c,a),i(l,p,a),s(p,h),s(h,T),s(h,g),s(g,D),s(h,L),s(p,N),s(p,$),s(p,j),s(p,F),i(l,y,a),_e(m,l,a),i(l,K,a),i(l,H,a),i(l,C,a),i(l,k,a),s(k,P),s(k,je),s(k,pe),s(pe,W),s(W,Ue),s(W,ke),s(ke,ve),s(W,Ee),s(W,ge),s(k,Ie),q&&q.m(k,null),i(l,we,a),i(l,ee,a),i(l,Te,a),i(l,te,a),i(l,Ce,a),i(l,le,a),i(l,Se,a),i(l,I,a),s(I,$e),s(I,Je),s(I,J),R&&R.m(J,null),s(J,Oe);for(let v=0;v<A.length;v+=1)A[v]&&A[v].m(J,null);i(l,Me,a),i(l,ne,a),i(l,qe,a),i(l,V,a),s(V,Re),s(V,xe),s(V,ae),s(ae,x),s(x,De),s(x,Qe),s(x,He),s(x,ze),s(x,O),s(O,Ke),_e(Y,O,null),s(O,We),s(O,Ye),s(O,Ge),s(O,Pe),s(O,Xe),s(O,Le),s(O,Ze),s(O,Fe),s(O,et),s(ae,tt),_e(G,ae,null),i(l,Ae,a),i(l,se,a),i(l,Be,a),i(l,Q,a),s(Q,ie);for(let v=0;v<B.length;v+=1)B[v]&&B[v].m(ie,null);s(Q,nt),s(Q,de);for(let v=0;v<M.length;v+=1)M[v]&&M[v].m(de,null);z=!0},p(l,[a]){var ut,ft,bt,mt;(!z||a&1)&&n!==(n=l[0].name+"")&&E(f,n),(!z||a&1)&&S!==(S=l[0].name+"")&&E(D,S);const v={};a&25&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},l[3],U.dummyCollectionSchemaData(l[0])),null,4)};

const record = await pb.collection('${(ut=l[0])==null?void 0:ut.name}').update('RECORD_ID', data);
    `),a&25&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},l[3],U.dummyCollectionSchemaData(l[0])),null,2)};

final record = await pb.collection('${(ft=l[0])==null?void 0:ft.name}').update('RECORD_ID', body: body);
    `),m.$set(v),(!z||a&1)&&ue!==(ue=l[0].name+"")&&E(ve,ue),l[5]?q||(q=vt(),q.c(),q.m(k,null)):q&&(q.d(1),q=null),(bt=l[0])!=null&&bt.isAuth?R||(R=gt(),R.c(),R.m(J,Oe)):R&&(R.d(1),R=null),a&1&&(fe=X((mt=l[0])==null?void 0:mt.schema),A=Ne(A,a,st,1,l,fe,Ve,J,_t,wt,null,kt)),a&6&&(be=X(l[2]),B=Ne(B,a,it,1,l,be,lt,ie,_t,Tt,null,yt)),a&6&&(oe=X(l[2]),qt(),M=Ne(M,a,dt,1,l,oe,at,de,Rt,Ct,null,ht),Dt())},i(l){if(!z){re(m.$$.fragment,l),re(Y.$$.fragment,l),re(G.$$.fragment,l);for(let a=0;a<oe.length;a+=1)re(M[a]);z=!0}},o(l){ce(m.$$.fragment,l),ce(Y.$$.fragment,l),ce(G.$$.fragment,l);for(let a=0;a<M.length;a+=1)ce(M[a]);z=!1},d(l){l&&(d(e),d(c),d(p),d(y),d(K),d(H),d(C),d(k),d(we),d(ee),d(Te),d(te),d(Ce),d(le),d(Se),d(I),d(Me),d(ne),d(qe),d(V),d(Ae),d(se),d(Be),d(Q)),he(m,l),q&&q.d(),R&&R.d();for(let a=0;a<A.length;a+=1)A[a].d();he(Y),he(G);for(let a=0;a<B.length;a+=1)B[a].d();for(let a=0;a<M.length;a+=1)M[a].d()}}}function Kt(r,e,t){let n,f,{collection:u=new Ht}=e,c=200,p=[],h={};const T=g=>t(1,c=g.code);return r.$$set=g=>{"collection"in g&&t(0,u=g.collection)},r.$$.update=()=>{var g,S;r.$$.dirty&1&&t(5,n=(u==null?void 0:u.updateRule)===null),r.$$.dirty&1&&t(2,p=[{code:200,body:JSON.stringify(U.dummyCollectionRecord(u),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${(S=(g=u==null?void 0:u.schema)==null?void 0:g[0])==null?void 0:S.name}": {
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
            `}]),r.$$.dirty&1&&(u.$isAuth?t(3,h={username:"test_username_update",emailVisibility:!1,password:"87654321",passwordConfirm:"87654321",oldPassword:"12345678"}):t(3,h={}))},t(4,f=U.getApiExampleUrl(Pt.baseUrl)),[u,c,p,h,f,n,T]}class Xt extends $t{constructor(e){super(),Ot(this,e,Kt,zt,Mt,{collection:0})}}export{Xt as default};
