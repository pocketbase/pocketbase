import{S as Lt,i as Ht,s as Pt,C as Q,O as At,e as a,w as k,b as m,c as Pe,f as h,g as r,h as n,m as Re,x,P as He,Q as ht,k as Rt,R as gt,n as Bt,t as fe,a as pe,o as d,d as ge,L as Ft,p as jt,r as ue,u as Dt,y as le}from"./index.2d20c7a4.js";import{S as Nt}from"./SdkTabs.dcaa443a.js";function wt(o,e,l){const s=o.slice();return s[7]=e[l],s}function Ct(o,e,l){const s=o.slice();return s[7]=e[l],s}function St(o,e,l){const s=o.slice();return s[12]=e[l],s}function $t(o){let e;return{c(){e=a("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",h(e,"class","txt-hint txt-sm txt-right")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function Tt(o){let e,l,s,b,p,c,f,y,T,w,O,F,D,V,H,I,j,g,S,N,q,C,_;function M(u,$){var ee,K;return(K=(ee=u[0])==null?void 0:ee.options)!=null&&K.requireEmail?It:Vt}let z=M(o),P=z(o);return{c(){e=a("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',l=m(),s=a("tr"),s.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>username</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The username of the auth record.
                    <br/>
                    If not set, it will be auto generated.</td>`,b=m(),p=a("tr"),c=a("td"),f=a("div"),P.c(),y=m(),T=a("span"),T.textContent="email",w=m(),O=a("td"),O.innerHTML='<span class="label">String</span>',F=m(),D=a("td"),D.textContent="Auth record email address.",V=m(),H=a("tr"),H.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>emailVisibility</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Whether to show/hide the auth record email when fetching the record data.</td>`,I=m(),j=a("tr"),j.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                        <span>password</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Auth record password.</td>`,g=m(),S=a("tr"),S.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                        <span>passwordConfirm</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Auth record password confirmation.</td>`,N=m(),q=a("tr"),q.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>verified</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,C=m(),_=a("tr"),_.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>',h(f,"class","inline-flex")},m(u,$){r(u,e,$),r(u,l,$),r(u,s,$),r(u,b,$),r(u,p,$),n(p,c),n(c,f),P.m(f,null),n(f,y),n(f,T),n(p,w),n(p,O),n(p,F),n(p,D),r(u,V,$),r(u,H,$),r(u,I,$),r(u,j,$),r(u,g,$),r(u,S,$),r(u,N,$),r(u,q,$),r(u,C,$),r(u,_,$)},p(u,$){z!==(z=M(u))&&(P.d(1),P=z(u),P&&(P.c(),P.m(f,y)))},d(u){u&&d(e),u&&d(l),u&&d(s),u&&d(b),u&&d(p),P.d(),u&&d(V),u&&d(H),u&&d(I),u&&d(j),u&&d(g),u&&d(S),u&&d(N),u&&d(q),u&&d(C),u&&d(_)}}}function Vt(o){let e;return{c(){e=a("span"),e.textContent="Optional",h(e,"class","label label-warning")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function It(o){let e;return{c(){e=a("span"),e.textContent="Required",h(e,"class","label label-success")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function Jt(o){let e;return{c(){e=a("span"),e.textContent="Optional",h(e,"class","label label-warning")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function Et(o){let e;return{c(){e=a("span"),e.textContent="Required",h(e,"class","label label-success")},m(l,s){r(l,e,s)},d(l){l&&d(e)}}}function Ut(o){var p;let e,l=((p=o[12].options)==null?void 0:p.maxSelect)===1?"id":"ids",s,b;return{c(){e=k("Relation record "),s=k(l),b=k(".")},m(c,f){r(c,e,f),r(c,s,f),r(c,b,f)},p(c,f){var y;f&1&&l!==(l=((y=c[12].options)==null?void 0:y.maxSelect)===1?"id":"ids")&&x(s,l)},d(c){c&&d(e),c&&d(s),c&&d(b)}}}function Qt(o){let e,l,s,b,p;return{c(){e=k("File object."),l=a("br"),s=k(`
                        Set to `),b=a("code"),b.textContent="null",p=k(" to delete already uploaded file(s).")},m(c,f){r(c,e,f),r(c,l,f),r(c,s,f),r(c,b,f),r(c,p,f)},p:le,d(c){c&&d(e),c&&d(l),c&&d(s),c&&d(b),c&&d(p)}}}function zt(o){let e;return{c(){e=k("URL address.")},m(l,s){r(l,e,s)},p:le,d(l){l&&d(e)}}}function Kt(o){let e;return{c(){e=k("Email address.")},m(l,s){r(l,e,s)},p:le,d(l){l&&d(e)}}}function Wt(o){let e;return{c(){e=k("JSON array or object.")},m(l,s){r(l,e,s)},p:le,d(l){l&&d(e)}}}function Yt(o){let e;return{c(){e=k("Number value.")},m(l,s){r(l,e,s)},p:le,d(l){l&&d(e)}}}function Gt(o){let e;return{c(){e=k("Plain text value.")},m(l,s){r(l,e,s)},p:le,d(l){l&&d(e)}}}function qt(o,e){let l,s,b,p,c,f=e[12].name+"",y,T,w,O,F=Q.getFieldValueType(e[12])+"",D,V,H,I;function j(_,M){return _[12].required?Et:Jt}let g=j(e),S=g(e);function N(_,M){if(_[12].type==="text")return Gt;if(_[12].type==="number")return Yt;if(_[12].type==="json")return Wt;if(_[12].type==="email")return Kt;if(_[12].type==="url")return zt;if(_[12].type==="file")return Qt;if(_[12].type==="relation")return Ut}let q=N(e),C=q&&q(e);return{key:o,first:null,c(){l=a("tr"),s=a("td"),b=a("div"),S.c(),p=m(),c=a("span"),y=k(f),T=m(),w=a("td"),O=a("span"),D=k(F),V=m(),H=a("td"),C&&C.c(),I=m(),h(b,"class","inline-flex"),h(O,"class","label"),this.first=l},m(_,M){r(_,l,M),n(l,s),n(s,b),S.m(b,null),n(b,p),n(b,c),n(c,y),n(l,T),n(l,w),n(w,O),n(O,D),n(l,V),n(l,H),C&&C.m(H,null),n(l,I)},p(_,M){e=_,g!==(g=j(e))&&(S.d(1),S=g(e),S&&(S.c(),S.m(b,p))),M&1&&f!==(f=e[12].name+"")&&x(y,f),M&1&&F!==(F=Q.getFieldValueType(e[12])+"")&&x(D,F),q===(q=N(e))&&C?C.p(e,M):(C&&C.d(1),C=q&&q(e),C&&(C.c(),C.m(H,null)))},d(_){_&&d(l),S.d(),C&&C.d()}}}function Ot(o,e){let l,s=e[7].code+"",b,p,c,f;function y(){return e[6](e[7])}return{key:o,first:null,c(){l=a("button"),b=k(s),p=m(),h(l,"class","tab-item"),ue(l,"active",e[1]===e[7].code),this.first=l},m(T,w){r(T,l,w),n(l,b),n(l,p),c||(f=Dt(l,"click",y),c=!0)},p(T,w){e=T,w&4&&s!==(s=e[7].code+"")&&x(b,s),w&6&&ue(l,"active",e[1]===e[7].code)},d(T){T&&d(l),c=!1,f()}}}function Mt(o,e){let l,s,b,p;return s=new At({props:{content:e[7].body}}),{key:o,first:null,c(){l=a("div"),Pe(s.$$.fragment),b=m(),h(l,"class","tab-item"),ue(l,"active",e[1]===e[7].code),this.first=l},m(c,f){r(c,l,f),Re(s,l,null),n(l,b),p=!0},p(c,f){e=c;const y={};f&4&&(y.content=e[7].body),s.$set(y),(!p||f&6)&&ue(l,"active",e[1]===e[7].code)},i(c){p||(fe(s.$$.fragment,c),p=!0)},o(c){pe(s.$$.fragment,c),p=!1},d(c){c&&d(l),ge(s)}}}function Xt(o){var st,it,at,ot,rt,dt,ct,ft;let e,l,s=o[0].name+"",b,p,c,f,y,T,w,O=o[0].name+"",F,D,V,H,I,j,g,S,N,q,C,_,M,z,P,u,$,ee,K=o[0].name+"",me,Be,Fe,be,ne,_e,W,ke,je,J,ye,De,ve,E=[],Ne=new Map,he,se,we,Y,Ce,Ve,Se,G,$e,Ie,Te,Je,A,Ee,te,Ue,Qe,ze,qe,Ke,Oe,We,Ye,Ge,Me,Xe,Ae,ie,Le,X,ae,U=[],Ze=new Map,xe,oe,B=[],et=new Map,Z;S=new Nt({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${o[4]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},o[3],Q.dummyCollectionSchemaData(o[0])),null,4)};

const record = await pb.collection('${(st=o[0])==null?void 0:st.name}').create(data);
`+((it=o[0])!=null&&it.isAuth?`
// (optional) send an email verification request
await pb.collection('${(at=o[0])==null?void 0:at.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${o[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},o[3],Q.dummyCollectionSchemaData(o[0])),null,2)};

final record = await pb.collection('${(ot=o[0])==null?void 0:ot.name}').create(body: body);
`+((rt=o[0])!=null&&rt.isAuth?`
// (optional) send an email verification request
await pb.collection('${(dt=o[0])==null?void 0:dt.name}').requestVerification('test@example.com');
`:"")}});let R=o[5]&&$t(),L=((ct=o[0])==null?void 0:ct.isAuth)&&Tt(o),de=(ft=o[0])==null?void 0:ft.schema;const tt=t=>t[12].name;for(let t=0;t<de.length;t+=1){let i=St(o,de,t),v=tt(i);Ne.set(v,E[t]=qt(v,i))}te=new At({props:{content:"?expand=relField1,relField2.subRelField"}});let ce=o[2];const lt=t=>t[7].code;for(let t=0;t<ce.length;t+=1){let i=Ct(o,ce,t),v=lt(i);Ze.set(v,U[t]=Ot(v,i))}let re=o[2];const nt=t=>t[7].code;for(let t=0;t<re.length;t+=1){let i=wt(o,re,t),v=nt(i);et.set(v,B[t]=Mt(v,i))}return{c(){e=a("h3"),l=k("Create ("),b=k(s),p=k(")"),c=m(),f=a("div"),y=a("p"),T=k("Create a new "),w=a("strong"),F=k(O),D=k(" record."),V=m(),H=a("p"),H.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,I=m(),j=a("p"),j.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/#uploading-files" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,g=m(),Pe(S.$$.fragment),N=m(),q=a("h6"),q.textContent="API details",C=m(),_=a("div"),M=a("strong"),M.textContent="POST",z=m(),P=a("div"),u=a("p"),$=k("/api/collections/"),ee=a("strong"),me=k(K),Be=k("/records"),Fe=m(),R&&R.c(),be=m(),ne=a("div"),ne.textContent="Body Parameters",_e=m(),W=a("table"),ke=a("thead"),ke.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,je=m(),J=a("tbody"),ye=a("tr"),ye.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                    <span>id</span></div></td> 
            <td><span class="label">String</span></td> 
            <td><strong>15 characters string</strong> to store as record ID.
                <br/>
                If not set, it will be auto generated.</td>`,De=m(),L&&L.c(),ve=m();for(let t=0;t<E.length;t+=1)E[t].c();he=m(),se=a("div"),se.textContent="Query parameters",we=m(),Y=a("table"),Ce=a("thead"),Ce.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Ve=m(),Se=a("tbody"),G=a("tr"),$e=a("td"),$e.textContent="expand",Ie=m(),Te=a("td"),Te.innerHTML='<span class="label">String</span>',Je=m(),A=a("td"),Ee=k(`Auto expand relations when returning the created record. Ex.:
                `),Pe(te.$$.fragment),Ue=k(`
                Supports up to 6-levels depth nested relations expansion. `),Qe=a("br"),ze=k(`
                The expanded relations will be appended to the record under the
                `),qe=a("code"),qe.textContent="expand",Ke=k(" property (eg. "),Oe=a("code"),Oe.textContent='"expand": {"relField1": {...}, ...}',We=k(`).
                `),Ye=a("br"),Ge=k(`
                Only the relations to which the request user has permissions to `),Me=a("strong"),Me.textContent="view",Xe=k(" will be expanded."),Ae=m(),ie=a("div"),ie.textContent="Responses",Le=m(),X=a("div"),ae=a("div");for(let t=0;t<U.length;t+=1)U[t].c();xe=m(),oe=a("div");for(let t=0;t<B.length;t+=1)B[t].c();h(e,"class","m-b-sm"),h(f,"class","content txt-lg m-b-sm"),h(q,"class","m-b-xs"),h(M,"class","label label-primary"),h(P,"class","content"),h(_,"class","alert alert-success"),h(ne,"class","section-title"),h(W,"class","table-compact table-border m-b-base"),h(se,"class","section-title"),h(Y,"class","table-compact table-border m-b-base"),h(ie,"class","section-title"),h(ae,"class","tabs-header compact left"),h(oe,"class","tabs-content"),h(X,"class","tabs")},m(t,i){r(t,e,i),n(e,l),n(e,b),n(e,p),r(t,c,i),r(t,f,i),n(f,y),n(y,T),n(y,w),n(w,F),n(y,D),n(f,V),n(f,H),n(f,I),n(f,j),r(t,g,i),Re(S,t,i),r(t,N,i),r(t,q,i),r(t,C,i),r(t,_,i),n(_,M),n(_,z),n(_,P),n(P,u),n(u,$),n(u,ee),n(ee,me),n(u,Be),n(_,Fe),R&&R.m(_,null),r(t,be,i),r(t,ne,i),r(t,_e,i),r(t,W,i),n(W,ke),n(W,je),n(W,J),n(J,ye),n(J,De),L&&L.m(J,null),n(J,ve);for(let v=0;v<E.length;v+=1)E[v].m(J,null);r(t,he,i),r(t,se,i),r(t,we,i),r(t,Y,i),n(Y,Ce),n(Y,Ve),n(Y,Se),n(Se,G),n(G,$e),n(G,Ie),n(G,Te),n(G,Je),n(G,A),n(A,Ee),Re(te,A,null),n(A,Ue),n(A,Qe),n(A,ze),n(A,qe),n(A,Ke),n(A,Oe),n(A,We),n(A,Ye),n(A,Ge),n(A,Me),n(A,Xe),r(t,Ae,i),r(t,ie,i),r(t,Le,i),r(t,X,i),n(X,ae);for(let v=0;v<U.length;v+=1)U[v].m(ae,null);n(X,xe),n(X,oe);for(let v=0;v<B.length;v+=1)B[v].m(oe,null);Z=!0},p(t,[i]){var pt,ut,mt,bt,_t,kt,yt,vt;(!Z||i&1)&&s!==(s=t[0].name+"")&&x(b,s),(!Z||i&1)&&O!==(O=t[0].name+"")&&x(F,O);const v={};i&25&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${t[4]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},t[3],Q.dummyCollectionSchemaData(t[0])),null,4)};

const record = await pb.collection('${(pt=t[0])==null?void 0:pt.name}').create(data);
`+((ut=t[0])!=null&&ut.isAuth?`
// (optional) send an email verification request
await pb.collection('${(mt=t[0])==null?void 0:mt.name}').requestVerification('test@example.com');
`:"")),i&25&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${t[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},t[3],Q.dummyCollectionSchemaData(t[0])),null,2)};

final record = await pb.collection('${(bt=t[0])==null?void 0:bt.name}').create(body: body);
`+((_t=t[0])!=null&&_t.isAuth?`
// (optional) send an email verification request
await pb.collection('${(kt=t[0])==null?void 0:kt.name}').requestVerification('test@example.com');
`:"")),S.$set(v),(!Z||i&1)&&K!==(K=t[0].name+"")&&x(me,K),t[5]?R||(R=$t(),R.c(),R.m(_,null)):R&&(R.d(1),R=null),(yt=t[0])!=null&&yt.isAuth?L?L.p(t,i):(L=Tt(t),L.c(),L.m(J,ve)):L&&(L.d(1),L=null),i&1&&(de=(vt=t[0])==null?void 0:vt.schema,E=He(E,i,tt,1,t,de,Ne,J,ht,qt,null,St)),i&6&&(ce=t[2],U=He(U,i,lt,1,t,ce,Ze,ae,ht,Ot,null,Ct)),i&6&&(re=t[2],Rt(),B=He(B,i,nt,1,t,re,et,oe,gt,Mt,null,wt),Bt())},i(t){if(!Z){fe(S.$$.fragment,t),fe(te.$$.fragment,t);for(let i=0;i<re.length;i+=1)fe(B[i]);Z=!0}},o(t){pe(S.$$.fragment,t),pe(te.$$.fragment,t);for(let i=0;i<B.length;i+=1)pe(B[i]);Z=!1},d(t){t&&d(e),t&&d(c),t&&d(f),t&&d(g),ge(S,t),t&&d(N),t&&d(q),t&&d(C),t&&d(_),R&&R.d(),t&&d(be),t&&d(ne),t&&d(_e),t&&d(W),L&&L.d();for(let i=0;i<E.length;i+=1)E[i].d();t&&d(he),t&&d(se),t&&d(we),t&&d(Y),ge(te),t&&d(Ae),t&&d(ie),t&&d(Le),t&&d(X);for(let i=0;i<U.length;i+=1)U[i].d();for(let i=0;i<B.length;i+=1)B[i].d()}}}function Zt(o,e,l){let s,b,{collection:p=new Ft}=e,c=200,f=[],y={};const T=w=>l(1,c=w.code);return o.$$set=w=>{"collection"in w&&l(0,p=w.collection)},o.$$.update=()=>{var w,O;o.$$.dirty&1&&l(5,s=(p==null?void 0:p.createRule)===null),o.$$.dirty&1&&l(2,f=[{code:200,body:JSON.stringify(Q.dummyCollectionRecord(p),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${(O=(w=p==null?void 0:p.schema)==null?void 0:w[0])==null?void 0:O.name}": {
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
            `}]),o.$$.dirty&1&&(p.isAuth?l(3,y={username:"test_username",email:"test@example.com",emailVisibility:!0,password:"12345678",passwordConfirm:"12345678"}):l(3,y={}))},l(4,b=Q.getApiExampleUrl(jt.baseUrl)),[p,c,f,y,b,s,T]}class tl extends Lt{constructor(e){super(),Ht(this,e,Zt,Xt,Pt,{collection:0})}}export{tl as default};
