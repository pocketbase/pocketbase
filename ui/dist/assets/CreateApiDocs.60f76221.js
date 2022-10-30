import{S as Ct,i as St,s as Tt,C as Q,O as wt,e as a,w as k,b,c as Pe,f as h,g as o,h as n,m as Re,x as Z,P as Ae,Q as pt,k as $t,R as Ot,n as Mt,t as fe,a as pe,o as r,d as ge,L as Lt,p as Ht,r as ue,u as qt,y as le}from"./index.97f016a1.js";import{S as At}from"./SdkTabs.88269ae0.js";function ut(d,e,l){const s=d.slice();return s[7]=e[l],s}function bt(d,e,l){const s=d.slice();return s[7]=e[l],s}function mt(d,e,l){const s=d.slice();return s[12]=e[l],s}function _t(d){let e;return{c(){e=a("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",h(e,"class","txt-hint txt-sm txt-right")},m(l,s){o(l,e,s)},d(l){l&&r(e)}}}function kt(d){let e,l,s,m,p,c,f,y,$,w,M,F,D,I,A,J,j,g,S,N,O,C,_;function L(u,T){var ee,z;return(z=(ee=u[0])==null?void 0:ee.options)!=null&&z.requireEmail?Rt:Pt}let x=L(d),P=x(d);return{c(){e=a("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',l=b(),s=a("tr"),s.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>username</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>The username of the auth record.
                    <br/>
                    If not set, it will be auto generated.</td>`,m=b(),p=a("tr"),c=a("td"),f=a("div"),P.c(),y=b(),$=a("span"),$.textContent="email",w=b(),M=a("td"),M.innerHTML='<span class="label">String</span>',F=b(),D=a("td"),D.textContent="Auth record email address.",I=b(),A=a("tr"),A.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>emailVisibility</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Whether to show/hide the auth record email when fetching the record data.</td>`,J=b(),j=a("tr"),j.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                        <span>password</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Auth record password.</td>`,g=b(),S=a("tr"),S.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                        <span>passwordConfirm</span></div></td> 
                <td><span class="label">String</span></td> 
                <td>Auth record password confirmation.</td>`,N=b(),O=a("tr"),O.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                        <span>verified</span></div></td> 
                <td><span class="label">Boolean</span></td> 
                <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,C=b(),_=a("tr"),_.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>',h(f,"class","inline-flex")},m(u,T){o(u,e,T),o(u,l,T),o(u,s,T),o(u,m,T),o(u,p,T),n(p,c),n(c,f),P.m(f,null),n(f,y),n(f,$),n(p,w),n(p,M),n(p,F),n(p,D),o(u,I,T),o(u,A,T),o(u,J,T),o(u,j,T),o(u,g,T),o(u,S,T),o(u,N,T),o(u,O,T),o(u,C,T),o(u,_,T)},p(u,T){x!==(x=L(u))&&(P.d(1),P=x(u),P&&(P.c(),P.m(f,y)))},d(u){u&&r(e),u&&r(l),u&&r(s),u&&r(m),u&&r(p),P.d(),u&&r(I),u&&r(A),u&&r(J),u&&r(j),u&&r(g),u&&r(S),u&&r(N),u&&r(O),u&&r(C),u&&r(_)}}}function Pt(d){let e;return{c(){e=a("span"),e.textContent="Optional",h(e,"class","label label-warning")},m(l,s){o(l,e,s)},d(l){l&&r(e)}}}function Rt(d){let e;return{c(){e=a("span"),e.textContent="Required",h(e,"class","label label-success")},m(l,s){o(l,e,s)},d(l){l&&r(e)}}}function gt(d){let e;return{c(){e=a("span"),e.textContent="Optional",h(e,"class","label label-warning")},m(l,s){o(l,e,s)},d(l){l&&r(e)}}}function Bt(d){let e;return{c(){e=a("span"),e.textContent="Required",h(e,"class","label label-success")},m(l,s){o(l,e,s)},d(l){l&&r(e)}}}function Ft(d){var p;let e,l=((p=d[12].options)==null?void 0:p.maxSelect)===1?"id":"ids",s,m;return{c(){e=k("Relation record "),s=k(l),m=k(".")},m(c,f){o(c,e,f),o(c,s,f),o(c,m,f)},p(c,f){var y;f&1&&l!==(l=((y=c[12].options)==null?void 0:y.maxSelect)===1?"id":"ids")&&Z(s,l)},d(c){c&&r(e),c&&r(s),c&&r(m)}}}function jt(d){let e,l,s,m,p;return{c(){e=k("File object."),l=a("br"),s=k(`
                        Set to `),m=a("code"),m.textContent="null",p=k(" to delete already uploaded file(s).")},m(c,f){o(c,e,f),o(c,l,f),o(c,s,f),o(c,m,f),o(c,p,f)},p:le,d(c){c&&r(e),c&&r(l),c&&r(s),c&&r(m),c&&r(p)}}}function Dt(d){let e;return{c(){e=k("URL address.")},m(l,s){o(l,e,s)},p:le,d(l){l&&r(e)}}}function Nt(d){let e;return{c(){e=k("Email address.")},m(l,s){o(l,e,s)},p:le,d(l){l&&r(e)}}}function It(d){let e;return{c(){e=k("JSON array or object.")},m(l,s){o(l,e,s)},p:le,d(l){l&&r(e)}}}function Jt(d){let e;return{c(){e=k("Number value.")},m(l,s){o(l,e,s)},p:le,d(l){l&&r(e)}}}function Et(d){let e;return{c(){e=k("Plain text value.")},m(l,s){o(l,e,s)},p:le,d(l){l&&r(e)}}}function yt(d,e){let l,s,m,p,c,f=e[12].name+"",y,$,w,M,F=Q.getFieldValueType(e[12])+"",D,I,A,J;function j(_,L){return _[12].required?Bt:gt}let g=j(e),S=g(e);function N(_,L){if(_[12].type==="text")return Et;if(_[12].type==="number")return Jt;if(_[12].type==="json")return It;if(_[12].type==="email")return Nt;if(_[12].type==="url")return Dt;if(_[12].type==="file")return jt;if(_[12].type==="relation")return Ft}let O=N(e),C=O&&O(e);return{key:d,first:null,c(){l=a("tr"),s=a("td"),m=a("div"),S.c(),p=b(),c=a("span"),y=k(f),$=b(),w=a("td"),M=a("span"),D=k(F),I=b(),A=a("td"),C&&C.c(),J=b(),h(m,"class","inline-flex"),h(M,"class","label"),this.first=l},m(_,L){o(_,l,L),n(l,s),n(s,m),S.m(m,null),n(m,p),n(m,c),n(c,y),n(l,$),n(l,w),n(w,M),n(M,D),n(l,I),n(l,A),C&&C.m(A,null),n(l,J)},p(_,L){e=_,g!==(g=j(e))&&(S.d(1),S=g(e),S&&(S.c(),S.m(m,p))),L&1&&f!==(f=e[12].name+"")&&Z(y,f),L&1&&F!==(F=Q.getFieldValueType(e[12])+"")&&Z(D,F),O===(O=N(e))&&C?C.p(e,L):(C&&C.d(1),C=O&&O(e),C&&(C.c(),C.m(A,null)))},d(_){_&&r(l),S.d(),C&&C.d()}}}function vt(d,e){let l,s=e[7].code+"",m,p,c,f;function y(){return e[6](e[7])}return{key:d,first:null,c(){l=a("button"),m=k(s),p=b(),h(l,"class","tab-item"),ue(l,"active",e[1]===e[7].code),this.first=l},m($,w){o($,l,w),n(l,m),n(l,p),c||(f=qt(l,"click",y),c=!0)},p($,w){e=$,w&4&&s!==(s=e[7].code+"")&&Z(m,s),w&6&&ue(l,"active",e[1]===e[7].code)},d($){$&&r(l),c=!1,f()}}}function ht(d,e){let l,s,m,p;return s=new wt({props:{content:e[7].body}}),{key:d,first:null,c(){l=a("div"),Pe(s.$$.fragment),m=b(),h(l,"class","tab-item"),ue(l,"active",e[1]===e[7].code),this.first=l},m(c,f){o(c,l,f),Re(s,l,null),n(l,m),p=!0},p(c,f){e=c;const y={};f&4&&(y.content=e[7].body),s.$set(y),(!p||f&6)&&ue(l,"active",e[1]===e[7].code)},i(c){p||(fe(s.$$.fragment,c),p=!0)},o(c){pe(s.$$.fragment,c),p=!1},d(c){c&&r(l),ge(s)}}}function Ut(d){var st,it,at,ot;let e,l,s=d[0].name+"",m,p,c,f,y,$,w,M=d[0].name+"",F,D,I,A,J,j,g,S,N,O,C,_,L,x,P,u,T,ee,z=d[0].name+"",be,Be,Fe,me,ne,_e,K,ke,je,E,ye,De,ve,U=[],Ne=new Map,he,se,we,W,Ce,Ie,Se,Y,Te,Je,$e,Ee,H,Ue,te,Ve,Qe,xe,Oe,ze,Me,Ke,We,Ye,Le,Ge,He,ie,qe,G,ae,V=[],Xe=new Map,Ze,oe,B=[],et=new Map,X;S=new At({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${d[4]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},d[3],Q.dummyCollectionSchemaData(d[0])),null,4)};

const record = await pb.collection('${(st=d[0])==null?void 0:st.name}').create(data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${d[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},d[3],Q.dummyCollectionSchemaData(d[0])),null,2)};

final record = await pb.collection('${(it=d[0])==null?void 0:it.name}').create(body: body);
    `}});let R=d[5]&&_t(),q=((at=d[0])==null?void 0:at.isAuth)&&kt(d),de=(ot=d[0])==null?void 0:ot.schema;const tt=t=>t[12].name;for(let t=0;t<de.length;t+=1){let i=mt(d,de,t),v=tt(i);Ne.set(v,U[t]=yt(v,i))}te=new wt({props:{content:"?expand=relField1,relField2.subRelField"}});let ce=d[2];const lt=t=>t[7].code;for(let t=0;t<ce.length;t+=1){let i=bt(d,ce,t),v=lt(i);Xe.set(v,V[t]=vt(v,i))}let re=d[2];const nt=t=>t[7].code;for(let t=0;t<re.length;t+=1){let i=ut(d,re,t),v=nt(i);et.set(v,B[t]=ht(v,i))}return{c(){e=a("h3"),l=k("Create ("),m=k(s),p=k(")"),c=b(),f=a("div"),y=a("p"),$=k("Create a new "),w=a("strong"),F=k(M),D=k(" record."),I=b(),A=a("p"),A.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,J=b(),j=a("p"),j.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/#uploading-files" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,g=b(),Pe(S.$$.fragment),N=b(),O=a("h6"),O.textContent="API details",C=b(),_=a("div"),L=a("strong"),L.textContent="POST",x=b(),P=a("div"),u=a("p"),T=k("/api/collections/"),ee=a("strong"),be=k(z),Be=k("/records"),Fe=b(),R&&R.c(),me=b(),ne=a("div"),ne.textContent="Body Parameters",_e=b(),K=a("table"),ke=a("thead"),ke.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,je=b(),E=a("tbody"),ye=a("tr"),ye.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                    <span>id</span></div></td> 
            <td><span class="label">String</span></td> 
            <td><strong>15 characters string</strong> to store as record ID.
                <br/>
                If not set, it will be auto generated.</td>`,De=b(),q&&q.c(),ve=b();for(let t=0;t<U.length;t+=1)U[t].c();he=b(),se=a("div"),se.textContent="Query parameters",we=b(),W=a("table"),Ce=a("thead"),Ce.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Ie=b(),Se=a("tbody"),Y=a("tr"),Te=a("td"),Te.textContent="expand",Je=b(),$e=a("td"),$e.innerHTML='<span class="label">String</span>',Ee=b(),H=a("td"),Ue=k(`Auto expand relations when returning the created record. Ex.:
                `),Pe(te.$$.fragment),Ve=k(`
                Supports up to 6-levels depth nested relations expansion. `),Qe=a("br"),xe=k(`
                The expanded relations will be appended to the record under the
                `),Oe=a("code"),Oe.textContent="expand",ze=k(" property (eg. "),Me=a("code"),Me.textContent='"expand": {"relField1": {...}, ...}',Ke=k(`).
                `),We=a("br"),Ye=k(`
                Only the relations to which the account has permissions to `),Le=a("strong"),Le.textContent="view",Ge=k(" will be expanded."),He=b(),ie=a("div"),ie.textContent="Responses",qe=b(),G=a("div"),ae=a("div");for(let t=0;t<V.length;t+=1)V[t].c();Ze=b(),oe=a("div");for(let t=0;t<B.length;t+=1)B[t].c();h(e,"class","m-b-sm"),h(f,"class","content txt-lg m-b-sm"),h(O,"class","m-b-xs"),h(L,"class","label label-primary"),h(P,"class","content"),h(_,"class","alert alert-success"),h(ne,"class","section-title"),h(K,"class","table-compact table-border m-b-base"),h(se,"class","section-title"),h(W,"class","table-compact table-border m-b-base"),h(ie,"class","section-title"),h(ae,"class","tabs-header compact left"),h(oe,"class","tabs-content"),h(G,"class","tabs")},m(t,i){o(t,e,i),n(e,l),n(e,m),n(e,p),o(t,c,i),o(t,f,i),n(f,y),n(y,$),n(y,w),n(w,F),n(y,D),n(f,I),n(f,A),n(f,J),n(f,j),o(t,g,i),Re(S,t,i),o(t,N,i),o(t,O,i),o(t,C,i),o(t,_,i),n(_,L),n(_,x),n(_,P),n(P,u),n(u,T),n(u,ee),n(ee,be),n(u,Be),n(_,Fe),R&&R.m(_,null),o(t,me,i),o(t,ne,i),o(t,_e,i),o(t,K,i),n(K,ke),n(K,je),n(K,E),n(E,ye),n(E,De),q&&q.m(E,null),n(E,ve);for(let v=0;v<U.length;v+=1)U[v].m(E,null);o(t,he,i),o(t,se,i),o(t,we,i),o(t,W,i),n(W,Ce),n(W,Ie),n(W,Se),n(Se,Y),n(Y,Te),n(Y,Je),n(Y,$e),n(Y,Ee),n(Y,H),n(H,Ue),Re(te,H,null),n(H,Ve),n(H,Qe),n(H,xe),n(H,Oe),n(H,ze),n(H,Me),n(H,Ke),n(H,We),n(H,Ye),n(H,Le),n(H,Ge),o(t,He,i),o(t,ie,i),o(t,qe,i),o(t,G,i),n(G,ae);for(let v=0;v<V.length;v+=1)V[v].m(ae,null);n(G,Ze),n(G,oe);for(let v=0;v<B.length;v+=1)B[v].m(oe,null);X=!0},p(t,[i]){var rt,dt,ct,ft;(!X||i&1)&&s!==(s=t[0].name+"")&&Z(m,s),(!X||i&1)&&M!==(M=t[0].name+"")&&Z(F,M);const v={};i&25&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${t[4]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},t[3],Q.dummyCollectionSchemaData(t[0])),null,4)};

const record = await pb.collection('${(rt=t[0])==null?void 0:rt.name}').create(data);
    `),i&25&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${t[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},t[3],Q.dummyCollectionSchemaData(t[0])),null,2)};

final record = await pb.collection('${(dt=t[0])==null?void 0:dt.name}').create(body: body);
    `),S.$set(v),(!X||i&1)&&z!==(z=t[0].name+"")&&Z(be,z),t[5]?R||(R=_t(),R.c(),R.m(_,null)):R&&(R.d(1),R=null),(ct=t[0])!=null&&ct.isAuth?q?q.p(t,i):(q=kt(t),q.c(),q.m(E,ve)):q&&(q.d(1),q=null),i&1&&(de=(ft=t[0])==null?void 0:ft.schema,U=Ae(U,i,tt,1,t,de,Ne,E,pt,yt,null,mt)),i&6&&(ce=t[2],V=Ae(V,i,lt,1,t,ce,Xe,ae,pt,vt,null,bt)),i&6&&(re=t[2],$t(),B=Ae(B,i,nt,1,t,re,et,oe,Ot,ht,null,ut),Mt())},i(t){if(!X){fe(S.$$.fragment,t),fe(te.$$.fragment,t);for(let i=0;i<re.length;i+=1)fe(B[i]);X=!0}},o(t){pe(S.$$.fragment,t),pe(te.$$.fragment,t);for(let i=0;i<B.length;i+=1)pe(B[i]);X=!1},d(t){t&&r(e),t&&r(c),t&&r(f),t&&r(g),ge(S,t),t&&r(N),t&&r(O),t&&r(C),t&&r(_),R&&R.d(),t&&r(me),t&&r(ne),t&&r(_e),t&&r(K),q&&q.d();for(let i=0;i<U.length;i+=1)U[i].d();t&&r(he),t&&r(se),t&&r(we),t&&r(W),ge(te),t&&r(He),t&&r(ie),t&&r(qe),t&&r(G);for(let i=0;i<V.length;i+=1)V[i].d();for(let i=0;i<B.length;i+=1)B[i].d()}}}function Vt(d,e,l){let s,m,{collection:p=new Lt}=e,c=200,f=[],y={};const $=w=>l(1,c=w.code);return d.$$set=w=>{"collection"in w&&l(0,p=w.collection)},d.$$.update=()=>{var w,M;d.$$.dirty&1&&l(5,s=(p==null?void 0:p.createRule)===null),d.$$.dirty&1&&l(2,f=[{code:200,body:JSON.stringify(Q.dummyCollectionRecord(p),null,2)},{code:400,body:`
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
            `}]),d.$$.dirty&1&&(p.isAuth?l(3,y={username:"test_username",email:"test@exampe.com",emailVisibility:!0,password:"12345678",passwordConfirm:"12345678"}):l(3,y={}))},l(4,m=Q.getApiExampleUrl(Ht.baseUrl)),[p,c,f,y,m,s,$]}class zt extends Ct{constructor(e){super(),St(this,e,Vt,Ut,Tt,{collection:0})}}export{zt as default};
