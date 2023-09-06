import{S as Ct,i as St,s as Ot,C as E,O as X,N as Tt,e as r,w as _,b as f,c as me,f as g,g as d,h as s,m as _e,x as U,P as je,Q as bt,k as $t,R as Mt,n as qt,t as re,a as ce,o,d as he,p as Rt,r as ye,u as Dt,y as Z}from"./index-1320cae2.js";import{S as Ht}from"./SdkTabs-da4c9c1b.js";import{F as Pt}from"./FieldsQueryParam-534c6ac5.js";function mt(c,e,t){const n=c.slice();return n[8]=e[t],n}function _t(c,e,t){const n=c.slice();return n[8]=e[t],n}function ht(c,e,t){const n=c.slice();return n[13]=e[t],n}function yt(c){let e;return{c(){e=r("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",g(e,"class","txt-hint txt-sm txt-right")},m(t,n){d(t,e,n)},d(t){t&&o(e)}}}function kt(c){let e,t,n,u,m,i,p,y,T,C,w,O,L,j,$,A,F;return{c(){e=r("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',t=f(),n=r("tr"),n.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>username</span></div></td> <td><span class="label">String</span></td> <td>The username of the auth record.</td>',u=f(),m=r("tr"),m.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address.
                    <br/>
                    This field can be updated only by admins or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,i=f(),p=r("tr"),p.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>emailVisibility</span></div></td> <td><span class="label">Boolean</span></td> <td>Whether to show/hide the auth record email when fetching the record data.</td>',y=f(),T=r("tr"),T.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>oldPassword</span></div></td> <td><span class="label">String</span></td> <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Admins and auth records with
                    &quot;Manage&quot; access can skip this field.</td>`,C=f(),w=r("tr"),w.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>New auth record password.</td>',O=f(),L=r("tr"),L.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>New auth record password confirmation.</td>',j=f(),$=r("tr"),$.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,A=f(),F=r("tr"),F.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>'},m(h,b){d(h,e,b),d(h,t,b),d(h,n,b),d(h,u,b),d(h,m,b),d(h,i,b),d(h,p,b),d(h,y,b),d(h,T,b),d(h,C,b),d(h,w,b),d(h,O,b),d(h,L,b),d(h,j,b),d(h,$,b),d(h,A,b),d(h,F,b)},d(h){h&&(o(e),o(t),o(n),o(u),o(m),o(i),o(p),o(y),o(T),o(C),o(w),o(O),o(L),o(j),o($),o(A),o(F))}}}function Lt(c){let e;return{c(){e=r("span"),e.textContent="Optional",g(e,"class","label label-warning")},m(t,n){d(t,e,n)},d(t){t&&o(e)}}}function Ft(c){let e;return{c(){e=r("span"),e.textContent="Required",g(e,"class","label label-success")},m(t,n){d(t,e,n)},d(t){t&&o(e)}}}function Bt(c){var m;let e,t=((m=c[13].options)==null?void 0:m.maxSelect)>1?"ids":"id",n,u;return{c(){e=_("User "),n=_(t),u=_(".")},m(i,p){d(i,e,p),d(i,n,p),d(i,u,p)},p(i,p){var y;p&1&&t!==(t=((y=i[13].options)==null?void 0:y.maxSelect)>1?"ids":"id")&&U(n,t)},d(i){i&&(o(e),o(n),o(u))}}}function Nt(c){var m;let e,t=((m=c[13].options)==null?void 0:m.maxSelect)>1?"ids":"id",n,u;return{c(){e=_("Relation record "),n=_(t),u=_(".")},m(i,p){d(i,e,p),d(i,n,p),d(i,u,p)},p(i,p){var y;p&1&&t!==(t=((y=i[13].options)==null?void 0:y.maxSelect)>1?"ids":"id")&&U(n,t)},d(i){i&&(o(e),o(n),o(u))}}}function jt(c){let e,t,n,u,m;return{c(){e=_("File object."),t=r("br"),n=_(`
                        Set to `),u=r("code"),u.textContent="null",m=_(" to delete already uploaded file(s).")},m(i,p){d(i,e,p),d(i,t,p),d(i,n,p),d(i,u,p),d(i,m,p)},p:Z,d(i){i&&(o(e),o(t),o(n),o(u),o(m))}}}function At(c){let e;return{c(){e=_("URL address.")},m(t,n){d(t,e,n)},p:Z,d(t){t&&o(e)}}}function Et(c){let e;return{c(){e=_("Email address.")},m(t,n){d(t,e,n)},p:Z,d(t){t&&o(e)}}}function Ut(c){let e;return{c(){e=_("JSON array or object.")},m(t,n){d(t,e,n)},p:Z,d(t){t&&o(e)}}}function It(c){let e;return{c(){e=_("Number value.")},m(t,n){d(t,e,n)},p:Z,d(t){t&&o(e)}}}function Jt(c){let e;return{c(){e=_("Plain text value.")},m(t,n){d(t,e,n)},p:Z,d(t){t&&o(e)}}}function vt(c,e){let t,n,u,m,i,p=e[13].name+"",y,T,C,w,O=E.getFieldValueType(e[13])+"",L,j,$,A;function F(k,P){return k[13].required?Ft:Lt}let h=F(e),b=h(e);function K(k,P){if(k[13].type==="text")return Jt;if(k[13].type==="number")return It;if(k[13].type==="json")return Ut;if(k[13].type==="email")return Et;if(k[13].type==="url")return At;if(k[13].type==="file")return jt;if(k[13].type==="relation")return Nt;if(k[13].type==="user")return Bt}let H=K(e),S=H&&H(e);return{key:c,first:null,c(){t=r("tr"),n=r("td"),u=r("div"),b.c(),m=f(),i=r("span"),y=_(p),T=f(),C=r("td"),w=r("span"),L=_(O),j=f(),$=r("td"),S&&S.c(),A=f(),g(u,"class","inline-flex"),g(w,"class","label"),this.first=t},m(k,P){d(k,t,P),s(t,n),s(n,u),b.m(u,null),s(u,m),s(u,i),s(i,y),s(t,T),s(t,C),s(C,w),s(w,L),s(t,j),s(t,$),S&&S.m($,null),s(t,A)},p(k,P){e=k,h!==(h=F(e))&&(b.d(1),b=h(e),b&&(b.c(),b.m(u,m))),P&1&&p!==(p=e[13].name+"")&&U(y,p),P&1&&O!==(O=E.getFieldValueType(e[13])+"")&&U(L,O),H===(H=K(e))&&S?S.p(e,P):(S&&S.d(1),S=H&&H(e),S&&(S.c(),S.m($,null)))},d(k){k&&o(t),b.d(),S&&S.d()}}}function gt(c,e){let t,n=e[8].code+"",u,m,i,p;function y(){return e[7](e[8])}return{key:c,first:null,c(){t=r("button"),u=_(n),m=f(),g(t,"class","tab-item"),ye(t,"active",e[1]===e[8].code),this.first=t},m(T,C){d(T,t,C),s(t,u),s(t,m),i||(p=Dt(t,"click",y),i=!0)},p(T,C){e=T,C&4&&n!==(n=e[8].code+"")&&U(u,n),C&6&&ye(t,"active",e[1]===e[8].code)},d(T){T&&o(t),i=!1,p()}}}function wt(c,e){let t,n,u,m;return n=new Tt({props:{content:e[8].body}}),{key:c,first:null,c(){t=r("div"),me(n.$$.fragment),u=f(),g(t,"class","tab-item"),ye(t,"active",e[1]===e[8].code),this.first=t},m(i,p){d(i,t,p),_e(n,t,null),s(t,u),m=!0},p(i,p){e=i;const y={};p&4&&(y.content=e[8].body),n.$set(y),(!m||p&6)&&ye(t,"active",e[1]===e[8].code)},i(i){m||(re(n.$$.fragment,i),m=!0)},o(i){ce(n.$$.fragment,i),m=!1},d(i){i&&o(t),he(n)}}}function Vt(c){var ot,rt,ct;let e,t,n=c[0].name+"",u,m,i,p,y,T,C,w=c[0].name+"",O,L,j,$,A,F,h,b,K,H,S,k,P,Ae,pe,W,Ee,ke,ue=c[0].name+"",ve,Ue,ge,Ie,we,ee,Te,te,Ce,le,Se,I,Oe,Je,J,$e,B=[],Ve=new Map,Me,ne,qe,V,Re,xe,ae,x,De,Qe,He,ze,M,Ke,Y,We,Ye,Ge,Pe,Xe,Le,Ze,Fe,et,tt,G,Be,se,Ne,Q,ie,N=[],lt=new Map,nt,de,q=[],at=new Map,z;b=new Ht({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${c[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},c[3],E.dummyCollectionSchemaData(c[0])),null,4)};

const record = await pb.collection('${(ot=c[0])==null?void 0:ot.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${c[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},c[3],E.dummyCollectionSchemaData(c[0])),null,2)};

final record = await pb.collection('${(rt=c[0])==null?void 0:rt.name}').update('RECORD_ID', body: body);
    `}});let R=c[5]&&yt(),D=c[6]&&kt(),fe=X((ct=c[0])==null?void 0:ct.schema);const st=l=>l[13].name;for(let l=0;l<fe.length;l+=1){let a=ht(c,fe,l),v=st(a);Ve.set(v,B[l]=vt(v,a))}Y=new Tt({props:{content:"?expand=relField1,relField2.subRelField21"}}),G=new Pt({});let be=X(c[2]);const it=l=>l[8].code;for(let l=0;l<be.length;l+=1){let a=_t(c,be,l),v=it(a);lt.set(v,N[l]=gt(v,a))}let oe=X(c[2]);const dt=l=>l[8].code;for(let l=0;l<oe.length;l+=1){let a=mt(c,oe,l),v=dt(a);at.set(v,q[l]=wt(v,a))}return{c(){e=r("h3"),t=_("Update ("),u=_(n),m=_(")"),i=f(),p=r("div"),y=r("p"),T=_("Update a single "),C=r("strong"),O=_(w),L=_(" record."),j=f(),$=r("p"),$.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,A=f(),F=r("p"),F.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,h=f(),me(b.$$.fragment),K=f(),H=r("h6"),H.textContent="API details",S=f(),k=r("div"),P=r("strong"),P.textContent="PATCH",Ae=f(),pe=r("div"),W=r("p"),Ee=_("/api/collections/"),ke=r("strong"),ve=_(ue),Ue=_("/records/"),ge=r("strong"),ge.textContent=":id",Ie=f(),R&&R.c(),we=f(),ee=r("div"),ee.textContent="Path parameters",Te=f(),te=r("table"),te.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to update.</td></tr></tbody>',Ce=f(),le=r("div"),le.textContent="Body Parameters",Se=f(),I=r("table"),Oe=r("thead"),Oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',Je=f(),J=r("tbody"),D&&D.c(),$e=f();for(let l=0;l<B.length;l+=1)B[l].c();Me=f(),ne=r("div"),ne.textContent="Query parameters",qe=f(),V=r("table"),Re=r("thead"),Re.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',xe=f(),ae=r("tbody"),x=r("tr"),De=r("td"),De.textContent="expand",Qe=f(),He=r("td"),He.innerHTML='<span class="label">String</span>',ze=f(),M=r("td"),Ke=_(`Auto expand relations when returning the updated record. Ex.:
                `),me(Y.$$.fragment),We=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ye=r("br"),Ge=_(`
                The expanded relations will be appended to the record under the
                `),Pe=r("code"),Pe.textContent="expand",Xe=_(" property (eg. "),Le=r("code"),Le.textContent='"expand": {"relField1": {...}, ...}',Ze=_(`). Only
                the relations that the user has permissions to `),Fe=r("strong"),Fe.textContent="view",et=_(" will be expanded."),tt=f(),me(G.$$.fragment),Be=f(),se=r("div"),se.textContent="Responses",Ne=f(),Q=r("div"),ie=r("div");for(let l=0;l<N.length;l+=1)N[l].c();nt=f(),de=r("div");for(let l=0;l<q.length;l+=1)q[l].c();g(e,"class","m-b-sm"),g(p,"class","content txt-lg m-b-sm"),g(H,"class","m-b-xs"),g(P,"class","label label-primary"),g(pe,"class","content"),g(k,"class","alert alert-warning"),g(ee,"class","section-title"),g(te,"class","table-compact table-border m-b-base"),g(le,"class","section-title"),g(I,"class","table-compact table-border m-b-base"),g(ne,"class","section-title"),g(V,"class","table-compact table-border m-b-lg"),g(se,"class","section-title"),g(ie,"class","tabs-header compact combined left"),g(de,"class","tabs-content"),g(Q,"class","tabs")},m(l,a){d(l,e,a),s(e,t),s(e,u),s(e,m),d(l,i,a),d(l,p,a),s(p,y),s(y,T),s(y,C),s(C,O),s(y,L),s(p,j),s(p,$),s(p,A),s(p,F),d(l,h,a),_e(b,l,a),d(l,K,a),d(l,H,a),d(l,S,a),d(l,k,a),s(k,P),s(k,Ae),s(k,pe),s(pe,W),s(W,Ee),s(W,ke),s(ke,ve),s(W,Ue),s(W,ge),s(k,Ie),R&&R.m(k,null),d(l,we,a),d(l,ee,a),d(l,Te,a),d(l,te,a),d(l,Ce,a),d(l,le,a),d(l,Se,a),d(l,I,a),s(I,Oe),s(I,Je),s(I,J),D&&D.m(J,null),s(J,$e);for(let v=0;v<B.length;v+=1)B[v]&&B[v].m(J,null);d(l,Me,a),d(l,ne,a),d(l,qe,a),d(l,V,a),s(V,Re),s(V,xe),s(V,ae),s(ae,x),s(x,De),s(x,Qe),s(x,He),s(x,ze),s(x,M),s(M,Ke),_e(Y,M,null),s(M,We),s(M,Ye),s(M,Ge),s(M,Pe),s(M,Xe),s(M,Le),s(M,Ze),s(M,Fe),s(M,et),s(ae,tt),_e(G,ae,null),d(l,Be,a),d(l,se,a),d(l,Ne,a),d(l,Q,a),s(Q,ie);for(let v=0;v<N.length;v+=1)N[v]&&N[v].m(ie,null);s(Q,nt),s(Q,de);for(let v=0;v<q.length;v+=1)q[v]&&q[v].m(de,null);z=!0},p(l,[a]){var pt,ut,ft;(!z||a&1)&&n!==(n=l[0].name+"")&&U(u,n),(!z||a&1)&&w!==(w=l[0].name+"")&&U(O,w);const v={};a&25&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[4]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},l[3],E.dummyCollectionSchemaData(l[0])),null,4)};

const record = await pb.collection('${(pt=l[0])==null?void 0:pt.name}').update('RECORD_ID', data);
    `),a&25&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},l[3],E.dummyCollectionSchemaData(l[0])),null,2)};

final record = await pb.collection('${(ut=l[0])==null?void 0:ut.name}').update('RECORD_ID', body: body);
    `),b.$set(v),(!z||a&1)&&ue!==(ue=l[0].name+"")&&U(ve,ue),l[5]?R||(R=yt(),R.c(),R.m(k,null)):R&&(R.d(1),R=null),l[6]?D||(D=kt(),D.c(),D.m(J,$e)):D&&(D.d(1),D=null),a&1&&(fe=X((ft=l[0])==null?void 0:ft.schema),B=je(B,a,st,1,l,fe,Ve,J,bt,vt,null,ht)),a&6&&(be=X(l[2]),N=je(N,a,it,1,l,be,lt,ie,bt,gt,null,_t)),a&6&&(oe=X(l[2]),$t(),q=je(q,a,dt,1,l,oe,at,de,Mt,wt,null,mt),qt())},i(l){if(!z){re(b.$$.fragment,l),re(Y.$$.fragment,l),re(G.$$.fragment,l);for(let a=0;a<oe.length;a+=1)re(q[a]);z=!0}},o(l){ce(b.$$.fragment,l),ce(Y.$$.fragment,l),ce(G.$$.fragment,l);for(let a=0;a<q.length;a+=1)ce(q[a]);z=!1},d(l){l&&(o(e),o(i),o(p),o(h),o(K),o(H),o(S),o(k),o(we),o(ee),o(Te),o(te),o(Ce),o(le),o(Se),o(I),o(Me),o(ne),o(qe),o(V),o(Be),o(se),o(Ne),o(Q)),he(b,l),R&&R.d(),D&&D.d();for(let a=0;a<B.length;a+=1)B[a].d();he(Y),he(G);for(let a=0;a<N.length;a+=1)N[a].d();for(let a=0;a<q.length;a+=1)q[a].d()}}}function xt(c,e,t){let n,u,m,{collection:i}=e,p=200,y=[],T={};const C=w=>t(1,p=w.code);return c.$$set=w=>{"collection"in w&&t(0,i=w.collection)},c.$$.update=()=>{var w,O;c.$$.dirty&1&&t(6,n=(i==null?void 0:i.type)==="auth"),c.$$.dirty&1&&t(5,u=(i==null?void 0:i.updateRule)===null),c.$$.dirty&1&&t(2,y=[{code:200,body:JSON.stringify(E.dummyCollectionRecord(i),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${(O=(w=i==null?void 0:i.schema)==null?void 0:w[0])==null?void 0:O.name}": {
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
            `}]),c.$$.dirty&1&&(i.type==="auth"?t(3,T={username:"test_username_update",emailVisibility:!1,password:"87654321",passwordConfirm:"87654321",oldPassword:"12345678"}):t(3,T={}))},t(4,m=E.getApiExampleUrl(Rt.baseUrl)),[i,p,y,T,m,u,n,C]}class Wt extends Ct{constructor(e){super(),St(this,e,xt,Vt,Ot,{collection:0})}}export{Wt as default};
