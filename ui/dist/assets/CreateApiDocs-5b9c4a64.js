import{S as Pt,i as gt,s as Ft,C as Q,M as Lt,e as a,w as k,b as m,c as be,f as h,g as r,h as n,m as _e,x,N as Be,P as $t,k as Bt,Q as Rt,n as jt,t as fe,a as pe,o as d,d as ke,T as Dt,p as Nt,r as ye,u as Vt,y as ne}from"./index-a65ca895.js";import{S as Jt}from"./SdkTabs-ad912c8f.js";import{F as Et}from"./FieldsQueryParam-ba250473.js";function Ct(o,e,l){const s=o.slice();return s[7]=e[l],s}function St(o,e,l){const s=o.slice();return s[7]=e[l],s}function Tt(o,e,l){const s=o.slice();return s[12]=e[l],s}function qt(o){let e;return{c(){e=a("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",h(e,"class","txt-hint txt-sm txt-right")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function Mt(o){let e,l,s,b,p,c,f,v,T,w,M,R,D,E,L,I,j,B,C,N,q,$,_;function O(u,S){var ee,K;return(K=(ee=u[0])==null?void 0:ee.options)!=null&&K.requireEmail?Ut:It}let z=O(o),P=z(o);return{c(){e=a("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',l=m(),s=a("tr"),s.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>username</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The username of the auth record.
                    <br/>
                    If not set, it will be auto generated.</td>`,b=m(),p=a("tr"),c=a("td"),f=a("div"),P.c(),v=m(),T=a("span"),T.textContent="email",w=m(),M=a("td"),M.innerHTML='<span class="label">String</span>',R=m(),D=a("td"),D.textContent="Auth record email address.",E=m(),L=a("tr"),L.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>emailVisibility</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Whether to show/hide the auth record email when fetching the record data.</td>`,I=m(),j=a("tr"),j.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                        <span>password</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Auth record password.</td>`,B=m(),C=a("tr"),C.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                        <span>passwordConfirm</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Auth record password confirmation.</td>`,N=m(),q=a("tr"),q.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>verified</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,$=m(),_=a("tr"),_.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>',h(f,"class","inline-flex")},m(u,S){r(u,e,S),r(u,l,S),r(u,s,S),r(u,b,S),r(u,p,S),n(p,c),n(c,f),P.m(f,null),n(f,v),n(f,T),n(p,w),n(p,M),n(p,R),n(p,D),r(u,E,S),r(u,L,S),r(u,I,S),r(u,j,S),r(u,B,S),r(u,C,S),r(u,N,S),r(u,q,S),r(u,$,S),r(u,_,S)},p(u,S){z!==(z=O(u))&&(P.d(1),P=z(u),P&&(P.c(),P.m(f,v)))},d(u){u&&d(e),u&&d(l),u&&d(s),u&&d(b),u&&d(p),P.d(),u&&d(E),u&&d(L),u&&d(I),u&&d(j),u&&d(B),u&&d(C),u&&d(N),u&&d(q),u&&d($),u&&d(_)}}}function It(o){let e;return{c(){e=a("span"),e.textContent="Optional",h(e,"class","label label-warning")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function Ut(o){let e;return{c(){e=a("span"),e.textContent="Required",h(e,"class","label label-success")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function Qt(o){let e;return{c(){e=a("span"),e.textContent="Optional",h(e,"class","label label-warning")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function zt(o){let e;return{c(){e=a("span"),e.textContent="Required",h(e,"class","label label-success")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function Kt(o){var p;let e,l=((p=o[12].options)==null?void 0:p.maxSelect)===1?"id":"ids",s,b;return{c(){e=k("Relation record "),s=k(l),b=k(".")},m(c,f){r(c,e,f),r(c,s,f),r(c,b,f)},p(c,f){var v;f&1&&l!==(l=((v=c[12].options)==null?void 0:v.maxSelect)===1?"id":"ids")&&x(s,l)},d(c){c&&d(e),c&&d(s),c&&d(b)}}}function Wt(o){let e,l,s,b,p;return{c(){e=k("File object."),l=a("br"),s=k(`
                        Set to `),b=a("code"),b.textContent="null",p=k(" to delete already uploaded file(s).")},m(c,f){r(c,e,f),r(c,l,f),r(c,s,f),r(c,b,f),r(c,p,f)},p:ne,d(c){c&&d(e),c&&d(l),c&&d(s),c&&d(b),c&&d(p)}}}function Yt(o){let e;return{c(){e=k("URL address.")},m(l,s){r(l,e,s)},p:ne,d(l){l&&d(e)}}}function Gt(o){let e;return{c(){e=k("Email address.")},m(l,s){r(l,e,s)},p:ne,d(l){l&&d(e)}}}function Xt(o){let e;return{c(){e=k("JSON array or object.")},m(l,s){r(l,e,s)},p:ne,d(l){l&&d(e)}}}function Zt(o){let e;return{c(){e=k("Number value.")},m(l,s){r(l,e,s)},p:ne,d(l){l&&d(e)}}}function xt(o){let e;return{c(){e=k("Plain text value.")},m(l,s){r(l,e,s)},p:ne,d(l){l&&d(e)}}}function Ot(o,e){let l,s,b,p,c,f=e[12].name+"",v,T,w,M,R=Q.getFieldValueType(e[12])+"",D,E,L,I;function j(_,O){return _[12].required?zt:Qt}let B=j(e),C=B(e);function N(_,O){if(_[12].type==="text")return xt;if(_[12].type==="number")return Zt;if(_[12].type==="json")return Xt;if(_[12].type==="email")return Gt;if(_[12].type==="url")return Yt;if(_[12].type==="file")return Wt;if(_[12].type==="relation")return Kt}let q=N(e),$=q&&q(e);return{key:o,first:null,c(){l=a("tr"),s=a("td"),b=a("div"),C.c(),p=m(),c=a("span"),v=k(f),T=m(),w=a("td"),M=a("span"),D=k(R),E=m(),L=a("td"),$&&$.c(),I=m(),h(b,"class","inline-flex"),h(M,"class","label"),this.first=l},m(_,O){r(_,l,O),n(l,s),n(s,b),C.m(b,null),n(b,p),n(b,c),n(c,v),n(l,T),n(l,w),n(w,M),n(M,D),n(l,E),n(l,L),$&&$.m(L,null),n(l,I)},p(_,O){e=_,B!==(B=j(e))&&(C.d(1),C=B(e),C&&(C.c(),C.m(b,p))),O&1&&f!==(f=e[12].name+"")&&x(v,f),O&1&&R!==(R=Q.getFieldValueType(e[12])+"")&&x(D,R),q===(q=N(e))&&$?$.p(e,O):($&&$.d(1),$=q&&q(e),$&&($.c(),$.m(L,null)))},d(_){_&&d(l),C.d(),$&&$.d()}}}function At(o,e){let l,s=e[7].code+"",b,p,c,f;function v(){return e[6](e[7])}return{key:o,first:null,c(){l=a("button"),b=k(s),p=m(),h(l,"class","tab-item"),ye(l,"active",e[1]===e[7].code),this.first=l},m(T,w){r(T,l,w),n(l,b),n(l,p),c||(f=Vt(l,"click",v),c=!0)},p(T,w){e=T,w&4&&s!==(s=e[7].code+"")&&x(b,s),w&6&&ye(l,"active",e[1]===e[7].code)},d(T){T&&d(l),c=!1,f()}}}function Ht(o,e){let l,s,b,p;return s=new Lt({props:{content:e[7].body}}),{key:o,first:null,c(){l=a("div"),be(s.$$.fragment),b=m(),h(l,"class","tab-item"),ye(l,"active",e[1]===e[7].code),this.first=l},m(c,f){r(c,l,f),_e(s,l,null),n(l,b),p=!0},p(c,f){e=c;const v={};f&4&&(v.content=e[7].body),s.$set(v),(!p||f&6)&&ye(l,"active",e[1]===e[7].code)},i(c){p||(fe(s.$$.fragment,c),p=!0)},o(c){pe(s.$$.fragment,c),p=!1},d(c){c&&d(l),ke(s)}}}function el(o){var at,ot,rt,dt,ct,ft,pt,ut;let e,l,s=o[0].name+"",b,p,c,f,v,T,w,M=o[0].name+"",R,D,E,L,I,j,B,C,N,q,$,_,O,z,P,u,S,ee,K=o[0].name+"",ve,Re,je,he,se,we,W,$e,De,U,Ce,Ne,Se,V=[],Ve=new Map,Te,ie,qe,Y,Me,Je,ae,G,Oe,Ee,Ae,Ie,A,Ue,te,Qe,ze,Ke,He,We,Le,Ye,Ge,Xe,Pe,Ze,xe,le,ge,oe,Fe,X,re,J=[],et=new Map,tt,de,g=[],lt=new Map,Z;C=new Jt({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${o[4]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},o[3],Q.dummyCollectionSchemaData(o[0])),null,4)};

const record = await pb.collection('${(at=o[0])==null?void 0:at.name}').create(data);
`+((ot=o[0])!=null&&ot.isAuth?`
// (optional) send an email verification request
await pb.collection('${(rt=o[0])==null?void 0:rt.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${o[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},o[3],Q.dummyCollectionSchemaData(o[0])),null,2)};

final record = await pb.collection('${(dt=o[0])==null?void 0:dt.name}').create(body: body);
`+((ct=o[0])!=null&&ct.isAuth?`
// (optional) send an email verification request
await pb.collection('${(ft=o[0])==null?void 0:ft.name}').requestVerification('test@example.com');
`:"")}});let F=o[5]&&qt(),H=((pt=o[0])==null?void 0:pt.isAuth)&&Mt(o),ue=(ut=o[0])==null?void 0:ut.schema;const nt=t=>t[12].name;for(let t=0;t<ue.length;t+=1){let i=Tt(o,ue,t),y=nt(i);Ve.set(y,V[t]=Ot(y,i))}te=new Lt({props:{content:"?expand=relField1,relField2.subRelField"}}),le=new Et({});let me=o[2];const st=t=>t[7].code;for(let t=0;t<me.length;t+=1){let i=St(o,me,t),y=st(i);et.set(y,J[t]=At(y,i))}let ce=o[2];const it=t=>t[7].code;for(let t=0;t<ce.length;t+=1){let i=Ct(o,ce,t),y=it(i);lt.set(y,g[t]=Ht(y,i))}return{c(){e=a("h3"),l=k("Create ("),b=k(s),p=k(")"),c=m(),f=a("div"),v=a("p"),T=k("Create a new "),w=a("strong"),R=k(M),D=k(" record."),E=m(),L=a("p"),L.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,I=m(),j=a("p"),j.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,B=m(),be(C.$$.fragment),N=m(),q=a("h6"),q.textContent="API details",$=m(),_=a("div"),O=a("strong"),O.textContent="POST",z=m(),P=a("div"),u=a("p"),S=k("/api/collections/"),ee=a("strong"),ve=k(K),Re=k("/records"),je=m(),F&&F.c(),he=m(),se=a("div"),se.textContent="Body Parameters",we=m(),W=a("table"),$e=a("thead"),$e.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,De=m(),U=a("tbody"),Ce=a("tr"),Ce.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                    <span>id</span></div></td> 
            <td><span class="label">String</span></td> 
            <td><strong>15 characters string</strong> to store as record ID.
                <br/>
                If not set, it will be auto generated.</td>`,Ne=m(),H&&H.c(),Se=m();for(let t=0;t<V.length;t+=1)V[t].c();Te=m(),ie=a("div"),ie.textContent="Query parameters",qe=m(),Y=a("table"),Me=a("thead"),Me.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Je=m(),ae=a("tbody"),G=a("tr"),Oe=a("td"),Oe.textContent="expand",Ee=m(),Ae=a("td"),Ae.innerHTML='<span class="label">String</span>',Ie=m(),A=a("td"),Ue=k(`Auto expand relations when returning the created record. Ex.:
                `),be(te.$$.fragment),Qe=k(`
                Supports up to 6-levels depth nested relations expansion. `),ze=a("br"),Ke=k(`
                The expanded relations will be appended to the record under the
                `),He=a("code"),He.textContent="expand",We=k(" property (eg. "),Le=a("code"),Le.textContent='"expand": {"relField1": {...}, ...}',Ye=k(`).
                `),Ge=a("br"),Xe=k(`
                Only the relations to which the request user has permissions to `),Pe=a("strong"),Pe.textContent="view",Ze=k(" will be expanded."),xe=m(),be(le.$$.fragment),ge=m(),oe=a("div"),oe.textContent="Responses",Fe=m(),X=a("div"),re=a("div");for(let t=0;t<J.length;t+=1)J[t].c();tt=m(),de=a("div");for(let t=0;t<g.length;t+=1)g[t].c();h(e,"class","m-b-sm"),h(f,"class","content txt-lg m-b-sm"),h(q,"class","m-b-xs"),h(O,"class","label label-primary"),h(P,"class","content"),h(_,"class","alert alert-success"),h(se,"class","section-title"),h(W,"class","table-compact table-border m-b-base"),h(ie,"class","section-title"),h(Y,"class","table-compact table-border m-b-base"),h(oe,"class","section-title"),h(re,"class","tabs-header compact left"),h(de,"class","tabs-content"),h(X,"class","tabs")},m(t,i){r(t,e,i),n(e,l),n(e,b),n(e,p),r(t,c,i),r(t,f,i),n(f,v),n(v,T),n(v,w),n(w,R),n(v,D),n(f,E),n(f,L),n(f,I),n(f,j),r(t,B,i),_e(C,t,i),r(t,N,i),r(t,q,i),r(t,$,i),r(t,_,i),n(_,O),n(_,z),n(_,P),n(P,u),n(u,S),n(u,ee),n(ee,ve),n(u,Re),n(_,je),F&&F.m(_,null),r(t,he,i),r(t,se,i),r(t,we,i),r(t,W,i),n(W,$e),n(W,De),n(W,U),n(U,Ce),n(U,Ne),H&&H.m(U,null),n(U,Se);for(let y=0;y<V.length;y+=1)V[y]&&V[y].m(U,null);r(t,Te,i),r(t,ie,i),r(t,qe,i),r(t,Y,i),n(Y,Me),n(Y,Je),n(Y,ae),n(ae,G),n(G,Oe),n(G,Ee),n(G,Ae),n(G,Ie),n(G,A),n(A,Ue),_e(te,A,null),n(A,Qe),n(A,ze),n(A,Ke),n(A,He),n(A,We),n(A,Le),n(A,Ye),n(A,Ge),n(A,Xe),n(A,Pe),n(A,Ze),n(ae,xe),_e(le,ae,null),r(t,ge,i),r(t,oe,i),r(t,Fe,i),r(t,X,i),n(X,re);for(let y=0;y<J.length;y+=1)J[y]&&J[y].m(re,null);n(X,tt),n(X,de);for(let y=0;y<g.length;y+=1)g[y]&&g[y].m(de,null);Z=!0},p(t,[i]){var mt,bt,_t,kt,yt,vt,ht,wt;(!Z||i&1)&&s!==(s=t[0].name+"")&&x(b,s),(!Z||i&1)&&M!==(M=t[0].name+"")&&x(R,M);const y={};i&25&&(y.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${t[4]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},t[3],Q.dummyCollectionSchemaData(t[0])),null,4)};

const record = await pb.collection('${(mt=t[0])==null?void 0:mt.name}').create(data);
`+((bt=t[0])!=null&&bt.isAuth?`
// (optional) send an email verification request
await pb.collection('${(_t=t[0])==null?void 0:_t.name}').requestVerification('test@example.com');
`:"")),i&25&&(y.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${t[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},t[3],Q.dummyCollectionSchemaData(t[0])),null,2)};

final record = await pb.collection('${(kt=t[0])==null?void 0:kt.name}').create(body: body);
`+((yt=t[0])!=null&&yt.isAuth?`
// (optional) send an email verification request
await pb.collection('${(vt=t[0])==null?void 0:vt.name}').requestVerification('test@example.com');
`:"")),C.$set(y),(!Z||i&1)&&K!==(K=t[0].name+"")&&x(ve,K),t[5]?F||(F=qt(),F.c(),F.m(_,null)):F&&(F.d(1),F=null),(ht=t[0])!=null&&ht.isAuth?H?H.p(t,i):(H=Mt(t),H.c(),H.m(U,Se)):H&&(H.d(1),H=null),i&1&&(ue=(wt=t[0])==null?void 0:wt.schema,V=Be(V,i,nt,1,t,ue,Ve,U,$t,Ot,null,Tt)),i&6&&(me=t[2],J=Be(J,i,st,1,t,me,et,re,$t,At,null,St)),i&6&&(ce=t[2],Bt(),g=Be(g,i,it,1,t,ce,lt,de,Rt,Ht,null,Ct),jt())},i(t){if(!Z){fe(C.$$.fragment,t),fe(te.$$.fragment,t),fe(le.$$.fragment,t);for(let i=0;i<ce.length;i+=1)fe(g[i]);Z=!0}},o(t){pe(C.$$.fragment,t),pe(te.$$.fragment,t),pe(le.$$.fragment,t);for(let i=0;i<g.length;i+=1)pe(g[i]);Z=!1},d(t){t&&d(e),t&&d(c),t&&d(f),t&&d(B),ke(C,t),t&&d(N),t&&d(q),t&&d($),t&&d(_),F&&F.d(),t&&d(he),t&&d(se),t&&d(we),t&&d(W),H&&H.d();for(let i=0;i<V.length;i+=1)V[i].d();t&&d(Te),t&&d(ie),t&&d(qe),t&&d(Y),ke(te),ke(le),t&&d(ge),t&&d(oe),t&&d(Fe),t&&d(X);for(let i=0;i<J.length;i+=1)J[i].d();for(let i=0;i<g.length;i+=1)g[i].d()}}}function tl(o,e,l){let s,b,{collection:p=new Dt}=e,c=200,f=[],v={};const T=w=>l(1,c=w.code);return o.$$set=w=>{"collection"in w&&l(0,p=w.collection)},o.$$.update=()=>{var w,M;o.$$.dirty&1&&l(5,s=(p==null?void 0:p.createRule)===null),o.$$.dirty&1&&l(2,f=[{code:200,body:JSON.stringify(Q.dummyCollectionRecord(p),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${(M=(w=p==null?void 0:p.schema)==null?void 0:w[0])==null?void 0:M.name}": {
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
            `}]),o.$$.dirty&1&&(p.$isAuth?l(3,v={username:"test_username",email:"test@example.com",emailVisibility:!0,password:"12345678",passwordConfirm:"12345678"}):l(3,v={}))},l(4,b=Q.getApiExampleUrl(Nt.baseUrl)),[p,c,f,v,b,s,T]}class il extends Pt{constructor(e){super(),gt(this,e,tl,el,Ft,{collection:0})}}export{il as default};
