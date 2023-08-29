import{S as qt,i as Ot,s as Mt,C as Q,N as Tt,e as a,w as k,b as u,c as be,f as h,g as d,h as n,m as _e,x,O as Be,P as _t,k as Ht,Q as Lt,n as Pt,t as fe,a as pe,o as c,d as ke,p as gt,r as ye,u as Ft,y as ne}from"./index-d577ac2a.js";import{S as At}from"./SdkTabs-7ce69c38.js";import{F as Bt}from"./FieldsQueryParam-23e11274.js";function kt(o,e,l){const s=o.slice();return s[8]=e[l],s}function yt(o,e,l){const s=o.slice();return s[8]=e[l],s}function vt(o,e,l){const s=o.slice();return s[13]=e[l],s}function ht(o){let e;return{c(){e=a("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",h(e,"class","txt-hint txt-sm txt-right")},m(l,s){d(l,e,s)},d(l){l&&c(e)}}}function wt(o){let e,l,s,m,b,r,f,v,T,q,$,L,D,E,g,I,j,R,C,N,O,w,_;function M(p,S){var ee,K;return(K=(ee=p[0])==null?void 0:ee.options)!=null&&K.requireEmail?jt:Rt}let z=M(o),F=z(o);return{c(){e=a("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',l=u(),s=a("tr"),s.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>username</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The username of the auth record.
                    <br/>
                    If not set, it will be auto generated.</td>`,m=u(),b=a("tr"),r=a("td"),f=a("div"),F.c(),v=u(),T=a("span"),T.textContent="email",q=u(),$=a("td"),$.innerHTML='<span class="label">String</span>',L=u(),D=a("td"),D.textContent="Auth record email address.",E=u(),g=a("tr"),g.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>emailVisibility</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Whether to show/hide the auth record email when fetching the record data.</td>`,I=u(),j=a("tr"),j.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                        <span>password</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Auth record password.</td>`,R=u(),C=a("tr"),C.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                        <span>passwordConfirm</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Auth record password confirmation.</td>`,N=u(),O=a("tr"),O.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>verified</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,w=u(),_=a("tr"),_.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>',h(f,"class","inline-flex")},m(p,S){d(p,e,S),d(p,l,S),d(p,s,S),d(p,m,S),d(p,b,S),n(b,r),n(r,f),F.m(f,null),n(f,v),n(f,T),n(b,q),n(b,$),n(b,L),n(b,D),d(p,E,S),d(p,g,S),d(p,I,S),d(p,j,S),d(p,R,S),d(p,C,S),d(p,N,S),d(p,O,S),d(p,w,S),d(p,_,S)},p(p,S){z!==(z=M(p))&&(F.d(1),F=z(p),F&&(F.c(),F.m(f,v)))},d(p){p&&c(e),p&&c(l),p&&c(s),p&&c(m),p&&c(b),F.d(),p&&c(E),p&&c(g),p&&c(I),p&&c(j),p&&c(R),p&&c(C),p&&c(N),p&&c(O),p&&c(w),p&&c(_)}}}function Rt(o){let e;return{c(){e=a("span"),e.textContent="Optional",h(e,"class","label label-warning")},m(l,s){d(l,e,s)},d(l){l&&c(e)}}}function jt(o){let e;return{c(){e=a("span"),e.textContent="Required",h(e,"class","label label-success")},m(l,s){d(l,e,s)},d(l){l&&c(e)}}}function Dt(o){let e;return{c(){e=a("span"),e.textContent="Optional",h(e,"class","label label-warning")},m(l,s){d(l,e,s)},d(l){l&&c(e)}}}function Nt(o){let e;return{c(){e=a("span"),e.textContent="Required",h(e,"class","label label-success")},m(l,s){d(l,e,s)},d(l){l&&c(e)}}}function Vt(o){var b;let e,l=((b=o[13].options)==null?void 0:b.maxSelect)===1?"id":"ids",s,m;return{c(){e=k("Relation record "),s=k(l),m=k(".")},m(r,f){d(r,e,f),d(r,s,f),d(r,m,f)},p(r,f){var v;f&1&&l!==(l=((v=r[13].options)==null?void 0:v.maxSelect)===1?"id":"ids")&&x(s,l)},d(r){r&&c(e),r&&c(s),r&&c(m)}}}function Jt(o){let e,l,s,m,b;return{c(){e=k("File object."),l=a("br"),s=k(`
                        Set to `),m=a("code"),m.textContent="null",b=k(" to delete already uploaded file(s).")},m(r,f){d(r,e,f),d(r,l,f),d(r,s,f),d(r,m,f),d(r,b,f)},p:ne,d(r){r&&c(e),r&&c(l),r&&c(s),r&&c(m),r&&c(b)}}}function Et(o){let e;return{c(){e=k("URL address.")},m(l,s){d(l,e,s)},p:ne,d(l){l&&c(e)}}}function It(o){let e;return{c(){e=k("Email address.")},m(l,s){d(l,e,s)},p:ne,d(l){l&&c(e)}}}function Ut(o){let e;return{c(){e=k("JSON array or object.")},m(l,s){d(l,e,s)},p:ne,d(l){l&&c(e)}}}function Qt(o){let e;return{c(){e=k("Number value.")},m(l,s){d(l,e,s)},p:ne,d(l){l&&c(e)}}}function zt(o){let e;return{c(){e=k("Plain text value.")},m(l,s){d(l,e,s)},p:ne,d(l){l&&c(e)}}}function $t(o,e){let l,s,m,b,r,f=e[13].name+"",v,T,q,$,L=Q.getFieldValueType(e[13])+"",D,E,g,I;function j(_,M){return _[13].required?Nt:Dt}let R=j(e),C=R(e);function N(_,M){if(_[13].type==="text")return zt;if(_[13].type==="number")return Qt;if(_[13].type==="json")return Ut;if(_[13].type==="email")return It;if(_[13].type==="url")return Et;if(_[13].type==="file")return Jt;if(_[13].type==="relation")return Vt}let O=N(e),w=O&&O(e);return{key:o,first:null,c(){l=a("tr"),s=a("td"),m=a("div"),C.c(),b=u(),r=a("span"),v=k(f),T=u(),q=a("td"),$=a("span"),D=k(L),E=u(),g=a("td"),w&&w.c(),I=u(),h(m,"class","inline-flex"),h($,"class","label"),this.first=l},m(_,M){d(_,l,M),n(l,s),n(s,m),C.m(m,null),n(m,b),n(m,r),n(r,v),n(l,T),n(l,q),n(q,$),n($,D),n(l,E),n(l,g),w&&w.m(g,null),n(l,I)},p(_,M){e=_,R!==(R=j(e))&&(C.d(1),C=R(e),C&&(C.c(),C.m(m,b))),M&1&&f!==(f=e[13].name+"")&&x(v,f),M&1&&L!==(L=Q.getFieldValueType(e[13])+"")&&x(D,L),O===(O=N(e))&&w?w.p(e,M):(w&&w.d(1),w=O&&O(e),w&&(w.c(),w.m(g,null)))},d(_){_&&c(l),C.d(),w&&w.d()}}}function Ct(o,e){let l,s=e[8].code+"",m,b,r,f;function v(){return e[7](e[8])}return{key:o,first:null,c(){l=a("button"),m=k(s),b=u(),h(l,"class","tab-item"),ye(l,"active",e[2]===e[8].code),this.first=l},m(T,q){d(T,l,q),n(l,m),n(l,b),r||(f=Ft(l,"click",v),r=!0)},p(T,q){e=T,q&8&&s!==(s=e[8].code+"")&&x(m,s),q&12&&ye(l,"active",e[2]===e[8].code)},d(T){T&&c(l),r=!1,f()}}}function St(o,e){let l,s,m,b;return s=new Tt({props:{content:e[8].body}}),{key:o,first:null,c(){l=a("div"),be(s.$$.fragment),m=u(),h(l,"class","tab-item"),ye(l,"active",e[2]===e[8].code),this.first=l},m(r,f){d(r,l,f),_e(s,l,null),n(l,m),b=!0},p(r,f){e=r;const v={};f&8&&(v.content=e[8].body),s.$set(v),(!b||f&12)&&ye(l,"active",e[2]===e[8].code)},i(r){b||(fe(s.$$.fragment,r),b=!0)},o(r){pe(s.$$.fragment,r),b=!1},d(r){r&&c(l),ke(s)}}}function Kt(o){var at,ot,rt,dt,ct;let e,l,s=o[0].name+"",m,b,r,f,v,T,q,$=o[0].name+"",L,D,E,g,I,j,R,C,N,O,w,_,M,z,F,p,S,ee,K=o[0].name+"",ve,Re,je,he,se,we,W,$e,De,U,Ce,Ne,Se,V=[],Ve=new Map,Te,ie,qe,Y,Oe,Je,ae,G,Me,Ee,He,Ie,H,Ue,te,Qe,ze,Ke,Le,We,Pe,Ye,Ge,Xe,ge,Ze,xe,le,Fe,oe,Ae,X,re,J=[],et=new Map,tt,de,A=[],lt=new Map,Z;C=new At({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${o[5]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},o[4],Q.dummyCollectionSchemaData(o[0])),null,4)};

const record = await pb.collection('${(at=o[0])==null?void 0:at.name}').create(data);
`+(o[1]?`
// (optional) send an email verification request
await pb.collection('${(ot=o[0])==null?void 0:ot.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${o[5]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},o[4],Q.dummyCollectionSchemaData(o[0])),null,2)};

final record = await pb.collection('${(rt=o[0])==null?void 0:rt.name}').create(body: body);
`+(o[1]?`
// (optional) send an email verification request
await pb.collection('${(dt=o[0])==null?void 0:dt.name}').requestVerification('test@example.com');
`:"")}});let B=o[6]&&ht(),P=o[1]&&wt(o),ue=(ct=o[0])==null?void 0:ct.schema;const nt=t=>t[13].name;for(let t=0;t<ue.length;t+=1){let i=vt(o,ue,t),y=nt(i);Ve.set(y,V[t]=$t(y,i))}te=new Tt({props:{content:"?expand=relField1,relField2.subRelField"}}),le=new Bt({});let me=o[3];const st=t=>t[8].code;for(let t=0;t<me.length;t+=1){let i=yt(o,me,t),y=st(i);et.set(y,J[t]=Ct(y,i))}let ce=o[3];const it=t=>t[8].code;for(let t=0;t<ce.length;t+=1){let i=kt(o,ce,t),y=it(i);lt.set(y,A[t]=St(y,i))}return{c(){e=a("h3"),l=k("Create ("),m=k(s),b=k(")"),r=u(),f=a("div"),v=a("p"),T=k("Create a new "),q=a("strong"),L=k($),D=k(" record."),E=u(),g=a("p"),g.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,I=u(),j=a("p"),j.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,R=u(),be(C.$$.fragment),N=u(),O=a("h6"),O.textContent="API details",w=u(),_=a("div"),M=a("strong"),M.textContent="POST",z=u(),F=a("div"),p=a("p"),S=k("/api/collections/"),ee=a("strong"),ve=k(K),Re=k("/records"),je=u(),B&&B.c(),he=u(),se=a("div"),se.textContent="Body Parameters",we=u(),W=a("table"),$e=a("thead"),$e.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,De=u(),U=a("tbody"),Ce=a("tr"),Ce.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                    <span>id</span></div></td> 
            <td><span class="label">String</span></td> 
            <td><strong>15 characters string</strong> to store as record ID.
                <br/>
                If not set, it will be auto generated.</td>`,Ne=u(),P&&P.c(),Se=u();for(let t=0;t<V.length;t+=1)V[t].c();Te=u(),ie=a("div"),ie.textContent="Query parameters",qe=u(),Y=a("table"),Oe=a("thead"),Oe.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Je=u(),ae=a("tbody"),G=a("tr"),Me=a("td"),Me.textContent="expand",Ee=u(),He=a("td"),He.innerHTML='<span class="label">String</span>',Ie=u(),H=a("td"),Ue=k(`Auto expand relations when returning the created record. Ex.:
                `),be(te.$$.fragment),Qe=k(`
                Supports up to 6-levels depth nested relations expansion. `),ze=a("br"),Ke=k(`
                The expanded relations will be appended to the record under the
                `),Le=a("code"),Le.textContent="expand",We=k(" property (eg. "),Pe=a("code"),Pe.textContent='"expand": {"relField1": {...}, ...}',Ye=k(`).
                `),Ge=a("br"),Xe=k(`
                Only the relations to which the request user has permissions to `),ge=a("strong"),ge.textContent="view",Ze=k(" will be expanded."),xe=u(),be(le.$$.fragment),Fe=u(),oe=a("div"),oe.textContent="Responses",Ae=u(),X=a("div"),re=a("div");for(let t=0;t<J.length;t+=1)J[t].c();tt=u(),de=a("div");for(let t=0;t<A.length;t+=1)A[t].c();h(e,"class","m-b-sm"),h(f,"class","content txt-lg m-b-sm"),h(O,"class","m-b-xs"),h(M,"class","label label-primary"),h(F,"class","content"),h(_,"class","alert alert-success"),h(se,"class","section-title"),h(W,"class","table-compact table-border m-b-base"),h(ie,"class","section-title"),h(Y,"class","table-compact table-border m-b-base"),h(oe,"class","section-title"),h(re,"class","tabs-header compact left"),h(de,"class","tabs-content"),h(X,"class","tabs")},m(t,i){d(t,e,i),n(e,l),n(e,m),n(e,b),d(t,r,i),d(t,f,i),n(f,v),n(v,T),n(v,q),n(q,L),n(v,D),n(f,E),n(f,g),n(f,I),n(f,j),d(t,R,i),_e(C,t,i),d(t,N,i),d(t,O,i),d(t,w,i),d(t,_,i),n(_,M),n(_,z),n(_,F),n(F,p),n(p,S),n(p,ee),n(ee,ve),n(p,Re),n(_,je),B&&B.m(_,null),d(t,he,i),d(t,se,i),d(t,we,i),d(t,W,i),n(W,$e),n(W,De),n(W,U),n(U,Ce),n(U,Ne),P&&P.m(U,null),n(U,Se);for(let y=0;y<V.length;y+=1)V[y]&&V[y].m(U,null);d(t,Te,i),d(t,ie,i),d(t,qe,i),d(t,Y,i),n(Y,Oe),n(Y,Je),n(Y,ae),n(ae,G),n(G,Me),n(G,Ee),n(G,He),n(G,Ie),n(G,H),n(H,Ue),_e(te,H,null),n(H,Qe),n(H,ze),n(H,Ke),n(H,Le),n(H,We),n(H,Pe),n(H,Ye),n(H,Ge),n(H,Xe),n(H,ge),n(H,Ze),n(ae,xe),_e(le,ae,null),d(t,Fe,i),d(t,oe,i),d(t,Ae,i),d(t,X,i),n(X,re);for(let y=0;y<J.length;y+=1)J[y]&&J[y].m(re,null);n(X,tt),n(X,de);for(let y=0;y<A.length;y+=1)A[y]&&A[y].m(de,null);Z=!0},p(t,[i]){var ft,pt,ut,mt,bt;(!Z||i&1)&&s!==(s=t[0].name+"")&&x(m,s),(!Z||i&1)&&$!==($=t[0].name+"")&&x(L,$);const y={};i&51&&(y.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${t[5]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},t[4],Q.dummyCollectionSchemaData(t[0])),null,4)};

const record = await pb.collection('${(ft=t[0])==null?void 0:ft.name}').create(data);
`+(t[1]?`
// (optional) send an email verification request
await pb.collection('${(pt=t[0])==null?void 0:pt.name}').requestVerification('test@example.com');
`:"")),i&51&&(y.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${t[5]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},t[4],Q.dummyCollectionSchemaData(t[0])),null,2)};

final record = await pb.collection('${(ut=t[0])==null?void 0:ut.name}').create(body: body);
`+(t[1]?`
// (optional) send an email verification request
await pb.collection('${(mt=t[0])==null?void 0:mt.name}').requestVerification('test@example.com');
`:"")),C.$set(y),(!Z||i&1)&&K!==(K=t[0].name+"")&&x(ve,K),t[6]?B||(B=ht(),B.c(),B.m(_,null)):B&&(B.d(1),B=null),t[1]?P?P.p(t,i):(P=wt(t),P.c(),P.m(U,Se)):P&&(P.d(1),P=null),i&1&&(ue=(bt=t[0])==null?void 0:bt.schema,V=Be(V,i,nt,1,t,ue,Ve,U,_t,$t,null,vt)),i&12&&(me=t[3],J=Be(J,i,st,1,t,me,et,re,_t,Ct,null,yt)),i&12&&(ce=t[3],Ht(),A=Be(A,i,it,1,t,ce,lt,de,Lt,St,null,kt),Pt())},i(t){if(!Z){fe(C.$$.fragment,t),fe(te.$$.fragment,t),fe(le.$$.fragment,t);for(let i=0;i<ce.length;i+=1)fe(A[i]);Z=!0}},o(t){pe(C.$$.fragment,t),pe(te.$$.fragment,t),pe(le.$$.fragment,t);for(let i=0;i<A.length;i+=1)pe(A[i]);Z=!1},d(t){t&&c(e),t&&c(r),t&&c(f),t&&c(R),ke(C,t),t&&c(N),t&&c(O),t&&c(w),t&&c(_),B&&B.d(),t&&c(he),t&&c(se),t&&c(we),t&&c(W),P&&P.d();for(let i=0;i<V.length;i+=1)V[i].d();t&&c(Te),t&&c(ie),t&&c(qe),t&&c(Y),ke(te),ke(le),t&&c(Fe),t&&c(oe),t&&c(Ae),t&&c(X);for(let i=0;i<J.length;i+=1)J[i].d();for(let i=0;i<A.length;i+=1)A[i].d()}}}function Wt(o,e,l){let s,m,b,{collection:r}=e,f=200,v=[],T={};const q=$=>l(2,f=$.code);return o.$$set=$=>{"collection"in $&&l(0,r=$.collection)},o.$$.update=()=>{var $,L;o.$$.dirty&1&&l(1,s=r.type==="auth"),o.$$.dirty&1&&l(6,m=(r==null?void 0:r.createRule)===null),o.$$.dirty&1&&l(3,v=[{code:200,body:JSON.stringify(Q.dummyCollectionRecord(r),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${(L=($=r==null?void 0:r.schema)==null?void 0:$[0])==null?void 0:L.name}": {
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
            `}]),o.$$.dirty&2&&(s?l(4,T={username:"test_username",email:"test@example.com",emailVisibility:!0,password:"12345678",passwordConfirm:"12345678"}):l(4,T={}))},l(5,b=Q.getApiExampleUrl(gt.baseUrl)),[r,s,f,v,T,b,m,q]}class Zt extends qt{constructor(e){super(),Ot(this,e,Wt,Kt,Mt,{collection:0})}}export{Zt as default};
